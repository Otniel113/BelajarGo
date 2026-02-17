# Progress Week 5 - Complete

Proyek ini adalah versi final dari API Kasir (Cashier API) yang telah dilengkapi dengan fitur production-ready:
1.  **Autentikasi & Otorisasi:** Implementasi Login dan Register menggunakan JWT (JSON Web Tokens).
2.  **Role-Based Access Control (RBAC):** Hak akses berbeda untuk `admin` dan `member`.
3.  **Routing Modern:** Menggunakan library `chi` yang tetap kompatibel dengan middleware standar `net/http`.
4.  **Keamanan:** Hashing password menggunakan `bcrypt` dan middleware validasi JWT.
5.  **Pendaftaran Otomatis:** User baru akan otomatis terdaftar sebagai `member` secara default.

## Persiapan Go dan Environment

Ikuti langkah-langkah berikut untuk menjalankan proyek `cashier` secara lokal.

1.  **Versi Go:** Pastikan Anda menggunakan Go versi 1.22 ke atas (diperlukan untuk pola routing `net/http` terbaru dan `chi`).
2.  **Masuk ke folder project:**
    ```bash
    cd "Week05-Complete/cashier"
    ```
3.  **Install Dependency:**
    ```bash
    go mod tidy
    ```
4.  **Konfigurasi Environment (`.env`):**
    Salin file `.env.example` menjadi `.env` dan sesuaikan nilainya.

    **Windows (PowerShell):**
    ```powershell
    Copy-Item .env.example .env
    ```
    **macOS / Linux:**
    ```bash
    cp .env.example .env
    ```

    Isi file `.env` dengan detail database Anda:
    ```dotenv
    DB_CONN=postgresql://user:pass@host:port/dbname?sslmode=require
    PORT=8080
    JWT_SECRET=kode_rahasia_jwt_anda
    ```
    *Catatan: API_KEY sudah tidak digunakan lagi dan digantikan sepenuhnya oleh JWT.*

5.  **Migrasi Database:**
    Jalankan file `init.sql` untuk membuat tabel yang diperlukan (`users`, `categories`, `products`, `transactions`, `transaction_details`).

    Jika menggunakan `psql` CLI:
    ```bash
    psql "${DB_CONN}" -f init.sql
    ```

6.  **Menjalankan Aplikasi:**
    ```bash
    go run main.go
    ```

## Aturan Autentikasi & Otorisasi (RBAC)

### Hak Akses Per Role
- **Admin:** Akses penuh untuk mengelola (Create, Update, Delete) kategori dan produk, serta melihat laporan keuangan.
- **Member:** Bisa melihat kategori/produk dan melakukan transaksi (checkout).
- **Public (Non-Login):** Hanya bisa melihat daftar kategori/produk dan melakukan pendaftaran/login.

### Persyaratan Request
Selain endpoint publik, semua request wajib menyertakan header:
`Authorization: Bearer <token_jwt_anda>`

## Daftar API Endpoint

### Autentikasi
| Method | Endpoint    | Deskripsi                     | Body                                    | Query | Akses  |
| ------ | ----------- | ----------------------------- | --------------------------------------- | ----- | ------ |
| POST   | `/register` | Daftar sebagai `member`       | `{"username", "email", "password"}`     | -     | Publik |
| POST   | `/login`    | Login (mendapat token 30m)    | `{"identity", "password"}`              | -     | Publik |
| POST   | `/logout`   | Logout (client-side discard)  | -                                       | -     | Publik |

### Kategori & Produk
| Method | Endpoint           | Deskripsi                    | Body                                            | Query        | Hak Akses   |
| ------ | ------------------ | ---------------------------- | ----------------------------------------------- | ------------ | ----------- |
| GET    | `/categories`      | Ambil semua kategori         | -                                               | -            | Publik      |
| GET    | `/categories/{id}` | Detail kategori              | -                                               | -            | Publik      |
| POST   | `/categories`      | Tambah kategori baru         | `{"name", "description"}`                       | -            | Admin       |
| PUT    | `/categories/{id}` | Update kategori              | `{"name", "description"}`                       | -            | Admin       |
| DELETE | `/categories/{id}` | Hapus kategori               | -                                               | -            | Admin       |
| GET    | `/products`        | Ambil semua produk           | -                                               | `name`       | Publik      |
| POST   | `/products`        | Tambah produk baru           | `{"name", "price", "stock", "category_id"}`      | -            | Admin       |
| PUT    | `/products/{id}`   | Update produk                | `{"name", "price", "stock", "category_id"}`      | -            | Admin       |
| DELETE | `/products/{id}`   | Hapus produk                 | -                                               | -            | Admin       |

### Transaksi & Laporan
| Method | Endpoint         | Deskripsi                          | Body                                           | Query                     | Hak Akses    |
| ------ | ---------------- | ---------------------------------- | ---------------------------------------------- | ------------------------- | ------------ |
| POST   | `/checkout`      | Buat transaksi baru (Checkout)     | `{"items": [{"product_id", "quantity"}]}`      | -                         | Member       |
| GET    | `/report`        | Laporan lengkap berdasarkan range  | -                                              | `start_date`, `end_date`  | Admin        |
| GET    | `/report/today`  | Laporan transaksi hari ini         | -                                              | -                         | Admin        |

## Validasi & Fitur Tambahan
- **Registrasi:** Email melalui validasi format regex. Password harus minimal 6 karakter.
- **Login:** Token JWT berlaku selama **30 menit**.
- **Logout:** Endpoint `/logout` tersedia untuk konfirmasi aksi logout di sisi client.
- **Checkout:** User ID akan otomatis tersimpan dalam database transaksi berdasarkan token JWT yang dilampirkan.

## Koleksi Postman
Koleksi Postman dapat diakses di file: [Cashier.postman_collection.json](Cashier.postman_collection.json)

### Cara Import ke Postman:
1. Buka aplikasi **Postman**.
2. Klik tombol **Import** di pojok kiri atas.
3. Pilih file `Cashier.postman_collection.json` dari folder ini.
4. Sesuaikan variabel `goURL` di tab Variables koleksi tersebut.
5. Gunakan tipe **Bearer Token** pada tab Authorization dan masukkan JWT yang didapat dari `/login` untuk mengakses endpoint terproteksi.
