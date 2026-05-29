package gengo

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"testing"
	"time"
)

// fingerprint exercises every Generator method in a fixed order and returns a
// combined string capturing the full produced sequence, so two fingerprints are
// equal iff the two generators produced identical output.
func fingerprint(g *Generator) string {
	var b strings.Builder
	start := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 100; i++ {
		fmt.Fprintf(&b, "%d|%d|%d|%d|%d|%d|%d|%d|%d|%d|",
			g.Int8(), g.Int8Between(-100, 100), g.Int16(), g.Int16Between(-1000, 1000),
			g.Int32(), g.Int32Between(-1000000, 1000000), g.Int(), g.IntBetween(1, 100),
			g.Int64(), g.Int64Between(0, 1000000000))
		fmt.Fprintf(&b, "%d|%d|%d|%d|%d|%d|%d|%d|",
			g.Uint8(), g.Uint8Between(0, 200), g.Byte(), g.Uint16(),
			g.Uint16Between(0, 5000), g.Uint32(), g.Uint32Between(0, 1000000), g.Uint64())
		fmt.Fprintf(&b, "%d|%v|%v|%v|%v|",
			g.Uint64Between(0, 1000000000000), g.Float32(), g.Float32Between(-1, 1),
			g.Float64(), g.Float64Between(-1, 1))
		fmt.Fprintf(&b, "%v|%v|%v|%v|",
			g.Complex64(), g.Complex64Between(-1, 1, -1, 1), g.Complex128(),
			g.Complex128Between(-1, 1, -1, 1))
		fmt.Fprintf(&b, "%s|%s|%s|%s|%s|%s|%s|%s|%s|%s|",
			g.String(8, Alphanumeric), g.StringBetween(3, 12, Alphabetic), g.StringAllChars(6),
			g.StringAlphanumeric(5), g.StringAlphabetic(5), g.StringAlphabeticUppercase(5),
			g.StringAlphabeticLowercase(5), g.StringNumeric(5), g.StringHexadecimal(5),
			g.StringSymbols(5))
		fmt.Fprintf(&b, "%d|%d|%d|%v|%s|%s|%v|",
			g.Date().UnixNano(), g.UnixDate().UnixNano(), g.DateBetween(start, end).UnixNano(),
			g.Bool(), g.Word(), g.WordByLengthType(MediumLengthWords), g.Words(3))
	}
	return b.String()
}

func TestGeneratorReproducible(t *testing.T) {
	// Same seed must yield identical sequences across every method.
	first := fingerprint(New(99))
	second := fingerprint(New(99))
	if first != second {
		t.Fatal("New(99) must produce identical sequences for the same seed")
	}
	// Different seeds must (with overwhelming probability) diverge.
	different := fingerprint(New(100))
	if first == different {
		t.Fatal("different seeds should produce different sequences")
	}
}

func TestNewSourceReproducible(t *testing.T) {
	a := NewSource(rand.NewPCG(7, 11))
	b := NewSource(rand.NewPCG(7, 11))
	for i := 0; i < loop; i++ {
		if a.Uint64() != b.Uint64() {
			t.Fatal("NewSource with identical PCG seeds must reproduce the sequence")
		}
	}
}

func TestGeneratorRanges(t *testing.T) {
	g := New(1)
	for i := 0; i < loop; i++ {
		if v := g.IntBetween(10, 20); v < 10 || v > 20 {
			t.Fatalf("IntBetween out of range: %d", v)
		}
		if v := g.Float64Between(-5, 5); v < -5 || v > 5 {
			t.Fatalf("Float64Between out of range: %v", v)
		}
		if v := g.Uint8Between(100, 150); v < 100 || v > 150 {
			t.Fatalf("Uint8Between out of range: %d", v)
		}
		s := g.String(16, Alphanumeric)
		if len(s) != 16 {
			t.Fatalf("String length = %d, want 16", len(s))
		}
		for _, ch := range s {
			if !strings.ContainsRune(Alphanumeric, ch) {
				t.Fatalf("String char %q not in Alphanumeric", ch)
			}
		}
	}
}

func TestGeneratorIntSwapBranches(t *testing.T) {
	// Signed integer Between methods must swap reversed bounds (min > max).
	g := New(7)
	if v := g.Int8Between(100, -100); v < -100 || v > 100 {
		t.Fatalf("Int8Between swap: %d", v)
	}
	if v := g.Int16Between(1000, -1000); v < -1000 || v > 1000 {
		t.Fatalf("Int16Between swap: %d", v)
	}
	if v := g.Int32Between(1000, -1000); v < -1000 || v > 1000 {
		t.Fatalf("Int32Between swap: %d", v)
	}
	if v := g.IntBetween(100, 1); v < 1 || v > 100 {
		t.Fatalf("IntBetween swap: %d", v)
	}
	if v := g.Int64Between(1000, 0); v < 0 || v > 1000 {
		t.Fatalf("Int64Between swap: %d", v)
	}
}

func TestGeneratorUintSwapAndFullRange(t *testing.T) {
	g := New(7)
	// Unsigned Between methods must swap reversed bounds.
	if v := g.Uint8Between(200, 0); v > 200 {
		t.Fatalf("Uint8Between swap: %d", v)
	}
	if v := g.Uint16Between(5000, 0); v > 5000 {
		t.Fatalf("Uint16Between swap: %d", v)
	}
	if v := g.Uint32Between(1000000, 0); v > 1000000 {
		t.Fatalf("Uint32Between swap: %d", v)
	}
	if v := g.Uint64Between(1000, 0); v > 1000 {
		t.Fatalf("Uint64Between swap: %d", v)
	}
	// Full-range special branches (rangeSize overflows to 0 / exceeds MaxUint32).
	_ = g.IntBetween(math.MinInt, math.MaxInt)
	_ = g.Int64Between(math.MinInt64, math.MaxInt64)
	_ = g.Uint32Between(0, math.MaxUint32)
	_ = g.Uint64Between(0, math.MaxUint64)
}

func TestGeneratorFloatBranches(t *testing.T) {
	g := New(7)
	// Swap branches.
	if v := g.Float32Between(1, -1); v < -1 || v > 1 {
		t.Fatalf("Float32Between swap: %v", v)
	}
	if v := g.Float64Between(1, -1); v < -1 || v > 1 {
		t.Fatalf("Float64Between swap: %v", v)
	}
	// Wide-range (overflow) branches must stay finite and in range.
	for i := 0; i < loop; i++ {
		if v := g.Float32Between(-math.MaxFloat32, math.MaxFloat32); math.IsInf(float64(v), 0) || math.IsNaN(float64(v)) {
			t.Fatalf("Float32Between wide produced non-finite: %v", v)
		}
		if v := g.Float64Between(-math.MaxFloat64, math.MaxFloat64); math.IsInf(v, 0) || math.IsNaN(v) {
			t.Fatalf("Float64Between wide produced non-finite: %v", v)
		}
	}
}

func TestGeneratorStringEdges(t *testing.T) {
	g := New(7)
	if s := g.String(0, "abc"); s != "" {
		t.Fatalf("String(0) = %q", s)
	}
	if s := g.String(10, ""); s != "" {
		t.Fatalf("String empty source = %q", s)
	}
	if s := g.String(20, "a"); s != strings.Repeat("a", 20) {
		t.Fatalf("String single-char = %q", s)
	}
	if s := g.String(MaxStringLength+1, "ab"); uint32(len(s)) != MaxStringLength {
		t.Fatalf("String cap = %d", len(s))
	}
	if s := g.StringBetween(12, 3, Alphabetic); len(s) < 3 || len(s) > 12 {
		t.Fatalf("StringBetween swap len = %d", len(s))
	}
}

func TestGeneratorDateAndWordEdges(t *testing.T) {
	g := New(7)
	// DateBetween swap + same-second branches.
	start := time.Date(2001, 2, 3, 4, 5, 6, 0, time.UTC)
	end := time.Date(2009, 8, 7, 6, 5, 4, 0, time.UTC)
	if d := g.DateBetween(end, start); d.Before(start) || d.After(end) {
		t.Fatalf("DateBetween swap: %v", d)
	}
	if d := g.DateBetween(start, start); !d.Equal(start) {
		t.Fatalf("DateBetween same-second = %v, want %v", d, start)
	}
	// WordByLengthType all categories + default.
	if w := g.WordByLengthType(SmallLengthWord); len(w) < 1 || len(w) > 4 {
		t.Fatalf("Small word len = %d", len(w))
	}
	if w := g.WordByLengthType(BigLengthWords); len(w) < 9 || len(w) > 30 {
		t.Fatalf("Big word len = %d", len(w))
	}
	if w := g.WordByLengthType(LengthTypeWords(99)); len(w) < 1 || len(w) > 30 {
		t.Fatalf("default word len = %d", len(w))
	}
	// Words non-positive lengths return an empty slice.
	if w := g.Words(0); len(w) != 0 {
		t.Fatalf("Words(0) len = %d", len(w))
	}
	if w := g.Words(-3); len(w) != 0 {
		t.Fatalf("Words(-3) len = %d", len(w))
	}
}

func BenchmarkGeneratorInt64(b *testing.B) {
	g := New(1)
	for i := 0; i < b.N; i++ {
		_ = g.Int64()
	}
}

func BenchmarkGeneratorString(b *testing.B) {
	g := New(1)
	for i := 0; i < b.N; i++ {
		_ = g.String(8, Alphanumeric)
	}
}
