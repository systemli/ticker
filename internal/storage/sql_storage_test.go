package storage

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	pagination "github.com/systemli/ticker/internal/api/pagination"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqlStorageTestSuite struct {
	db    *gorm.DB
	store *SqlStorage
	suite.Suite
}

func (s *SqlStorageTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open("file:testdatabase?mode=memory&cache=shared"), &gorm.Config{})
	s.NoError(err)

	s.db = db
	s.store = NewSqlStorage(db, "/uploads")

	err = db.AutoMigrate(
		&Ticker{},
		&TickerTelegram{},
		&TickerMastodon{},
		&TickerBluesky{},
		&TickerSignalGroup{},
		&TickerWebsite{},
		&User{},
		&Message{},
		&Upload{},
		&Attachment{},
		&Setting{},
	)
	s.NoError(err)
}

func (s *SqlStorageTestSuite) BeforeTest(suiteName, testName string) {
	s.NoError(s.db.Exec("DELETE FROM users").Error)
	s.NoError(s.db.Exec("DELETE FROM messages").Error)
	s.NoError(s.db.Exec("DELETE FROM attachments").Error)
	s.NoError(s.db.Exec("DELETE FROM tickers").Error)
	s.NoError(s.db.Exec("DELETE FROM ticker_mastodons").Error)
	s.NoError(s.db.Exec("DELETE FROM ticker_telegrams").Error)
	s.NoError(s.db.Exec("DELETE FROM ticker_blueskies").Error)
	s.NoError(s.db.Exec("DELETE FROM ticker_signal_groups").Error)
	s.NoError(s.db.Exec("DELETE FROM ticker_websites").Error)
	s.NoError(s.db.Exec("DELETE FROM settings").Error)
	s.NoError(s.db.Exec("DELETE FROM uploads").Error)
}

func (s *SqlStorageTestSuite) TestFindUsers() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when no users exist", func() {
		filter := NewUserFilter(nil)
		users, err := s.store.FindUsers(filter)
		s.NoError(err)
		s.Empty(users)
	})

	s.Run("when users exist", func() {
		user.Tickers = []Ticker{{ID: 1}}
		err := s.db.Create(&user).Error
		s.NoError(err)

		s.Run("without preload", func() {
			filter := NewUserFilter(nil)
			users, err := s.store.FindUsers(filter)
			s.NoError(err)
			s.Len(users, 1)
			s.Empty(users[0].Tickers)
		})

		s.Run("with preload", func() {
			filter := NewUserFilter(nil)
			users, err := s.store.FindUsers(filter, WithTickers())
			s.NoError(err)
			s.Len(users, 1)
			s.Len(users[0].Tickers, 1)
		})

		s.Run("with filters", func() {
			email := "user@example.org"
			isSuperAdmin := false
			filter := UserFilter{Email: &email, IsSuperAdmin: &isSuperAdmin, OrderBy: "id", Sort: "asc"}
			users, err := s.store.FindUsers(filter)
			s.NoError(err)
			s.Len(users, 1)

			email = "user@example.com"
			filter = UserFilter{Email: &email, OrderBy: "id", Sort: "asc"}
			users, err = s.store.FindUsers(filter)
			s.NoError(err)
			s.Empty(users)

			isSuperAdmin = true
			filter = UserFilter{IsSuperAdmin: &isSuperAdmin, OrderBy: "id", Sort: "asc"}
			users, err = s.store.FindUsers(filter)
			s.NoError(err)
			s.Empty(users)
		})
	})

}

func (s *SqlStorageTestSuite) TestFindUserByID() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when user does not exist", func() {
		_, err := s.store.FindUserByID(1)
		s.Error(err)
	})

	s.Run("when user exists", func() {
		user.Tickers = []Ticker{{ID: 1}}
		err := s.db.Create(&user).Error
		s.NoError(err)

		s.Run("without preload", func() {
			user, err := s.store.FindUserByID(user.ID)
			s.NoError(err)
			s.NotNil(user)
			s.Empty(user.Tickers)
		})

		s.Run("with preload", func() {
			user, err := s.store.FindUserByID(user.ID, WithTickers())
			s.NoError(err)
			s.NotNil(user)
			s.Len(user.Tickers, 1)
		})
	})
}

func (s *SqlStorageTestSuite) TestFindUsersByIDs() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when no users exist", func() {
		users, err := s.store.FindUsersByIDs([]int{1, 2})
		s.NoError(err)
		s.Empty(users)
	})

	s.Run("when empty ids", func() {
		users, err := s.store.FindUsersByIDs([]int{})
		s.NoError(err)
		s.Empty(users)
	})

	s.Run("when users exist", func() {
		user.Tickers = []Ticker{{ID: 1}}
		err := s.db.Create(&user).Error
		s.NoError(err)

		s.Run("without preload", func() {
			users, err := s.store.FindUsersByIDs([]int{user.ID})
			s.NoError(err)
			s.Len(users, 1)
			s.Empty(users[0].Tickers)
		})

		s.Run("with preload", func() {
			users, err := s.store.FindUsersByIDs([]int{user.ID}, WithTickers())
			s.NoError(err)
			s.Len(users, 1)
			s.Len(users[0].Tickers, 1)
		})
	})
}

func (s *SqlStorageTestSuite) TestFindUserByEmail() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when user does not exist", func() {
		_, err := s.store.FindUserByEmail("user@example.org")
		s.Error(err)
	})

	s.Run("when user exists", func() {
		user.Tickers = []Ticker{{ID: 1}}
		err := s.db.Create(&user).Error
		s.NoError(err)

		s.Run("without preload", func() {
			user, err := s.store.FindUserByEmail("user@example.org")
			s.NoError(err)
			s.NotNil(user)
			s.Empty(user.Tickers)
		})

		s.Run("with preload", func() {
			user.Tickers = []Ticker{{ID: 1}}
			user, err := s.store.FindUserByEmail("user@example.org", WithTickers())
			s.NoError(err)
			s.NotNil(user)
			s.Len(user.Tickers, 1)
		})
	})
}

func (s *SqlStorageTestSuite) TestFindUsersByTicker() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when no users exist", func() {
		users, err := s.store.FindUsersByTicker(Ticker{ID: 1})
		s.NoError(err)
		s.Empty(users)
	})

	s.Run("when users exist", func() {
		user.Tickers = []Ticker{{ID: 1}}
		err := s.db.Create(&user).Error
		s.NoError(err)

		s.Run("without preload", func() {
			users, err := s.store.FindUsersByTicker(Ticker{ID: 1})
			s.NoError(err)
			s.Len(users, 1)
			s.Empty(users[0].Tickers)
		})

		s.Run("with preload", func() {
			users, err := s.store.FindUsersByTicker(Ticker{ID: 1}, WithTickers())
			s.NoError(err)
			s.Len(users, 1)
			s.Len(users[0].Tickers, 1)
		})
	})
}

func (s *SqlStorageTestSuite) TestSaveUser() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	s.Run("when user is new", func() {
		err = s.store.SaveUser(&user)
		s.NoError(err)

		var count int64
		err = s.db.Model(&User{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when email is empty", func() {
		user.Email = ""
		err = s.store.SaveUser(&user)
		s.Error(err)
	})

	s.Run("when encrypted password is empty", func() {
		user := User{Email: "user@example.org"}
		user.EncryptedPassword = ""
		err = s.store.SaveUser(&user)
		s.Error(err)
	})

	s.Run("when user is existing", func() {
		user.Email = "update@example.org"

		err = s.store.SaveUser(&user)
		s.NoError(err)
	})

	s.Run("when user is existing with tickers", func() {
		ticker := Ticker{}
		err = s.store.SaveTicker(&ticker)
		s.NoError(err)

		user.Tickers = append(user.Tickers, ticker)
		err = s.store.SaveUser(&user)
		s.NoError(err)

		user, err = s.store.FindUserByID(user.ID, WithTickers())
		s.NoError(err)
		s.Len(user.Tickers, 1)
	})

	s.Run("when user removes tickers", func() {
		user.Tickers = []Ticker{}
		err = s.store.SaveUser(&user)
		s.NoError(err)

		user, err = s.store.FindUserByID(user.ID, WithTickers())
		s.NoError(err)
		s.Empty(user.Tickers)
	})
}

func (s *SqlStorageTestSuite) TestDeleteUser() {
	s.Run("when user does not exist", func() {
		user := User{ID: 1}
		err := s.store.DeleteUser(user)
		s.NoError(err)
	})

	s.Run("when user exists", func() {
		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		err = s.store.DeleteUser(user)
		s.NoError(err)

		var count int64
		err = s.db.Model(&User{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteTickerUsers() {
	s.Run("when ticker does not exist", func() {
		ticker := &Ticker{ID: 1}
		err := s.store.DeleteTickerUsers(ticker)
		s.NoError(err)
	})

	s.Run("when ticker exists", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		count := s.db.Model(&ticker).Association("Users").Count()
		s.Equal(int64(0), count)

		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		err = s.db.Model(&ticker).Association("Users").Append(&user)
		s.NoError(err)

		count = s.db.Model(&ticker).Association("Users").Count()
		s.Equal(int64(1), count)

		err = s.store.DeleteTickerUsers(ticker)
		s.NoError(err)

		count = s.db.Model(&ticker).Association("Users").Count()
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteTickerUser() {
	s.Run("when ticker does not exist", func() {
		ticker := &Ticker{ID: 1}
		user := &User{ID: 1}
		err := s.store.DeleteTickerUser(ticker, user)
		s.NoError(err)
	})

	s.Run("when ticker exists", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		err = s.db.Model(&ticker).Association("Users").Append(&user)
		s.NoError(err)

		count := s.db.Model(&ticker).Association("Users").Count()
		s.Equal(int64(1), count)

		err = s.store.DeleteTickerUser(ticker, &user)
		s.NoError(err)

		count = s.db.Model(&ticker).Association("Users").Count()
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestAddTickerUser() {
	ticker := &Ticker{}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	err = s.db.Create(&user).Error
	s.NoError(err)

	err = s.store.AddTickerUser(ticker, &user)
	s.NoError(err)

	count := s.db.Model(&ticker).Association("Users").Count()
	s.Equal(int64(1), count)
}

func (s *SqlStorageTestSuite) TestFindTickerByID() {
	s.Run("when ticker does not exist", func() {
		_, err := s.store.FindTickerByID(1)
		s.Error(err)
	})

	err := s.db.Create(&Ticker{ID: 1}).Error
	s.NoError(err)

	s.Run("when ticker exists", func() {
		ticker, err := s.store.FindTickerByID(1)
		s.NoError(err)
		s.NotNil(ticker)
	})

	s.Run("when ticker exists with users", func() {
		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		err = s.db.Model(&Ticker{ID: 1}).Association("Users").Append(&user)
		s.NoError(err)

		ticker, err := s.store.FindTickerByID(1, WithPreload())
		s.NoError(err)
		s.NotNil(ticker)
		s.Len(ticker.Users, 1)
	})
}

func (s *SqlStorageTestSuite) TestFindTickersByIDs() {
	s.Run("when no tickers exist", func() {
		tickers, err := s.store.FindTickersByIDs([]int{1, 2})
		s.NoError(err)
		s.Empty(tickers)
	})
	err := s.db.Create(&Ticker{ID: 1}).Error
	s.NoError(err)

	s.Run("when tickers exist", func() {
		tickers, err := s.store.FindTickersByIDs([]int{1})
		s.NoError(err)
		s.Len(tickers, 1)
	})

	s.Run("when tickers exist with users", func() {
		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		err = s.db.Model(&Ticker{ID: 1}).Association("Users").Append(&user)
		s.NoError(err)

		tickers, err := s.store.FindTickersByIDs([]int{1}, WithPreload())
		s.NoError(err)
		s.Len(tickers, 1)
		s.Len(tickers[0].Users, 1)
	})
}

func (s *SqlStorageTestSuite) TestFindTickerByOrigin() {
	s.Run("when ticker does not exist", func() {
		_, err := s.store.FindTickerByOrigin("https://systemli.org")
		s.Error(err)
	})

	ticker := Ticker{Websites: []TickerWebsite{{Origin: "https://systemli.org"}}}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when ticker exists", func() {
		ticker, err := s.store.FindTickerByOrigin("https://systemli.org")
		s.NoError(err)
		s.NotNil(ticker)
	})

	s.Run("when ticker exists with preload", func() {
		ticker.Mastodon = TickerMastodon{Active: true}
		ticker.Telegram = TickerTelegram{Active: true}

		err = s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&ticker).Error
		s.NoError(err)

		ticker, err := s.store.FindTickerByOrigin("https://systemli.org", WithPreload())
		s.NoError(err)
		s.NotNil(ticker)
		s.True(ticker.Mastodon.Active)
		s.True(ticker.Telegram.Active)
		s.Equal("https://systemli.org", ticker.Websites[0].Origin)
	})
}

func (s *SqlStorageTestSuite) TestFindTickersByUser() {
	s.Run("when no tickers exist", func() {
		filter := TickerFilter{OrderBy: "id", Sort: "desc"}
		tickers, err := s.store.FindTickersByUser(User{ID: 1}, filter)
		s.NoError(err)
		s.Empty(tickers)
	})

	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	err = s.db.Create(&user).Error
	s.NoError(err)

	ticker := Ticker{Users: []User{user}, Active: false, Title: "title", Websites: []TickerWebsite{{Origin: "http://localhost"}}}
	err = s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when tickers exist", func() {
		filter := TickerFilter{OrderBy: "id", Sort: "desc"}
		tickers, err := s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Len(tickers, 1)
	})

	s.Run("when tickers exist with preload", func() {
		filter := TickerFilter{OrderBy: "id", Sort: "desc"}
		tickers, err := s.store.FindTickersByUser(user, filter, WithPreload())
		s.NoError(err)
		s.Len(tickers, 1)
		s.Len(tickers[0].Users, 1)
	})

	s.Run("when super admin", func() {
		filter := TickerFilter{OrderBy: "id", Sort: "desc"}
		tickers, err := s.store.FindTickersByUser(User{IsSuperAdmin: true}, filter)
		s.NoError(err)
		s.Len(tickers, 1)
	})

	s.Run("when filter is set", func() {
		active := true
		filter := TickerFilter{OrderBy: "id", Sort: "desc", Active: &active}
		tickers, err := s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Empty(tickers)

		active = false
		filter = TickerFilter{OrderBy: "id", Sort: "desc", Active: &active}
		tickers, err = s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Len(tickers, 1)

		title := "title"
		filter = TickerFilter{OrderBy: "id", Sort: "desc", Title: &title}
		tickers, err = s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Len(tickers, 1)

		origin := "localhost"
		filter = TickerFilter{OrderBy: "id", Sort: "desc", Origin: &origin}
		tickers, err = s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Len(tickers, 1)

		origin = "systemli.org"
		filter = TickerFilter{OrderBy: "id", Sort: "desc", Origin: &origin}
		tickers, err = s.store.FindTickersByUser(user, filter)
		s.NoError(err)
		s.Empty(tickers)
	})
}

func (s *SqlStorageTestSuite) TestFindTickerByUserAndID() {
	user, err := NewUser("user@example.org", "password")
	s.NoError(err)

	err = s.db.Create(&user).Error
	s.NoError(err)

	ticker := Ticker{Users: []User{user}}
	err = s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when ticker exists", func() {
		ticker, err := s.store.FindTickerByUserAndID(user, ticker.ID)
		s.NoError(err)
		s.NotNil(ticker)
	})

	s.Run("when ticker exists with preload", func() {
		ticker, err := s.store.FindTickerByUserAndID(user, ticker.ID, WithPreload())
		s.NoError(err)
		s.NotNil(ticker)
		s.Len(ticker.Users, 1)
	})

	s.Run("when super admin", func() {
		ticker, err := s.store.FindTickerByUserAndID(User{IsSuperAdmin: true}, ticker.ID)
		s.NoError(err)
		s.NotNil(ticker)
	})
}

func (s *SqlStorageTestSuite) TestSaveTicker() {
	ticker := Ticker{}

	s.Run("when ticker is new", func() {
		err := s.store.SaveTicker(&ticker)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Ticker{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when ticker is existing", func() {
		ticker.Domain = "systemli.org"
		err := s.store.SaveTicker(&ticker)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Ticker{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when ticker is existing and properties are updated", func() {
		ticker.Active = true
		ticker.Information = TickerInformation{
			Author:   "Author",
			Email:    "Email",
			Twitter:  "Twitter",
			Facebook: "Facebook",
			Telegram: "Telegram",
			Mastodon: "Mastodon",
			Bluesky:  "Bluesky",
		}
		ticker.Location = TickerLocation{Lat: 1, Lon: 1}

		s.NoError(s.store.SaveTicker(&ticker))
		s.True(ticker.Active)
		s.Equal("Author", ticker.Information.Author)
		s.Equal(float64(1), ticker.Location.Lat)
		s.Equal(float64(1), ticker.Location.Lon)

		ticker.Active = false
		ticker.Information.Author = ""
		ticker.Information.Email = ""
		ticker.Information.Twitter = ""
		ticker.Information.Facebook = ""
		ticker.Information.Telegram = ""
		ticker.Information.Mastodon = ""
		ticker.Information.Bluesky = ""
		ticker.Location.Lat = 0
		ticker.Location.Lon = 0

		s.NoError(s.store.SaveTicker(&ticker))

		ticker, err := s.store.FindTickerByID(ticker.ID)
		s.NoError(err)
		s.Equal("", ticker.Information.Author)
		s.Equal("", ticker.Information.Email)
		s.Equal("", ticker.Information.Twitter)
		s.Equal("", ticker.Information.Facebook)
		s.Equal("", ticker.Information.Telegram)
		s.Equal("", ticker.Information.Mastodon)
		s.Equal("", ticker.Information.Bluesky)
		s.Equal(float64(0), ticker.Location.Lat)
		s.Equal(float64(0), ticker.Location.Lon)
	})

	s.Run("when ticker is existing with users", func() {
		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		ticker.Users = append(ticker.Users, user)
		err = s.store.SaveTicker(&ticker)
		s.NoError(err)

		ticker, err = s.store.FindTickerByID(ticker.ID, WithPreload())
		s.NoError(err)
		s.Len(ticker.Users, 1)
	})

	s.Run("when ticker removes users", func() {
		ticker.Users = []User{}
		err := s.store.SaveTicker(&ticker)
		s.NoError(err)

		ticker, err = s.store.FindTickerByID(ticker.ID, WithPreload())
		s.NoError(err)
		s.Empty(ticker.Users)
	})
}

func (s *SqlStorageTestSuite) TestDeleteIntegrations() {
	s.Run("when ticker is resetting integrations", func() {
		var telegramCount, mastodonCount, blueskyCount, signalGroupCount int64

		ticker := Ticker{}
		ticker.Telegram = TickerTelegram{Active: true, ChannelName: "channel"}
		ticker.Mastodon = TickerMastodon{Active: true, Server: "server", Token: "token", AccessToken: "access_token"}
		ticker.Bluesky = TickerBluesky{Active: true, AppKey: "app_key"}
		ticker.SignalGroup = TickerSignalGroup{Active: true, GroupID: "group_id", GroupInviteLink: "group_invite_link"}

		err := s.store.SaveTicker(&ticker)
		s.NoError(err)

		s.store.DB.Model(&TickerTelegram{}).Where("ticker_id = ?", ticker.ID).Count(&telegramCount)
		s.store.DB.Model(&TickerMastodon{}).Where("ticker_id = ?", ticker.ID).Count(&mastodonCount)
		s.store.DB.Model(&TickerBluesky{}).Where("ticker_id = ?", ticker.ID).Count(&blueskyCount)
		s.store.DB.Model(&TickerSignalGroup{}).Where("ticker_id = ?", ticker.ID).Count(&signalGroupCount)

		s.Equal(int64(1), telegramCount)
		s.Equal(int64(1), mastodonCount)
		s.Equal(int64(1), blueskyCount)
		s.Equal(int64(1), signalGroupCount)

		s.True(ticker.Telegram.Active)
		s.Equal("channel", ticker.Telegram.ChannelName)
		s.True(ticker.Mastodon.Active)
		s.Equal("server", ticker.Mastodon.Server)
		s.Equal("token", ticker.Mastodon.Token)
		s.Equal("access_token", ticker.Mastodon.AccessToken)
		s.True(ticker.Bluesky.Active)
		s.Equal("app_key", ticker.Bluesky.AppKey)
		s.True(ticker.SignalGroup.Active)
		s.Equal("group_id", ticker.SignalGroup.GroupID)
		s.Equal("group_invite_link", ticker.SignalGroup.GroupInviteLink)

		err = s.store.DeleteIntegrations(&ticker)
		s.NoError(err)

		s.False(ticker.Telegram.Active)
		s.Empty(ticker.Telegram.ChannelName)
		s.False(ticker.Mastodon.Active)
		s.Empty(ticker.Mastodon.Server)
		s.Empty(ticker.Mastodon.Token)
		s.Empty(ticker.Mastodon.AccessToken)
		s.False(ticker.Bluesky.Active)
		s.Empty(ticker.Bluesky.AppKey)
		s.False(ticker.SignalGroup.Active)
		s.Empty(ticker.SignalGroup.GroupID)
		s.Empty(ticker.SignalGroup.GroupInviteLink)

		s.store.DB.Model(&TickerTelegram{}).Where("ticker_id = ?", ticker.ID).Count(&telegramCount)
		s.store.DB.Model(&TickerMastodon{}).Where("ticker_id = ?", ticker.ID).Count(&mastodonCount)
		s.store.DB.Model(&TickerBluesky{}).Where("ticker_id = ?", ticker.ID).Count(&blueskyCount)
		s.store.DB.Model(&TickerSignalGroup{}).Where("ticker_id = ?", ticker.ID).Count(&signalGroupCount)
		s.Equal(int64(0), telegramCount)
		s.Equal(int64(0), mastodonCount)
		s.Equal(int64(0), blueskyCount)
		s.Equal(int64(0), signalGroupCount)
	})
}

func (s *SqlStorageTestSuite) TestResetTicker() {
	ticker := &Ticker{
		Title:       "title",
		Domain:      "example.org",
		Active:      true,
		Description: "description",
		Information: TickerInformation{
			Author:   "author",
			Email:    "email",
			Twitter:  "twitter",
			Facebook: "facebook",
			Telegram: "telegram",
			Mastodon: "mastodon",
			Bluesky:  "bluesky",
		},
		Location: TickerLocation{
			Lat: 1,
			Lon: 1,
		},
	}

	err := s.store.SaveTicker(ticker)
	s.NoError(err)
	s.NotZero(ticker.ID)

	s.Run("basic reset", func() {
		err = s.store.ResetTicker(ticker)

		s.NoError(err)
		s.Equal("Ticker", ticker.Title)
		s.Equal("example.org", ticker.Domain)
		s.False(ticker.Active)
		s.Empty(ticker.Description)
		s.Empty(ticker.Information.Author)
		s.Empty(ticker.Information.Email)
		s.Empty(ticker.Information.Twitter)
		s.Empty(ticker.Information.Facebook)
		s.Empty(ticker.Information.Telegram)
		s.Empty(ticker.Information.Mastodon)
		s.Empty(ticker.Information.Bluesky)
		s.Zero(ticker.Location.Lat)
		s.Zero(ticker.Location.Lon)
	})

	s.Run("with integrations", func() {
		var telegramCount, mastodonCount, blueskyCount, signalGroupCount int64
		ticker.Telegram = TickerTelegram{Active: true, ChannelName: "channel"}
		ticker.Mastodon = TickerMastodon{Active: true, Server: "server", Token: "token", AccessToken: "access_token"}
		ticker.Bluesky = TickerBluesky{Active: true, AppKey: "app_key"}
		ticker.SignalGroup = TickerSignalGroup{Active: true, GroupID: "group_id", GroupInviteLink: "group_invite_link"}

		err = s.store.SaveTicker(ticker)
		s.NoError(err)

		s.store.DB.Model(&TickerTelegram{}).Where("ticker_id = ?", ticker.ID).Count(&telegramCount)
		s.store.DB.Model(&TickerMastodon{}).Where("ticker_id = ?", ticker.ID).Count(&mastodonCount)
		s.store.DB.Model(&TickerBluesky{}).Where("ticker_id = ?", ticker.ID).Count(&blueskyCount)
		s.store.DB.Model(&TickerSignalGroup{}).Where("ticker_id = ?", ticker.ID).Count(&signalGroupCount)

		s.Equal(int64(1), telegramCount)
		s.Equal(int64(1), mastodonCount)
		s.Equal(int64(1), blueskyCount)
		s.Equal(int64(1), signalGroupCount)

		err = s.store.ResetTicker(ticker)
		s.NoError(err)
		s.False(ticker.Telegram.Active)
		s.False(ticker.Mastodon.Active)
		s.False(ticker.Bluesky.Active)
		s.False(ticker.SignalGroup.Active)
		s.Empty(ticker.Telegram.ChannelName)
		s.Empty(ticker.Mastodon.Server)
		s.Empty(ticker.Mastodon.Token)
		s.Empty(ticker.Mastodon.AccessToken)
		s.Empty(ticker.Bluesky.AppKey)
		s.Empty(ticker.SignalGroup.GroupID)
		s.Empty(ticker.SignalGroup.GroupInviteLink)

		s.store.DB.Model(&TickerTelegram{}).Where("ticker_id = ?", ticker.ID).Count(&telegramCount)
		s.store.DB.Model(&TickerMastodon{}).Where("ticker_id = ?", ticker.ID).Count(&mastodonCount)
		s.store.DB.Model(&TickerBluesky{}).Where("ticker_id = ?", ticker.ID).Count(&blueskyCount)
		s.store.DB.Model(&TickerSignalGroup{}).Where("ticker_id = ?", ticker.ID).Count(&signalGroupCount)

		s.Equal(int64(0), telegramCount)
		s.Equal(int64(0), mastodonCount)
		s.Equal(int64(0), blueskyCount)
		s.Equal(int64(0), signalGroupCount)
	})
}

func (s *SqlStorageTestSuite) TestDeleteTicker() {
	s.Run("when ticker does not exist", func() {
		ticker := &Ticker{ID: 1}
		err := s.store.DeleteTicker(ticker)
		s.NoError(err)
	})

	s.Run("when ticker exists", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(ticker).Error
		s.NoError(err)

		err = s.store.DeleteTicker(ticker)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Ticker{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})

	s.Run("when ticker exists with user", func() {
		user, err := NewUser("user@example.org", "password")
		s.NoError(err)

		err = s.db.Create(&user).Error
		s.NoError(err)

		ticker := &Ticker{ID: 1, Users: []User{user}}
		err = s.db.Create(&ticker).Error
		s.NoError(err)

		err = s.store.DeleteTicker(ticker)
		s.NoError(err)

		var tickerCount int64
		err = s.db.Model(&Ticker{}).Count(&tickerCount).Error
		s.NoError(err)
		s.Equal(int64(0), tickerCount)

		var userCount int64
		err = s.db.Model(&User{}).Count(&userCount).Error
		s.NoError(err)
		s.Equal(int64(1), userCount)
	})

	s.Run("when ticker exists with uploads", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		upload := Upload{TickerID: ticker.ID}
		err = s.db.Create(&upload).Error
		s.NoError(err)

		err = s.store.DeleteTicker(ticker)
		s.NoError(err)

		var tickerCount int64
		err = s.db.Model(&Ticker{}).Count(&tickerCount).Error
		s.NoError(err)
		s.Equal(int64(0), tickerCount)

		var uploadCount int64
		err = s.db.Model(&Upload{}).Count(&uploadCount).Error
		s.NoError(err)
		s.Equal(int64(0), uploadCount)
	})

	s.Run("when ticker exists with messages", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		message := Message{TickerID: ticker.ID}
		err = s.db.Create(&message).Error
		s.NoError(err)

		err = s.store.DeleteTicker(ticker)
		s.NoError(err)

		var tickerCount int64
		err = s.db.Model(&Ticker{}).Count(&tickerCount).Error
		s.NoError(err)
		s.Equal(int64(0), tickerCount)

		var messageCount int64
		err = s.db.Model(&Message{}).Count(&messageCount).Error
		s.NoError(err)
		s.Equal(int64(0), messageCount)
	})

	s.Run("when ticker exists with integrations", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		mastodon := TickerMastodon{TickerID: ticker.ID}
		err = s.db.Create(&mastodon).Error
		s.NoError(err)

		telegram := TickerTelegram{TickerID: ticker.ID}
		err = s.db.Create(&telegram).Error
		s.NoError(err)

		bluesky := TickerBluesky{TickerID: ticker.ID}
		err = s.db.Create(&bluesky).Error
		s.NoError(err)

		signalGroup := TickerSignalGroup{TickerID: ticker.ID}
		err = s.db.Create(&signalGroup).Error
		s.NoError(err)

		err = s.store.DeleteTicker(ticker)
		s.NoError(err)

		var tickerCount int64
		err = s.db.Model(&Ticker{}).Count(&tickerCount).Error
		s.NoError(err)
		s.Equal(int64(0), tickerCount)

		var mastodonCount int64
		err = s.db.Model(&TickerMastodon{}).Count(&mastodonCount).Error
		s.NoError(err)
		s.Equal(int64(0), mastodonCount)

		var telegramCount int64
		err = s.db.Model(&TickerTelegram{}).Count(&telegramCount).Error
		s.NoError(err)
		s.Equal(int64(0), telegramCount)

		var blueskyCount int64
		err = s.db.Model(&TickerBluesky{}).Count(&blueskyCount).Error
		s.NoError(err)
		s.Equal(int64(0), blueskyCount)

		var signalGroupCount int64
		err = s.db.Model(&TickerSignalGroup{}).Count(&signalGroupCount).Error
		s.NoError(err)
		s.Equal(int64(0), signalGroupCount)
	})
}

func (s *SqlStorageTestSuite) TestSaveTickerWebsites() {
	ticker := Ticker{}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when ticker website is new", func() {
		err = s.store.SaveTickerWebsites(&ticker, []TickerWebsite{{Origin: "http://localhost"}})
		s.NoError(err)

		var count int64
		err = s.db.Model(&TickerWebsite{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when ticker website is existing", func() {
		ticker.Websites = []TickerWebsite{{Origin: "http://localhost:3000"}}
		err := s.db.Save(&ticker).Error
		s.NoError(err)

		err = s.store.SaveTickerWebsites(&ticker, ticker.Websites)
		s.NoError(err)
	})

	s.Run("when ticker website is removed", func() {
		err = s.store.SaveTickerWebsites(&ticker, []TickerWebsite{})
		s.NoError(err)

		var count int64
		err = s.db.Model(&TickerWebsite{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestFindUploadByUUID() {
	s.Run("when upload does not exist", func() {
		_, err := s.store.FindUploadByUUID("uuid")
		s.Error(err)
	})

	s.Run("when upload exists", func() {
		err := s.db.Create(&Upload{UUID: "uuid"}).Error
		s.NoError(err)

		upload, err := s.store.FindUploadByUUID("uuid")
		s.NoError(err)
		s.NotNil(upload)
	})
}

func (s *SqlStorageTestSuite) TestFindUploadsByIDs() {
	s.Run("when no uploads exist", func() {
		uploads, err := s.store.FindUploadsByIDs([]int{1, 2})
		s.NoError(err)
		s.Empty(uploads)
	})

	s.Run("when uploads exist", func() {
		err := s.db.Create(&Upload{ID: 1}).Error
		s.NoError(err)

		uploads, err := s.store.FindUploadsByIDs([]int{1})
		s.NoError(err)
		s.Len(uploads, 1)
	})
}

func (s *SqlStorageTestSuite) TestSaveUpload() {
	upload := Upload{}

	s.Run("when upload is new", func() {
		err := s.store.SaveUpload(&upload)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when upload is existing", func() {
		upload.UUID = "uuid"
		err := s.store.SaveUpload(&upload)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
		s.Equal("uuid", upload.UUID)
	})
}

func (s *SqlStorageTestSuite) TestDeleteUpload() {
	s.Run("when upload does not exist", func() {
		upload := Upload{ID: 1}
		err := s.store.DeleteUpload(upload)
		s.NoError(err)
	})

	s.Run("when upload exists", func() {
		upload := Upload{ID: 1}
		err := s.db.Create(&upload).Error
		s.NoError(err)

		err = s.store.DeleteUpload(upload)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteUploads() {
	s.Run("when uploads do not exist", func() {
		uploads := []Upload{{ID: 1}}
		s.store.DeleteUploads(uploads)

		var count int64
		err := s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})

	s.Run("when uploads exist", func() {
		var count int64
		uploads := []Upload{{ID: 1}}
		err := s.db.Create(&uploads).Error
		s.NoError(err)

		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)

		s.store.DeleteUploads(uploads)
		s.NoError(err)

		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteUploadsByTicker() {
	s.Run("when uploads do not exist", func() {
		ticker := &Ticker{ID: 1}
		err := s.store.DeleteUploadsByTicker(ticker)
		s.NoError(err)
	})

	s.Run("when uploads exist", func() {
		ticker := &Ticker{ID: 1}
		err := s.db.Create(&ticker).Error
		s.NoError(err)

		upload := Upload{TickerID: ticker.ID}
		err = s.db.Create(&upload).Error
		s.NoError(err)

		var count int64
		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)

		err = s.store.DeleteUploadsByTicker(ticker)
		s.NoError(err)

		err = s.db.Model(&Upload{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestFindMessage() {
	s.Run("when message does not exist", func() {
		_, err := s.store.FindMessage(1, 1)
		s.Error(err)
	})

	message := Message{ID: 1, TickerID: 1}
	err := s.db.Create(&message).Error
	s.NoError(err)

	s.Run("when message exists", func() {
		message, err := s.store.FindMessage(1, 1)
		s.NoError(err)
		s.NotNil(message)
	})

	s.Run("when message exists with attachments", func() {
		message.Attachments = []Attachment{{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}
		err := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&message).Error
		s.NoError(err)

		message, err := s.store.FindMessage(1, 1, WithAttachments())
		s.NoError(err)
		s.NotNil(message)
		s.Len(message.Attachments, 1)
		s.Equal("uuid", message.Attachments[0].UUID)
	})
}

func (s *SqlStorageTestSuite) TestFindMessagesByTicker() {
	ticker := Ticker{ID: 1}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when no messages exist", func() {
		messages, err := s.store.FindMessagesByTicker(ticker)
		s.NoError(err)
		s.Empty(messages)
	})

	s.Run("when messages exist", func() {
		messages, err := s.store.FindMessagesByTicker(ticker)
		s.NoError(err)
		s.Empty(messages)
	})

	s.Run("when messages exist with attachments", func() {
		message := Message{
			TickerID: ticker.ID,
			Text:     "Text",
			Attachments: []Attachment{
				{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"},
			},
		}
		err := s.db.Create(&message).Error
		s.NoError(err)

		messages, err := s.store.FindMessagesByTicker(ticker, WithAttachments())
		s.NoError(err)
		s.Len(messages, 1)

		s.Len(messages[0].Attachments, 1)
		s.Equal("uuid", messages[0].Attachments[0].UUID)
	})
}

func (s *SqlStorageTestSuite) TestFindMessagesByTickerAndPagination() {
	ticker := Ticker{ID: 1}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	s.Run("when no messages exist", func() {
		p := pagination.NewPagination(&gin.Context{})
		messages, err := s.store.FindMessagesByTickerAndPagination(ticker, *p)
		s.NoError(err)
		s.Empty(messages)
	})

	err = s.db.Create(&[]Message{
		{TickerID: ticker.ID, ID: 1, Attachments: []Attachment{{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}},
		{TickerID: ticker.ID, ID: 2, Attachments: []Attachment{{ID: 2, MessageID: 2, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}},
		{TickerID: ticker.ID, ID: 3, Attachments: []Attachment{{ID: 3, MessageID: 3, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}},
		{TickerID: ticker.ID, ID: 4, Attachments: []Attachment{{ID: 4, MessageID: 4, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}},
	}).Error
	s.NoError(err)

	s.Run("when messages exist with attachments", func() {
		p := pagination.NewPagination(&gin.Context{})
		messages, err := s.store.FindMessagesByTickerAndPagination(ticker, *p, WithAttachments())
		s.NoError(err)
		s.Len(messages, 4)
		s.Equal("uuid", messages[0].Attachments[0].UUID)
		s.Equal("uuid", messages[1].Attachments[0].UUID)
		s.Equal("uuid", messages[2].Attachments[0].UUID)
		s.Equal("uuid", messages[3].Attachments[0].UUID)
	})

	s.Run("when messages exist with limit set", func() {
		c := &gin.Context{}
		c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2"}}
		p := pagination.NewPagination(c)
		messages, err := s.store.FindMessagesByTickerAndPagination(ticker, *p)
		s.NoError(err)
		s.Len(messages, 2)
		s.Equal(4, messages[0].ID)
		s.Equal(3, messages[1].ID)
	})

	s.Run("when messages exist with limit and after set", func() {
		c := &gin.Context{}
		c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2&after=2"}}
		p := pagination.NewPagination(c)
		messages, err := s.store.FindMessagesByTickerAndPagination(ticker, *p)
		s.NoError(err)
		s.Len(messages, 2)
		s.Equal(4, messages[0].ID)
		s.Equal(3, messages[1].ID)
	})

	s.Run("when messages exist with limit and before set", func() {
		c := &gin.Context{}
		c.Request = &http.Request{URL: &url.URL{RawQuery: "limit=2&before=4"}}
		p := pagination.NewPagination(c)
		messages, err := s.store.FindMessagesByTickerAndPagination(ticker, *p)
		s.NoError(err)
		s.Len(messages, 2)
		s.Equal(3, messages[0].ID)
		s.Equal(2, messages[1].ID)
	})
}

func (s *SqlStorageTestSuite) TestSaveMessage() {
	message := Message{Attachments: []Attachment{{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}}

	s.Run("when message is new", func() {
		err := s.store.SaveMessage(&message)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Message{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)

		err = s.db.Model(&Attachment{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when message is existing", func() {
		message.TickerID = 1
		message.Attachments = append(message.Attachments, Attachment{ID: 2, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"})
		err := s.store.SaveMessage(&message)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Message{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)

		err = s.db.Model(&Attachment{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(2), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteMessage() {
	s.Run("when message does not exist", func() {
		message := Message{ID: 1}
		err := s.store.DeleteMessage(message)
		s.NoError(err)
	})

	s.Run("when message exists", func() {
		message := Message{ID: 1, Attachments: []Attachment{{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}}
		err := s.db.Create(&message).Error
		s.NoError(err)

		err = s.store.DeleteMessage(message)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Message{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)

		err = s.db.Model(&Attachment{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestDeleteMessages() {
	ticker := &Ticker{ID: 1}
	err := s.db.Create(&ticker).Error
	s.NoError(err)

	message := Message{ID: 1, TickerID: ticker.ID, Attachments: []Attachment{{ID: 1, MessageID: 1, UUID: "uuid", ContentType: "image/jpg", Extension: "jpg"}}}
	err = s.db.Create(&message).Error
	s.NoError(err)

	s.Run("when messages do not exist", func() {
		err := s.store.DeleteMessages(&Ticker{ID: 2})
		s.NoError(err)
	})

	s.Run("when messages exist", func() {
		err := s.store.DeleteMessages(ticker)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Message{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)

		err = s.db.Model(&Attachment{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(0), count)
	})
}

func (s *SqlStorageTestSuite) TestGetInactiveSettings() {
	s.Run("when no settings exist", func() {
		settings := s.store.GetInactiveSettings()
		s.Equal(DefaultInactiveSettings().Author, settings.Author)
	})

	s.Run("when settings exist", func() {
		setting := Setting{Name: SettingInactiveName, Value: `{"author":"test"}`}
		err := s.db.Create(&setting).Error
		s.NoError(err)

		settings := s.store.GetInactiveSettings()
		s.Equal("test", settings.Author)
	})
}

func (s *SqlStorageTestSuite) TestSaveInactiveSettings() {
	settings := InactiveSettings{Author: "test"}

	s.Run("when settings are new", func() {
		err := s.store.SaveInactiveSettings(settings)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Setting{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when settings are existing", func() {
		settings.Author = "test2"
		err := s.store.SaveInactiveSettings(settings)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Setting{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
		s.Equal("test2", settings.Author)
	})
}

func (s *SqlStorageTestSuite) TestGetRefreshIntervalSettings() {
	s.Run("when no settings exist", func() {
		settings := s.store.GetRefreshIntervalSettings()
		s.Equal(DefaultRefreshIntervalSettings().RefreshInterval, settings.RefreshInterval)
	})

	s.Run("when settings exist", func() {
		setting := Setting{Name: SettingRefreshInterval, Value: `{"refreshInterval":1000}`}
		err := s.db.Create(&setting).Error
		s.NoError(err)

		settings := s.store.GetRefreshIntervalSettings()
		s.Equal(1000, settings.RefreshInterval)
	})
}

func (s *SqlStorageTestSuite) TestSaveRefreshIntervalSettings() {
	settings := RefreshIntervalSettings{RefreshInterval: 1000}

	s.Run("when settings are new", func() {
		err := s.store.SaveRefreshIntervalSettings(settings)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Setting{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
	})

	s.Run("when settings are existing", func() {
		settings.RefreshInterval = 2000
		err := s.store.SaveRefreshIntervalSettings(settings)
		s.NoError(err)

		var count int64
		err = s.db.Model(&Setting{}).Count(&count).Error
		s.NoError(err)
		s.Equal(int64(1), count)
		s.Equal(2000, settings.RefreshInterval)
	})
}

func TestSqlStorageTestSuite(t *testing.T) {
	suite.Run(t, new(SqlStorageTestSuite))
}
