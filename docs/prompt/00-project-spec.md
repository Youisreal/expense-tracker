## 1. Informasi Aplikasi

**Nama aplikasi:**

`Expense Tracker`

**Deskripsi aplikasi:**

Aplikasi berbasis web untuk membantu pengguna mencatat, memantau, mengelompokkan, dan mengelola transaksi pemasukan serta pengeluaran keuangan bulanan secara personal dan terstruktur.

**Target pengguna:**

Individu, pekerja, atau mahasiswa yang ingin mengelola dan melacak keuangan pribadi mereka setiap bulan.

**Masalah yang ingin diselesaikan:**

Kesulitan dalam melacak arus kas (cash flow) harian dan bulanan secara manual, kurangnya fleksibilitas pengkategorian transaksi, sehingga pengguna kesulitan mengetahui ke mana saja uang mereka pergi dan berapa sisa saldo yang tersedia.

---

## 2. Kebutuhan Aplikasi

Fitur utama aplikasi:

1. **Autentikasi Pengguna**: Registrasi dan Login pengguna dengan sistem keamanan JWT yang disimpan dalam Cookie (HttpOnly).
2. **Landing Page**: Halaman selamat datang / pengenalan aplikasi Expence Tracker bagi calon pengguna baru.
3. **Pencatatan Transaksi**: Menambahkan catatan transaksi baru (pemasukan atau pengeluaran) melalui tombol cepat **"Tambah Transaksi"** beserta nominal, kategori, tanggal, dan keterangan.
4. **Manajemen Kategori Custom**: Kategori bawaan (misal: Makanan, Tagihan, Hiburan) serta kemampuan bagi pengguna untuk menambah, mengedit, atau menghapus kategori custom sesuai kebutuhan mereka.
5. **Filter Transaksi**: Menyaring daftar transaksi berdasarkan kategori (misal: Makanan, Tagihan, Hiburan, atau kategori buatan sendiri) dan periode/bulan.
6. **Dashboard Keuangan**: Menampilkan riwayat transaksi bulanan serta ringkasan total pemasukan, total pengeluaran, dan saldo bersih.
7. **Grafik Perbandingan**: Menampilkan grafik perbandingan pemasukan dan pengeluaran.
8. **CRUD Transaksi**: Mengedit dan menghapus catatan transaksi yang sudah dibuat.
9. **Logout**: Mengakhiri sesi pengguna dan membersihkan cookie JWT.

**Halaman aplikasi:**

- Landing Page (Pengenalan Aplikasi & Call-to-Action)
- Login
- Register
- Main Page — Dashboard Keuangan (Pencatatan, Filter Kategori, & Riwayat Transaksi)

---

## 3. Pilihan Teknologi

### A. Backend Framework

- [ ] 1. Go Standard Library (`net/http`)
- [x] 2. Gin
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `2`

### B. ORM (Object-Relational Mapping)

- [ ] 1. Tidak menggunakan ORM (SQL manual)
- [x] 2. GORM
- [ ] 3. Lainnya sebutkan

**Pilihan saya:** `2`

**Ketentuan:** Menggunakan GORM untuk pemetaan objek relasional ke database PostgreSQL (Neon).

### C. Database

- [x] 1. PostgreSQL (Neon Cloud Database)
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** `1`

### D. Frontend

- [x] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** `1`

### E. Sistem Autentikasi

- [ ] 1. Session
- [x] 2. JWT

**Pilihan saya:** `2`

**Ketentuan:** Gunakan JSON Web Token (JWT) yang dihasilkan oleh backend untuk mengelola status login pengguna.

### F. Penyimpanan Sesi atau Token

- [x] 1. Cookie
- [ ] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** `1`

**Ketentuan:** Simpan token JWT di dalam HttpOnly Cookie yang aman untuk mencegah serangan XSS.

### G. API

- [x] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** `1`

**Ketentuan:** Backend menyediakan REST API untuk autentikasi, manajemen kategori, dan pengolahan transaksi, diakses oleh frontend menggunakan `fetch()`.

### H. Version Control

- [x] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** `1`

---

## 4. Detail Autentikasi & Otorisasi

**Registrasi akun:** Ya

**Login:** Ya

**Logout:** Ya

**Role pengguna:** User saja. Tidak ada role Admin.

**Aturan akses:**

1. Halaman **Landing Page**, **Login**, dan **Register** terbuka untuk umum (publik).
2. Pengguna harus login untuk mengakses Dashboard Keuangan, fitur pencatatan, dan kelola kategori.
3. Setiap transaksi dan kategori custom harus memiliki `user_id` sebagai pemiliknya (kategori default/bawaan disiapkan secara publik atau otomatis di-generate per user).
4. Backend harus memverifikasi kepemilikan data pada setiap operasi CRUD transaksi dan kategori.
5. Password wajib di-hash menggunakan algoritma aman (`bcrypt`) sebelum disimpan ke database.
6. Setelah logout, cookie JWT akan dibersihkan oleh server.

---

## 5. Konfigurasi Environment dan Database

Contoh `.env.example`:

```dotenv
APP_PORT=8080
DATABASE_URL=postgres://user:password@ep-sample-pooler.region.neon.tech/neondb?sslmode=require
JWT_SECRET=replace-with-a-secure-random-secret

```

---

## 6. Pembaruan Struktur Folder Proyek

```text
expence-tracker/
├── .env.example
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── main.go                  # Entry point aplikasi
├── config/                  # Koneksi DB PostgreSQL Neon & GORM
│   └── database.go
├── controllers/             # Logika bisnis
│   ├── auth_controller.go
│   ├── category_controller.go     # Logika Kategori (Default & Custom)
│   └── transaction_controller.go  # Logika Transaksi & Filter
├── models/                  # Struct GORM
│   ├── user.go
│   ├── category.go          # Struct Kategori Custom/Default
│   └── transaction.go       # Struct Transaksi
├── middleware/              # Middleware JWT Auth
│   └── auth_middleware.go
└── static/                  # Frontend HTML/CSS/JS murni
    ├── css/
    │   └── style.css
    ├── js/
    │   ├── app.js           # Logika Dashboard & Filter
    │   └── categories.js    # Logika Modal & Kategori Custom
    ├── landing.html         # Landing Page Pengenalan Expence Tracker
    ├── login.html
    ├── register.html
    └── index.html           # Main Dashboard (Tombol cepat & Filter)

```

---

### Penjelasan Singkat Folder & File:

- **`cmd/server/main.go`**: Titik masuk utama aplikasi untuk menjalankan server Gin, memuat konfigurasi environment, melakukan koneksi database, dan mendaftarkan rute API.
- **`internal/config/`**: Mengatur dan menginisialisasi koneksi GORM ke database Neon PostgreSQL menggunakan variabel lingkungan (`DATABASE_URL`).
- **`internal/models/`**: Berisi definisi struktur data (struct) GORM untuk entitas pengguna (`User`), kategori (`Category`), dan transaksi (`Transaction`) guna keperluan migrasi otomatis tabel.
- **`internal/handlers/`**: Berisi logika endpoint (controller) untuk menangani autentikasi, manajemen kategori custom, serta CRUD dan filter transaksi keuangan.
- **`internal/middleware/`**: Berisi middleware keamanan untuk memverifikasi validitas token JWT yang dikirimkan melalui Cookie pada endpoint privat.
- **`web/`**: Berisi seluruh aset antarmuka pengguna frontend (HTML, CSS, dan JavaScript murni) termasuk Landing Page, halaman Login, Register, dan Dashboard Utama.

---
