package insdict

import (
	"math/rand"
	"testing"
)

func scalarGroupMasks(ctrl *[groupWidth]byte, tag byte) (mask uint64) {
	for i, c := range ctrl {
		if c == tag {
			mask |= 1 << i
		}
		if c == tagEmpty {
			mask |= 1 << (i + 16)
		}
		if c < tagUsed {
			mask |= 1 << (i + 32)
		}
	}
	return
}

func TestGroupMasks(t *testing.T) {
	check := func(ctrl *[groupWidth]byte, tag byte) {
		t.Helper()
		want := scalarGroupMasks(ctrl, tag)
		if got := groupMasks(ctrl, tag); got != want {
			t.Fatalf("groupMasks(%x, %x) = %012x; want %012x", ctrl, tag, got, want)
		}
		if got := groupMasksPortable(ctrl, tag); got != want {
			t.Fatalf("portable(%x, %x) = %012x; want %012x", ctrl, tag, got, want)
		}
	}
	// Exercise every byte in every lane, including 0/1 adjacency and repeated
	// fingerprints differing by one bit, which catch SWAR borrow mistakes.
	for tag := 0; tag < 256; tag++ {
		for lane := 0; lane < groupWidth; lane++ {
			for value := 0; value < 256; value++ {
				var ctrl [groupWidth]byte
				for i := range ctrl {
					ctrl[i] = byte(tag) ^ byte(i&1)
				}
				ctrl[lane] = byte(value)
				check(&ctrl, byte(tag))
			}
		}
	}
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 10000; i++ {
		var ctrl [groupWidth]byte
		rng.Read(ctrl[:])
		check(&ctrl, byte(rng.Intn(256)))
	}
}

func TestGroupedProbeCoverage(t *testing.T) {
	for groups := 1; groups <= 1024; groups *= 2 {
		for start := 0; start < groups; start++ {
			seen := make([]bool, groups)
			base, stride := start*groupWidth, 0
			for i := 0; i < groups; i++ {
				if seen[base/groupWidth] {
					t.Fatalf("groups=%d start=%d revisited group %d", groups, start, base/groupWidth)
				}
				seen[base/groupWidth] = true
				stride += groupWidth
				base = (base + stride) & (groups*groupWidth - 1)
			}
		}
	}
}

func TestGroupedTombstones(t *testing.T) {
	if !groupedProbing {
		t.Skip("grouped-probing experiment")
	}
	// Begin in the last group, fill it, and wrap into the first group.
	d := NewDictFuncSize[int, int](func(int) uint64 { return 63 }, 40)
	for k := 0; k < 40; k++ {
		d.Put(k, k+1)
	}
	for k := 0; k < groupWidth; k++ {
		d.Del(k)
	}
	for k := groupWidth; k < 40; k++ {
		if v, ok := d.Get2(k); !ok || v != k+1 {
			t.Fatalf("key %d beyond deleted group: (%d, %v)", k, v, ok)
		}
		if d.Put(k, k+1) {
			t.Fatalf("key %d duplicated across a tombstone", k)
		}
	}
	for k := 40; k < 56; k++ {
		d.Put(k, k+1)
	}
	checkInvariants(t, d)
	clone := d.Clone()
	clone.Pack(true)
	checkInvariants(t, clone)
}

var groupMaskSink uint64

func BenchmarkGroupMasks(b *testing.B) {
	var ctrl [groupWidth]byte
	for i := range ctrl {
		ctrl[i] = byte(i) | tagUsed
	}
	for _, impl := range []string{"selected", "portable"} {
		b.Run(impl, func(b *testing.B) {
			var sum uint64
			if impl == "selected" {
				for i := 0; i < b.N; i++ {
					sum += groupMasks(&ctrl, byte(i)|tagUsed)
				}
			} else {
				for i := 0; i < b.N; i++ {
					sum += groupMasksPortable(&ctrl, byte(i)|tagUsed)
				}
			}
			groupMaskSink = sum
		})
	}
}
