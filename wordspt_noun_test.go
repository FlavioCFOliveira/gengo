package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the public pt-PT noun generators (NounPT, NounPTOf, and
// their *Generator methods). The checks map one-to-one to the task acceptance
// criteria: orthographic conformance over a large sample (the Sprint 7 string
// oracles plus valid, lowercase UTF-8), gender/number/length honoring with Any
// producing varied output, upward length normalization when a category is too
// short for an ending, and reproducibility of a seeded Generator.

// nounSampleSize is the sampled-run size for the conformance tests. It exceeds the
// task's >=1000 floor comfortably while keeping the suite fast.
const nounSampleSize = 20000

// assertNounOrthographyConformant fails t when word is not an orthographically
// well-formed, lowercase, valid-UTF-8 pt-PT word, per the Sprint 7 string oracles
// and the case/encoding rules.
func assertNounOrthographyConformant(t *testing.T, word string) {
	t.Helper()
	if word == "" {
		t.Fatalf("generated an empty noun")
	}
	if !utf8.ValidString(word) {
		t.Fatalf("noun %q is not valid UTF-8", word)
	}
	if !wordConformsToCedilla(word) {
		t.Errorf("noun %q violates the pt-PT cedilla rule", word)
	}
	if wordHasForbiddenDiaeresis(word) {
		t.Errorf("noun %q contains a forbidden diaeresis", word)
	}
	if wordContainsGrave(word) {
		t.Errorf("noun %q contains a grave accent", word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("noun %q contains an uppercase rune %q", word, r)
		}
	}
}

// nounGenderCompat reports the grammatical genders a SINGULAR noun surface is
// compatible with, by matching its class ending (most specific suffix first). A
// gender-inflecting or gender-fixed ending yields a single gender; the
// common-gender -ista yields both. It is a test oracle for gender honoring.
func nounGenderCompat(w string) (masc, fem bool) {
	switch {
	case strings.HasSuffix(w, "ção"):
		return false, true
	case strings.HasSuffix(w, "dade"):
		return false, true
	case strings.HasSuffix(w, "agem"):
		return false, true
	case strings.HasSuffix(w, "mento"):
		return true, false
	case strings.HasSuffix(w, "eiro"):
		return true, false
	case strings.HasSuffix(w, "eira"):
		return false, true
	case strings.HasSuffix(w, "ista"):
		return true, true
	case strings.HasSuffix(w, "ora"):
		return false, true
	case strings.HasSuffix(w, "or"):
		return true, false
	case strings.HasSuffix(w, "o"):
		return true, false
	case strings.HasSuffix(w, "a"):
		return false, true
	}
	return false, false
}

// TestNounPTOrthographicConformance samples NounPT and every gender/number/length
// combination of NounPTOf and asserts that every output is orthographically
// conformant, lowercase, and valid UTF-8.
func TestNounPTOrthographicConformance(t *testing.T) {
	for i := 0; i < nounSampleSize; i++ {
		assertNounOrthographyConformant(t, NounPT())
	}

	genders := []Gender{AnyGender, Masculine, Feminine}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	g := New(0xB0B)
	for _, ge := range genders {
		for _, nu := range numbers {
			for _, le := range lengths {
				for i := 0; i < 2000; i++ {
					assertNounOrthographyConformant(t, NounPTOf(ge, nu, le))
					assertNounOrthographyConformant(t, g.NounPTOf(ge, nu, le))
				}
			}
		}
	}
}

// TestNounPTGenderHonored verifies that a requested concrete gender is realized:
// a masculine request never yields a feminine-only ending, and vice versa. It
// tests singular forms, where the surface gender marker is unambiguous.
func TestNounPTGenderHonored(t *testing.T) {
	g := New(7)
	for i := 0; i < nounSampleSize; i++ {
		m := g.NounPTOf(Masculine, Singular, AnyLengthWord)
		if masc, _ := nounGenderCompat(m); !masc {
			t.Fatalf("Masculine request produced feminine-only noun %q", m)
		}
		f := g.NounPTOf(Feminine, Singular, AnyLengthWord)
		if _, fem := nounGenderCompat(f); !fem {
			t.Fatalf("Feminine request produced masculine-only noun %q", f)
		}
	}
}

// TestNounPTNumberHonored verifies that number is honored. Every pt-PT noun
// plural in the ending inventory ends in -s, and no singular ending does, so the
// final -s is a sound discriminator.
func TestNounPTNumberHonored(t *testing.T) {
	g := New(11)
	for i := 0; i < nounSampleSize; i++ {
		s := g.NounPTOf(AnyGender, Singular, AnyLengthWord)
		if strings.HasSuffix(s, "s") {
			t.Fatalf("Singular request produced a plural-looking noun %q", s)
		}
		p := g.NounPTOf(AnyGender, Plural, AnyLengthWord)
		if !strings.HasSuffix(p, "s") {
			t.Fatalf("Plural request produced a singular-looking noun %q", p)
		}
	}
}

// TestNounPTLengthHonored verifies that every output of a concrete length
// category lands within the expected character window. A singular noun fits its
// requested category exactly. A Small plural cannot fit four characters (the
// shortest plural noun already exceeds a comfortable Small word), so it normalizes
// upward to the Medium window, per the plural minimum-viable-length rule; Medium
// and Big plurals fit their category. Length is measured in runes.
func TestNounPTLengthHonored(t *testing.T) {
	g := New(29)
	cases := []struct {
		l        LengthTypeWords
		n        Number
		lo, hi   int
		describe string
	}{
		{SmallLengthWord, Singular, 1, 4, "Small singular fits Small"},
		{SmallLengthWord, Plural, 5, 8, "Small plural normalizes up to Medium"},
		{MediumLengthWords, Singular, 5, 8, "Medium singular fits Medium"},
		{MediumLengthWords, Plural, 5, 8, "Medium plural fits Medium"},
		{BigLengthWords, Singular, 9, 30, "Big singular fits Big"},
		{BigLengthWords, Plural, 9, 30, "Big plural fits Big"},
	}
	for _, tc := range cases {
		for i := 0; i < nounSampleSize; i++ {
			w := g.NounPTOf(AnyGender, tc.n, tc.l)
			n := utf8.RuneCountInString(w)
			if n < tc.lo || n > tc.hi {
				t.Fatalf("%s: noun %q has %d runes, outside window [%d,%d]",
					tc.describe, w, n, tc.lo, tc.hi)
			}
		}
	}
}

// TestNounPTAnyGenderVaried verifies that AnyGender yields both masculine- and
// feminine-compatible nouns over a sampled run.
func TestNounPTAnyGenderVaried(t *testing.T) {
	g := New(101)
	sawMasc, sawFem := false, false
	for i := 0; i < nounSampleSize; i++ {
		masc, fem := nounGenderCompat(g.NounPTOf(AnyGender, Singular, AnyLengthWord))
		if masc && !fem {
			sawMasc = true
		}
		if fem && !masc {
			sawFem = true
		}
	}
	if !sawMasc || !sawFem {
		t.Errorf("Any gender not varied: masc=%v fem=%v", sawMasc, sawFem)
	}
}

// TestNounPTAnyNumberVaried verifies that AnyNumber yields both singular and
// plural nouns over a sampled run.
func TestNounPTAnyNumberVaried(t *testing.T) {
	g := New(202)
	sawSing, sawPlur := false, false
	for i := 0; i < nounSampleSize; i++ {
		if strings.HasSuffix(g.NounPTOf(AnyGender, AnyNumber, AnyLengthWord), "s") {
			sawPlur = true
		} else {
			sawSing = true
		}
	}
	if !sawSing || !sawPlur {
		t.Errorf("Any number not varied: singular=%v plural=%v", sawSing, sawPlur)
	}
}

// TestNounPTAnyLengthVaried verifies that AnyLengthWord yields nouns across the
// small, medium and big character buckets over a sampled run.
func TestNounPTAnyLengthVaried(t *testing.T) {
	g := New(303)
	counts := map[int]int{}
	for i := 0; i < nounSampleSize; i++ {
		n := utf8.RuneCountInString(g.NounPTOf(AnyGender, Singular, AnyLengthWord))
		switch {
		case n <= 4:
			counts[1]++
		case n <= 8:
			counts[2]++
		default:
			counts[3]++
		}
	}
	if counts[1] == 0 || counts[2] == 0 || counts[3] == 0 {
		t.Errorf("Any length not varied: small=%d medium=%d big=%d", counts[1], counts[2], counts[3])
	}
}

// TestNounPTReproducible verifies that a seeded Generator is deterministic: two
// generators created with the same seed produce identical noun sequences, for
// both the sugar method and the fully specified method.
func TestNounPTReproducible(t *testing.T) {
	a := New(0xC0FFEE)
	b := New(0xC0FFEE)
	for i := 0; i < nounSampleSize; i++ {
		if x, y := a.NounPT(), b.NounPT(); x != y {
			t.Fatalf("NounPT diverged at %d: %q != %q", i, x, y)
		}
	}

	genders := []Gender{AnyGender, Masculine, Feminine}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	c := New(42)
	d := New(42)
	for _, ge := range genders {
		for _, nu := range numbers {
			for _, le := range lengths {
				for i := 0; i < 200; i++ {
					if x, y := c.NounPTOf(ge, nu, le), d.NounPTOf(ge, nu, le); x != y {
						t.Fatalf("NounPTOf(%d,%d,%d) diverged at %d: %q != %q", ge, nu, le, i, x, y)
					}
				}
			}
		}
	}
}

// TestNounEndingLengthAwareSelection verifies the length-aware ending selection
// that makes strict length honoring possible: a Small request only ever selects
// endings that a Small word can contain (the short thematic -o and -a), and drops
// the longer endings; and when a category is too short for an ending, that
// ending's minimum viable length normalizes the category upward.
func TestNounEndingLengthAwareSelection(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))

	// Only -o (masc) and -a (fem) fit the Small ceiling of 4 characters.
	smallMasc := map[string]bool{}
	smallFem := map[string]bool{}
	for i := 0; i < 5000; i++ {
		smallMasc[selectNounEnding(r, Masculine, 4).label] = true
		smallFem[selectNounEnding(r, Feminine, 4).label] = true
	}
	if len(smallMasc) != 1 || !smallMasc["-o"] {
		t.Errorf("Small masculine selection = %v, want only -o", smallMasc)
	}
	if len(smallFem) != 1 || !smallFem["-a"] {
		t.Errorf("Small feminine selection = %v, want only -a", smallFem)
	}

	// The long endings must be reachable once the ceiling is the Medium/Big range.
	bigMasc := map[string]bool{}
	for i := 0; i < 5000; i++ {
		bigMasc[selectNounEnding(r, Masculine, 30).label] = true
	}
	for _, want := range []string{"-o", "-or", "-mento", "-eiro", "-ista"} {
		if !bigMasc[want] {
			t.Errorf("Big masculine selection missing ending %q (got %v)", want, bigMasc)
		}
	}

	// An ending whose minimum viable length exceeds a category ceiling normalizes
	// that category upward (never downward).
	minViable := endMascMento.endingMinViable() // 7 = affix 5 + CV floor 2
	if eff := normalizeLength(SmallLengthWord, minViable); eff == SmallLengthWord {
		t.Errorf("-mento (minViable %d) did not normalize Small upward", minViable)
	}
}

// TestNounPTConcurrentSafe exercises the package-level NounPT from many
// goroutines. Run with -race, it demonstrates that the package-level surface is
// safe for concurrent use, as documented: its shared global source delegates each
// draw to the concurrency-safe global math/rand/v2 generator.
func TestNounPTConcurrentSafe(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 4000
	done := make(chan struct{}, goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for i := 0; i < perGoroutine; i++ {
				if NounPT() == "" {
					t.Error("NounPT returned an empty string")
					return
				}
			}
		}()
	}
	for g := 0; g < goroutines; g++ {
		<-done
	}
}

// accentedSyllables returns a copy of a stem's syllables with the graphic accent
// applied to the tonic nucleus, so the decomposition-level accentuation oracle
// (stemConformsToAccentuation, which expects accented nuclei) can be applied. It
// mirrors the substitution assembleAccentedStem performs while assembling.
func accentedSyllables(stem syllabicStem) []syllable {
	out := make([]syllable, len(stem.syllables))
	copy(out, stem.syllables)
	plan := accentuateStem(stem)
	if plan.tonic >= 0 {
		nuc := []rune(out[plan.tonic].nucleus)
		nuc[0] = plan.accented
		out[plan.tonic].nucleus = string(nuc)
	}
	return out
}

// TestNounStemConformsToOracles builds combined noun stems for every ending and
// asserts, over a sampled run, that each stem satisfies the Sprint 7
// decomposition oracles: the phonotactic oracle (validating the leading/ending
// junction this task introduces) and, after applying the accent, the accentuation
// oracle (validating that the ending's stress lands the accent correctly on the
// combined word). This is stronger than the string oracles, which do not inspect
// syllable structure or accent placement.
func TestNounStemConformsToOracles(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5EED, 0xF00D))
	endings := []*nounEnding{
		&endMascO, &endFemA, &endMascEiro, &endFemEira, &endMascOr,
		&endFemOra, &endMascMento, &endFemCao, &endFemDade, &endFemAgem, &endComumIsta,
	}
	var buf []syllable
	for _, end := range endings {
		for i := 0; i < 3000; i++ {
			leading := 1 + int(r.Uint32N(6))
			stem := sampleNounStem(r, buf, leading, end)
			buf = stem.syllables
			if !stemConformsToPhonotactics(stem.syllables) {
				t.Fatalf("ending %s: stem %q fails phonotactic oracle", end.label, accentedWord(stem))
			}
			if !stemConformsToAccentuation(accentedSyllables(stem), stem.tonic) {
				t.Fatalf("ending %s: stem %q fails accentuation oracle (tonic %d)",
					end.label, accentedWord(stem), stem.tonic)
			}
		}
	}
}

// BenchmarkNounPTOfSingular measures the singular noun path and its allocation
// budget: a singular noun must assemble in exactly one string allocation.
func BenchmarkNounPTOfSingular(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.NounPTOf(AnyGender, Singular, AnyLengthWord)
	}
}

// BenchmarkNounPTOfPlural measures the plural noun path, which adds one tail
// concatenation in the shared pluralizer over the singular assembly.
func BenchmarkNounPTOfPlural(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.NounPTOf(AnyGender, Plural, AnyLengthWord)
	}
}

// BenchmarkNounPT measures the default package-level noun path over the global
// source (random gender, number and length).
func BenchmarkNounPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = NounPT()
	}
}
