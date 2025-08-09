# pet4u-go

Feature-oriented DDD Lite monolith for Pet4U: find pet-friendly places, book specialist appointments, keep pet medical info (patient card).

Architecture follows Clean Architecture (Ports & Adapters) with DDD per feature to maximize developer experience.

## Architecture

- Features in `internal/modules/<feature>`:
  - `domain/` – Entities, value objects, domain ports (interfaces), domain errors
  - `application/` – Use cases (orchestrate domain). No framework imports
  - `infrastructure/` – Adapters (e.g., Postgres, in-memory)
  - `interfaces/http/` – HTTP handlers, DTOs, mappers, router
- Cross-cutting utilities in `internal/pkg/*` (e.g., `httpx`)
- Composition in `internal/server/*`

References: [DDD Lite Intro](https://threedots.tech/post/ddd-lite-in-go-introduction/), [Avoid DRY in Go](https://threedots.tech/post/things-to-know-about-dry/), [Clean Architecture](https://threedots.tech/post/introducing-clean-architecture/)

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

## API Sketch

- `GET /api/v1/places` – list demo places
- `GET /api/v1/places/search?q=...&tag=park&tag=cafe` – search
- `POST /api/v1/appointments` – create booking
- `GET /api/v1/appointments` – list bookings
- `POST /api/v1/pets` – create pet
- `GET /api/v1/pets/{id}` – get pet

## Notes

- Domain stays pure; all IO in adapters. Prefer small explicit mappers over shared structs to avoid API-storage coupling, per [When to avoid DRY in Go](https://threedots.tech/post/things-to-know-about-dry/).

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
