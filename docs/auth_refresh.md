## Refresh Session Auth

### Refresh Session

`POST /api/v1/auth/refresh`

Endpoint ini digunakan untuk memperbarui access token ketika refresh token masih valid. Request tidak perlu body. Cookie auth akan dibaca otomatis oleh browser.

### Headers

- `Accept: application/json`
- `Content-Type: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location --request POST 'http://localhost:9100/api/v1/auth/refresh' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --header 'Accept-Language: id'
```

### Success Response

```json
{
  "message": "Session berhasil diperbarui.",
  "data": {
    "id": 1,
    "name": "Super Admin",
    "email": "admin@example.com",
    "status": 1
  }
}
```
