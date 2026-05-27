# Usar binario precompilado
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copiar binario compilado
COPY tablero .

# Copiar frontend dist
COPY dist ./frontend/dist

EXPOSE 8080

CMD ["./tablero"]
