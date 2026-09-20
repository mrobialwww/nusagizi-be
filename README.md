# NusaGizi Backend API

NusaGizi Backend is a RESTful API built in Go using the Gin framework. It serves as the core backend for the NusaGizi application, an integrated platform designed to monitor and support children's nutrition, growth, and development, connecting mothers with caregivers and providing AI-powered nutrition recommendations.

## 🚀 Tech Stack

- **Language:** Go (1.25.0+)
- **Web Framework:** [Gin Web Framework](https://gin-gonic.com/)
- **Database:** PostgreSQL (using `pgx/v5`)
- **Authentication:** Auth0 (JWT validation and Machine-to-Machine)
- **Cloud Storage:** Cloudflare R2 (S3-compatible via `aws-sdk-go-v2`)
- **Environment Management:** `godotenv`

## 📁 Project Structure

Following standard Go project layout:

```text
nusagizi_be/
├── cmd/
│   └── api/
│       └── main.go         # Application entry point and route registrations
├── internal/
│   ├── auth/               # Auth0 JWT validation and verifiers
│   ├── config/             # Environment configuration loader
│   ├── database/           # PostgreSQL connection pool setup
│   ├── handlers/           # HTTP controllers/handlers
│   ├── infrastructure/     # External services integrations (e.g., Cloudflare R2)
│   ├── middleware/         # Gin middlewares (Auth, logging, etc.)
│   ├── models/             # Domain entities and DTOs (Data Transfer Objects)
│   ├── repository/         # Database access layer
│   ├── services/           # Business logic layer
│   ├── utils/              # Helper functions
│   └── who/                # (Specialized modules/interfaces if any)
├── migrations/             # Database migration scripts
├── scripts/                # Utility scripts
├── seed/                   # SQL seed files for dummy data
├── endpoint.md             # Detailed endpoint documentation
└── fix-image.md / r2...    # Architectural decision records and markdown guides
```

## 🧩 Core Modules

The API exposes the following main modules:

1. **Auth & User Profile:** Setup Mother and Caregiver profiles, update user details, confirm email via OTP.
2. **Child Profile:** Manage children profiles associated with a mother.
3. **Child Growth:** Log and analyze child growth metrics.
4. **Child Development (KPSP):** Process KPSP assessments, checklist milestones, and serve developmental recommendations.
5. **Child Nutrition & Menu:** Fetch today's nutrition reports, bookmarked recipes, daily shopping list, and an AI-powered menu generator integration.
6. **Photos & Social Contacts:** Manage caregiver contacts, share child photos, and manage photo galleries.
7. **Caregiver Engagements:** Link and unlink caregivers to mothers' accounts.
8. **Medical Notes:** Create and manage customized medical records/notes for children.
9. **Dashboard:** Unified dashboard summary endpoint for quick app overviews.
10. **Notifications:** Fetch system and user-specific notifications.
11. **Check-in:** QR Check-in system between caregivers and mothers.

## 🛠️ Prerequisites

To run this project locally, ensure you have:

- [Go](https://golang.org/doc/install) (version 1.25.0 or newer)
- PostgreSQL database (or access to Neon DB as currently configured)
- Auth0 Account with a configured Application/API
- Cloudflare Account with an R2 Bucket set up

## ⚙️ Configuration (.env)

Create a `.env` file in the root directory based on the following template:

```env
# Database
DATABASE_URL=postgresql://user:password@host:port/dbname?sslmode=require...
PORT=3000

# External Services
AI_SERVICE_BASE_URL=https://nusagizi-ai-production.up.railway.app

# Auth0 Configuration
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_AUDIENCE=https://api.nusagizi.com
CLIENT_ID=your-m2m-client-id
CLIENT_SECRET=your-m2m-client-secret
GRANT_TYPE=client_credentials
AUTH0_M2M_CLIENT_ID=your-m2m-client-id
AUTH0_M2M_CLIENT_SECRET=your-m2m-client-secret

# Auth0 Roles
AUTH0_ROLE_ID_MOTHER=role-id-mother
AUTH0_ROLE_ID_CAREGIVER=role-id-caregiver
AUTH0_ROLE_ID_DOCTOR=role-id-doctor

# Cloudflare R2 Settings
R2_ACCOUNT_ID=your-r2-account-id
R2_ACCESS_KEY_ID=your-r2-access-key
R2_SECRET_ACCESS_KEY=your-r2-secret-key
R2_BUCKET=your-bucket-name
R2_PUBLIC_URL=https://your-custom-r2-domain.com
```

## 🛠️ Spesifikasi Lingkungan Pengujian (Testing Environment)

Untuk memastikan aplikasi berjalan dengan optimal, kami merekomendasikan lingkungan pengujian (testing environment) sebagai berikut:

- **Sistem Operasi:** Windows 10/11, macOS, atau Linux (Ubuntu 20.04+ direkomendasikan)
- **Go Version:** Go v1.25.0 atau terbaru
- **Database:** PostgreSQL v15+ (Local) atau terkoneksi ke NeonDB Cloud
- **Koneksi Internet:** Diperlukan (untuk integrasi external seperti Auth0, Cloudflare R2, dan sistem AI)
- **Tools Tambahan:**
    - [Postman](https://www.postman.com/) atau [Insomnia](https://insomnia.rest/) untuk _API Testing_
    - [Air](https://github.com/cosmtrek/air) (opsional) untuk fitur _hot-reload_ saat masa development

## 🔧 Panduan Instalasi

Ikuti langkah-langkah di bawah ini untuk menginstal dan menyiapkan proyek di mesin lokal Anda:

1. **Clone repository:**
   Buka terminal atau command prompt, lalu jalankan perintah berikut:

    ```bash
    git clone <repository-url>
    cd nusagizi_be
    ```

2. **Unduh Dependensi:**
   Pastikan Anda sudah menginstal Go, lalu unduh semua module yang dibutuhkan:

    ```bash
    go mod download
    ```

    _Jika ada package yang belum tersinkronisasi, Anda juga dapat menjalankan `go mod tidy`._

3. **Atur Environment Variables:**
    - Duplikat file `.env.example` (jika ada) atau buat file `.env` baru di _root_ direktori.
    - Isi kredensial API dan database di file `.env` berdasarkan spesifikasi yang dijelaskan pada bagian **Configuration** di atas.

## 🚀 Cara Menjalankan Aplikasi

Anda dapat menjalankan _NusaGizi Backend_ dengan 2 cara, bergantung pada kebutuhan Anda (Development atau Production).

### Cara 1: Menggunakan Go Run (Standar)

Sangat cocok untuk pengujian biasa atau environment production.

```bash
go run cmd/api/main.go
```

### Cara 2: Menggunakan Air (Development Mode)

Cocok digunakan saat masa _development_ karena mendukung _live-reloading_ (aplikasi me-restart otomatis ketika ada perubahan kode).

1. Pastikan Anda menginstal `air` terlebih dahulu: `go install github.com/cosmtrek/air@latest`
2. Jalankan perintah:

```bash
air
```

### Cara 3: Build & Eksekusi Binari

Apabila ingin dites sebagai sebuah binary _executable_ (contoh: sebelum deploy server):

```bash
go build -o main.exe cmd/api/main.go
./main.exe
```

### Verifikasi Server

Secara _default_, _backend server_ akan berjalan di alamat `http://localhost:3000`. Anda bisa membukanya melalui browser atau Postman.
Jika berhasil, maka Anda akan menerima _response_ seperti berikut:

```json
{
    "message": "NusaGizi API is running!",
    "status": "success",
    "database": "connected"
}
```

## 📚 API Documentation

Detail spesifikasi _endpoint_, _request_, dan _response_ telah kami dokumentasikan secara komprehensif pada file [endpoint.md](./endpoint.md). Kunjungi file tersebut untuk melihat struktur pertukaran data API lebih lanjut.

---

_Built for ProjectLomba - NusaGizi._
