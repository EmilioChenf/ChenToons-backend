# ChenToons Backend

ChenToons es una API REST sencilla para llevar un tracker de series y caricaturas como Pocoyo, Escandalosos, Snoopy, Bluey y otras parecidas. El proyecto esta hecho con Go, Fiber, PostgreSQL y Docker, pensando en que sea facil de correr, explicar y conectar con un frontend universitario.

## Tecnologias

- Go 1.22
- Fiber
- PostgreSQL
- Docker y Docker Compose
- OpenAPI / Swagger UI
- CSV para exportar series y abrirlas en Excel

## Variables de entorno

El backend funciona con variables separadas para Docker local o con `DATABASE_URL` para Render.

Variables principales:

- `PORT`: puerto que usa Render.
- `APP_PORT`: puerto local opcional. Si existe, tiene prioridad sobre `PORT`.
- `APP_ENV`: `development` o `production`.
- `UPLOAD_DIR`: carpeta donde se guardan imagenes. En Docker local se usa `/app/uploads`; en Render se monta como disco persistente.
- `DATABASE_URL`: URL completa de PostgreSQL para Render.
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`: alternativa local si no usas `DATABASE_URL`.
- `CORS_ORIGIN`: origen permitido para el frontend. Para pruebas puede ser `*`.

## Correr con Docker

1. Crea el archivo `.env` desde el ejemplo:

```bash
cp .env.example .env
```

En Windows tambien puedes copiarlo manualmente.

2. Levanta backend y base de datos:

```bash
docker compose up --build
```

3. Prueba que vive:

```bash
curl http://localhost:8080/health
```

PostgreSQL corre en el puerto `5432` y el backend en `8080`.

## Endpoints principales

- `GET /health`
- `GET /series`
- `GET /series/:id`
- `POST /series`
- `PUT /series/:id`
- `DELETE /series/:id`
- `GET /personajes`
- `POST /personajes`
- `GET /episodios`
- `POST /episodios`
- `GET /series/:id/ratings`
- `POST /series/:id/ratings`
- `GET /series/:id/promedio-rating`
- `POST /uploads`
- `GET /uploads/:filename`
- `GET /export/series.csv`

## Filtros de series

`GET /series` acepta:

- `search`
- `q` (alias de `search`)
- `genero`
- `categoria`
- `estado`
- `sort`
- `order`
- `page`
- `limit`

Ejemplo:

```bash
curl "http://localhost:8080/series?search=snoopy&page=1&limit=8&sort=nombre&order=asc"
```

Tambien puedes usar:

```bash
curl "http://localhost:8080/series?q=snoopy"
```

## Swagger

Abre esta ruta en el navegador:

```text
http://localhost:8080/docs
```

Tambien funciona:

```text
http://localhost:8080/swagger
```

El archivo OpenAPI esta en `docs/openapi.yaml`.

## Probar CORS

El backend permite CORS usando la variable `CORS_ORIGIN`. Para desarrollo viene como `*`.

Desde un frontend local puedes llamar:

```js
fetch("http://localhost:8080/series")
  .then((res) => res.json())
  .then(console.log);
```

## Subir imagenes

El campo del formulario debe llamarse `imagen`.
El backend acepta imagenes `.jpg`, `.jpeg`, `.png`, `.webp` y `.gif` de hasta 1MB.

```bash
curl -X POST http://localhost:8080/uploads \
  -F "imagen=@snoopy.jpg"
```

La respuesta devuelve una ruta como:

```json
{
  "mensaje": "imagen subida",
  "ruta": "/uploads/1710000000-snoopy.jpg"
}
```

Esa ruta se puede guardar en `imagen` al crear o actualizar una serie.

## Exportar CSV

```bash
curl -o series_chentoons.csv http://localhost:8080/export/series.csv
```

El CSV se puede abrir manualmente con Excel o usarlo desde el frontend como descarga.

## Seeds

Al iniciar, el backend asegura datos iniciales relacionados con Pocoyo, Escandalosos, Snoopy, Bluey, Doraemon, Gravity Falls y otras caricaturas. El seed es idempotente: reutiliza series existentes por nombre y solo completa personajes, episodios y ratings que falten.

## Render

El repo incluye `render.yaml` para desplegar el backend como servicio Docker, crear PostgreSQL en Render y montar un disco persistente para `uploads`.

Nota: Render Persistent Disk requiere un servicio web de pago. Por eso `render.yaml` usa `plan: starter` para el backend. Si usas plan gratis, los uploads funcionaran, pero se perderan al redeploy porque el filesystem es efimero.

Pasos sugeridos:

1. Sube este repo backend a GitHub.
2. En Render, crea un Blueprint desde el repo.
3. Render leera `render.yaml`.
4. Verifica que `DATABASE_URL` quede conectado a `chentoons-postgres`.
5. Verifica que el disco `chentoons-uploads` este montado en `/app/uploads`.
6. Configura `CORS_ORIGIN` con la URL de tu frontend publicado, o deja `*` durante pruebas.
7. Cuando despliegue, prueba `/health`, `/series`, `/uploads/archivo.jpg` y `/docs`.

El backend tambien acepta `DATABASE_URL` directamente, por si prefieres configurar un Web Service y una base Postgres manualmente.

## Reflection breve

Este backend busca ser claro antes que enorme. Las migraciones estan en Go para no depender de herramientas extra, los handlers son directos y la estructura separa lo justo: config, database, models, handlers y routes. Para un video de explicacion, lo mas importante es mostrar Docker Compose, los endpoints CRUD, la carpeta `uploads`, Swagger y el CSV.

## Screenshots pendientes

- Captura de `/docs`
- Captura de respuesta de `/series`
- Captura de subida de imagen
- Captura del CSV abierto en Excel
