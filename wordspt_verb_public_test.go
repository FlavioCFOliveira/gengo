package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
	"unicode/utf8"
)

// This file verifies the public pt-PT verb generators (VerbPT, VerbPTOf, and their
// *Generator methods, task #31). The checks map to the task and specification
// acceptance criteria: orthographic conformance over a large sample, mood/tense/
// person/number honoring with Any producing varied output, the deterministic
// normalization matrix (N1, N2, N3), upward length normalization for unviable
// categories, reproducibility of a seeded Generator, the single-allocation budget,
// and the two resolved design flags (infinitive Any-person selection and the
// present-tense radical-accentuation conclusion).

// verbSampleSize is the sampled-run size for the conformance tests. It exceeds the
// task's >=1000 floor comfortably while keeping the suite fast.
const verbSampleSize = 20000

// allMoods lists the seven concrete moods, for exhaustive option loops.
var allMoods = []Mood{
	Indicative, Subjunctive, ImperativeAffirmative, ImperativeNegative,
	Infinitive, Gerund, Participle,
}

// concreteRand returns a source usable where the resolution consumes no randomness
// (every option is concrete), so verbDesinence is deterministic.
func concreteRand() *rand.Rand { return rand.New(rand.NewPCG(0, 0)) }

// verbSlotDesinences returns, for a fully concrete slot, the desinence each of the
// three conjugations appends to the bare radical, using the exact resolution the
// core uses (verbDesinence). A generated form of that slot must end in one of them.
func verbSlotDesinences(m Mood, t Tense, p Person, n Number) []string {
	r := concreteRand()
	out := make([]string, 0, 3)
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		out = append(out, verbDesinence(r, m, c, t, p, n))
	}
	return out
}

// ---------------------------------------------------------------------------
// Orthographic conformance (AC1, AC7, AC10)
// ---------------------------------------------------------------------------

// TestVerbPTOrthographicConformance samples VerbPT and every mood/tense/person/
// number/length combination of VerbPTOf and asserts that every output is
// orthographically conformant, lowercase, and valid UTF-8 (the Sprint 7 string
// oracles). This certifies, in particular, that no finite desinence attaching to a
// sampled radical ever produces a cedilla-rule violation (the reason the verb
// radical excludes the frontness-sensitive onsets).
func TestVerbPTOrthographicConformance(t *testing.T) {
	for i := 0; i < verbSampleSize; i++ {
		assertVerbOrthographyConformant(t, VerbPT())
	}

	tenses := []Tense{AnyTense, Present, Imperfect, Preterite, Future, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	g := New(0xB0B)
	for _, m := range allMoods {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for _, l := range lengths {
						for i := 0; i < 60; i++ {
							assertVerbOrthographyConformant(t, VerbPTOf(m, tt, p, n, l))
							assertVerbOrthographyConformant(t, g.VerbPTOf(m, tt, p, n, l))
						}
					}
				}
			}
		}
	}
}

// TestVerbDesinenceMatchesAssemblers proves that the desinence resolution used by
// the public core is faithful to the #28..#30 paradigm-exact assemblers: for the
// model radicals fal/com/part, radical + verbDesinence(...) equals the
// corresponding conjugateNonFinite/conjugateIndicative/conjugateSubjunctive/
// conjugateImperative output for the same concrete slot. This ties the public
// generator to the authoritative Cunha & Cintra tables.
func TestVerbDesinenceMatchesAssemblers(t *testing.T) {
	radicals := map[*verbConjugation]string{&conjAr: "fal", &conjEr: "com", &conjIr: "part"}
	r := concreteRand()
	persons := []struct {
		p  Person
		n  Number
		vp verbPerson
	}{
		{First, Singular, verbP1s},
		{Second, Singular, verbP2s},
		{Third, Singular, verbP3s},
		{First, Plural, verbP1p},
		{Third, Plural, verbP3p},
	}
	for _, c := range []*verbConjugation{&conjAr, &conjEr, &conjIr} {
		rad := radicals[c]
		for _, pr := range persons {
			// Indicative, every tense.
			for _, tt := range []Tense{Present, Imperfect, Preterite, Future, Conditional} {
				got := rad + verbDesinence(r, Indicative, c, tt, pr.p, pr.n)
				want := conjugateIndicative(rad, c, tt, pr.vp)
				if got != want {
					t.Errorf("indicative %s %s %d: %q != %q", c.label, tenseName(tt), pr.vp, got, want)
				}
			}
			// Subjunctive, every tense.
			for _, tt := range []Tense{Present, Imperfect, Future} {
				got := rad + verbDesinence(r, Subjunctive, c, tt, pr.p, pr.n)
				want := conjugateSubjunctive(rad, c, tt, pr.vp)
				if got != want {
					t.Errorf("subjunctive %s %s %d: %q != %q", c.label, subjunctiveTenseName(tt), pr.vp, got, want)
				}
			}
			// Personal infinitive (a concrete person forces the personal form).
			got := rad + verbDesinence(r, Infinitive, c, AnyTense, pr.p, pr.n)
			want := conjugateNonFinite(rad, c, formPersonalInfinitive, pr.vp)
			if got != want {
				t.Errorf("personal infinitive %s %d: %q != %q", c.label, pr.vp, got, want)
			}
		}
		// Gerund and participle (invariable).
		if got, want := rad+verbDesinence(r, Gerund, c, AnyTense, AnyPerson, AnyNumber),
			conjugateNonFinite(rad, c, formGerund, verbP1s); got != want {
			t.Errorf("gerund %s: %q != %q", c.label, got, want)
		}
		if got, want := rad+verbDesinence(r, Participle, c, AnyTense, AnyPerson, AnyNumber),
			conjugateNonFinite(rad, c, formParticiple, verbP1s); got != want {
			t.Errorf("participle %s: %q != %q", c.label, got, want)
		}
		// Imperative, both polarities, every valid person.
		for _, vp := range imperativePersons {
			p, n := verbPersonToPersonNumber(vp)
			if got, want := rad+verbDesinence(r, ImperativeAffirmative, c, AnyTense, p, n),
				conjugateImperative(rad, c, imperativeAff, vp); got != want {
				t.Errorf("affirmative imperative %s %d: %q != %q", c.label, vp, got, want)
			}
			if got, want := rad+verbDesinence(r, ImperativeNegative, c, AnyTense, p, n),
				conjugateImperative(rad, c, imperativeNeg, vp); got != want {
				t.Errorf("negative imperative %s %d: %q != %q", c.label, vp, got, want)
			}
		}
	}
}

// verbPersonToPersonNumber maps a verbPerson slot back to a public Person/Number
// pair, for driving the assembler-equivalence test.
func verbPersonToPersonNumber(vp verbPerson) (Person, Number) {
	switch vp {
	case verbP1s:
		return First, Singular
	case verbP2s:
		return Second, Singular
	case verbP3s:
		return Third, Singular
	case verbP1p:
		return First, Plural
	default: // verbP3p
		return Third, Plural
	}
}

// ---------------------------------------------------------------------------
// Mood / tense / person / number honoring (AC2, AC5)
// ---------------------------------------------------------------------------

// TestVerbPTSlotHonored verifies that a concrete finite/non-finite slot produces a
// form ending in a desinence valid for that exact slot, across all three
// conjugations. This ties the public output to the resolved mood, tense, person and
// number.
func TestVerbPTSlotHonored(t *testing.T) {
	g := New(7)
	type slot struct {
		m Mood
		t Tense
		p Person
		n Number
	}
	slots := []slot{
		{Indicative, Present, First, Singular},
		{Indicative, Present, Second, Singular},
		{Indicative, Preterite, First, Plural},
		{Indicative, Imperfect, Third, Plural},
		{Indicative, Future, Third, Singular},
		{Indicative, Conditional, First, Plural},
		{Subjunctive, Present, First, Singular},
		{Subjunctive, Imperfect, First, Plural},
		{Subjunctive, Future, Third, Plural},
		{ImperativeAffirmative, AnyTense, Second, Singular},
		{ImperativeAffirmative, AnyTense, Third, Plural},
		{ImperativeNegative, AnyTense, Second, Singular},
		{Infinitive, AnyTense, Second, Singular},
		{Infinitive, AnyTense, First, Plural},
		{Gerund, AnyTense, AnyPerson, AnyNumber},
		{Participle, AnyTense, AnyPerson, AnyNumber},
	}
	for _, s := range slots {
		des := verbSlotDesinences(s.m, s.t, s.p, s.n)
		for i := 0; i < 3000; i++ {
			w := g.VerbPTOf(s.m, s.t, s.p, s.n, AnyLengthWord)
			if !hasAnySuffix(w, des...) {
				t.Fatalf("mood %d tense %d person %d number %d: %q ends in none of %v",
					s.m, s.t, s.p, s.n, w, des)
			}
		}
	}
}

// TestVerbPTPresentTenseHonored confirms the present indicative and present
// subjunctive are distinguishable and honored: present indicative singular ends in
// a present-indicative desinence, and present subjunctive ends in a
// present-subjunctive desinence, which for the first person singular differ (falo
// vs fale).
func TestVerbPTPresentTenseHonored(t *testing.T) {
	g := New(13)
	indDes := verbSlotDesinences(Indicative, Present, First, Singular)  // {o,o,o}
	subDes := verbSlotDesinences(Subjunctive, Present, First, Singular) // {e,a,a}
	for i := 0; i < verbSampleSize; i++ {
		if w := g.VerbPTOf(Indicative, Present, First, Singular, AnyLengthWord); !hasAnySuffix(w, indDes...) {
			t.Fatalf("present indicative 1s %q ends in none of %v", w, indDes)
		}
		if w := g.VerbPTOf(Subjunctive, Present, First, Singular, AnyLengthWord); !hasAnySuffix(w, subDes...) {
			t.Fatalf("present subjunctive 1s %q ends in none of %v", w, subDes)
		}
	}
}

// TestVerbPTAnyMoodVaried verifies that AnyMood, over a sampled run, reaches every
// mood. Because a bare surface cannot always be classified into a single mood, the
// test drives resolveMood (the exact resolver verbCore uses) directly, asserting
// every mood is reached and the finite indicative dominates.
func TestVerbPTAnyMoodVaried(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 5))
	counts := map[Mood]int{}
	const n = 200000
	for i := 0; i < n; i++ {
		counts[resolveMood(r, AnyMood)]++
	}
	for _, m := range allMoods {
		if counts[m] == 0 {
			t.Errorf("AnyMood never reached mood %d", m)
		}
	}
	for _, m := range allMoods {
		if m == Indicative {
			continue
		}
		if counts[Indicative] <= counts[m] {
			t.Errorf("indicative not dominant: indicative=%d, mood %d=%d", counts[Indicative], m, counts[m])
		}
	}
}

// TestResolveMoodPassthrough verifies that a concrete mood is returned unchanged.
func TestResolveMoodPassthrough(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 1))
	for _, m := range allMoods {
		if got := resolveMood(r, m); got != m {
			t.Errorf("resolveMood(%d) = %d, want passthrough", m, got)
		}
	}
}

// TestVerbPTIndicativeTenseVaried verifies that AnyTense, in the indicative mood,
// reaches all five tenses over a sampled run (via the present/imperfect/preterite/
// future/conditional desinence families), and TestVerbPTPersonNumberVaried checks
// the person/number coverage.
func TestVerbPTIndicativeTenseVaried(t *testing.T) {
	g := New(101)
	// Distinctive tense discriminators for a concrete person (3s) across tenses.
	// 3s: present a/e, imperfect ava/ia, preterite ou/eu/iu, future ará/erá/irá,
	// conditional aria/eria/iria.
	sawPret, sawFut, sawCond, sawImperf := false, false, false, false
	for i := 0; i < verbSampleSize; i++ {
		w := g.VerbPTOf(Indicative, AnyTense, Third, Singular, AnyLengthWord)
		switch {
		case hasAnySuffix(w, "ou", "eu", "iu"):
			sawPret = true
		case strings.HasSuffix(w, "á"):
			sawFut = true
		case hasAnySuffix(w, "aria", "eria", "iria"):
			sawCond = true
		case hasAnySuffix(w, "ava", "ia"):
			sawImperf = true
		}
	}
	if !sawPret || !sawFut || !sawCond || !sawImperf {
		t.Errorf("AnyTense not varied: preterite=%v future=%v conditional=%v imperfect=%v",
			sawPret, sawFut, sawCond, sawImperf)
	}
}

// TestVerbPTPersonNumberVaried verifies that AnyPerson/AnyNumber, in the present
// indicative, reaches singular and plural forms (plural 1p/3p end in -mos/-am/-em;
// singular ends otherwise).
func TestVerbPTPersonNumberVaried(t *testing.T) {
	g := New(202)
	sawSingular, sawPlural := false, false
	for i := 0; i < verbSampleSize; i++ {
		w := g.VerbPTOf(Indicative, Present, AnyPerson, AnyNumber, AnyLengthWord)
		if hasAnySuffix(w, "mos", "am", "em") {
			sawPlural = true
		} else {
			sawSingular = true
		}
	}
	if !sawSingular || !sawPlural {
		t.Errorf("Any person/number not varied: singular=%v plural=%v", sawSingular, sawPlural)
	}
}

// ---------------------------------------------------------------------------
// Flag #1: infinitive Any-person selection (impersonal vs personal)
// ---------------------------------------------------------------------------

// TestResolveInfinitiveForm verifies the infinitive-selection rule: a concrete
// person or number forces the inflected personal infinitive, while both-Any yields
// a weighted choice between the impersonal infinitive (the common default) and a
// random personal infinitive, with both reachable and the impersonal favored.
func TestResolveInfinitiveForm(t *testing.T) {
	r := rand.New(rand.NewPCG(23, 29))

	// A concrete person forces the personal infinitive.
	for _, p := range []Person{First, Second, Third} {
		for _, n := range []Number{AnyNumber, Singular, Plural} {
			for i := 0; i < 1000; i++ {
				form, _ := resolveInfinitiveForm(r, p, n)
				if form != formPersonalInfinitive {
					t.Fatalf("resolveInfinitiveForm(%d,%d) form=%d, want personal", p, n, form)
				}
			}
		}
	}
	// A concrete number (with Any person) forces the personal infinitive.
	for _, n := range []Number{Singular, Plural} {
		for i := 0; i < 1000; i++ {
			if form, _ := resolveInfinitiveForm(r, AnyPerson, n); form != formPersonalInfinitive {
				t.Fatalf("resolveInfinitiveForm(Any,%d) form=%d, want personal", n, form)
			}
		}
	}

	// Both Any: impersonal (default) and personal both reachable; impersonal favored.
	impersonal, personal := 0, 0
	const n = 200000
	for i := 0; i < n; i++ {
		form, _ := resolveInfinitiveForm(r, AnyPerson, AnyNumber)
		if form == formImpersonalInfinitive {
			impersonal++
		} else {
			personal++
		}
	}
	if impersonal == 0 || personal == 0 {
		t.Fatalf("both-Any infinitive not varied: impersonal=%d personal=%d", impersonal, personal)
	}
	if impersonal <= personal {
		t.Errorf("impersonal not favored: impersonal=%d personal=%d", impersonal, personal)
	}
}

// TestVerbPTInfinitiveConcretePersonIsPersonal verifies end-to-end that a concrete
// person on the infinitive yields the inflected personal infinitive: for the second
// person singular, the form ends in -ares/-eres/-ires (never the bare -ar/-er/-ir),
// and for the first plural in -armos/-ermos/-irmos.
func TestVerbPTInfinitiveConcretePersonIsPersonal(t *testing.T) {
	g := New(31)
	for i := 0; i < verbSampleSize; i++ {
		if w := g.VerbPTOf(Infinitive, AnyTense, Second, Singular, AnyLengthWord); !hasAnySuffix(w, "ares", "eres", "ires") {
			t.Fatalf("infinitive 2s %q is not the personal -ares/-eres/-ires form", w)
		}
		if w := g.VerbPTOf(Infinitive, AnyTense, First, Plural, AnyLengthWord); !hasAnySuffix(w, "armos", "ermos", "irmos") {
			t.Fatalf("infinitive 1p %q is not the personal -armos/-ermos/-irmos form", w)
		}
	}
}

// ---------------------------------------------------------------------------
// Flag #2: present-tense radical accentuation (empirical conclusion)
// ---------------------------------------------------------------------------

// TestVerbPTPresentNeedsNoAccent is the empirical backing of the resolved flag #2:
// over a large sample, every generated present indicative and present subjunctive
// form carries NO acute or circumflex accent AND ends in a default paroxytone
// desinence (one of -o, -a, -e, -as, -es, -am, -em, -amos, -emos, -imos). Because a
// default paroxytone ending requires no graphic accent under the pt-PT rules
// regardless of the stressed vowel, a bare radical plus an unaccented present
// desinence is orthographically correct for every sampled radical: no present form
// requires an accent that the generator omits, so no Layer 3 call is needed.
func TestVerbPTPresentNeedsNoAccent(t *testing.T) {
	g := New(0xACCE)
	// The complete set of present desinences (indicative and subjunctive), all of
	// which are default paroxytone endings that take no graphic accent.
	presentDesinences := []string{"o", "a", "e", "as", "es", "am", "em", "amos", "emos", "imos"}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	for _, m := range []Mood{Indicative, Subjunctive} {
		for _, p := range persons {
			for _, n := range numbers {
				for i := 0; i < 2000; i++ {
					w := g.VerbPTOf(m, Present, p, n, AnyLengthWord)
					if hasAcuteOrCircumflexRune(w) {
						t.Fatalf("present %s form %q carries an acute/circumflex accent", moodName(m), w)
					}
					if !hasAnySuffix(w, presentDesinences...) {
						t.Fatalf("present %s form %q ends in no default paroxytone desinence", moodName(m), w)
					}
				}
			}
		}
	}
}

// hasAcuteOrCircumflexRune reports whether word contains any rune bearing an acute
// or circumflex accent (á, â, é, ê, í, ó, ô, ú), reusing the rune helpers.
func hasAcuteOrCircumflexRune(word string) bool {
	for _, r := range word {
		if hasAcute(r) || hasCircumflex(r) {
			return true
		}
	}
	return false
}

// moodName returns a readable mood label for test failure messages only.
func moodName(m Mood) string {
	switch m {
	case Indicative:
		return "indicative"
	case Subjunctive:
		return "subjunctive"
	case ImperativeAffirmative:
		return "imperative-affirmative"
	case ImperativeNegative:
		return "imperative-negative"
	case Infinitive:
		return "infinitive"
	case Gerund:
		return "gerund"
	case Participle:
		return "participle"
	default:
		return "any"
	}
}

// ---------------------------------------------------------------------------
// Normalization matrix (AC6): N1, N2, N3
// ---------------------------------------------------------------------------

// TestVerbPTNormalizationN1 verifies N1 (second person plural -> third person
// plural). For a seeded Generator, a Second+Plural request and a Third+Plural
// request draw identical randomness (both are concrete, so person/number resolution
// consumes none) and resolve to the same slot, so they must produce byte-identical
// output in every finite mood.
func TestVerbPTNormalizationN1(t *testing.T) {
	finite := []Mood{Indicative, Subjunctive, ImperativeAffirmative, ImperativeNegative}
	tenses := []Tense{AnyTense, Present, Imperfect, Future}
	for _, m := range finite {
		for _, tt := range tenses {
			a := New(0x4E31)
			b := New(0x4E31)
			for i := 0; i < 2000; i++ {
				second := a.VerbPTOf(m, tt, Second, Plural, AnyLengthWord)
				third := b.VerbPTOf(m, tt, Third, Plural, AnyLengthWord)
				if second != third {
					t.Fatalf("N1 %s tense %d: Second+Plural %q != Third+Plural %q", moodName(m), tt, second, third)
				}
			}
		}
	}
}

// TestVerbPTNormalizationN2 verifies N2 (first person singular imperative -> first
// person plural). A First+Singular imperative request and a First+Plural request
// resolve to the same slot with identical randomness, so they produce
// byte-identical output for both polarities.
func TestVerbPTNormalizationN2(t *testing.T) {
	for _, m := range []Mood{ImperativeAffirmative, ImperativeNegative} {
		a := New(0x4E32)
		b := New(0x4E32)
		for i := 0; i < 4000; i++ {
			singular := a.VerbPTOf(m, AnyTense, First, Singular, AnyLengthWord)
			plural := b.VerbPTOf(m, AnyTense, First, Plural, AnyLengthWord)
			if singular != plural {
				t.Fatalf("N2 %s: First+Singular %q != First+Plural %q", moodName(m), singular, plural)
			}
		}
	}
}

// TestVerbPTNormalizationN3Subjunctive verifies N3 for the subjunctive: a Preterite
// request normalizes to Imperfect and a Conditional request to Future. The
// normalized and target tenses resolve identically (both concrete), so they produce
// byte-identical output.
func TestVerbPTNormalizationN3Subjunctive(t *testing.T) {
	persons := []struct {
		p Person
		n Number
	}{{First, Singular}, {Third, Plural}, {First, Plural}}
	pairs := []struct{ requested, target Tense }{
		{Preterite, Imperfect},
		{Conditional, Future},
	}
	for _, pr := range persons {
		for _, tp := range pairs {
			a := New(0x4E33)
			b := New(0x4E33)
			for i := 0; i < 3000; i++ {
				req := a.VerbPTOf(Subjunctive, tp.requested, pr.p, pr.n, AnyLengthWord)
				tar := b.VerbPTOf(Subjunctive, tp.target, pr.p, pr.n, AnyLengthWord)
				if req != tar {
					t.Fatalf("N3 subjunctive %s->%s (person %d number %d): %q != %q",
						tenseName(tp.requested), tenseName(tp.target), pr.p, pr.n, req, tar)
				}
			}
		}
	}
}

// TestVerbPTNormalizationN3TenseIgnored verifies N3 for the moods without tense:
// the non-finite moods (Infinitive, Gerund, Participle) and the imperative moods
// ignore a requested tense entirely, so varying only the tense (with identical
// randomness) produces byte-identical output.
func TestVerbPTNormalizationN3TenseIgnored(t *testing.T) {
	tenseless := []Mood{Infinitive, Gerund, Participle, ImperativeAffirmative, ImperativeNegative}
	tenses := []Tense{Present, Imperfect, Preterite, Future, Conditional}
	for _, m := range tenseless {
		for _, tt := range tenses {
			a := New(0x4E34)
			b := New(0x4E34)
			for i := 0; i < 2000; i++ {
				base := a.VerbPTOf(m, AnyTense, First, Plural, AnyLengthWord)
				withTense := b.VerbPTOf(m, tt, First, Plural, AnyLengthWord)
				if base != withTense {
					t.Fatalf("N3 %s: tense %s changed the output %q != %q", moodName(m), tenseName(tt), base, withTense)
				}
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Length conformance and normalization (AC8)
// ---------------------------------------------------------------------------

// TestVerbPTLengthHonored verifies that every output of a concrete length category
// lands within the expected character window, including the upward normalization of
// unviable categories: a short slot (present 1s) fits every category, while a long
// slot (conditional 1p, imperfect subjunctive 1p) and the -mente-like gerund
// normalize the smaller categories upward. Length is measured in runes.
func TestVerbPTLengthHonored(t *testing.T) {
	g := New(29)
	cases := []struct {
		name   string
		m      Mood
		t      Tense
		p      Person
		n      Number
		l      LengthTypeWords
		lo, hi int
	}{
		{"small present 1s", Indicative, Present, First, Singular, SmallLengthWord, 1, 4},
		{"medium present 1s", Indicative, Present, First, Singular, MediumLengthWords, 5, 8},
		{"big present 1s", Indicative, Present, First, Singular, BigLengthWords, 9, 30},
		{"small conditional 1p -> big", Indicative, Conditional, First, Plural, SmallLengthWord, 9, 30},
		{"medium conditional 1p -> big", Indicative, Conditional, First, Plural, MediumLengthWords, 9, 30},
		{"big conditional 1p", Indicative, Conditional, First, Plural, BigLengthWords, 9, 30},
		{"small gerund -> medium", Gerund, AnyTense, AnyPerson, AnyNumber, SmallLengthWord, 5, 8},
		{"medium gerund", Gerund, AnyTense, AnyPerson, AnyNumber, MediumLengthWords, 5, 8},
		{"big imperfect subjunctive 1p", Subjunctive, Imperfect, First, Plural, BigLengthWords, 9, 30},
		{"small imperative aff 2s", ImperativeAffirmative, AnyTense, Second, Singular, SmallLengthWord, 1, 4},
	}
	for _, tc := range cases {
		for i := 0; i < verbSampleSize; i++ {
			w := g.VerbPTOf(tc.m, tc.t, tc.p, tc.n, tc.l)
			r := utf8.RuneCountInString(w)
			if r < tc.lo || r > tc.hi {
				t.Fatalf("%s: verb %q has %d runes, outside window [%d,%d]", tc.name, w, r, tc.lo, tc.hi)
			}
		}
	}
}

// TestVerbPTAnyLengthVaried verifies that AnyLengthWord yields verbs across the
// small, medium and big character buckets over a sampled run, using a slot that is
// viable across all three (present 1s).
func TestVerbPTAnyLengthVaried(t *testing.T) {
	g := New(303)
	counts := map[int]int{}
	for i := 0; i < verbSampleSize; i++ {
		r := utf8.RuneCountInString(g.VerbPTOf(Indicative, Present, First, Singular, AnyLengthWord))
		switch {
		case r <= 4:
			counts[1]++
		case r <= 8:
			counts[2]++
		default:
			counts[3]++
		}
	}
	if counts[1] == 0 || counts[2] == 0 || counts[3] == 0 {
		t.Errorf("Any length not varied: small=%d medium=%d big=%d", counts[1], counts[2], counts[3])
	}
}

// ---------------------------------------------------------------------------
// Reproducibility, allocation, concurrency (AC12)
// ---------------------------------------------------------------------------

// TestVerbPTReproducible verifies that a seeded Generator is deterministic: two
// generators created with the same seed produce identical verb sequences, for both
// the sugar method and the fully specified method across every option combination.
func TestVerbPTReproducible(t *testing.T) {
	a := New(0xC0FFEE)
	b := New(0xC0FFEE)
	for i := 0; i < verbSampleSize; i++ {
		if x, y := a.VerbPT(), b.VerbPT(); x != y {
			t.Fatalf("VerbPT diverged at %d: %q != %q", i, x, y)
		}
	}

	tenses := []Tense{AnyTense, Present, Imperfect, Preterite, Future, Conditional}
	persons := []Person{AnyPerson, First, Second, Third}
	numbers := []Number{AnyNumber, Singular, Plural}
	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	c := New(42)
	d := New(42)
	for _, m := range allMoods {
		for _, tt := range tenses {
			for _, p := range persons {
				for _, n := range numbers {
					for _, l := range lengths {
						for i := 0; i < 8; i++ {
							if x, y := c.VerbPTOf(m, tt, p, n, l), d.VerbPTOf(m, tt, p, n, l); x != y {
								t.Fatalf("VerbPTOf(%d,%d,%d,%d,%d) diverged at %d: %q != %q",
									m, tt, p, n, l, i, x, y)
							}
						}
					}
				}
			}
		}
	}
}

// TestVerbPTSingleAllocation asserts the specification's allocation budget (AC12):
// generating a single verb performs exactly one string allocation, across a
// representative set of moods (finite and non-finite, including the accented
// imperfect subjunctive first plural).
func TestVerbPTSingleAllocation(t *testing.T) {
	g := New(1)
	cases := []struct {
		name string
		m    Mood
		t    Tense
		p    Person
		n    Number
	}{
		{"present indicative", Indicative, Present, AnyPerson, AnyNumber},
		{"conditional indicative", Indicative, Conditional, First, Plural},
		{"imperfect subjunctive 1p", Subjunctive, Imperfect, First, Plural},
		{"imperative negative", ImperativeNegative, AnyTense, AnyPerson, AnyNumber},
		{"infinitive", Infinitive, AnyTense, AnyPerson, AnyNumber},
		{"gerund", Gerund, AnyTense, AnyPerson, AnyNumber},
		{"participle", Participle, AnyTense, AnyPerson, AnyNumber},
		{"fully random", AnyMood, AnyTense, AnyPerson, AnyNumber},
	}
	for _, tc := range cases {
		allocs := testing.AllocsPerRun(1000, func() {
			_ = g.VerbPTOf(tc.m, tc.t, tc.p, tc.n, AnyLengthWord)
		})
		if allocs != 1 {
			t.Errorf("%s allocated %.0f times, want 1", tc.name, allocs)
		}
	}
}

// TestVerbPTConcurrentSafe exercises the package-level VerbPT from many goroutines.
// Run with -race, it demonstrates that the package-level surface is safe for
// concurrent use, as documented: its shared global source delegates each draw to
// the concurrency-safe global math/rand/v2 generator.
func TestVerbPTConcurrentSafe(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 4000
	done := make(chan struct{}, goroutines)
	for g := 0; g < goroutines; g++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for i := 0; i < perGoroutine; i++ {
				if VerbPT() == "" {
					t.Error("VerbPT returned an empty string")
					return
				}
			}
		}()
	}
	for g := 0; g < goroutines; g++ {
		<-done
	}
}

// ---------------------------------------------------------------------------
// Verb radical onset restriction
// ---------------------------------------------------------------------------

// TestVerbRadicalOnsetInventoryExcludesSensitive verifies that the verb radical
// onset inventory excludes exactly the frontness-sensitive onsets (c, g, ç, qu, gu)
// and keeps the frontness-insensitive clusters (cr, cl, gr, gl) and every other
// onset, so a stripped radical stays orthographically valid before any desinence
// vowel.
func TestVerbRadicalOnsetInventoryExcludesSensitive(t *testing.T) {
	forms := map[string]bool{}
	for i := range verbRadicalOnsetInv.forms {
		forms[verbRadicalOnsetInv.forms[i].form] = true
	}
	for _, bad := range []string{"c", "g", "ç", "qu", "gu", ""} {
		if forms[bad] {
			t.Errorf("verb radical onset inventory must not contain %q", bad)
		}
	}
	for _, good := range []string{"cr", "cl", "gr", "gl", "p", "b", "t", "d", "r", "l", "m", "n", "s", "lh", "nh", "ch"} {
		if !forms[good] {
			t.Errorf("verb radical onset inventory must contain %q", good)
		}
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

// BenchmarkVerbPT measures the default package-level verb path over the global
// source (random mood, tense, person, number and length).
func BenchmarkVerbPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = VerbPT()
	}
}

// BenchmarkVerbPTOfPresentIndicative measures the finite present-indicative path and
// its allocation budget (exactly one string allocation).
func BenchmarkVerbPTOfPresentIndicative(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.VerbPTOf(Indicative, Present, AnyPerson, AnyNumber, AnyLengthWord)
	}
}

// BenchmarkVerbPTOfImperfectSubjunctive measures the accented imperfect-subjunctive
// first-plural path (a multi-byte baked accent in the desinence).
func BenchmarkVerbPTOfImperfectSubjunctive(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.VerbPTOf(Subjunctive, Imperfect, First, Plural, AnyLengthWord)
	}
}

// BenchmarkVerbPTOfInfinitive measures the non-finite infinitive path.
func BenchmarkVerbPTOfInfinitive(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.VerbPTOf(Infinitive, AnyTense, AnyPerson, AnyNumber, AnyLengthWord)
	}
}
