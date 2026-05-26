---
title: Paridad Dev/Prod y Verificación Paso a Paso
version: 1.0
date_created: 2026-05-12
tags: process
---

# Paridad Dev/Prod y Verificación Paso a Paso

## 1. Propósito

Evitar que cambios funcionen en desarrollo pero no en producción por diferencias de entorno. No lanzar oleadas de código sin verificar cada paso.

## 2. Reglas

- **R1 — Un cambio, una verificación**: No commitear 5 archivos sin probar. Por cada cambio funcional (endpoint, DB, frontend), verificarlo antes de pasar al siguiente.

- **R2 — El entorno dev ES producción**: Si producción corre en Docker con PostgreSQL, desarrollo también. `docker compose up` con los mismos servicios. Prohibido desarrollar contra SQLite si producción usa PostgreSQL. Prohibido depender de binarios del host (gcc, node) que no estén en producción.

- **R3 — Smoke test mínimo antes de cada deploy**: Verificar manualmente:
  1. `GET /api/health` → 200
  2. Login funciona → JWT válido
  3. CRUD básico (crear tarea, mover columna)
  4. Frontend carga sin errores en consola del navegador

- **R4 — Nunca forzar builds completos sin verificar**: No rebuildear todo con `--no-cache` hasta entender qué cambió. Usar `--no-cache-filter=frontend-builder` para solo refrescar una etapa.

- **R5 — Los archivos generados (dist, .next, build/) no se eliminan del pipeline sin verificar el reemplazo**: Si se quita una etapa del Dockerfile que genera estos archivos, probar localmente primero.

- **R6 — La database URL se verifica desde dentro del contenedor**: `docker exec` o log del server al arrancar. No asumir que la variable está bien solo porque está en el dashboard.

## 3. Stack objetivo

| Capa | Dev | Producción |
|------|-----|------------|
| Servidor | `go run .` en contenedor | Mismo binario |
| DB | PostgreSQL en compose | PostgreSQL externo |
| Frontend | npm run build en contenedor | npm run build en contenedor |
| Build | Docker multi-stage sin CGO | Mismo Dockerfile |

## 4. Checklist pre-deploy

- [ ] `docker compose build` sin errores
- [ ] Contenedor arranca sin panics
- [ ] `curl localhost/api/health` → 200
- [ ] Login y JWT funcionan
- [ ] Frontend carga con estilos
- [ ] Variables de entorno existen en el hosting

## 5. Anti-patrones (no hacer)

- Usar SQLite en dev y PostgreSQL en prod
- Tener CGO en prod pero no en dev (o viceversa)
- Confiar en archivos locales (dist/) que no existen en CI/CD
- Hacer `--no-cache` nuclear sin saber qué se está solucionando
- Commitear código sin probar el endpoint o la conexión a DB
