FROM golang:1.18 AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/redis-exporter .

FROM alpine:3.16
COPY --from=builder /src/bin /bin
CMD ["redis-exporter"]