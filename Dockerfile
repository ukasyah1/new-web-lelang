FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -o /server ./cmd/api

FROM alpine:3.20
RUN adduser -D -H appuser && mkdir /data && chown appuser:appuser /data
USER appuser
COPY --from=builder /server /server
ENV SQLITE_PATH=/data/lelang.db
ENV PORT=8080
VOLUME ["/data"]
EXPOSE 8080
ENTRYPOINT ["/server"]
