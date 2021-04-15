FROM ghcr.io/jenkins-x/jx-go:3.2.42

ENTRYPOINT ["jx-updatebot"]

RUN jx gitops plugin upgrade

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot
