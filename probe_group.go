package insdict

import (
	"encoding/binary"
	"math/bits"
)

const groupWidth = 16

// Groups are aligned within the control array, so a 16-byte load never crosses
// its end, including for the smallest table. Triangular probing visits every
// group in a power-of-two table: g, g+1, g+3, g+6, ... (modulo group count).
// Empty bytes terminate lookup only AFTER checking every candidate in a group.
// Deletions remain tombstones, so they cannot break reachability of later groups.
func (d *Dict[K, V]) findGrouped(k K, h uint64, tag byte) (slot, ix int) {
	base := (h & d.mask) &^ uint64(groupWidth-1)
	var step uint64
	for {
		masks := groupMasks((*[groupWidth]byte)(d.tags[base:]), tag)
		for matches := uint16(masks); matches != 0; matches &= matches - 1 {
			slot := int(base) + bits.TrailingZeros16(matches)
			ix := d.indices[slot]
			e := &d.entries[ix]
			if e.hash == h && (e.key == k || (e.key != e.key && k != k && equalNaNKeys(e.key, k))) {
				return slot, ix
			}
		}
		if empty := uint16(masks >> 16); empty != 0 {
			return int(base) + bits.TrailingZeros16(empty), -1
		}
		step += groupWidth
		base = (base + step) & d.mask
	}
}

// Only called after find has established that the key is absent. Returning an
// earlier tombstone is then safe even if the key's probe path spans many groups.
func (d *Dict[K, V]) freeSlotGrouped(h uint64) uint64 {
	base := (h & d.mask) &^ uint64(groupWidth-1)
	var step uint64
	for {
		masks := groupMasks((*[groupWidth]byte)(d.tags[base:]), 0)
		if available := uint16(masks >> 32); available != 0 {
			return base + uint64(bits.TrailingZeros16(available))
		}
		step += groupWidth
		base = (base + step) & d.mask
	}
}

// Each 16-bit field contains one bit per slot: fingerprint matches (low),
// empty slots (middle), and empty-or-deleted slots (high). Bits 48..63 are zero.
// The portable implementation uses SWAR (SIMD within a register), with explicit
// little-endian loads to keep lane order identical on every architecture.
func groupMasksPortable(ctrl *[groupWidth]byte, tag byte) uint64 {
	a := binary.LittleEndian.Uint64(ctrl[:8])
	b := binary.LittleEndian.Uint64(ctrl[8:])
	want := uint64(tag) * 0x0101010101010101
	matches := byteMask(zeroBytes(a^want)) | byteMask(zeroBytes(b^want))<<8
	empty := byteMask(zeroBytes(a)) | byteMask(zeroBytes(b))<<8
	available := byteMask(^a) | byteMask(^b)<<8
	return matches | empty<<16 | available<<32
}

// Unlike (x-ones)&^x&highBits, this exact zero-byte test has no borrow-induced
// false positives in neighboring lanes (notably an empty byte followed by 1).
func zeroBytes(x uint64) uint64 {
	const lowBits = uint64(0x7f7f7f7f7f7f7f7f)
	return ^(((x & lowBits) + lowBits) | x | lowBits)
}

func byteMask(x uint64) uint64 {
	return (((x & 0x8080808080808080) >> 7) * 0x0102040810204080) >> 56
}
