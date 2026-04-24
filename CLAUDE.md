# Application Overview

This application is Web API server made by golang and Echo. Other tools used include Docker, Air, and OpenAPI.

## Development Commands

### Docker (primary workflow)
```bash
make up              # Start all containers (app + postgres)
make down            # Stop containers
make logs            # Tail container logs
make build           # Rebuild docker images
make clean           # Remove volumes and containers
```

### Database Migrations
```bash
make migrate-up                          # Apply all pending migrations
make migrate-down                        # Roll back last migration
make migrate-create name=<description>   # Create a new migration pair (up/down)
```

### Running Directly (without Docker)
```bash
go run ./cmd/server/main.go
```

## Important Notes

- When developing an API, please write the API specification in the `openapi` directory first and get it reviewed. After that, please proceed with the Go implementation.
- If an ID exists, please implement it in UUID format.
- Please format the timestamp as “2024-01-01T18:00:00+09:00”.
