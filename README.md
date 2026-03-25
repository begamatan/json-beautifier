# JSON Beautifier

A production-quality MVP JSON tool that lets you beautify, minify, and validate JSON right in your browser.

---

## Project Overview

| Feature | Description |
|---|---|
| Beautify | Pretty-print JSON with 2 or 4-space indent |
| Minify | Strip all unnecessary whitespace |
| Validate | Check whether input is valid JSON |
| Copy | One-click copy of output to clipboard |
| Download | Save output as `output.json` |
| Size limit | Requests capped at 5 MB |

---

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go 1.24, standard library HTTP server |
| CORS | `github.com/rs/cors` |
| Frontend | Vue 3 + TypeScript + Vite |
| Unit tests (BE) | Go `testing` (table-driven) |
| Integration tests (BE) | Go `net/http/httptest` |
| Unit tests (FE) | Vitest + `@vue/test-utils` |
| E2E tests | Playwright (Chromium) |
| Containerization | Docker + Docker Compose |

---

## Folder Structure

```
json-beautifier/
├── backend/
│   ├── cmd/server/main.go          # HTTP server entry point
│   ├── internal/
│   │   ├── formatter/
│   │   │   ├── formatter.go        # Beautify / Minify / Validate logic
│   │   │   └── formatter_test.go   # Unit tests (table-driven)
│   │   ├── handler/
│   │   │   ├── handler.go          # HTTP handlers
│   │   │   └── handler_test.go     # Handler integration tests
│   │   └── middleware/
│   │       └── middleware.go       # Request size limiter
│   ├── go.mod
│   ├── go.sum
│   └── Dockerfile                  # Multi-stage Go image
├── frontend/
│   ├── src/
│   │   ├── api/index.ts            # API client
│   │   ├── components/
│   │   │   ├── JsonTool.vue        # Main tool component
│   │   │   └── __tests__/
│   │   │       └── JsonTool.spec.ts # Component unit tests
│   │   ├── App.vue
│   │   └── main.ts
│   ├── e2e/
│   │   └── vue.spec.ts             # Playwright E2E tests
│   ├── nginx.conf                  # Nginx SPA + proxy config
│   ├── Dockerfile                  # Multi-stage frontend image
│   ├── package.json
│   ├── vite.config.ts
│   └── playwright.config.ts
├── docs/
│   ├── architecture.md
│   ├── api-contract.md
│   ├── testing.md
│   └── deployment.md
├── docker-compose.yml
├── .gitignore
└── README.md
```

---

## Local Setup

### Prerequisites

- Go 1.24+
- Node 20+
- Docker & Docker Compose (for container flow)

### Backend

```bash
cd backend
go mod download
go run ./cmd/server
# API available at http://localhost:8080
```

### Frontend (dev server)

```bash
cd frontend
npm install
npm run dev
# UI available at http://localhost:5173
```

> Set `VITE_API_BASE_URL=http://localhost:8080/api/v1` (or create a `.env.local` file) so the frontend dev server talks to your local Go server.

---

## Docker Setup

### Production build

```bash
docker compose up --build
```

- Frontend: http://localhost:80
- Backend health: http://localhost:8080/api/v1/health

### Development (hot reload)

```bash
docker compose -f docker-compose.dev.yml up --build
```

- Frontend dev server (Vite): http://localhost:5173
- Backend API: http://localhost:8080/api/v1/health
- Source changes are reflected immediately — no rebuild required.

### Stop

```bash
docker compose down
# or, for the dev stack:
docker compose -f docker-compose.dev.yml down
```

---

## Environment Variables

### Backend

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Listen port |
| `CORS_ALLOWED_ORIGINS` | `*` | Allowed CORS origins |

### Frontend (Vite build-time)

| Variable | Default | Description |
|---|---|---|
| `VITE_API_BASE_URL` | `/api/v1` | API base URL |

### Frontend (Vite dev server)

| Variable | Default | Description |
|---|---|---|
| `BACKEND_URL` | `http://localhost:8080` | Proxy target for `/api` requests |

---

## Run Commands

```bash
# Backend
cd backend && go run ./cmd/server

# Frontend dev
cd frontend && npm run dev

# Frontend production build
cd frontend && npm run build-only

# Docker (all services)
docker compose up --build
```

---

## Test Commands

```bash
# Backend unit + integration tests
cd backend && go test ./...

# Frontend unit tests
cd frontend && npm run test:unit

# Frontend E2E (Playwright — starts dev server automatically)
cd frontend && npm run test:e2e
```

---

## API Quickstart

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Beautify
curl -s -X POST http://localhost:8080/api/v1/beautify \
  -H 'Content-Type: application/json' \
  -d '{"json":"{\"name\":\"Alice\",\"age\":30}","indent":2}'

# Minify
curl -s -X POST http://localhost:8080/api/v1/minify \
  -H 'Content-Type: application/json' \
  -d '{"json":"{\n  \"name\": \"Alice\"\n}"}'

# Validate
curl -s -X POST http://localhost:8080/api/v1/validate \
  -H 'Content-Type: application/json' \
  -d '{"json":"{\"name\":\"Alice\"}"}'
```

---

## Troubleshooting

| Problem | Solution |
|---|---|
| `CORS_ALLOWED_ORIGINS` mismatch | Set the env var to your frontend origin, e.g. `http://localhost:5173` |
| Frontend can't reach API | Ensure `VITE_API_BASE_URL` is set correctly and the backend is running |
| Port 8080 already in use | Change `PORT` env var or update `docker-compose.yml` host port mapping |
| Clipboard copy fails | Clipboard API requires HTTPS or `localhost`; use HTTP only in local dev |
| E2E tests time out | Ensure port 5173 (dev) or 4173 (preview) is free before running Playwright |

---

## Docs

See the [`docs/`](docs/) folder for:

- [Architecture overview](docs/architecture.md)
- [API contract & examples](docs/api-contract.md)
- [Testing strategy](docs/testing.md)
- [Docker & deployment notes](docs/deployment.md)

---

## End-to-End Verification Checklist

- [ ] `cd backend && go test ./...` → all pass
- [ ] `cd frontend && npm run test:unit -- --run` → all pass
- [ ] `cd frontend && npm run test:e2e` → 2 Playwright tests pass
- [ ] `cd frontend && npm run build-only` → build succeeds
- [ ] `docker compose up --build` → both containers start
- [ ] `curl http://localhost:8080/api/v1/health` → `{"status":"ok"}`
- [ ] Beautify, Minify, Validate buttons work in browser
- [ ] Copy button copies output to clipboard
- [ ] Download button downloads `output.json`
- [ ] Oversized payload (>5 MB) returns error
