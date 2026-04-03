# URL Shorten API

A high-performance URL shortening service built with Go, Gin, and PostgreSQL.

## Features

- **URL Shortening** - Create short, shareable links with optional custom codes
- **Expiration** - Set expiration times for links (optional)
- **Statistics** - Track click counts and view link analytics
- **Authentication** - JWT-based auth for protected operations
- **Rate Limiting** - Built-in rate limiting to prevent abuse

## Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose

### Running with Docker

```bash
cp .env.example .env
make docker-run-build
```

The server will start at `http://localhost:8080`

### Running Locally

```bash
cp .env.example .env
go build -o main ./cmd/server
./main
```

## API Endpoints

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/auth/register` | Register new user | No |
| POST | `/auth/login` | Login user | No |
| POST | `/shorten` | Create short URL | Optional |
| GET | `/:shortCode` | Redirect to original URL | No |
| GET | `/stats/:shortCode` | Get URL statistics | Yes |
| DELETE | `/:shortCode` | Delete a URL | Yes |

### Create Short URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example.com", "customCode": "my-link"}'
```

### With Authentication

```bash
# Login first
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'

# Use token for authenticated requests
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"url": "https://example.com"}'
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | Database user |
| DB_PASSWORD | postgres | Database password |
| DB_NAME | urlshortener | Database name |
| JWT_SECRET | your-secret-key | JWT signing key |
| SERVER_PORT | 8080 | HTTP server port |
| BASE_URL | http://localhost:8080 | Base URL for short links |
| RATE_LIMIT_RPM | 100 | Requests per minute limit |

## Tech Stack

- **Go** - Programming language
- **Gin** - Web framework
- **GORM** - ORM for PostgreSQL
- **JWT** - Authentication
- **Docker** - Containerization
