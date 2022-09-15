FROM golang:1.17-alpine AS builder

RUN apk update && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*
RUN apk add git openssh

ENV GOPATH=/go

ENV GOOS="linux"
ENV GOARCH="amd64"
ENV GO111MODULE=on

COPY . $GOPATH/src/github.com/s4kibs4mi/twilfe
WORKDIR $GOPATH/src/github.com/s4kibs4mi/twilfe

RUN go get -v .
RUN go build -v -o twilfe
RUN mv twilfe /go/bin/twilfe

FROM alpine

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN apk add bash

WORKDIR /root

COPY --from=builder /go/bin/twilfe /usr/local/bin/twilfe

ENTRYPOINT ["twilfe"]
