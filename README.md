# goping

[![Go Reference](https://pkg.go.dev/badge/github.com/mameen7/goping.svg)](https://pkg.go.dev/github.com/mameen7/goping)
[![Go Report Card](https://goreportcard.com/badge/github.com/mameen7/goping)](https://goreportcard.com/report/github.com/mameen7/goping)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](LICENSE)

A fast and lightweight ping utility in Go that supports **ICMP ping** and gracefully **falls back to TCP ping** when ICMP is not permitted.

Use it as:
- A **library** in your Go applications.
- A **CLI tool** as a drop-in replacement for `ping`.

---

## ✨ Features

- ICMP ping (with `setcap` support for non-root use)
- Automatic fallback to TCP ping when ICMP is blocked
- CLI tool with configurable options
- Library with a clean, typed API
- Zero-allocation loop trick for performance

---

## 📦 Installation

### CLI
```bash
go install github.com/mameen7/goping/cmd/goping@latest
```

You might need to add Go bin directory to your path
``` bash
export PATH=$PATH:$(go env GOPATH)/bin
```


### Library
```bash
go get github.com/mameen7/goping
```

---

## 🚀 Usage

### CLI

Run ping with a count:
```bash
goping -c 3 google.com
```

Output:
```
PING google.com (142.250.64.78): 56 data bytes
64 bytes from 142.250.64.78: icmp_seq=0 time=21.3ms
64 bytes from 142.250.64.78: icmp_seq=1 time=20.7ms
64 bytes from 142.250.64.78: icmp_seq=2 time=19.9ms

--- google.com ping statistics ---
3 packets transmitted, 3 received, 0.0% packet loss
RTT MIN/AVG/MAX = 19.9/20.6/21.3 ms
```

Available flags:
```
-c int
    Number of packets to send (default 5)
-i Duration
    delay time per request (default 1s)
-t Duration
    request timeout duration (default 10s)
-s int
    packet byte size to send per request (default 64)
```

### Library

Import and use inside your Go code:
```go
package main

import (
    "fmt"
    "github.com/mameen7/goping"
)

func main() {
    p := goping.NewPinger(
		host,                               // e.g google.com
		goping.WithCount(count),            // optional count [int] number of payloads to send (default 5)
		goping.WithInterval(interval),      // optional interval [time.Duration] delay time per request (default 1s)
		goping.WithTimeout(timeout),        // optional timeout [time.Duration] request timeout duration (default 10s)
		goping.WithSize(size),              // optional size [int] packet byte size to send per request (default 64)
	)
    stats, err := p.Ping()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Sent: %d, Received: %d, Loss: %.2f%%\n", stats.Sent, stats.Recv, stats.Loss)
    fmt.Printf("RTT Min/Avg/Max: %v / %v / %v\n", stats.MinRTT, stats.AvgRTT, stats.MaxRTT)
}
```

### Stats struct
```go
type Stats struct {
    Sent   int
    Recv   int
    MinRTT time.Duration
    AvgRTT time.Duration
    MaxRTT time.Duration
    Loss   float64
    RTTs   []time.Duration
}
```

---

## 🧪 Running Tests

```bash
go test ./...
```

---

## 🤝 Contributing

Contributions are welcome! Please open an issue or PR to improve features, documentation, or testing.