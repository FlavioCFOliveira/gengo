package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the public pt-PT adjective generators (AdjectivePT,
// AdjectivePTOf, and their *Generator methods). The checks map one-to-one to the
// task acceptance criteria: orthographic conformance over a large sample across
// ALL degrees (the Sprint 7 string oracles plus valid, lowercase UTF-8),
// gender/number/degree/length honoring with Any producing varied output, the
// correctness of the -íssimo superlative (its c->qu / g->gu adjustments, its
// single acute, and its gender/number inflection), and reproducibility of a seeded
// Generator.

// adjSampleSize is the sampled-run size for the conformance tests. It exceeds the
// task's >=1000 floor comfortably while keeping the suite fast.
const adjSampleSize = 20000

// assertAdjectiveOrthographyConformant fails t when word is not an
// orthographically well-formed, lowercase, valid-UTF-8 pt-PT word, per the Sprint
// 7 string oracles and the case/encoding rules.
func assertAdjectiveOrthographyConformant(t *testing.T, word string) {
	t.Helper()
	if word == "" {
		t.Fatalf("generated an empty adjective")
	}
	if !utf8.ValidString(word) {
		t.Fatalf("adjective %q is not valid UTF-8", word)
	}
	if !wordConformsToCedilla(word) {
		t.Errorf("adjective %q violates the pt-PT cedilla rule", word)
	}
	if wordHasForbiddenDiaeresis(word) {
		t.Errorf("adjective %q contains a forbidden diaeresis", word)
	}
	if wordContainsGrave(word) {
		t.Errorf("adjective %q contains a grave accent", word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("adjective %q contains an uppercase rune %q", word, r)
		}
	}
}

// countAcuteOrCircumflex returns how many runes of word bear an acute or
// circumflex accent. A superlative must carry exactly one (the í of -íssimo); a
// positive carries at most one (its tonic, when the rules require it).
func countAcuteOrCircumflex(word string) int {
	n := 0
	for _, r := range word {
		if hasAcute(r) || hasCircumflex(r) {
			n++
		}
	}
	return n
}

// adjSuperlativeSuffixes are the four inflected -íssimo endings, used to detect a
// superlative surface and its gender/number.
var adjSuperlativeSuffixes = []string{"íssimo", "íssima", "íssimos", "íssimas"}

// isSuperlativeSurface reports whether word ends in one of the four -íssimo forms.
func isSuperlativeSurface(word string) bool {
	for _, s := range adjSuperlativeSuffixes {
		if strings.HasSuffix(word, s) {
			return true
		}
	}
	return false
}

// TestAdjectivePTOrthographicConformance samples AdjectivePT and every
// gender/number/degree/length combination of AdjectivePTOf and asserts that every
// output is orthographically conformant, lowercase, and valid UTF-8. It exercises
// both the package-level function (global source) and a seeded Generator.
func TestAdjectivePTOrthographicConformance(t *testing.T) {
	for i := 0; i < adjSampleSize; i++ {
		assertAdjectiveOrthographyConformant(t, AdjectivePT())
	}

	genders := []Gender{AnyGender, Masculine, Feminine}
	numbers := []Number{AnyNumber, Singular, Plural}
	degrees := []Degree{AnyDegree, Positive, Superlative}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	g := New(0xADE)
	for _, ge := range genders {
		for _, nu := range numbers {
			for _, de := range degrees {
				for _, le := range lengths {
					for i := 0; i < 600; i++ {
						assertAdjectiveOrthographyConformant(t, AdjectivePTOf(ge, nu, de, le))
						assertAdjectiveOrthographyConformant(t, g.AdjectivePTOf(ge, nu, de, le))
					}
				}
			}
		}
	}
}

// TestAdjectivePTGenderHonored verifies that a requested concrete gender is
// realized. For a positive singular adjective the surface vowel is an unambiguous
// gender marker: an -o ending (oso/ico/ivo) is masculine-only and an -a ending
// (osa/ica/iva) is feminine-only, while the invariable endings (-al, -ável,
// -ível, -ente, -ante) serve either gender. For a superlative the gender shows in
// the -íssimo vs -íssima suffix.
func TestAdjectivePTGenderHonored(t *testing.T) {
	g := New(7)
	for i := 0; i < adjSampleSize; i++ {
		m := g.AdjectivePTOf(Masculine, Singular, Positive, AnyLengthWord)
		if strings.HasSuffix(m, "a") {
			t.Fatalf("Masculine positive request produced a feminine-only adjective %q", m)
		}
		f := g.AdjectivePTOf(Feminine, Singular, Positive, AnyLengthWord)
		if strings.HasSuffix(f, "o") {
			t.Fatalf("Feminine positive request produced a masculine-only adjective %q", f)
		}

		ms := g.AdjectivePTOf(Masculine, Singular, Superlative, AnyLengthWord)
		if !strings.HasSuffix(ms, "íssimo") {
			t.Fatalf("Masculine superlative request produced non-masculine %q", ms)
		}
		fs := g.AdjectivePTOf(Feminine, Singular, Superlative, AnyLengthWord)
		if !strings.HasSuffix(fs, "íssima") {
			t.Fatalf("Feminine superlative request produced non-feminine %q", fs)
		}
	}
}

// TestAdjectivePTNumberHonored verifies that number is honored across both degrees.
// Every pt-PT adjective plural in the inventory ends in -s (osos/osas, ais, áveis,
// entes, antes, íssimos/íssimas), and no singular form does, so the final -s is a
// sound discriminator.
func TestAdjectivePTNumberHonored(t *testing.T) {
	g := New(11)
	degrees := []Degree{Positive, Superlative}
	for _, de := range degrees {
		for i := 0; i < adjSampleSize; i++ {
			s := g.AdjectivePTOf(AnyGender, Singular, de, AnyLengthWord)
			if strings.HasSuffix(s, "s") {
				t.Fatalf("Singular request (degree %d) produced a plural-looking adjective %q", de, s)
			}
			p := g.AdjectivePTOf(AnyGender, Plural, de, AnyLengthWord)
			if !strings.HasSuffix(p, "s") {
				t.Fatalf("Plural request (degree %d) produced a singular-looking adjective %q", de, p)
			}
		}
	}
}

// TestAdjectivePTDegreeHonored verifies that degree is honored: a Positive request
// never yields a -íssimo form, and a Superlative request always does.
func TestAdjectivePTDegreeHonored(t *testing.T) {
	g := New(13)
	for i := 0; i < adjSampleSize; i++ {
		p := g.AdjectivePTOf(AnyGender, AnyNumber, Positive, AnyLengthWord)
		if isSuperlativeSurface(p) {
			t.Fatalf("Positive request produced a superlative %q", p)
		}
		s := g.AdjectivePTOf(AnyGender, AnyNumber, Superlative, AnyLengthWord)
		if !isSuperlativeSurface(s) {
			t.Fatalf("Superlative request produced a non-superlative %q", s)
		}
	}
}

// TestAdjectivePTSuperlativeInflection verifies that the synthetic absolute
// superlative inflects for gender and number through its own suffix: -íssimo,
// -íssima, -íssimos, -íssimas.
func TestAdjectivePTSuperlativeInflection(t *testing.T) {
	g := New(17)
	cases := []struct {
		gen    Gender
		num    Number
		suffix string
	}{
		{Masculine, Singular, "íssimo"},
		{Feminine, Singular, "íssima"},
		{Masculine, Plural, "íssimos"},
		{Feminine, Plural, "íssimas"},
	}
	for _, tc := range cases {
		for i := 0; i < adjSampleSize; i++ {
			w := g.AdjectivePTOf(tc.gen, tc.num, Superlative, AnyLengthWord)
			if !strings.HasSuffix(w, tc.suffix) {
				t.Fatalf("superlative(%d,%d) = %q, want suffix %q", tc.gen, tc.num, w, tc.suffix)
			}
		}
	}
}

// TestAdjectivePTSuperlativeOrthography verifies the deeper orthographic
// guarantees of the -íssimo superlative over a large sample: it carries exactly
// one acute/circumflex accent (the í of the suffix, the base having lost its own
// accent), it is a valid conformant word, and the c->qu adjustment is reachable
// through the -ico/-ica family (the riquíssimo pattern).
func TestAdjectivePTSuperlativeOrthography(t *testing.T) {
	g := New(19)
	sawHardC := false
	for i := 0; i < adjSampleSize; i++ {
		w := g.AdjectivePTOf(AnyGender, AnyNumber, Superlative, AnyLengthWord)
		assertAdjectiveOrthographyConformant(t, w)
		if got := countAcuteOrCircumflex(w); got != 1 {
			t.Fatalf("superlative %q has %d acute/circumflex accents, want exactly 1", w, got)
		}
		if !isSuperlativeSurface(w) {
			t.Fatalf("superlative %q does not end in an -íssimo form", w)
		}
		if strings.Contains(w, "quíss") {
			sawHardC = true
		}
	}
	if !sawHardC {
		t.Errorf("no c->qu superlative (quíss...) observed; the riquíssimo pattern is unreachable")
	}
}

// TestAdjectiveSuperlativeHardeningRule verifies the c->qu and g->gu orthographic
// adjustments directly on the superlative builder, using the specification's exact
// examples. It builds bare stems for rico, longo, banal and amável and checks that
// the superlative preserves the hard sound before the front í (riquíssimo,
// longuíssimo) and attaches directly to a consonant-final base (banalíssimo), and
// that a -vel base takes the regular (flagged) superlative (amavelíssimo).
func TestAdjectiveSuperlativeHardeningRule(t *testing.T) {
	cases := []struct {
		name string
		stem syllabicStem
		gen  Gender
		num  Number
		want string
	}{
		{
			name: "rico->riquíssimo (c->qu)",
			stem: syllabicStem{syllables: []syllable{{onset: "r", nucleus: "i"}, {onset: "c", nucleus: "o"}}, tonic: 0},
			gen:  Masculine, num: Singular, want: "riquíssimo",
		},
		{
			name: "longo->longuíssimo (g->gu)",
			stem: syllabicStem{syllables: []syllable{{onset: "l", nucleus: "o", coda: "n"}, {onset: "g", nucleus: "o"}}, tonic: 0},
			gen:  Masculine, num: Singular, want: "longuíssimo",
		},
		{
			name: "rica->riquíssima (c->qu, feminine)",
			stem: syllabicStem{syllables: []syllable{{onset: "r", nucleus: "i"}, {onset: "c", nucleus: "a"}}, tonic: 0},
			gen:  Feminine, num: Singular, want: "riquíssima",
		},
		{
			name: "longo->longuíssimos (g->gu, plural)",
			stem: syllabicStem{syllables: []syllable{{onset: "l", nucleus: "o", coda: "n"}, {onset: "g", nucleus: "o"}}, tonic: 0},
			gen:  Masculine, num: Plural, want: "longuíssimos",
		},
		{
			name: "banal->banalíssimo (consonant-final, direct)",
			stem: syllabicStem{syllables: []syllable{{onset: "b", nucleus: "a"}, {onset: "n", nucleus: "a", coda: "l"}}, tonic: 1},
			gen:  Masculine, num: Singular, want: "banalíssimo",
		},
		{
			name: "amável->amavelíssimo (regular -vel, flagged)",
			stem: syllabicStem{syllables: []syllable{{nucleus: "a"}, {onset: "m", nucleus: "a"}, {onset: "v", nucleus: "e", coda: "l"}}, tonic: 1},
			gen:  Masculine, num: Singular, want: "amavelíssimo",
		},
		{
			name: "famoso->famosíssima (drop vowel, feminine)",
			stem: syllabicStem{syllables: []syllable{{onset: "f", nucleus: "a"}, {onset: "m", nucleus: "o"}, {onset: "s", nucleus: "o"}}, tonic: 1},
			gen:  Feminine, num: Singular, want: "famosíssima",
		},
	}
	for _, tc := range cases {
		got := superlativeWord(tc.stem, tc.gen, tc.num)
		if got != tc.want {
			t.Errorf("%s: got %q, want %q", tc.name, got, tc.want)
		}
		assertAdjectiveOrthographyConformant(t, got)
		if got := countAcuteOrCircumflex(got); got != 1 {
			t.Errorf("%s: %q has %d acute/circumflex accents, want 1", tc.name, tc.want, got)
		}
	}
}

// TestAdjectivePTLengthHonored verifies that every output of a concrete
// length/degree/number request lands within the expected character window. No
// pt-PT adjective is shorter than five characters (the shortest ending, -al, plus
// the CV stem floor), so a Small positive request normalizes upward to the Medium
// window. Every superlative exceeds the Medium ceiling, so any Superlative request
// lands in the Big window. Length is measured in runes.
func TestAdjectivePTLengthHonored(t *testing.T) {
	g := New(29)
	cases := []struct {
		l        LengthTypeWords
		n        Number
		d        Degree
		lo, hi   int
		describe string
	}{
		{SmallLengthWord, Singular, Positive, 5, 8, "Small positive normalizes up to Medium"},
		{SmallLengthWord, Plural, Positive, 5, 8, "Small positive plural normalizes up to Medium"},
		{MediumLengthWords, Singular, Positive, 5, 8, "Medium positive fits Medium"},
		{MediumLengthWords, Plural, Positive, 5, 8, "Medium positive plural fits Medium"},
		{BigLengthWords, Singular, Positive, 9, 30, "Big positive fits Big"},
		{BigLengthWords, Plural, Positive, 9, 30, "Big positive plural fits Big"},
		{SmallLengthWord, Singular, Superlative, 9, 30, "Small superlative normalizes up to Big"},
		{MediumLengthWords, Singular, Superlative, 9, 30, "Medium superlative normalizes up to Big"},
		{BigLengthWords, Singular, Superlative, 9, 30, "Big superlative fits Big"},
		{BigLengthWords, Plural, Superlative, 9, 30, "Big superlative plural fits Big"},
	}
	for _, tc := range cases {
		for i := 0; i < adjSampleSize; i++ {
			w := g.AdjectivePTOf(AnyGender, tc.n, tc.d, tc.l)
			n := utf8.RuneCountInString(w)
			if n < tc.lo || n > tc.hi {
				t.Fatalf("%s: adjective %q has %d runes, outside window [%d,%d]",
					tc.describe, w, n, tc.lo, tc.hi)
			}
		}
	}
}

// TestAdjectivePTAnyGenderVaried verifies that AnyGender yields both masculine-
// and feminine-marked positive adjectives over a sampled run (the invariable
// endings, which mark neither, are ignored by this check).
func TestAdjectivePTAnyGenderVaried(t *testing.T) {
	g := New(101)
	sawMasc, sawFem := false, false
	for i := 0; i < adjSampleSize; i++ {
		w := g.AdjectivePTOf(AnyGender, Singular, Positive, AnyLengthWord)
		switch {
		case strings.HasSuffix(w, "o"):
			sawMasc = true
		case strings.HasSuffix(w, "a"):
			sawFem = true
		}
	}
	if !sawMasc || !sawFem {
		t.Errorf("Any gender not varied: masc=%v fem=%v", sawMasc, sawFem)
	}
}

// TestAdjectivePTAnyNumberVaried verifies that AnyNumber yields both singular and
// plural adjectives over a sampled run.
func TestAdjectivePTAnyNumberVaried(t *testing.T) {
	g := New(202)
	sawSing, sawPlur := false, false
	for i := 0; i < adjSampleSize; i++ {
		if strings.HasSuffix(g.AdjectivePTOf(AnyGender, AnyNumber, Positive, AnyLengthWord), "s") {
			sawPlur = true
		} else {
			sawSing = true
		}
	}
	if !sawSing || !sawPlur {
		t.Errorf("Any number not varied: singular=%v plural=%v", sawSing, sawPlur)
	}
}

// TestAdjectivePTAnyDegreeVaried verifies that AnyDegree yields both positive and
// superlative adjectives over a sampled run.
func TestAdjectivePTAnyDegreeVaried(t *testing.T) {
	g := New(303)
	sawPos, sawSuper := false, false
	for i := 0; i < adjSampleSize; i++ {
		if isSuperlativeSurface(g.AdjectivePTOf(AnyGender, AnyNumber, AnyDegree, AnyLengthWord)) {
			sawSuper = true
		} else {
			sawPos = true
		}
	}
	if !sawPos || !sawSuper {
		t.Errorf("Any degree not varied: positive=%v superlative=%v", sawPos, sawSuper)
	}
}

// TestAdjectivePTAnyLengthVaried verifies that AnyLengthWord yields positive
// adjectives spread across the length buckets, in particular the medium (5-8) and
// big (9+) ranges. The structural minimum adjective is four characters (a
// bare-vowel leading syllable plus the -al ending, for example "egal"): no pt-PT
// adjective in the inventory is shorter than that, so the smallest words are
// four-character -al forms and no output is ever three characters or fewer.
func TestAdjectivePTAnyLengthVaried(t *testing.T) {
	g := New(404)
	sawMedium, sawBig := false, false
	for i := 0; i < adjSampleSize; i++ {
		n := utf8.RuneCountInString(g.AdjectivePTOf(AnyGender, Singular, Positive, AnyLengthWord))
		switch {
		case n < 4:
			t.Fatalf("observed an adjective shorter than the four-character structural floor: %d runes", n)
		case n <= 8:
			sawMedium = true
		default:
			sawBig = true
		}
	}
	if !sawMedium || !sawBig {
		t.Errorf("Any length not varied: medium=%v big=%v", sawMedium, sawBig)
	}
}

// TestAdjectivePTEndingVariety verifies that, over a broad (Any length) sample,
// every adjective ending is reachable for each gender, confirming the ending
// inventory is fully wired and the weighted selection reaches all forms.
func TestAdjectivePTEndingVariety(t *testing.T) {
	g := New(505)
	mascSuffixes := []string{"oso", "ico", "ivo", "al", "ável", "ível", "ente", "ante"}
	femSuffixes := []string{"osa", "ica", "iva", "al", "ável", "ível", "ente", "ante"}
	seenMasc := map[string]bool{}
	seenFem := map[string]bool{}
	for i := 0; i < adjSampleSize; i++ {
		m := g.AdjectivePTOf(Masculine, Singular, Positive, AnyLengthWord)
		for _, s := range mascSuffixes {
			if strings.HasSuffix(m, s) {
				seenMasc[s] = true
				break
			}
		}
		f := g.AdjectivePTOf(Feminine, Singular, Positive, AnyLengthWord)
		for _, s := range femSuffixes {
			if strings.HasSuffix(f, s) {
				seenFem[s] = true
				break
			}
		}
	}
	for _, s := range mascSuffixes {
		if !seenMasc[s] {
			t.Errorf("masculine ending %q never observed", s)
		}
	}
	for _, s := range femSuffixes {
		if !seenFem[s] {
			t.Errorf("feminine ending %q never observed", s)
		}
	}
}

// TestAdjectivePTReproducible verifies that a seeded Generator is deterministic:
// two generators created with the same seed produce identical adjective sequences,
// for both the sugar method and the fully specified method.
func TestAdjectivePTReproducible(t *testing.T) {
	a := New(0xC0FFEE)
	b := New(0xC0FFEE)
	for i := 0; i < adjSampleSize; i++ {
		if x, y := a.AdjectivePT(), b.AdjectivePT(); x != y {
			t.Fatalf("AdjectivePT diverged at %d: %q != %q", i, x, y)
		}
	}

	genders := []Gender{AnyGender, Masculine, Feminine}
	numbers := []Number{AnyNumber, Singular, Plural}
	degrees := []Degree{AnyDegree, Positive, Superlative}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	c := New(42)
	d := New(42)
	for _, ge := range genders {
		for _, nu := range numbers {
			for _, de := range degrees {
				for _, le := range lengths {
					for i := 0; i < 100; i++ {
						if x, y := c.AdjectivePTOf(ge, nu, de, le), d.AdjectivePTOf(ge, nu, de, le); x != y {
							t.Fatalf("AdjectivePTOf(%d,%d,%d,%d) diverged at %d: %q != %q",
								ge, nu, de, le, i, x, y)
						}
					}
				}
			}
		}
	}
}

// TestAdjectiveEndingLengthAwareSelection verifies the length-aware ending
// selection: when a category is too short for any adjective ending (Small, whose
// ceiling of four characters is below every ending's minimum viable length), the
// selector falls back to the shortest ending (-al), and the core then normalizes
// the category upward; with a Big ceiling every ending is reachable.
func TestAdjectiveEndingLengthAwareSelection(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))

	// No adjective ending fits the Small ceiling of 4 characters, so the fallback
	// is always the shortest ending, -al.
	smallMasc := map[string]bool{}
	for i := 0; i < 5000; i++ {
		smallMasc[selectFittingEnding(r, masculineAdjectiveEndings, 4).label] = true
	}
	if len(smallMasc) != 1 || !smallMasc["-al"] {
		t.Errorf("Small masculine selection = %v, want only -al", smallMasc)
	}

	// Every masculine ending must be reachable once the ceiling is the Big range.
	bigMasc := map[string]bool{}
	for i := 0; i < 5000; i++ {
		bigMasc[selectFittingEnding(r, masculineAdjectiveEndings, 30).label] = true
	}
	for _, want := range []string{"-oso", "-ico", "-ivo", "-al", "-ável", "-ível", "-ente", "-ante"} {
		if !bigMasc[want] {
			t.Errorf("Big masculine selection missing ending %q (got %v)", want, bigMasc)
		}
	}
}

// TestAdjectiveStemConformsToOracles builds combined adjective stems for every
// ending and asserts, over a sampled run, that each positive stem satisfies the
// Sprint 7 decomposition oracles: the phonotactic oracle (validating the
// leading/ending junction) and, after applying the accent, the accentuation oracle
// (validating that the ending's stress — including the -ico proparoxytone — lands
// the accent correctly on the combined word). This is stronger than the string
// oracles, which do not inspect syllable structure or accent placement.
func TestAdjectiveStemConformsToOracles(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5EED, 0xF00D))
	endings := []*nounEnding{
		&endAdjMascOso, &endAdjFemOsa, &endAdjMascIco, &endAdjFemIca,
		&endAdjMascIvo, &endAdjFemIva, &endAdjAl, &endAdjAvel, &endAdjIvel,
		&endAdjEnte, &endAdjAnte,
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

// TestAdjectivePTConcurrentSafe exercises the package-level AdjectivePT from many
// goroutines. Run with -race, it demonstrates that the package-level surface is
// safe for concurrent use, as documented: its shared global source delegates each
// draw to the concurrency-safe global math/rand/v2 generator.
func TestAdjectivePTConcurrentSafe(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 4000
	done := make(chan struct{}, goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for i := 0; i < perGoroutine; i++ {
				if AdjectivePT() == "" {
					t.Error("AdjectivePT returned an empty string")
					return
				}
			}
		}()
	}
	for g := 0; g < goroutines; g++ {
		<-done
	}
}

// BenchmarkAdjectivePTOfPositiveSingular measures the positive singular path and
// its allocation budget: a positive singular adjective must assemble in exactly
// one string allocation.
func BenchmarkAdjectivePTOfPositiveSingular(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.AdjectivePTOf(AnyGender, Singular, Positive, AnyLengthWord)
	}
}

// BenchmarkAdjectivePTOfPositivePlural measures the positive plural path, which
// adds one tail concatenation in the shared pluralizer over the singular assembly.
func BenchmarkAdjectivePTOfPositivePlural(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.AdjectivePTOf(AnyGender, Plural, Positive, AnyLengthWord)
	}
}

// BenchmarkAdjectivePTOfSuperlativeSingular measures the -íssimo superlative path,
// which is a single allocation: the word is written straight from the bare stem
// into one pre-sized builder, with no separate base assembly.
func BenchmarkAdjectivePTOfSuperlativeSingular(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.AdjectivePTOf(AnyGender, Singular, Superlative, AnyLengthWord)
	}
}

// BenchmarkAdjectivePT measures the default package-level adjective path over the
// global source (random gender, number, degree and length).
func BenchmarkAdjectivePT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = AdjectivePT()
	}
}
