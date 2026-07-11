FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w" \
    -o /app/application \
    ./cmd/api

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H appuser

COPY --from=builder /app/application ./application

USER appuser

ENV PORT=8080 \
    TZ=Asia/Jakarta

EXPOSE 8080

CMD ["./application"]
