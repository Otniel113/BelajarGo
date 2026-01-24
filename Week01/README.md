# Progress Week 1

CRUD dasar tanpa database menggunakan struct.

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
