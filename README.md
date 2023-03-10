# discordgpt

[![Build](https://github.com/pycabbage/discordgpt/actions/workflows/build.yml/badge.svg)](https://github.com/pycabbage/discordgpt/actions/workflows/build.yml)

Discord bot using ChatGPT API

[Prebuilt binary](https://nightly.link/pycabbage/discordgpt/workflows/build/main)

## Usage

first, write token/channelID to `.env` file.

```env
DISCORD_TOKEN=[discord bot token]
DISCORD_CHANNELID=[space-separated channel ids to use the bot]
GPT_SECRET_KEY=[OpenAI APIKEY]
```

### Use [prebuilt binary](https://nightly.link/pycabbage/discordgpt/workflows/build/main)/[prebuilt image](https://github.com/users/pycabbage/packages/container/package/discordgpt)

#### Run [prebuilt binary](https://nightly.link/pycabbage/discordgpt/workflows/build/main) locally

```bash
curl -kLO https://nightly.link/pycabbage/discordgpt/workflows/build/main/discordgpt-linux-amd64.zip
unzip discordgpt-linux-amd64.zip
chmod +x discordgpt-linux-amd64
./discordgpt-linux-amd64
```

#### Run [prebuilt image](https://github.com/users/pycabbage/packages/container/package/discordgpt)

```bash
docker run --rm -it --env-file .env ghcr.io/pycabbage/discordgpt:latest
# or specify token/channelID
docker run --rm -it \
  -e DISCORD_TOKEN=[discord bot token] \
  -e DISCORD_CHANNELID=[space-separated channel ids to use the bot] \
  -e GPT_SECRET_KEY=[OpenAI APIKEY] \
  ghcr.io/pycabbage/discordgpt:latest
```

### Build Mannualy

#### Run locally

```bash
# Build binary
go build -o discordgpt
# Run
./discordgpt
```

#### Run on container

```bash
# Build binary
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o discordgpt-linux-amd64
# Build container
docker build . -t discordgpt:latest --build-arg BINARY=discordgpt-linux-amd64
# Run container
docker run --rm -it --env-file .env discordgpt:latest
```
