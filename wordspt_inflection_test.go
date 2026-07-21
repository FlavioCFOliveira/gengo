package gengo

import (
	"math/rand/v2"
	"testing"
	"unicode/utf8"
)

// The tests in this file exercise the WordsPT flexion enums, the Any-resolution
// helpers, and the length machinery (wordspt_inflection.go): the enum zero-value
// contract, uniform and reproducible Any resolution, the category character
// windows, the minimum-viable-length normalization, and the character-window to
// syllable-count mapping (including an end-to-end demonstration that the mapping
// lands generated bare stems inside the targeted window).

// TestFlexionEnumZeroValues locks the specification contract that every flexion
// enum's zero value is Any and that the concrete values are distinct.
func TestFlexionEnumZeroValues(t *testing.T) {
	if AnyGender != 0 || Masculine != 1 || Feminine != 2 {
		t.Errorf("Gender constants = {%d,%d,%d}, want {0,1,2}", AnyGender, Masculine, Feminine)
	}
	if AnyNumber != 0 || Singular != 1 || Plural != 2 {
		t.Errorf("Number constants = {%d,%d,%d}, want {0,1,2}", AnyNumber, Singular, Plural)
	}
	if AnyDegree != 0 || Positive != 1 || Superlative != 2 {
		t.Errorf("Degree constants = {%d,%d,%d}, want {0,1,2}", AnyDegree, Positive, Superlative)
	}
	if AnyLengthWord != 0 {
		t.Errorf("AnyLengthWord = %d, want 0", AnyLengthWord)
	}
	// The zero value of each enum type must equal its Any constant.
	var g Gender
	var n Number
	var d Degree
	var l LengthTypeWords
	if g != AnyGender || n != AnyNumber || d != AnyDegree || l != AnyLengthWord {
		t.Error("the zero value of a flexion enum must be its Any constant")
	}
}

// TestResolveConcretePassthrough verifies that a concrete request is returned
// unchanged, never re-randomized.
func TestResolveConcretePassthrough(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 100; i++ {
		if got := resolveGender(r, Masculine); got != Masculine {
			t.Fatalf("resolveGender(Masculine) = %d", got)
		}
		if got := resolveGender(r, Feminine); got != Feminine {
			t.Fatalf("resolveGender(Feminine) = %d", got)
		}
		if got := resolveNumber(r, Singular); got != Singular {
			t.Fatalf("resolveNumber(Singular) = %d", got)
		}
		if got := resolveNumber(r, Plural); got != Plural {
			t.Fatalf("resolveNumber(Plural) = %d", got)
		}
		if got := resolveDegree(r, Positive); got != Positive {
			t.Fatalf("resolveDegree(Positive) = %d", got)
		}
		if got := resolveDegree(r, Superlative); got != Superlative {
			t.Fatalf("resolveDegree(Superlative) = %d", got)
		}
	}
}

// TestResolveAnyUniformAndValid verifies that Any (and any out-of-range value)
// resolves only to valid concrete values, that both values are reachable, and
// that the split is close to uniform.
func TestResolveAnyUniformAndValid(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	const n = 200000

	countGender := map[Gender]int{}
	countNumber := map[Number]int{}
	countDegree := map[Degree]int{}
	for i := 0; i < n; i++ {
		countGender[resolveGender(r, AnyGender)]++
		countNumber[resolveNumber(r, AnyNumber)]++
		countDegree[resolveDegree(r, AnyDegree)]++
	}

	assertBinaryUniform(t, "Gender", countGender[Masculine], countGender[Feminine], n)
	if len(countGender) != 2 {
		t.Errorf("resolveGender(Any) produced %d distinct values, want 2", len(countGender))
	}
	assertBinaryUniform(t, "Number", countNumber[Singular], countNumber[Plural], n)
	if len(countNumber) != 2 {
		t.Errorf("resolveNumber(Any) produced %d distinct values, want 2", len(countNumber))
	}
	assertBinaryUniform(t, "Degree", countDegree[Positive], countDegree[Superlative], n)
	if len(countDegree) != 2 {
		t.Errorf("resolveDegree(Any) produced %d distinct values, want 2", len(countDegree))
	}

	// Out-of-range values normalize to a valid concrete value (never error/panic).
	for i := 0; i < 1000; i++ {
		if g := resolveGender(r, Gender(200)); g != Masculine && g != Feminine {
			t.Fatalf("resolveGender(out-of-range) = %d, want a valid gender", g)
		}
		if nn := resolveNumber(r, Number(200)); nn != Singular && nn != Plural {
			t.Fatalf("resolveNumber(out-of-range) = %d, want a valid number", nn)
		}
		if d := resolveDegree(r, Degree(200)); d != Positive && d != Superlative {
			t.Fatalf("resolveDegree(out-of-range) = %d, want a valid degree", d)
		}
	}
}

// assertBinaryUniform fails if the two counts diverge from a 50/50 split by more
// than three percentage points.
func assertBinaryUniform(t *testing.T, name string, a, b, total int) {
	t.Helper()
	fa := float64(a) / float64(total)
	if fa < 0.47 || fa > 0.53 {
		t.Errorf("%s Any split not uniform: %d vs %d (%.3f)", name, a, b, fa)
	}
}

// TestResolveReproducible confirms that Any resolution is reproducible through a
// seeded source, so the *Generator surface is deterministic.
func TestResolveReproducible(t *testing.T) {
	a := New(99)
	b := New(99)
	for i := 0; i < 1000; i++ {
		if resolveGender(a.r, AnyGender) != resolveGender(b.r, AnyGender) {
			t.Fatalf("resolveGender diverged at i=%d", i)
		}
		if resolveNumber(a.r, AnyNumber) != resolveNumber(b.r, AnyNumber) {
			t.Fatalf("resolveNumber diverged at i=%d", i)
		}
		if resolveDegree(a.r, AnyDegree) != resolveDegree(b.r, AnyDegree) {
			t.Fatalf("resolveDegree diverged at i=%d", i)
		}
	}
}

// TestCharRangeOf covers the category character windows, including the Any and
// out-of-range full-range default.
func TestCharRangeOf(t *testing.T) {
	tests := []struct {
		l                LengthTypeWords
		wantMin, wantMax int
	}{
		{SmallLengthWord, 1, 4},
		{MediumLengthWords, 5, 8},
		{BigLengthWords, 9, 30},
		{AnyLengthWord, 1, 30},
		{LengthTypeWords(200), 1, 30},
	}
	for _, tc := range tests {
		gotMin, gotMax := charRangeOf(tc.l)
		if gotMin != tc.wantMin || gotMax != tc.wantMax {
			t.Errorf("charRangeOf(%d) = (%d,%d), want (%d,%d)", tc.l, gotMin, gotMax, tc.wantMin, tc.wantMax)
		}
	}
}

// TestNormalizeLength covers the minimum-viable-length rule: the category is
// raised to the smallest one that can hold a minViable-character word and is never
// lowered; Any is preserved for any realistic minimum.
func TestNormalizeLength(t *testing.T) {
	tests := []struct {
		l        LengthTypeWords
		minVia   int
		want     LengthTypeWords
		describe string
	}{
		{SmallLengthWord, 1, SmallLengthWord, "small holds 1"},
		{SmallLengthWord, 4, SmallLengthWord, "small holds 4"},
		{SmallLengthWord, 5, MediumLengthWords, "small too short for 5 -> medium"},
		{SmallLengthWord, 8, MediumLengthWords, "small too short for 8 -> medium"},
		{SmallLengthWord, 9, BigLengthWords, "small too short for 9 -> big"},
		{MediumLengthWords, 8, MediumLengthWords, "medium holds 8"},
		{MediumLengthWords, 9, BigLengthWords, "medium too short for 9 -> big"},
		{MediumLengthWords, 4, MediumLengthWords, "never lowered below request"},
		{BigLengthWords, 30, BigLengthWords, "big holds 30"},
		{BigLengthWords, 31, BigLengthWords, "big clamps at big"},
		{AnyLengthWord, 8, AnyLengthWord, "any holds 8"},
		{AnyLengthWord, 30, AnyLengthWord, "any holds 30"},
		{AnyLengthWord, 31, BigLengthWords, "any beyond 30 -> big"},
	}
	for _, tc := range tests {
		if got := normalizeLength(tc.l, tc.minVia); got != tc.want {
			t.Errorf("normalizeLength(%d, %d) = %d, want %d [%s]", tc.l, tc.minVia, got, tc.want, tc.describe)
		}
	}
}

// TestSyllablesForCharRange verifies that the character-window to syllable-count
// mapping stays within the empirically calibrated bounds for each category,
// always returns at least one syllable, and is monotonic across categories.
func TestSyllablesForCharRange(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 6))
	const n = 100000
	type bound struct {
		l              LengthTypeWords
		loWant, hiWant int
	}
	// Bounds derived from runesPerSyllable = 2.48 and each category window.
	bounds := []bound{
		{SmallLengthWord, 1, 2},   // targets 1..4 -> counts 1..2
		{MediumLengthWords, 2, 3}, // targets 5..8 -> counts 2..3
		{BigLengthWords, 4, 12},   // targets 9..30 -> counts 4..12
		{AnyLengthWord, 1, 12},    // targets 1..30 -> counts 1..12
	}
	prevMax := 0
	for _, b := range bounds {
		minC, maxC := charRangeOf(b.l)
		observedMin, observedMax := 1<<30, 0
		for i := 0; i < n; i++ {
			k := syllablesForCharRange(r, minC, maxC)
			if k < 1 {
				t.Fatalf("syllablesForCharRange(%d,%d) = %d < 1", minC, maxC, k)
			}
			if k < observedMin {
				observedMin = k
			}
			if k > observedMax {
				observedMax = k
			}
		}
		if observedMin < b.loWant || observedMax > b.hiWant {
			t.Errorf("category %d: observed counts [%d,%d], want within [%d,%d]",
				b.l, observedMin, observedMax, b.loWant, b.hiWant)
		}
		if observedMax < prevMax {
			t.Errorf("mapping not monotonic across categories: category %d max %d < previous %d",
				b.l, observedMax, prevMax)
		}
		prevMax = observedMax
	}
}

// TestStemSyllablesForNormalizesAndSizes verifies that stemSyllablesFor honors
// the minimum-viable normalization (a too-short category is sized as the
// normalized-up category) and that a larger fixed-affix cost shrinks the stem.
func TestStemSyllablesForNormalizesAndSizes(t *testing.T) {
	r := rand.New(rand.NewPCG(8, 9))
	const n = 100000

	// A Small request with minViable 6 must be sized as Medium (counts 2..3),
	// not as Small (counts 1..2).
	sawMediumLike := false
	for i := 0; i < n; i++ {
		k := stemSyllablesFor(r, SmallLengthWord, 6, 0)
		if k < 1 {
			t.Fatalf("stemSyllablesFor returned %d < 1", k)
		}
		if k >= 2 {
			sawMediumLike = true
		}
	}
	if !sawMediumLike {
		t.Error("stemSyllablesFor(Small, minViable=6) did not normalize up to a Medium-sized stem")
	}

	// A larger affix cost must, on average, reduce the stem syllable count for the
	// same category.
	var sum0, sumAffix int
	for i := 0; i < n; i++ {
		sum0 += stemSyllablesFor(r, BigLengthWords, 1, 0)
		sumAffix += stemSyllablesFor(r, BigLengthWords, 1, 8)
	}
	if sumAffix >= sum0 {
		t.Errorf("affix cost did not shrink the stem: affix-8 total %d >= affix-0 total %d", sumAffix, sum0)
	}

	// The stem count is never below one, even when the affix exceeds the window.
	for i := 0; i < 1000; i++ {
		if k := stemSyllablesFor(r, SmallLengthWord, 1, 100); k < 1 {
			t.Fatalf("stemSyllablesFor with huge affix returned %d < 1", k)
		}
	}
}

// TestStemSyllablesForLandsInWindow is the end-to-end demonstration that the
// length-category mapping targets the requested window: with no fixed affix
// (affixRunes = 0), a bare accented stem built at the mapped syllable count lands
// inside the category's character window for the large majority of samples. The
// residual is the inherent per-syllable length variance (1..5 runes), which the
// class enforces exactly on the assembled word.
func TestStemSyllablesForLandsInWindow(t *testing.T) {
	r := rand.New(rand.NewPCG(0xF00D, 0x1234))
	const n = 50000
	cases := []struct {
		l       LengthTypeWords
		minFrac float64
	}{
		{SmallLengthWord, 0.70},
		{MediumLengthWords, 0.60},
		{BigLengthWords, 0.75},
	}
	var dst []syllable
	for _, tc := range cases {
		lo, hi := charRangeOf(tc.l)
		inWindow := 0
		for i := 0; i < n; i++ {
			k := stemSyllablesFor(r, tc.l, 1, 0)
			st := sampleSyllabicStem(r, dst, k, i%k)
			dst = st.syllables
			runes := utf8.RuneCountInString(accentedWord(st))
			if runes >= lo && runes <= hi {
				inWindow++
			}
		}
		frac := float64(inWindow) / float64(n)
		if frac < tc.minFrac {
			t.Errorf("category %d: only %.3f of bare stems landed in window [%d,%d], want >= %.2f",
				tc.l, frac, lo, hi, tc.minFrac)
		}
	}
}

// BenchmarkResolveGender measures the Any-resolution path (a single draw).
func BenchmarkResolveGender(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = resolveGender(r, AnyGender)
	}
}

// BenchmarkStemSyllablesFor measures the length mapping (normalization plus the
// calibrated inversion); it allocates nothing.
func BenchmarkStemSyllablesFor(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = stemSyllablesFor(r, MediumLengthWords, 5, 2)
	}
}
