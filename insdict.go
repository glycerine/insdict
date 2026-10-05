package insdict

import (
	"iter"

	"github.com/cespare/xxhash/v2"
)

const (
	slotEmpty int64 = -1
	slotDummy int64 = -2
	minSize         = 8 // power of two
	tagEmpty  byte  = 0
	tagDummy  byte  = 1
	tagUsed   byte  = 0x80
)

type entry[K comparable, V any] struct {
	hash uint64
	key  K
	val  V
	live bool
}

// Dict is an insertion-ordered hash map modeled on CPython's compact dict.
// The zero value is not usable for hashing custom key types; use NewDict or
// NewDictFunc. Not safe for concurrent use.
type Dict[K comparable, V any] struct {
	hash    func(K) uint64 // nil => defaultHash
	indices []int64        // slotEmpty, slotDummy, or index into entries
	entries []entry[K, V]  // dense, insertion-ordered, may contain holes
	live    int            // live entry count
	mask    uint64
	tags    []byte // empty, dummy, or tagUsed | high 7 hash bits
}

func NewDict[K comparable, V any]() *Dict[K, V] {
	return &Dict[K, V]{}
}

// NewDictFunc lets callers supply a hash for key types the default doesn't know.
func NewDictFunc[K comparable, V any](hash func(K) uint64) *Dict[K, V] {
	return &Dict[K, V]{hash: hash}
}

// NewDictSize returns a Dict with room for hint entries, so inserting up to
// hint distinct keys triggers no rebuild and no reallocation.
func NewDictSize[K comparable, V any](hint int) *Dict[K, V] {
	return NewDictFuncSize[K, V](nil, hint)
}

// NewDictFuncSize is NewDictFunc with a capacity hint. A nil hash selects the
// default hash.
func NewDictFuncSize[K comparable, V any](hash func(K) uint64, hint int) *Dict[K, V] {
	d := &Dict[K, V]{hash: hash}
	if hint > 0 {
		d.rebuild(presizeFor(hint))
	}
	return d
}

// presizeFor returns the smallest power-of-two table size whose usable entry
// capacity (2/3 of the table) holds n entries.
func presizeFor(n int) int {
	size := minSize
	for size*2/3 < n {
		size <<= 1
	}
	return size
}

func (d *Dict[K, V]) hashOf(k K) uint64 {
	if d.hash != nil {
		return d.hash(k)
	}
	return defaultHash(k)
}

// Mix64 is the splitmix64 finalizer: a good, cheap integer mixer.
// Exported so custom hash functions can use it to combine fields.
func Mix64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

func defaultHash[K comparable](k K) uint64 {
	switch v := any(k).(type) {
	case string:
		return xxhash.Sum64String(v)
	case int:
		return Mix64(uint64(v))
	case int8:
		return Mix64(uint64(v))
	case int16:
		return Mix64(uint64(v))
	case int32:
		return Mix64(uint64(v))
	case int64:
		return Mix64(uint64(v))
	case uint:
		return Mix64(uint64(v))
	case uint8:
		return Mix64(uint64(v))
	case uint16:
		return Mix64(uint64(v))
	case uint32:
		return Mix64(uint64(v))
	case uint64:
		return Mix64(v)
	case bool:
		if v {
			return Mix64(1)
		}
		return Mix64(0)
	}
	panic("insdict: no default hash for key type; use NewDictFunc")
}

// Len returns the number of live entries.
func (d *Dict[K, V]) Len() int { return d.live }

// Get returns the value for k, or the zero value if absent.
func (d *Dict[K, V]) Get(k K) (v V) {
	v, _ = d.Get2(k)
	return
}

// Get2 returns the value for k and whether it was present.
func (d *Dict[K, V]) Get2(k K) (v V, found bool) {
	if d.live == 0 {
		return
	}
	var h uint64
	if d.hash != nil {
		h = d.hash(k)
	} else if v, ok := any(k).(int); ok {
		// Keep the common integer hash inline instead of dispatching through
		// hashOf and the full defaultHash type switch on every lookup.
		h = Mix64(uint64(v))
	} else {
		h = defaultHash(k)
	}
	tag := hashTag(h)
	i, perturb := h&d.mask, h
	for {
		// Check compact metadata first: most misses never read either array.
		ctrl := d.tags[i]
		if ctrl == tag {
			ix := d.indices[i]
			e := &d.entries[ix]
			if e.hash == h && e.key == k {
				return e.val, true
			}
		} else if ctrl == tagEmpty {
			return
		}
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
}

// Occupied tags use the high bit to distinguish them from empty and dummy
// slots. Use hash bits independent of the low bits selecting the initial slot.
func hashTag(h uint64) byte { return tagUsed | byte(h>>57) }

// find returns the indices slot and entries index for k, or ix == -1 and the
// empty slot where the probe ended if k is absent. Requires d.indices != nil
// and tag == hashTag(h). Passing the tag keeps this probe loop inlineable.
func (d *Dict[K, V]) find(k K, h uint64, tag byte) (slot int, ix int64) {
	i, perturb := h&d.mask, h
	for d.tags[i] != tagEmpty {
		if d.tags[i] == tag {
			ix = d.indices[i]
			e := &d.entries[ix]
			if e.hash == h && e.key == k {
				return int(i), ix
			}
		}
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
	return int(i), -1
}

// entries capacity is 2/3 of the index table, counting holes
func (d *Dict[K, V]) usable() int { return len(d.indices) * 2 / 3 }

// first empty-or-dummy slot on h's probe path (only call when the key is known absent)
func (d *Dict[K, V]) freeSlot(h uint64) uint64 {
	i, perturb := h&d.mask, h
	for d.indices[i] >= 0 {
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
	return i
}

// rebuild into a fresh table of `size` slots,
// compacting out holes, preserving order
func (d *Dict[K, V]) rebuild(size int) {
	old := d.entries
	d.indices = make([]int64, size)
	d.tags = make([]byte, size)
	for i := range d.indices {
		d.indices[i] = slotEmpty
	}
	d.mask = uint64(size - 1)
	d.entries = make([]entry[K, V], 0, d.usable())
	for i := range old {
		if !old[i].live {
			continue
		}
		slot := d.freeSlot(old[i].hash)
		d.indices[slot] = int64(len(d.entries))
		d.tags[slot] = hashTag(old[i].hash)
		d.entries = append(d.entries, old[i])
	}
}

// size such that live entries fill at most ~1/3 of usable capacity after rebuild
func sizeFor(live int) int {
	size := minSize
	for size < live*3 {
		size <<= 1
	}
	return size
}

// Put sets k to v. Overwriting an existing key keeps its original position.
func (d *Dict[K, V]) Put(k K, v V) {
	if d.indices == nil {
		d.rebuild(minSize)
	}
	h := d.hashOf(k)
	tag := hashTag(h)

	// existing key: overwrite in place, order position unchanged
	if _, ix := d.find(k, h, tag); ix >= 0 {
		d.entries[ix].val = v
		return
	}

	// new key: make room if the entries array is full (holes included)
	if len(d.entries) >= d.usable() {
		d.rebuild(sizeFor(d.live))
	}

	ix := int64(len(d.entries))
	d.entries = append(d.entries, entry[K, V]{hash: h, key: k, val: v, live: true})
	slot := d.freeSlot(h)
	d.indices[slot] = ix
	d.tags[slot] = tag
	d.live++
}

// Del removes k and reports whether it was present.
// Del may compact the table, which invalidates in-progress iteration.
func (d *Dict[K, V]) Del(k K) bool {
	if d.live == 0 {
		return false
	}
	h := d.hashOf(k)

	slot, ix := d.find(k, h, hashTag(h))
	if ix < 0 {
		return false
	}

	d.indices[slot] = slotDummy
	d.tags[slot] = tagDummy
	d.entries[ix] = entry[K, V]{} // zero it so the GC can release K and V
	d.live--

	if len(d.entries) > 32 && d.live < len(d.entries)/4 {
		d.rebuild(sizeFor(d.live))
	}
	return true
}

// All iterates entries in insertion order. Do not call Del during iteration
// (it may compact and renumber entries). Put of new keys during iteration
// is also not safe since it too could also provoke a resize.
func (d *Dict[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i := 0; i < len(d.entries); i++ {
			e := &d.entries[i]
			if !e.live {
				continue
			}
			if !yield(e.key, e.val) {
				return
			}
		}
	}
}
