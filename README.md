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

| Method   | Path                | Description        |
| -------- | ------------------- | ------------------ |
| `GET`    | `/v1/products`      | List all products  |
| `GET`    | `/v1/products/{id}` | Get single product |
| `POST`   | `/v1/products`      | Create product     |
| `DELETE` | `/v1/products/{id}` | Delete product     |

## Testing with curl

```bash
## Create product
curl -X POST http://localhost:8080/v1/products \
  -d '{"id":"1","name":"MacBook Pro","price":2499.99}'

curl -X POST http://localhost:8080/v1/products \
  -d '{"id":"2","name":"iPhone 16","price":999.99}'

## List all products
curl http://localhost:8080/v1/products

## Get single product
curl http://localhost:8080/v1/products/1

## Delete product
curl -X DELETE http://localhost:8080/v1/products/1
```
