FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o chentoons ./cmd

FROM alpine:3.20

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/chentoons .
COPY --from=builder /app/docs ./docs
COPY --from=builder /app/uploads ./uploads

EXPOSE 8080
CMD ["./chentoons"]
