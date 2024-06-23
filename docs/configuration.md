# Configuration

```yaml
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# log_format sets log format for logrus (default: json)
log_format: "json"
# configuration for the database
database:
  type: "sqlite" # postgres, mysql, sqlite
  dsn: "ticker.db" # postgres: "host=localhost port=5432 user=ticker dbname=ticker password=ticker sslmode=disable"
# secret used for JSON Web Tokens
secret: "slorp-panfil-becall-dorp-hashab-incus-biter-lyra-pelage-sarraf-drunk"
# telegram configuration
telegram:
  token: "" # telegram bot token
# signal group configuration
signal_group:
  api_url: "" # URL to your signal cli (https://github.com/AsamK/signal-cli)
  avatar: "" # URL to the avatar for the signal group
  account: "" # phone number for the signal account
# listen port for prometheus metrics exporter
metrics_listen: ":8181"
upload:
  # path where to store the uploaded files
  path: "uploads"
  # base url for uploaded assets
  url: "http://localhost:8080"
```

!!! note
    All configuration options can be set via environment variables.

The following env vars can be used:

* `TICKER_LISTEN`
* `TICKER_LOG_FORMAT`
* `TICKER_LOG_LEVEL`
* `TICKER_DATABASE_TYPE`
* `TICKER_DATABASE_DSN`
* `TICKER_LOG_LEVEL`
* `TICKER_INITIATOR`
* `TICKER_SECRET`
* `TICKER_TELEGRAM_TOKEN`
* `TICKER_SIGNAL_GROUP_API_URL`
* `TICKER_SIGNAL_GROUP_AVATAR`
* `TICKER_SIGNAL_GROUP_ACCOUNT`
* `TICKER_METRICS_LISTEN`
* `TICKER_UPLOAD_PATH`
* `TICKER_UPLOAD_URL`
