# Systemli Ticker

**Service to distribute short messages in support of events, demonstrations, or other time-sensitive events.**

This repository contains the API for the [Systemli Ticker Project](https://www.systemli.org/en/service/ticker.html).

!!! note "Requirements"

    The project is written in Go. You should be familiar with the structure and organisation of the code. If not, there are
    some [good guides](https://golang.org/doc/code.html).

## First run

1. Clone the project

        git clone https://github.com/systemli/ticker.git

2. Start the project

        cd ticker
        go run . run

3. Check the API

        curl http://localhost:8080/healthz

4. Create a user

        go run . user create --email <email-address> --password <password> --super-admin

## Testing

```shell
go test ./...
```
