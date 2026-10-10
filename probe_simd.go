//go:build insdict_simd

package insdict

// Opt in to the grouped-probing experiment for the whole package.
const (
	groupedProbing = true
	minSize        = groupWidth
)
