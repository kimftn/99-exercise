# Property API Boilerplate

Simple Go + Fiber starter for three service areas:

- Listing Services
- User Services
- Public API Layer Services

## Architecture

This project uses a layered structure:

`handler -> service -> repository`

For the public listing APIs, the flow is:

`handler -> service -> http client -> external listing service`

DTOs are used between the HTTP layer and service layer to keep request and response contracts separate from the domain model.

### Layer Responsibilities

- `handler`
  Receives HTTP requests, parses request bodies and params, and returns HTTP responses.
- `service`
  Contains business flow and coordinates between handlers, repositories, and DTO mapping.
- `repository`
  Handles data access for local services.
- `client`
  Handles outbound HTTP calls to external services.
- `dto`
  Defines request and response payloads for each service.

### Current Project Structure

```text
cmd/api/main.go
cmd/migrate/main.go
db/migrations/
internal/app/router.go
internal/client/
internal/domain/
internal/dto/
internal/http/handlers/
internal/http/helpers/
internal/repository/
internal/service/
```

### Request Flow Example

For local `POST /listings`:

1. `ListingHandler` parses the incoming JSON into `dto.CreateListingRequest`.
2. `ListingService` applies business rules and maps the DTO into a domain model.
3. `ListingRepository` saves the data into PostgreSQL.
4. `ListingService` maps the saved model into `dto.ListingResponse`.
5. `ListingHandler` returns the JSON response to the client.

For `GET /public-api/listings`:

1. `PublicHandler` parses query params.
2. `PublicService` delegates to a listing HTTP client.
3. `ListingClient` calls the external Listing Service, for example `http://localhost:6000/listings`.
4. `PublicHandler` returns the external listing data in the public API response.

The same pattern is used for:

- Listing Services
- User Services
- Public API Layer Services

## Run

```bash
go mod tidy
go run ./cmd/api
```

Server starts on `:3000`.

The app automatically loads `.env` from the project root if the file exists.

## PostgreSQL Connection

This project now includes a PostgreSQL connection package using `pgxpool`.

Files:

- `internal/database/postgres/config.go`
- `internal/database/postgres/postgres.go`

Supported environment variables:

- `DATABASE_URL`
- `PGHOST`
- `PGPORT`
- `PGUSER`
- `PGPASSWORD`
- `PGDATABASE`
- `PGSSLMODE`
- `LISTING_SERVICE_BASE_URL`

Example:

```bash
.env

PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGPASSWORD=postgres
PGDATABASE=property_api
PGSSLMODE=disable
LISTING_SERVICE_BASE_URL=http://localhost:6000
```

Then run:

```bash
go run ./cmd/api
```

The app opens a PostgreSQL connection pool at startup and also uses `LISTING_SERVICE_BASE_URL` for the public listing APIs.

## Migrations

Migration files now live in:

- `db/migrations/001_create_users_and_listings.up.sql`
- `db/migrations/001_create_users_and_listings.down.sql`
- `cmd/migrate/main.go`
- `internal/database/migration/migration.go`

The initial migration creates:

- `users`
  Columns: `id`, `name`, `created_at`, `updated_at`
- `listings`
  Columns: `id`, `user_id`, `price`, `listing_type`, `created_at`, `updated_at`
- `listing_type` enum
  Values: `rent`, `sale`

Notes:

- `created_at` and `updated_at` use `BIGINT` because microsecond timestamps will overflow a normal PostgreSQL `INT`.
- `id` and `user_id` also use `BIGINT` to stay consistent with PostgreSQL auto-increment keys.
- `listings.user_id` has a foreign key to `users.id`.

Run migrations with the built-in command:

```bash
go run ./cmd/migrate up
```

Roll back one migration:

```bash
go run ./cmd/migrate down
```

Roll back multiple migration steps:

```bash
go run ./cmd/migrate down 2
```

Check the current version:

```bash
go run ./cmd/migrate version
```

The migration command also loads `.env` from the project root automatically before connecting to PostgreSQL.

## Endpoints

### Listing Services

- `POST /listings`
- `GET /listings`

### User Services

- `GET /users`
- `GET /users/:id`
- `POST /users`

### Public API Layer Services

- `GET /public-api/listings`
- `POST /public-api/users`
- `POST /public-api/listings`

## Sample Request

```bash
curl -X POST http://localhost:3000/listings \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Modern Loft",
    "city":"Surabaya",
    "price":7800000,
    "category":"rent"
  }'
```
