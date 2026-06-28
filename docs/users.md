## User CRUD

Seluruh endpoint pada modul ini hanya dapat diakses oleh user dengan status `super_admin`. Cookie auth akan dibaca otomatis oleh browser.

### List Users

`GET /api/v1/users`

Endpoint ini mendukung pagination, global search, dan sorting. Default sorting adalah `id desc`.

List user tidak menampilkan akun admin utama yang dikonfigurasi pada sistem.

### Query Params

- `page` default `1`
- `limit` default `10`
- `search` opsional, global search ke `name`, `email`, `address`, dan `telp`
- `sort_by` opsional: `name`, `email`, atau `status`
- `sort_direction` opsional: `asc` atau `desc`

### Keterangan Status User

- `0` = `not_active`: user nonaktif dan tidak boleh login
- `1` = `super_admin`: user aktif dengan role super admin dan boleh login
- `2` = `free_member`: user aktif member gratis dan boleh login
- `3` = `premium_member`: user aktif member premium dan boleh login

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location 'http://localhost:9100/api/v1/users?page=1&limit=10&search=john&sort_by=status&sort_direction=desc' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: id'
```

### Success Response

```json
{
  "message": "Data user berhasil diambil.",
  "data": {
    "items": [
      {
        "id": 2,
        "name": "John Doe",
        "email": "john@example.com",
        "address": "Jakarta",
        "telp": "08123456789",
        "status": 2,
        "created_by": 1,
        "updated_by": 1,
        "created_at": "2026-06-28T09:30:00+07:00",
        "updated_at": "2026-06-28T09:30:00+07:00"
      },
      {
        "id": 3,
        "name": "Jane Doe",
        "email": "jane@example.com",
        "address": "Bandung",
        "telp": "082233445566",
        "status": 1,
        "created_by": 1,
        "updated_by": 1,
        "created_at": "2026-06-28T09:45:00+07:00",
        "updated_at": "2026-06-28T09:45:00+07:00"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 10,
      "total_items": 2,
      "total_pages": 1
    }
  }
}
```

### Catatan Validasi Status

- field `status` wajib diisi
- nilai yang diizinkan hanya `0`, `1`, `2`, atau `3`
- jika nilai di luar daftar tersebut maka request ditolak

### Create User

`POST /api/v1/users`

### Headers

- `Accept: application/json`
- `Content-Type: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Body

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "address": "Jakarta",
  "telp": "08123456789",
  "status": 2
}
```

Keterangan field:

- `name` wajib diisi, maksimal `100` karakter
- `email` wajib diisi, format email valid, maksimal `255` karakter, dan harus unik untuk user yang belum dihapus
- `status` wajib diisi dengan nilai `0`, `1`, `2`, atau `3`
- `address` opsional
- `telp` opsional, maksimal `30` karakter

Jika email sudah dipakai oleh user yang belum dihapus, endpoint mengembalikan error `400`.

### Error Response Example

```json
{
  "errors": [
    "Email sudah digunakan."
  ]
}
```

### Curl

```bash
curl --location 'http://localhost:9100/api/v1/users' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --header 'Accept-Language: id' \
  --data-raw '{
    "name": "John Doe",
    "email": "john@example.com",
    "address": "Jakarta",
    "telp": "08123456789",
    "status": 2
  }'
```

### Detail User

`GET /api/v1/users/:id`

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location 'http://localhost:9100/api/v1/users/1' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: id'
```

### Update User

`PUT /api/v1/users/:id`

### Headers

- `Accept: application/json`
- `Content-Type: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Body

```json
{
  "name": "John Doe Updated",
  "email": "john.updated@example.com",
  "address": "Bandung",
  "telp": "08123450000",
  "status": 3
}
```

Keterangan field:

- `name` wajib diisi, maksimal `100` karakter
- `email` wajib diisi, format email valid, maksimal `255` karakter, dan harus unik untuk user yang belum dihapus
- `status` wajib diisi dengan nilai `0`, `1`, `2`, atau `3`
- `address` opsional
- `telp` opsional, maksimal `30` karakter

Jika email sudah dipakai oleh user yang belum dihapus, endpoint mengembalikan error `400`.

### Curl

```bash
curl --location --request PUT 'http://localhost:9100/api/v1/users/1' \
  --header 'Accept: application/json' \
  --header 'Content-Type: application/json' \
  --header 'Accept-Language: id' \
  --data-raw '{
    "name": "John Doe Updated",
    "email": "john.updated@example.com",
    "address": "Bandung",
    "telp": "08123450000",
    "status": 3
  }'
```

### Delete User

`DELETE /api/v1/users/:id`

Endpoint ini melakukan soft delete dengan mengisi `deleted_at`.
Email milik user yang sudah dihapus dapat digunakan lagi saat membuat user baru.

### Headers

- `Accept: application/json`
- `Accept-Language: id` atau `Accept-Language: en`

### Curl

```bash
curl --location --request DELETE 'http://localhost:9100/api/v1/users/1' \
  --header 'Accept: application/json' \
  --header 'Accept-Language: id'
```
