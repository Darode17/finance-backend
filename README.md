# Finance Data Processing and Access Control Backend

A RESTful backend for a finance dashboard system built with **Go (Gin + GORM + SQLite)**.

---

## Tech Stack

| Layer | Choice | Reason |
|---|---|---|
| Language | Go 1.21 | Performant, typed, excellent stdlib |
| Framework | Gin | Lightweight, fast HTTP router with middleware support |
| ORM | GORM | Clean data modeling; easy migrations |
| Database | SQLite | Zero-config; portable for assessment; swap to PostgreSQL trivially |
| Auth | JWT (golang-jwt/jwt) | Stateless token auth, simple to reason about |
| Password hashing | bcrypt | Industry-standard hashing |

---

## Project Structure

```
finance-backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point, router setup
├── internal/
│   ├── database/
│   │   └── db.go            # DB init, auto-migrate, seed admin user
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── record_handler.go
│   │   └── dashboard_handler.go
│   ├── middleware/
│   │   └── auth.go          # JWT validation + role guard middleware
│   ├── models/
│   │   └── models.go        # DB models + request/response DTOs
│   └── services/
│       ├── auth_service.go
│       ├── user_service.go
│       ├── record_service.go
│       └── dashboard_service.go
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## Setup and Running

### Prerequisites
- Go 1.21+
- CGO enabled (required by `go-sqlite3`; comes with standard Go install)

### Steps

```bash
# 1. Clone the repo
git clone <your-repo-url>
cd finance-backend

# 2. Install dependencies
go mod tidy

# 3. (Optional) copy and edit env
cp .env.example .env

# 4. Run the server
make run
# or: go run ./cmd/server
```

The server starts on **http://localhost:8080**.

On first run, a default **admin** user is seeded automatically:
```
Email:    admin@finance.local
Password: admin@123
```

---

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | HTTP port |
| `JWT_SECRET` | `finance-backend-secret-key-...` | JWT signing secret |
| `DB_PATH` | `finance.db` | SQLite file path |
| `GIN_MODE` | `debug` | `debug` or `release` |

---

## Roles and Permissions

| Action | Viewer | Analyst | Admin |
|---|:---:|:---:|:---:|
| Login | ✅ | ✅ | ✅ |
| View own profile (`/me`) | ✅ | ✅ | ✅ |
| View financial records | ✅ | ✅ | ✅ |
| Create financial records | ❌ | ✅ | ✅ |
| Update financial records | ❌ | ❌ | ✅ |
| Delete financial records | ❌ | ❌ | ✅ |
| View dashboard summary | ❌ | ✅ | ✅ |
| Manage users (CRUD) | ❌ | ❌ | ✅ |

Enforced via `middleware.RequireRole(...)` on each route group.

---

## API Reference

All endpoints return a consistent envelope:
```json
{ "success": true, "message": "...", "data": { ... } }
{ "success": false, "error": "..." }
```

### Auth

#### `POST /api/auth/login`
```json
{ "email": "admin@finance.local", "password": "admin@123" }
```
Returns `{ "token": "...", "user": { ... } }` — use the token as `Authorization: Bearer <token>`.

---

### Users (Admin only)

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/users` | Create user |
| GET | `/api/users` | List users (paginated) |
| GET | `/api/users/:id` | Get user by ID |
| PUT | `/api/users/:id` | Update user |
| DELETE | `/api/users/:id` | Soft-delete user |

**Create user body:**
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com",
  "password": "secret123",
  "role": "analyst"
}
```

**Update user body** (all fields optional):
```json
{
  "name": "Jane Smith",
  "role": "viewer",
  "is_active": false
}
```

**List users query params:** `?page=1&page_size=20`

---

### Financial Records

| Method | Endpoint | Allowed Roles |
|---|---|---|
| GET | `/api/records` | All |
| GET | `/api/records/:id` | All |
| POST | `/api/records` | Analyst, Admin |
| PUT | `/api/records/:id` | Admin |
| DELETE | `/api/records/:id` | Admin |

**Create/Update body:**
```json
{
  "amount": 5000.00,
  "type": "income",
  "category": "Salary",
  "date": "2024-03-15",
  "description": "March salary"
}
```

**GET /api/records query params:**
```
?type=income           # Filter by type: income | expense
?category=Salary       # Filter by category
?start_date=2024-01-01 # Filter from date (YYYY-MM-DD)
?end_date=2024-03-31   # Filter to date (YYYY-MM-DD)
?page=1                # Pagination
?page_size=20          # Items per page (max 100)
```

---

### Dashboard (Analyst + Admin)

#### `GET /api/dashboard/summary`
Returns aggregated data:
```json
{
  "total_income": 50000.00,
  "total_expenses": 32000.00,
  "net_balance": 18000.00,
  "total_records": 42,
  "category_totals": [
    { "category": "Salary", "total": 50000.00, "count": 2 }
  ],
  "monthly_trends": [
    { "month": "2024-03", "income": 25000, "expense": 16000, "net": 9000 }
  ],
  "recent_activity": [ ... ]
}
```

---

### My Profile (All authenticated users)

#### `GET /api/me`
Returns the currently authenticated user's profile.

---

## Design Decisions and Assumptions

1. **Soft deletes for records** — Financial records are never hard-deleted (IS_DELETED flag). This preserves audit trail integrity. Users are hard-deleted (GORM soft-delete via `DeletedAt`).

2. **SQLite for simplicity** — Swapping to PostgreSQL requires only changing the GORM driver import and DSN; no business logic changes needed.

3. **Stateless JWT** — Tokens expire in 24 hours. No refresh token mechanism (out of scope); this can be added with a `refresh_tokens` table.

4. **Analyst can create records but not edit/delete** — A deliberate design choice: analysts should be able to enter new data but modifications to existing records require admin oversight.

5. **Dashboard queries run against live DB** — No caching layer for this scope. For production, summaries could be pre-computed or cached with Redis.

6. **Input validation** — Uses Gin's `binding` tags backed by `go-playground/validator`. All DTOs validate types, required fields, enum values, and numeric ranges.

---

## Optional Enhancements Implemented

- ✅ JWT Authentication
- ✅ Pagination on all list endpoints
- ✅ Filtering on records (type, category, date range)
- ✅ Soft deletes on financial records
- ✅ Consistent API response envelope
- ✅ Seed admin user on first run
- ✅ `.env.example` for configuration

## Not Implemented (Out of Scope)

- Unit/integration tests
- Rate limiting
- Search (full-text)
- Swagger/OpenAPI spec generation

---

## Quick cURL Examples

```bash
# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@finance.local","password":"admin@123"}'

# Create a record (replace TOKEN)
curl -X POST http://localhost:8080/api/records \
  -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"amount":1500,"type":"expense","category":"Utilities","date":"2024-04-01"}'

# Get dashboard summary
curl http://localhost:8080/api/dashboard/summary \
  -H "Authorization: Bearer TOKEN"
```
