# Tablero Kanban Fullstack - Go + Astro

Una aplicación de tablero Kanban personal responsiva y fácil de usar, construida con Go (backend) y Astro (frontend).

## Características

✨ **Tablero Kanban completo**
- Columnas configurables: "Por hacer", "En progreso", "Hecho"
- Tarjetas con título, descripción, prioridad y fecha límite
- Drag & drop entre columnas
- Crear, editar y eliminar tareas

🔒 **Seguridad robusta**
- Login con contraseña única
- Autenticación JWT
- Headers de seguridad HTTP
- Rate limiting en login

💾 **Base de datos flexible**
- SQLite para desarrollo local (automático)
- PostgreSQL para producción (vía DATABASE_URL)
- Migraciones automáticas con GORM

📦 **Un único binario**
- Frontend Astro embebido en Go
- No requiere dependencias externas
- Compatible con Docker

## Stack Técnico

- **Backend**: Go 1.26 + Gin + GORM
- **Frontend**: Astro + Vanilla JavaScript
- **Base de datos**: SQLite / PostgreSQL
- **Autenticación**: JWT
- **Despliegue**: Docker, Render, Railway

## Instalación Local

### Requisitos
- Go 1.26+
- Node.js 20+
- npm o yarn

### Pasos

1. **Clonar el repositorio**
```bash
git clone <repo-url>
cd tablero
```

2. **Instalar dependencias**
```bash
npm install --prefix frontend
go mod download
```

3. **Compilar frontend**
```bash
npm run build --prefix frontend
```

4. **Compilar y ejecutar backend**
```bash
APP_PASSWORD=test JWT_SECRET=secret123 go run main.go
```

O en un solo comando:
```bash
APP_PASSWORD=test JWT_SECRET=secret123 npm run build --prefix frontend && go run main.go
```

5. **Acceder a la aplicación**
Abre http://localhost:8080 en tu navegador

Usa contraseña: `test`

## Variables de Entorno

```bash
# Requeridas
APP_PASSWORD=tu_contraseña_secreta        # Contraseña para acceder
JWT_SECRET=tu_jwt_secret_muy_largo        # Clave para firmar JWTs

# Opcionales
PORT=8080                                  # Puerto (default: 8080)
DATABASE_URL=                              # PostgreSQL URL (si no está → SQLite)
GIN_MODE=release                           # release/debug (default: debug)
```

### Ejemplo con PostgreSQL

```bash
APP_PASSWORD=mypass \
JWT_SECRET=mysecret \
DATABASE_URL="postgresql://user:password@localhost/tablero" \
PORT=8080 \
go run main.go
```

## Estructura del Proyecto

```
tablero/
├── main.go                    # Servidor principal con embedding
├── go.mod / go.sum           # Dependencias Go
├── internal/
│   ├── models/               # Modelos GORM (Task, Column)
│   ├── handlers/             # Handlers de rutas API
│   ├── middleware/           # JWT, seguridad
│   ├── database/             # Inicialización de BD
│   └── utils/                # Utilidades (JWT)
├── frontend/
│   ├── src/
│   │   ├── pages/            # Páginas Astro (index.astro)
│   │   ├── components/       # Componentes reutilizables
│   │   ├── layouts/          # Layouts base
│   │   └── styles/           # CSS global
│   ├── package.json          # Dependencias frontend
│   ├── astro.config.mjs      # Configuración Astro
│   └── dist/                 # Salida compilada (embebida en Go)
├── Dockerfile                # Multi-stage para Docker
├── render.yaml               # Configuración Render
├── railway.toml              # Configuración Railway
├── .env.example              # Ejemplo de variables
└── README.md                 # Este archivo
```

## API Reference

### Autenticación
- `POST /api/login` - Login con contraseña
- `POST /api/logout` - Cierra sesión

### Tareas
- `GET /api/tasks` - Obtener todas las tareas
- `POST /api/tasks` - Crear tarea
- `PUT /api/tasks/:id` - Actualizar tarea
- `DELETE /api/tasks/:id` - Eliminar tarea

### Columnas
- `GET /api/columns` - Obtener columnas
- `POST /api/columns` - Crear columna
- `PUT /api/columns/:id` - Actualizar columna
- `DELETE /api/columns/:id` - Eliminar columna

## Despliegue

### Docker Local

```bash
docker build -t tablero .
docker run -e APP_PASSWORD=test -e JWT_SECRET=secret123 -p 8080:8080 tablero
```

### Render

1. Conecta tu repositorio a Render
2. Crea nuevo servicio Web
3. Establece variables de entorno:
   - `APP_PASSWORD`
   - `JWT_SECRET`
4. Render usará automáticamente DATABASE_URL con PostgreSQL free

### Railway

1. Conecta tu repositorio a Railway
2. Railway detectará el archivo `railway.toml`
3. Establece secretos:
   - `APP_PASSWORD`
   - `JWT_SECRET`
4. Database: Agrega PostgreSQL plugin

## Desarrollo

### Frontend en modo watch
```bash
cd frontend
npm run dev
```

### Backend con auto-reload
```bash
# Instalar air: https://github.com/cosmtrek/air
go install github.com/cosmtrek/air@latest

# Ejecutar con watch
APP_PASSWORD=test JWT_SECRET=secret123 air
```

## Modelado de Datos

### Column
```go
type Column struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"not null;unique"`
    Order     int       `gorm:"default:0"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Tasks     []Task    `gorm:"foreignKey:ColumnID"`
}
```

### Task
```go
type Task struct {
    ID          uint      `gorm:"primaryKey"`
    Title       string    `gorm:"not null"`
    Description string
    Priority    string    `gorm:"default:'media'"` // baja, media, alta
    DueDate     *time.Time
    ColumnID    uint      `gorm:"not null"`
    Order       int       `gorm:"default:0"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

## Seguridad

- **CORS**: Habilitado con configuración permisiva (puede ajustarse en producción)
- **CSP**: Content Security Policy básica
- **Headers**: X-Frame-Options, X-Content-Type-Options, X-XSS-Protection
- **JWT**: HS256 con expiración de 24 horas
- **Rate Limiting**: 5 intentos de login por minuto por IP

## Solución de Problemas

### "Port already in use"
```bash
# Cambiar puerto
PORT=3000 go run main.go
```

### "Could not connect to database"
- Verificar DATABASE_URL si usas PostgreSQL
- SQLite usa data.db (creado automáticamente)

### "Invalid token"
- Token expiró después de 24 horas
- Cierra sesión y vuelve a acceder

### "Unauthorized"
- Falta header Authorization
- Token inválido o expirado

## Licencia

MIT

## Autor

Tablero Kanban - 2026
