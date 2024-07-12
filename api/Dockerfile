FROM golang:1.22-alpine AS builder
WORKDIR /api

COPY go.mod go.sum *.go ./
COPY routes ./

RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN go generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /opdb-api

FROM alpine

COPY --from=builder /opdb-api /opdb-api

CMD ["/opdb-api"]
