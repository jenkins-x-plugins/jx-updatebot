FROM gcr.io/jenkinsxio/jx-cli-base:0.0.23

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot