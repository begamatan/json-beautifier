# API Contract

Base URL: `http://localhost:8080/api/v1`

## GET /api/v1/health

**Response 200**
```json
{ "status": "ok" }
```

---

## POST /api/v1/beautify

**Request**
```json
{
  "json": "{\"name\":\"Alice\"}",
  "indent": 2
}
```

| Field    | Type    | Required | Notes                     |
|----------|---------|----------|---------------------------|
| `json`   | string  | yes      | Raw JSON string           |
| `indent` | integer | no       | `2` or `4`; defaults to `2` |

**Response 200**
```json
{
  "result": "{\n  \"name\": \"Alice\"\n}"
}
```

**Response 422 – invalid JSON**
```json
{
  "code": "INVALID_JSON",
  "message": "input is not valid JSON",
  "details": ""
}
```

**Response 400 – bad indent value**
```json
{
  "code": "BAD_REQUEST",
  "message": "indent must be 2 or 4, got 3",
  "details": ""
}
```

---

## POST /api/v1/minify

**Request**
```json
{
  "json": "{\n  \"name\": \"Alice\"\n}"
}
```

**Response 200**
```json
{
  "result": "{\"name\":\"Alice\"}"
}
```

---

## POST /api/v1/validate

**Request**
```json
{
  "json": "{\"name\":\"Alice\"}"
}
```

**Response 200 – valid**
```json
{
  "valid": true,
  "message": "JSON is valid"
}
```

**Response 200 – invalid**
```json
{
  "valid": false,
  "message": "input is not valid JSON"
}
```

---

## Error Schema

All error responses share the same shape:

```json
{
  "code":    "STRING_CODE",
  "message": "Human-readable description",
  "details": "Optional extra context"
}
```

| Code             | HTTP Status | Cause                          |
|------------------|-------------|--------------------------------|
| `BAD_REQUEST`    | 400         | Malformed body or invalid param|
| `INVALID_JSON`   | 422         | Input is not parseable JSON    |
| `INTERNAL_ERROR` | 500         | Unexpected server error        |
