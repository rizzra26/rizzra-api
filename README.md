# Rizzra API

Backend API for the [rizzra](https://rizzra.com) portfolio, built with Go and PostgreSQL.

## Stack

- **Go 1.26** — runtime
- [Fiber v3](https://docs.gofiber.io/) — HTTP framework
- [pgx v5](https://github.com/jackc/pgx) — PostgreSQL driver
- [golang-jwt v5](https://github.com/golang-jwt/jwt) — JWT auth
- [caarlos0/env](https://github.com/caarlos0/env) — env config

## Project Structure

```
api/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/config.go        # Env-based configuration
│   ├── database/
│   │   ├── postgres.go         # Connection, migrations, seeding
│   │   └── migrations/         # SQL migration files (embedded)
│   ├── handlers/               # HTTP handlers per domain
│   ├── middleware/
│   │   ├── auth.go             # JWT Bearer token validation
│   │   └── cors.go             # CORS settings
│   ├── models/                 # Domain structs (User, Letter, Project, StackCategory, StackItem)
│   ├── repository/             # Data access layer (SQL queries)
│   ├── router/router.go        # Route definitions
│   └── util/                   # JWT helpers, password hashing, response helpers
└── uploads/                    # Uploaded files (covers)
```

## Getting Started

### Prerequisites

- Go 1.26+
- PostgreSQL 15+

### Setup

```bash
cp .env.example .env
# Fill in .env with your database and JWT secret configuration
```

### Environment Variables

| Variable      | Default                  | Description            |
|-------------- |--------------------------|------------------------|
| `DB_HOST`     | `localhost`              | PostgreSQL host        |
| `DB_PORT`     | `5432`                   | PostgreSQL port        |
| `DB_USER`     | `postgres`               | Database user          |
| `DB_PASSWORD` | `postgrespw!`            | Database password      |
| `DB_NAME`     | `rizzra_dev`             | Database name          |
| `JWT_SECRET`  | `change-me-in-production`| JWT signing secret     |
| `PORT`        | `8888`                   | Server port            |
| `UPLOAD_DIR`  | `./uploads`              | Upload directory       |

### Development

```bash
# Run with Go
go run ./cmd/server

# or with air (hot reload)
air
```

The server starts at `http://localhost:8888`. Database migrations and admin user seeding run automatically on startup.

Default admin credentials:
- Email: `admin@rizzra.dev`
- Password: `admin123`

## API Endpoints

### Public (no authentication required)

| Method   | Endpoint                     | Description                              |
|----------|------------------------------|------------------------------------------|
| `POST`   | `/api/v1/auth/login`         | Login, returns access + refresh tokens   |
| `POST`   | `/api/v1/auth/refresh`       | Refresh access token                     |
| `GET`    | `/api/v1/letters`            | List all letters (paginated)             |
| `GET`    | `/api/v1/letters/:id`        | Get letter detail                        |
| `GET`    | `/api/v1/projects`           | List all projects (paginated)            |
| `GET`    | `/api/v1/projects/:id`       | Get project detail                       |
| `GET`    | `/api/v1/stack/categories`   | List stack categories with their items   |

### Protected (Bearer token required)

All endpoints below require the header: `Authorization: Bearer <access_token>`

| Method  | Endpoint                              | Description                |
|---------|---------------------------------------|----------------------------|
| `GET`   | `/api/v1/dashboard/stats`             | Admin dashboard stats      |
| `POST`  | `/api/v1/letters`                     | Create a new letter        |
| `PUT`   | `/api/v1/letters/:id`                 | Update a letter            |
| `DELETE`| `/api/v1/letters/:id`                 | Delete a letter (soft)     |
| `POST`  | `/api/v1/projects`                    | Create a new project       |
| `PUT`   | `/api/v1/projects/:id`                | Update a project           |
| `DELETE`| `/api/v1/projects/:id`                | Delete a project (soft)    |
| `POST`  | `/api/v1/projects/reorder`            | Reorder projects           |
| `POST`  | `/api/v1/stack/categories`            | Create a stack category    |
| `PUT`   | `/api/v1/stack/categories/:id`        | Update a stack category    |
| `DELETE`| `/api/v1/stack/categories/:id`        | Delete a category (soft)   |
| `POST`  | `/api/v1/stack/items`                 | Create a stack item        |
| `PUT`   | `/api/v1/stack/items/:id`             | Update a stack item        |
| `DELETE`| `/api/v1/stack/items/:id`             | Delete a stack item (soft) |
| `POST`  | `/api/v1/upload/cover`                | Upload project cover (multipart) |

### Auth

**Login**

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "admin@rizzra.dev",
  "password": "admin123"
}
```

Response:

```json
{
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "expires_in": 1710000000,
    "user": {
      "id": "uuid",
      "email": "admin@rizzra.dev",
      "username": "rizzra",
      "role": "admin"
    }
  }
}
```

Access token expires in 1 hour, refresh token in 7 days.

### Pagination

`GET /letters` and `GET /projects` support these query parameters:

| Param      | Default | Description        |
|------------|---------|--------------------|
| `page`     | 1       | Current page       |
| `per_page` | 20      | Items per page     |

The response includes `meta`:

```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 42,
    "total_pages": 3
  }
}
```

### Standard Response Format

All responses follow this structure:

```json
{
  "data": {},
  "meta": {},
  "message": "",
  "error": ""
}
```

---

## Related

- `../web` — Next.js frontend (portfolio website)
- `../admin` — Nuxt.js admin panel
