FROM golang:1.19-alpine as builder
RUN apk add build-base
COPY . /app
WORKDIR /app
RUN go env -w GOPROXY="https://goproxy.cn,direct"
RUN go build -ldflags="-w -s" -o zura main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/zura /app/zura
COPY ./config.yaml /app/config.yaml
ENTRYPOINT ["app/zura"]