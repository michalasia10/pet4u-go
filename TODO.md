## ToDo (pet4u-go)

A concise, living checklist to keep backend aligned with the frontend (`~/Projects/pet4u`). Update continuously.

### Conventions and naming
- [x] Feature-oriented DDD layout: `internal/modules/<feature>/{domain,application,infrastructure,interfaces/http}`
- [x] Domain ports in `domain/` only; adapters in `infrastructure/` and `interfaces/http/`
- [x] DTOs and mappers live in `interfaces/http/` only; no JSON tags in domain



---

### Module: Auth (missing on backend; frontend is mocked)
Endpoints to expose as used by the frontend (`features/auth/infrastructure/api/authApi.ts`):
- [ ] POST `/api/v1/auth/login` — req: `{ email, pass }` → res: `{ user, token }`
- [ ] POST `/api/v1/auth/refresh` — req: `{ refreshToken }` → res: `{ accessToken, refreshToken }`
- [ ] POST `/api/v1/auth/logout` — 204
- [ ] GET `/api/v1/auth/me` — res: `user`

Layers and files:
- [ ] `internal/modules/auth/domain/entities/user.go`
- [ ] `internal/modules/auth/domain/valueobjects/{token.go,credentials.go}`
- [ ] `internal/modules/auth/domain/repositories/user_repository.go` (port)
- [ ] `internal/modules/auth/domain/services/auth_service.go` (port: credentials verification)
- [ ] `internal/modules/auth/application/{login.go,refresh_token.go,get_me.go,logout.go}`
- [ ] `internal/modules/auth/infrastructure/memory/{user_repository.go,token_provider.go}`
- [ ] `internal/modules/auth/interfaces/http/{dto.go,mapper.go,router.go}`
- [ ] `internal/server/routes.go` — mount: `r.Mount("/auth", authhttp.NewRouter())`

DTO shapes (match frontend):
- [ ] `UserDTO`: `{ id: string, email: string, firstName: string, lastName: string }`
- [ ] `AuthTokenDTO`: `{ accessToken: string, refreshToken: string }`

Security (MVP):
- [ ] Simple in-memory token provider; later JWT + persistent user repo
- [ ] Bearer verification middleware (sufficient for `/auth/me` in MVP)

---

### Module: Places (extend to match frontend)
Endpoints used on frontend (`features/places/infrastructure/api/placeApi.ts`):
- [ ] GET `/api/v1/places` — support filters: `q`, `tag`, `species`
- [ ] GET `/api/v1/places/{id}` — place details
- [ ] POST `/api/v1/places` — create place (body: `AddPlaceDTO`)

Layers and files:
- [ ] `domain/repository.go` — add `GetByID(id string)` and `Create(p Place)`
- [ ] `application/list_places.go` — params: `Query`, `Tags`, `Species`
- [ ] `application/get_place.go`
- [ ] `application/create_place.go`
- [ ] `infrastructure/memory/place_repository.go` — method impls + seed
- [ ] `interfaces/http/dto.go` — `PlaceDTO`, `PlaceDetailsDTO`, `CreatePlaceRequestDTO`
- [ ] `interfaces/http/mapper.go`
- [ ] `interfaces/http/router.go` — add GET `/`, GET `/{id}`, POST `/`

Domain notes:
- [ ] Consider VO `Location` in `domain` (lat/lng), keep domain free of JSON tags
- [ ] Details DTO may include optional fields (gallery, reviews) — start MVP with a minimal set

---

### Module: Pets (optional extensions)
- [ ] (Optional) List `/api/v1/pets` and search — if frontend requires it
- [ ] (Optional) Update/archive

### Module: Appointments (optional extensions)
- [ ] (Optional) Get single appointment `GET /api/v1/appointments/{id}`
- [ ] (Optional) Cancel/Reschedule — per domain rules

---

### Tests
- [ ] Unit: use cases (`application/`) with in-memory repositories
- [ ] Integration: `interfaces/http` + `internal/server` (chi router) — happy path and 422 validations
- [ ] (Optional) `testcontainers` once a DB appears

### DX
- [ ] Update `README.md` with new endpoints and curl examples
- [ ] Provide example frontend `.env`: `VITE_API_URL=http://localhost:8080/api/v1`
- [ ] Verify CORS (credentials, headers) — already globally configured, double-check

---

### Milestones
1) MVP compatibility with frontend (no JWT):
   - [ ] Auth (login/refresh/me/logout) + Places (list/get/create) + routing
2) Contract stabilization and DTO validations (422):
   - [ ] Align with `react-hook-form` and Zod on the frontend
3) Security:
   - [ ] JWT, refresh rotation, middleware, tests

---

### Open questions
- [ ] Should `GET /places/{id}` return full details (gallery, reviews, contact)? — MVP: yes, but optional fields
- [ ] Should `POST /places` set `averageRating` (0) and `createdAt` immediately? — MVP: yes
- [ ] Target storage for `auth` and `places` (Postgres?) and migrations — after MVP


