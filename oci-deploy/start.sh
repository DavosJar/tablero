#!/bin/bash
set -e

echo "🚀 Iniciando Tablero..."

# Verificar .env
if [ ! -f .env ]; then
    echo "⚠️  Falta .env, copiando desde .env.example..."
    cp .env.example .env
    echo "✏️  Edita .env antes de continuar"
    exit 1
fi

echo "🐳 Levantando servicios..."
docker-compose up -d

echo ""
echo "✨ Tablero corriendo"
echo "📍 http://localhost:8080"
