FROM golang:1.13-alpine AS builder

RUN apk update && apk add --no-cache --update git gcc musl-dev

WORKDIR /app

ADD . /app

RUN cd /app

RUN go get -u ./...

RUN CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -ldflags="-w -s" -o bot . ; mv bot /app/


FROM alpine:latest

RUN apk update && apk add --no-cache --update ca-certificates tzdata

WORKDIR /root

COPY --from=builder /app/bot .

ADD ./database/migrations ./database/migrations

CMD ["./bot"]
