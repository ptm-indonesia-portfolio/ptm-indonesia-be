Context7:
- Wajib pakai Context7 saat butuh pembuatan kode, langkah setup/konfigurasi, atau dokumentasi library/API.

Aturan Umum:
- Jangan running aplikasi (server) karena sudah dijalankan; cukup jalankan integration test atau unit test yang relevan, kalau ada error langsung diperbaiki.
- kalau buat service baru utamakan bikin dokumentasi terlebih dahulu di folder docs
- Pattern lebih baik banyak file yang spesifik dan mudah di-maintain/debug.
- Hindari penggunaan console yang tidak perlu karena cepat penuh log docker; kecuali console untuk error.
- Field di DB gunakan snake_case (contoh created_at)
- Kalau menambahkan table baru wajib tambahkan field created_by dan updated_by (default 0 untuk id by system).
- Isi file .env harus sama seperti env.example (jangan menambahkan nilai sensitif ke repo).
- Isi file .env.testing/env.testing hanya untuk konfigurasi testing lokal dan tidak boleh berisi credential sensitif/production; gunakan nilai dummy atau khusus testing.
- kalau ada script yang sering duplicate utamakan buat helper
- interface dan struct taruh di model
- Untuk Depedency Injection kita pakainya google wire
- Untuk database migration kita pakainya golang-migration
- Untuk perubahan config runtime/local development, jangan ubah 1 file saja.
  Wajib cek konsistensi minimal di `.env`, `.env.testing`, `env.example`, `docker-compose.yml`, dan kode config yang membaca env.
- Jangan mengandalkan mount file `.env` di container sebagai requirement runtime.
  Container dan aplikasi harus tetap bisa jalan dari environment variable; file `.env` hanya fallback untuk local bila memang dibutuhkan.
- Kalau perubahan menyentuh URL/origin/callback/auth frontend-backend:
  wajib cek dan sinkronkan `APP_PORT`, `BASEURL`, `FRONTEND_URL`, `GOOGLE_REDIRECT_URI`, dan `CORS_ORIGINS`.
- Jangan mengasumsikan origin frontend, callback URL, atau port yang aktif.
  Selalu cek env/config yang dipakai project terkait sebelum mengubah auth, CORS, cookie, redirect, atau Docker.
- Untuk local development di Windows, pastikan tooling yang menghasilkan binary sementara (contoh `air`) punya konfigurasi yang kompatibel dengan Windows.
- Tool auto reload untuk local development tidak boleh menjalankan migration atau seed otomatis, kecuali memang diminta user secara eksplisit.
- Saat ubah startup/runtime behavior, pastikan aplikasi tetap memberi feedback yang jelas saat dijalankan lokal
  (minimal info bahwa server berhasil start dan bind ke address yang dipakai).

Logging:
- Untuk logging kita pakainya library Logrus dan untuk logging kita pakai file
- logging level minimunnya error agar log file tidak cepat penuh
- Untuk formatnya tolong pakai json formatter

Pagination dan Sorting:
- Kalau ada endpoint list dengan paging: wajib ada sorting stabil dengan id desc (untuk mencegah data pindah halaman/duplikasi).
- Jika sudah ada sort utama lain: id desc tetap dipakai sebagai tie-breaker.

Validation dan Error:
- Untuk validasi kita pakai github.com/go-playground/validator
- Semua validation error harus status 400.
- Bentuk error validation harus seperti ini:
{
  "errors": [
    "The email is required!",
    "The name is required!",
    "The role is required!"
  ]
}

Upload File:
- Untuk update yang ada upload file: request body wajib punya status_file.
- status_file = 0: tidak ada perubahan file.
- status_file = 1 + ada upload file: ganti file.
- status_file = 1 + tidak ada upload file: hapus file.

Dokumentasi API:
- Jangan asal membuat file .md (repo public, hindari info sensitif).
- Boleh membuat/mengubah docs hanya untuk dokumentasi API di folder docs.
- Jika tambah endpoint baru: wajib buat file docs baru khusus modul tersebut di folder docs (jangan menambahkannya ke file docs modul lain).
- Jika ubah endpoint yang sudah ada: wajib update file docs modul terkait (contoh curl tanpa Cookie; header minimal Accept + Content-Type bila perlu).

Environment:
- Anggap semua environment non-development (production, staging, preprod, dll) sebagai "production-like".
- Untuk perubahan auth/login/session/cookie pada local FE-BE, selalu pikirkan perilaku browser lintas origin:
  CORS, cookie flags, callback URL, redirect URL, dan `withCredentials` harus konsisten.
- Untuk cookie auth:
  bedakan dengan jelas cookie yang hanya untuk backend (`HttpOnly`) dan cookie yang memang perlu dibaca frontend.
- Perubahan config auth/cors/cookie tidak boleh berhenti di backend saja;
  cek apakah docs, env, Docker, dan flow frontend yang relevan ikut konsisten.

Verifikasi:
- Jika mengubah Docker, env, CORS, auth, session, cookie, redirect, atau startup config:
  minimal jalankan `go test ./...`.
- Jika perubahan menyentuh login, refresh token, cookie, auth middleware, atau Google auth:
  jalankan juga integration test yang relevan.
- Jika perubahan menyentuh runtime/local config:
  verifikasi health endpoint dan, bila relevan, header CORS/origin yang aktif.

Arsitektur dan Struktur Folder:
- docs: dokumentasi API (contoh curl). Tidak perlu pakai header Cookie.
- config: untuk config yang dibutuhkan
- controllers: untuk controller
- helper: untuk helper
- lang: untuk data multi language
- model: untuk entity, DTO, request, response, custom constraint, type sets atau apapun untuk tipe data yang biasanya isinya struct atau interface
- repository: untuk yang berhubungan query db
    - contract: untuk contractnya
- services: untuk service
- validation: untuk file yang berhubungan validation
