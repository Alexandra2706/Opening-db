FROM golang:1.22
WORKDIR /api

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /opdb-api

CMD ["/opdb-api"]
