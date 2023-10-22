# Configuration

```yaml
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# log_format sets log format for logrus (default: json)
log_format: "json"
# initiator is the email for the first admin user (see password in logs)
initiator: "admin@systemli.org"
# configuration for the database
database:
    type: "sqlite" # postgres, mysql, sqlite
    dsn: "ticker.db" # postgres: "host=localhost port=5432 user=ticker dbname=ticker password=ticker sslmode=disable"
# secret used for JSON Web Tokens
secret: "slorp-panfil-becall-dorp-hashab-incus-biter-lyra-pelage-sarraf-drunk"
# telegram configuration
telegram_bot_token: ""
# listen port for prometheus metrics exporter
metrics_listen: ":8181"
# path where to store the uploaded files
upload_path: "uploads"
# base url for uploaded assets
upload_url: "http://localhost:8080"
```

!!! note
    We use [viper](https://github.com/spf13/viper). That means you can use any of the supported file formats. Env variables
    will overwrite existing config file values. Note that every env variable MUST be prefixed by: `TICKER_`.
    E.g. `TICKER_DATABASE`.

The following env vars can be used:

* `TICKER_LISTEN`
* `TICKER_LOG_FORMAT`
* `TICKER_LOG_LEVEL`
* `TICKER_DATABASE_TYPE`
* `TICKER_DATABASE_DSN`
* `TICKER_LOG_LEVEL`
* `TICKER_INITIATOR`
* `TICKER_SECRET`
* `TICKER_TELEGRAM_BOT_TOKEN`
* `TICKER_METRICS_LISTEN`
* `TICKER_UPLOAD_PATH`
* `TICKER_UPLOAD_URL`
