# Build
FROM golang:1.20 AS builder
WORKDIR /alpacon-cli
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o alpacon-cli


# Deploy
FROM alpine:latest
COPY --from=builder /alpacon-cli/alpacon-cli /usr/local/bin/alpacon
ENTRYPOINT ["alpacon"]