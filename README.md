# pokedex-api

A Go-based API that provides information about Pokemon, with support for translations.

## Prerequisites

- If you run it locally, [Go](https://go.dev/) 1.24 or higher 
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/) for containerized execution

## Getting Started

### 1. Environment Variables

The application requires the following environment variables:

- `PORT`: The port the server will listen on (default: `8080`).
- `POKEMON_API_URL`: The base URL for the PokeAPI (e.g., `https://pokeapi.co`).

### 2. Running Locally

To run the application directly on your machine:

```bash
export POKEMON_API_URL=https://pokeapi.co
go run ./cmd/api
```

Alternatively, using the `Makefile`:

```bash
make build
./bin/pokedex-api
```

### 3. Running with Docker Compose (and WireMock)

To run the application along with a mocked PokeAPI (WireMock) for testing:

```bash
docker-compose up --build
```

- The API will be available at `http://localhost:8080`.
- WireMock admin/mock interface will be at `http://localhost:8081`.

### 4. Interacting with the API

You can reach the API by making http requests.
For example:

```bash
curl http://localhost:8080/api/pokemon/mewtwo
```

## Testing

### Unit Tests

Run the standard test suite:

```bash
make test
```

### E2E Tests

The project includes in-process E2E tests that spin up the application and a mock PokeAPI server:

```bash
make test-e2e
```

## API Endpoints

- `GET /health`: Health check endpoint.
- `GET /api/pokemon/{name}`: Get basic information about a Pokemon.
- `GET /api/pokemon/translated/{name}`: Get information about a Pokemon with translated descriptions.
