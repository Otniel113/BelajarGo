# Progress Week 2

CRUD dengan database PostgreSQL, Modularize, dan Setup environment file.

## Setup Go dan Environment

Persiapan singkat untuk menjalankan proyek Week02 (`cashier`) secara lokal.

- Prasyarat:
    - Install Go (versi 1.20+ direkomendasikan). Cek versi dengan:

```bash
go version
```

- Masuk ke folder project:

```bash
cd "BelajarGo/Week02/cashier"
```

- Install dependency dan bersihkan module:

```bash
go get github.com/jackc/pgx/v5/stdlib
go get github.com/spf13/viper
go mod tidy
```

- Salin file environment dan sesuaikan nilainya (copy ` .env.example` ke `.env`):

Windows (PowerShell):
```powershell
Copy-Item .env.example .env
notepad .env
```

macOS / Linux:
```bash
cp .env.example .env
nano .env
```

- Isi `DB_CONN` pada `.env` dengan connection string Postgres kamu. Contoh format:

```
DB_CONN=postgresql://<DB_USER>:<DB_PASSWORD>@<HOST>:<PORT>/<DB_NAME>?sslmode=require
PORT=8080
```

Catatan:
- Jika password mengandung karakter khusus (mis. `!`, `@`, `:`), pastikan URL-encode karakter tersebut (mis. `!` -> `%21`).
- Pastikan `DB_CONN` sesuai dengan nama variable yang digunakan aplikasi (`DB_CONN`).

- Buat tabel `categories` jika belum ada. Contoh menjalankan `init.sql`:

Jika pakai `psql` CLI:
```bash
psql "${DB_CONN}" -f init.sql
```

Atau jalankan query SQL berikut di PgAdmin / Supabase SQL editor:

```sql
CREATE TABLE IF NOT EXISTS categories (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        description TEXT
);
```

- Menjalankan aplikasi:

```bash
go run main.go
```

## API Endpoint

| Method | Endpoint | Description | Request Body | Response |
|---|---|---|---|---|
| GET | `/categories` | Get all categories | - | `[]Category` |
| POST | `/categories` | Create category | `{"name": "...", "description": "..."}` | `Category` (Created) |
| GET | `/categories/{id}` | Get category by ID | - | `Category` |
| PUT | `/categories/{id}` | Update category | `{"name": "...", "description": "..."}` | `Category` |
| DELETE | `/categories/{id}` | Delete category | - | `{"message": "success delete category"}` |

### Request & Response Examples

#### 1. GET /categories
**Description:** Mengambil semua data category.
**Response (200 OK):**
```json
[
    {
        "id": 1,
        "name": "Elektronik",
        "description": "Barang-barang elektronik"
    },
    {
        "id": 2,
        "name": "Pakaian",
        "description": "Berbagai jenis pakaian"
    }
]
```

#### 2. POST /categories
**Description:** Menambahkan category baru.
**Request Body:**
```json
{
    "name": "Makanan",
    "description": "Segala jenis makanan dan minuman"
}
```
**Response (201 Created):**
```json
{
    "id": 3,
    "name": "Makanan",
    "description": "Segala jenis makanan dan minuman"
}
```

#### 3. GET /categories/{id}
**Description:** Mengambil satu data category berdasarkan ID.
**Example:** `/categories/1`
**Response (200 OK):**
```json
{
    "id": 1,
    "name": "Elektronik",
    "description": "Barang-barang elektronik"
}
```
**Response (404 Not Found):**
```text
Category not found
```

#### 4. PUT /categories/{id}
**Description:** Mengupdate data category berdasarkan ID.
**Example:** `/categories/1`
**Request Body:**
```json
{
    "name": "Elektronik & Gadget",
    "description": "Smartphone, Laptop, dan Aksesoris"
}
```
**Response (200 OK):**
```json
{
    "id": 1,
    "name": "Elektronik & Gadget",
    "description": "Smartphone, Laptop, dan Aksesoris"
}
```

#### 5. DELETE /categories/{id}
**Description:** Menghapus category berdasarkan ID.
**Example:** `/categories/1`
**Response (200 OK):**
```json
{
    "message": "success delete category"
}
```

## Postman Collections

Koleksi Postman dapat diakses di file: [Categories-API-Collections.json](Categories-API-Collections.json)

### Panduan Import ke Postman:

1.  Buka aplikasi **Postman**.
2.  Klik tombol **Import** di pojok kiri atas (di bawah nama Workspace).
3.  Pilih tab **File** dan klik **files** atau drag-and-drop file `Categories-API-Collections.json`.
4.  Klik **Import** untuk konfirmasi.
5.  Pastikan Server Go kamu sudah berjalan di `localhost:8080` jika ingin menjalankan di local.
6.  Gunakan environment variable `{{goURL}}` untuk mengganti linknya jika ingin melihatnya di deployment.
