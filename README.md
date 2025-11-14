# Core Idea

**The Problem I Solved:**

Every time you build a new application, you rebuild authentication from scratch. Integrating auth across different tech stacks (Node.js, Python, Go, Java, etc.) is painful, time-consuming, and error-prone.

**My Solution:**

A **standalone, language-agnostic authentication microservice** that:

1. **Centralized Auth Service** — One Go microservice handles all authentication

   - User registration, login, password management
   - JWT token generation and validation
   - User management (admin operations)

2. **Simple RBAC** — Roles defined in `.env` file

   - Define roles once in environment: `ROLES=["Super Admin", "Manager", "User"]`
   - Quick to setup, no database configuration needed
   - Easy to modify roles for different deployments
   - Each role has specific permissions based on endpoint access

3. **Shared Secret Key** — Single JWT_SECRET shared across all backends

   - Auth Service generates JWT tokens (includes user role)
   - Any backend (Node.js, Python, Go, Java, etc.) validates tokens locally using the same secret
   - No callback to auth service needed for token validation
   - Backend extracts role from JWT and implements own authorization

4. **Plug & Play Integration** — Simply drop it into any tech stack

   - Copy the JWT validation logic to your backend
   - Validate `Authorization: Bearer <token>` header
   - Extract role from JWT claims
   - Check if user's role can access the endpoint

5. **One Deploy, Multiple Apps** — Deploy auth service once, use everywhere
   - All your apps share the same authentication
   - Consistent user identity across services
   - Consistent roles across all backends
   - No redundant auth rebuilding

**The Flow:**

```
┌─────────────────────────────────────┐
│   Auth Service (Go)                 │
│                                     │
│ Roles (from .env):                  │
│ - Super Admin                       │
│ - Manager                           │
│ - User                              │
│                                     │
│ Features:                           │
│ - Register, Login, Tokens           │
│ - User Management (CRUD)            │
│ - Role Assignment                   │
│ - Issues JWT with role              │
└────────────┬────────────────────────┘
             │ (shares JWT_SECRET + role info)
   ┌─────────┼──────────┬──────────┐
   │         │          │          │
┌──▼───┐  ┌──▼───┐  ┌───▼───┐  ┌───▼───┐
│Node  │  │Pyth  │  │ Go    │  │Java   │
│App   │  │ App  │  │ App   │  │ App   │
│      │  │      │  │       │  │       │
│Vali- │  │Vali- │  │Vali-  │  │Vali-  │
│date  │  │date  │  │dates  │  │dates  │
│JWT   │  │JWT   │  │JWT    │  │JWT    │
│+ Role│  │+ Role│  │+ Role │  │+ Role │
└──────┘  └──────┘  └───────┘  └───────┘
```

**Key Features:**

- ✅ **User Management** — Register, login, change password
- ✅ **Admin Operations** — Create, read, update, delete users with role assignment
- ✅ **Simple RBAC** — Roles defined in `.env` for quick setup
- ✅ **Reusable** — No more rebuilding auth
- ✅ **Scalable** — One auth service, infinite apps
- ✅ **Lightweight** — Minimal dependencies, high performance
- ✅ **Open Source** — Community contribution and trust
- ✅ **Framework Agnostic** — Works with any backend tech stack
- ✅ **Secure** — JWT validation, bcrypt hashing, role-based access

**Configuration (`.env`):**

```env
# Roles - Define once, use everywhere
ROLES=["Super Admin", "Manager", "User"]
DEFAULT_REGISTRATION_ROLE=User

# Super Admin Setup
SUPER_ADMIN_EMAIL=superadmin@web.com
SUPER_ADMIN_PASSWORD=superadminpass123

# JWT Configuration
JWT_SECRET=your-shared-secret-key
```

**What Makes It Different:**

1. **Not a framework** — It's a standalone service
2. **Not tied to database** — Works with your existing tech stack
3. **Not complex** — Simple role configuration in `.env`
4. **Not reinventing** — Uses proven JWT standards
5. **Not restrictive** — Each backend implements its own authorization logic

## Auth Service - Microservice

A lightweight, standalone authentication microservice built with Go. Provides user registration, login, token refresh, password management, role-based access control, and admin user management using JWT tokens and UUID identifiers.

## Features

- **User Management**: Register, login, password change, profile management
- **RBAC (Role-Based Access Control)**: Super Admin, Manager, User roles with dynamic permission system
- **Admin User Management**: Create, read, update, delete users with role assignment
- **JWT Authentication**: Access tokens (15 min) and refresh tokens (7 days)
- **Secure Password Hashing**: bcrypt for password security
- **Soft Deletes**: Users marked as deleted, not permanently removed
- **UUID Identifiers**: Scalable, globally unique user IDs
- **PostgreSQL Integration**: Raw SQL queries for performance
- **Docker Ready**: Pre-built binary support with Docker
- **Simple & Minimal**: No ORM, pure standard library + essential libraries
- **In-House Auth**: Deploy once, use across multiple microservices

## Tech Stack

- **Language**: Go 1.24.4
- **Database**: PostgreSQL
- **Authentication**: JWT (HS256)
- **Password Hashing**: bcrypt
- **ID Generation**: UUID v4
- **Dependencies**:
  - `github.com/golang-jwt/jwt/v5` - JWT token generation/verification
  - `github.com/lib/pq` - PostgreSQL driver
  - `golang.org/x/crypto/bcrypt` - Password hashing
  - `github.com/joho/godotenv` - Environment variable management

## Project Structure

```
.
├── Dockerfile
├── README.md
├── cmd
│   └── auth
│       └── main.go
├── config
│   └── config.go
├── db
│   ├── Queries
│   │   ├── Common.go
│   │   ├── CreateRole.go
│   │   ├── CreateUser.go
│   │   ├── DeleteUser.go (soft delete)
│   │   ├── GetAllUsers.go
│   │   ├── GetRoleByName.go
│   │   ├── GetUserByEmail.go
│   │   ├── GetUserByID.go
│   │   ├── GetUserByIDAdmin.go
│   │   ├── UpdatePassword.go
│   │   ├── UpdateUser.go
│   │   └── UpdateUserRole.go
│   ├── Seeder
│   │   ├── SeedRoles.go
│   │   └── SeedSuperAdmin.go
│   └── database
│       ├── Close.go
│       ├── CreateTables.go
│       └── InitDB.go
├── go.mod
├── go.sum
├── handlers
│   ├── admin
│   │   ├── CreateUserHandler.go
│   │   ├── DeleteUserHandler.go
│   │   ├── GetAllUsersHandler.go
│   │   ├── GetUserHandler.go
│   │   ├── UpdateUserHandler.go
│   │   └── UpdateUserRoleHandler.go
│   ├── auth
│   │   ├── ChangePasswordHandler.go
│   │   ├── LoginHandler.go
│   │   ├── RefreshTokenHandler.go
│   │   └── RegisterHandler.go
│   ├── common.go
│   ├── handler.go
│   └── user
│       └── GetProfileHandler.go
├── middleware
│   ├── auth
│   │   ├── AuthMiddleware.go
│   │   └── GetClaimsFromContext.go
│   ├── constants
│   │   ├── Constants.go
│   │   └── RespondError.go
│   └── role.go
├── models
│   ├── admin.go
│   ├── token.go
│   └── user.go
├── test
│   └── test.py
└── utils
    ├── jwt
    │   ├── Common.go
    │   ├── GenerateToken.go
    │   └── VerifyToken.go
    └── password
        ├── Common.go
        ├── HashPassword.go
        └── VerifyPassword.go
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

   # RBAC Configuration
   ROLES=["Super Admin", "Manager", "User"]
   DEFAULT_REGISTRATION_ROLE=User
   SUPER_ADMIN_EMAIL=superadmin@web.com
   SUPER_ADMIN_PASSWORD=superadminpass123
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
     -e ROLES='["Super Admin", "Manager", "User"]' \
     -e DEFAULT_REGISTRATION_ROLE=User \
     -e SUPER_ADMIN_EMAIL=superadmin@web.com \
     -e SUPER_ADMIN_PASSWORD=superadminpass123 \
     auth-service:latest
   ```

## API Endpoints

### How It Works

The auth service provides a complete authentication and user management system:

1. **Registration & Login**: New users register with default role, existing users login
2. **Token Generation**: Access tokens (15 min) and refresh tokens (7 days) are issued
3. **Protected Endpoints**: All endpoints except `/health`, `/register`, `/login`, `/refresh` require valid access token
4. **Role-Based Access**: Admin endpoints check for "Super Admin" role
5. **UUID Identifiers**: All users identified by UUID for scalability
6. **Soft Deletes**: Deleted users retain history (deleted_at timestamp)

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

- User gets default role from `DEFAULT_REGISTRATION_ROLE`
- Returns UUID as user ID

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

- Returns tokens and user info with role

#### Refresh Token

```bash
POST /refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

Response: `{access_token}`

- Generates new access token using refresh token

### Protected Endpoints (Require Access Token)

#### Get Profile

```bash
GET /profile
Authorization: Bearer your-access-token
```

Response: `{id, username, email, role, deleted_at, created_at, updated_at}`

- Returns authenticated user's profile
- Shows deletion status if soft deleted

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

- Requires old password verification

### Admin Endpoints (Super Admin Only)

#### Get All Users

```bash
GET /admin/users
Authorization: Bearer super-admin-access-token
```

Response: `{total: number, users: [...]}`

- Lists all non-deleted users
- Excludes soft-deleted users

#### Get Specific User

```bash
GET /admin/users/get/{uuid}
Authorization: Bearer super-admin-access-token
```

Response: `{id, username, email, role, deleted_at, created_at, updated_at}`

- Shows user details including deletion status
- Can view soft-deleted users

#### Create User (Admin)

```bash
POST /admin/users/create
Authorization: Bearer super-admin-access-token
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "password123",
  "role": "Manager"
}
```

Response: `{id, username, email, role, created_at, updated_at}`

- Admin creates user with specific role
- Does not send welcome email (you add this later)

#### Update User

```bash
PATCH /admin/users/update/{uuid}
Authorization: Bearer super-admin-access-token
Content-Type: application/json

{
  "username": "updatedname",
  "email": "updated@example.com"
}
```

Response: `{id, username, email, role, updated_at, ...}`

- Update username and/or email
- At least one field required

#### Delete User (Soft Delete)

```bash
DELETE /admin/users/delete/{uuid}
Authorization: Bearer super-admin-access-token
```

Response: `{"message":"user deleted successfully"}`

- Marks user as deleted (sets deleted_at timestamp)
- User cannot login after deletion
- Prevents self-deletion

#### Update User Role

```bash
PUT /admin/users/role/{uuid}
Authorization: Bearer super-admin-access-token
Content-Type: application/json

{
  "role": "Manager"
}
```

Response: `{message: "user role updated successfully", user: {...}}`

- Changes user's role
- Prevents self-role-change
- Prevents removing last Super Admin

## Token Details

- **Access Token Duration**: 15 minutes
- **Refresh Token Duration**: 7 days
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Header Format**: `Authorization: Bearer <token>`
- **Token Type in JWT**: Includes "access" or "refresh" to distinguish token types

## Database Schema

### Users Table

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(100) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(100) DEFAULT 'User',
  deleted_at TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Roles Table

```sql
CREATE TABLE roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(100) UNIQUE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

Tables are automatically created on first run.

## Usage Examples

### Complete Flow

```bash
# 1. Register new user
REGISTER=$(curl -s -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john",
    "email": "john@example.com",
    "password": "secure123"
  }')

ACCESS_TOKEN=$(echo $REGISTER | jq -r '.access_token')
REFRESH_TOKEN=$(echo $REGISTER | jq -r '.refresh_token')

# 2. Get profile
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

# 4. Refresh token
NEW_TOKEN=$(curl -s -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}" | jq -r '.access_token')

# 5. Login with new password
LOGIN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "newsecure456"
  }')
```

### Admin Operations

```bash
# Login as super admin first
ADMIN_LOGIN=$(curl -s -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@web.com",
    "password": "superadminpass123"
  }')

ADMIN_TOKEN=$(echo $ADMIN_LOGIN | jq -r '.access_token')

# Create new manager
curl -X POST http://localhost:8080/admin/users/create \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "manager1",
    "email": "manager@example.com",
    "password": "manager123",
    "role": "Manager"
  }'

# Get all users
curl -X GET http://localhost:8080/admin/users \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Update user role
curl -X PUT http://localhost:8080/admin/users/role/{user-uuid} \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"role": "Manager"}'
```

## Security Considerations

- **Never** commit `.env` file with real secrets
- Change `JWT_SECRET` in production to a strong, random value (minimum 32 characters)
- Change default `SUPER_ADMIN_EMAIL` and `SUPER_ADMIN_PASSWORD` immediately
- Use HTTPS in production
- Store tokens securely on client side (HttpOnly cookies recommended)
- Implement token blacklisting for logout (future enhancement)
- Use environment-specific configuration for different deployments
- Regularly rotate `JWT_SECRET`
- Monitor failed login attempts (future enhancement)

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
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found (resource doesn't exist)
- `409` - Conflict (user already exists, duplicate email)
- `500` - Internal Server Error

## Testing

A comprehensive test suite is included in `test/test.py`:

```bash
pip install requests python-dotenv
python3 test/test.py
```

Tests cover:

- User registration and role assignment
- Login and token generation
- Token refresh
- Profile management
- Password changes
- Admin operations (create, read, update, delete, role change)
- Authorization and access control
- Unauthorized access attempts

## Future Enhancements

- Token blacklisting for logout
- Email verification on registration
- Two-factor authentication
- Password reset via email
- Activity logging and audit trail
- Rate limiting on authentication endpoints
- OAuth2 integration
- API key authentication for service-to-service communication
- Permission management system (dynamic groups and permissions)
- GraphQL API alternative

## Deployment Checklist

- [ ] Change `JWT_SECRET` to strong random value
- [ ] Change `SUPER_ADMIN_EMAIL` and `SUPER_ADMIN_PASSWORD`
- [ ] Use HTTPS in production
- [ ] Set `DB_SOURCE` to production database
- [ ] Enable proper logging
- [ ] Set up monitoring and alerting
- [ ] Backup database regularly
- [ ] Review security considerations section

## License

MIT
