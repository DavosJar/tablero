# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Build backend (NO CGO, no gcc, no sqlite)
FROM golang:1.26-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/dist ./frontend/dist
RUN go build -o tablero .

# Stage 3: Runtime
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=backend-builder /app/tablero .
COPY --from=backend-builder /app/frontend/dist ./frontend/dist
EXPOSE 8080
CMD ["./tablero"]
