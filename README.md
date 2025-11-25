# ERP Meble Backend

Backend aplikasi ERP Meble menggunakan Go (Golang) dengan dukungan real-time updates melalui WebSocket.

## Framework dan Package yang Digunakan

### Framework Utama
- **Gin** (`github.com/gin-gonic/gin`) - Web framework yang ringan dan cepat untuk HTTP server
- **GORM** (`gorm.io/gorm`) - ORM untuk database operations
- **PostgreSQL Driver** (`gorm.io/driver/postgres`) - Driver untuk PostgreSQL database

### Real-time Communication
- **Gorilla WebSocket** (`github.com/gorilla/websocket`) - Library untuk WebSocket connections, memungkinkan real-time updates ke semua client yang terhubung

### Authentication & Security
- **JWT** (`github.com/golang-jwt/jwt/v5`) - JSON Web Tokens untuk authentication
- **Bcrypt** (`golang.org/x/crypto/bcrypt`) - Password hashing

### Validation & Utilities
- **Validator** (`github.com/go-playground/validator/v10`) - Struct validation
- **Godotenv** (`github.com/joho/godotenv`) - Environment variables management
- **CORS** (`github.com/gin-contrib/cors`) - Cross-Origin Resource Sharing middleware

## Struktur Project

```
be/
├── cmd/
│   └── server/
│       └── main.go          # Entry point aplikasi
├── internal/
│   ├── config/              # Konfigurasi aplikasi
│   ├── database/            # Database connection
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Middleware (auth, etc)
│   ├── models/              # Database models
│   └── websocket/           # WebSocket hub dan handlers
├── pkg/                     # Shared packages
├── .env.example             # Contoh environment variables
├── go.mod                   # Go modules
└── README.md
```

## Setup

1. **Install dependencies:**
   ```bash
   go mod download
   ```

2. **Setup environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env sesuai dengan konfigurasi Anda
   ```

3. **Setup database:**
   - Pastikan PostgreSQL sudah terinstall dan running
   - Buat database: `CREATE DATABASE erp_meble;`
   - Update konfigurasi di `.env`

4. **Run migrations (akan ditambahkan nanti):**
   ```bash
   # Migrations akan dibuat untuk auto-migrate models
   ```

5. **Run server:**
   ```bash
   go run cmd/server/main.go
   ```

## Real-time Updates

Aplikasi menggunakan WebSocket untuk real-time updates. Semua client yang terhubung akan menerima update secara real-time ketika ada perubahan data.

### WebSocket Endpoint
- **URL:** `ws://localhost:8080/ws?user_id=<user_id>`
- **Usage:** Client dapat connect ke endpoint ini untuk menerima real-time updates

### Broadcasting Updates
Untuk mengirim update ke semua client yang terhubung:
```go
import "real-erp-mebel/be/internal/websocket"

// Broadcast message
hub.BroadcastMessage([]byte(`{"type": "update", "data": {...}}`))
```

## API Endpoints

### Public Endpoints
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/register` - Register user baru

### Protected Endpoints (require JWT token)
- `GET /api/v1/users/me` - Get current user info

### Health Check
- `GET /health` - Check server status

## Development

### Menambahkan Model Baru
1. Buat file di `internal/models/`
2. Model akan di-migrate otomatis saat aplikasi start (akan ditambahkan)

### Menambahkan Handler Baru
1. Buat handler function di `internal/handlers/`
2. Register route di `cmd/server/main.go`

### Real-time Updates
Gunakan WebSocket hub untuk broadcast updates:
```go
// Di handler atau service Anda
hub.BroadcastMessage([]byte(jsonData))
```

## Next Steps

1. Tambahkan auto-migration untuk models
2. Implementasi CRUD operations untuk entities ERP
3. Tambahkan real-time notifications untuk events penting
4. Implementasi role-based access control (RBAC)
5. Tambahkan logging dan monitoring
6. Setup testing

