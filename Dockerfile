FROM golang:1.15

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot

ENV XDG_CONFIG_HOME /home/.config
