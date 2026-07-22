package gengo

import (
	"math/rand/v2"
	"strings"
	"unicode/utf8"
)

// This file establishes the WordsPT verb-conjugation FOUNDATION: the public verb
// flexion enums (Mood, Tense, Person) shared by every verb slot, and the internal,
// reusable regular-verb machinery for the NON-FINITE forms of all three
// conjugations — the impersonal infinitive, the inflected personal infinitive, the
// gerund and the participle. Later tasks build the finite indicative, subjunctive
// and imperative forms, and the public VerbPT/VerbPTOf API, on top of this
// machinery; only the enums are public here.
//
// Every generated verb is REGULAR (pt-PT pseudo-verbs are always regular), so the
// regular paradigm of the chosen conjugation applies uniformly and the forms are
// correct-by-construction: no candidate is generated and rejected. A verb form is
// assembled as (radical) + (thematic vowel) + (desinence), exactly the classical
// decomposition of a Portuguese verb (Cunha & Cintra), where:
//
//   - the RADICAL is a bare pt-PT stem reusing the Sprint 7 syllabic sampler
//     (through sampleNounStem, wordspt_noun.go), terminated by a sampled consonant
//     onset that carries the following thematic vowel — so the radical/theme
//     junction is a legal onset+nucleus and never a vowel hiatus;
//   - the THEMATIC VOWEL is the conjugation's theme vowel (a, e, i for the first,
//     second and third conjugations), realized as the nucleus of that final
//     syllable; the participle takes its own theme vowel (a for the first
//     conjugation, i for the second and third: falado, comido, partido);
//   - the DESINENCE is a fixed, ASCII string suffix from the authoritative Cunha &
//     Cintra conjugation tables (the infinitive marker r, the gerund -ndo, the
//     participle -do, and the personal-infinitive person endings), appended after
//     the thematic vowel.
//
// No non-finite form carries a graphic accent on the radical — the impersonal and
// personal infinitives, the gerund and the participle are all default-stressed
// (falar, falares, falando, falado) — so the non-finite assembly writes the bare
// stem and never invokes the Layer 3 accentuation (wordspt_accent.go). Because the
// desinences are pure ASCII and the thematic vowel joins the radical as a legal
// nucleus, the assembled form is validated against the Sprint 7 STRING oracles
// (cedilla, diaeresis, grave, valid lowercase UTF-8), like the plural, the
// superlative and the -mente adverb.
//
// Persons: current pt-PT usage, WITHOUT the archaic second person plural ("vós").
// The paradigm has exactly five person/number slots — 1s, 2s, 3s, 1p, 3p — and a
// requested second person plural normalizes to the third person plural (the vocês
// address, Normalization rule N1 of the specification).
//
// Allocation: the syllable buffer is a fixed-size array on the caller's stack, so
// it does not allocate on the heap; a rejected length attempt is measured from the
// buffer with no allocation. A non-finite verb therefore performs exactly one
// string allocation (assembleNonFiniteVerb writes the bare stem and the desinence
// into a single pre-sized builder).
//
// Randomness surfaces: the shared core draws only from an injectable *rand.Rand, so
// a seeded Generator is fully reproducible; the package-level verb functions (a
// later task) will feed it the global source through wordsPTGlobalRand.
//
// Sources: Cunha, C. & Cintra, L., "Nova Gramática do Português Contemporâneo" (the
// three regular conjugations, the theme vowels, and the non-finite desinence
// tables); Acordo Ortográfico da Língua Portuguesa (graphic accent, cedilla and
// nasal spelling). The conjugation sampling weights are ordinal estimates
// (flagged), consistent with the other WordsPT weights; they change only how often
// each conjugation appears, never the conformance of an output.

// ---------------------------------------------------------------------------
// Public verb flexion enums
// ---------------------------------------------------------------------------

// Mood is the grammatical mood option of a generated pt-PT verb. Its zero value,
// [AnyMood], means "choose a valid mood at random". The non-finite moods
// ([Infinitive], [Gerund], [Participle]) take no tense; the finite moods
// ([Indicative], [Subjunctive], [ImperativeAffirmative], [ImperativeNegative]) are
// built by later tasks.
type Mood uint8

const (
	// AnyMood is the zero value of [Mood] and selects a random valid mood.
	AnyMood Mood = iota
	// Indicative selects the indicative mood (finite).
	Indicative
	// Subjunctive selects the subjunctive (conjuntivo) mood (finite).
	Subjunctive
	// ImperativeAffirmative selects the affirmative imperative mood (finite).
	ImperativeAffirmative
	// ImperativeNegative selects the negative imperative mood. In pt-PT the negated
	// command is periphrastic ("não" + present subjunctive); WordsPT, a single-word
	// generator, produces only the verb form, without the negating particle.
	ImperativeNegative
	// Infinitive selects the infinitive mood: the impersonal infinitive (falar) or,
	// for a concrete person and number, the inflected personal infinitive.
	Infinitive
	// Gerund selects the gerund (gerúndio): the invariable -ndo form (falando).
	Gerund
	// Participle selects the past participle: the invariable masculine singular form
	// (falado, comido, partido).
	Participle
)

// Tense is the grammatical tense option of a generated pt-PT verb. Its zero value,
// [AnyTense], means "choose a valid tense at random for the resolved mood". Tense
// applies only to the finite moods; on the non-finite and imperative moods a
// requested tense is ignored.
type Tense uint8

const (
	// AnyTense is the zero value of [Tense] and selects a random valid tense.
	AnyTense Tense = iota
	// Present selects the present tense.
	Present
	// Imperfect selects the imperfect past tense.
	Imperfect
	// Preterite selects the preterite (perfect past) tense.
	Preterite
	// Future selects the future tense.
	Future
	// Conditional selects the conditional tense.
	Conditional
)

// Person is the grammatical person option of a generated pt-PT verb. Its zero
// value, [AnyPerson], means "choose a valid person at random". WordsPT follows
// current pt-PT usage without the archaic second person plural ("vós"): a second
// person plural request normalizes to the third person plural (see the package
// documentation, Normalization rule N1).
type Person uint8

const (
	// AnyPerson is the zero value of [Person] and selects a random valid person.
	AnyPerson Person = iota
	// First selects the first person (eu / nós).
	First
	// Second selects the second person (tu).
	Second
	// Third selects the third person (singular ele/ela, plural elas/vocês).
	Third
)

// ---------------------------------------------------------------------------
// Regular conjugation classes
// ---------------------------------------------------------------------------

// verbConjugation describes one of the three regular pt-PT conjugation classes. It
// carries the two theme vowels a form needs: themeVowel is the citation theme vowel
// used by the infinitive, the personal infinitive and the gerund (a, e, i), while
// participleTheme is the theme vowel of the regular past participle, which the
// second conjugation realizes as i (comido), not e — so it is a for the first
// conjugation and i for the second and third (falado, comido, partido). Both are
// single lowercase ASCII vowels.
type verbConjugation struct {
	// label is the human-readable conjugation name ("-ar", "-er", "-ir"); it is used
	// in documentation and tests, never in generation.
	label string
	// themeVowel is the theme vowel of the infinitive, personal infinitive and
	// gerund (a, e, i).
	themeVowel string
	// participleTheme is the theme vowel of the regular past participle (a for the
	// first conjugation, i for the second and third).
	participleTheme string
}

// The three regular conjugation classes, identified by their infinitive ending.
// Model verbs: falar (-ar), comer (-er), partir (-ir).
var (
	// conjAr is the first conjugation (-ar), theme vowel a, participle -ado (falar,
	// falando, falado). It is by far the most productive class in pt-PT.
	conjAr = verbConjugation{label: "-ar", themeVowel: "a", participleTheme: "a"}
	// conjEr is the second conjugation (-er), theme vowel e, participle -ido (comer,
	// comendo, comido — the participle theme is i, not e).
	conjEr = verbConjugation{label: "-er", themeVowel: "e", participleTheme: "i"}
	// conjIr is the third conjugation (-ir), theme vowel i, participle -ido (partir,
	// partindo, partido).
	conjIr = verbConjugation{label: "-ir", themeVowel: "i", participleTheme: "i"}
)

// weightedConjugation pairs a conjugation class with its ordinal sampling weight.
type weightedConjugation struct {
	conjugation *verbConjugation
	weight      uint32
}

// verbConjugations lists the three regular conjugations with their sampling
// weights. Weights are ordinal estimates (flagged): the first conjugation (-ar) is
// the dominant, and the only productive, class in pt-PT — new verbs are coined in
// -ar — so it far outweighs the second (-er) and third (-ir), which are of similar,
// smaller size. The values change only how often each conjugation is drawn, never
// the conformance of an output.
var verbConjugations = []weightedConjugation{
	{&conjAr, 75},
	{&conjEr, 13},
	{&conjIr, 12},
}

// verbConjugationsTotal is the total conjugation weight, the exclusive upper bound
// for the single bounded draw in [pickConjugation]. It is computed once, at package
// initialization.
var verbConjugationsTotal = func() uint32 {
	var total uint32
	for i := range verbConjugations {
		total += verbConjugations[i].weight
	}
	return total
}()

// pickConjugation draws a regular conjugation from the injectable source r using
// the ordinal weights, in a single bounded draw and a short linear scan (no
// allocation, no rejection). The first conjugation dominates, matching pt-PT.
func pickConjugation(r *rand.Rand) *verbConjugation {
	v := r.Uint32N(verbConjugationsTotal)
	var acc uint32
	for i := range verbConjugations {
		acc += verbConjugations[i].weight
		if v < acc {
			return verbConjugations[i].conjugation
		}
	}
	return &conjAr // unreachable: v < verbConjugationsTotal always matches an entry
}

// ---------------------------------------------------------------------------
// Person/number paradigm slots (no vós) and non-finite forms
// ---------------------------------------------------------------------------

// verbPerson is an internal person/number slot of the pt-PT verb paradigm. WordsPT
// omits the archaic second person plural ("vós"), so there are exactly five slots.
// The order matches the conjugation tables (1s, 2s, 3s, 1p, 3p) and indexes the
// desinence tables directly.
type verbPerson uint8

const (
	verbP1s verbPerson = iota // eu (first singular)
	verbP2s                   // tu (second singular)
	verbP3s                   // ele / ela / você (third singular)
	verbP1p                   // nós (first plural)
	verbP3p                   // elas / vocês (third plural)
)

// resolveVerbPerson maps a public [Person] and [Number] to a concrete verbPerson
// paradigm slot, drawing from the injectable source r for any unspecified trait.
// AnyPerson (or any undefined value) resolves to a uniformly random First, Second
// or Third; AnyNumber resolves through [resolveNumber]. A second person plural
// request normalizes to the third person plural (Normalization rule N1: current
// pt-PT expresses second person plural address with vocês, which governs third
// person plural agreement). The function never errors and never panics.
func resolveVerbPerson(r *rand.Rand, p Person, n Number) verbPerson {
	rp := p
	switch rp {
	case First, Second, Third:
	default:
		// First, Second or Third with equal probability (a switch avoids any
		// numeric narrowing of the draw into the Person type).
		switch r.Uint32N(3) {
		case 0:
			rp = First
		case 1:
			rp = Second
		default:
			rp = Third
		}
	}
	if resolveNumber(r, n) == Singular {
		switch rp {
		case First:
			return verbP1s
		case Second:
			return verbP2s
		default:
			return verbP3s
		}
	}
	// Plural. Second person plural (vós) does not exist and normalizes to third
	// person plural; the first person plural is nós, the third the vocês form.
	if rp == First {
		return verbP1p
	}
	return verbP3p
}

// nonFiniteForm identifies which non-finite paradigm slot to assemble.
type nonFiniteForm uint8

const (
	// formImpersonalInfinitive is the impersonal infinitive (falar, comer, partir).
	formImpersonalInfinitive nonFiniteForm = iota
	// formPersonalInfinitive is the inflected personal infinitive (falar, falares,
	// falar, falarmos, falarem), selected by a verbPerson.
	formPersonalInfinitive
	// formGerund is the invariable gerund (falando, comendo, partindo).
	formGerund
	// formParticiple is the invariable past participle (falado, comido, partido).
	formParticiple
)

// ---------------------------------------------------------------------------
// Non-finite desinence data (Cunha & Cintra conjugation tables)
// ---------------------------------------------------------------------------

// infinitiveMarker is the infinitive desinence attached after the thematic vowel:
// -ar, -er, -ir. The personal infinitive appends the person endings after it.
const infinitiveMarker = "r"

// gerundMarker is the gerund desinence attached after the thematic vowel: -ando,
// -endo, -indo.
const gerundMarker = "ndo"

// participleMarker is the regular past-participle desinence attached after the
// (participle) thematic vowel: -ado, -ido.
const participleMarker = "do"

// personalInfinitiveSuffixes are the complete suffixes of the inflected personal
// infinitive, appended after the thematic vowel, indexed by verbPerson. Each begins
// with the infinitive marker r; the second singular adds -es, the first plural
// -mos and the third plural -em, while the first and third singular coincide with
// the impersonal infinitive (falar, falares, falar, falarmos, falarem). Source:
// Cunha & Cintra, personal-infinitive paradigm.
var personalInfinitiveSuffixes = [5]string{
	verbP1s: infinitiveMarker,         // -ar
	verbP2s: infinitiveMarker + "es",  // -ares
	verbP3s: infinitiveMarker,         // -ar
	verbP1p: infinitiveMarker + "mos", // -armos
	verbP3p: infinitiveMarker + "em",  // -arem
}

// themeVowelFor returns the thematic vowel a non-finite form realizes for
// conjugation c: the participle uses the conjugation's participle theme (a, i, i),
// every other non-finite form uses the citation theme vowel (a, e, i).
func themeVowelFor(c *verbConjugation, form nonFiniteForm) string {
	if form == formParticiple {
		return c.participleTheme
	}
	return c.themeVowel
}

// nonFiniteSuffix returns the desinence appended after the thematic vowel for a
// non-finite form: the infinitive marker for the impersonal infinitive, the
// person-specific suffix for the personal infinitive, the gerund marker, or the
// participle marker. The returned string is a compile-time constant (no allocation).
func nonFiniteSuffix(form nonFiniteForm, p verbPerson) string {
	switch form {
	case formPersonalInfinitive:
		return personalInfinitiveSuffixes[p]
	case formGerund:
		return gerundMarker
	case formParticiple:
		return participleMarker
	default: // formImpersonalInfinitive
		return infinitiveMarker
	}
}

// conjugateNonFinite assembles a non-finite verb form from a bare RADICAL string in
// exactly one allocation: radical + thematic vowel + desinence. It is the
// paradigm-exact assembler — given the model radical "fal", "com" or "part", the
// conjugation, the form and (for the personal infinitive) the person, it produces
// the textbook form (falar, comeres, partindo, falado, comido, partido) — and is
// what the unit tests verify against the authoritative tables. The generation path
// uses the equivalent stem-based [assembleNonFiniteVerb]; both share
// [nonFiniteSuffix] and [themeVowelFor], so they cannot diverge. The result is
// lowercase and valid UTF-8, because the radical, the thematic vowel and every
// desinence are lowercase, valid UTF-8. p is ignored for every form except the
// personal infinitive.
func conjugateNonFinite(radical string, c *verbConjugation, form nonFiniteForm, p verbPerson) string {
	theme := themeVowelFor(c, form)
	suffix := nonFiniteSuffix(form, p)

	var b strings.Builder
	b.Grow(len(radical) + len(theme) + len(suffix))
	b.WriteString(radical)
	b.WriteString(theme)
	b.WriteString(suffix)
	return b.String()
}

// ---------------------------------------------------------------------------
// Thematic-vowel endings for radical sampling
// ---------------------------------------------------------------------------

// The thematic-vowel endings reuse the noun ending machinery (wordspt_noun.go) to
// sample the radical's final consonant together with the thematic vowel it carries:
// each is a single-syllable [nounEnding] whose onset is drawn at build time (so the
// radical always ends in a consonant before the thematic vowel — no vowel hiatus)
// and whose nucleus is the thematic vowel. The front/back onset inventory follows
// the vowel: e and i are front (they license qu, gu and forbid ç), a is back (it
// licenses ç and forbids qu, gu). stressFromEnd is irrelevant because non-finite
// forms are assembled bare, without invoking the accentuation layer. affixRunes is
// two — a sampled onset (one rune, the common case) plus the thematic vowel.
var (
	// themeEndingA carries the back thematic vowel a (first conjugation, and the
	// first-conjugation participle).
	themeEndingA = nounEnding{
		label: "-a-", syllables: []syllable{{nucleus: "a"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 0, affixRunes: 2, pluralRuneDelta: 0,
	}
	// themeEndingE carries the front thematic vowel e (second-conjugation infinitive
	// and gerund).
	themeEndingE = nounEnding{
		label: "-e-", syllables: []syllable{{nucleus: "e"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 0, affixRunes: 2, pluralRuneDelta: 0,
	}
	// themeEndingI carries the front thematic vowel i (third-conjugation infinitive
	// and gerund, and the second- and third-conjugation participle).
	themeEndingI = nounEnding{
		label: "-i-", syllables: []syllable{{nucleus: "i"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 0, affixRunes: 2, pluralRuneDelta: 0,
	}
)

// themeEndingFor returns the thematic-vowel ending for a thematic vowel, so that
// [sampleNounStem] draws a radical terminated by a consonant onset carrying that
// vowel. The vowel is always one of a, e or i (from [themeVowelFor]); a is the
// default.
func themeEndingFor(themeVowel string) *nounEnding {
	switch themeVowel {
	case "e":
		return &themeEndingE
	case "i":
		return &themeEndingI
	default: // "a"
		return &themeEndingA
	}
}

// assembleNonFiniteVerb renders a non-finite verb form from a bare combined stem
// (the radical, terminated by the sampled onset and the thematic vowel) and a
// desinence suffix, in exactly one allocation. It writes every syllable straight
// from the bare (unaccented) stem — so the form carries no graphic accent, as every
// regular non-finite form is default-stressed — and appends the ASCII desinence.
// The builder is pre-sized to the exact byte length so the writes never reallocate
// and Builder.String returns the bytes without a second copy. The result is
// lowercase, valid UTF-8.
func assembleNonFiniteVerb(stem syllabicStem, suffix string) string {
	syllables := stem.syllables

	total := len(suffix)
	for i := range syllables {
		total += len(syllables[i].onset) + len(syllables[i].nucleus) + len(syllables[i].coda)
	}

	var b strings.Builder
	b.Grow(total)
	for i := range syllables {
		b.WriteString(syllables[i].onset)
		b.WriteString(syllables[i].nucleus)
		b.WriteString(syllables[i].coda)
	}
	b.WriteString(suffix)
	return b.String()
}

// ---------------------------------------------------------------------------
// Shared generation core (non-finite forms)
// ---------------------------------------------------------------------------

// nonFiniteVerbCore is the single generation engine for the non-finite verb forms.
// It draws only from r, so a seeded Generator is fully reproducible, and reuses a
// caller-provided syllable buffer (backed by a stack array) so that only the final
// assembly allocates.
//
// It resolves the conjugation (weighted, first-conjugation dominant), resolves the
// person/number slot for the personal infinitive (Person and Number are ignored for
// the impersonal infinitive, the gerund and the participle), sizes the leading stem
// for the thematic vowel plus the desinence (normalizing the length category up to
// the minimum viable length), then samples a radical whose assembled form falls in
// the category window. The output is always lowercase, valid UTF-8,
// orthographically well-formed, and unaccented.
func nonFiniteVerbCore(r *rand.Rand, buf []syllable, form nonFiniteForm, p Person, n Number, l LengthTypeWords) string {
	c := pickConjugation(r)

	vp := verbP1s // the impersonal infinitive, gerund and participle ignore the person
	if form == formPersonalInfinitive {
		vp = resolveVerbPerson(r, p, n)
	}

	theme := themeVowelFor(c, form)
	end := themeEndingFor(theme)
	suffix := nonFiniteSuffix(form, vp)
	suffixRunes := utf8.RuneCountInString(suffix)

	// The thematic vowel (in the ending) and the desinence are the fixed affix cost
	// on top of the bare leading stem, for both length normalization and stem sizing.
	affixCost := end.affixRunes + suffixRunes
	minViable := affixCost + nounStemFloor
	lo, hi := charRangeOf(normalizeLength(l, minViable))

	var stem syllabicStem
	for attempt := 0; attempt < nounMaxLengthAttempts; attempt++ {
		leadingCount := stemSyllablesFor(r, l, minViable, affixCost)
		if leadingCount+len(end.syllables) > cap(buf) {
			leadingCount = cap(buf) - len(end.syllables)
			if leadingCount < 1 {
				leadingCount = 1
			}
		}
		stem = sampleNounStem(r, buf, leadingCount, end)
		buf = stem.syllables

		runes := stemRuneLen(stem.syllables) + suffixRunes
		if runes >= lo && runes <= hi {
			break
		}
	}

	return assembleNonFiniteVerb(stem, suffix)
}
