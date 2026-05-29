// Package gengo provides a simple and practical way to generate random data
// structures. It offers a zero-dependency API for generating random values of
// standard Go types including integers, floats, complex numbers, strings,
// booleans, dates, and words.
//
// The package-level functions draw from the global math/rand/v2 source and are
// safe for concurrent use. For reproducible output, use a seeded [Generator]
// created with [New] or [NewSource].
//
// gengo is not cryptographically secure: it is built on math/rand/v2, whose
// output is predictable. Do not use it for passwords, tokens, keys, or any other
// secret; use crypto/rand instead.
package gengo
