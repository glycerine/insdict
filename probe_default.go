//go:build !insdict_simd

package insdict

const (
	groupedProbing = false
	minSize        = 8
)
