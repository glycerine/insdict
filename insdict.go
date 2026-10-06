package insdict

import (
	"iter"
	"unsafe"

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

// Dict is an insertion-ordered hash map modeled on CPython's 3.7+ compact dict.
//
// For built-in comparable key types, the zero value Dict
// is perfectly usable and needs no NewDict() call. The built-in
// comparable key types are: string, int, int8, int16,
// int32, int64, uint, uint8, uint16, uint32, uint64, bool. Other
// key types need the user to supply the hash function, and so require
// a call to NewDictFunc to set up. The internal defaultHash() panics to enforce this.
//
// Pre-allocating with NewDictSize or NewDictFuncSize can save
// time and memory by avoiding table rebuilds on growth; benchmark your use.
//
// Just like the built-in Go map, we are not safe for concurrent use by default,
// and require external synchronization when a writer can race with readers.
// Readers do not modify the data structure and so do not race with each other.
// Any number of read-only goroutines can access a Dict concurrently (those
// that do no Put, no Del, and no Pack; only Get, Get2, Len, or All).
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
	d := &Dict[K, V]{}
	if hint > 0 {
		d.initSize(hint)
	}
	return d
}

// NewDictFuncSize is NewDictFunc with a capacity hint. A nil hash selects the
// default hash.
func NewDictFuncSize[K comparable, V any](hash func(K) uint64, hint int) *Dict[K, V] {
	d := &Dict[K, V]{hash: hash}
	if hint > 0 {
		d.initSize(hint)
	}
	return d
}

// Keep the constructor small enough to inline, so a local Dict can stay on
// the stack. Entry capacity follows the requested hint rather than spare slots.
func (d *Dict[K, V]) initSize(hint int) {
	d.rebuildEntries(presizeFor(hint), hint)
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

// EasyHashString provides a default hash for strings. Currently this
// is based on cespare/xxhash, but this is subject to change.
func EasyHashString(key string) uint64 { return xxhash.Sum64String(key) }

// EasyHash* functions are a set of convenience hash functions
// to make it easy for users to call NewDictFunc
// for a given key type. It is probably faster to set the d.hash function once
// rather than to dispatch to defaultHash on every Get and doing
// reflection. After all, the type of the key is known and fixed.
// .
func EasyHashBool(key bool) uint64 {
	if key {
		return Mix64(1)
	}
	return Mix64(0)
}

func EasyHashInt(key int) uint64     { return Mix64(uint64(key)) }
func EasyHashInt8(key int8) uint64   { return Mix64(uint64(key)) }
func EasyHashInt16(key int16) uint64 { return Mix64(uint64(key)) }
func EasyHashInt32(key int32) uint64 { return Mix64(uint64(key)) }
func EasyHashInt64(key int64) uint64 { return Mix64(uint64(key)) }

func EasyHashUint(key uint) uint64     { return Mix64(uint64(key)) }
func EasyHashUint8(key uint8) uint64   { return Mix64(uint64(key)) }
func EasyHashUint16(key uint16) uint64 { return Mix64(uint64(key)) }
func EasyHashUint32(key uint32) uint64 { return Mix64(uint64(key)) }
func EasyHashUint64(key uint64) uint64 { return Mix64(key) }

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

// usable is the maximum entry count (holes included) before a table rebuild.
func (d *Dict[K, V]) usable() int { return len(d.indices) * 2 / 3 }

// first empty-or-dummy slot on h's probe path (only call when the key is known absent)
func (d *Dict[K, V]) freeSlot(h uint64) uint64 {
	i, perturb := h&d.mask, h
	for d.tags[i] >= tagUsed {
		perturb >>= 5
		i = (i*5 + perturb + 1) & d.mask
	}
	return i
}

// Rebuild the table to size slots, compacting holes and preserving order.
// Same-size compaction reuses existing storage when capacity permits.
func (d *Dict[K, V]) rebuild(size int) {
	d.rebuildEntries(size, size*2/3)
}

func (d *Dict[K, V]) rebuildEntries(size, capacity int) {
	old := d.entries
	reuse := size == len(d.indices)
	if reuse {
		clear(d.tags)
	} else {
		// Sizes are powers of two >= 8. Allocate the int64 indexes followed by
		// exactly size control bytes in one pointer-free, aligned allocation.
		// Limit the index slice's capacity so it cannot overlap the controls.
		storage := make([]int64, size+size/8)
		d.indices = storage[:size:size]
		d.tags = unsafe.Slice((*byte)(unsafe.Pointer(&storage[size])), size)
	}
	for i := range d.indices {
		d.indices[i] = slotEmpty
	}
	d.mask = uint64(size - 1)
	reuseEntries := reuse && cap(old) >= capacity
	if reuseEntries {
		d.entries = old[:0]
	} else {
		d.entries = make([]entry[K, V], 0, capacity)
	}
	for i := range old {
		if !old[i].live {
			continue
		}
		slot := d.freeSlot(old[i].hash)
		d.indices[slot] = int64(len(d.entries))
		d.tags[slot] = hashTag(old[i].hash)
		d.entries = append(d.entries, old[i])
	}
	if reuseEntries {
		// Compaction can leave duplicate pointer-bearing entries past the new
		// length. Clear them so later deletions can release their keys/values.
		clear(old[len(d.entries):])
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
// Put may rebuild and re-pack the underlying array. Interleaving Put with
// range All() iteration is not recommended, as it may make iteration miss elements. See
// the All docs for more information.
func (d *Dict[K, V]) Put(k K, v V) {
	if d.indices == nil {
		d.rebuild(minSize)
	}
	var h uint64
	if d.hash != nil {
		h = d.hash(k)
	} else if v, ok := any(k).(int); ok {
		h = Mix64(uint64(v))
	} else {
		h = defaultHash(k)
	}
	tag := hashTag(h)

	// existing key: overwrite in place, order position unchanged
	slot, ix := d.find(k, h, tag)
	if ix >= 0 {
		d.entries[ix].val = v
		return
	}

	// New key: rebuild when the table's entry budget is full (holes included).
	if len(d.entries) >= d.usable() {
		d.rebuild(sizeFor(d.live))
		slot = int(d.freeSlot(h))
	} else if d.live != len(d.entries) {
		// Deletions leave tombstones: prefer the first one on the probe path.
		// Without holes, find already returned the insertion slot.
		slot = int(d.freeSlot(h))
	}

	ix = int64(len(d.entries))
	d.entries = append(d.entries, entry[K, V]{hash: h, key: k, val: v, live: true})
	d.indices[slot] = ix
	d.tags[slot] = tag
	d.live++
}

// Del removes k and reports whether it was present.
//
// Del will not automatically re-pack the underlying table, even
// if many tombstones are present, and thus it is safe to delete
// with Del during an All iteration.
//
// After many deletions, to vacuum tombstones, you should call Pack
// manually, or use juse regularly use DelPackFalse instead of Del.
//
// DelPackFalse is a convenient alternative to Del
// that will automatically compact based on heuristics -- if you
// don't want to think about this very hard about when to Pack but still want
// your table to get Packed at some point.
func (d *Dict[K, V]) Del(k K) (found bool) {
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

	return true
}

// DelPackFalse is a convenience wrapper that calls
// Del(k) followed by Pack(force=false).
// As a replacement for Del, it can save the user from having
// to think too hard about when to Pack away their tombstones.
// However it cannot be intermixed with All iteration safely
// as it calls Pack; see the comments on All.
func (d *Dict[K, V]) DelPackFalse(k K) (found bool) {
	found = d.Del(k)
	d.Pack(false)
	return
}

// Pack may vacuum and re-pack the underlying array, removing tombstones.
// If force is false then heuristics are used, currently 75% tombstones,
// to decide whether to re-pack. If force is true then we always repack
// if there is a single tombstone. If there are no tombstones then
// Pack is always a very fast no-op. This enables preparing for
// Put during iteration (an uncommon pattern):
//
// You must call Pack(true) to eliminate all tombstones before doing a
// range All() in the special case of interleaving Put calls with iteration -- otherwise
// your iteration may miss keys after a Put grows the table and
// shrinks the indexes of keys that had tombstones before them.
// Do not do both Put and Del during All iteration unless you can
// tolerate skipping over some keys unknowingly. See the All docs for more.
func (d *Dict[K, V]) Pack(force bool) {
	if d.live == len(d.entries) {
		// no tombstones, do nothing.
		return
	}
	if !force && len(d.entries) > 32 && d.live < len(d.entries)/4 {
		force = true
	}
	if force {
		d.rebuild(sizeFor(d.live))
	}
}

// All iterates entries in insertion order.
//
// It is safe to call Del() during iteration
// since it does not auto-repack the array, but
// instead only writes a tombstone. (Compaction only happens
// when the user calls Pack() manually or when Put
// grows the array and we compact during the copy over).
//
// Put of new keys during All iteration is not recommended.
// Put could provoke a resize of the underlying array.
// The copy to the new larger array will omit tombstones. This will
// lower the index of elements that were after the tombstones. You risk missing some
// elements during the iteration (without knowing it), if there were
// tombstones present before the current iteration point.
//
// If you really must Put during iteration,
// be sure to call Pack(true) before starting All so as to force vacuuming out of all
// tombstones beforehand; and forbid Del during such iterations (that also Put).
// As noted, a Del followed by a Put can result in an internal copy
// and compaction that will lower the index of all keys that had
// tombstones before them in the array. The iterator's held index integer can
// become too large, causing some Dict entries to be missed. Since
// this is not expected to be a common use pattern, we do not contort the code to
// accommodate it. You have been warned.
//
// A simple alternative approach that will not skip over any of the original
// keys while supporting both Put and Del during iteration is to
// Clone the Dict and iterate one copy while modifying the other.
//
// Note that if you only need to Put (and not Del), then Pack(true)
// once before All suffices to avoid skipped keys and the need to Clone.
// Lacking tombstones, the underlying array can be grown during
// iteration without changing any of the original index positions.
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

// Clone creates and returns an independent and identical copy of d.
// Of course if V contains pointers then r will contain
// an identical copy of those pointers.
func (d *Dict[K, V]) Clone() (r *Dict[K, V]) {
	r = &Dict[K, V]{
		hash:    d.hash,
		indices: append([]int64(nil), d.indices...),
		entries: append([]entry[K, V](nil), d.entries...),
		live:    d.live,
		mask:    d.mask,
		tags:    append([]byte(nil), d.tags...),
	}
	return
}
