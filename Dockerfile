FROM golang:1.15

ENTRYPOINT ["jx-updatebot"]

# helm 3
ENV HELM3_VERSION 3.5.0
RUN curl -f -L https://get.helm.sh/helm-v${HELM3_VERSION}-linux-386.tar.gz | tar xzv && \
    mv linux-386/helm /usr/bin

COPY ./build/linux/jx-updatebot /usr/bin/jx-updatebot

ENV XDG_CONFIG_HOME /home/.config
