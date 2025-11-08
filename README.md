# ParseInsight

**HTTP Parser Performance Analyzer** - A web service where users submit their actual HTTP messages and test how different parsing libraries perform.

## Features

- **Real-time Benchmarking**: Test your actual HTTP messages against multiple parser libraries
- **Comparative Analysis**: Visual comparison of parser performance with charts and tables
- **Multiple Parsers**: Currently supports `net/http` (stdlib) and `fasthttp`
- **Detailed Metrics**: Measures throughput, latency, memory allocation, and allocations per operation
- **Concurrent Testing**: Support for concurrent benchmarks to simulate real-world scenarios
- **Auto-detection**: Automatically detects if your message is a request or response
- **Web Interface**: Clean, responsive UI for easy testing

## Architecture

```
┌─────────────────────────────────────────────┐
│         User Input Interface                │
│  - Raw HTTP message textarea/file upload    │
│  - Message validation                       │
│  - Library selection                        │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│      Performance Testing Engine             │
│  - Parse user's HTTP message                │
│  - Run benchmarks across libraries          │
│  - Measure: throughput, latency, memory     │
└─────────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────────┐
│         Results Presentation                │
│  - Comparative performance metrics          │
│  - Charts/graphs                            │
│  - Winner recommendation                    │
│  - Detailed breakdown per library           │
└─────────────────────────────────────────────┘
```

## Project Structure

```
ParseInsight/
├── cmd/
│   └── server/
│       └── main.go                 # HTTP server entry point
├── internal/
│   ├── adapters/
│   │   ├── adapter.go              # Interface definition
│   │   ├── nethttp.go              # stdlib adapter
│   │   ├── fasthttp.go             # fasthttp adapter
│   │   └── registry.go             # Adapter registry
│   ├── benchmark/
│   │   ├── models.go               # Data models
│   │   ├── runner.go               # Benchmark execution
│   │   └── validator.go            # Input validation
│   └── api/
│       ├── handlers.go             # HTTP handlers
│       └── models.go               # Request/response models
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── styles.css          # Application styles
│   │   └── js/
│   │       └── app.js              # Frontend logic
│   └── index.html                  # Main UI
├── go.mod
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Modern web browser

### Installation

1. Clone the repository:
```bash
git clone https://github.com/ABAlosaimi/ParseInsight.git
cd ParseInsight
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run cmd/server/main.go
```

4. Open your browser and navigate to:
```
http://localhost:8080
```

## Usage

### Web Interface

1. **Enter HTTP Message**: Paste your raw HTTP request or response
2. **Configure Test**: Set iterations, concurrency, and select parsers to test
3. **Run Benchmark**: Click "Run Benchmark" to start testing
4. **View Results**: See comparative charts, tables, and recommendations

### Example HTTP Request

```http
GET /api/users HTTP/1.1
Host: example.com
User-Agent: Mozilla/5.0
Accept: application/json
Authorization: Bearer token123
```

### Example HTTP Response

```http
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 27

{"message":"Hello World"}
```

### API Endpoints

#### POST /api/benchmark

Run a benchmark test on an HTTP message.

**Request:**
```json
{
  "message": "GET /api/users HTTP/1.1\r\nHost: example.com\r\n\r\n",
  "message_type": "request",
  "iterations": 10000,
  "concurrency": 1,
  "libraries": ["net/http", "fasthttp"]
}
```

**Response:**
```json
{
  "results": [
    {
      "library": "fasthttp",
      "ops_per_second": 1250000,
      "avg_time_per_parse": 800,
      "memory_allocated": 0,
      "allocs_per_op": 0,
      "success": true,
      "winner": true
    },
    {
      "library": "net/http",
      "ops_per_second": 450000,
      "avg_time_per_parse": 2222,
      "memory_allocated": 416,
      "allocs_per_op": 1,
      "success": true
    }
  ],
  "recommendation": "fasthttp is 2.78x faster than net/http (1250000 vs 450000 ops/sec)",
  "message_type": "request"
}
```

#### GET /api/libraries

Get list of available parser libraries.

**Response:**
```json
{
  "libraries": ["net/http", "fasthttp"]
}
```

## Configuration Limits

- **Max Message Size**: 1 MB
- **Max Iterations**: 10,000,000
- **Max Concurrency**: 100
- **Max Duration**: 30 seconds

## Metrics Explained

- **Ops/Second**: Number of parse operations per second (higher is better)
- **Avg Time/Parse**: Average time to parse a single message (lower is better)
- **Memory Allocated**: Total memory allocated during benchmark
- **Allocs/Op**: Number of memory allocations per operation (lower is better)

## Adding New Parsers

To add a new HTTP parser library:

1. Create a new adapter in `internal/adapters/`:

```go
type MyParserAdapter struct {
    messageType string
}

func (a *MyParserAdapter) Name() string {
    return "myparser"
}

func (a *MyParserAdapter) Parse(raw []byte) error {
    // Implement parsing logic
}

func (a *MyParserAdapter) BenchmarkParse(raw []byte, iterations int) BenchmarkResult {
    // Implement benchmarking logic
}
```

2. Register it in `internal/adapters/registry.go`:

```go
r.Register("myparser", func(mt string) ParserAdapter {
    return NewMyParserAdapter(mt)
})
```

## Development

### Run Tests
```bash
go test ./...
```

### Build
```bash
go build -o parseinsight cmd/server/main.go
```

### Run Production Build
```bash
./parseinsight
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - feel free to use this project for any purpose.

## Acknowledgments

- Built with [Go](https://golang.org/)
- Uses [fasthttp](https://github.com/valyala/fasthttp) for high-performance HTTP parsing
- Charts powered by [Chart.js](https://www.chartjs.org/)

## Support

For issues, questions, or suggestions, please open an issue on GitHub.

---

Built with ❤️ for HTTP parser enthusiasts
