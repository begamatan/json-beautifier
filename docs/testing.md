# Testing Strategy

## Test Pyramid

```
         /\
        /E2E\        ← 2 Playwright tests (happy path, invalid JSON)
       /------\
      /  Integ \     ← Handler tests (httptest.NewRecorder)
     /----------\
    /  Unit      \   ← Formatter package (table-driven)
   /--------------\
  /  Frontend Unit \  ← Vitest + vue-test-utils (mocked API)
 /------------------\
```

## Backend Unit Tests (`backend/internal/formatter`)

- Pure functions, no I/O.
- Table-driven tests covering: valid inputs, invalid JSON, edge cases (empty input, wrong indent).
- Run with: `go test ./internal/formatter/...`

## Backend Handler / Integration Tests (`backend/internal/handler`)

- Use `net/http/httptest` – no real TCP socket needed.
- Each endpoint tested for: 200 success, 422 invalid JSON, 400 bad request, 413 oversized body.
- Run with: `go test ./internal/handler/...`

## Frontend Unit Tests (`frontend/src`)

- Vitest + `@vue/test-utils`.
- `api/index.ts` is fully mocked with `vi.mock('../../api')`.
- Tests cover: rendering, button actions, success/error messages, copy/download disabled state.
- Run with: `npm run test:unit`

## E2E Tests (`frontend/e2e`)

- Playwright with Chromium.
- API routes are intercepted with `page.route()` so no running backend is required.
- Tests: happy-path beautify, invalid JSON error display.
- Run with: `npm run test:e2e`

## TDD Workflow

1. Write failing test that describes the desired behaviour.
2. Write the minimal implementation to make the test pass.
3. Refactor while keeping tests green.
4. Repeat for the next behaviour.
