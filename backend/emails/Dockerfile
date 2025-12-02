FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
COPY internal/ ./internal/

WORKDIR /build/emails

COPY emails/go.mod emails/go.sum ./
RUN go mod download

COPY emails/ ./

RUN GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

FROM alpine:latest

WORKDIR /app

# Install AWS Lambda Runtime Interface Emulator for local development
ADD https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie /usr/local/bin/aws-lambda-rie
RUN chmod +x /usr/local/bin/aws-lambda-rie

COPY --from=builder /build/emails/bootstrap /app/bootstrap

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/aws-lambda-rie", "/app/bootstrap"]
