FROM golang:1.24-alpine AS builder

WORKDIR /app

ENV GO111MODULE=on
ENV GOOS=linux
ENV CGO_ENABLED=0

RUN apk add --no-cache build-base ca-certificates git

COPY configs/development.env .
COPY go.mod .
COPY go.sum .

RUN go install github.com/cosmtrek/air@v1.28.0
RUN go install github.com/swaggo/swag/cmd/swag@v1.16.2
# RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go mod download -x

COPY . .

EXPOSE 9000 40000

CMD ["air", "-c", ".air.toml"]
