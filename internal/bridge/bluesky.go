package bridge

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	lexutil "github.com/bluesky-social/indigo/lex/util"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/systemli/ticker/internal/bluesky"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
	"github.com/systemli/ticker/internal/util"
)

type BlueskyBridge struct {
	config  config.Config
	storage storage.Storage
}

func (bb *BlueskyBridge) Update(ticker storage.Ticker) error {
	return nil
}

func (bb *BlueskyBridge) Send(ticker storage.Ticker, message *storage.Message) error {
	if !ticker.Bluesky.Connected() || !ticker.Bluesky.Active {
		return nil
	}

	client, err := bluesky.Authenticate(ticker.Bluesky.Handle, ticker.Bluesky.AppKey)
	if err != nil {
		log.WithError(err).Error("failed to create client")
		return err
	}

	post := &bsky.FeedPost{
		Text:      message.Text,
		CreatedAt: time.Now().Local().Format(time.RFC3339),
		Facets:    []*bsky.RichtextFacet{},
	}

	links := util.ExtractURLs(message.Text)
	for _, link := range links {
		startIndex := strings.Index(message.Text, link)
		endIndex := startIndex + len(link)
		post.Facets = append(post.Facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Link: &bsky.RichtextFacet_Link{
						Uri: link,
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(startIndex),
				ByteEnd:   int64(endIndex),
			},
		})
	}

	hashtags := util.ExtractHashtags(message.Text)
	for _, hashtag := range hashtags {
		startIndex := strings.Index(message.Text, hashtag)
		endIndex := startIndex + len(hashtag)
		post.Facets = append(post.Facets, &bsky.RichtextFacet{
			Features: []*bsky.RichtextFacet_Features_Elem{
				{
					RichtextFacet_Tag: &bsky.RichtextFacet_Tag{
						Tag: hashtag[1:],
					},
				},
			},
			Index: &bsky.RichtextFacet_ByteSlice{
				ByteStart: int64(startIndex),
				ByteEnd:   int64(endIndex),
			},
		})
	}

	if len(message.Attachments) > 0 {
		var images []*bsky.EmbedImages_Image

		for _, attachment := range message.Attachments {
			upload, err := bb.storage.FindUploadByUUID(attachment.UUID)
			if err != nil {
				log.WithError(err).Error("failed to find upload")
				continue
			}

			b, err := os.ReadFile(upload.FullPath(bb.config.Upload.Path))
			if err != nil {
				log.WithError(err).Error("failed to read file")
				continue
			}

			resp, err := comatproto.RepoUploadBlob(context.TODO(), client, bytes.NewReader(b))
			if err != nil {
				log.WithError(err).Error("failed to upload blob")
				continue
			}

			images = append(images, &bsky.EmbedImages_Image{
				Image: &lexutil.LexBlob{
					Ref:      resp.Blob.Ref,
					MimeType: http.DetectContentType(b),
					Size:     resp.Blob.Size,
				},
			})
		}

		if post.Embed == nil {
			post.Embed = &bsky.FeedPost_Embed{}
		}

		post.Embed.EmbedImages = &bsky.EmbedImages{
			Images: images,
		}
	}

	resp, err := comatproto.RepoCreateRecord(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.post",
		Repo:       client.Auth.Did,
		Record: &lexutil.LexiconTypeDecoder{
			Val: post,
		},
	})
	if err != nil {
		log.WithError(err).Error("failed to create post")
		return err
	}

	// Create thread gate if reply restriction is configured
	if ticker.Bluesky.ReplyRestriction != "" {
		err = bb.createThreadGate(client, resp.Uri, ticker.Bluesky.ReplyRestriction)
		if err != nil {
			log.WithError(err).Warn("failed to create thread gate")
		}
	}

	message.Bluesky = storage.BlueskyMeta{
		Handle: ticker.Bluesky.Handle,
		Uri:    resp.Uri,
		Cid:    resp.Cid,
	}

	return nil
}

func (bb *BlueskyBridge) Delete(ticker storage.Ticker, message *storage.Message) error {
	if !ticker.Bluesky.Connected() {
		return nil
	}

	if message.Bluesky.Uri == "" {
		return nil
	}

	client, err := bluesky.Authenticate(ticker.Bluesky.Handle, ticker.Bluesky.AppKey)
	if err != nil {
		log.WithError(err).Error("failed to create client")
		return err
	}

	uri := message.Bluesky.Uri
	if !strings.HasPrefix(uri, "at://did:plc:") {
		uri = "at://did:plc:" + uri
	}

	parts := strings.Split(uri, "/")
	if len(parts) < 3 {
		log.WithError(err).WithField("uri", uri).Error("invalid post uri")
		return fmt.Errorf("invalid post uri")
	}
	rkey := parts[len(parts)-1]
	schema := parts[len(parts)-2]

	// Delete thread gate first (if it exists, ignore errors)
	_, _ = comatproto.RepoDeleteRecord(context.TODO(), client, &comatproto.RepoDeleteRecord_Input{
		Repo:       client.Auth.Did,
		Collection: "app.bsky.feed.threadgate",
		Rkey:       rkey,
	})

	_, err = comatproto.RepoDeleteRecord(context.TODO(), client, &comatproto.RepoDeleteRecord_Input{
		Repo:       client.Auth.Did,
		Collection: schema,
		Rkey:       rkey,
	})
	if err != nil {
		log.WithError(err).Error("failed to delete post")
	}

	return err
}

// createThreadGate creates a thread gate record for the given post URI.
func (bb *BlueskyBridge) createThreadGate(client *xrpc.Client, postUri string, replyRestriction string) error {
	parts := strings.Split(postUri, "/")
	if len(parts) < 3 {
		return fmt.Errorf("invalid post uri: %s", postUri)
	}
	rkey := parts[len(parts)-1]

	threadgate := &bsky.FeedThreadgate{
		Post:      postUri,
		CreatedAt: time.Now().Local().Format(time.RFC3339),
		Allow:     buildAllowRules(replyRestriction),
	}

	_, err := comatproto.RepoCreateRecord(context.TODO(), client, &comatproto.RepoCreateRecord_Input{
		Collection: "app.bsky.feed.threadgate",
		Repo:       client.Auth.Did,
		Rkey:       &rkey,
		Record: &lexutil.LexiconTypeDecoder{
			Val: threadgate,
		},
	})

	return err
}

// buildAllowRules converts a reply restriction string to thread gate allow rules.
// Valid values: "followers", "following", "mentioned", "nobody"
// Multiple values can be combined with commas, e.g. "followers,mentioned"
// "nobody" results in an empty allow list (no one can reply).
func buildAllowRules(replyRestriction string) []*bsky.FeedThreadgate_Allow_Elem {
	if replyRestriction == "nobody" {
		return []*bsky.FeedThreadgate_Allow_Elem{}
	}

	var rules []*bsky.FeedThreadgate_Allow_Elem
	for _, restriction := range strings.Split(replyRestriction, ",") {
		restriction = strings.TrimSpace(restriction)
		switch restriction {
		case "followers":
			rules = append(rules, &bsky.FeedThreadgate_Allow_Elem{
				FeedThreadgate_FollowerRule: &bsky.FeedThreadgate_FollowerRule{},
			})
		case "following":
			rules = append(rules, &bsky.FeedThreadgate_Allow_Elem{
				FeedThreadgate_FollowingRule: &bsky.FeedThreadgate_FollowingRule{},
			})
		case "mentioned":
			rules = append(rules, &bsky.FeedThreadgate_Allow_Elem{
				FeedThreadgate_MentionRule: &bsky.FeedThreadgate_MentionRule{},
			})
		}
	}

	return rules
}
