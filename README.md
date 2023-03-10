# discordgpt

Discord bot using ChatGPT API

## Setup

```bash
GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -o discordgpt-linux-amd64
docker build . -t discordgpt:latest
```

## Run

```bash
docker run --rm -it --env-file .env discordgpt:latest
```
