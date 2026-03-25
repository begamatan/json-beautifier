# Architecture Overview

## System Diagram

```
Browser (Vue 3 SPA)
      │
      │  HTTP/JSON
      ▼
Nginx (port 80)
      │
      │  /api/*  → proxy
      ▼
Go HTTP Server (port 8080)
      │
      ├── POST /api/v1/beautify  → formatter.Beautify()
      ├── POST /api/v1/minify    → formatter.Minify()
      ├── POST /api/v1/validate  → formatter.Validate()
      └── GET  /api/v1/health
```

## Components

| Layer      | Technology               | Responsibility                     |
|------------|--------------------------|------------------------------------|
| Frontend   | Vue 3 + TypeScript + Vite| UI, user input, display output     |
| Proxy      | Nginx                    | Serve SPA, proxy `/api/*` to backend|
| Backend    | Go std-lib HTTP server   | Business logic, JSON processing    |

## Data Flow

1. User pastes JSON into the textarea and clicks an action button.
2. Vue calls the corresponding `api/index.ts` function which sends a `POST` request.
3. The Go handler reads the body (max 5 MB), calls the formatter package, and returns a JSON response.
4. Vue updates the output textarea or shows an error banner.

## Key Design Decisions

- **No external Go router**: the standard library `net/http` `ServeMux` (with method patterns) is sufficient and avoids dependencies.
- **CORS via `github.com/rs/cors`**: simple, well-tested, single-purpose library.
- **Formatter as a pure package**: no state, fully table-driven testable, zero side-effects.
- **Frontend API layer**: thin wrapper (`src/api/index.ts`) that can be easily swapped for a mock in tests.
