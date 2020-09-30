FROM gcr.io/jenkinsxio-labs-private/jxl-base:0.0.52

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot