# Auth Service - Microservice

A lightweight, standalone authentication microservice built with Go. Provides user registration, login, token refresh, and password management using JWT tokens.

## Features

- **User Management**: Register, login, and password change functionality
- **JWT Authentication**: Access tokens (15 min) and refresh tokens (7 days)
- **Secure Password Hashing**: bcrypt for password security
- **PostgreSQL Integration**: Raw SQL queries for performance
- **Docker Ready**: Pre-built binary support with Docker
- **Simple & Minimal**: No ORM, pure standard library + essential libraries
- **In-House Auth**: Deploy once, use across multiple microservices

## Tech Stack

- **Language**: Go 1.24.4
- **Database**: PostgreSQL
- **Authentication**: JWT (HS256)
- **Password Hashing**: bcrypt
- **Dependencies**:
  - `github.com/golang-jwt/jwt/v5` - JWT token generation/verification
  - `github.com/lib/pq` - PostgreSQL driver
  - `golang.org/x/crypto/bcrypt` - Password hashing
  - `github.com/joho/godotenv` - Environment variable management

## Project Structure

```
.
├── cmd/auth/
│   └── main.go              # Entry point
├── config/
│   └── config.go            # Configuration management
├── db/
│   ├── db.go                # Database connection
│   └── queries.go           # SQL queries
├── handlers/
│   ├── auth.go              # Auth endpoints
│   └── user.go              # User endpoints
├── middleware/
│   └── auth.go              # JWT middleware
├── models/
│   ├── user.go              # User data structures
│   └── token.go             # Token data structures
├── utils/
│   ├── jwt.go               # JWT utilities
│   └── password.go          # Password utilities
├── go.mod
├── go.sum
├── Dockerfile
└── .env
```

## Setup

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- Docker (optional)

### Local Development

1. **Clone and install dependencies**:

   ```bash
   go mod download
   ```

2. **Create `.env` file**:

   ```env
   DB_DRIVER=postgres
   DB_SOURCE=postgresql://user:password@localhost:5432/auth_db?sslmode=disable
   JWT_SECRET=your-super-secret-key-change-this-in-production
   SERVER_PORT=8080
   ```

3. **Build the binary** (choose based on your OS):

   **macOS**:

   ```bash
   GOOS=darwin GOARCH=amd64 go build -o auth-service cmd/auth/main.go
   ```

   **Linux**:

   ```bash
   GOOS=linux GOARCH=amd64 go build -o auth-service cmd/auth/main.go
   ```

   **Windows**:

   ```bash
   GOOS=windows GOARCH=amd64 go build -o auth-service.exe cmd/auth/main.go
   ```

   **Or cross-compile for all platforms**:

   ```bash
   # macOS
   GOOS=darwin GOARCH=amd64 go build -o auth-service-darwin cmd/auth/main.go

   # Linux
   GOOS=linux GOARCH=amd64 go build -o auth-service-linux cmd/auth/main.go

   # Windows
   GOOS=windows GOARCH=amd64 go build -o auth-service-windows.exe cmd/auth/main.go
   ```

4. **Run the service**:
   ```bash
   ./auth-service
   ```

The service will start on `http://localhost:8080`

### Docker Deployment

1. **Build binary for Linux**:

   ```bash
   GOOS=linux GOARCH=amd64 go build -o auth-service cmd/auth/main.go
   ```

2. **Build Docker image**:

   ```bash
   docker build -t auth-service:latest .
   ```

3. **Run container**:
   ```bash
   docker run -p 8080:8080 \
     -e DB_DRIVER=postgres \
     -e DB_SOURCE="postgresql://user:password@host:5432/auth_db?sslmode=require" \
     -e JWT_SECRET="your-secret-key" \
     -e SERVER_PORT=8080 \
     auth-service:latest
   ```

## API Endpoints

### Public Endpoints

#### Health Check

```bash
GET /health
```

Response: `{"status":"healthy"}`

#### Register User

```bash
POST /register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "password123"
}
```

Response: `{access_token, refresh_token, user}`

#### Login

```bash
POST /login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

Response: `{access_token, refresh_token, user}`

#### Refresh Token

```bash
POST /refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

Response: `{access_token}`

### Protected Endpoints (Require Access Token)

#### Get Profile

```bash
GET /profile
Authorization: Bearer your-access-token
```

Response: `{id, username, email, created_at, updated_at}`

#### Change Password

```bash
POST /change-password
Authorization: Bearer your-access-token
Content-Type: application/json

{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

Response: `{"message":"password changed successfully"}`

## Token Details

- **Access Token Duration**: 15 minutes
- **Refresh Token Duration**: 7 days
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Header Format**: `Authorization: Bearer <token>`

## Usage Example

```bash
# 1. Register
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "secure123"
  }')

ACCESS_TOKEN=$(echo $REGISTER_RESPONSE | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

# 2. Get profile using access token
curl -X GET http://localhost:8080/profile \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Change password
curl -X POST http://localhost:8080/change-password \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "secure123",
    "new_password": "newsecure456"
  }'
```

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

The table is automatically created on first run.

## Security Considerations

- **Never** commit `.env` file with real secrets
- Change `JWT_SECRET` in production to a strong, random value
- Use HTTPS in production
- Store tokens securely on the client side (HttpOnly cookies recommended)
- Implement token blacklisting for logout functionality (future enhancement)
- Use environment-specific configuration for different deployments

## Error Responses

All errors follow a consistent format:

```json
{
  "error": "error message"
}
```

Common status codes:

- `400` - Bad Request (validation error)
- `401` - Unauthorized (invalid/expired token)
- `409` - Conflict (user already exists)
- `500` - Internal Server Error

## Future Enhancements

- Token blacklisting for logout
- Email verification on registration
- Two-factor authentication
- Role-based access control (RBAC)
- Audit logging
- Rate limiting
- OAuth2 integration

## License

MIT
