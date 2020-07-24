FROM golang:1.14-alpine AS builder

ENV GO111MODULE=on

WORKDIR /bullshit
COPY . .
RUN go mod download
RUN go build -o bullshit

FROM alpine:latest

COPY --from=builder /bullshit/bullshit /
RUN mkdir generator
COPY generator/data.json generator/

EXPOSE 10000
ENTRYPOINT ["/bullshit"]
