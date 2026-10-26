# 💰 Expense Tracker

Aplikasi berbasis web untuk mencatat, memantau, mengelompokkan, dan mengelola transaksi keuangan pemasukan serta pengeluaran bulanan secara personal dan terstruktur.

---

## 🚀 Fitur Utama

1. **Autentikasi Aman:** Registrasi dan Login dengan enkripsi kata sandi `bcrypt` dan sesi JWT yang disimpan di dalam **HttpOnly Cookie** (`SameSite=Lax`) guna mencegah serangan XSS.
2. **Landing Page Publik:** Halaman selamat datang dan pengenalan aplikasi bagi calon pengguna baru.
3. **Pencatatan Transaksi Cepat:** Tombol cepat *"Tambah Transaksi"* untuk mencatat nominal (Rupiah integer), kategori, tanggal, dan deskripsi.
4. **Manajemen Kategori:** Kategori bawaan sistem (*read-only*) dan kemampuan bagi pengguna untuk menambah, mengedit, atau menghapus kategori custom.
5. **Dashboard Keuangan:** Menampilkan kartu ringkasan total pemasukan, total pengeluaran, dan saldo bersih bulanan secara real-time.
6. **Visualisasi Grafik:** Grafik perbandingan pemasukan vs pengeluaran serta proporsi kategori pengeluaran interaktif menggunakan `Chart.js`.
7. **Filter Riwayat Transaksi:** Menyaring riwayat berdasarkan periode bulan/tahun, jenis transaksi (pemasukan/pengeluaran), dan kategori.
8. **CRUD Transaksi Lengkap:** Menambah, melihat rincian, mengedit, dan menghapus transaksi dengan verifikasi hak kepemilikan akun.

---

## 🛠️ Tumpukan Teknologi (Tech Stack)

- **Backend:** Go (Golang) v1.22+
- **Web Framework:** [Gin Web Framework](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/)
- **Database:** PostgreSQL (Neon Serverless PostgreSQL)
- **Frontend:** HTML5, CSS3, JavaScript murni (Vanilla JS) & Chart.js (via CDN)
- **Keamanan:** JWT (`golang-jwt/jwt/v5`), Bcrypt (`golang.org/x/crypto/bcrypt`)

---

## 📁 Struktur Direktori

```text
expense-tracker/
├── .env.example                 # Template konfigurasi environment (aman di Git)
├── .gitignore                   # Mengabaikan file sensitif (.env, bin, cookies)
├── README.md                    # Dokumentasi lengkap proyek
├── go.mod                       # Definisi modul dan dependensi Go
├── go.sum                       # Checksum dependensi Go
├── main.go                      # Entry point server Gin & rute aplikasi
├── config/
│   └── database.go              # Koneksi GORM ke PostgreSQL Neon & seeder kategori default
├── controllers/
│   ├── auth_controller.go       # Handler Registrasi, Login, Logout, & Profil
│   ├── category_controller.go   # Handler CRUD Kategori (default & kustom)
│   └── transaction_controller.go# Handler CRUD Transaksi, ringkasan saldo, & filter
├── middleware/
│   └── auth_middleware.go       # Middleware validasi JWT dari HttpOnly Cookie
├── models/
│   ├── user.go                  # Struct User & method hash bcrypt
│   ├── category.go              # Struct Category (bawaan & kustom)
│   └── transaction.go           # Struct Transaction (nominal integer Rupiah)
└── static/                      # Antarmuka frontend (Vanilla HTML/CSS/JS)
    ├── css/
    │   └── style.css            # Styling responsif antarmuka
    ├── js/
    │   ├── auth.js              # Logika autentikasi & sesi
    │   ├── app.js               # Logika Dashboard, Transaksi, & Grafik
    │   └── categories.js        # Logika Modal & Kategori Custom
    ├── landing.html             # Halaman publik pengenalan aplikasi
    ├── login.html               # Halaman masuk
    ├── register.html            # Halaman daftar
    └── index.html               # Halaman Dashboard Keuangan
```

---

## ⚙️ Panduan Menjalankan Proyek

### 1. Prasyarat
- [Go](https://go.dev/dl/) versi 1.21 atau yang lebih baru.
- Akun dan database [Neon PostgreSQL](https://neon.tech/).

### 2. Konfigurasi Lingkungan (`.env`)
Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Buka file `.env` dan masukkan konfigurasi Anda:
```dotenv
APP_PORT=8080
DATABASE_URL=postgres://<user>:<password>@<ep-domain>.neon.tech/neondb?sslmode=require
JWT_SECRET=rahasia-kunci-jwt-sangat-aman-12345
```
> **Catatan Keamanan:** File `.env` sudah masuk ke dalam `.gitignore` dan tidak boleh di-commit ke repositori Git. Nilai rahasia tidak pernah ditampilkan di log aplikasi maupun respons API.

### 3. Unduh Dependensi
```bash
go mod tidy
```

### 4. Jalankan Aplikasi
```bash
go run main.go
```

Aplikasi akan berjalan di:
- **Landing Page:** [http://localhost:8080/](http://localhost:8080/)
- **Halaman Login:** [http://localhost:8080/login](http://localhost:8080/login)
- **Halaman Daftar:** [http://localhost:8080/register](http://localhost:8080/register)
- **Dashboard Keuangan:** [http://localhost:8080/dashboard](http://localhost:8080/dashboard)

---

## 🧪 Panduan Pengujian Aplikasi (Testing Guide)

Pengujian dapat dilakukan melalui **Antarmuka Web (Browser)** untuk pengujian alur interaktif pengguna (*End-to-End*) maupun melalui **Terminal (cURL)** untuk pengujian endpoint REST API.

### A. Pengujian Melalui Browser (UI / E2E)

1. **Pengujian Halaman Publik & Pendaftaran Akun:**
   - Akses `http://localhost:8080/`.
   - Klik tombol **"Mulai Sekarang"** atau buka `http://localhost:8080/register`.
   - Masukkan Username (min. 3 karakter), Email valid, Kata Sandi (min. 6 karakter), dan Konfirmasi Sandi.
   - *Verifikasi:* Jika konfirmasi sandi tidak sesuai, muncul peringatan instan. Setelah sukses, Anda otomatis dialihkan ke halaman login.

2. **Pengujian Login & Proteksi Sesi (HttpOnly Cookie):**
   - Di halaman `http://localhost:8080/login`, masukkan kredensial yang baru didaftarkan.
   - *Verifikasi:* Setelah berhasil masuk, browser otomatis diarahkan ke `http://localhost:8080/dashboard`.
   - Buka tab DevTools Browser (`F12` -> *Application* -> *Cookies*): periksa bahwa cookie `auth_token` terpasang dengan atribut `HttpOnly` aktif untuk mencegah serangan XSS.

3. **Pengujian Kategori Custom & Proteksi Kategori Bawaan:**
   - Di Dashboard, klik tombol **"🏷️ Kelola Kategori"**.
   - Tambahkan kategori baru, misalnya: nama `Langganan Hosting`, tipe `Pengeluaran`.
   - *Verifikasi:* Kategori custom bertambah pada daftar dengan badge `Custom` serta tombol `Edit` dan `Hapus`. Kategori bawaan sistem (seperti `Makanan & Minuman`, `Gaji & Upah`) berstatus `Read-only` dan diproteksi sehingga tidak dapat diubah/dihapus.

4. **Pengujian Pencatatan Transaksi & Kalkulasi Real-Time:**
   - Klik tombol **"➕ Tambah Transaksi"**.
   - Catat transaksi pemasukan: pilih tab **Pemasukan**, isi nominal `10.000.000`, pilih kategori `Gaji & Upah`, tanggal hari ini, lalu simpan.
   - Catat transaksi pengeluaran: pilih tab **Pengeluaran**, isi nominal `150.000`, pilih kategori `Makanan & Minuman`, tanggal hari ini, lalu simpan.
   - *Verifikasi:* 
     - Kartu **Total Pemasukan** bertambah `Rp 10.000.000`.
     - Kartu **Total Pengeluaran** bertambah `Rp 150.000`.
     - Kartu **Sisa Saldo Bersih** otomatis menghitung selisihnya (`Rp 9.850.000`).
     - Grafik batang dan grafik donut Chart.js langsung memvisualisasikan data arus kas.
     - Riwayat tabel menampilkan baris baru lengkap dengan format Rupiah.

5. **Pengujian Filter Transaksi:**
   - Ubah dropdown filter di atas tabel menjadi **"Hanya Pengeluaran"** atau pilih kategori tertentu.
   - *Verifikasi:* Baris tabel menyaring hanya data yang sesuai dengan filter yang aktif.
   - Ganti pemilih bulan/tahun di bagian atas dashboard untuk memverifikasi isolasi periode laporan bulanan.

6. **Pengujian Edit & Hapus Transaksi:**
   - Klik tombol **"✏️ Edit"** pada salah satu baris transaksi, ubah nominal atau keterangan, lalu klik simpan.
   - Klik tombol **"🗑️ Hapus"** dan setujui konfirmasi browser.
   - *Verifikasi:* Kartu ringkasan saldo, visualisasi grafik, dan tabel langsung ter-update secara otomatis.

7. **Pengujian Logout & Route Guard:**
   - Klik tombol **"Keluar"** di navbar.
   - Coba akses kembali `http://localhost:8080/dashboard` secara langsung di browser.
   - *Verifikasi:* Sesi telah terhapus dan Anda otomatis dialihkan kembali ke halaman login.

---

### B. Pengujian REST API Melalui Terminal (cURL)

Gunakan perintah cURL berikut secara berurutan di terminal (Command Prompt / PowerShell / Git Bash):

#### 1. Uji Health Check
```bash
curl http://localhost:8080/api/health
```

#### 2. Uji Registrasi Akun
```bash
curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d "{\"username\": \"tester\", \"email\": \"tester@example.com\", \"password\": \"password123\"}"
```

#### 3. Uji Login (Menyimpan Cookie HttpOnly ke file `cookies.txt`)
```bash
curl -i -X POST http://localhost:8080/api/auth/login -c cookies.txt -H "Content-Type: application/json" -d "{\"identifier\": \"tester@example.com\", \"password\": \"password123\"}"
```

#### 4. Uji Verifikasi Sesi Profil (`/api/auth/me`)
```bash
curl -b cookies.txt http://localhost:8080/api/auth/me
```

#### 5. Uji Tambah Kategori Custom
```bash
curl -X POST http://localhost:8080/api/categories -b cookies.txt -H "Content-Type: application/json" -d "{\"name\": \"Tagihan Internet\", \"type\": \"expense\"}"
```

#### 6. Uji Proteksi Kategori Bawaan (Dilarang Dihapus)
```bash
curl -i -X DELETE http://localhost:8080/api/categories/1 -b cookies.txt
```
*(Status respons yang diharapkan: `403 Forbidden`)*

#### 7. Uji Tambah Transaksi Pengeluaran
```bash
curl -X POST http://localhost:8080/api/transactions -b cookies.txt -H "Content-Type: application/json" -d "{\"category_id\": 1, \"amount\": 50000, \"type\": \"expense\", \"date\": \"2026-10-26\", \"description\": \"Makan Siang\"}"
```

#### 8. Uji Ambil Ringkasan Keuangan (Dashboard Summary)
```bash
curl -b cookies.txt "http://localhost:8080/api/transactions/summary?month=10&year=2026"
```

#### 9. Uji Logout
```bash
curl -i -X POST http://localhost:8080/api/auth/logout -b cookies.txt -c cookies.txt
```

---

## 🧪 Panduan Pengujian Aplikasi (Testing Guide)

Pengujian dapat dilakukan melalui **Antarmuka Web (Browser)** untuk pengujian alur interaktif pengguna (*End-to-End*) maupun melalui **Terminal (cURL)** untuk pengujian endpoint REST API.

### A. Pengujian Melalui Browser (UI / E2E)

1. **Pengujian Halaman Publik & Pendaftaran Akun:**
   - Akses `http://localhost:8080/`.
   - Klik tombol **"Mulai Sekarang"** atau buka `http://localhost:8080/register`.
   - Masukkan Username (min. 3 karakter), Email valid, Kata Sandi (min. 6 karakter), dan Konfirmasi Sandi.
   - *Verifikasi:* Jika konfirmasi sandi tidak sesuai, muncul peringatan instan. Setelah sukses, Anda otomatis dialihkan ke halaman login.

2. **Pengujian Login & Proteksi Sesi (HttpOnly Cookie):**
   - Di halaman `http://localhost:8080/login`, masukkan kredensial yang baru didaftarkan.
   - *Verifikasi:* Setelah berhasil masuk, browser otomatis diarahkan ke `http://localhost:8080/dashboard`.
   - Buka tab DevTools Browser (`F12` -> *Application* -> *Cookies*): periksa bahwa cookie `auth_token` terpasang dengan atribut `HttpOnly` aktif untuk mencegah serangan XSS.

3. **Pengujian Kategori Custom & Proteksi Kategori Bawaan:**
   - Di Dashboard, klik tombol **"🏷️ Kelola Kategori"**.
   - Tambahkan kategori baru, misalnya: nama `Langganan Hosting`, tipe `Pengeluaran`.
   - *Verifikasi:* Kategori custom bertambah pada daftar dengan badge `Custom` serta tombol `Edit` dan `Hapus`. Kategori bawaan sistem (seperti `Makanan & Minuman`, `Gaji & Upah`) berstatus `Read-only` dan diproteksi sehingga tidak dapat diubah/dihapus.

4. **Pengujian Pencatatan Transaksi & Kalkulasi Real-Time:**
   - Klik tombol **"➕ Tambah Transaksi"**.
   - Catat transaksi pemasukan: pilih tab **Pemasukan**, isi nominal `10.000.000`, pilih kategori `Gaji & Upah`, tanggal hari ini, lalu simpan.
   - Catat transaksi pengeluaran: pilih tab **Pengeluaran**, isi nominal `150.000`, pilih kategori `Makanan & Minuman`, tanggal hari ini, lalu simpan.
   - *Verifikasi:* 
     - Kartu **Total Pemasukan** bertambah `Rp 10.000.000`.
     - Kartu **Total Pengeluaran** bertambah `Rp 150.000`.
     - Kartu **Sisa Saldo Bersih** otomatis menghitung selisihnya (`Rp 9.850.000`).
     - Grafik batang dan grafik donut Chart.js langsung memvisualisasikan data arus kas.
     - Riwayat tabel menampilkan baris baru lengkap dengan format Rupiah.

5. **Pengujian Filter Transaksi:**
   - Ubah dropdown filter di atas tabel menjadi **"Hanya Pengeluaran"** atau pilih kategori tertentu.
   - *Verifikasi:* Baris tabel menyaring hanya data yang sesuai dengan filter yang aktif.
   - Ganti pemilih bulan/tahun di bagian atas dashboard untuk memverifikasi isolasi periode laporan bulanan.

6. **Pengujian Edit & Hapus Transaksi:**
   - Klik tombol **"✏️ Edit"** pada salah satu baris transaksi, ubah nominal atau keterangan, lalu klik simpan.
   - Klik tombol **"🗑️ Hapus"** dan setujui konfirmasi browser.
   - *Verifikasi:* Kartu ringkasan saldo, visualisasi grafik, dan tabel langsung ter-update secara otomatis.

7. **Pengujian Logout & Route Guard:**
   - Klik tombol **"Keluar"** di navbar.
   - Coba akses kembali `http://localhost:8080/dashboard` secara langsung di browser.
   - *Verifikasi:* Sesi telah terhapus dan Anda otomatis dialihkan kembali ke halaman login.

---

### B. Pengujian REST API Melalui Terminal (cURL)

Gunakan perintah cURL berikut secara berurutan di terminal (Command Prompt / PowerShell / Git Bash):

#### 1. Uji Health Check
```bash
curl http://localhost:8080/api/health
```

#### 2. Uji Registrasi Akun
```bash
curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d "{\"username\": \"tester\", \"email\": \"tester@example.com\", \"password\": \"password123\"}"
```

#### 3. Uji Login (Menyimpan Cookie HttpOnly ke file `cookies.txt`)
```bash
curl -i -X POST http://localhost:8080/api/auth/login -c cookies.txt -H "Content-Type: application/json" -d "{\"identifier\": \"tester@example.com\", \"password\": \"password123\"}"
```

#### 4. Uji Verifikasi Sesi Profil (`/api/auth/me`)
```bash
curl -b cookies.txt http://localhost:8080/api/auth/me
```

#### 5. Uji Tambah Kategori Custom
```bash
curl -X POST http://localhost:8080/api/categories -b cookies.txt -H "Content-Type: application/json" -d "{\"name\": \"Tagihan Internet\", \"type\": \"expense\"}"
```

#### 6. Uji Proteksi Kategori Bawaan (Dilarang Dihapus)
```bash
curl -i -X DELETE http://localhost:8080/api/categories/1 -b cookies.txt
```
*(Status respons yang diharapkan: `403 Forbidden`)*

#### 7. Uji Tambah Transaksi Pengeluaran
```bash
curl -X POST http://localhost:8080/api/transactions -b cookies.txt -H "Content-Type: application/json" -d "{\"category_id\": 1, \"amount\": 50000, \"type\": \"expense\", \"date\": \"2026-10-26\", \"description\": \"Makan Siang\"}"
```

#### 8. Uji Ambil Ringkasan Keuangan (Dashboard Summary)
```bash
curl -b cookies.txt "http://localhost:8080/api/transactions/summary?month=10&year=2026"
```

#### 9. Uji Logout
```bash
curl -i -X POST http://localhost:8080/api/auth/logout -b cookies.txt -c cookies.txt
```

---

## 🔒 Praktik Keamanan yang Diterapkan

1. **HttpOnly Cookie:** Token JWT disimpan di cookie dengan atribut `HttpOnly`, sehingga skrip JavaScript di browser tidak dapat membacanya, mencegah risiko pencurian token melalui serangan Cross-Site Scripting (XSS).
2. **SameSite=Lax:** Cookie dikonfigurasi dengan `SameSite=Lax` untuk melindungi dari serangan CSRF.
3. **Password Hashing:** Kata sandi di-hash menggunakan algoritma `bcrypt` dengan salt dinamis sebelum disimpan ke basis data. Field password pada model struct diabaikan saat serialisasi JSON (`json:"-"`).
4. **Isolasi Data (Multi-tenancy per User):** Setiap kueri transaksi dan kategori custom selalu memvalidasi `user_id` dari token JWT sesi pengguna, sehingga pengguna tidak dapat membaca, mengedit, atau menghapus data milik pengguna lain.
5. **Proteksi Kategori Sistem:** Kategori bawaan (`user_id IS NULL`) dilindungi di sisi backend sehingga tidak dapat dimodifikasi atau dihapus oleh pengguna biasa.
