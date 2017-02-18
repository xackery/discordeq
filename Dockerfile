FROM golang:1.6.2
MAINTAINER Xackery <xackery@gmail.com>

ENV GOPATH /go
ENV USER root

# pre-install known dependencies before the source, so we don't redownload them whenever the source changes
RUN go get github.com/xackery/eqemuconfig \
	&& go get gopkg.in/natefinch/lumberjack.v2 \
	&& go get github.com/bwmarrin/discordgo \
	&& go get github.com/ziutek/telnet

COPY . /go/src/github.com/xackery/discordeq

RUN cd /go/src/github.com/xackery/discordeq \
	&& go get -d -v \
	&& go install \
	&& go test github.com/xackery/discordeq...
