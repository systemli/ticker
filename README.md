# ticker [![Integration](https://github.com/systemli/ticker/workflows/Integration/badge.svg)](https://github.com/systemli/ticker/actions) [![Quality](https://github.com/systemli/ticker/workflows/Quality/badge.svg)](https://github.com/systemli/ticker/actions) [![codecov](https://codecov.io/gh/systemli/ticker/branch/master/graph/badge.svg)](https://codecov.io/gh/systemli/ticker) [![Docker Automated build](https://img.shields.io/docker/automated/systemli/ticker.svg)](https://hub.docker.com/r/systemli/ticker/) [![MicroBadger Size](https://img.shields.io/microbadger/image-size/systemli/ticker.svg)](https://hub.docker.com/r/systemli/ticker/)

This repository contains the API and storage for the [Systemli Ticker Project](https://www.systemli.org/en/service/ticker.html).

## Requirements

The project is written in Go. You should be familiar with the structure and organisation of the code. If not, there are some [good guides](https://golang.org/doc/code.html).

## First run 

  * we provide a `Makefile` for clean, build, test and release the software

```
➜  ticker git:(master) ✗ make run
go build -o build/ticker -v
cp config.yml.dist build/config.yml
./build/ticker -config build/config.yml
INFO[0000] admin user created (change password now!)     email=admin@systemli.org password="5O.AVsHDd@Y23<aGWlxpwKiS"
INFO[0000] starting ticker at localhost:8080
```

Now you have a running ticker API!

## Configuration

  * Example config.yml.dist

```
# listen binds ticker to specific address and port
listen: "localhost:8080"
# log_level sets log level for logrus
log_level: "error"
# initiator is the email for the first admin user (see password in logs)
initiator: "admin@systemli.org"
# database is the path to the bolt file
database: "ticker.db"
# secret used for JSON Web Tokens
secret: "slorp-panfil-becall-dorp-hashab-incus-biter-lyra-pelage-sarraf-drunk"
# twitter configuration
twitter_consumer_key: ""
twitter_consumer_secret: ""
# listen port for prometheus metrics exporter
metrics_listen: ":8181"
```

We use [viper](https://github.com/spf13/viper). That means you can use any of the supported
file formats. Env variables will overwrite existing config file values.
Note that every env variable MUST be prefixed by: `TICKER_`. E.g. `TICKER_DATABASE`.

The following env vars can be used: 
* TICKER_DATABASE
* TICKER_LISTEN
* TICKER_LOG_LEVEL
* TICKER_INITIATOR
* TICKER_SECRET
* TICKER_TWITTER_CONSUMER_KEY
* TICKER_TWITTER_CONSUMER_SECRET
* TICKER_METRICS_LISTEN

## Testing

```
make test
```
