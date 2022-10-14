FROM golang:1.19.2 AS builder
ENV LC_ALL=C.UTF-8
WORKDIR /app/
COPY . .
RUN go build -o pkgtop cmd/pkgtop.go

FROM debian:buster-slim
WORKDIR /root/
COPY --from=builder /app/pkgtop .
CMD ["./pkgtop"]
