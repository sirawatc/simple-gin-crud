# Simple Gin CRUD API

<div align="center">
  <img src="https://img.shields.io/badge/Go-1.24.5-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/Gin-1.10.1-green.svg" alt="Gin Version">
  <img src="https://img.shields.io/badge/GORM-1.30.1-orange.svg" alt="GORM Version">
  <img src="https://img.shields.io/badge/coverage-90%25-brightgreen.svg" alt="PostgreSQL">
  <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="License">
</div>

A simple boilerplate CRUD API built with Go Gin framework following clean architecture principles.

## 🎯 Features

- **RESTful API**: Complete CRUD operations example.
- **Database**: PostgreSQL with GORM and auto-migrations.
- **Testing**: Unit tests with 90%+ coverage.
- **Validation**: Input validation with detailed error messages.
- **Middleware**: CORS support and request id injection.
- **Logging**: Structured logging with request id tracking.
- **Error Handling**: Standardized error responses with custom codes.
- **Documentation**: Complete Postman collection for API testing.
- **Docker Support**: Containerized application with Docker Compose.
- **CI Pipeline**: Automated code analysis and testing.

## 📋 Prerequisites

- 🐹 **Go 1.24.5** or higher  
- 🐳 **Docker Desktop** (for containerized setup)
- 🐘 **PostgreSQL** (if running locally)
- 🛠️ **Make** (optional, for convenience commands)

## 🚀 Quick Start

Choose your preferred setup method. Begin by cloning the repository:

```bash
git clone https://github.com/sirawatc/simple-gin-crud.git
cd simple-gin-crud
```

### Option 1: Using Docker

Configure environment variables in `docker-compose.yml` if needed.

```bash
docker-compose up -d
```

### Option 2: Local Development

Ensure your database is running and create a `.env` file with your configuration.

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Run the server**
   ```bash
   make dev
   # or
   go run cmd/main/main.go
   ```

The API will be available at `http://localhost:8080`

## 🧪 Testing

### Unit Test

Run the test using:

```bash
make test
# or
go test ./...
```

### Postman Collection

This project includes a Postman collection for testing.

For more detail, see [API Documentation](doc/README.md).

## 🏗️ Project Structure

```
simple-gin-crud/
├── cmd/main/           # Server entry point
├── database/           # Database initialization and migrations
├── internal/           # Internal application code
│   ├── .../            # Domain folders
│   └── shared/         # Shared components
├── pkg/                # Reusable packages
├── server/             # Server configuration and routing
├── doc/                # Documentation and Postman collection
├── docker-compose.yml  # Docker services configuration
└── Dockerfile          # Application container definition
```

## 🔄 CI Pipeline

The CI pipeline includes:

- **🔍 Code Quality**: Linting with golangci-lint
- **🧪 Testing**: Unit tests with coverage reporting
- **🐳 Docker Build**: Container image building and testing
- **📊 Coverage**: Test coverage analysis and reporting
- **🔒 Security**: Dependency vulnerability scanning

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE.md](LICENSE.md) file for details.