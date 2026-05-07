# Property API Boilerplate

Simple Go + Fiber starter for three service areas:

- Listing Services
- User Services
- Public API Layer Services

## Architecture

This project uses a layered structure:

`handler -> service -> repository -> store`

DTOs are used between the HTTP layer and service layer to keep request and response contracts separate from the domain model.

### Layer Responsibilities

- `handler`
  Receives HTTP requests, parses request bodies and params, and returns HTTP responses.
- `service`
  Contains business flow and coordinates between handlers, repositories, and DTO mapping.
- `repository`
  Handles data access. In this boilerplate it uses an in-memory store, but it can later be replaced with a database implementation.
- `dto`
  Defines request and response payloads for each service.
- `store`
  Holds the dummy in-memory data used by the repositories.

### Current Project Structure

```text
cmd/api/main.go
internal/app/router.go
internal/domain/
internal/dto/
internal/http/handlers/
internal/repository/
internal/service/
internal/store/
```

### Request Flow Example

For `POST /listings`:

1. `ListingHandler` parses the incoming JSON into `dto.CreateListingRequest`.
2. `ListingService` applies business rules and maps the DTO into a domain model.
3. `ListingRepository` saves the data into the in-memory `store`.
4. `ListingService` maps the saved model into `dto.ListingResponse`.
5. `ListingHandler` returns the JSON response to the client.

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
