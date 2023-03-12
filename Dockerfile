FROM ghcr.io/jenkins-x/jx-go:latest

LABEL org.opencontainers.image.source https://github.com/jenkins-x-plugins/jx-updatebot

ENTRYPOINT ["jx-updatebot"]

RUN apk --no-cache add git
RUN jx upgrade plugins --boot --path /usr/bin

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot
