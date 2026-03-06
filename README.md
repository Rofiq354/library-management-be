# 📚 Library Management System

Aplikasi web perpustakaan digital berbasis **Golang** dengan tampilan card-based yang clean dan modern. Sistem ini dirancang untuk memudahkan pengelolaan perpustakaan dengan struktur role yang jelas dan aman.

---

## 📖 Deskripsi Aplikasi

**Library Management System** adalah platform perpustakaan digital yang memungkinkan tiga tingkat akses berbeda:

### 👥 Stakeholder & Role

#### 1. **Superadmin** 🔐

Pengguna dengan akses penuh ke seluruh sistem untuk mengelola sumber daya.

**Akses & Fitur:**

- 👤 **Manajemen User** - Create, read, update, delete user (Siswa & Admin)
- 📊 **Dashboard Kontrol Penuh** - Monitoring keseluruhan sistem
- ⚙️ **Pengaturan Sistem** - Konfigurasi global aplikasi

#### 2. **Admin** 📚

Pengguna yang bertanggung jawab untuk pengelolaan perpustakaan dan konten.

**Akses & Fitur:**

- 📖 **Manajemen Buku** - Create, read, update, delete buku perpustakaan
- 🏷️ **Manajemen Kategori** - Create, read, update, delete kategori buku
- 📊 **Dashboard Admin** - Melihat statistik perpustakaan dan aktivitas pengguna
- 📤 **Upload PDF** - Mengunggah file PDF untuk setiap buku

#### 3. **Siswa** 📖

Pengguna akhir yang dapat mengakses dan memanfaatkan koleksi perpustakaan.

**Akses & Fitur:**

- 🔍 **Jelajahi Koleksi** - Melihat semua buku dalam tampilan card-based yang menarik
- 🔎 **Search & Filter** - Pencarian buku berdasarkan:
  - 📓 Judul buku
  - ✍️ Nama author/penulis
  - 🏷️ Kategori buku
- 📄 **Baca Buku** - Membaca preview atau konten buku secara online
- ⬇️ **Download PDF** - Mengunduh file PDF buku untuk dibaca offline

---

## 🛠️ Tech Stack

| Komponen         | Teknologi            | Fungsi                                   | Dependency                               |
| ---------------- | -------------------- | ---------------------------------------- | ---------------------------------------- |
| Language         | Go (Golang)          | Backend language                         | `go 1.25.0`                              |
| Web Framework    | Gin                  | HTTP web framework yang lightweight      | `github.com/gin-gonic/gin`               |
| ORM              | GORM                 | Object-Relational Mapping untuk database | `gorm.io/gorm`                           |
| Database         | SQLite               | Database relasional yang ringan          | `github.com/glebarez/sqlite`             |
| Auth             | JWT (JSON Web Token) | Secure authentication & session          | `github.com/golang-jwt/jwt/v5`           |
| Password Hashing | bcrypt               | Enkripsi password yang aman              | `golang.org/x/crypto/bcrypt`             |
| File Storage     | Cloudinary           | Cloud storage untuk image & PDF          | `github.com/cloudinary/cloudinary-go/v2` |
| CORS             | gin-contrib/cors     | Cross-Origin Resource Sharing middleware | `github.com/gin-contrib/cors`            |
| Env Management   | godotenv             | Load environment variables dari .env     | `github.com/joho/godotenv`               |

---

## 📁 Folder Structure

```
library-app/
├── main.go                  ← Entry point aplikasi
├── config/
│   ├── cloudinary.go        ← Setup koneksi Cloudinary untuk upload image & PDF
│   ├── config.go            ← Load environment variables dari .env
│   └── database.go          ← Koneksi SQLite + AutoMigrate
├── handlers/                ← Logic CRUD semua entity
├── middleware/              ← JWT Authentication & middleware lainnya
├── models/                  ← Struct data (User, Book, Category, dll.)
├── routes/                  ← Definisi semua endpoint aplikasi
├── templates/               ← File HTML untuk server-side rendering
├── static/                  ← CSS, JavaScript, dan asset gambar
├── tmp/                     ← Temporary files
├── go.mod
├── go.sum
├── .env                     ← Environment variables (jangan commit)
├── .gitignore
├── library.db               ← Database SQLite
└── README.md
```

## ⚙️ Setup & Instalasi

<details>
<summary><strong>1️⃣ Prasyarat</strong></summary>

Pastikan sudah install **Go** di komputer kamu:

```bash
go version
# Output: go version go1.22.0 ...
```

Belum install? Download di → https://golang.org/dl/

**Versi minimum yang direkomendasikan:** Go 1.20 atau lebih baru

</details>

<details>
<summary><strong>2️⃣ Clone / Buat Project</strong></summary>

```bash
mkdir library-management-be
cd library-management-be
go mod init library-app
```

**Hasil struktur:**
```
library-management-be/
├── go.mod
└── (folder lainnya akan dibuat kemudian)
```

</details>

<details>
<summary><strong>3️⃣ Install Dependencies</strong></summary>

```bash
go mod tidy
```

Perintah `go mod tidy` akan secara otomatis:

- ✅ Download semua dependencies yang diperlukan
- ✅ Menghapus dependencies yang tidak digunakan
- ✅ Update `go.mod` dan `go.sum` dengan versi terbaru

**Dependencies yang akan diinstall:**

| Nama | Fungsi |
| ---- | ------ |
| **Gin** | Web framework untuk REST API |
| **GORM** | ORM untuk interaksi database |
| **SQLite** | Database driver |
| **JWT** | Authentication & session |
| **bcrypt** | Password hashing |
| **Cloudinary** | Cloud storage untuk image & PDF |
| **CORS** | Cross-Origin Resource Sharing middleware |
| **godotenv** | Load environment variables |

</details>

<details>
<summary><strong>4️⃣ Setup Environment Variables</strong></summary>

Buat file `.env` di root folder project:

```bash
touch .env
```

Edit `.env` dan masukkan konfigurasi:

```env
PORT=8080
JWT_SECRET=your_secret_key_here
FRONTEND_ORIGIN=http://localhost:3000
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

**Catatan:** Jangan commit file `.env` ke repository! Tambahkan ke `.gitignore`

</details>

<details>
<summary><strong>5️⃣ Jalankan Aplikasi</strong></summary>

```bash
go run main.go
```

**Output yang muncul jika berhasil:**

```
✅ Database berhasil terkoneksi!
✅ Tabel berhasil dibuat!
✅ Superadmin berhasil dibuat!
        Email: superadmin@library.com
        Password: superadmin123
[GIN-debug] Listening and serving HTTP on :8080
```

**Server berjalan di:** http://localhost:8080

**Catatan:**
- File database `library.db` akan dibuat otomatis di folder project saat pertama kali dijalankan
- Tidak perlu install server database apapun (SQLite berbasis file)
- Gunakan credentials default superadmin di atas untuk login pertama kali
- **PENTING:** Ubah password superadmin setelah login untuk keamanan!

</details>

---

## 🗄️ Database

Project ini menggunakan **SQLite** — database berbasis file, tidak perlu install server terpisah.

File database `library.db` dibuat otomatis saat aplikasi pertama kali dijalankan. GORM akan membuat tabel secara otomatis lewat **AutoMigrate** berdasarkan struct yang sudah didefinisikan.

### Schema

<details>
<summary><strong>📚 Tabel `books`</strong></summary>

| Field            | Type      | Constraint           | Deskripsi                                      |
| ---------------- | --------- | -------------------- | ---------------------------------------------- |
| `id`             | `UINT`    | PRIMARY KEY          | ID unik buku                                   |
| `title`          | `STRING`  | NOT NULL             | Judul buku                                     |
| `author`         | `STRING`  | NOT NULL             | Nama penulis                                   |
| `description`    | `TEXT`    |                      | Deskripsi/sinopsis buku                        |
| `cover_url`      | `STRING`  |                      | URL cover dari Cloudinary                      |
| `pdf_url`        | `STRING`  |                      | URL PDF dari Cloudinary                        |
| `cover_public_id` | `STRING` |                      | Public ID cover di Cloudinary (untuk delete)   |
| `pdf_public_id`  | `STRING`  |                      | Public ID PDF di Cloudinary (untuk delete)     |
| `stock`          | `INT`     | NOT NULL             | Jumlah stok buku                               |
| `published_year` | `INT`     |                      | Tahun terbit                                   |
| `category_id`    | `UINT`    | FOREIGN KEY          | Referensi ke tabel `categories`                |
| `created_at`     | `TIMESTAMP` | AUTO                | Waktu pembuatan record                         |
| `updated_at`     | `TIMESTAMP` | AUTO                | Waktu update terakhir                          |
| `deleted_at`     | `TIMESTAMP` | NULLABLE (soft delete) | Waktu penghapusan (soft delete)             |

</details>

<details>
<summary><strong>🏷️ Tabel `categories`</strong></summary>

| Field      | Type      | Constraint  | Deskripsi               |
| ---------- | --------- | ----------- | ----------------------- |
| `id`       | `UINT`    | PRIMARY KEY | ID unik kategori        |
| `name`     | `STRING`  | NOT NULL    | Nama kategori (Fiksi, Non-Fiksi, dll) |
| `created_at` | `TIMESTAMP` | AUTO    | Waktu pembuatan record  |
| `updated_at` | `TIMESTAMP` | AUTO    | Waktu update terakhir   |
| `deleted_at` | `TIMESTAMP` | NULLABLE (soft delete) | Waktu penghapusan |

</details>

<details>
<summary><strong>👤 Tabel `users`</strong></summary>

| Field       | Type      | Constraint  | Deskripsi                              |
| ----------- | --------- | ----------- | -------------------------------------- |
| `id`        | `UINT`    | PRIMARY KEY | ID unik user                           |
| `name`      | `STRING`  | NOT NULL    | Nama lengkap user                      |
| `email`     | `STRING`  | UNIQUE      | Email user (untuk login)               |
| `password`  | `STRING`  | NOT NULL    | Password yang sudah di-hash dengan bcrypt |
| `role`      | `STRING`  | NOT NULL    | Role user: `superadmin`, `admin`, `siswa` |
| `created_at` | `TIMESTAMP` | AUTO    | Waktu pembuatan record                 |
| `updated_at` | `TIMESTAMP` | AUTO    | Waktu update terakhir                  |
| `deleted_at` | `TIMESTAMP` | NULLABLE (soft delete) | Waktu penghapusan |

</details>

<details>
<summary><strong>👨‍🎓 Tabel `siswa`</strong></summary>

| Field          | Type      | Constraint           | Deskripsi                          |
| -------------- | --------- | -------------------- | ---------------------------------- |
| `id`           | `UINT`    | PRIMARY KEY          | ID unik siswa                      |
| `user_id`      | `UINT`    | FOREIGN KEY, UNIQUE  | Referensi ke tabel `users` (one-to-one) |
| `nis`          | `STRING`  | UNIQUE               | Nomor Induk Siswa                  |
| `kelas`        | `STRING`  | NOT NULL             | Kelas siswa (X, XI, XII, dst)      |
| `jurusan`      | `STRING`  | NOT NULL             | Jurusan (IPA, IPS, dll)            |
| `no_telepon`   | `STRING`  |                      | Nomor telepon siswa                |
| `alamat`       | `TEXT`    |                      | Alamat lengkap siswa               |
| `tanggal_lahir` | `STRING` |                      | Tanggal lahir siswa                |
| `created_at`   | `TIMESTAMP` | AUTO                | Waktu pembuatan record             |
| `updated_at`   | `TIMESTAMP` | AUTO                | Waktu update terakhir              |
| `deleted_at`   | `TIMESTAMP` | NULLABLE (soft delete) | Waktu penghapusan             |

</details>

### Relasi Database

```
Categories ──────────── Books (1:N)
    ↑                      ↓
    └─ kategori_id ← category_id

Users ────────────── Siswa (1:1)
  ↑                    ↓
  └─ id (FK) ← user_id (UNIQUE)
  
  Relasi: CASCADE DELETE (hapus user → otomatis hapus siswa)
```

### Tools

Untuk melihat isi database secara visual, gunakan **DB Browser for SQLite** → https://sqlitebrowser.org/dl/

---

## 🛣️ API Endpoints

### 🔓 Public Routes (Bebas diakses)

<details>
<summary><strong>🔐 Authentication</strong></summary>

| Method | Endpoint     | Deskripsi                 |
| ------ | ------------ | ------------------------- |
| `POST` | `/api/login` | Login, mendapat JWT token |

</details>

<details>
<summary><strong>📚 Buku</strong></summary>

| Method | Endpoint         | Deskripsi               |
| ------ | ---------------- | ----------------------- |
| `GET`  | `/api/books`     | Ambil semua buku        |
| `GET`  | `/api/books/:id` | Ambil detail satu buku  |

</details>

<details>
<summary><strong>🏷️ Kategori</strong></summary>

| Method | Endpoint            | Deskripsi             |
| ------ | ------------------- | --------------------- |
| `GET`  | `/api/categories`   | Ambil semua kategori  |
| `GET`  | `/api/categories/:id` | Ambil detail kategori |

</details>

### 🔒 Protected Routes (Wajib JWT Token)

<details>
<summary><strong>🚪 Logout</strong></summary>

| Method | Endpoint      | Deskripsi   |
| ------ | ------------- | ----------- |
| `POST` | `/api/logout` | Logout user |

</details>

#### Admin Routes (Manajemen Buku & Kategori)

<details>
<summary><strong>📚 Buku Management</strong></summary>

| Method   | Endpoint               | Deskripsi        |
| -------- | ---------------------- | ---------------- |
| `POST`   | `/api/admin/books`     | Tambah buku baru |
| `PUT`    | `/api/admin/books/:id` | Update buku      |
| `DELETE` | `/api/admin/books/:id` | Hapus buku       |

</details>

<details>
<summary><strong>🏷️ Kategori Management</strong></summary>

| Method   | Endpoint                    | Deskripsi            |
| -------- | --------------------------- | -------------------- |
| `POST`   | `/api/admin/categories`     | Tambah kategori baru |
| `PUT`    | `/api/admin/categories/:id` | Update kategori      |
| `DELETE` | `/api/admin/categories/:id` | Hapus kategori       |

</details>

<details>
<summary><strong>👨‍🎓 Siswa Management</strong></summary>

| Method   | Endpoint               | Deskripsi          |
| -------- | ---------------------- | ------------------ |
| `GET`    | `/api/admin/siswa`     | Ambil semua siswa  |
| `GET`    | `/api/admin/siswa/:id` | Ambil detail siswa |
| `PUT`    | `/api/admin/siswa/:id` | Update data siswa  |
| `DELETE` | `/api/admin/siswa/:id` | Hapus siswa        |

</details>

#### Superadmin Routes (Manajemen User)

<details>
<summary><strong>👥 User Management</strong></summary>

| Method   | Endpoint                    | Deskripsi         |
| -------- | --------------------------- | ----------------- |
| `GET`    | `/api/superadmin/users`     | Ambil semua user  |
| `GET`    | `/api/superadmin/users/:id` | Ambil detail user |
| `POST`   | `/api/superadmin/users`     | Tambah user baru  |
| `PUT`    | `/api/superadmin/users/:id` | Update user       |
| `DELETE` | `/api/superadmin/users/:id` | Hapus user        |

</details>

## 🧪 Cara Test API

Gunakan **Thunder Client** (VS Code Extension) atau **Postman**.

<details>
<summary><strong>📝 Step 1 — Login</strong></summary>

**Endpoint:** `POST http://localhost:8080/api/login`

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
    "email": "superadmin@library.com",
    "password": "superadmin123"
}
```

**Response (200 OK):**
```json
{
  "message": "Login berhasil!",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "name": "Super Admin",
    "email": "superadmin@library.com",
    "role": "superadmin"
  }
}
```

**Catatan:**
- Token ini akan digunakan untuk semua request yang memerlukan authentication
- Token berlaku selama sesi (sesuai konfigurasi JWT)
- Copy token untuk step berikutnya

</details>

<details>
<summary><strong>🔐 Step 2 — Setup Authorization Header</strong></summary>

Untuk mengakses protected routes (admin & superadmin), tambahkan header berikut di setiap request:

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**Cara di Thunder Client:**
1. Buka tab "Headers"
2. Tambah header baru: `Authorization`
3. Value: `Bearer <paste_token_anda>`

**Cara di Postman:**
1. Buka tab "Authorization"
2. Pilih type: "Bearer Token"
3. Paste token di field "Token"

**Contoh pada gambar:**
```
┌─ Headers
├─ Authorization: Bearer eyJhbGci...
├─ Content-Type: application/json
└─ (headers lainnya)
```

</details>

<details>
<summary><strong>➕ Step 3 — Tambah Kategori Baru (Admin)</strong></summary>

**Endpoint:** `POST /api/admin/categories`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
    "name": "Fiksi"
}
```

**Response (201 Created):**
```json
{
  "ID": 1,
  "CreatedAt": "2026-03-04T13:34:26.1158963+07:00",
  "UpdatedAt": "2026-03-04T13:34:26.1158963+07:00",
  "DeletedAt": null,
  "name": "Fiksi",
  "books": null
}
```

**Catatan:**
- Category ID yang dibuat akan digunakan saat menambah buku
- Simpan ID kategori untuk reference step berikutnya

</details>

<details>
<summary><strong>📚 Step 4 — Tambah Buku Baru (Admin)</strong></summary>

**Endpoint:** `POST /api/admin/books`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: multipart/form-data
```

**Form Fields:**

| Field | Type | Contoh Value |
| ----- | ---- | ------------ |
| `title` | text | Dune |
| `author` | text | Frank Herbert |
| `description` | text | Kisah epik perebutan planet gurun Arrakis yang kaya akan rempah 'spice'. |
| `stock` | number | 30 |
| `published_year` | number | 1965 |
| `category_id` | number | 1 |
| `cover` | file | cover_image.jpg |
| `pdf` | file | book.pdf |

**Response (201 Created):**
```json
{
  "ID": 1,
  "CreatedAt": "2026-03-04T13:39:39.2090296+07:00",
  "UpdatedAt": "2026-03-04T13:39:39.2090296+07:00",
  "DeletedAt": null,
  "title": "Dune",
  "author": "Frank Herbert",
  "description": "Kisah epik perebutan planet gurun Arrakis yang kaya akan rempah 'spice'.",
  "cover_url": "https://res.cloudinary.com/.../.jpg",
  "pdf_url": "https://res.cloudinary.com/.../.pdf",
  "cover_public_id": "library/books/covers/...",
  "pdf_public_id": "library/books/pdfs/...",
  "stock": 30,
  "published_year": 1965,
  "category_id": 1,
  "category": {
    "ID": 1,
    "CreatedAt": "2026-03-04T13:34:26.1158963+07:00",
    "UpdatedAt": "2026-03-04T13:34:26.1158963+07:00",
    "DeletedAt": null,
    "name": "Fiksi",
    "books": null
  }
}
```

**Catatan:**
- File `cover` dan `pdf` di-upload ke **Cloudinary** secara otomatis
- Response berisi `cover_url` dan `pdf_url` yang dapat langsung digunakan di frontend
- `cover_public_id` dan `pdf_public_id` disimpan untuk update/delete file di Cloudinary kemudian hari
- Pastikan format file: `.jpg`/`.png` untuk cover, `.pdf` untuk document

</details>

<details>
<summary><strong>👨‍🎓 Step 5 — Tambah Siswa Baru (Superadmin)</strong></summary>

**Endpoint:** `POST /api/superadmin/users` (User dulu)

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body (Create User):**
```json
{
    "name": "Ahmad Rizki",
    "email": "ahmad@student.com",
    "password": "password123",
    "role": "siswa"
}
```

**Response (201 Created):**
```json
{
  "ID": 2,
  "CreatedAt": "2026-03-05T10:30:00+07:00",
  "UpdatedAt": "2026-03-05T10:30:00+07:00",
  "DeletedAt": null,
  "name": "Ahmad Rizki",
  "email": "ahmad@student.com",
  "role": "siswa"
}
```

**Catatan:**
- User yange talah dibuat akan membuat Siswa profile
- User ID yang digunakan adalah ID dari user yang baru dibuat

</details>

<details>
<summary><strong>📖 Step 6 — Akses Public Routes (Siswa)</strong></summary>

Siswa dapat mengakses routes public tanpa token:

**Get All Books:**
```
GET http://localhost:8080/api/books
```

**Response:**
```json
[
  {
    "ID": 1,
    "title": "Dune",
    "author": "Frank Herbert",
    "cover_url": "https://res.cloudinary.com/...",
    "stock": 30,
    "category": {
      "ID": 1,
      "name": "Fiksi"
    }
  },
  ...
]
```

**Get Book Detail:**
```
GET http://localhost:8080/api/books/1
```

**Search Books (dengan query parameter):**
```
GET http://localhost:8080/api/books?search=dune
GET http://localhost:8080/api/books?author=Frank
GET http://localhost:8080/api/books?category_id=1
```

**Get All Categories:**
```
GET http://localhost:8080/api/categories
```

**Catatan:**
- Routes ini bebas diakses tanpa token
- Siswa dapat melihat semua informasi buku dan category
- Di tahap development berikutnya, akan ada fitur peminjaman buku

</details>

---

## 📖 Konsep Penting yang Dipelajari

### `r` vs `c` di Gin

|               | Penjelasan                                                                                  |
| ------------- | ------------------------------------------------------------------------------------------- |
| `r` (Router)  | Dibuat sekali saat app start. Bertugas mendaftarkan semua route.                            |
| `c` (Context) | Dibuat per request. Bertugas melayani satu user — ambil data dari request & kirim response. |

### Cara Ambil Data dari Request

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

## Konfigurasi Environment

### Mengapa Environment Variables?

Environment variables memungkinkan aplikasi bekerja di berbagai environment (development, staging, production) **tanpa mengubah kode**.

### File `.env`

Semua konfigurasi sensitif disimpan di file `.env` yang **tidak di-commit** ke repository.

### Config Struct & LoadConfig()

Struct `Config` mengumpulkan semua environment variables menjadi satu object yang mudah diakses:

```go
type Config struct {
	Port                    string
	JWTSecret               string
	FrontendOrigin          string
	CloudinaryCloudName     string
	CloudinaryAPIKey        string
	CloudinaryAPISecret     string
}

func LoadConfig() (*Config, error) {
	godotenv.Load(".env") // Baca file .env
	return &Config{
		Port:   os.Getenv("PORT"),
		// ... field lainnya
	}, nil
}
```

## Database Schema & Migration

### AutoMigrate

GORM secara otomatis membuat/update table berdasarkan struct model:

```go
db.AutoMigrate(
	&models.Category{},
	&models.User{},
	&models.Book{},
	&models.Siswa{},
)
```

**Yang terjadi:**

- ✅ Table dibuat jika belum ada
- ✅ Kolom ditambah jika ada field baru di struct
- ✅ Index otomatis dibuat (contoh: `gorm:"uniqueIndex"`)

### Seed Data (Default Admin)

Saat aplikasi pertama kali jalan, `seedAdmin()` membuat Superadmin default:

```go
Email: superadmin@library.com
Password: superadmin123 (di-hash dengan bcrypt)
Role: superadmin
```

**Catatan:** Function hanya membuat admin jika table masih kosong — tidak akan duplicate.

## Integrasi Cloudinary

### Setup Koneksi

Cloudinary diinisialkan dengan URL yang berisi credentials:

```go
url := fmt.Sprintf(
	"cloudinary://%s:%s@%s",
	cfg.CloudinaryAPIKey,
	cfg.CloudinaryAPISecret,
	cfg.CloudinaryCloudName,
)
cld, _ := cloudinary.NewFromURL(url)
cld.Admin.Ping(context.Background()) // Test koneksi
```

### Upload Image

Untuk upload cover buku:

```go
func (h *AuthHandler) uploadToCloudinary(file *multipart.FileHeader, folder string) (string, string, error) {
	src, _ := file.Open()
	defer src.Close()

	result, _ := h.Cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder: folder,
	})

	return result.SecureURL, result.PublicID, nil
	// SecureURL → URL HTTPS yang bisa langsung ditampilkan
	// PublicID → ID untuk update/delete di kemudian hari
}
```

### Upload PDF

Untuk upload file PDF buku dengan validasi format:

```go
func (h *AuthHandler) uploadPDFToCloudinary(file *multipart.FileHeader) (string, string, error) {
	// Validasi extension
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".pdf") {
		return "", "", fmt.Errorf("file harus berformat PDF")
	}

	result, _ := h.Cloudinary.Upload.Upload(context.Background(), src, uploader.UploadParams{
		Folder:       "library/books/pdfs",
		ResourceType: "raw", // Untuk non-image file (PDF, doc, dll)
	})

	return result.SecureURL, result.PublicID, nil
}
```

**Perbedaan:**

- **Image upload** → `ResourceType` default (image), ditampilkan thumbnail preview
- **PDF upload** → `ResourceType: "raw"`, disimpan sebagai file tanpa processing

---

## 🚀 Roadmap Pengembangan

### ✅ Phase 1 - Core Setup & Database

- [x] Setup project & folder structure
- [x] Koneksi database SQLite + AutoMigrate
- [x] Environment configuration (.env)
- [x] Cloudinary integration

### ✅ Phase 2 - Authentication & User Management

- [x] CRUD User (Superadmin)
- [x] CRUD Siswa (Superadmin & Admin)
- [x] Autentikasi JWT (Login/Register)
- [x] Middleware auth untuk proteksi route
- [x] Role-based access control (Superadmin, Admin, Siswa)

### ✅ Phase 3 - Library Management

- [x] CRUD Buku (Admin)
- [x] CRUD Kategori (Admin)
- [x] Upload PDF & Cover buku (Cloudinary)
- [x] Search & filter buku (title, author, category)

### ✅ Phase 4 - User Interface

- [x] Admin dashboard
- [x] Tampilan daftar buku (card-based)
- [x] Tampilan detail buku
- [x] Tampilan manajemen siswa

### 🔄 Phase 5 - Book Reservation & Loan System (Coming Soon)

- [ ] Model `Reservation` untuk peminjaman buku
- [ ] CRUD Reservation (Siswa create, Admin approve/reject)
- [ ] Validasi stok buku saat peminjaman
- [ ] Durasi peminjaman & tanggal kembali otomatis
- [ ] Notifikasi peminjaman berhasil/ditolak
- [ ] History peminjaman per siswa
- [ ] Reminder pengembalian buku (otomatis/manual)
- [ ] Denda keterlambatan (optional)
- [ ] Status peminjaman: pending → approved → completed → returned

### 🎯 Phase 6 - Enhanced Features (Future)

- [ ] Rating & review buku dari siswa
- [ ] Wishlist buku untuk siswa
- [ ] Statistik peminjaman per buku & per siswa
- [ ] Export report peminjaman (PDF/Excel)
- [ ] Email notification untuk admin & siswa
- [ ] QR code untuk buku (untuk tracking)
- [ ] Mobile app version

---

> Project ini dibuat sebagai **Custom Project** untuk keperluan pembelajaran Golang.
> Deadline: **6 Maret 2026**
