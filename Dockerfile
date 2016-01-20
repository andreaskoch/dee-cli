FROM golang:latest
MAINTAINER Andreas Koch <andy@ak7.io>

# Add sources
ADD . /go/src/github.com/andreaskoch/dnsimple-cli
WORKDIR /go/src/github.com/andreaskoch/dnsimple-cli

# Build
RUN make crosscompile
