# Docker & Deployment Notes

## Images

| Image     | Base           | Strategy               |
|-----------|----------------|------------------------|
| backend   | `scratch`      | Multi-stage; Go binary only (~6 MB) |
| frontend  | `nginx:alpine` | Multi-stage; built static files     |

## Local Docker Compose

```bash
docker compose up --build
```

- Frontend: http://localhost:80
- Backend: http://localhost:8080/api/v1/health

## Environment Variables

| Variable               | Default   | Description                        |
|------------------------|-----------|------------------------------------|
| `PORT`                 | `8080`    | Backend listen port                |
| `CORS_ALLOWED_ORIGINS` | `*`       | Comma-separated allowed origins    |
| `VITE_API_BASE_URL`    | `/api/v1` | Frontend API base URL              |

## Production Considerations

- Replace `CORS_ALLOWED_ORIGINS=*` with the actual frontend origin.
- Run behind a reverse proxy (e.g. Nginx, Caddy) with TLS termination.
- Set resource limits in Docker / Kubernetes (not included in MVP).
- The Go binary is built with `-trimpath -ldflags="-s -w"` for minimal size.

## Known Limitations & Future Improvements

- No authentication or rate limiting (MVP).
- No Kubernetes manifests (out of scope).
- Indent option limited to 2 or 4 spaces (by design; easily extendable).
- Frontend does not persist state across page reloads (no localStorage).
- No dark mode support.
- Clipboard API requires HTTPS in production browsers.
