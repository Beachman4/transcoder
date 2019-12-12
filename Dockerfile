FROM golang:1.11-alpine as builder

WORKDIR /go/src/github.com/flocasts/transcoding

ADD . .

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/* && apk add git

RUN go get -d ./...

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o transcoding

FROM jrottenberg/ffmpeg:4.1-alpine

RUN apk add bash

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/flocasts/transcoding/transcoding ./transcoding
COPY --from=builder /go/src/github.com/flocasts/transcoding/hls ./hls

EXPOSE 8080

ENTRYPOINT ["/app/transcoding"]