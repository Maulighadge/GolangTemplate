# 🚀 Golang Production-Ready MVC Backend Template

A clean, modular, and scalable **MVC (Model-View-Controller) Architecture Template** for Golang backend services, built with **Gin** and **GORM**.

This template is designed to be easily cloned and reused as the foundation for any production-grade RESTful API or microservice.

---

## 🌟 Key Features

- **Architecture**: Clean MVC with decoupled layers (Models, Views/DTOs, Repositories, Services, Controllers).
- **Web Framework**: [Gin Web Framework](https://github.com/gin-gonic/gin) (ultra-fast, lightweight HTTP router).
- **ORM**: [GORM v2](https://gorm.io/) with multi-database support:
  - 🗄️ **SQLite** (Default, zero setup for instant local development)
  - 🐘 **PostgreSQL** (Production-ready with connection pooling)
  - 🐬 **MySQL** (Full support)
- **Database Auto-Migration**: Built-in automatic schema migrations on startup.
- **Authentication & Security**:
  - Secure password hashing using **Bcrypt**.
  - Stateless authentication using **JWT (JSON Web Tokens)** with role-based access control (RBAC).
- **Graceful Shutdown**: Intercepts `SIGINT`/`SIGTERM` to allow active requests and database pools to finish cleanly.
- **Unified API Response Standard**: Uniform JSON envelope structure across success, error, and paginated responses.
- **Cross-Origin Resource Sharing (CORS)**: Configurable CORS middleware.
- **Request Logging & Panic Recovery**: Colored request telemetry and centralized recovery middleware.
- **Pagination & Filtering**: Built-in reusable pagination, search, and category query filters.
- **Containerization**: Multi-stage lightweight `Dockerfile` and `docker-compose.yml` with PostgreSQL.
- **Developer Experience**: Complete `Makefile` with commands for building, running, tidying, and testing.

---

## 📁 Project Architecture & Directory Structure

```text
├── cmd/
│   └── server/
│       └── main.go               # Application entry point & graceful shutdown
├── config/
│   └── config.go                 # Environment configuration loader (.env)
├── internal/
│   ├── controllers/              # [C] HTTP request handlers & input validation
│   │   ├── auth_controller.go
│   │   ├── health_controller.go
│   │   ├── product_controller.go
│   │   └── user_controller.go
│   ├── dto/                      # [V] Data Transfer Objects & API View Models
│   │   ├── auth_dto.go
│   │   ├── product_dto.go
│   │   └── user_dto.go
│   ├── middlewares/              # HTTP Middlewares (JWT Auth, RBAC, CORS, Logger, Recovery)
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   ├── logger_middleware.go
│   │   └── recovery_middleware.go
│   ├── models/                   # [M] Database entities & GORM schemas
│   │   ├── base.go               # Base UUID model with timestamps and soft delete
│   │   ├── product.go
│   │   └── user.go
│   ├── repositories/             # Data access layer (GORM database queries)
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   ├── routes/                   # Routing configuration & dependency injection
│   │   └── routes.go
│   └── services/                 # Core business logic layer
│       ├── auth_service.go
│       ├── product_service.go
│       └── user_service.go
├── pkg/                          # Reusable shared packages across projects
│   ├── database/                 # GORM connection manager & AutoMigrate
│   │   ├── database.go
│   │   └── migration.go
│   ├── logger/                   # Structured logger
│   │   └── logger.go
│   ├── response/                 # Standard JSON API responses
│   │   └── response.go
│   └── utils/                    # Utility helpers
│       ├── jwt.go                # JWT token signing and validation
│       ├── pagination.go         # Pagination calculator and query extractor
│       └── password.go           # Bcrypt password hashing
├── .env.example                  # Sample environment variables
├── .gitignore                    # Git ignore file
├── docker-compose.yml            # Docker Compose setup (App + PostgreSQL)
├── Dockerfile                    # Production multi-stage Docker build
├── Makefile                      # Make CLI commands
├── go.mod                        # Go module dependencies
└── README.md                     # Documentation
```

---

## ⚡ Quick Start

### 1. Prerequisites
- **Go**: 1.22+ (or newer)
- **Make** (optional, for running convenience commands)
- **Docker & Docker Compose** (optional, for containerized deployment)

### 2. Clone and Setup Environment

```bash
# Clone the repository
git clone https://github.com/Maulighadge/GolangTemplate.git
cd GolangTemplate

# Copy sample environment configuration
cp .env.example .env
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Run Locally (Zero Setup with SQLite)
By default, `DB_DRIVER` in `.env` is set to `sqlite`, which will automatically create `app.db` locally without needing to install or run a separate database server.

```bash
go run cmd/server/main.go
# OR using make
make run
```

The server will start at `http://localhost:8080`.

---

## ⚙️ Environment Configuration

Edit the `.env` file to configure your application:

| Variable | Default | Description |
| :--- | :--- | :--- |
| `APP_NAME` | `Golang-MVC-API` | Application name |
| `APP_ENV` | `development` | Environment (`development` or `production`) |
| `APP_PORT` | `8080` | Port on which the HTTP server listens |
| `DB_DRIVER` | `sqlite` | Database driver (`sqlite`, `postgres`, `mysql`) |
| `DB_HOST` | `localhost` | Database host (Postgres/MySQL) |
| `DB_PORT` | `5432` | Database port (Postgres/MySQL) |
| `DB_USER` | `postgres` | Database username |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `golang_template_db`| Database name |
| `DB_SQLITE_PATH` | `./app.db` | Path to SQLite database file |
| `JWT_SECRET` | `secret` | Secret key used for signing JWT tokens |
| `JWT_EXPIRATION_HOURS` | `72` | Expiration time for JWT tokens |
| `CORS_ALLOWED_ORIGINS` | `*` | Allowed CORS origins |

---

## 📡 API Endpoints

### 🩺 Health Check
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/health` | Health & Database status | Public |

### 🔐 Authentication
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/register` | Register a new user | Public |
| `POST` | `/api/v1/auth/login` | Login and get JWT token | Public |

### 👤 User Profile
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/users/me` | Get current user profile | Bearer Token |
| `PUT` | `/api/v1/users/me` | Update current user profile | Bearer Token |
| `PUT` | `/api/v1/users/change-password`| Change user password | Bearer Token |

### 📦 Products (CRUD Example Domain)
| Method | Endpoint | Description | Auth |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/products` | List products (with pagination & search) | Public |
| `GET` | `/api/v1/products/:id` | Get single product by UUID | Public |
| `POST` | `/api/v1/products` | Create a new product | Bearer Token |
| `PUT` | `/api/v1/products/:id` | Update product (Owner or Admin) | Bearer Token |
| `DELETE` | `/api/v1/products/:id` | Delete product (Owner or Admin) | Bearer Token |

---

## 📋 API Request & Response Examples

### 1. Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "email": "jane@example.com",
    "password": "Password123!"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsIn...",
    "user": {
      "id": "e9b5f3a0-47b2-4d7a-8d19-58d7c07b6699",
      "name": "Jane Doe",
      "email": "jane@example.com",
      "role": "user",
      "is_active": true,
      "created_at": "2026-10-21T18:00:00Z",
      "updated_at": "2026-10-21T18:00:00Z"
    }
  }
}
```

### 2. Create Product (Authenticated)
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <YOUR_JWT_TOKEN>" \
  -d '{
    "name": "Mechanical Keyboard",
    "description": "RGB Backlit Custom Mechanical Keyboard",
    "sku": "KB-RGB-001",
    "price": 89.99,
    "stock": 50,
    "category": "Electronics"
  }'
```

### 3. List Products (with Search & Pagination)
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10&search=Keyboard&category=Electronics"
```

**Response with Pagination Metadata:**
```json
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": [
    {
      "id": "7b8956be-57ff-4340-9a2c-9828e833f47e",
      "name": "Mechanical Keyboard",
      "description": "RGB Backlit Custom Mechanical Keyboard",
      "sku": "KB-RGB-001",
      "price": 89.99,
      "stock": 50,
      "category": "Electronics",
      "user_id": "e9b5f3a0-47b2-4d7a-8d19-58d7c07b6699",
      "created_at": "2026-10-21T18:05:00Z",
      "updated_at": "2026-10-21T18:05:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 10,
    "total_rows": 1,
    "total_pages": 1
  }
}
```

---

## 🛠️ How to Reuse This Template for a New Project

To create a new microservice or API using this template:

1. **Clone or use as GitHub template:**
   ```bash
   git clone https://github.com/Maulighadge/GolangTemplate.git my-new-project
   cd my-new-project
   rm -rf .git
   git init
   ```

2. **Rename the Go Module:**
   Find and replace `github.com/Maulighadge/GolangTemplate` with your new module path (e.g. `github.com/your-org/my-new-project`):
   ```bash
   # In go.mod and across all imports
   go mod edit -module github.com/your-org/my-new-project
   find . -type f -name "*.go" -exec sed -i 's|github.com/Maulighadge/GolangTemplate|github.com/your-org/my-new-project|g' {} +
   go mod tidy
   ```

3. **Add your new entity (e.g., `Order`):**
   - Create Model: `internal/models/order.go`
   - Register in AutoMigrate: `pkg/database/migration.go`
   - Create DTO: `internal/dto/order_dto.go`
   - Create Repository: `internal/repositories/order_repository.go`
   - Create Service: `internal/services/order_service.go`
   - Create Controller: `internal/controllers/order_controller.go`
   - Register routes: `internal/routes/routes.go`

---

## 🐳 Docker & Docker Compose

Run the entire application along with PostgreSQL using Docker Compose:

```bash
# Build and run containers in background
docker compose up -d --build

# View logs
docker compose logs -f app

# Stop containers
docker compose down
```

---

## 🧪 Makefile Commands

```bash
make run          # Run the application locally
make build        # Build executable binary in bin/server
make tidy         # Download & cleanup go module dependencies
make test         # Run unit tests
make docker-up    # Spin up Docker containers
make docker-down  # Stop and remove Docker containers
make clean        # Clean up build artifacts and sqlite db
```

---

## 📄 License
This project is open-source and available under the [MIT License](LICENSE).
