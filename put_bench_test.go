package insdict

import (
	"fmt"
	"runtime"
	"testing"
)

// Keys and values are prepared outside timing; each operation constructs a
// fresh dictionary. The existing integer benchmarks also cover larger sizes.
func BenchmarkPutStrings(b *testing.B) {
	benchmarkPutValues(b, nil, func(i int) string { return fmt.Sprintf("key-%016x", i) }, func(i int) int { return i })
}

func BenchmarkPutLargeValues(b *testing.B) {
	benchmarkPutValues(b, nil, func(i int) int { return i }, func(i int) [128]byte { return [128]byte{byte(i), byte(i >> 8)} })
}

type putPair struct{ A, B int }

func BenchmarkPutCustomHash(b *testing.B) {
	benchmarkPutValues(b, func(k putPair) uint64 { return Mix64(uint64(k.A)) ^ Mix64(uint64(k.B)) },
		func(i int) putPair { return putPair{i, i + 1} }, func(i int) int { return i })
}

func benchmarkPutValues[K comparable, V any](b *testing.B, hash func(K) uint64, key func(int) K, value func(int) V) {
	for _, n := range []int{1024, 65536} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			keys, values := make([]K, n), make([]V, n)
			for i := 0; i < n; i++ {
				keys[i], values[i] = key(i), value(i)
			}
			for _, presized := range []bool{false, true} {
				b.Run(fmt.Sprintf("presized=%v", presized), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						d := NewDictFunc[K, V](hash)
						if presized {
							d = NewDictFuncSize[K, V](hash, n)
						}
						for j, k := range keys {
							d.Put(k, values[j])
						}
						benchmarkSink = d.Len()
					}
					reportNsPerKey(b, "put_ns/key", n)
				})
			}
		})
	}
}

// Maintain a fixed live set while deletes leave holes and repeated inserts
// eventually compact the table. Each operation is one delete plus one Put.
func BenchmarkPutChurn(b *testing.B) {
	for _, n := range []int{1024, 65536} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			d := NewDictSize[int, int](n)
			for k := 0; k < n; k++ {
				d.Put(k, k)
			}
			// Finish initial growth before measuring steady compaction/reuse.
			for i := 0; i < 2*n; i++ {
				d.Del(i & (2*n - 1))
				d.Put((i+n)&(2*n-1), i)
			}
			var before, after runtime.MemStats
			runtime.ReadMemStats(&before)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				old, next := i&(2*n-1), (i+n)&(2*n-1)
				d.Del(old)
				d.Put(next, i)
			}
			b.StopTimer()
			runtime.ReadMemStats(&after)
			// Built-in allocs/op rounds down these infrequent allocations.
			b.ReportMetric(float64(after.Mallocs-before.Mallocs)/float64(b.N), "allocs/cycle")
			benchmarkSink = d.Len()
		})
	}
}
