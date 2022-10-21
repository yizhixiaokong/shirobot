#build stage
FROM golang:1.18-alpine AS builder
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
      && apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:latest
RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories \
      && apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=shirobot Version=0.0.1
EXPOSE 3000
