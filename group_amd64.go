//go:build insdict_simd && amd64 && !purego && !insdict_swar

package insdict

// SSE2 is part of the amd64 baseline on both Darwin and Linux. No CPU feature
// dispatch, cgo, GOEXPERIMENT, or AVX requirement is needed.
//
//go:noescape
func groupMasks(ctrl *[groupWidth]byte, tag byte) uint64
