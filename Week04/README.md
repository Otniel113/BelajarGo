# Progress Week 4

Menambahkan Middleware = API Key, CORS, dan Logger

## Setup Go dan Environment

Persiapan singkat untuk menjalankan proyek Week04 (`cashier`) secara lokal.

- Prasyarat:
    - Install Go (versi 1.20+ direkomendasikan). Cek versi dengan:

```bash
go version
```

- Masuk ke folder project:

```bash
cd "BelajarGo/Week04/cashier"
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
API_KEY=your_secret_api_key
```

Catatan:
- Pastikan semua request menyertakan header `X-API-Key` dengan nilai yang sesuai.
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
| POST | `/checkout` | Create a transaction (checkout) from multiple items | `{"items":[{"product_id":1,"quantity":2}]}` | `Transaction` with `id`, `total_amount`, `created_at`, `details` |

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

1.  Setup terlebih dahulu untuk nama variable URL dan API Key pada [Postman Collection](Cashier.postman_collection.json).
```json
"variable": [
		{
			"key": "goURL",
			"value": "http://localhost:8080",
			"type": "string"
		},
		{
			"key": "api_key",
			"value": "",
			"type": "string"
		}
]
```

2.  Buka aplikasi **Postman**.
3.  Pergi ke tab **Workspaces** di kiri atas dan pilih Workspaces yang diinginkan jika belum berada di Wokkspaces
3.  Klik tombol **Import** di pojok kiri atas (di bawah nama Workspace).
4.  Pilih tab **File** dan klik **files** atau drag-and-drop file `Cashier.postman_collection.json`.
5.  Klik **Import** untuk konfirmasi.
6.  Jika ingin ada perubahan variabel, klik **Collections** di kiri, pilih **Cashier**, klik titik tiga, pilih **Edit**, lalu tab **Variables**
7.  Ubah isi variabel dan klik **Save**
