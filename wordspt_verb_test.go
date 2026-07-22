package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the WordsPT verb-conjugation foundation (task #28): the
// public Mood/Tense/Person enums, the regular-conjugation machinery, the
// radical/thematic-vowel/desinence assembly, and the NON-FINITE forms (impersonal
// infinitive, inflected personal infinitive, gerund, participle) of all three
// conjugations. The exact-form checks are the primary acceptance criterion: they
// build every non-finite slot from the model radicals fal/com/part and assert the
// textbook form, so the desinences are proven against the Cunha & Cintra paradigm.

// TestNonFiniteExactForms is the core acceptance test: for each conjugation, built
// from its model radical (fal/com/part), every non-finite form must equal the
// textbook value exactly. This proves the thematic vowels (including the
// participle's a/i/i split: falado, comido, partido) and every desinence against
// the authoritative paradigm.
func TestNonFiniteExactForms(t *testing.T) {
	cases := []struct {
		conj    *verbConjugation
		radical string
		// personal infinitive, indexed by verbPerson (1s, 2s, 3s, 1p, 3p)
		personalInf [5]string
		impersonal  string
		gerund      string
		participle  string
	}{
		{
			conj: &conjAr, radical: "fal",
			impersonal:  "falar",
			personalInf: [5]string{"falar", "falares", "falar", "falarmos", "falarem"},
			gerund:      "falando",
			participle:  "falado",
		},
		{
			conj: &conjEr, radical: "com",
			impersonal:  "comer",
			personalInf: [5]string{"comer", "comeres", "comer", "comermos", "comerem"},
			gerund:      "comendo",
			participle:  "comido",
		},
		{
			conj: &conjIr, radical: "part",
			impersonal:  "partir",
			personalInf: [5]string{"partir", "partires", "partir", "partirmos", "partirem"},
			gerund:      "partindo",
			participle:  "partido",
		},
	}

	persons := [5]verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, tc := range cases {
		if got := conjugateNonFinite(tc.radical, tc.conj, formImpersonalInfinitive, verbP1s); got != tc.impersonal {
			t.Errorf("%s impersonal infinitive = %q, want %q", tc.conj.label, got, tc.impersonal)
		}
		for i, vp := range persons {
			if got := conjugateNonFinite(tc.radical, tc.conj, formPersonalInfinitive, vp); got != tc.personalInf[i] {
				t.Errorf("%s personal infinitive[%d] = %q, want %q", tc.conj.label, vp, got, tc.personalInf[i])
			}
		}
		if got := conjugateNonFinite(tc.radical, tc.conj, formGerund, verbP1s); got != tc.gerund {
			t.Errorf("%s gerund = %q, want %q", tc.conj.label, got, tc.gerund)
		}
		if got := conjugateNonFinite(tc.radical, tc.conj, formParticiple, verbP1s); got != tc.participle {
			t.Errorf("%s participle = %q, want %q", tc.conj.label, got, tc.participle)
		}
	}
}

// TestPersonalInfinitiveIgnoresPersonWhereApplicable confirms that the person is
// consumed only by the personal infinitive: the impersonal infinitive, the gerund
// and the participle produce the same form for every person, per the paradigm
// (they are invariable for person).
func TestPersonalInfinitiveIgnoresPersonWhereApplicable(t *testing.T) {
	invariant := []nonFiniteForm{formImpersonalInfinitive, formGerund, formParticiple}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, form := range invariant {
		want := conjugateNonFinite("fal", &conjAr, form, verbP1s)
		for _, vp := range persons {
			if got := conjugateNonFinite("fal", &conjAr, form, vp); got != want {
				t.Errorf("form %d not invariant for person: person %d gave %q, want %q", form, vp, got, want)
			}
		}
	}
}

// assertVerbOrthographyConformant fails t when word is not an orthographically
// well-formed, lowercase, valid-UTF-8 pt-PT word, per the Sprint 7 string oracles
// and the case/encoding rules. It is the string-oracle check the task requires for
// the non-finite forms (which are assembled with an ASCII desinence, like the
// plural and the -mente adverb).
func assertVerbOrthographyConformant(t *testing.T, word string) {
	t.Helper()
	if word == "" {
		t.Fatalf("generated an empty verb")
	}
	if !utf8.ValidString(word) {
		t.Fatalf("verb %q is not valid UTF-8", word)
	}
	if !wordConformsToCedilla(word) {
		t.Errorf("verb %q violates the pt-PT cedilla rule", word)
	}
	if wordHasForbiddenDiaeresis(word) {
		t.Errorf("verb %q contains a forbidden diaeresis", word)
	}
	if wordContainsGrave(word) {
		t.Errorf("verb %q contains a grave accent", word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("verb %q contains an uppercase rune %q", word, r)
		}
	}
}

// verbNonFiniteEndings maps a non-finite form to the set of surface endings a
// generated form of that form may carry, across all three conjugations. It is a
// paradigm oracle: it ties the generated core output to the desinence tables, a
// stronger check than the string oracles alone.
func hasAnySuffix(word string, suffixes ...string) bool {
	for _, s := range suffixes {
		if strings.HasSuffix(word, s) {
			return true
		}
	}
	return false
}

// TestNonFiniteVerbOrthographicConformance samples the non-finite core over every
// form, person, number and length and asserts that each output is orthographically
// conformant, lowercase, valid UTF-8, and carries a desinence from the paradigm.
func TestNonFiniteVerbOrthographicConformance(t *testing.T) {
	r := rand.New(rand.NewPCG(0xB0B, 0xCAFE))
	var buf [nounMaxSyllables]syllable

	forms := []nonFiniteForm{
		formImpersonalInfinitive, formPersonalInfinitive, formGerund, formParticiple,
	}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}

	for _, form := range forms {
		for _, p := range persons {
			for _, n := range numbers {
				for _, l := range lengths {
					for i := 0; i < 200; i++ {
						w := nonFiniteVerbCore(r, buf[:0], form, p, n, l)
						assertVerbOrthographyConformant(t, w)
						assertNonFiniteEnding(t, form, w)
					}
				}
			}
		}
	}
}

// assertNonFiniteEnding checks that a generated non-finite form ends with a
// desinence the paradigm permits for its form, across all three conjugations.
func assertNonFiniteEnding(t *testing.T, form nonFiniteForm, word string) {
	t.Helper()
	switch form {
	case formImpersonalInfinitive:
		if !hasAnySuffix(word, "ar", "er", "ir") {
			t.Errorf("impersonal infinitive %q does not end in -ar/-er/-ir", word)
		}
	case formPersonalInfinitive:
		// -ar/-er/-ir (1s,3s), -res (2s), -rmos (1p), -rem (3p), across conjugations.
		if !hasAnySuffix(word,
			"ar", "er", "ir",
			"ares", "eres", "ires",
			"armos", "ermos", "irmos",
			"arem", "erem", "irem") {
			t.Errorf("personal infinitive %q has no valid personal ending", word)
		}
	case formGerund:
		if !hasAnySuffix(word, "ando", "endo", "indo") {
			t.Errorf("gerund %q does not end in -ando/-endo/-indo", word)
		}
	case formParticiple:
		if !hasAnySuffix(word, "ado", "ido") {
			t.Errorf("participle %q does not end in -ado/-ido", word)
		}
	}
}

// TestNonFiniteVerbLengthHonored verifies that every output of a concrete length
// category lands within the expected character window. The shortest non-finite
// verb exceeds four characters (a CV stem plus a sampled consonant, the thematic
// vowel and the desinence), so a Small request normalizes upward to Medium, per
// the minimum-viable-length rule; Medium and Big fit their category. Length is
// measured in runes.
func TestNonFiniteVerbLengthHonored(t *testing.T) {
	r := rand.New(rand.NewPCG(29, 31))
	var buf [nounMaxSyllables]syllable

	forms := []nonFiniteForm{
		formImpersonalInfinitive, formPersonalInfinitive, formGerund, formParticiple,
	}
	cases := []struct {
		l      LengthTypeWords
		lo, hi int
	}{
		{SmallLengthWord, 5, 8}, // normalizes up to Medium
		{MediumLengthWords, 5, 8},
		{BigLengthWords, 9, 30},
	}
	for _, form := range forms {
		for _, tc := range cases {
			for i := 0; i < 4000; i++ {
				w := nonFiniteVerbCore(r, buf[:0], form, First, Plural, tc.l)
				n := utf8.RuneCountInString(w)
				if n < tc.lo || n > tc.hi {
					t.Fatalf("form %d length %d: verb %q has %d runes, outside window [%d,%d]",
						form, tc.l, w, n, tc.lo, tc.hi)
				}
			}
		}
	}
}

// TestNonFiniteVerbReproducible verifies that the shared core is deterministic: two
// sources created with the same seed produce identical verb sequences for the same
// arguments.
func TestNonFiniteVerbReproducible(t *testing.T) {
	forms := []nonFiniteForm{
		formImpersonalInfinitive, formPersonalInfinitive, formGerund, formParticiple,
	}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}

	a := rand.New(rand.NewPCG(42, 42))
	b := rand.New(rand.NewPCG(42, 42))
	var bufA, bufB [nounMaxSyllables]syllable
	for _, form := range forms {
		for _, p := range persons {
			for _, n := range numbers {
				for _, l := range lengths {
					for i := 0; i < 50; i++ {
						x := nonFiniteVerbCore(a, bufA[:0], form, p, n, l)
						y := nonFiniteVerbCore(b, bufB[:0], form, p, n, l)
						if x != y {
							t.Fatalf("core diverged (form %d, p %d, n %d, l %d) at %d: %q != %q",
								form, p, n, l, i, x, y)
						}
					}
				}
			}
		}
	}
}

// TestPickConjugationWeighted verifies that the weighted conjugation draw reaches
// all three conjugations and that the first conjugation (-ar) dominates, as pt-PT
// requires.
func TestPickConjugationWeighted(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 9))
	counts := map[string]int{}
	const n = 60000
	for i := 0; i < n; i++ {
		counts[pickConjugation(r).label]++
	}
	for _, label := range []string{"-ar", "-er", "-ir"} {
		if counts[label] == 0 {
			t.Errorf("conjugation %q never drawn", label)
		}
	}
	if counts["-ar"] <= counts["-er"] || counts["-ar"] <= counts["-ir"] {
		t.Errorf("first conjugation not dominant: -ar=%d -er=%d -ir=%d",
			counts["-ar"], counts["-er"], counts["-ir"])
	}
}

// TestResolveVerbPerson verifies the person/number resolution, including the N1
// normalization (second person plural -> third person plural, no vós) and the
// Any-resolution of an unspecified person and number.
func TestResolveVerbPerson(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))

	// Concrete mappings.
	concrete := []struct {
		p    Person
		n    Number
		want verbPerson
	}{
		{First, Singular, verbP1s},
		{Second, Singular, verbP2s},
		{Third, Singular, verbP3s},
		{First, Plural, verbP1p},
		{Third, Plural, verbP3p},
		{Second, Plural, verbP3p}, // N1: vós -> vocês (third plural)
	}
	for _, tc := range concrete {
		if got := resolveVerbPerson(r, tc.p, tc.n); got != tc.want {
			t.Errorf("resolveVerbPerson(%d,%d) = %d, want %d", tc.p, tc.n, got, tc.want)
		}
	}

	// N1 holds regardless of the resolved number when the second person plural is
	// requested explicitly: it must never yield a second-plural slot (there is none).
	for i := 0; i < 10000; i++ {
		if got := resolveVerbPerson(r, Second, Plural); got != verbP3p {
			t.Fatalf("Second+Plural resolved to %d, want verbP3p (N1)", got)
		}
	}

	// Any person and any number reach every one of the five slots.
	seen := map[verbPerson]bool{}
	for i := 0; i < 20000; i++ {
		seen[resolveVerbPerson(r, AnyPerson, AnyNumber)] = true
	}
	for _, vp := range []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p} {
		if !seen[vp] {
			t.Errorf("Any person/number never reached slot %d", vp)
		}
	}
}

// TestVerbRadicalStemConformsToOracle builds combined verb stems (leading stem plus
// the sampled thematic-vowel syllable) for every thematic-vowel ending and asserts,
// over a sampled run, that each stem satisfies the Sprint 7 phonotactic oracle. This
// validates the radical/thematic-vowel junction the verb machinery introduces
// (a consonant onset always precedes the thematic vowel, so there is no hiatus).
func TestVerbRadicalStemConformsToOracle(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5EED, 0xF00D))
	endings := []*nounEnding{&themeEndingA, &themeEndingE, &themeEndingI}
	var buf []syllable
	for _, end := range endings {
		for i := 0; i < 5000; i++ {
			leading := 1 + int(r.Uint32N(6))
			stem := sampleNounStem(r, buf, leading, end)
			buf = stem.syllables
			if !stemConformsToPhonotactics(stem.syllables) {
				t.Fatalf("theme ending %s: stem %q fails phonotactic oracle",
					end.label, assembleStem(stem.syllables))
			}
			// The last syllable must carry the thematic vowel with a real consonant
			// onset (never empty), so the radical never ends in a vowel hiatus.
			last := stem.syllables[len(stem.syllables)-1]
			if last.onset == "" {
				t.Fatalf("theme ending %s: thematic syllable has empty onset in %q",
					end.label, assembleStem(stem.syllables))
			}
		}
	}
}

// TestConjugateNonFiniteSingleAllocation asserts that the paradigm-exact assembler
// performs exactly one allocation (the returned string).
func TestConjugateNonFiniteSingleAllocation(t *testing.T) {
	forms := []nonFiniteForm{
		formImpersonalInfinitive, formPersonalInfinitive, formGerund, formParticiple,
	}
	for _, form := range forms {
		allocs := testing.AllocsPerRun(1000, func() {
			_ = conjugateNonFinite("part", &conjIr, form, verbP3p)
		})
		if allocs != 1 {
			t.Errorf("conjugateNonFinite (form %d) allocated %.0f times, want 1", form, allocs)
		}
	}
}

// BenchmarkNonFiniteVerbInfinitive measures the impersonal-infinitive path and its
// allocation budget: a non-finite verb must assemble in exactly one string
// allocation.
func BenchmarkNonFiniteVerbInfinitive(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	var buf [nounMaxSyllables]syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = nonFiniteVerbCore(r, buf[:0], formImpersonalInfinitive, AnyPerson, AnyNumber, AnyLengthWord)
	}
}

// BenchmarkNonFiniteVerbPersonalInfinitive measures the inflected personal
// infinitive path (which resolves a person/number slot).
func BenchmarkNonFiniteVerbPersonalInfinitive(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	var buf [nounMaxSyllables]syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = nonFiniteVerbCore(r, buf[:0], formPersonalInfinitive, AnyPerson, AnyNumber, AnyLengthWord)
	}
}

// BenchmarkNonFiniteVerbGerund measures the gerund path.
func BenchmarkNonFiniteVerbGerund(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	var buf [nounMaxSyllables]syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = nonFiniteVerbCore(r, buf[:0], formGerund, AnyPerson, AnyNumber, AnyLengthWord)
	}
}

// BenchmarkNonFiniteVerbParticiple measures the participle path.
func BenchmarkNonFiniteVerbParticiple(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	var buf [nounMaxSyllables]syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = nonFiniteVerbCore(r, buf[:0], formParticiple, AnyPerson, AnyNumber, AnyLengthWord)
	}
}
