# edata

![Status](https://img.shields.io/badge/Status-In_Development-yellow)
![Version](https://img.shields.io/badge/Version-0.1.0-blue)
![Stability](https://img.shields.io/badge/Stability-Experimental-orange)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-00ADD8)
![License](https://img.shields.io/badge/License-MIT-green)

A Go application demonstrating subscription service implementation with best practices.

> **Note**: This project is currently under active development. Features and documentation may be incomplete or subject to change.

## 🎯 Overview

This educational project showcases how to build applications in Go, featuring:

- Clean architecture and project organization
- PostgreSQL database management
- Docker containerization
- Database migrations
- JWT authentication
- RESTful API design

## 🚧 Development Status

- **Stage**: Alpha
- **API Stability**: Experimental
- **Features**: In Progress

### Current Focus Areas
- [ ] Core API implementation
- [ ] Authentication & authorization
- [ ] Database schema design
- [ ] Testing infrastructure
- [ ] API documentation

## 🏗️ Project Structure

```
edata/
├── bin/                  # Compiled binaries
├── cmd/                  # Application entrypoints
│   ├── api/             # API server
│   └── migrate/         # DB migration tool
├── config/              # Configuration management
├── db/                  # Database layer
├── infra/               # Infrastructure setup
├── service/             # Business logic
├── types/               # Core types/models
└── utils/               # Shared utilities
```

[View detailed structure](#detailed-structure)

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or higher
- Docker and Docker Compose
- Make

### Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/sudarakas/edata.git
   cd edata
   ```

2. **Set up environment**
   ```bash
   cp .env.example .env
   ```

3. **Start database**
   ```bash
   docker-compose -f infra/db.yaml up -d
   ```

4. **Run migrations**
   ```bash
   make migrate-up
   ```

5. **Start server**
   ```bash
   make run
   ```

## ⚙️ Configuration

### Environment Variables

```env
# Server
PUBLIC_HOST=http://localhost
PORT=8080

# Database
DB_USER=postgres
DB_PASSWORD=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=edata_db
```

## 📋 Available Commands

| Command | Description |
|---------|------------|
| `make run` | Start the application |
| `make migrate-up` | Run database migrations |
| `make migrate-down` | Rollback migrations |
| `make test` | Run tests |
| `make lint` | Run linters |

## 🔍 Detailed Structure

```
edata/
├── .env                  # Environment variables
├── bin/                  # Compiled binaries
│   └── edata
├── cmd/                  # Application entrypoints
│   ├── api/             # API server
│   │   └── main.go
│   └── migrate/         # Migration tool
│       ├── main.go
│       └── migrations/
├── config/              # Configuration
├── db/                  # Database layer
├── go.mod              # Go modules
├── go.sum              # Module checksums
├── infra/              # Infrastructure
│   └── db.yaml         # DB container config
├── Makefile            # Build automation
├── service/            # Business logic
│   ├── auth/          # Authentication
│   └── user/          # User management
├── types/              # Core types
│   └── types.go
└── utils/             # Shared utilities
```

## 🤝 Contributing

Currently, this project is in initial development and not accepting contributions. Once stable, contribution guidelines will be provided.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---
Built with ❤️ using Go