# ERP Meble Backend

Backend aplikasi ERP Meble menggunakan Go (Golang) dengan dukungan real-time updates melalui WebSocket.

## 🚀 Quick Start

### 1. Setup Database

**Buat database di pgAdmin:**
```sql
CREATE DATABASE mebel_db;
```

**Buat file `.env` di folder `be/`:**
```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password  # GANTI INI!
DB_NAME=mebel_db
DB_SSLMODE=disable
```

**Test koneksi:**
```bash
go run test_db.go
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Run Server

```bash
go run cmd/server/main.go
```

Server akan berjalan di `http://localhost:8080`

---

## 📋 Setup Database Lengkap

Lihat dokumentasi:
- **`QUICK_SETUP_DB.md`** - Setup cepat (5 menit)
- **`SETUP_DATABASE.md`** - Panduan lengkap + troubleshooting

---

## 🏗️ Struktur Project

```
be/
├── cmd/server/main.go      # Entry point
├── internal/
│   ├── config/             # Konfigurasi
│   ├── database/           # Database connection
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Middleware (auth, dll)
│   ├── models/             # Database models
│   └── websocket/          # WebSocket hub
├── pkg/utils/              # Utilities
└── test_db.go              # Test koneksi database
```

---

## 📡 Endpoints

### Public
- `GET /health` - Health check
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/register` - Register

### Protected (require JWT)
- `GET /api/v1/users/me` - Get current user

### WebSocket
- `GET /ws` - WebSocket connection untuk real-time updates

---

## 🔧 Konfigurasi

Semua konfigurasi di file `.env`:
- Database connection
- JWT secret
- CORS settings
- Server port

---

## 📚 Dokumentasi

- **`SETUP_DATABASE.md`** - Setup database lengkap
- **`QUICK_SETUP_DB.md`** - Setup cepat
- **`FEATURE_LIST.md`** (di root) - Daftar fitur yang perlu dibuat
- **`DEVELOPMENT_ROADMAP.md`** (di root) - Roadmap development

---

## 🎯 Next Steps

Setelah database terhubung:
1. ✅ Database sudah connect
2. ✅ Auto-migration sudah jalan
3. ✅ Siap untuk mulai membuat fitur-fitur ERP

**Mulai dari:** Master Data Barang (lihat FEATURE_LIST.md)

---

## 🛠️ Tech Stack

- **Go 1.25+** - Programming language
- **Gin** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **Gorilla WebSocket** - Real-time updates
- **JWT** - Authentication

---

**Selamat coding!** 🚀

