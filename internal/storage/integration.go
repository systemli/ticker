package storage

// Integration identifies an outbound platform configured on a Ticker
// (Telegram / Mastodon / Bluesky / SignalGroup).
type Integration string

const (
	IntegrationTelegram    Integration = "telegram"
	IntegrationMastodon    Integration = "mastodon"
	IntegrationBluesky     Integration = "bluesky"
	IntegrationSignalGroup Integration = "signal_group"
)

// AllIntegrations lists every supported integration; useful for ClearIntegrations
// and for iterating across known platforms.
var AllIntegrations = []Integration{
	IntegrationTelegram,
	IntegrationMastodon,
	IntegrationBluesky,
	IntegrationSignalGroup,
}
