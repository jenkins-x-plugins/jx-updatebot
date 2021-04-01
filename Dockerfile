FROM ghcr.io/jenkins-x/jx-go:3.1.353

ENTRYPOINT ["jx-updatebot"]

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot

ENV XDG_CONFIG_HOME /home/.config
