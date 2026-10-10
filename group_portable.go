//go:build !insdict_simd || !amd64 || purego || insdict_swar

package insdict

func groupMasks(ctrl *[groupWidth]byte, tag byte) uint64 {
	return groupMasksPortable(ctrl, tag)
}
