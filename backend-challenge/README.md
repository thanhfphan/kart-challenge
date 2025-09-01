# Backend Challenge - Order Food Online API

A high-performance e-commerce API built with Go, featuring advanced coupon processing capabilities, Redis caching, and OpenAPI-first design.

## Prerequisites

Before running this project, ensure you have the following installed:

- **Go 1.24+** - The application is built with Go 1.24
- **Docker & Docker Compose** - For running infrastructure services (MySQL, Redis)
- **Make** - For running build commands
- **oapi-codegen** - For OpenAPI code generation (install with `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest`)

## Getting Started

```bash
# MariaDB database on port 3306
# Redis cache on port 6379
make up

# Run database migration(create tables, seed data)
make migrate

# This critical step processes large coupon data files from S3:
# - Downloads 3 compressed coupon files from AWS S3
# - Performs external sorting and merging operations
# - Generates valid coupon codes and stores them in the database
# - Creates optimized binary files for fast lookups(optional)
make preprocess

# Start the API server
# The API server starts on port 4010 with the following endpoints available:
# - Health check: `http://localhost:4010/health-check`
# - API documentation: `http://localhost:4010/api-docs`
# - OpenAPI spec: `http://localhost:4010/openapi.json`
make start
```
## Other useful Commands

- `make tests` - Run unit tests
- `make down` - Stop infrastructure services
- `make generate` - Generate code from OpenAPI specification

## Core Feature Implementation: Coupon Processing System

The system handles retrieving valid coupons from 3 large data files through a sophisticated multi-stage pipeline:

### Step-by-Step Process Flow

#### 1. **Data Ingestion Phase**
- **Source**: 3 compressed files from S3 (`couponbase1.gz`, `couponbase2.gz`, `couponbase3.gz`)
- **Parallel Processing**: Downloads are processed concurrently using goroutines
- **Format Conversion**: Each gzipped file is converted to binary pairs format for efficient processing

#### 2. **External Sorting Phase**
- **Challenge**: Files are too large to fit in memory (1M+ records each)
- **Solution**: External merge sort with configurable chunk size (1,000,000 records)
- **Process**: 
  - Split large files into memory-manageable chunks
  - Sort each chunk individually
  - Merge sorted chunks back together
- **Output**: 3 sorted binary files (`pairs1.sorted.bin`, `pairs2.sorted.bin`, `pairs3.sorted.bin`)

#### 3. **Three-Way Merge Phase**
- **Algorithm**: Efficient 3-way merge to find common coupon codes across all files
- **Logic**: Only coupons present in at least 2 files are considered valid
- **Output**: Single `valid_coupons.txt` file containing verified coupon codes

#### 4. **Database Storage Phase**
- **Batch Processing**: Coupons are inserted in batches of 1,000 for optimal performance
- **Streaming**: Uses `bufio.Scanner` for memory-efficient file reading
- **Validation**: Each coupon is validated before database insertion

## Architecture & Design

### OpenAPI-First Approach
- **Specification**: Complete OpenAPI 3.1 specification in `openapi.yaml`
- **Code Generation**: Server interfaces and models generated from spec using `oapi-codegen`
- **Documentation**: Auto-generated interactive documentation
- **Validation**: Request/response validation based on OpenAPI schema

### Clean Architecture Principles
- **Separation of Concerns**: Clear boundaries between layers
- **Dependency Injection**: Configurable dependencies through interfaces
- **Repository Pattern**: Data access abstraction
- **Use Case Layer**: Business logic encapsulation

## Codebase Structure

```
backend-challenge/
├── app/                    # Application layer
│   ├── delivery/          # HTTP handlers and OpenAPI implementation
│   ├── dto/               # Data transfer objects
│   ├── models/            # Domain models
│   ├── repos/             # Repository implementations
│   └── usecases/          # Business logic layer
├── cmd/                   # Application entry points
│   ├── api/               # Main API server
│   ├── migrate/           # Database migration tool
│   └── preprocess/        # Coupon preprocessing pipeline
├── pkg/                   # Shared packages
│   ├── cache/             # Redis caching implementation
│   ├── infras/            # Infrastructure utilities
│   ├── logging/           # Structured logging
│   └── validation/        # Request validation
├── config/                # Configuration management
├── env/                   # Environment setup
├── migrations/            # Database schema migrations
└── builders/              # Docker and build configurations
```
