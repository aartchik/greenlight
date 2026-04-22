# Greenlight

`Greenlight` is a REST API for managing a movie catalog. The project is built in Go and focuses on a clean application structure, predictable HTTP behavior, and a deployment flow that is easy to run locally with Docker and PostgreSQL.

At the current stage, the repository already contains the core API foundation: movie CRUD, filtering and pagination, database migrations, validation, rate limiting, structured logging, and the initial user model. Authentication and authorization are planned as the next iteration and the codebase is already prepared for that expansion.

## What This Project Does

- Stores and manages a collection of movies in PostgreSQL.
- Exposes JSON endpoints for creating, reading, updating, deleting, and listing movies.
- Supports filtering by title and genres.
- Supports pagination and sorting for list endpoints.
- Uses optimistic concurrency for movie updates through record versioning.
- Includes a user table and password hashing foundation for upcoming auth flows.

## Stack

- Go
- PostgreSQL
- Docker Compose
- `httprouter`
- `alice`
- `bcrypt`

## Current API

### Health and utility

- `GET /ping`
- `GET /v1/healthcheck`

### Movies

- `GET /v1/movies`
- `POST /v1/movies`
- `GET /v1/movies/:id`
- `PATCH /v1/movies/:id`
- `DELETE /v1/movies/:id`

Example list query:

```text
/v1/movies?title=alien&genres=sci-fi,horror&page=1&page_size=20&sort=-year
```

Supported sorting values:

- `id`
- `-id`
- `title`
- `-title`
- `year`
- `-year`
- `runtime`
- `-runtime`

## Project Layout

```text
cmd/api/           HTTP server, routes, handlers, middleware
internal/data/     Models, queries, filters, domain types
internal/jsonlog/  Structured logging helpers
internal/validator/ Input validation helpers
migrations/        PostgreSQL schema migrations
docker-compose.yml Local app + database startup
dockerfile         Production-style container build
```

## Running Locally

### Option 1: Docker Compose

```bash
docker compose up --build
```

This starts:

- PostgreSQL
- database migrations
- the API on `http://localhost:4000`

### Option 2: Run the API directly

Start PostgreSQL separately, then run:

```bash
go run ./cmd/api -db-dsn="postgres://greenlight_user:1234@localhost:5432/greenlight?sslmode=disable"
```

Default server port:

```text
:4000
```

## Example Requests

Create a movie:

```bash
curl -X POST http://localhost:4000/v1/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Blade Runner",
    "year": 1982,
    "runtime": 117,
    "genres": ["sci-fi", "thriller"]
  }'
```

List movies:

```bash
curl "http://localhost:4000/v1/movies?sort=-year&page=1&page_size=10"
```

Check health:

```bash
curl http://localhost:4000/v1/healthcheck
```

## Roadmap

- Authentication via user registration and login
- Authorization and protected routes
- Token-based access flow
- Better operational tooling and developer commands
- Expanded tests and CI pipeline

## Status

The project is in active development. The base API and data layer are already in place, and the next major step is completing the authentication and authorization layer around the existing user model.
