###Build container
FROM golang:1.8-alpine as builder
WORKDIR /go/src/github.com/nawa/back-friend/
COPY ./ .

RUN apk --no-cache add git
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -o backfriend .

###Main container
FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/nawa/back-friend/backfriend .
VOLUME /app/config
EXPOSE 8080

CMD ["./backfriend","--config","./config/config.yml"]