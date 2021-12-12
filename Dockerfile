FROM golang:1.13-alpine as build

RUN apk add --no-cache bash
RUN apk add build-base

WORKDIR /go/src/location

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build

FROM alpine:latest

RUN apk update
RUN apk upgrade
RUN apk add curl
RUN apk add wget

COPY --from=build /go/src/location /usr/bin

ENTRYPOINT ["location"]