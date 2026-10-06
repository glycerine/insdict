package insdict

import (
	"fmt"
	"testing"
)

func TestIndexTagStorage(t *testing.T) {
	for _, size := range []int{8, 16, 64, 256, 1024} {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			d := NewDict[int, int]()
			d.rebuild(size)
			if len(d.tags) != size || cap(d.tags) != size || cap(d.indices) != size {
				t.Fatal("metadata slices have incorrect lengths or capacities")
			}
			// Exercise the entire tag tail, including its last byte. With
			// checkptr enabled, an undersized allocation fails at construction.
			for i := range d.tags {
				d.tags[i] = 0xff
			}
			for i, ix := range d.indices {
				if ix != slotEmpty {
					t.Fatalf("tag writes overlap index %d: got %d", i, ix)
				}
			}
			clear(d.tags)
			for k := 0; k < d.usable(); k++ {
				d.Put(k, k+1)
			}
			checkInvariants(t, d)
			for k := 0; k < d.Len(); k++ {
				if got := d.Get(k); got != k+1 {
					t.Fatalf("Get(%d) = %d, want %d", k, got, k+1)
				}
			}
		})
	}
}

func TestEntryLimitSizing(t *testing.T) {
	// Test the boundary arithmetic without allocating multi-gigabyte tables.
	for _, n := range []int{0, 1, 5, 1000, 1_000_000, maxEntries - 1, maxEntries} {
		size := presizeFor(n)
		if size > maxTableSize || size&(size-1) != 0 || usableFor(size) < n {
			t.Fatalf("presizeFor(%d) = %d, usable=%d", n, size, usableFor(size))
		}
		if size > minSize && usableFor(size/2) >= n {
			t.Fatalf("presizeFor(%d) did not choose the smallest table", n)
		}
		grown := sizeFor(n)
		if grown > maxTableSize || grown&(grown-1) != 0 || usableFor(grown) < n {
			t.Fatalf("sizeFor(%d) = %d, usable=%d", n, grown, usableFor(grown))
		}
	}
	if got := usableFor(maxTableSize); got != maxEntries {
		t.Fatalf("maximum usable capacity = %d, want %d", got, maxEntries)
	}
	for _, tc := range []struct {
		name string
		call func()
	}{
		{"presized", func() { NewDictSize[int, int](maxEntries + 1) }},
		{"custom_hash", func() { NewDictFuncSize[int, int](EasyHashInt, maxEntries+1) }},
		{"rebuild_capacity", func() { NewDict[int, int]().rebuildEntries(minSize, maxEntries+1) }},
		{"rebuild_size", func() { NewDict[int, int]().rebuildEntries(maxTableSize+1, 0) }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("oversized table request did not panic")
				}
			}()
			tc.call()
		})
	}
}
