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

// ---------------------------------------------------------------------------
// Indicative mood (task #29)
// ---------------------------------------------------------------------------

// tenseName returns a readable label for an indicative Tense, for test failure
// messages only.
func tenseName(t Tense) string {
	switch t {
	case Present:
		return "present"
	case Imperfect:
		return "imperfect"
	case Preterite:
		return "preterite"
	case Future:
		return "future"
	case Conditional:
		return "conditional"
	default:
		return "any"
	}
}

// TestIndicativeExactForms is the core acceptance test for the indicative mood.
// For each conjugation, built from its model radical (fal/com/part), every
// indicative cell — 5 tenses x 5 persons per conjugation, 75 in total — must equal
// the textbook value exactly. This proves every desinence, and every baked-in
// graphic accent, against the Cunha & Cintra paradigm, including the pt-PT
// falámos (preterite 1p) versus falamos (present 1p) acute distinction, which is
// also asserted explicitly below.
func TestIndicativeExactForms(t *testing.T) {
	type tenseRow struct {
		tense Tense
		forms [5]string // indexed 1s, 2s, 3s, 1p, 3p
	}
	cases := []struct {
		conj    *verbConjugation
		radical string
		rows    []tenseRow
	}{
		{
			conj: &conjAr, radical: "fal",
			rows: []tenseRow{
				{Present, [5]string{"falo", "falas", "fala", "falamos", "falam"}},
				{Imperfect, [5]string{"falava", "falavas", "falava", "falávamos", "falavam"}},
				{Preterite, [5]string{"falei", "falaste", "falou", "falámos", "falaram"}},
				{Future, [5]string{"falarei", "falarás", "falará", "falaremos", "falarão"}},
				{Conditional, [5]string{"falaria", "falarias", "falaria", "falaríamos", "falariam"}},
			},
		},
		{
			conj: &conjEr, radical: "com",
			rows: []tenseRow{
				{Present, [5]string{"como", "comes", "come", "comemos", "comem"}},
				{Imperfect, [5]string{"comia", "comias", "comia", "comíamos", "comiam"}},
				{Preterite, [5]string{"comi", "comeste", "comeu", "comemos", "comeram"}},
				{Future, [5]string{"comerei", "comerás", "comerá", "comeremos", "comerão"}},
				{Conditional, [5]string{"comeria", "comerias", "comeria", "comeríamos", "comeriam"}},
			},
		},
		{
			conj: &conjIr, radical: "part",
			rows: []tenseRow{
				{Present, [5]string{"parto", "partes", "parte", "partimos", "partem"}},
				{Imperfect, [5]string{"partia", "partias", "partia", "partíamos", "partiam"}},
				{Preterite, [5]string{"parti", "partiste", "partiu", "partimos", "partiram"}},
				{Future, [5]string{"partirei", "partirás", "partirá", "partiremos", "partirão"}},
				{Conditional, [5]string{"partiria", "partirias", "partiria", "partiríamos", "partiriam"}},
			},
		},
	}

	persons := [5]verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	asserted := 0
	for _, tc := range cases {
		for _, row := range tc.rows {
			for i, vp := range persons {
				got := conjugateIndicative(tc.radical, tc.conj, row.tense, vp)
				want := row.forms[i]
				if got != want {
					t.Errorf("%s %s person %d = %q, want %q",
						tc.conj.label, tenseName(row.tense), vp, got, want)
				}
				asserted++
			}
		}
	}
	if asserted != 75 {
		t.Fatalf("asserted %d indicative cells, want 75 (3 conjugations x 5 tenses x 5 persons)", asserted)
	}
}

// TestIndicativePreteriteVsPresentFirstPlural pins the pt-PT distinction between
// the first-conjugation present first person plural falamos (no accent) and the
// preterite first person plural falámos (acute on the theme vowel): they differ
// only by that accent, and both must be produced as shown.
func TestIndicativePreteriteVsPresentFirstPlural(t *testing.T) {
	present := conjugateIndicative("fal", &conjAr, Present, verbP1p)
	preterite := conjugateIndicative("fal", &conjAr, Preterite, verbP1p)
	if present != "falamos" {
		t.Errorf("present 1p = %q, want %q", present, "falamos")
	}
	if preterite != "falámos" {
		t.Errorf("preterite 1p = %q, want %q", preterite, "falámos")
	}
	if present == preterite {
		t.Errorf("present and preterite 1p must differ (falamos vs falámos), both = %q", present)
	}
}

// TestConjugateIndicativeSingleAllocation asserts that the paradigm-exact
// indicative assembler performs exactly one allocation (the returned string) for
// every tense and person.
func TestConjugateIndicativeSingleAllocation(t *testing.T) {
	tenses := []Tense{Present, Imperfect, Preterite, Future, Conditional}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, tt := range tenses {
		for _, vp := range persons {
			allocs := testing.AllocsPerRun(1000, func() {
				_ = conjugateIndicative("part", &conjIr, tt, vp)
			})
			if allocs != 1 {
				t.Errorf("conjugateIndicative(%s, person %d) allocated %.0f times, want 1",
					tenseName(tt), vp, allocs)
			}
		}
	}
}

// TestIndicativeOrthographicConformance asserts that every indicative cell of the
// three model conjugations is an orthographically well-formed, lowercase,
// valid-UTF-8 pt-PT word (the Sprint 7 string oracles). This certifies the
// baked-in accents: only the acute (á, í) and the nasal tilde (ã of -ão) appear,
// and no cedilla, diaeresis or grave is ever produced.
func TestIndicativeOrthographicConformance(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{Present, Imperfect, Preterite, Future, Conditional}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, vp := range persons {
				w := conjugateIndicative(radicals[c], c, tt, vp)
				assertVerbOrthographyConformant(t, w)
				if !strings.HasPrefix(w, radicals[c]) {
					t.Errorf("%s %s person %d = %q does not start with radical %q",
						c.label, tenseName(tt), vp, w, radicals[c])
				}
			}
		}
	}
}

// TestResolveIndicativeTense verifies that a concrete tense passes through
// unchanged and that AnyTense reaches all five indicative tenses.
func TestResolveIndicativeTense(t *testing.T) {
	r := rand.New(rand.NewPCG(5, 7))
	for _, tt := range []Tense{Present, Imperfect, Preterite, Future, Conditional} {
		if got := resolveIndicativeTense(r, tt); got != tt {
			t.Errorf("resolveIndicativeTense(%s) = %s, want passthrough", tenseName(tt), tenseName(got))
		}
	}
	seen := map[Tense]bool{}
	for i := 0; i < 20000; i++ {
		seen[resolveIndicativeTense(r, AnyTense)] = true
	}
	for _, tt := range []Tense{Present, Imperfect, Preterite, Future, Conditional} {
		if !seen[tt] {
			t.Errorf("AnyTense never resolved to %s", tenseName(tt))
		}
	}
}

// TestIndicativeVerbFormConformance exercises the injectable-*rand.Rand assembler
// over every conjugation, tense (including AnyTense), person and number, asserting
// that each form is orthographically conformant and starts with the supplied
// radical. It validates the option-resolution path (tense, person/number with the
// N1 normalization) on top of the exact desinences.
func TestIndicativeVerbFormConformance(t *testing.T) {
	r := rand.New(rand.NewPCG(0x1D, 0x1CA))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{AnyTense, Present, Imperfect, Preterite, Future, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 200; i++ {
						w := indicativeVerbForm(r, radicals[c], c, tt, p, n)
						assertVerbOrthographyConformant(t, w)
						if !strings.HasPrefix(w, radicals[c]) {
							t.Fatalf("%s form %q does not start with radical %q", c.label, w, radicals[c])
						}
					}
				}
			}
		}
	}
}

// TestIndicativeVerbFormReproducible verifies that the injectable-*rand.Rand
// assembler is deterministic: two sources with the same seed produce identical
// indicative sequences for the same arguments.
func TestIndicativeVerbFormReproducible(t *testing.T) {
	a := rand.New(rand.NewPCG(99, 99))
	b := rand.New(rand.NewPCG(99, 99))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{AnyTense, Present, Imperfect, Preterite, Future, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 50; i++ {
						x := indicativeVerbForm(a, radicals[c], c, tt, p, n)
						y := indicativeVerbForm(b, radicals[c], c, tt, p, n)
						if x != y {
							t.Fatalf("indicativeVerbForm diverged (%s, %s, p %d, n %d) at %d: %q != %q",
								c.label, tenseName(tt), p, n, i, x, y)
						}
					}
				}
			}
		}
	}
}

// BenchmarkConjugateIndicative measures the paradigm-exact indicative assembler
// and its allocation budget: an indicative form must assemble in exactly one
// string allocation.
func BenchmarkConjugateIndicative(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = conjugateIndicative("fal", &conjAr, Preterite, verbP1p)
	}
}

// BenchmarkIndicativeVerbForm measures the injectable-*rand.Rand assembler,
// including the tense and person/number resolution, for a fully unspecified
// request.
func BenchmarkIndicativeVerbForm(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = indicativeVerbForm(r, "fal", &conjAr, AnyTense, AnyPerson, AnyNumber)
	}
}

// ---------------------------------------------------------------------------
// Subjunctive and imperative moods (task #30)
// ---------------------------------------------------------------------------

// The five imperative slots reduce to the four the mood actually has (no first
// person singular); this list drives the imperative exact-form and conformance
// tests.
var imperativePersons = [4]verbPerson{verbP2s, verbP3s, verbP1p, verbP3p}

// subjunctiveTenseName returns a readable label for a subjunctive Tense, for test
// failure messages only.
func subjunctiveTenseName(t Tense) string {
	switch t {
	case Present:
		return "present"
	case Imperfect:
		return "imperfect"
	case Future:
		return "future"
	default:
		return "any"
	}
}

// polarityName returns a readable label for an imperativePolarity, for test
// failure messages only.
func polarityName(pol imperativePolarity) string {
	if pol == imperativeNeg {
		return "negative"
	}
	return "affirmative"
}

// TestSubjunctiveExactForms is the core acceptance test for the subjunctive mood.
// For each conjugation, built from its model radical (fal/com/part), every
// subjunctive cell — 3 tenses x 5 persons per conjugation, 45 in total — must
// equal the textbook value exactly. This proves every desinence, and every
// baked-in graphic accent, against the Cunha & Cintra paradigm: the theme-vowel
// swap of the present (fale vs coma/parta), the proparoxytone imperfect first
// plural (falássemos, comêssemos, partíssemos) and the personal-infinitive future
// (falarmos, comeres, partir).
func TestSubjunctiveExactForms(t *testing.T) {
	type tenseRow struct {
		tense Tense
		forms [5]string // indexed 1s, 2s, 3s, 1p, 3p
	}
	cases := []struct {
		conj    *verbConjugation
		radical string
		rows    []tenseRow
	}{
		{
			conj: &conjAr, radical: "fal",
			rows: []tenseRow{
				{Present, [5]string{"fale", "fales", "fale", "falemos", "falem"}},
				{Imperfect, [5]string{"falasse", "falasses", "falasse", "falássemos", "falassem"}},
				{Future, [5]string{"falar", "falares", "falar", "falarmos", "falarem"}},
			},
		},
		{
			conj: &conjEr, radical: "com",
			rows: []tenseRow{
				{Present, [5]string{"coma", "comas", "coma", "comamos", "comam"}},
				{Imperfect, [5]string{"comesse", "comesses", "comesse", "comêssemos", "comessem"}},
				{Future, [5]string{"comer", "comeres", "comer", "comermos", "comerem"}},
			},
		},
		{
			conj: &conjIr, radical: "part",
			rows: []tenseRow{
				{Present, [5]string{"parta", "partas", "parta", "partamos", "partam"}},
				{Imperfect, [5]string{"partisse", "partisses", "partisse", "partíssemos", "partissem"}},
				{Future, [5]string{"partir", "partires", "partir", "partirmos", "partirem"}},
			},
		},
	}

	persons := [5]verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	asserted := 0
	for _, tc := range cases {
		for _, row := range tc.rows {
			for i, vp := range persons {
				got := conjugateSubjunctive(tc.radical, tc.conj, row.tense, vp)
				want := row.forms[i]
				if got != want {
					t.Errorf("%s subjunctive %s person %d = %q, want %q",
						tc.conj.label, subjunctiveTenseName(row.tense), vp, got, want)
				}
				asserted++
			}
		}
	}
	if asserted != 45 {
		t.Fatalf("asserted %d subjunctive cells, want 45 (3 conjugations x 3 tenses x 5 persons)", asserted)
	}
}

// TestSubjunctiveImperfectFirstPlural pins the accent on the proparoxytone
// imperfect subjunctive first person plural, which is the one subjunctive cell
// that carries a graphic accent. In particular the second conjugation takes the
// CIRCUMFLEX ê (comêssemos, closed quality), not the acute é or a bare e, while
// the first and third take the acute á and í (falássemos, partíssemos). A revert
// to the wrong accent (or none) fails this guard.
func TestSubjunctiveImperfectFirstPlural(t *testing.T) {
	cases := []struct {
		conj    *verbConjugation
		radical string
		want    string
		accent  rune // the graphic accent the form must carry
	}{
		{&conjAr, "fal", "falássemos", 'á'},
		{&conjEr, "com", "comêssemos", 'ê'}, // circumflex, not acute
		{&conjIr, "part", "partíssemos", 'í'},
	}
	for _, tc := range cases {
		got := conjugateSubjunctive(tc.radical, tc.conj, Imperfect, verbP1p)
		if got != tc.want {
			t.Errorf("%s imperfect subjunctive 1p = %q, want %q", tc.conj.label, got, tc.want)
		}
		if !strings.ContainsRune(got, tc.accent) {
			t.Errorf("%s imperfect subjunctive 1p %q lacks the expected accent %q",
				tc.conj.label, got, tc.accent)
		}
	}
	// Explicitly reject the two wrong spellings of the second-conjugation cell: the
	// bare e (comessemos) and the acute é (coméssemos). Only the circumflex is correct.
	com := conjugateSubjunctive("com", &conjEr, Imperfect, verbP1p)
	if com == "comessemos" || com == "coméssemos" {
		t.Errorf("second-conjugation imperfect subjunctive 1p = %q, want the circumflex form %q", com, "comêssemos")
	}
}

// TestSubjunctiveFutureEqualsPersonalInfinitive certifies the grammatical identity
// that the regular future subjunctive coincides with the inflected personal
// infinitive, so the two assemblers must agree cell by cell for every conjugation
// and person.
func TestSubjunctiveFutureEqualsPersonalInfinitive(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, vp := range persons {
			fut := conjugateSubjunctive(radicals[c], c, Future, vp)
			inf := conjugateNonFinite(radicals[c], c, formPersonalInfinitive, vp)
			if fut != inf {
				t.Errorf("%s future subjunctive person %d = %q, but personal infinitive = %q (must coincide)",
					c.label, vp, fut, inf)
			}
		}
	}
}

// TestImperativeExactForms is the core acceptance test for the imperative mood.
// For each conjugation, built from its model radical (fal/com/part), every
// imperative cell — 2 polarities x 4 valid persons (no first singular) per
// conjugation, 24 in total — must equal the textbook value exactly. This proves
// the derivation rule (affirmative 2s = present indicative 3s; every other person,
// and all negatives = present subjunctive) against the authoritative paradigm.
func TestImperativeExactForms(t *testing.T) {
	cases := []struct {
		conj    *verbConjugation
		radical string
		// indexed by verbPerson; the verbP1s slot is unused (no imperative 1s).
		affirmative [5]string
		negative    [5]string
	}{
		{
			conj: &conjAr, radical: "fal",
			affirmative: [5]string{verbP2s: "fala", verbP3s: "fale", verbP1p: "falemos", verbP3p: "falem"},
			negative:    [5]string{verbP2s: "fales", verbP3s: "fale", verbP1p: "falemos", verbP3p: "falem"},
		},
		{
			conj: &conjEr, radical: "com",
			affirmative: [5]string{verbP2s: "come", verbP3s: "coma", verbP1p: "comamos", verbP3p: "comam"},
			negative:    [5]string{verbP2s: "comas", verbP3s: "coma", verbP1p: "comamos", verbP3p: "comam"},
		},
		{
			conj: &conjIr, radical: "part",
			affirmative: [5]string{verbP2s: "parte", verbP3s: "parta", verbP1p: "partamos", verbP3p: "partam"},
			negative:    [5]string{verbP2s: "partas", verbP3s: "parta", verbP1p: "partamos", verbP3p: "partam"},
		},
	}

	asserted := 0
	for _, tc := range cases {
		for _, vp := range imperativePersons {
			if got := conjugateImperative(tc.radical, tc.conj, imperativeAff, vp); got != tc.affirmative[vp] {
				t.Errorf("%s affirmative imperative person %d = %q, want %q",
					tc.conj.label, vp, got, tc.affirmative[vp])
			}
			asserted++
			if got := conjugateImperative(tc.radical, tc.conj, imperativeNeg, vp); got != tc.negative[vp] {
				t.Errorf("%s negative imperative person %d = %q, want %q",
					tc.conj.label, vp, got, tc.negative[vp])
			}
			asserted++
		}
	}
	if asserted != 24 {
		t.Fatalf("asserted %d imperative cells, want 24 (3 conjugations x 2 polarities x 4 persons)", asserted)
	}
}

// TestImperativeDerivationRule verifies, independently of the hard-coded reference
// table above, that the derivation rule holds against the present indicative and
// present subjunctive assemblers: the affirmative second singular equals the
// present indicative third singular, and every other imperative cell equals the
// present subjunctive for that person.
func TestImperativeDerivationRule(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		radical := radicals[c]
		// Affirmative 2s = present indicative 3s.
		if got, want := conjugateImperative(radical, c, imperativeAff, verbP2s),
			conjugateIndicative(radical, c, Present, verbP3s); got != want {
			t.Errorf("%s affirmative imperative 2s = %q, want present indicative 3s %q", c.label, got, want)
		}
		// Affirmative 3s/1p/3p = present subjunctive.
		for _, vp := range []verbPerson{verbP3s, verbP1p, verbP3p} {
			if got, want := conjugateImperative(radical, c, imperativeAff, vp),
				conjugateSubjunctive(radical, c, Present, vp); got != want {
				t.Errorf("%s affirmative imperative person %d = %q, want present subjunctive %q",
					c.label, vp, got, want)
			}
		}
		// All negative persons = present subjunctive.
		for _, vp := range imperativePersons {
			if got, want := conjugateImperative(radical, c, imperativeNeg, vp),
				conjugateSubjunctive(radical, c, Present, vp); got != want {
				t.Errorf("%s negative imperative person %d = %q, want present subjunctive %q",
					c.label, vp, got, want)
			}
		}
	}
}

// TestImperativeNormalizesFirstSingular verifies the N2 normalization baked into
// the assembler: the imperative has no first person singular, so a verbP1s request
// must produce the first person plural form for both polarities.
func TestImperativeNormalizesFirstSingular(t *testing.T) {
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		radical := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}[c]
		for _, pol := range []imperativePolarity{imperativeAff, imperativeNeg} {
			got := conjugateImperative(radical, c, pol, verbP1s)
			want := conjugateImperative(radical, c, pol, verbP1p)
			if got != want {
				t.Errorf("%s %s imperative 1s = %q, want the 1p form %q (N2)",
					c.label, polarityName(pol), got, want)
			}
		}
	}
}

// TestConjugateSubjunctiveSingleAllocation asserts that the paradigm-exact
// subjunctive assembler performs exactly one allocation (the returned string) for
// every tense and person.
func TestConjugateSubjunctiveSingleAllocation(t *testing.T) {
	tenses := []Tense{Present, Imperfect, Future}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, tt := range tenses {
		for _, vp := range persons {
			allocs := testing.AllocsPerRun(1000, func() {
				_ = conjugateSubjunctive("part", &conjIr, tt, vp)
			})
			if allocs != 1 {
				t.Errorf("conjugateSubjunctive(%s, person %d) allocated %.0f times, want 1",
					subjunctiveTenseName(tt), vp, allocs)
			}
		}
	}
}

// TestConjugateImperativeSingleAllocation asserts that the paradigm-exact
// imperative assembler performs exactly one allocation (the returned string) for
// every polarity and person.
func TestConjugateImperativeSingleAllocation(t *testing.T) {
	persons := []verbPerson{verbP2s, verbP3s, verbP1p, verbP3p}
	for _, pol := range []imperativePolarity{imperativeAff, imperativeNeg} {
		for _, vp := range persons {
			allocs := testing.AllocsPerRun(1000, func() {
				_ = conjugateImperative("part", &conjIr, pol, vp)
			})
			if allocs != 1 {
				t.Errorf("conjugateImperative(%s, person %d) allocated %.0f times, want 1",
					polarityName(pol), vp, allocs)
			}
		}
	}
}

// TestSubjunctiveOrthographicConformance asserts that every subjunctive cell of the
// three model conjugations is an orthographically well-formed, lowercase,
// valid-UTF-8 pt-PT word (the Sprint 7 string oracles). This certifies the
// baked-in accents: only the acute (á, í) and the circumflex (ê) appear, and no
// cedilla, diaeresis or grave is ever produced.
func TestSubjunctiveOrthographicConformance(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{Present, Imperfect, Future}
	persons := []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, vp := range persons {
				w := conjugateSubjunctive(radicals[c], c, tt, vp)
				assertVerbOrthographyConformant(t, w)
				if !strings.HasPrefix(w, radicals[c]) {
					t.Errorf("%s subjunctive %s person %d = %q does not start with radical %q",
						c.label, subjunctiveTenseName(tt), vp, w, radicals[c])
				}
			}
		}
	}
}

// TestImperativeOrthographicConformance asserts that every imperative cell of the
// three model conjugations is an orthographically well-formed, lowercase,
// valid-UTF-8 pt-PT word and starts with the model radical. The regular imperative
// carries no graphic accent, so these are plain ASCII forms.
func TestImperativeOrthographicConformance(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, pol := range []imperativePolarity{imperativeAff, imperativeNeg} {
			for _, vp := range imperativePersons {
				w := conjugateImperative(radicals[c], c, pol, vp)
				assertVerbOrthographyConformant(t, w)
				if !strings.HasPrefix(w, radicals[c]) {
					t.Errorf("%s %s imperative person %d = %q does not start with radical %q",
						c.label, polarityName(pol), vp, w, radicals[c])
				}
			}
		}
	}
}

// TestResolveSubjunctiveTense verifies the subjunctive tense resolution: the three
// existing tenses pass through unchanged, the N3 normalization maps the two
// indicative-only tenses to the nearest subjunctive tense (Preterite -> Imperfect,
// Conditional -> Future), and AnyTense reaches all three subjunctive tenses and
// never yields Preterite or Conditional.
func TestResolveSubjunctiveTense(t *testing.T) {
	r := rand.New(rand.NewPCG(11, 13))
	for _, tt := range []Tense{Present, Imperfect, Future} {
		if got := resolveSubjunctiveTense(r, tt); got != tt {
			t.Errorf("resolveSubjunctiveTense(%s) = %s, want passthrough",
				subjunctiveTenseName(tt), subjunctiveTenseName(got))
		}
	}
	if got := resolveSubjunctiveTense(r, Preterite); got != Imperfect {
		t.Errorf("resolveSubjunctiveTense(Preterite) = %s, want Imperfect (N3)", tenseName(got))
	}
	if got := resolveSubjunctiveTense(r, Conditional); got != Future {
		t.Errorf("resolveSubjunctiveTense(Conditional) = %s, want Future (N3)", tenseName(got))
	}
	seen := map[Tense]bool{}
	for i := 0; i < 20000; i++ {
		got := resolveSubjunctiveTense(r, AnyTense)
		if got == Preterite || got == Conditional {
			t.Fatalf("AnyTense resolved to %s, which the subjunctive lacks", tenseName(got))
		}
		seen[got] = true
	}
	for _, tt := range []Tense{Present, Imperfect, Future} {
		if !seen[tt] {
			t.Errorf("AnyTense never resolved to %s", subjunctiveTenseName(tt))
		}
	}
}

// TestResolveImperativePerson verifies the imperative person resolution: the N2
// normalization (First + Singular -> First + Plural, since there is no imperative
// first singular), the inherited N1 normalization (Second + Plural -> Third
// Plural), that verbP1s is never produced, and that AnyPerson/AnyNumber reaches
// all four imperative slots.
func TestResolveImperativePerson(t *testing.T) {
	r := rand.New(rand.NewPCG(17, 19))

	concrete := []struct {
		p    Person
		n    Number
		want verbPerson
	}{
		{First, Singular, verbP1p}, // N2: no imperative 1s -> 1p
		{Second, Singular, verbP2s},
		{Third, Singular, verbP3s},
		{First, Plural, verbP1p},
		{Third, Plural, verbP3p},
		{Second, Plural, verbP3p}, // N1: vós -> vocês (third plural)
	}
	for _, tc := range concrete {
		if got := resolveImperativePerson(r, tc.p, tc.n); got != tc.want {
			t.Errorf("resolveImperativePerson(%d,%d) = %d, want %d", tc.p, tc.n, got, tc.want)
		}
	}

	// verbP1s must never be produced, and every imperative slot must be reachable.
	seen := map[verbPerson]bool{}
	for i := 0; i < 40000; i++ {
		got := resolveImperativePerson(r, AnyPerson, AnyNumber)
		if got == verbP1s {
			t.Fatalf("resolveImperativePerson produced verbP1s, which the imperative lacks")
		}
		seen[got] = true
	}
	for _, vp := range imperativePersons {
		if !seen[vp] {
			t.Errorf("Any person/number never reached imperative slot %d", vp)
		}
	}
}

// TestSubjunctiveVerbFormConformance exercises the injectable-*rand.Rand assembler
// over every conjugation, tense (including AnyTense and the two indicative-only
// tenses that normalize under N3), person and number, asserting that each form is
// orthographically conformant and starts with the supplied radical.
func TestSubjunctiveVerbFormConformance(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5B, 0x5C))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{AnyTense, Present, Imperfect, Future, Preterite, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 200; i++ {
						w := subjunctiveVerbForm(r, radicals[c], c, tt, p, n)
						assertVerbOrthographyConformant(t, w)
						if !strings.HasPrefix(w, radicals[c]) {
							t.Fatalf("%s form %q does not start with radical %q", c.label, w, radicals[c])
						}
					}
				}
			}
		}
	}
}

// TestImperativeVerbFormConformance exercises the injectable-*rand.Rand imperative
// assembler over every conjugation, polarity, person (including the First +
// Singular request that normalizes under N2) and number, asserting that each form
// is orthographically conformant and starts with the supplied radical.
func TestImperativeVerbFormConformance(t *testing.T) {
	r := rand.New(rand.NewPCG(0x1319, 0x2337))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, pol := range []imperativePolarity{imperativeAff, imperativeNeg} {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 200; i++ {
						w := imperativeVerbForm(r, radicals[c], c, pol, p, n)
						assertVerbOrthographyConformant(t, w)
						if !strings.HasPrefix(w, radicals[c]) {
							t.Fatalf("%s form %q does not start with radical %q", c.label, w, radicals[c])
						}
					}
				}
			}
		}
	}
}

// TestSubjunctiveVerbFormReproducible verifies that the injectable-*rand.Rand
// subjunctive assembler is deterministic: two sources with the same seed produce
// identical sequences for the same arguments.
func TestSubjunctiveVerbFormReproducible(t *testing.T) {
	a := rand.New(rand.NewPCG(77, 77))
	b := rand.New(rand.NewPCG(77, 77))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	tenses := []Tense{AnyTense, Present, Imperfect, Future, Preterite, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 50; i++ {
						x := subjunctiveVerbForm(a, radicals[c], c, tt, p, n)
						y := subjunctiveVerbForm(b, radicals[c], c, tt, p, n)
						if x != y {
							t.Fatalf("subjunctiveVerbForm diverged (%s, %s, p %d, n %d) at %d: %q != %q",
								c.label, subjunctiveTenseName(tt), p, n, i, x, y)
						}
					}
				}
			}
		}
	}
}

// TestImperativeVerbFormReproducible verifies that the injectable-*rand.Rand
// imperative assembler is deterministic: two sources with the same seed produce
// identical sequences for the same arguments.
func TestImperativeVerbFormReproducible(t *testing.T) {
	a := rand.New(rand.NewPCG(88, 88))
	b := rand.New(rand.NewPCG(88, 88))
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		for _, pol := range []imperativePolarity{imperativeAff, imperativeNeg} {
			for _, p := range persons {
				for _, n := range numbers {
					for i := 0; i < 50; i++ {
						x := imperativeVerbForm(a, radicals[c], c, pol, p, n)
						y := imperativeVerbForm(b, radicals[c], c, pol, p, n)
						if x != y {
							t.Fatalf("imperativeVerbForm diverged (%s, %s, p %d, n %d) at %d: %q != %q",
								c.label, polarityName(pol), p, n, i, x, y)
						}
					}
				}
			}
		}
	}
}

// BenchmarkConjugateSubjunctive measures the paradigm-exact subjunctive assembler
// and its allocation budget: a subjunctive form must assemble in exactly one
// string allocation. The accented proparoxytone imperfect first plural is the
// worst case (a multi-byte accent in the desinence).
func BenchmarkConjugateSubjunctive(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = conjugateSubjunctive("com", &conjEr, Imperfect, verbP1p)
	}
}

// BenchmarkSubjunctiveVerbForm measures the injectable-*rand.Rand subjunctive
// assembler, including the tense and person/number resolution, for a fully
// unspecified request.
func BenchmarkSubjunctiveVerbForm(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = subjunctiveVerbForm(r, "fal", &conjAr, AnyTense, AnyPerson, AnyNumber)
	}
}

// BenchmarkConjugateImperative measures the paradigm-exact imperative assembler
// and its allocation budget: an imperative form must assemble in exactly one
// string allocation.
func BenchmarkConjugateImperative(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = conjugateImperative("fal", &conjAr, imperativeNeg, verbP2s)
	}
}

// BenchmarkImperativeVerbForm measures the injectable-*rand.Rand imperative
// assembler, including the person/number resolution, for a fully unspecified
// request.
func BenchmarkImperativeVerbForm(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 1))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = imperativeVerbForm(r, "fal", &conjAr, imperativeAff, AnyPerson, AnyNumber)
	}
}
