FROM golang:1.25-alpine

WORKDIR /build

COPY go.mod go.sum ./
COPY internal/ ./internal/

WORKDIR /build/eventform

COPY eventform/go.mod eventform/go.sum ./
RUN go mod download

COPY eventform/ ./

CMD ["go", "run", "main.go"]
