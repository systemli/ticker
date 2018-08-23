FROM alpine:latest

EXPOSE 8080

ENV GOPATH=/go
ENV APPPATH=/$GOPATH/src/git.codecoop.org/systemli/ticker

COPY . $APPPATH

RUN apk add --update -t go \
    && apk add -t musl-dev \
    && cd $APPPATH \
    && go build -o /ticker \
    && apk del --purge go \
    && rm -rf $GOPATH

ENTRYPOINT ["/ticker"]
