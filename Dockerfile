FROM ubuntu:jammy
ARG BINARY=discordgpt-linux-amd64
ARG DEBIAN_FRONTEND=noninteractive

RUN \
  apt update && \
  apt install ca-certificates --no-install-recommends -y && \
  apt clean && \
  rm -rf /var/lib/apt/lists/*

COPY ${BINARY} /usr/local/bin/discordgpt

CMD ["/usr/local/bin/discordgpt"]
