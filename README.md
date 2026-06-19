# Chirpy Server

Chirpy Server is a Go-based API for managing chirps, users, authentication, and premium upgrade webhooks. The server serves the frontend static assets from the `static/` directory and exposes JSON endpoints for creating, reading, updating, and deleting chirps.

## Features

- User registration and profile updates
- Login flow with access and refresh tokens
- Chirp creation, listing, fetching, and deletion
- Admin metrics endpoint for counting static file requests
- Optional database reset endpoint for development environments
- Polka webhook support for premium user events

## Requirements

- Go 1.26+
- PostgreSQL database
- `sqlc` (for generating query code)
- A `.env` file or exported environment variables

## Installation

1. Clone the repository:

   ```bash
   git clone <your-repo-url>
   cd bootdev-chirpy-go
   ```

2. Install Go dependencies:

   ```bash
   go mod tidy
   ```

3. Install `sqlc` if it is not already available:

   ```bash
   go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
   ```

4. Create a PostgreSQL database and update your environment variables.

## Environment Variables

The application expects the following variables to be set before running:

```bash
DB_URL=postgres://username:password@localhost:5432/chirpy
JWT_SECRET=your-secret-key
PLATFORM=dev
POLKA_KEY=your-polka-webhook-key
```

- `DB_URL`: connection string for your PostgreSQL database
- `JWT_SECRET`: secret key used for signing JWTs
- `PLATFORM`: set to `dev` to allow the admin reset route
- `POLKA_KEY`: API key required for the webhook endpoint

## Running the Project

### Option 1: Run directly

```bash
go run .
```

The server listens on port `8080` by default.

### Option 2: Use the helper script

```bash
./build-and-serve.sh
```

This script generates SQL code, builds the binary, and starts the server.

## API Overview

### Health

- `GET /api/healthz` — returns a simple health check response

### Chirps

- `POST /api/chirps` — create a chirp (requires a valid bearer token)
- `GET /api/chirps` — list all chirps or filter by `author_id`
- `GET /api/chirps/{chirpId}` — fetch a single chirp
- `DELETE /api/chirps/{chirpId}` — delete a chirp if you own it

### Users

- `POST /api/users` — create a new user account
- `PUT /api/users` — update a user profile (requires a valid bearer token)

### Auth

- `POST /api/login` — log in and receive tokens
- `POST /api/refresh` — refresh an access token with a refresh token
- `POST /api/revoke` — revoke a refresh token

### Webhooks

- `POST /api/polka/webhooks` — updates premium status for a user

### Admin

- `GET /admin/metrics` — displays a simple metrics page for static asset requests
- `POST /admin/reset` — resets the database and hit count (only works when `PLATFORM=dev`)

## Development Notes

- SQL query definitions live in the `sql/queries` folder.
- Database schema files are stored in `sql/schema`.
- Generated Go database code is output to `internal/database`.
- After changing SQL files, run `sqlc generate` to refresh the generated code.

## Example Workflow

1. Start the database and set your environment variables.
2. Run `sqlc generate` if needed.
3. Start the server with `go run .`.
4. Use tools like `curl` or Postman to test the endpoints.
