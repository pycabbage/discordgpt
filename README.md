# discordgpt

[![Build](https://github.com/pycabbage/discordgpt/actions/workflows/build.yml/badge.svg)](https://github.com/pycabbage/discordgpt/actions/workflows/build.yml)

Discord bot using ChatGPT API

## Run locally

```bash
# Build binary
go build -o discordgpt
# Run
./discordgpt
```

## Run on container

```bash
# Build binary
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o discordgpt-linux-amd64
# Build container
docker build . -t discordgpt:latest --build-arg BINARY=discordgpt-linux-amd64
# Run container
docker run --rm -it --env-file .env discordgpt:latest
```
