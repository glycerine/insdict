package insdict

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"unsafe"
)

var pointGetSink int

// Every query cycle visits n distinct keys in a deterministic random order.
// Masking avoids division; queries and tables are built before timing.
// Unlike a short query ring, this still exercises large-table cache misses.
// Each ns/op represents one point lookup. Run individual leaves for pprof:
// go test -run '^$' -bench '^BenchmarkPointGetInt$/^n=1024$/^grown$/^hit$/^Dict$' -benchtime=10s -cpuprofile=cpu.out
func BenchmarkPointGetInt(b *testing.B) {
	benchmarkPointGet(b, func(i int) int { return i })
}

func BenchmarkPointGetString(b *testing.B) {
	benchmarkPointGet(b, func(i int) string { return fmt.Sprintf("key-%016x", i) })
}

func benchmarkPointGet[K comparable](b *testing.B, key func(int) K) {
	for _, n := range []int{1024, 16384, 1048576} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for _, layout := range []string{"grown", "sparse"} {
				b.Run(layout, func(b *testing.B) {
					d := NewDict[K, int]()
					m := make(map[K]int, n)
					for i := 0; i < n; i++ {
						k := key(i)
						// Nonzero values make hits distinguishable from misses.
						d.Put(k, i+1)
						m[k] = i + 1
					}
					// An intentionally larger index table tests fewer collisions
					// against a larger cache footprint (25% versus 50% load).
					if layout == "sparse" {
						d.rebuild(len(d.indices) * 2)
					}
					// Benchmark-only sidecar fingerprints preserve full int64 indexes.
					tags := make([]byte, len(d.indices))
					for slot, ix := range d.indices {
						if ix >= 0 {
							tags[slot] = byte(d.entries[ix].hash >> 56)
						}
					}
					order := rand.New(rand.NewSource(1)).Perm(n)
					for _, workload := range []string{"hit", "miss", "mixed"} {
						b.Run(workload, func(b *testing.B) {
							queries := make([]K, n)
							for i, k := range order {
								if workload == "miss" || workload == "mixed" && i&1 != 0 {
									k += n
								}
								queries[i] = key(k)
								// Equal strings have distinct storage, avoiding the
								// pointer-equality shortcut on successful lookups.
								if s, ok := any(queries[i]).(string); ok {
									queries[i] = any(strings.Clone(s)).(K)
								}
							}
							probes := pointGetProbes(d, queries)
							for _, impl := range []string{"Dict", "Dict2", "Map", "Map2", "Fused", "TypedInt", "Tagged"} {
								b.Run(impl, func(b *testing.B) {
									intDict, isInt := any(d).(*Dict[int, int])
									intQueries, _ := any(queries).([]int)
									if impl == "TypedInt" && !isInt {
										b.Skip("integer specialization experiment")
									}
									b.ReportAllocs()
									mask, sum := n-1, 0
									b.ResetTimer()
									switch impl {
									case "Dict":
										for i := 0; i < b.N; i++ {
											sum += d.Get(queries[i&mask])
										}
									case "Map":
										for i := 0; i < b.N; i++ {
											sum += m[queries[i&mask]]
										}
									case "Dict2":
										for i := 0; i < b.N; i++ {
											v, ok := d.Get2(queries[i&mask])
											sum += v
											if ok {
												sum++
											}
										}
									case "Map2":
										for i := 0; i < b.N; i++ {
											v, ok := m[queries[i&mask]]
											sum += v
											if ok {
												sum++
											}
										}
									case "Fused":
										for i := 0; i < b.N; i++ {
											sum += pointGetFused(d, queries[i&mask])
										}
									case "TypedInt":
										for i := 0; i < b.N; i++ {
											sum += pointGetTypedInt(intDict, intQueries[i&mask])
										}
									case "Tagged":
										for i := 0; i < b.N; i++ {
											sum += pointGetTagged(d, tags, queries[i&mask])
										}
									}
									b.StopTimer()
									pointGetSink = sum
									if impl != "Map" && impl != "Map2" {
										b.ReportMetric(probes, "probes/op")
										b.ReportMetric(float64(d.live)/float64(len(d.indices)), "load")
										bytes := float64(cap(d.indices)*8+cap(d.tags)) + float64(cap(d.entries))*float64(unsafe.Sizeof(entry[K, int]{}))
										if impl == "Tagged" {
											bytes += float64(len(tags))
										}
										b.ReportMetric(bytes, "table_B")
									}
								})
							}
						})
					}
				})
			}
		})
	}
}

// Benchmark-only scalar fingerprint filter. An 8-bit tag beside the index
// rejects most collisions before the dependent entries read. This is a simpler
// experiment than grouped Swiss-table probing, and uses the original sequence.
func pointGetTagged[K comparable](d *Dict[K, int], tags []byte, k K) int {
	if d.live == 0 {
		return 0
	}
	var h uint64
	if d.hash != nil {
		h = d.hash(k)
	} else {
		h = defaultHash(k)
	}
	tag := byte(h >> 56)
	i, perturb := h&d.mask, h
	for {
		ix := d.indices[i]
		if ix == slotEmpty {
			return 0
		}
		if ix >= 0 && tags[i] == tag {
			e := &d.entries[ix]
			if e.hash == h && e.key == k {
				return e.val
			}
		}
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
}

// Benchmark-only concrete int specialization. Compared with Fused, this removes
// the defaultHash call and type switch, while retaining custom-hash support.
func pointGetTypedInt(d *Dict[int, int], k int) int {
	if d.live == 0 {
		return 0
	}
	var h uint64
	if d.hash != nil {
		h = d.hash(k)
	} else {
		h = Mix64(uint64(k))
	}
	i, perturb := h&d.mask, h
	for {
		ix := d.indices[i]
		if ix == slotEmpty {
			return 0
		}
		if ix >= 0 {
			e := &d.entries[ix]
			if e.hash == h && e.key == k {
				return e.val
			}
		}
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
}

// Benchmark-only candidate: inline the hash selection and return the value
// directly from the matching entry, avoiding hashOf and repeated entry indexing.
// Keep the same hashing, table layout, and probe sequence as production Get.
func pointGetFused[K comparable](d *Dict[K, int], k K) int {
	if d.live == 0 {
		return 0
	}
	var h uint64
	if d.hash != nil {
		h = d.hash(k)
	} else {
		h = defaultHash(k)
	}
	i, perturb := h&d.mask, h
	for {
		ix := d.indices[i]
		if ix == slotEmpty {
			return 0
		}
		if ix >= 0 {
			e := &d.entries[ix]
			if e.hash == h && e.key == k {
				return e.val
			}
		}
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
}

// Count probes outside the timed loop to separate collision behavior from CPU
// overhead. This deliberately duplicates the production probe sequence.
func pointGetProbes[K comparable](d *Dict[K, int], queries []K) float64 {
	var probes int
	for _, k := range queries {
		h := d.hashOf(k)
		i, perturb := h&d.mask, h
		for {
			probes++
			ix := d.indices[i]
			if ix == slotEmpty {
				break
			}
			if ix >= 0 && d.entries[ix].hash == h && d.entries[ix].key == k {
				break
			}
			perturb >>= 5
			i = (i*5 + perturb + 1) & d.mask
		}
	}
	return float64(probes) / float64(len(queries))
}

func TestPointGetCandidates(t *testing.T) {
	for _, hash := range []func(int) uint64{nil, func(int) uint64 { return 0 }} {
		d := NewDictFunc[int, int](hash)
		if got := pointGetFused(d, 0); got != 0 {
			t.Fatalf("empty: %d", got)
		}
		if pointGetTypedInt(d, 0) != 0 || pointGetTagged(d, nil, 0) != 0 {
			t.Fatal("empty candidate lookup")
		}
		for i := 0; i < 512; i++ {
			d.Put(i, i+1)
		}
		for i := 0; i < 512; i += 2 {
			d.Del(i)
		}
		tags := make([]byte, len(d.indices))
		for slot, ix := range d.indices {
			if ix >= 0 {
				tags[slot] = byte(d.entries[ix].hash >> 56)
			}
		}
		for k := -10; k < 600; k++ {
			if got, want := pointGetFused(d, k), d.Get(k); got != want {
				t.Fatalf("key %d: got %d, want %d", k, got, want)
			}
			if got, want := pointGetTypedInt(d, k), d.Get(k); got != want {
				t.Fatalf("typed key %d: got %d, want %d", k, got, want)
			}
			if got, want := pointGetTagged(d, tags, k), d.Get(k); got != want {
				t.Fatalf("tagged key %d: got %d, want %d", k, got, want)
			}
		}
	}
	d := NewDict[string, int]()
	d.Put("", 0)
	d.Put("a", 2)
	d.Put("b", 3)
	d.Del("a")
	tags := make([]byte, len(d.indices))
	for slot, ix := range d.indices {
		if ix >= 0 {
			tags[slot] = byte(d.entries[ix].hash >> 56)
		}
	}
	for _, k := range []string{"", "a", "b", "missing"} {
		if got, want := pointGetFused(d, k), d.Get(k); got != want {
			t.Fatalf("key %q: got %d, want %d", k, got, want)
		}
		if got, want := pointGetTagged(d, tags, k), d.Get(k); got != want {
			t.Fatalf("tagged key %q: got %d, want %d", k, got, want)
		}
	}
}

// Isolate the scan from construction: otherwise the larger Dict stack frame
// changes the inlined loop's instruction alignment in existing benchmarks.
func BenchmarkDictScanValues(b *testing.B) {
	for _, n := range []int{1000, 1000000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			for _, holes := range []bool{false, true} {
				b.Run(fmt.Sprintf("holes=%v", holes), func(b *testing.B) {
					d := NewDict[int, int]()
					for i := 0; i < n; i++ {
						d.Put(i, i+1)
					}
					if holes {
						for i := 0; i < n; i += 4 {
							d.Del(i)
						}
					}
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						pointGetSink = sumDictForBenchmark(d)
					}
					b.StopTimer()
					reportIterNsPerKey(b, d.Len())
				})
			}
		})
	}
}

//go:noinline
func sumDictForBenchmark(d *Dict[int, int]) int {
	sum := 0
	for _, v := range d.All() {
		sum += v
	}
	return sum
}
