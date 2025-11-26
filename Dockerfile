FROM golang:1.25.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/proxio main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=builder /app/proxio .
COPY --from=builder /app/.env .

EXPOSE 1337
ENTRYPOINT [ "/proxio" ]
