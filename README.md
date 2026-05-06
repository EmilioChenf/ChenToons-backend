# ChenToons Backend

ChenToons es una API REST sencilla para llevar un tracker de series y caricaturas como Pocoyo, Escandalosos, Snoopy, Bluey y otras parecidas. El proyecto esta hecho con Go, Fiber, PostgreSQL y Docker, pensando en que sea facil de correr, explicar y conectar con un frontend universitario.

## Tecnologias

- Go 1.22
- Fiber
- PostgreSQL
- Docker y Docker Compose
- OpenAPI / Swagger UI
- CSV para exportar series y abrirlas en Excel

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

Al iniciar, si la tabla `series` esta vacia, el backend crea datos iniciales relacionados con Pocoyo, Escandalosos, Snoopy, Bluey, Doraemon, Gravity Falls y otras caricaturas. Tambien crea algunos personajes, episodios y ratings.

## Reflection breve

Este backend busca ser claro antes que enorme. Las migraciones estan en Go para no depender de herramientas extra, los handlers son directos y la estructura separa lo justo: config, database, models, handlers y routes. Para un video de explicacion, lo mas importante es mostrar Docker Compose, los endpoints CRUD, la carpeta `uploads`, Swagger y el CSV.

## Screenshots pendientes

- Captura de `/docs`
- Captura de respuesta de `/series`
- Captura de subida de imagen
- Captura del CSV abierto en Excel
