# reference:
#   https://hub.docker.com/_/golang
#   https://hub.docker.com/_/alpine

# build
FROM golang:1.20 as builder
WORKDIR /build
COPY . .
RUN go env -w GO111MODULE=on && \
    go env -w CGO_ENABLED=0 && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o main .

# run
FROM alpine:3.17
WORKDIR /
COPY --from=builder /build/main .
ENTRYPOINT ["./main"]