# stage 1: build the go service
FROM golang:1.26.1-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o ledger .

# stage 2: run it
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/ledger .
CMD ["./ledger"]
