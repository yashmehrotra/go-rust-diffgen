package main

import (
	"fmt"
	"math"
	_ "net/http/pprof"
	"os"
	"runtime"
	"testing"
	"time"
)

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

// Bytes produces a human readable representation of an SI size.
//
// See also: ParseBytes.
//
// Bytes(82854982) -> 83 MB
func Bytes(s uint64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

func BenchmarkRustDiffGenerator(b *testing.B) {
	var (
		maxTotal uint64 = 0
		maxInUse uint64 = 0
		maxAlloc uint64 = 0
	)

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			time.Sleep(10 * time.Millisecond)

			go func(a, b, c uint64) {
				maxTotal = max(maxTotal, a)
				maxInUse = max(maxInUse, b)
				maxAlloc = max(maxAlloc, c)
			}(m.TotalAlloc, m.HeapAlloc, m.HeapInuse)
		}
	}()

	before, err := os.ReadFile("echo_server.json")
	if err != nil {
		b.Fatalf("failed to open file echo-server.json: %v", err)
	}

	after, err := os.ReadFile("echo_server_new.json")
	if err != nil {
		b.Fatalf("failed to open file echo-server.json: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		RustDiff(string(before), string(after))
	}

	fmt.Printf("\nTotalAlloc: %s\nHeapAlloc: %s\nHeapInuse: %s\n", Bytes(maxTotal), Bytes(maxAlloc), Bytes(maxInUse))
}

func BenchmarkGoDiffGenerator(b *testing.B) {
	var (
		maxTotal uint64 = 0
		maxInUse uint64 = 0
		maxAlloc uint64 = 0
	)

	go func() {
		for {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			time.Sleep(10 * time.Millisecond)

			go func(a, b, c uint64) {
				maxTotal = max(maxTotal, a)
				maxInUse = max(maxInUse, b)
				maxAlloc = max(maxAlloc, c)
			}(m.TotalAlloc, m.HeapAlloc, m.HeapInuse)
		}
	}()

	before, err := os.ReadFile("echo_server.json")
	if err != nil {
		b.Fatalf("failed to open file echo_server.json: %v", err)
	}

	after, err := os.ReadFile("echo_server_new.json")
	if err != nil {
		b.Fatalf("failed to open file echo_server_new.json: %v", err)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		GoDiff(string(before), string(after))
	}

	fmt.Printf("\nTotalAlloc: %s\nHeapAlloc: %s\nHeapInuse: %s\n", Bytes(maxTotal), Bytes(maxAlloc), Bytes(maxInUse))
}
