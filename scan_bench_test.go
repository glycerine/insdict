package insdict

import (
	"fmt"
	"testing"
)

// BenchmarkScanCompare isolates each summation in a non-inlined function so
// construction does not affect the scan loop's stack frame. The slice is a
// sequential-access reference; both containers sum identical int values.
//
//	go test -run '^$' -bench '^BenchmarkScanCompare$' -benchmem -count=3
//
// For Linux perf, build once and select one leaf with a fixed scan count:
//
//	go test -c -o /tmp/insdict-scan.test
//	GOMAXPROCS=1 perf stat -e cycles,instructions,branches,branch-misses -- \
//	  /tmp/insdict-scan.test -test.run='^$' \
//	  -test.bench='^BenchmarkScanCompare$/^n=1000000$/^map$' \
//	  -test.benchtime=1000x -test.cpu=1
//
// External perf includes startup and construction; use enough scans to make
// these costs small. Benchmark timings exclude construction. CPU pinning with
// taskset can reduce run-to-run variation.
func BenchmarkScanCompare(b *testing.B) {
	for _, n := range []int{1000, 1_000_000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.Run("dict", func(b *testing.B) {
				d := NewDictSize[int, int](n)
				for i := 0; i < n; i++ {
					d.Put(i, i+1)
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkSink = sumDictForBenchmark(d)
				}
				b.StopTimer()
				reportIterNsPerKey(b, n)
			})
			b.Run("map", func(b *testing.B) {
				m := make(map[int]int, n)
				for i := 0; i < n; i++ {
					m[i] = i + 1
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkSink = sumMapForBenchmark(m)
				}
				b.StopTimer()
				reportIterNsPerKey(b, n)
			})
			b.Run("slice", func(b *testing.B) {
				values := make([]int, n)
				for i := range values {
					values[i] = i + 1
				}
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkSink = sumSliceForBenchmark(values)
				}
				b.StopTimer()
				reportIterNsPerKey(b, n)
			})
		})
	}
}

//go:noinline
func sumMapForBenchmark(m map[int]int) int {
	sum := 0
	for _, v := range m {
		sum += v
	}
	return sum
}

//go:noinline
func sumSliceForBenchmark(values []int) int {
	sum := 0
	for _, v := range values {
		sum += v
	}
	return sum
}
