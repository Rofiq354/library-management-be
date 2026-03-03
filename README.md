# 📚 Library Management System

Aplikasi web perpustakaan digital berbasis **Golang** dengan tampilan card-based yang clean dan modern. Pengunjung bisa menjelajahi koleksi buku, sementara Admin bisa mengelola semua data melalui panel yang terlindungi JWT.

---

## 🛠️ Tech Stack

| Komponen         | Teknologi            |
| ---------------- | -------------------- |
| Language         | Go (Golang)          |
| Web Framework    | Gin                  |
| ORM              | GORM                 |
| Database         | SQLite               |
| Auth             | JWT (JSON Web Token) |
| Password Hashing | bcrypt               |

---

## 📁 Folder Structure

```
library-app/
├── main.go                  ← Entry point, sesimpel mungkin
├── config/
│   └── database.go          ← Koneksi SQLite + AutoMigrate
├── models/
│   ├── book.go              ← Struct Book
│   ├── user.go              ← Struct User
│   └── category.go          ← Struct Category
├── handlers/
│   ├── book_handler.go      ← CRUD Buku
│   ├── user_handler.go      ← CRUD User
│   └── auth_handler.go      ← Login & Register
├── routes/
│   └── routes.go            ← Semua route dikumpulkan di sini
├── middleware/
│   └── auth.go              ← JWT Auth Middleware
├── templates/               ← File HTML (server-side rendering)
├── static/                  ← CSS, JS, gambar
├── go.mod
└── go.sum
```

---

## ⚙️ Setup & Instalasi

### 1. Prasyarat

Pastikan sudah install **Go** di komputer kamu:

```bash
go version
# Output: go version go1.22.0 ...
```

Belum install? Download di → https://golang.org/dl/

### 2. Clone / Buat Project

```bash
mkdir library-app
cd library-app
go mod init library-app
```

### 3. Install Dependencies

```bash
# Web framework
go get github.com/gin-gonic/gin

# ORM + SQLite driver
go get gorm.io/gorm
go get github.com/glebarez/sqlite

# JWT
go get github.com/golang-jwt/jwt/v5

# Password hashing
go get golang.org/x/crypto/bcrypt
```

### 4. Jalankan Aplikasi

```bash
go run main.go
```

Output yang muncul jika berhasil:

```
✅ Database berhasil terkoneksi!
✅ Tabel berhasil dibuat!
[GIN-debug] Listening and serving HTTP on :8080
```

Server berjalan di → **http://localhost:8080**

> **Catatan:** File database `library.db` akan dibuat otomatis di folder project saat pertama kali dijalankan. Tidak perlu install server database apapun!

---

## 🗄️ Database

Project ini menggunakan **SQLite** — database berbasis file, tidak perlu install server terpisah.

File database `library.db` dibuat otomatis saat aplikasi pertama kali dijalankan. GORM akan membuat tabel secara otomatis lewat **AutoMigrate** berdasarkan struct yang sudah didefinisikan.

### Schema

**Tabel `books`**

```
id, title, author, description, cover_url, stock, published_year,
category_id, created_at, updated_at, deleted_at
```

**Tabel `categories`**

```
id, name, created_at, updated_at, deleted_at
```

**Tabel `users`**

```
id, name, email, password (hashed), role, created_at, updated_at, deleted_at
```

Untuk melihat isi database secara visual, gunakan **DB Browser for SQLite** → https://sqlitebrowser.org/dl/

---

## 🛣️ API Endpoints

### 🔓 Public Routes (Bebas diakses)

| Method | Endpoint     | Deskripsi                 |
| ------ | ------------ | ------------------------- |
| `POST` | `/register`  | Buat akun admin baru      |
| `POST` | `/login`     | Login, mendapat JWT token |
| `GET`  | `/books`     | Ambil semua buku          |
| `GET`  | `/books/:id` | Ambil detail satu buku    |

### 🔒 Admin Routes (Wajib JWT Token)

| Method   | Endpoint           | Deskripsi         |
| -------- | ------------------ | ----------------- |
| `POST`   | `/admin/books`     | Tambah buku baru  |
| `PUT`    | `/admin/books/:id` | Update data buku  |
| `DELETE` | `/admin/books/:id` | Hapus buku        |
| `GET`    | `/admin/users`     | Ambil semua user  |
| `GET`    | `/admin/users/:id` | Ambil detail user |
| `POST`   | `/admin/users`     | Tambah user baru  |
| `PUT`    | `/admin/users/:id` | Update data user  |
| `DELETE` | `/admin/users/:id` | Hapus user        |

---

## 🧪 Cara Test API

Gunakan **Thunder Client** (VS Code Extension) atau **Postman**.

### Step 1 — Register akun admin

```
POST http://localhost:8080/register
Content-Type: application/json

{
    "name": "Admin",
    "email": "admin@library.com",
    "password": "rahasia123"
}
```

### Step 2 — Login

```
POST http://localhost:8080/login
Content-Type: application/json

{
    "email": "admin@library.com",
    "password": "rahasia123"
}
```

Response:

```json
{
  "message": "Login berhasil!",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "name": "Admin",
    "email": "admin@library.com",
    "role": "admin"
  }
}
```

### Step 3 — Gunakan token untuk akses Admin routes

Tambahkan header berikut di setiap request ke `/admin/...`:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Step 4 — Tambah buku baru

```
POST http://localhost:8080/admin/books
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "Fire and Blood",
    "author": "George R.R. Martin",
    "description": "Sejarah House Targaryen di Westeros.",
    "cover_url": "https://example.com/cover.jpg",
    "stock": 5,
    "published_year": 2018,
    "category_id": 1
}
```

---

## 🔐 Cara Kerja Autentikasi

```
Super Admin kirim email + password
         │
         ▼
Cari admin di database berdasarkan email
         │
         ▼
Bandingkan password pakai bcrypt
         │
         ▼
Buat JWT Token (berlaku 24 jam)
         │
         ▼
Kirim token ke client

── Untuk akses Admin ──────────────────

Request ke /admin/...
         │
         ▼
Middleware cek header Authorization
         │
    ┌────┴────┐
Tidak ada    Ada token
    │              │
  401 ❌        Validasi JWT
                   │
              ┌────┴────┐
           Invalid    Valid
              │          │
            401 ❌     Lanjut ke handler ✅
```

---

## 📖 Konsep Penting yang Dipelajari

### `r` vs `c` di Gin

|               | Penjelasan                                                                                  |
| ------------- | ------------------------------------------------------------------------------------------- |
| `r` (Router)  | Dibuat sekali saat app start. Bertugas mendaftarkan semua route.                            |
| `c` (Context) | Dibuat per request. Bertugas melayani satu user — ambil data dari request & kirim response. |

### 3 Cara ambil data dari request

```go
c.Param("id")            // dari URL path  → /books/:id
c.Query("search")        // dari URL query → /books?search=narnia
c.ShouldBindJSON(&input) // dari body JSON → POST/PUT request
```

### Pointer `*` dan `&`

```go
// & = "alamat" dari variabel
// * = "pergi ke alamat itu"

c.ShouldBindJSON(&input)
// Gin dapat alamat input → langsung isi input yang asli ✅
// Bukan fotokopi → perubahan benar-benar tersimpan
```

### Struct + JSON Tag

```go
type Book struct {
    Title  string `json:"title"`   // JSON "title" → Go Title
    Author string `json:"author"`  // JSON "author" → Go Author
}
```

### Middleware

```go
// Dipasang di group route
admin := r.Group("/admin", middleware.AuthMiddleware())

// Semua route di dalam group ini otomatis terlindungi
// Request masuk → cek token dulu → baru masuk ke handler
```

---

## 🚀 Roadmap Pengembangan

- [x] Setup project & folder structure
- [x] Koneksi database SQLite + AutoMigrate
- [x] CRUD Buku
- [x] CRUD User
- [x] Autentikasi JWT (Login/Register)
- [x] Middleware auth untuk proteksi route admin
- [ ] Templates HTML (tampilan card-based)
- [ ] Search & filter buku
- [ ] Admin dashboard
- [ ] Filter buku per kategori
- [ ] Upload gambar cover buku

---

## 📦 Dependencies

```
github.com/gin-gonic/gin        — Web framework
gorm.io/gorm                    — ORM
github.com/glebarez/sqlite           — SQLite driver untuk GORM
github.com/golang-jwt/jwt/v5    — JWT auth
golang.org/x/crypto/bcrypt      — Password hashing
```

---

> Project ini dibuat sebagai **Custom Project** untuk keperluan pembelajaran Golang.
> Deadline: **6 Maret 2026**
