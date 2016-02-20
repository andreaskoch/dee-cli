FROM golang:latest
MAINTAINER Andreas Koch <andy@ak7.io>

# Add sources
ADD . /go/src/github.com/andreaskoch/dee-cli
WORKDIR /go/src/github.com/andreaskoch/dee-cli

# Cross-Compile
RUN make crosscompile

# Compile
RUN make install

# Link binary to public folder
RUN ln -s `pwd`/bin/dee-cli /bin/dee

ENTRYPOINT ["/bin/dee"]
