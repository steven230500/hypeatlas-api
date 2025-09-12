# HypeAtlas – Deploy en Droplet (API + Worker)

Este stack levanta **API** (Go + chi) y **Worker** (ingesta periódica) en Docker,
con reverse proxy y TLS servidos por un Caddy central en otra pila.
La API publica documentación Swagger y está detrás de `https://api.hypeatlas.app`.

## Rutas públicas

- **Docs (Swagger UI):** https://api.hypeatlas.app/docs
- **Esquema OpenAPI:**  https://api.hypeatlas.app/openapi.yaml
- **Healthcheck:**      https://api.hypeatlas.app/healthz

> Si pruebas desde el droplet vía loopback:  
> `curl -H "Host: api.hypeatlas.app" http://127.0.0.1/healthz`

## Directorio y archivos

- **Ruta en servidor:** `/opt/hypeatlas`
- **Compose principal:** `/opt/hypeatlas/docker-compose.prod.yml`
- **Variables:** `/opt/hypeatlas/.env.prod` (compartido por API y Worker)

## Servicios

- `hypeatlas-api`
  - Healthcheck interno cada 15s (`/healthz`)
  - Conectado a redes: `hypeatlas` (interna), `edge` (externa compartida con Caddy)
  - No publica puertos al host; Caddy enruta por `api.hypeatlas.app -> hypeatlas-api:8080`
- `hypeatlas-worker`
  - Corre ciclos de ingesta cada 30s (Twitch, etc.)
  - Usa el mismo `.env.prod` (DB / API keys)
  - Redes: `hypeatlas`, `edge`

## Operaciones comunes

### Levantar / actualizar

```bash
cd /opt/hypeatlas
docker compose -f docker-compose.prod.yml pull
docker compose -f docker-compose.prod.yml up -d --remove-orphans
