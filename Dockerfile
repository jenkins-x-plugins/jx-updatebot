FROM ghcr.io/jenkins-x/jx-boot:latest

LABEL org.opencontainers.image.source https://github.com/jenkins-x-plugins/jx-updatebot

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot
