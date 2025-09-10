# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG VERSION="dev"

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN go build -o newsapi cmd/newsapi/main.go

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/newsapi ./
CMD ["./newsapi"]
ENTRYPOINT ["./newsapi"]