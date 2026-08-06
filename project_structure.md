# Struktur Folder Proyek NusaGizi Backend

Proyek `nusagizi_be` menggunakan pola arsitektur *Clean Architecture* atau *Layered Architecture* yang umum digunakan pada framework pengembangan web di Golang. Struktur folder ini memisahkan secara jelas antara pengaturan routing, logika bisnis, dan interaksi ke database.

Berikut adalah penjelasan fungsi dari masing-masing folder dan file di dalam proyek:

## 📁 Direktori Utama (Root)
Direktori root merupakan tempat bagi file-file konfigurasi level proyek, skrip, dan folder utama.

*   **`cmd/`**
    Tempat file entry point aplikasi berada. Folder `cmd/api/` umumnya berisi `main.go` yang tugasnya mengatur inisialisasi awal seperti *setup database*, memanggil konfigurasi, melakukan *routing* endpoint, dan menjalankan HTTP server.
*   **`internal/`**
    Folder paling penting. Berisi seluruh *core logic* aplikasi. Di Golang, folder bernama `internal` bersifat privat; kode di dalamnya tidak bisa di-import (digunakan) oleh proyek luar. Penjelasan isi folder `internal` ada di bagian bawah.
*   **`migrations/`**
    Menyimpan skrip SQL (seperti `.up.sql` dan `.down.sql`) untuk pembuatan dan perubahan skema database (tabel, kolom). Ini digunakan untuk melacak riwayat versi skema database aplikasi.
*   **`scripts/`**
    Menyimpan skrip utilitas yang dapat dieksekusi, contohnya `migrate.ps1` (skrip PowerShell) untuk mempermudah developer menjalankan proses migrasi database.
*   **`tmp/`**
    Folder sementara yang biasanya dibuat otomatis oleh *live reload tool* (seperti Air) sebagai tempat meletakkan file *build executable* selama proses *development*.
*   **File Konfigurasi Utama:**
    *   **`.env`**: File berisi Environment Variables rahasia, seperti kredensial database (username, password), *secret key* JWT, dan port server. File ini **tidak boleh** masuk ke repositori Git.
    *   **`.air.toml`**: Konfigurasi untuk package `Air`, yaitu alat untuk mereload server Go secara otomatis (Hot-Reload) setiap kali ada kode yang disimpan (disave).
    *   **`go.mod` & `go.sum`**: File *package manager* bawaan Go yang bertugas mencatat dan mengunci versi dependensi / library pihak ketiga yang digunakan oleh proyek ini.
    *   **`.gitignore`**: Daftar file atau folder (seperti `.env` dan folder `tmp/`) yang akan diabaikan oleh Git.

---

## 📁 `internal/` (Application Layer)
Pemisahan struktur folder di `internal` memastikan setiap file memiliki satu tanggung jawab (*Single Responsibility Principle*):

*   **`config/`**
    Bertugas memuat (me-load) variabel lingkungan dari file `.env` dan menyimpannya ke dalam struktur data (*struct*) terpusat agar mudah diakses di seluruh aplikasi.
*   **`database/`**
    Berisi fungsi untuk menginisialisasi koneksi aplikasi ke database (misalnya melakukan *ping* dan mengatur *connection pool*).
*   **`models/`**
    Menyimpan struktur data (*structs*). Ini mencakup definisi tabel database (Entitas) dan struktur *Request/Response* JSON (DTO - Data Transfer Object).
*   **`repository/`** (Layer Data Access)
    Satu-satunya bagian yang bertugas berbicara langsung dengan Database. Isinya adalah kumpulan fungsi untuk mengeksekusi sintaks CRUD (Create, Read, Update, Delete). Tidak ada logika bisnis di layer ini.
*   **`services/`** (Layer Business Logic)
    Pusat atau "otak" aplikasi. Semua logika bisnis (seperti kalkulasi dan validasi khusus) berjalan di sini. Layer `services` tidak tahu menahu soal query database, ia hanya memanggil fungsi-fungsi yang disediakan oleh layer `repository`.
*   **`handlers/`** (Layer Controller / Presentation)
    Bertugas menerima permintaan (Request) dari pengguna (klien/frontend). Mengambil data dari URL, parameter, maupun Body JSON, lalu meneruskannya ke `services`. Jika `services` selesai memproses, `handlers` akan merangkum datanya untuk dikirim kembali (*Response*) sebagai JSON ke klien.
*   **`middleware/`**
    Fungsi penengah yang akan dieksekusi sebelum request mencapai `handlers`. Biasanya digunakan untuk hal-hal yang repetitif, seperti: validasi keamanan (Authorization/JWT), *CORS*, atau *Logging* aktivitas klien.
*   **`auth/`**
    Modul yang mengurus mekanisme keamanan khusus seperti *generate* token JWT, verifikasi password, dan *hash* password.
*   **`who/`**
    Modul khusus (mengingat nama proyek ini NusaGizi). Kemungkinan besar digunakan untuk perhitungan status gizi, berat/tinggi badan ideal, Z-Score, maupun pengolahan data antropometri menggunakan referensi standar dari World Health Organization (WHO).

---

## 🔄 Alur Request Sederhana
Jika diringkas, saat Frontend (klien) meminta data, alur perjalanannya adalah:

1. **Request** dari klien ditangkap oleh rute di **`cmd/api/main.go`**.
2. Masuk ke **`middleware`** (misal: "Apakah tokennya valid?").
3. Masuk ke **`handlers`** (membaca data yang dikirim klien).
4. Masuk ke **`services`** (melakukan kalkulasi/aturan bisnis berdasarkan data klien).
5. Masuk ke **`repository`** (menyimpan atau mengambil data gizi dari database).
6. **Kembali lagi** dari *Repository ➡️ Services ➡️ Handlers ➡️ Middleware ➡️ Klien* sebagai **Response JSON**.
