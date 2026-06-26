## Health Check

### Endpoint

`GET /api/v1/health`

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl Indonesia

```bash
curl --location 'http://localhost:9100/api/v1/health' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: id'
```

### Curl English

```bash
curl --location 'http://localhost:9100/api/v1/health' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: en'
```

### Success Response

```json
{
  "message": "Layanan berjalan dengan baik.",
  "data": {
    "name": "PTM Indonesia API",
    "environment": "development",
    "database": "up",
    "default_language": "id",
    "supported_languages": [
      "id",
      "en"
    ],
    "timestamp": "2026-06-25T10:30:00+07:00"
  }
}
```
