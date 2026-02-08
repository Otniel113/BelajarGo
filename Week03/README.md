# Progress Week 3

Bisa mencari product berdasarkan nama dengan parameter query, membuat transaksi atau checkout, dan summary report hari ini dan rentang custom.

## Setup Go dan Environment

Persiapan singkat untuk menjalankan proyek Week03 (`cashier`) secara lokal.

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

- Buat tabel `categories` dan `products` jika belum ada. Contoh menjalankan `init.sql`:

Jika pakai `psql` CLI:
```bash
psql "${DB_CONN}" -f init.sql
```

Atau copy-paste query SQL di PgAdmin / Supabase SQL editor dari `init.sql`.

- Menjalankan aplikasi:

```bash
go run main.go
```

## API Endpoint

### Categories

| Method | Endpoint | Description | Request Body | Response |
|---|---|---|---|---|
| GET | `/categories` | Get all categories | - | `[]Category` |
| POST | `/categories` | Create category | `{"name": "...", "description": "..."}` | `Category` (Created) |
| GET | `/categories/{id}` | Get category by ID | - | `Category` |
| PUT | `/categories/{id}` | Update category | `{"name": "...", "description": "..."}` | `Category` |
| DELETE | `/categories/{id}` | Delete category | - | `{"message": "success delete category"}` |

### Products

| Method | Endpoint | Description | Request Body | Response |
|---|---|---|---|---|
| GET | `/products` | Get all products with Category (supports query params `name`) | Optional query params: `name` | `[]Product` |
| POST | `/products` | Create product | `{"name": "...", "price": 100, "stock": 10, "category_id": 1}` | `Product` (Created) |
| GET | `/products/{id}` | Get product by ID | - | `Product` |
| PUT | `/products/{id}` | Update product | `{"name": "...", "price": 200, "stock": 3, "category_id": 1}` | `Product` |
| DELETE | `/products/{id}` | Delete product | - | `{"message": "success delete product"}` |

### Transactions

| Method | Endpoint | Description | Request Body | Response |
|---|---|---|---|---|
| POST | `/api/checkout` | Create a transaction (checkout) from multiple items | `{"items":[{"product_id":1,"quantity":2}]}` | `Transaction` with `id`, `total_amount`, `created_at`, `details` |

- Headers: `Content-Type: application/json`
- Errors: `400` for invalid body, `500` for server/repo errors (e.g., insufficient stock)

### Reports

| Method | Endpoint | Description | Query Params | Response |
|---|---|---|---|---|
| GET | `/report/today` | Get report for today | - | `Report` |
| GET | `/report` | Get report for date range | `start_date`, `end_date` (YYYY-MM-DD) | `Report` |

**Response Example:**
```json
{
    "total_revenue": 45000,
    "total_transaction": 5,
    "most_sold_product":{
        "name": "Indomie Goreng",
        "sold_qty":12
    }
}
```

## Postman Collections

Koleksi Postman dapat diakses di file: [Cashier.postman_collection.json](Cashier.postman_collection.json)

### Panduan Import ke Postman:

1.  Buka aplikasi **Postman**.
2.  Klik tombol **Import** di pojok kiri atas (di bawah nama Workspace).
3.  Pilih tab **File** dan klik **files** atau drag-and-drop file `Cashier.postman_collection.json`.
4.  Klik **Import** untuk konfirmasi.
5.  Pastikan Server Go kamu sudah berjalan di `localhost:8080` jika ingin menjalankan di local.
6.  Gunakan environment variable `{{goURL}}` untuk mengganti linknya jika ingin melihatnya di deployment.
