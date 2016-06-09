FROM golang:1.6.2-alpine

RUN apk add --update git && apk add --update make && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/r3labs/definition-mapper
WORKDIR /go/src/github.com/r3labs/definition-mapper

RUN make deps && go install

ENTRYPOINT /go/bin/definition-mapper
