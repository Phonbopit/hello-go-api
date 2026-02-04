# Hello Go API

A learning project to understand Go

## Prerequisites

- Go 1.22+

Check your version:

```bash
go version
```

## Build & Run

### Option 1: Run directly (like `node index.js`)

```bash
go run .

# Build executable
go build -o hello-go-api

# Run the executable
./hello-go-api

# Build optimized binary (smaller size, no debug info)
go build -ldflags="-s -w" -o hello-go-api

# Build for different OS (cross-compile)
GOOS=linux GOARCH=amd64 go build -o hello-go-api-linux
```

## API Endpoints

| Method   | Path                | Auth Required | Description        |
| -------- | ------------------- | ------------- | ------------------ |
| `POST`   | `/admin/keys`       | No            | Create API key     |
| `GET`    | `/admin/keys`       | No            | List API keys      |
| `GET`    | `/v1/products`      | No            | List all products  |
| `GET`    | `/v1/products/{id}` | No            | Get single product |
| `POST`   | `/v1/products`      | **Yes**       | Create product     |
| `DELETE` | `/v1/products/{id}` | **Yes**       | Delete product     |

## Quick Start

```bash
# 1. Create API key
curl -X POST http://localhost:8080/admin/keys -d '{"name":"My App"}'
# Copy the "key" from response (shown once only)

# 2. Create product (requires key)
curl -X POST http://localhost:8080/v1/products \
  -H "X-API-Key: YOUR_KEY_HERE" \
  -d '{"id":"1","name":"MacBook","price":2499}'

# 3. List products (public)
curl http://localhost:8080/v1/products
```

## Test Script

```bash
./test-api.sh
```
