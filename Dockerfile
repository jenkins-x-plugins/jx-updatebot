FROM golang:1.15

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot