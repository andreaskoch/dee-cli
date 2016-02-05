FROM golang:latest
MAINTAINER Andreas Koch <andy@ak7.io>

# Add sources
ADD . /go/src/github.com/andreaskoch/dee-cli
WORKDIR /go/src/github.com/andreaskoch/dee-cli

# Build
RUN make crosscompile
