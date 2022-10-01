FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
RUN apk add git openssh

ENV GOPATH=/go

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GO111MODULE=on

COPY . $GOPATH/src/github.com/s4kibs4mi/jally-commerce-bot
WORKDIR $GOPATH/src/github.com/s4kibs4mi/jally-commerce-bot

RUN go get -v .
RUN go build -v -o jally-commerce-bot
RUN mv twilfe /go/bin/jally-commerce-bot

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add bash

WORKDIR /root

COPY --from=builder /go/bin/jally-commerce-bot /usr/local/bin/jally-commerce-bot

ENTRYPOINT ["jally-commerce-bot"]
