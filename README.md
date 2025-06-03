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
 docker run -p 8080:3401 -e SEED_HEX=0000000000000000000000000000000000000000000000000000000000000000 fasttrack/random grpc
```
```bash
 docker run -p 8081:3402 -e SEED_HEX=0000000000000000000000000000000000000000000000000000000000000000 fasttrack/random http
```

## Other

### Example HTTP Requests
```http
  GET http://localhost:8081/getRandomFloat64
```
```http
  GET http://localhost:8081/getRandomInt64?min=0&max=10

  Querystring parameters:  
  min - minimum number (inclusive)
  max - maximum number (inclusive)
```
```http
  GET http://localhost:8081/getDeterministicRandom?s=42&p=0.01,0.4,0.59
  
  Querystring parameters:
  (s)equence - the sequence number of the random number
  (p)robabilities - the set of probabilities to select an index from
```

### Generating a seed
There are several sites where a hex code can be generated.

Example: https://codebeautify.org/generate-random-hexadecimal-numbers

Simply set 'length of hex number' to 64 and generate one. 

### Validating deterministic results
The results from function DeterministicRandom can be tested for consistency by using the simulator to generate results
and then hashing the result of two runs with the same parameters.
```bash
 shasum -a 256 cmd/simulator/results/DeterministicRandom-X.csv
```
