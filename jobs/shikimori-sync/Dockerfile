FROM golang:1.22-alpine AS builder
WORKDIR /shikimori-sync

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY postgres ./postgres
COPY shikimori-api ./shikimori-api

RUN CGO_ENABLED=0 GOOS=linux go build -o /opdb-shikimori-sync

FROM alpine

COPY --from=builder /opdb-shikimori-sync /opdb-shikimori-sync
RUN chmod +x /opdb-shikimori-sync

CMD ["/opdb-shikimori-sync"]
