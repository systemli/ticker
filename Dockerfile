FROM golang:1.13-alpine AS build-env
ENV GO111MODULE=on
WORKDIR /go/src/github.com/systemli/ticker
ADD . /go/src/github.com/systemli/ticker
RUN apk update && apk add git gcc libc-dev
RUN go build -o /ticker

FROM alpine
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=build-env /ticker /ticker

EXPOSE 8080
ENTRYPOINT ["/ticker"]
