FROM golang:1.25.4
WORKDIR /app

COPY . .

RUN mkdir -p config

ENTRYPOINT ["go", "run", "-gcflags=all=-lang=go1.25", "cmd/sso/main.go", "--config=./config/local.yaml"]