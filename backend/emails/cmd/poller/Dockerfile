FROM golang:1.25-alpine

WORKDIR /app

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY ../.. .

WORKDIR /app/cmd/poller

CMD ["go", "run", "main.go"]
