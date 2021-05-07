FROM ghcr.io/jenkins-x/jx-go:latest

ENTRYPOINT ["jx-updatebot"]

RUN jx upgrade plugins --boot --path /usr/bin

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot
