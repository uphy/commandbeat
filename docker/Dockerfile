FROM ubuntu:14.04

ARG VERSION=
RUN apt-get update && \
    apt-get install -y curl && \
    apt-get clean && \
    rm -rf /var/cache/apt/archives/* /var/lib/apt/lists/*
RUN curl -L -O https://github.com/uphy/commandbeat/releases/download/$VERSION/commandbeat-7.0.0-alpha1-SNAPSHOT-amd64.deb && \
    dpkg -i commandbeat-*.deb && \
    rm -f commandbeat-*.deb

ENTRYPOINT [ "commandbeat" ]
CMD [ "-e"]
