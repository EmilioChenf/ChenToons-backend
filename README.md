# ChenToons Backend

ChenToons Backend es una API REST para administrar series animadas. El proyecto usa Go + Fiber + PostgreSQL, corre con Docker y expone documentacion con Swagger/OpenAPI.

La API maneja series, personajes, episodios, ratings, subida de imagenes y exportacion CSV. Esta pensada para trabajar con un frontend separado hecho en HTML, CSS y JavaScript vanilla.

## Tecnologias usadas

- Go
- Fiber
- PostgreSQL
- Docker
- Swagger/OpenAPI

## Como correr localmente

1. Clonar este repositorio:

```bash
git clone URL_DEL_REPOSITORIO_BACKEND
cd ChenToons-backend
```

2. Levantar backend y base de datos con Docker:

```bash
docker compose up --build
```

3. Probar las URLs principales:

```text
http://localhost:8080/health
http://localhost:8080/series
http://localhost:8080/swagger/index.html
```

Tambien puedes usar:

```text
http://localhost:8080/docs
```

## Variables de entorno

El backend puede conectarse a PostgreSQL usando `DATABASE_URL` o variables separadas.

Variables principales:

- `DATABASE_URL`: URL completa de PostgreSQL, usada normalmente en Render.
- `PORT`: puerto que usa Render.
- `APP_ENV`: ambiente del proyecto, por ejemplo `development` o `production`.
- `UPLOAD_DIR`: carpeta donde se guardan imagenes subidas.
- `CORS_ORIGIN`: origen permitido para el frontend.

Variables locales opcionales si no usas `DATABASE_URL`:

- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`

El archivo `.env.example` trae valores de referencia para desarrollo local.

## Docker

El proyecto incluye `Dockerfile` y `docker-compose.yml`.

Docker Compose levanta:

- `backend`: API de ChenToons en Go + Fiber.
- `postgres`: base de datos PostgreSQL.
- `postgres_data`: volumen para persistir la base de datos.
- `./uploads:/app/uploads`: carpeta local para imagenes subidas.

Comando principal:

```bash
docker compose up --build
```

Para reiniciar todo desde cero, incluyendo la base:

```bash
docker compose down -v
docker compose up --build
```

## Endpoints principales

| Metodo | Endpoint | Descripcion |
|---|---|---|
| GET | `/health` | Revisa si el backend esta activo |
| GET | `/series` | Lista series con filtros, busqueda, paginacion y ordenamiento |
| GET | `/series/:id` | Obtiene una serie por ID |
| POST | `/series` | Crea una serie |
| PUT | `/series/:id` | Edita una serie |
| DELETE | `/series/:id` | Elimina una serie |
| GET | `/personajes?serie_id=ID` | Lista personajes |
| POST | `/personajes` | Crea un personaje |
| GET | `/episodios?serie_id=ID` | Lista episodios |
| POST | `/episodios` | Crea un episodio |
| GET | `/series/:id/ratings` | Lista ratings de una serie |
| POST | `/series/:id/ratings` | Agrega rating y comentario |
| GET | `/series/:id/promedio-rating` | Calcula el promedio de rating |
| POST | `/uploads` | Sube una imagen |
| GET | `/uploads/:filename` | Sirve una imagen subida |
| GET | `/export/series.csv` | Exporta series en CSV |

`GET /series` acepta:

- `search` o `q`
- `genero`
- `categoria`
- `estado`
- `sort`
- `order`
- `page`
- `limit`

Ejemplo:

```bash
curl "http://localhost:8080/series?q=snoopy&page=1&limit=8&sort=nombre&order=asc"
```

## Swagger

Swagger UI esta disponible en:

```text
http://localhost:8080/swagger/index.html
```

Tambien funciona:

```text
http://localhost:8080/docs
```

El archivo OpenAPI esta en:

```text
docs/openapi.yaml
```

En Render sera la misma ruta con el dominio publicado:

```text
https://TU-BACKEND-RENDER.onrender.com/swagger/index.html
```

## CORS

El backend tiene CORS configurado para que el frontend pueda consumir la API desde otro puerto o dominio.

En desarrollo se puede usar:

```text
CORS_ORIGIN=*
```

En produccion conviene cambiarlo por la URL real del frontend publicado en Render.

## Challenges implementados

- Swagger/OpenAPI.
- Swagger UI servido desde el backend.
- Validaciones server-side.
- Codigos HTTP correctos.
- Paginacion en `/series`.
- Busqueda por `search` o `q`.
- Ordenamiento con `sort` y `order`.
- Sistema de rating con tabla propia.
- Upload de imagenes para series.
- Exportacion CSV.
- Docker con backend y PostgreSQL.
- Seeds iniciales coherentes e idempotentes.

## Deploy

El backend esta preparado para Render usando Docker y PostgreSQL de Render.

Archivo incluido:

```text
render.yaml
```

URLs del proyecto:

- Backend Render: `https://cheentoons-backend.onrender.com/`
- Frontend Render: `https://cheentoons-frontend.onrender.com/`
- Repositorio frontend: `https://github.com/EmilioChenf/ChenToons-frontend`

Notas para Render:

- Render usa `PORT`, y el backend ya lo soporta.
- Render puede usar `DATABASE_URL` para conectarse a PostgreSQL.
- Sin Persistent Disk, los archivos subidos a `/uploads` pueden perderse al redeploy.
- Las imagenes base del frontend pueden servirse desde `assets/images/`.

## Screenshot

![ChenToons Backend](screenshots/backend.png)

> Pendiente agregar captura real de Swagger o de una respuesta de la API.

## Reflexion

Con este backend practicamos como construir una API REST real usando Go y Fiber. Lo mas interesante fue conectar la API con PostgreSQL, manejar Docker y preparar el proyecto para correr igual en local y en Render. Tambien nos ayudo a entender mejor CORS, codigos HTTP, Swagger y la separacion entre backend y frontend.

Go + Fiber se sintio bastante directo para este tipo de proyecto. Docker costo un poco al inicio, pero al final facilita mucho probar la base de datos y el backend sin configurar todo manualmente. Si hicieramos otro proyecto parecido, si volveriamos a usar esta stack porque es rapida, simple y facil de explicar.

## Autor

Proyecto: ChenToons  
Autor: Emilio Josue Chen Borrayo
