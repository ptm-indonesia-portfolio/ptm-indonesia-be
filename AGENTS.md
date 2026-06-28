# Project Rules

Gunakan file ini sebagai ringkasan cepat. Aturan detail dipisah ke folder `rules`.

## Prioritas Penting

- Wajib pakai Context7 saat butuh pembuatan kode yang bergantung pada library/API, langkah setup/konfigurasi, atau dokumentasi library/API.
- Jangan menjalankan aplikasi/server secara otomatis. Anggap project biasanya sudah dijalankan user; cukup jalankan test yang relevan kecuali user meminta hal lain.
- Smoke test browser/manual tidak boleh dijalankan otomatis setelah perubahan. Jalankan hanya jika user memberi instruksi eksplisit untuk smoke test.
- Status code adalah kontrak penting project ini. Validation error dan business validation harus mengikuti mapping status code yang sudah ditetapkan di rules.
- Jika menambah service baru, utamakan buat/update dokumentasi API modul tersebut di folder `docs` terlebih dahulu.
- Jika controller mulai panjang karena route registration, grouping, atau middleware, pisahkan ke folder/package `routes`.

## Index Rules

- `rules/01-general.mdx`
  Aturan umum coding, struktur file, logging, DI, migration, naming, dan pola maintainability.
- `rules/02-api-contract.mdx`
  Kontrak API penting: validation, status code, pagination, sorting, dan upload file.
- `rules/03-docs-and-frontend.mdx`
  Aturan penulisan docs API dan laporan yang dibaca frontend/integrator.
- `rules/04-runtime-and-environment.mdx`
  Aturan env, Docker, auth/cookie/CORS lokal, dan runtime development.
- `rules/05-testing-and-verification.mdx`
  Kapan harus menjalankan test, aturan smoke test, dan cara verifikasi runtime aktif.
- `rules/06-architecture.mdx`
  Tanggung jawab folder/package dan pemisahan struktur codebase.

## Catatan

- File `AGENTS.md` ini tetap dipakai sebagai ringkasan cepat.
- Detail aturan harus dirujuk dari file `rules/*.mdx`.
