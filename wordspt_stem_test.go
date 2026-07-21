package gengo

import (
	"math/rand/v2"
	"strconv"
	"testing"
	"unicode"
	"unicode/utf8"
)

// The tests in this file exercise the WordsPT Layer 1 stem sampler and assembler
// (wordspt_stem.go). They prove, over a large sampled run, that every generated
// stem is correct by construction against the phonotactic oracle in wordspt.go,
// that requested syllable counts and tonic indices are honored, that the output
// is always valid lowercase UTF-8, and that the derived sampling inventories keep
// the correct-by-construction (no-rejection) property.

// newTestRand returns a deterministic source for reproducible sampling in tests.
func newTestRand() *rand.Rand { return rand.New(rand.NewPCG(0xC0FFEE, 0xBEEF)) }

// TestSampleSyllabicStemConformsToPhonotactics is the acceptance-criterion test:
// over at least 1000 iterations (the existing loop convention), across a range of
// syllable counts and tonic positions, every generated stem passes the
// phonotactic validator, the requested syllable count and tonic are honored, and
// the assembled string is valid lowercase UTF-8 that also satisfies the
// string-level orthographic oracles (cedilla, diaeresis, grave).
func TestSampleSyllabicStemConformsToPhonotactics(t *testing.T) {
	const loop = 1000
	r := newTestRand()

	var dst []syllable
	for i := 0; i < loop; i++ {
		for count := 1; count <= 8; count++ {
			// Exercise every tonic position, including out-of-range values,
			// to confirm the clamp always yields a valid tonic index.
			tonic := i % (count + 2)

			st := sampleSyllabicStem(r, dst, count, tonic)
			dst = st.syllables // reuse the buffer across iterations

			if len(st.syllables) != count {
				t.Fatalf("syllable count not honored: got %d, want %d", len(st.syllables), count)
			}
			if st.tonic < 0 || st.tonic >= count {
				t.Fatalf("tonic index %d out of range [0,%d)", st.tonic, count)
			}
			if !stemConformsToPhonotactics(st.syllables) {
				t.Fatalf("stem does not conform to phonotactics: %+v", st.syllables)
			}

			word := assembleStem(st.syllables)
			if !utf8.ValidString(word) {
				t.Fatalf("assembled stem is not valid UTF-8: %q", word)
			}
			if word == "" {
				t.Fatalf("assembled stem is empty for count %d", count)
			}
			if hasUpper(word) {
				t.Fatalf("assembled stem contains an uppercase letter: %q", word)
			}
			// The bare stem must also satisfy the string-level orthographic
			// oracles: cedilla placement, no abolished diaeresis, no grave accent.
			if !wordConformsToCedilla(word) {
				t.Fatalf("assembled stem violates the cedilla rule: %q", word)
			}
			if wordHasForbiddenDiaeresis(word) {
				t.Fatalf("assembled stem contains a forbidden diaeresis: %q", word)
			}
			if wordContainsGrave(word) {
				t.Fatalf("assembled stem contains a grave accent: %q", word)
			}
		}
	}
}

// hasUpper reports whether s contains any uppercase letter. WordsPT output is
// always lowercase.
func hasUpper(s string) bool {
	for _, r := range s {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// TestSampleSyllabicStemHonoursTonic checks that an in-range tonic is carried
// through unchanged and that out-of-range tonics are clamped into range.
func TestSampleSyllabicStemHonoursTonic(t *testing.T) {
	r := newTestRand()
	tests := []struct {
		count, tonic, want int
	}{
		{1, 0, 0},
		{3, 0, 0},
		{3, 1, 1},
		{3, 2, 2},
		{3, 5, 2},  // above range -> last syllable
		{3, -4, 0}, // below range -> first syllable
		{0, 0, 0},  // count clamped to 1, tonic clamped to 0
		{-2, 3, 0}, // count clamped to 1, tonic clamped to 0
	}
	for _, tc := range tests {
		st := sampleSyllabicStem(r, nil, tc.count, tc.tonic)
		if st.tonic != tc.want {
			t.Errorf("sampleSyllabicStem(count=%d, tonic=%d).tonic = %d, want %d",
				tc.count, tc.tonic, st.tonic, tc.want)
		}
		wantCount := tc.count
		if wantCount < 1 {
			wantCount = 1
		}
		if len(st.syllables) != wantCount {
			t.Errorf("sampleSyllabicStem(count=%d) produced %d syllables, want %d",
				tc.count, len(st.syllables), wantCount)
		}
	}
}

// TestSampleSyllabicStemReproducible confirms that two samplers seeded
// identically produce byte-for-byte identical stems, so the *Generator surface
// (which injects its own seeded source) is reproducible.
func TestSampleSyllabicStemReproducible(t *testing.T) {
	const loop = 1000
	a := New(20260721) // *Generator with a seeded source
	b := New(20260721)

	var da, db []syllable
	for i := 0; i < loop; i++ {
		count := 1 + i%6
		tonic := i % count

		sa := sampleSyllabicStem(a.r, da, count, tonic) // injects the Generator's source
		sb := sampleSyllabicStem(b.r, db, count, tonic)
		da, db = sa.syllables, sb.syllables

		wa, wb := assembleStem(sa.syllables), assembleStem(sb.syllables)
		if wa != wb {
			t.Fatalf("identical seeds diverged at i=%d: %q != %q", i, wa, wb)
		}
		if sa.tonic != sb.tonic {
			t.Fatalf("identical seeds diverged on tonic at i=%d: %d != %d", i, sa.tonic, sb.tonic)
		}
	}
}

// TestSampleSyllabicStemBufferReuse verifies that reusing the returned buffer
// does not corrupt later stems: every field is freshly written on each call, so
// a smaller request after a larger one carries no stale data.
func TestSampleSyllabicStemBufferReuse(t *testing.T) {
	r := newTestRand()
	var dst []syllable

	st := sampleSyllabicStem(r, dst, 8, 3)
	dst = st.syllables
	if !stemConformsToPhonotactics(st.syllables) {
		t.Fatalf("8-syllable stem not conformant: %+v", st.syllables)
	}

	// A shorter request reusing the same backing array must not leak the tail.
	for count := 1; count <= 8; count++ {
		st = sampleSyllabicStem(r, dst, count, 0)
		dst = st.syllables
		if len(st.syllables) != count {
			t.Fatalf("reuse: count not honored: got %d, want %d", len(st.syllables), count)
		}
		if !stemConformsToPhonotactics(st.syllables) {
			t.Fatalf("reuse: stem not conformant at count %d: %+v", count, st.syllables)
		}
		if cap(dst) < 8 {
			t.Fatalf("reuse: backing capacity shrank to %d", cap(dst))
		}
	}
}

// TestDerivedInventoriesNoRejection verifies that the derived sampling
// inventories share the correct-by-construction property of the base
// inventories: for every r in [0,total) the returned index partitions the
// cumulative table correctly, and every form is reachable. No draw is rejected.
func TestDerivedInventoriesNoRejection(t *testing.T) {
	inventories := []struct {
		name string
		inv  *weightedInventory
	}{
		{"onsetInitialInv", &onsetInitialInv},
		{"onsetMedialInv", &onsetMedialInv},
		{"nucleiFrontInv", &nucleiFrontInv},
		{"nucleiBackInv", &nucleiBackInv},
		{"codaBeforePBInv", &codaBeforePBInv},
		{"codaBeforeOtherInv", &codaBeforeOtherInv},
		{"codaFinalInv", &codaFinalInv},
	}
	for _, it := range inventories {
		t.Run(it.name, func(t *testing.T) {
			verifyInventorySampling(t, it.name, it.inv)
		})
	}
}

// TestOnsetInventoryMembership guards the word-position onset constraints: the
// word-initial inventory omits lh, nh and ç, the medial inventory includes them,
// and both offer the empty (vowel-initial) onset.
func TestOnsetInventoryMembership(t *testing.T) {
	for _, f := range []string{"lh", "nh", "ç"} {
		if inventoryContains(onsetInitialInv, f) {
			t.Errorf("onsetInitialInv must not contain %q", f)
		}
		if !inventoryContains(onsetMedialInv, f) {
			t.Errorf("onsetMedialInv must contain %q", f)
		}
	}
	if !inventoryContains(onsetInitialInv, "") || !inventoryContains(onsetMedialInv, "") {
		t.Error("onset inventories must offer the empty onset")
	}
}

// TestNucleiInventoryPartition guards that the front-nucleus inventory holds only
// front-starting nuclei, the back-nucleus inventory only back-starting nuclei,
// and that together they partition the full nucleus inventory.
func TestNucleiInventoryPartition(t *testing.T) {
	for i := range nucleiFrontInv.forms {
		if !nucleusStartsFront(nucleiFrontInv.forms[i].form) {
			t.Errorf("nucleiFrontInv holds non-front nucleus %q", nucleiFrontInv.forms[i].form)
		}
	}
	for i := range nucleiBackInv.forms {
		if !nucleusStartsBack(nucleiBackInv.forms[i].form) {
			t.Errorf("nucleiBackInv holds non-back nucleus %q", nucleiBackInv.forms[i].form)
		}
	}
	if len(nucleiFrontInv.forms)+len(nucleiBackInv.forms) != len(nuclei.forms) {
		t.Errorf("front(%d)+back(%d) nuclei != total(%d)",
			len(nucleiFrontInv.forms), len(nucleiBackInv.forms), len(nuclei.forms))
	}
}

// TestCodaInventoryMembership guards that each conditioned coda inventory omits
// the nasal that is illegal in its context and that every coda inventory offers
// the empty (open-syllable) coda.
func TestCodaInventoryMembership(t *testing.T) {
	if inventoryContains(codaBeforePBInv, "n") {
		t.Error("codaBeforePBInv must exclude n (nasal before p/b is m)")
	}
	if inventoryContains(codaBeforeOtherInv, "m") {
		t.Error("codaBeforeOtherInv must exclude m (coda m only before p/b)")
	}
	if inventoryContains(codaFinalInv, "n") {
		t.Error("codaFinalInv must exclude n (word-final nasal is m)")
	}
	for _, inv := range []weightedInventory{codaBeforePBInv, codaBeforeOtherInv, codaFinalInv} {
		if !inventoryContains(inv, "") {
			t.Error("coda inventories must offer the empty coda")
		}
	}
}

// inventoryContains reports whether inv holds the given form.
func inventoryContains(inv weightedInventory, form string) bool {
	for i := range inv.forms {
		if inv.forms[i].form == form {
			return true
		}
	}
	return false
}

// TestAssembleStemSingleAllocation asserts the assembler performs exactly one
// allocation for the returned string.
func TestAssembleStemSingleAllocation(t *testing.T) {
	r := newTestRand()
	st := sampleSyllabicStem(r, nil, 5, 2)
	if got := testing.AllocsPerRun(1000, func() { _ = assembleStem(st.syllables) }); got != 1 {
		t.Errorf("assembleStem allocations/op = %v, want 1", got)
	}
}

// TestAssembleStemContent checks that the assembler concatenates the syllable
// parts in order with no separators, and preserves multi-byte pt-PT graphemes.
func TestAssembleStemContent(t *testing.T) {
	tests := []struct {
		stem []syllable
		want string
	}{
		{[]syllable{syl("pr", "a", ""), syl("t", "o", "")}, "prato"},
		{[]syllable{syl("c", "o", ""), syl("r", "a", ""), syl("ç", "ão", "")}, "coração"},
		{[]syllable{syl("c", "a", "m"), syl("p", "o", "")}, "campo"},
		{[]syllable{syl("", "a", "")}, "a"},
	}
	for _, tc := range tests {
		if got := assembleStem(tc.stem); got != tc.want {
			t.Errorf("assembleStem(%+v) = %q, want %q", tc.stem, got, tc.want)
		}
	}
}

// TestFullStemAllocationBudget demonstrates the end-to-end allocation budget:
// with a reused syllable buffer, sampling allocates nothing and the assembled
// string is the single allocation per generated stem.
func TestFullStemAllocationBudget(t *testing.T) {
	r := newTestRand()
	dst := make([]syllable, 0, 16) // pre-sized reusable buffer
	got := testing.AllocsPerRun(1000, func() {
		st := sampleSyllabicStem(r, dst, 4, 1)
		dst = st.syllables
		_ = assembleStem(st.syllables)
	})
	if got != 1 {
		t.Errorf("full stem allocations/op = %v, want 1 (the assembled string)", got)
	}
}

// TestNoRejectionDeterministicDraws proves that the sampler consumes a fixed,
// deterministic number of random draws per call (three per syllable: onset,
// nucleus, coda), which is only possible without a generate-and-reject loop.
func TestNoRejectionDeterministicDraws(t *testing.T) {
	for count := 1; count <= 8; count++ {
		got := countUint32Draws(func(r *rand.Rand) {
			_ = sampleSyllabicStem(r, nil, count, 0)
		})
		want := uint64(3 * count)
		if got != want {
			t.Errorf("count=%d: consumed %d draws, want %d (3 per syllable, no rejection)", count, got, want)
		}
	}
}

// countUint32Draws runs fn with a source that tallies its Uint64 pulls, by
// wrapping a fixed-seed PCG source. The sampler issues exactly three weighted
// draws per syllable with no data-dependent redraw; a weighted draw is one
// Uint32N, which pulls the source once except in an unbiasing case whose
// probability is about n/2^64 (here below 1e-16) and which, for this fixed seed,
// deterministically does not occur. The tally therefore equals the sampler's
// weighted-draw count, proving the work per call is fixed (no rejection loop).
func countUint32Draws(fn func(*rand.Rand)) uint64 {
	c := &countingSource{inner: rand.NewPCG(1, 2)}
	fn(rand.New(c))
	return c.pulls
}

// countingSource is a rand.Source that counts Uint64 pulls, used to prove the
// sampler performs a deterministic number of draws.
type countingSource struct {
	inner *rand.PCG
	pulls uint64
}

func (c *countingSource) Uint64() uint64 {
	c.pulls++
	return c.inner.Uint64()
}

// BenchmarkSampleStem measures sampling into a reused buffer; after warm-up it
// allocates nothing (0 allocs/op), proving there is no per-call allocation and no
// rejection loop.
func BenchmarkSampleStem(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	var dst []syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		dst = sampleSyllabicStem(r, dst, 4, 1).syllables
	}
}

// BenchmarkAssembleStem measures assembly of a fixed stem; it performs exactly
// one allocation per operation (the returned string).
func BenchmarkAssembleStem(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	syllables := sampleSyllabicStem(r, nil, 4, 1).syllables
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = assembleStem(syllables)
	}
}

// BenchmarkBuildStem measures the full path (sample into a reused buffer, then
// assemble); it performs exactly one allocation per operation: the assembled
// string.
func BenchmarkBuildStem(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	var dst []syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		st := sampleSyllabicStem(r, dst, 4, 1)
		dst = st.syllables
		_ = assembleStem(st.syllables)
	}
}

// BenchmarkBuildStemByCount measures the full path across syllable counts, so the
// per-call work is visibly linear in the syllable count (deterministic, no
// rejection).
func BenchmarkBuildStemByCount(b *testing.B) {
	for _, count := range []int{1, 2, 3, 4, 6, 8} {
		b.Run("syllables="+strconv.Itoa(count), func(b *testing.B) {
			r := rand.New(rand.NewPCG(1, 2))
			var dst []syllable
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				st := sampleSyllabicStem(r, dst, count, 0)
				dst = st.syllables
				_ = assembleStem(st.syllables)
			}
		})
	}
}
