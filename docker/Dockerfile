FROM golang:latest AS builder
LABEL maintainer="keylo99 <keylo99official@gmail.com>"
ENV LC_ALL=C.UTF-8
WORKDIR /app/
COPY src .
RUN go get -d ./...
RUN go build -o pkgtop .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/pkgtop .
CMD ["./pkgtop"]