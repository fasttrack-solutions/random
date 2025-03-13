FROM golang:1.24-alpine3.21 AS build-env

# Update apk
RUN apk add --update --no-cache git gcc musl-dev openssh
RUN apk add build-base

# Use git with SSH instead of https
RUN git config --global url."git@github.com:".insteadOf https://github.com/
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/http ./cmd/http
RUN go build -o bin/grpc ./cmd/grpc

##
## Deploy
##
FROM alpine:3.21

WORKDIR /app

COPY --from=build-env /app/bin/http ./http
COPY --from=build-env /app/bin/grpc ./grpc

# dynamic entry point
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
