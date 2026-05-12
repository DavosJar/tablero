FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o tablero .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/tablero .
COPY --from=builder /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./tablero"]
