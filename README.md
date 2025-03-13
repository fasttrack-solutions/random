# Random Backend

## Development

### Generate GRPC

https://medium.com/@blackhorseya/streamlining-protobuf-code-generation-with-buf-in-golang-projects-b506316da7e2

```bash
 buf generate
```

### Running gosec
```bash
gosec -exclude-dir=pkg/pb ./...
```

## How to run locally

### Start simulator
```bash
 go run cmd/simulator/main.go
```

### Start GRPC Endpoint
```bash
 go run cmd/grpc/main.go
```

### Start HTTP Endpoint
```bash
 go run cmd/http/main.go
```

## How to run with Docker

### Build
```bash
 docker build -t fasttrack/random .
```

### Run
```bash
 docker run -p 8080:3401 fasttrack/random grpc
```
```bash
 docker run -p 8081:3402 fasttrack/random http
```
