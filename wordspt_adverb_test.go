package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the public pt-PT adverb generators (AdverbPT,
// AdverbPTByLengthType, and their *Generator methods). The checks map one-to-one
// to the task acceptance criteria: every adverb ends in -mente, carries NO graphic
// accent (the whole word is unaccented, proving the base lost its accent), is an
// orthographically well-formed lowercase valid-UTF-8 word (the Sprint 7 string
// oracles), is invariable, honors length with the small-category upward
// normalization, and is reproducible through a seeded Generator. It reuses
// countAcuteOrCircumflex from wordspt_adjective_test.go (same package).

// adverbSampleSize is the sampled-run size for the conformance tests. It exceeds
// the task's >=1000 floor comfortably while keeping the suite fast.
const adverbSampleSize = 20000

// assertAdverbOrthographyConformant fails t when word is not an orthographically
// well-formed, lowercase, valid-UTF-8 pt-PT word, per the Sprint 7 string oracles
// and the case/encoding rules.
func assertAdverbOrthographyConformant(t *testing.T, word string) {
	t.Helper()
	if word == "" {
		t.Fatalf("generated an empty adverb")
	}
	if !utf8.ValidString(word) {
		t.Fatalf("adverb %q is not valid UTF-8", word)
	}
	if !wordConformsToCedilla(word) {
		t.Errorf("adverb %q violates the pt-PT cedilla rule", word)
	}
	if wordHasForbiddenDiaeresis(word) {
		t.Errorf("adverb %q contains a forbidden diaeresis", word)
	}
	if wordContainsGrave(word) {
		t.Errorf("adverb %q contains a grave accent", word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("adverb %q contains an uppercase rune %q", word, r)
		}
	}
}

// assertAdverbConformant applies the full per-adverb contract: the orthography
// oracles, the -mente ending, and the absence of any acute or circumflex accent
// (the base having lost its graphic accent, and -mente carrying none).
func assertAdverbConformant(t *testing.T, word string) {
	t.Helper()
	assertAdverbOrthographyConformant(t, word)
	if !strings.HasSuffix(word, "mente") {
		t.Fatalf("adverb %q does not end in -mente", word)
	}
	if got := countAcuteOrCircumflex(word); got != 0 {
		t.Fatalf("adverb %q carries %d acute/circumflex accents, want 0 (unaccented)", word, got)
	}
}

// TestAdverbPTOrthographicConformance samples AdverbPT and every length category
// of AdverbPTByLengthType and asserts that every output is orthographically
// conformant, ends in -mente, is unaccented, lowercase, and valid UTF-8. It
// exercises both the package-level functions (global source) and a seeded
// Generator.
func TestAdverbPTOrthographicConformance(t *testing.T) {
	for i := 0; i < adverbSampleSize; i++ {
		assertAdverbConformant(t, AdverbPT())
	}

	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	g := New(0xADE)
	for _, le := range lengths {
		for i := 0; i < 5000; i++ {
			assertAdverbConformant(t, AdverbPTByLengthType(le))
			assertAdverbConformant(t, g.AdverbPTByLengthType(le))
		}
	}
}

// TestAdverbPTEndsInMente verifies, over a large sample, that every adverb ends in
// -mente (the productive adverbial suffix), across every length category.
func TestAdverbPTEndsInMente(t *testing.T) {
	g := New(3)
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	for _, le := range lengths {
		for i := 0; i < adverbSampleSize; i++ {
			w := g.AdverbPTByLengthType(le)
			if !strings.HasSuffix(w, "mente") {
				t.Fatalf("length %d: adverb %q does not end in -mente", le, w)
			}
		}
	}
}

// TestAdverbPTNoGraphicAccent verifies the core orthographic rule of -mente
// formation: because the primary stress moves to the suffix, the base loses its
// graphic accent and the whole adverb carries no acute or circumflex accent. It
// checks this over a large sample across every length category. (Nasal tildes, a
// mark of nasality rather than stress, are permitted and are not counted here.)
func TestAdverbPTNoGraphicAccent(t *testing.T) {
	g := New(5)
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	for _, le := range lengths {
		for i := 0; i < adverbSampleSize; i++ {
			w := g.AdverbPTByLengthType(le)
			if got := countAcuteOrCircumflex(w); got != 0 {
				t.Fatalf("length %d: adverb %q has %d acute/circumflex accents, want 0", le, w, got)
			}
		}
	}
}

// TestAdverbPTInvariable verifies that a -mente adverb is invariable: it always
// ends exactly in -mente and never carries a plural marker (-mentes) or any other
// inflection, confirming it does not vary for gender, number, or degree.
func TestAdverbPTInvariable(t *testing.T) {
	g := New(23)
	for i := 0; i < adverbSampleSize; i++ {
		w := g.AdverbPT()
		if !strings.HasSuffix(w, "mente") {
			t.Fatalf("adverb %q does not end in -mente", w)
		}
		if strings.HasSuffix(w, "mentes") {
			t.Fatalf("adverb %q is inflected for number (-mentes); adverbs are invariable", w)
		}
		// An invariable -mente adverb always ends in the vowel e, never in a plural
		// -s or a gender marker.
		if strings.HasSuffix(w, "s") {
			t.Fatalf("adverb %q ends in -s; a -mente adverb is invariable and ends in -e", w)
		}
	}
}

// TestAdverbWordFormation verifies the -mente composition directly on the
// single-allocation builder, using the specification's exact examples and the
// full feminine ending inventory. Each stem is BARE (unaccented), exactly as the
// composition stores it, so the builder demonstrates the accent-loss rule: a base
// that would be accented in isolation (rápida, fácil, básica, só, amável,
// possível) yields an unaccented adverb (rapidamente, facilmente, basicamente,
// somente, amavelmente, possivelmente).
func TestAdverbWordFormation(t *testing.T) {
	cases := []struct {
		name string
		stem syllabicStem
		want string
	}{
		{
			name: "bela->belamente (thematic -a)",
			stem: syllabicStem{syllables: []syllable{{onset: "b", nucleus: "e"}, {onset: "l", nucleus: "a"}}, tonic: 0},
			want: "belamente",
		},
		{
			name: "rápida->rapidamente (accent lost)",
			stem: syllabicStem{syllables: []syllable{{onset: "r", nucleus: "a"}, {onset: "p", nucleus: "i"}, {onset: "d", nucleus: "a"}}, tonic: 0},
			want: "rapidamente",
		},
		{
			name: "fácil->facilmente (accent lost, coda l)",
			stem: syllabicStem{syllables: []syllable{{onset: "f", nucleus: "a"}, {onset: "c", nucleus: "i", coda: "l"}}, tonic: 0},
			want: "facilmente",
		},
		{
			name: "só->somente (accent lost, monosyllable)",
			stem: syllabicStem{syllables: []syllable{{onset: "s", nucleus: "o"}}, tonic: 0},
			want: "somente",
		},
		{
			name: "famosa->famosamente (-osa)",
			stem: syllabicStem{syllables: []syllable{{onset: "f", nucleus: "a"}, {onset: "m", nucleus: "o"}, {onset: "s", nucleus: "a"}}, tonic: 1},
			want: "famosamente",
		},
		{
			name: "básica->basicamente (-ica, accent lost)",
			stem: syllabicStem{syllables: []syllable{{onset: "b", nucleus: "a"}, {onset: "s", nucleus: "i"}, {onset: "c", nucleus: "a"}}, tonic: 0},
			want: "basicamente",
		},
		{
			name: "ativa->ativamente (-iva, vowel-initial)",
			stem: syllabicStem{syllables: []syllable{{nucleus: "a"}, {onset: "t", nucleus: "i"}, {onset: "v", nucleus: "a"}}, tonic: 1},
			want: "ativamente",
		},
		{
			name: "legal->legalmente (-al, coda l)",
			stem: syllabicStem{syllables: []syllable{{onset: "l", nucleus: "e"}, {onset: "g", nucleus: "a", coda: "l"}}, tonic: 1},
			want: "legalmente",
		},
		{
			name: "amável->amavelmente (-ável, accent lost)",
			stem: syllabicStem{syllables: []syllable{{nucleus: "a"}, {onset: "m", nucleus: "a"}, {onset: "v", nucleus: "e", coda: "l"}}, tonic: 1},
			want: "amavelmente",
		},
		{
			name: "possível->possivelmente (-ível, accent lost)",
			stem: syllabicStem{syllables: []syllable{{onset: "p", nucleus: "o", coda: "s"}, {onset: "s", nucleus: "i"}, {onset: "v", nucleus: "e", coda: "l"}}, tonic: 1},
			want: "possivelmente",
		},
		{
			name: "elegante->elegantemente (-ante)",
			stem: syllabicStem{syllables: []syllable{{nucleus: "e"}, {onset: "l", nucleus: "e"}, {onset: "g", nucleus: "a", coda: "n"}, {onset: "t", nucleus: "e"}}, tonic: 2},
			want: "elegantemente",
		},
		{
			name: "presente->presentemente (-ente)",
			stem: syllabicStem{syllables: []syllable{{onset: "pr", nucleus: "e"}, {onset: "s", nucleus: "e", coda: "n"}, {onset: "t", nucleus: "e"}}, tonic: 1},
			want: "presentemente",
		},
	}
	for _, tc := range cases {
		got := adverbWord(tc.stem)
		if got != tc.want {
			t.Errorf("%s: adverbWord = %q, want %q", tc.name, got, tc.want)
		}
	}
}

// TestAdverbPTLengthHonored verifies that every output of a concrete length
// category lands within the expected character window. Because the shortest
// feminine adjective base plus -mente already exceeds the Small ceiling (4) and
// the Medium ceiling (8), a Small or Medium request normalizes upward to the Big
// window; a Big request fits Big directly; an Any request spans the full range.
// Length is measured in runes.
func TestAdverbPTLengthHonored(t *testing.T) {
	g := New(29)
	cases := []struct {
		l        LengthTypeWords
		lo, hi   int
		describe string
	}{
		{SmallLengthWord, 9, 30, "Small normalizes up to Big"},
		{MediumLengthWords, 9, 30, "Medium normalizes up to Big"},
		{BigLengthWords, 9, 30, "Big fits Big"},
		{AnyLengthWord, 1, 30, "Any spans the full range"},
	}
	for _, tc := range cases {
		for i := 0; i < adverbSampleSize; i++ {
			w := g.AdverbPTByLengthType(tc.l)
			n := utf8.RuneCountInString(w)
			if n < tc.lo || n > tc.hi {
				t.Fatalf("%s: adverb %q has %d runes, outside window [%d,%d]",
					tc.describe, w, n, tc.lo, tc.hi)
			}
		}
	}
}

// TestAdverbPTSmallNormalizesUp verifies the specification's minimum-viable-length
// consequence for adverbs explicitly: a -mente adverb never fits SmallLengthWord,
// so every Small request produces a word longer than the four-character Small
// ceiling.
func TestAdverbPTSmallNormalizesUp(t *testing.T) {
	g := New(31)
	for i := 0; i < adverbSampleSize; i++ {
		w := g.AdverbPTByLengthType(SmallLengthWord)
		if n := utf8.RuneCountInString(w); n <= 4 {
			t.Fatalf("Small adverb %q has %d runes, but a -mente adverb cannot fit the Small ceiling of 4", w, n)
		}
	}
}

// TestAdverbPTEndingVariety verifies that, over a broad (Any length) sample, every
// feminine adjective ending family is reachable through the adverb surface,
// confirming the ending inventory is fully wired into the adverb path. The
// derivational and invariable markers (-osamente, -icamente, -ivamente, -almente,
// -avelmente, -ivelmente, -entemente, -antemente) are checked, along with the
// plain thematic -a base (an -amente surface that is none of -osa/-ica/-iva).
func TestAdverbPTEndingVariety(t *testing.T) {
	g := New(505)
	markers := []string{
		"osamente", "icamente", "ivamente", "almente",
		"avelmente", "ivelmente", "entemente", "antemente",
	}
	seen := map[string]bool{}
	sawPlainA := false
	for i := 0; i < adverbSampleSize; i++ {
		w := g.AdverbPT()
		for _, m := range markers {
			if strings.HasSuffix(w, m) {
				seen[m] = true
			}
		}
		if strings.HasSuffix(w, "amente") &&
			!strings.HasSuffix(w, "osamente") &&
			!strings.HasSuffix(w, "icamente") &&
			!strings.HasSuffix(w, "ivamente") {
			sawPlainA = true
		}
	}
	for _, m := range markers {
		if !seen[m] {
			t.Errorf("adverb ending marker %q never observed", m)
		}
	}
	if !sawPlainA {
		t.Errorf("plain thematic -a adverb (-amente) never observed")
	}
}

// TestAdverbEndingLengthAwareSelection verifies the length-aware ending selection
// the adverb path relies on: only the feminine thematic -a fits the Small ceiling
// of four characters, so a Small ceiling selects only that ending; with a Big
// ceiling every feminine ending is reachable. (The adverb still normalizes its
// LENGTH upward out of Small, but the ending is selected against the raw ceiling,
// exactly as the adjective and noun paths do.)
func TestAdverbEndingLengthAwareSelection(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))

	small := map[string]bool{}
	for i := 0; i < 5000; i++ {
		small[selectFittingEnding(r, feminineAdjectiveEndings, 4).label] = true
	}
	if len(small) != 1 || !small["-a"] {
		t.Errorf("Small feminine selection = %v, want only -a", small)
	}

	big := map[string]bool{}
	for i := 0; i < 5000; i++ {
		big[selectFittingEnding(r, feminineAdjectiveEndings, 30).label] = true
	}
	for _, want := range []string{"-a", "-osa", "-ica", "-iva", "-al", "-ável", "-ível", "-ente", "-ante"} {
		if !big[want] {
			t.Errorf("Big feminine selection missing ending %q (got %v)", want, big)
		}
	}
}

// TestAdverbPTReproducible verifies that a seeded Generator is deterministic: two
// generators created with the same seed produce identical adverb sequences, for
// both the sugar method and the length-typed method across every length category.
func TestAdverbPTReproducible(t *testing.T) {
	a := New(0xC0FFEE)
	b := New(0xC0FFEE)
	for i := 0; i < adverbSampleSize; i++ {
		if x, y := a.AdverbPT(), b.AdverbPT(); x != y {
			t.Fatalf("AdverbPT diverged at %d: %q != %q", i, x, y)
		}
	}

	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	c := New(42)
	d := New(42)
	for _, le := range lengths {
		for i := 0; i < 2000; i++ {
			if x, y := c.AdverbPTByLengthType(le), d.AdverbPTByLengthType(le); x != y {
				t.Fatalf("AdverbPTByLengthType(%d) diverged at %d: %q != %q", le, i, x, y)
			}
		}
	}
}

// TestAdverbPTConcurrentSafe exercises the package-level AdverbPT from many
// goroutines. Run with -race, it demonstrates that the package-level surface is
// safe for concurrent use, as documented: its shared global source delegates each
// draw to the concurrency-safe global math/rand/v2 generator.
func TestAdverbPTConcurrentSafe(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 4000
	done := make(chan struct{}, goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for i := 0; i < perGoroutine; i++ {
				if AdverbPT() == "" {
					t.Error("AdverbPT returned an empty string")
					return
				}
			}
		}()
	}
	for g := 0; g < goroutines; g++ {
		<-done
	}
}

// BenchmarkAdverbPTByLengthType measures the adverb path and its allocation
// budget: a -mente adverb must assemble in exactly one string allocation (the bare
// stem and the suffix are written into a single pre-sized builder).
func BenchmarkAdverbPTByLengthType(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.AdverbPTByLengthType(BigLengthWords)
	}
}

// BenchmarkAdverbPT measures the default package-level adverb path over the global
// source (random length).
func BenchmarkAdverbPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = AdverbPT()
	}
}
