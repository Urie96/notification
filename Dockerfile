FROM golang:1.16.10-alpine AS builder
WORKDIR /root
ENV GOPROXY=https://goproxy.cn,direct
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -tags netgo -v -o service .

FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/service .
CMD ["./service"]