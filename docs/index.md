# Systemli Ticker

**Service to distribute short messages in support of events, demonstrations, or other time-sensitive events.**

This repository contains the API and storage for
the [Systemli Ticker Project](https://www.systemli.org/en/service/ticker.html).

!!! note "Requirements"

    The project is written in Go. You should be familiar with the structure and organisation of the code. If not, there are
    some [good guides](https://golang.org/doc/code.html).

## First run

- make sure you have pulled git submodules
```shell
git clone --recurse-submodules git@github.com:systemli/ticker.git
```

- or if you already cloned the repo
```shell
cd <path-to-ticker>
git submodule update --init --recursive
```

- start the ticker
```shell
cp config.yml.dist config.yml
go run .
```

Now you have a running ticker API!

## Testing

```shell
go test ./...
```