FROM golang:1.23.3-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY . .
RUN go build -ldflags "-w -s" -o /lumino ./cmd/lumino

FROM alpine:3.20

WORKDIR /app
COPY --from=builder /lumino /app/lumino

ENV LUMINO_ADDR=:8080
ENV LUMINO_DATA_DIR=/app/data

EXPOSE 8080
CMD ["/app/lumino"]
