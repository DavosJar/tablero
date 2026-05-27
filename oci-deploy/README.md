# OCI Deployment

Directorio optimizado para desplegar Tablero en OCI con Docker Compose.

## Setup Rápido

1. **Copiar variables de entorno:**
   ```bash
   cp .env.example .env
   # Editar .env con valores reales
   ```

2. **Iniciar servicios:**
   ```bash
   docker-compose up -d
   ```

3. **Verificar:**
   ```bash
   curl http://localhost:8080/api/health
   ```

## Estructura

- `tablero` - Binario compilado (Linux amd64)
- `Dockerfile` - Imagen minimal con Alpine
- `docker-compose.yml` - PostgreSQL + App
- `.env.example` - Template de variables de entorno

## Notas

- El binario está pre-compilado con `GOOS=linux GOARCH=amd64`
- PostgreSQL se ejecuta en un contenedor
- Los datos persisten en volumen `postgres_data`

## Logs

```bash
docker-compose logs -f
```

## Detener

```bash
docker-compose down
```
