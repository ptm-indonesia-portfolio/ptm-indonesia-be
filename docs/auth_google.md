## Google SSO Auth

### Login URL

`GET /api/v1/auth/google/login`

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location 'http://localhost:9100/api/v1/auth/google/login' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: id'
```

### Success Response

```json
{
  "message": "URL login Google berhasil dibuat.",
  "data": {
    "url": "https://accounts.google.com/o/oauth2/v2/auth?..."
  }
}
```

### Callback

`GET /api/v1/auth/google/callback`

Endpoint ini dipanggil oleh Google setelah user login. Backend akan:

- validasi state
- verifikasi akun Google
- cek email sudah terdaftar di tabel `users`
- set access token cookie `HttpOnly`
- set refresh token cookie `HttpOnly`
- set cookie `logged_in`
- redirect ke `FRONTEND_URL` dengan query `auth_status`

### Current Session

`GET /api/v1/auth/me`

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location 'http://localhost:9100/api/v1/auth/me' \
  --header 'Accept: application/json'
```

### Logout

`POST /api/v1/auth/logout`

Endpoint ini akan menghapus access token cookie, refresh token cookie, dan cookie `logged_in`.

### Headers

- `Accept: application/json`
- `Content-Type: application/json`

### Curl

```bash
curl --location --request POST 'http://localhost:9100/api/v1/auth/logout' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json'
```
