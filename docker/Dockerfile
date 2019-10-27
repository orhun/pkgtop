FROM golang:1.13.3 AS builder
ENV LC_ALL=C.UTF-8
WORKDIR /app/
COPY src .
RUN go get -d ./...
RUN go build -o pkgtop .

FROM alpine:3.10.2
WORKDIR /root/
COPY --from=builder /app/pkgtop .
CMD ["./pkgtop"]
