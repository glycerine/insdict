package insdict

import (
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"testing"
)

// ---------------------------------------------------------------------------
// Helpers: reference model and structural invariant checker
// ---------------------------------------------------------------------------

// model is the trivially-correct reference: a Go map plus an explicit order.
type model[K comparable, V any] struct {
	m     map[K]V
	order []K
}

func newModel[K comparable, V any]() *model[K, V] {
	return &model[K, V]{m: map[K]V{}}
}

func (m *model[K, V]) put(k K, v V) {
	if _, ok := m.m[k]; !ok {
		m.order = append(m.order, k)
	}
	m.m[k] = v
}

func (m *model[K, V]) del(k K) bool {
	if _, ok := m.m[k]; !ok {
		return false
	}
	delete(m.m, k)
	i := slices.Index(m.order, k)
	m.order = slices.Delete(m.order, i, i+1)
	return true
}

func keysOf[K comparable, V any](d *Dict[K, V]) []K {
	var ks []K
	for k := range d.All() {
		ks = append(ks, k)
	}
	return ks
}

// checkInvariants validates the internal structure of d.
func checkInvariants[K comparable, V any](t testing.TB, d *Dict[K, V]) {
	t.Helper()

	if d.indices == nil {
		if d.live != 0 || len(d.entries) != 0 {
			t.Fatalf("nil indices but live=%d entries=%d", d.live, len(d.entries))
		}
		return
	}

	size := len(d.indices)
	if size < minSize || size&(size-1) != 0 {
		t.Fatalf("index table size %d is not a power of two >= %d", size, minSize)
	}
	if d.mask != uint64(size-1) {
		t.Fatalf("mask %d does not match size %d", d.mask, size)
	}
	if len(d.entries) > d.usable() {
		t.Fatalf("len(entries)=%d exceeds usable=%d", len(d.entries), d.usable())
	}

	var zeroK K
	liveCount := 0
	for i := range d.entries {
		e := &d.entries[i]
		if !e.live {
			if e.hash != 0 || e.key != zeroK {
				t.Fatalf("dead entry %d not zeroed: %+v", i, *e)
			}
			continue
		}
		liveCount++
		if h := d.hashOf(e.key); h != e.hash {
			t.Fatalf("entry %d stored hash %x != recomputed %x", i, e.hash, h)
		}
		if _, ix := d.find(e.key, e.hash); int(ix) != i {
			t.Fatalf("find(%v) returned ix=%d, want %d", e.key, ix, i)
		}
	}
	if liveCount != d.live {
		t.Fatalf("live=%d but %d live entries", d.live, liveCount)
	}

	refs, empties := 0, 0
	for slot, ix := range d.indices {
		switch {
		case ix == slotEmpty:
			empties++
		case ix == slotDummy:
		case ix >= 0:
			refs++
			if int(ix) >= len(d.entries) || !d.entries[ix].live {
				t.Fatalf("slot %d points at bad/dead entry %d", slot, ix)
			}
		default:
			t.Fatalf("slot %d has invalid value %d", slot, ix)
		}
	}
	if refs != d.live {
		t.Fatalf("%d slots reference entries, want live=%d", refs, d.live)
	}
	if empties == 0 {
		t.Fatalf("no empty slot left: probing could loop forever")
	}
}

// checkAgainst verifies d matches the model in content and iteration order.
func checkAgainst[K comparable, V comparable](t testing.TB, d *Dict[K, V], m *model[K, V]) {
	t.Helper()
	checkInvariants(t, d)
	if d.Len() != len(m.m) {
		t.Fatalf("Len()=%d, want %d", d.Len(), len(m.m))
	}
	for k, want := range m.m {
		got, ok := d.Get2(k)
		if !ok || got != want {
			t.Fatalf("Get2(%v) = (%v,%v), want (%v,true)", k, got, ok, want)
		}
		if g := d.Get(k); g != want {
			t.Fatalf("Get(%v) = %v, want %v", k, g, want)
		}
	}
	if got := keysOf(d); !slices.Equal(got, m.order) {
		t.Fatalf("iteration order mismatch:\n got  %v\n want %v", got, m.order)
	}
}

// ---------------------------------------------------------------------------
// Basic unit tests
// ---------------------------------------------------------------------------

func TestEmptyDict(t *testing.T) {
	d := NewDict[string, int]()
	if d.Len() != 0 {
		t.Fatalf("Len() = %d", d.Len())
	}
	if v := d.Get("x"); v != 0 {
		t.Fatalf("Get on empty = %d", v)
	}
	if v, ok := d.Get2("x"); ok || v != 0 {
		t.Fatalf("Get2 on empty = (%d,%v)", v, ok)
	}
	if d.Del("x") {
		t.Fatal("Del on empty returned true")
	}
	if n := len(keysOf(d)); n != 0 {
		t.Fatalf("iterated %d entries on empty dict", n)
	}
	checkInvariants(t, d)
}

func TestZeroValueDict(t *testing.T) {
	var d Dict[string, int]
	if _, ok := d.Get2("a"); ok {
		t.Fatal("found key in zero-value dict")
	}
	d.Put("a", 1)
	if v, ok := d.Get2("a"); !ok || v != 1 {
		t.Fatalf("Get2 = (%d,%v)", v, ok)
	}
	checkInvariants(t, &d)
}

func TestPutGetBasic(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("a", 1)
	d.Put("b", 2)
	d.Put("c", 3)

	for k, want := range map[string]int{"a": 1, "b": 2, "c": 3} {
		if got := d.Get(k); got != want {
			t.Errorf("Get(%q) = %d, want %d", k, got, want)
		}
	}
	if _, ok := d.Get2("zzz"); ok {
		t.Error("found absent key")
	}
	if d.Len() != 3 {
		t.Errorf("Len() = %d", d.Len())
	}
	checkInvariants(t, d)
}

func TestGet2DistinguishesZeroValueFromAbsent(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("zero", 0)
	if v, ok := d.Get2("zero"); !ok || v != 0 {
		t.Fatalf("Get2(zero) = (%d,%v), want (0,true)", v, ok)
	}
	if _, ok := d.Get2("absent"); ok {
		t.Fatal("absent key reported as found")
	}
}

func TestOverwriteKeepsPosition(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("a", 1)
	d.Put("b", 2)
	d.Put("c", 3)
	d.Put("a", 100)
	d.Put("b", 200)

	if got, want := keysOf(d), []string{"a", "b", "c"}; !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	if d.Get("a") != 100 || d.Get("b") != 200 || d.Get("c") != 3 {
		t.Fatal("overwrite did not take effect")
	}
	if d.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", d.Len())
	}
	checkInvariants(t, d)
}

func TestInsertionOrder(t *testing.T) {
	d := NewDict[int, int]()
	rng := rand.New(rand.NewSource(1))
	perm := rng.Perm(500)
	for _, k := range perm {
		d.Put(k, k*2)
	}
	if got := keysOf(d); !slices.Equal(got, perm) {
		t.Fatal("iteration order differs from insertion order")
	}
	checkInvariants(t, d)
}

func TestDelete(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("a", 1)
	d.Put("b", 2)
	d.Put("c", 3)

	if !d.Del("b") {
		t.Fatal("Del(b) = false")
	}
	if d.Del("b") {
		t.Fatal("second Del(b) = true")
	}
	if d.Del("nope") {
		t.Fatal("Del(nope) = true")
	}
	if _, ok := d.Get2("b"); ok {
		t.Fatal("b still present")
	}
	if d.Len() != 2 {
		t.Fatalf("Len() = %d", d.Len())
	}
	if got, want := keysOf(d), []string{"a", "c"}; !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	checkInvariants(t, d)
}

func TestDeleteAllThenReuse(t *testing.T) {
	d := NewDict[int, int]()
	for i := 0; i < 100; i++ {
		d.Put(i, i)
	}
	for i := 0; i < 100; i++ {
		if !d.Del(i) {
			t.Fatalf("Del(%d) = false", i)
		}
	}
	if d.Len() != 0 {
		t.Fatalf("Len() = %d", d.Len())
	}
	if n := len(keysOf(d)); n != 0 {
		t.Fatalf("iterated %d entries after deleting all", n)
	}
	checkInvariants(t, d)

	d.Put(7, 70)
	if d.Get(7) != 70 || d.Len() != 1 {
		t.Fatal("dict not reusable after emptying")
	}
	checkInvariants(t, d)
}

func TestReinsertAfterDeleteGoesToEnd(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("a", 1)
	d.Put("b", 2)
	d.Put("c", 3)
	d.Del("a")
	d.Put("a", 10)

	if got, want := keysOf(d), []string{"b", "c", "a"}; !slices.Equal(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	if d.Get("a") != 10 {
		t.Fatalf("Get(a) = %d", d.Get("a"))
	}
	checkInvariants(t, d)
}

func TestDeleteZeroesEntryForGC(t *testing.T) {
	d := NewDict[string, *int]()
	x := 42
	d.Put("k", &x)
	slot, ix := d.find("k", d.hashOf("k"))
	_ = slot
	d.Del("k")
	if e := d.entries[ix]; e.live || e.val != nil || e.key != "" || e.hash != 0 {
		t.Fatalf("entry not zeroed after Del: %+v", e)
	}
	checkInvariants(t, d)
}

// ---------------------------------------------------------------------------
// Growth, compaction, churn
// ---------------------------------------------------------------------------

func TestGrowth(t *testing.T) {
	d := NewDict[int, int]()
	m := newModel[int, int]()
	for i := 0; i < 20000; i++ {
		k := i * 7919
		d.Put(k, i)
		m.put(k, i)
		if i%997 == 0 {
			checkAgainst(t, d, m)
		}
	}
	checkAgainst(t, d, m)
}

func TestShrinkAfterMassDelete(t *testing.T) {
	d := NewDict[int, int]()
	m := newModel[int, int]()
	const n = 1000
	for i := 0; i < n; i++ {
		d.Put(i, i)
		m.put(i, i)
	}
	// Delete all but every 100th key.
	for i := 0; i < n; i++ {
		if i%100 != 0 {
			d.Del(i)
			m.del(i)
		}
	}
	checkAgainst(t, d, m)
	if len(d.entries) >= n {
		t.Fatalf("entries not compacted: len=%d", len(d.entries))
	}
	if len(d.entries) > 32 && d.live < len(d.entries)/4 {
		t.Fatalf("shrink condition still true: live=%d entries=%d", d.live, len(d.entries))
	}
}

func TestOrderPreservedAcrossCompaction(t *testing.T) {
	d := NewDict[int, int]()
	m := newModel[int, int]()
	for i := 0; i < 200; i++ {
		d.Put(i, i)
		m.put(i, i)
	}
	// Remove the odd keys, forcing at least one rebuild along the way.
	for i := 1; i < 200; i += 2 {
		d.Del(i)
		m.del(i)
	}
	checkAgainst(t, d, m)
	// Add more after compaction.
	for i := 1000; i < 1100; i++ {
		d.Put(i, i)
		m.put(i, i)
	}
	checkAgainst(t, d, m)
}

// A rolling window of puts/deletes repeatedly fills the entries array with
// holes. If dummies were ever unbounded, find/freeSlot would spin forever.
func TestChurnTerminatesAndStaysConsistent(t *testing.T) {
	d := NewDict[int, int]()
	m := newModel[int, int]()
	const window = 10
	for i := 0; i < 200000; i++ {
		d.Put(i, i)
		m.put(i, i)
		if i >= window {
			d.Del(i - window)
			m.del(i - window)
		}
		if i%20011 == 0 {
			checkAgainst(t, d, m)
		}
	}
	checkAgainst(t, d, m)
	if d.Len() != window {
		t.Fatalf("Len() = %d, want %d", d.Len(), window)
	}
}

func TestRepeatedPutDeleteSameKey(t *testing.T) {
	d := NewDict[string, int]()
	for i := 0; i < 50000; i++ {
		d.Put("k", i)
		if !d.Del("k") {
			t.Fatalf("iteration %d: Del = false", i)
		}
	}
	checkInvariants(t, d)
	if d.Len() != 0 {
		t.Fatalf("Len() = %d", d.Len())
	}
}

func TestSizeFor(t *testing.T) {
	for live := 0; live < 5000; live++ {
		size := sizeFor(live)
		if size < minSize || size&(size-1) != 0 {
			t.Fatalf("sizeFor(%d) = %d: not a power of two >= minSize", live, size)
		}
		if size < live*3 {
			t.Fatalf("sizeFor(%d) = %d: too small", live, size)
		}
		// After rebuild, live entries must fit in usable capacity with room to spare.
		if live > size*2/3 {
			t.Fatalf("sizeFor(%d) = %d: usable=%d < live", live, size, size*2/3)
		}
	}
}

// ---------------------------------------------------------------------------
// Hashing
// ---------------------------------------------------------------------------

func TestCollidingHash(t *testing.T) {
	// Every key lands in one of 3 buckets: stresses probing, tombstones, and
	// full-key equality (as opposed to hash equality alone).
	d := NewDictFunc[int, int](func(k int) uint64 { return uint64(k % 3) })
	m := newModel[int, int]()
	for i := 0; i < 300; i++ {
		d.Put(i, i)
		m.put(i, i)
	}
	checkAgainst(t, d, m)
	for i := 0; i < 300; i += 2 {
		d.Del(i)
		m.del(i)
	}
	checkAgainst(t, d, m)
	for i := 0; i < 300; i += 4 {
		d.Put(i, -i)
		m.put(i, -i)
	}
	checkAgainst(t, d, m)
}

func TestConstantHash(t *testing.T) {
	d := NewDictFunc[string, int](func(string) uint64 { return 0 })
	m := newModel[string, int]()
	for i := 0; i < 100; i++ {
		k := fmt.Sprintf("key-%d", i)
		d.Put(k, i)
		m.put(k, i)
	}
	checkAgainst(t, d, m)
	for i := 0; i < 100; i += 3 {
		k := fmt.Sprintf("key-%d", i)
		d.Del(k)
		m.del(k)
	}
	checkAgainst(t, d, m)
}

func TestStructKeysWithCustomHash(t *testing.T) {
	type Point struct{ X, Y int32 }
	d := NewDictFunc[Point, string](func(p Point) uint64 {
		return Mix64(uint64(uint32(p.X))<<32 | uint64(uint32(p.Y)))
	})
	m := newModel[Point, string]()
	for x := int32(-20); x < 20; x++ {
		for y := int32(-20); y < 20; y++ {
			p := Point{x, y}
			s := fmt.Sprintf("%d,%d", x, y)
			d.Put(p, s)
			m.put(p, s)
		}
	}
	checkAgainst(t, d, m)
	if got := d.Get(Point{3, -4}); got != "3,-4" {
		t.Fatalf("Get = %q", got)
	}
}

func TestDefaultHashSupportedTypes(t *testing.T) {
	t.Run("string", func(t *testing.T) { roundTrip(t, []string{"", "a", "b", "héllo", "a\x00b"}) })
	t.Run("int", func(t *testing.T) { roundTrip(t, []int{0, 1, -1, 1 << 40, -(1 << 40)}) })
	t.Run("int8", func(t *testing.T) { roundTrip(t, []int8{0, 1, -1, 127, -128}) })
	t.Run("int16", func(t *testing.T) { roundTrip(t, []int16{0, 1, -1, 32767, -32768}) })
	t.Run("int32", func(t *testing.T) { roundTrip(t, []int32{0, 1, -1, 1 << 30, -(1 << 30)}) })
	t.Run("int64", func(t *testing.T) { roundTrip(t, []int64{0, 1, -1, 1 << 62, -(1 << 62)}) })
	t.Run("uint", func(t *testing.T) { roundTrip(t, []uint{0, 1, 1 << 40}) })
	t.Run("uint8", func(t *testing.T) { roundTrip(t, []uint8{0, 1, 255}) })
	t.Run("uint16", func(t *testing.T) { roundTrip(t, []uint16{0, 1, 65535}) })
	t.Run("uint32", func(t *testing.T) { roundTrip(t, []uint32{0, 1, 1 << 31}) })
	t.Run("uint64", func(t *testing.T) { roundTrip(t, []uint64{0, 1, 1 << 63}) })
	t.Run("bool", func(t *testing.T) { roundTrip(t, []bool{false, true}) })
}

func roundTrip[K comparable](t *testing.T, keys []K) {
	t.Helper()
	d := NewDict[K, int]()
	for i, k := range keys {
		d.Put(k, i+1)
	}
	for i, k := range keys {
		if got, ok := d.Get2(k); !ok || got != i+1 {
			t.Fatalf("Get2(%v) = (%d,%v), want (%d,true)", k, got, ok, i+1)
		}
	}
	if d.Len() != len(keys) {
		t.Fatalf("Len() = %d, want %d", d.Len(), len(keys))
	}
	checkInvariants(t, d)
}

func TestDefaultHashPanicsOnUnsupportedKey(t *testing.T) {
	type S struct{ A int }
	type ID uint64 // named types don't match the type switch

	mustPanic := func(name string, f func()) {
		t.Helper()
		defer func() {
			if recover() == nil {
				t.Errorf("%s: expected panic", name)
			}
		}()
		f()
	}
	mustPanic("struct", func() { NewDict[S, int]().Put(S{1}, 1) })
	mustPanic("named uint64", func() { NewDict[ID, int]().Put(ID(1), 1) })
	mustPanic("float64", func() { NewDict[float64, int]().Put(1.5, 1) })
}

func TestCustomHashOverridesDefault(t *testing.T) {
	// A named type that the default hash can't handle works with NewDictFunc.
	type ID uint64
	d := NewDictFunc[ID, string](func(id ID) uint64 { return Mix64(uint64(id)) })
	d.Put(1, "one")
	d.Put(2, "two")
	if d.Get(1) != "one" || d.Get(2) != "two" {
		t.Fatal("custom hash lookups failed")
	}
}

func TestMix64IsInjectiveOnSample(t *testing.T) {
	// splitmix64's finalizer is a bijection on uint64, so there must be no
	// collisions among distinct inputs.
	seen := make(map[uint64]uint64, 100000)
	for i := uint64(0); i < 100000; i++ {
		h := Mix64(i)
		if prev, dup := seen[h]; dup {
			t.Fatalf("Mix64 collision: %d and %d", prev, i)
		}
		seen[h] = i
	}
}

func TestHashDeterministic(t *testing.T) {
	for i := 0; i < 1000; i++ {
		k := fmt.Sprintf("key-%d", i)
		if defaultHash(k) != defaultHash(k) {
			t.Fatalf("defaultHash(%q) not stable", k)
		}
		if defaultHash(i) != defaultHash(i) {
			t.Fatalf("defaultHash(%d) not stable", i)
		}
	}
}

// Two dicts fed the same operation sequence must end up bit-for-bit identical.
// This is what makes the structure usable inside deterministic simulation.
func TestDeterministicLayout(t *testing.T) {
	build := func() *Dict[string, int] {
		d := NewDict[string, int]()
		rng := rand.New(rand.NewSource(99))
		for i := 0; i < 5000; i++ {
			k := fmt.Sprintf("k%d", rng.Intn(400))
			if rng.Intn(3) == 0 {
				d.Del(k)
			} else {
				d.Put(k, i)
			}
		}
		return d
	}
	a, b := build(), build()
	if !reflect.DeepEqual(a.indices, b.indices) {
		t.Fatal("index tables differ between identical runs")
	}
	if !reflect.DeepEqual(a.entries, b.entries) {
		t.Fatal("entry arrays differ between identical runs")
	}
}

// ---------------------------------------------------------------------------
// Iteration
// ---------------------------------------------------------------------------

func TestAllEarlyBreak(t *testing.T) {
	d := NewDict[int, int]()
	for i := 0; i < 10; i++ {
		d.Put(i, i)
	}
	var got []int
	for k := range d.All() {
		got = append(got, k)
		if len(got) == 3 {
			break
		}
	}
	if !slices.Equal(got, []int{0, 1, 2}) {
		t.Fatalf("got %v", got)
	}
}

func TestAllYieldsValues(t *testing.T) {
	d := NewDict[string, int]()
	d.Put("x", 10)
	d.Put("y", 20)
	d.Del("x")
	d.Put("z", 30)
	var ks []string
	var vs []int
	for k, v := range d.All() {
		ks = append(ks, k)
		vs = append(vs, v)
	}
	if !slices.Equal(ks, []string{"y", "z"}) || !slices.Equal(vs, []int{20, 30}) {
		t.Fatalf("got keys=%v vals=%v", ks, vs)
	}
}

func TestAllSkipsHoles(t *testing.T) {
	d := NewDict[int, int]()
	for i := 0; i < 20; i++ {
		d.Put(i, i)
	}
	// Few enough deletes that no compaction triggers (len(entries) <= 32).
	for i := 0; i < 20; i += 2 {
		d.Del(i)
	}
	want := []int{1, 3, 5, 7, 9, 11, 13, 15, 17, 19}
	if got := keysOf(d); !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestPutDuringIteration(t *testing.T) {
	d := NewDict[int, int]()
	const n = 10
	for i := 0; i < n; i++ {
		d.Put(i, i)
	}
	var visited []int
	first := true
	for k := range d.All() {
		visited = append(visited, k)
		if first {
			first = false
			for j := 100; j < 200; j++ { // forces table growth mid-iteration
				d.Put(j, j)
			}
		}
	}
	if len(visited) != n+100 {
		t.Fatalf("visited %d entries, want %d", len(visited), n+100)
	}
	for i := 0; i < n; i++ {
		if visited[i] != i {
			t.Fatalf("original keys visited out of order: %v", visited[:n])
		}
	}
	checkInvariants(t, d)
}

// ---------------------------------------------------------------------------
// Property-based tests (deterministic seeds, compared against a model)
// ---------------------------------------------------------------------------

func TestPropertyRandomOpsMatchModel(t *testing.T) {
	for seed := int64(0); seed < 50; seed++ {
		seed := seed
		t.Run(fmt.Sprintf("seed%d", seed), func(t *testing.T) {
			rng := rand.New(rand.NewSource(seed))
			keySpace := 1 + rng.Intn(500) // small spaces => heavy overwrite/delete churn
			d := NewDict[int, int]()
			m := newModel[int, int]()
			ops := 3000
			for i := 0; i < ops; i++ {
				k := rng.Intn(keySpace)
				switch rng.Intn(10) {
				case 0, 1, 2, 3, 4: // put
					d.Put(k, i)
					m.put(k, i)
				case 5, 6, 7: // delete
					if got, want := d.Del(k), m.del(k); got != want {
						t.Fatalf("op %d: Del(%d) = %v, want %v", i, k, got, want)
					}
				default: // get
					want, wok := m.m[k]
					got, ok := d.Get2(k)
					if ok != wok || got != want {
						t.Fatalf("op %d: Get2(%d) = (%d,%v), want (%d,%v)", i, k, got, ok, want, wok)
					}
				}
				if i%101 == 0 {
					checkAgainst(t, d, m)
				}
			}
			checkAgainst(t, d, m)
		})
	}
}

func TestPropertyStringKeysMatchModel(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	d := NewDict[string, int]()
	m := newModel[string, int]()
	randKey := func() string {
		n := rng.Intn(6)
		b := make([]byte, n)
		for i := range b {
			b[i] = byte('a' + rng.Intn(4))
		}
		return string(b)
	}
	for i := 0; i < 20000; i++ {
		k := randKey()
		if rng.Intn(3) == 0 {
			if got, want := d.Del(k), m.del(k); got != want {
				t.Fatalf("Del(%q) = %v, want %v", k, got, want)
			}
		} else {
			d.Put(k, i)
			m.put(k, i)
		}
	}
	checkAgainst(t, d, m)
}

func TestPropertyCollidingHashMatchesModel(t *testing.T) {
	for seed := int64(0); seed < 20; seed++ {
		rng := rand.New(rand.NewSource(seed))
		buckets := uint64(1 + rng.Intn(8))
		d := NewDictFunc[int, int](func(k int) uint64 { return uint64(k) % buckets })
		m := newModel[int, int]()
		for i := 0; i < 2000; i++ {
			k := rng.Intn(150)
			if rng.Intn(3) == 0 {
				if got, want := d.Del(k), m.del(k); got != want {
					t.Fatalf("seed %d: Del(%d) = %v, want %v", seed, k, got, want)
				}
			} else {
				d.Put(k, i)
				m.put(k, i)
			}
		}
		checkAgainst(t, d, m)
	}
}

// Put-then-Get, Put-twice-keeps-last, and Del-then-Get properties on random inputs.
func TestPropertyBasicLaws(t *testing.T) {
	rng := rand.New(rand.NewSource(123))
	for i := 0; i < 200; i++ {
		d := NewDict[int, int]()
		for j := 0; j < rng.Intn(100); j++ {
			d.Put(rng.Intn(1000), j)
		}
		k, v1, v2 := rng.Intn(1000), rng.Int(), rng.Int()

		d.Put(k, v1)
		if got, ok := d.Get2(k); !ok || got != v1 {
			t.Fatalf("Put/Get law violated: (%d,%v) vs %d", got, ok, v1)
		}
		n := d.Len()
		d.Put(k, v2)
		if d.Get(k) != v2 || d.Len() != n {
			t.Fatalf("overwrite law violated")
		}
		if !d.Del(k) {
			t.Fatalf("Del of present key returned false")
		}
		if _, ok := d.Get2(k); ok || d.Len() != n-1 {
			t.Fatalf("Del law violated")
		}
		if d.Del(k) {
			t.Fatalf("second Del returned true")
		}
		checkInvariants(t, d)
	}
}

// ---------------------------------------------------------------------------
// Fuzz tests
// ---------------------------------------------------------------------------

// runIntOps decodes data as a sequence of 2-byte operations and checks d
// against a model after every step.
//
//	byte0 low 2 bits: 0,1 = put, 2 = delete, 3 = get
//	byte0 high 6 bits and byte1: the key
func runIntOps(t *testing.T, d *Dict[int, int], data []byte) {
	t.Helper()
	if len(data) > 4096 {
		data = data[:4096]
	}
	m := newModel[int, int]()
	for i := 0; i+1 < len(data); i += 2 {
		op := data[i] & 3
		k := int(data[i]>>2)<<8 | int(data[i+1])
		switch op {
		case 0, 1:
			d.Put(k, i)
			m.put(k, i)
		case 2:
			if got, want := d.Del(k), m.del(k); got != want {
				t.Fatalf("step %d: Del(%d) = %v, want %v", i/2, k, got, want)
			}
		case 3:
			want, wok := m.m[k]
			got, ok := d.Get2(k)
			if ok != wok || got != want {
				t.Fatalf("step %d: Get2(%d) = (%d,%v), want (%d,%v)", i/2, k, got, ok, want, wok)
			}
		}
		checkInvariants(t, d)
	}
	checkAgainst(t, d, m)
	for _, absent := range []int{-1, 1 << 20} {
		if _, ok := d.Get2(absent); ok {
			t.Fatalf("found never-inserted key %d", absent)
		}
	}
}

func addIntSeeds(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0, 1})
	f.Add([]byte{0, 1, 2, 1, 3, 1})
	f.Add([]byte{0, 1, 0, 2, 0, 3, 2, 2, 0, 2, 3, 2})
	// 40 puts then 40 deletes: crosses growth and shrink thresholds.
	var grow []byte
	for i := 0; i < 40; i++ {
		grow = append(grow, 0, byte(i))
	}
	for i := 0; i < 40; i++ {
		grow = append(grow, 2, byte(i))
	}
	f.Add(grow)
}

func FuzzDictOps(f *testing.F) {
	addIntSeeds(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		runIntOps(t, NewDict[int, int](), data)
	})
}

// Same ops but with a pathologically weak hash, so almost every key collides.
func FuzzDictOpsCollidingHash(f *testing.F) {
	addIntSeeds(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		d := NewDictFunc[int, int](func(k int) uint64 { return uint64(k % 3) })
		runIntOps(t, d, data)
	})
}

// Same ops with a constant hash: every key shares one probe sequence.
func FuzzDictOpsConstantHash(f *testing.F) {
	addIntSeeds(f)
	f.Fuzz(func(t *testing.T, data []byte) {
		d := NewDictFunc[int, int](func(int) uint64 { return 0 })
		if len(data) > 1024 {
			data = data[:1024] // constant hash is O(n) per op
		}
		runIntOps(t, d, data)
	})
}

// FuzzDictStringKeys decodes variable-length string keys from the input.
//
//	byte0 % 3: 0 = put, 1 = delete, 2 = get
//	byte1 % 5: key length; the following bytes are the key
func FuzzDictStringKeys(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0, 2, 'a', 'b', 2, 2, 'a', 'b'})
	f.Add([]byte{0, 0, 0, 0, 1, 0})
	f.Add([]byte{0, 3, 'x', 'y', 'z', 0, 3, 'x', 'y', 'z', 1, 3, 'x', 'y', 'z'})
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 4096 {
			data = data[:4096]
		}
		d := NewDict[string, int]()
		m := newModel[string, int]()
		step := 0
		for i := 0; i+1 < len(data); step++ {
			op := data[i] % 3
			n := int(data[i+1]) % 5
			i += 2
			end := min(i+n, len(data))
			k := string(data[i:end])
			i = end

			switch op {
			case 0:
				d.Put(k, step)
				m.put(k, step)
			case 1:
				if got, want := d.Del(k), m.del(k); got != want {
					t.Fatalf("step %d: Del(%q) = %v, want %v", step, k, got, want)
				}
			case 2:
				want, wok := m.m[k]
				got, ok := d.Get2(k)
				if ok != wok || got != want {
					t.Fatalf("step %d: Get2(%q) = (%d,%v), want (%d,%v)", step, k, got, ok, want, wok)
				}
			}
			checkInvariants(t, d)
		}
		checkAgainst(t, d, m)
	})
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkDictPut(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		d := NewDict[int, int]()
		for j := 0; j < 1000; j++ {
			d.Put(j, j)
		}
	}
}

func BenchmarkMapPut(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m := map[int]int{}
		for j := 0; j < 1000; j++ {
			m[j] = j
		}
	}
}

func BenchmarkDictGet(b *testing.B) {
	d := NewDict[int, int]()
	for j := 0; j < 1000; j++ {
		d.Put(j, j)
	}
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		sink += d.Get(i % 1000)
	}
	_ = sink
}

func BenchmarkMapGet(b *testing.B) {
	m := map[int]int{}
	for j := 0; j < 1000; j++ {
		m[j] = j
	}
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		sink += m[i%1000]
	}
	_ = sink
}

func BenchmarkDictIterate(b *testing.B) {
	d := NewDict[int, int]()
	for j := 0; j < 1000; j++ {
		d.Put(j, j)
	}
	b.ResetTimer()
	var sink int
	for i := 0; i < b.N; i++ {
		for _, v := range d.All() {
			sink += v
		}
	}
	_ = sink
}
