package gengo

import (
	"math/rand/v2"
	"strings"
	"unicode/utf8"
)

// This file implements the public WordsPT pt-PT pseudo-verb generators: the
// package-level VerbPT and VerbPTOf, plus their mirrored *Generator methods. It
// wires the regular-conjugation engine built in wordspt_verb.go (tasks #28, #29
// and #30) into the public API defined in specification/wordspt.md.
//
//	VerbPT()                                             -> a random valid verb form
//	VerbPTOf(m Mood, t Tense, p Person, n Number, l LengthTypeWords)
//
// VerbPTOf addresses one paradigm slot: the mood m, the tense t, the person p and
// the number n, produced within length category l. Every option's Any zero value
// selects a random valid value for the resolved mood, and every grammatically
// impossible combination normalizes to the nearest valid form (the deterministic
// matrix of specification/wordspt-open-classes.md, Normalization). The function
// never returns an error and never panics.
//
// Assembly model
//
// A verb form is (radical) + (desinence), where the desinence carries the mood,
// tense, person, number marking and any graphic accent, exactly as the #28..#30
// paradigm-exact assemblers (conjugateNonFinite, conjugateIndicative,
// conjugateSubjunctive, conjugateImperative) define it. verbCore samples ONE bare
// radical and appends the resolved desinence, so the output is identical to what
// those assemblers produce for the same radical — with a single string
// allocation, because the radical is never materialized as a separate string:
// assembleVerbForm writes the radical's syllables and the desinence into one
// pre-sized builder.
//
// The radical is a bare, consonant-final pt-PT stem (fal-, com-, part- style),
// sampled hiatus-free by reusing sampleNounStem (wordspt_noun.go) with a
// verb-specific "theme" ending: a single sampled onset carrying a placeholder
// thematic vowel that is discarded, so the radical always ends in a real
// consonant onset (never a vowel), and no radical/desinence junction is a vowel
// hiatus. The onset is drawn from an inventory that excludes the
// frontness-sensitive onsets (see [frontnessSensitiveOnset]).
//
// Present-tense radical accentuation (resolved)
//
// The present indicative and present subjunctive stress the RADICAL (falo, fala;
// fale, coma). For a regular verb these forms are always default paroxytones
// ending in -o, -a, -e, -as, -es, -am, -em, -amos, -emos or -imos, and by the
// pt-PT paroxytone rule (Acordo Ortográfico) a default paroxytone ending takes NO
// graphic accent — regardless of which vowel is stressed. So a bare radical plus
// an unaccented present desinence is orthographically correct for every sampled
// radical, and the present path needs no accentuation (no Layer 3 call), exactly
// like the imperfect/preterite/future/conditional whose accents are baked into
// their desinences. This is verified empirically over a large sample in the tests.
//
// Randomness surfaces. The package-level VerbPT/VerbPTOf draw from the global,
// automatically seeded math/rand/v2 source (through wordsPTGlobalRand); the
// (*Generator) methods draw from the generator's own seeded source. Both call the
// single shared core verbCore with an injectable *rand.Rand, so there is no
// duplication of generation logic and a seeded Generator is fully reproducible.
//
// Sources: Cunha, C. & Cintra, L., "Nova Gramática do Português Contemporâneo"
// (the three regular conjugations and their paradigm tables); Acordo Ortográfico
// da Língua Portuguesa (graphic accent, cedilla and nasal spelling). The mood and
// infinitive-form sampling weights are ordinal estimates (flagged), consistent
// with the other WordsPT weights; they change only how often each valid form
// appears, never the conformance of an output.

// ---------------------------------------------------------------------------
// Verb radical onset inventory (frontness-insensitive)
// ---------------------------------------------------------------------------

// frontnessSensitiveOnset reports whether an onset's pt-PT spelling or legality
// depends on the frontness of the following vowel: c and g read soft before a
// front vowel (e, i) and hard before a back or central vowel (a, o, u); ç is legal
// only before a back or central vowel; qu and gu only before a front vowel.
//
// A finite verb desinence attaches directly to the BARE radical and may begin with
// a vowel of either frontness (the first-conjugation preterite -ei and present
// subjunctive -e endings are front; the present -o and the second/third
// conjugation present subjunctive -a endings are back). A radical ending in one of
// these onsets would therefore need an orthographic alternation to stay correct
// (caçar -> cacei, ficar -> fiquei, chegar -> cheguei, vencer -> venço, erguer ->
// ergo). To keep every regular form correct by construction, with a single
// allocation and no rejection, the verb radical never ends in one of these onsets.
// The obstruent+liquid clusters cr, cl, gr and gl are NOT sensitive (their c/g is
// always hard before the liquid), so they remain available.
//
// Flagged simplification (pseudo-word scope, consistent with the #28 "vowel-final
// roots are not generated" simplification): pseudo-verbs whose radical ends in a
// standalone c, g, ç, qu or gu (the caçar/ficar/vencer/erguer classes with
// orthographic alternations) are not generated. This trades a small amount of
// radical variety for guaranteed orthographic correctness across the whole
// paradigm.
func frontnessSensitiveOnset(form string) bool {
	switch form {
	case "c", "g", "ç", "qu", "gu":
		return true
	}
	return false
}

// notFrontnessSensitiveOnset is the deriveInventory predicate that keeps every
// onset except the frontness-sensitive ones.
func notFrontnessSensitiveOnset(form string) bool { return !frontnessSensitiveOnset(form) }

// verbRadicalOnsetInv samples the final (pre-desinence) onset of a verb radical:
// any single, digraph or cluster onset except the frontness-sensitive c, g, ç, qu
// and gu. Because those are the only onsets whose realization depends on the
// following vowel's frontness, the remaining onsets are valid before any desinence
// vowel, so a single inventory serves every conjugation (the front/back split the
// noun endings use collapses once the sensitive onsets are removed). The empty
// onset is excluded (weight 0), so the radical always ends in a consonant.
var verbRadicalOnsetInv = deriveInventory(0, notFrontnessSensitiveOnset, onsetSingles, onsetClusters)

// verbRadicalEnding is the "theme" ending sampleNounStem uses to terminate a verb
// radical with a sampled consonant onset. Its nucleus is a placeholder that
// assembleVerbForm discards (the real thematic vowel, where a form keeps one, is
// part of the desinence), so only the sampled onset survives into the radical.
// onsetInvOverride restricts that onset to the frontness-insensitive inventory.
// stressFromEnd is irrelevant (verb forms are assembled bare, without the Layer 3
// accentuation). affixRunes counts the sampled onset plus the placeholder vowel,
// matching the #28 theme endings, though verbCore computes its own affix cost.
var verbRadicalEnding = nounEnding{
	label: "-radical-", syllables: []syllable{{nucleus: "a"}},
	sampleOnset: true, onsetInvOverride: &verbRadicalOnsetInv,
	stressFromEnd: 0, affixRunes: 2, pluralRuneDelta: 0,
}

// ---------------------------------------------------------------------------
// Mood resolution (Any -> weighted random valid mood)
// ---------------------------------------------------------------------------

// weightedMood pairs a concrete mood with its ordinal sampling weight.
type weightedMood struct {
	mood   Mood
	weight uint32
}

// verbMoods lists the seven concrete moods with their sampling weights, used to
// resolve AnyMood. Weights are ordinal estimates (flagged): the finite indicative
// is by far the most frequent verb mood in running text, so it dominates; the
// subjunctive and the non-finite forms follow; the imperative is the least
// frequent. The values change only how often each mood is drawn, never the
// conformance of an output.
var verbMoods = []weightedMood{
	{Indicative, 50},
	{Subjunctive, 15},
	{Infinitive, 12},
	{Participle, 8},
	{Gerund, 6},
	{ImperativeAffirmative, 5},
	{ImperativeNegative, 4},
}

// verbMoodsTotal is the total mood weight, the exclusive upper bound for the single
// bounded draw in [resolveMood]. It is computed once, at package initialization.
var verbMoodsTotal = func() uint32 {
	var total uint32
	for i := range verbMoods {
		total += verbMoods[i].weight
	}
	return total
}()

// resolveMood resolves m to a concrete [Mood], drawing from the injectable source r
// for an unspecified request. A concrete mood (Indicative, Subjunctive,
// ImperativeAffirmative, ImperativeNegative, Infinitive, Gerund or Participle) is
// returned unchanged; AnyMood (and any undefined value) becomes one of the seven
// moods with the ordinal weighting above (finite indicative favored). The draw is a
// single bounded selection and a short linear scan, without allocation or rejection.
func resolveMood(r *rand.Rand, m Mood) Mood {
	switch m {
	case Indicative, Subjunctive, ImperativeAffirmative, ImperativeNegative,
		Infinitive, Gerund, Participle:
		return m
	default:
		v := r.Uint32N(verbMoodsTotal)
		var acc uint32
		for i := range verbMoods {
			acc += verbMoods[i].weight
			if v < acc {
				return verbMoods[i].mood
			}
		}
		return Indicative // unreachable: v < verbMoodsTotal always matches an entry
	}
}

// ---------------------------------------------------------------------------
// Infinitive form resolution (impersonal vs inflected personal)
// ---------------------------------------------------------------------------

// infinitiveImpersonalWeight and infinitivePersonalWeight are the ordinal weights
// (flagged) of the impersonal versus the inflected personal infinitive when the
// Infinitive mood is requested with BOTH person and number unspecified. The
// impersonal infinitive is the common citation form, so it is the default; a random
// personal infinitive is produced the rest of the time (specification
// wordspt-open-classes.md, "Selecting the infinitive form"). When a concrete person
// OR number is specified, the personal infinitive is always produced (this table is
// not consulted).
const (
	infinitiveImpersonalWeight = 3
	infinitivePersonalWeight   = 2
)

// resolveInfinitiveForm resolves the Infinitive mood to a concrete non-finite form
// and (for the personal infinitive) a person slot, per the specification's
// infinitive-selection rule. When a concrete person or number is requested, it is
// the inflected personal infinitive for the resolved person/number (which, for the
// first and third person singular, equals the impersonal infinitive). When both
// person and number are Any, it is a weighted choice: the impersonal infinitive
// (the common default) or a random inflected personal infinitive. It draws from r
// for every unspecified choice and never errors or panics.
func resolveInfinitiveForm(r *rand.Rand, p Person, n Number) (nonFiniteForm, verbPerson) {
	personSpecified := p == First || p == Second || p == Third
	numberSpecified := n == Singular || n == Plural
	if personSpecified || numberSpecified {
		return formPersonalInfinitive, resolveVerbPerson(r, p, n)
	}
	// Both Any: impersonal (default) or a random personal infinitive.
	if r.Uint32N(infinitiveImpersonalWeight+infinitivePersonalWeight) < infinitiveImpersonalWeight {
		return formImpersonalInfinitive, verbP1s
	}
	return formPersonalInfinitive, resolveVerbPerson(r, AnyPerson, AnyNumber)
}

// ---------------------------------------------------------------------------
// Desinence resolution (Mood -> the string appended to the bare radical)
// ---------------------------------------------------------------------------

// nonFiniteDesinences holds one conjugation's non-finite desinences (the thematic
// vowel folded into each form's suffix), precomputed once so that resolving a
// non-finite form is a table lookup rather than a per-call string concatenation.
// Materializing the desinence at run time would add a second allocation on top of
// the final assembly and break the single-allocation budget; storing the strings,
// exactly like the finite paradigm tables, keeps every verb form at one allocation.
type nonFiniteDesinences struct {
	impersonalInf string    // -ar, -er, -ir
	personalInf   [5]string // indexed by verbPerson: -ar/-ares/-ar/-armos/-arem, etc.
	gerund        string    // -ando, -endo, -indo
	participle    string    // -ado, -ido, -ido
}

// nonFiniteDesinenceRaw returns the thematic vowel plus a non-finite form's suffix
// (for example "a"+"r"="ar", "e"+"ndo"="endo", "i"+"do"="ido"). It mirrors
// conjugateNonFinite with an empty radical, so radical + nonFiniteDesinenceRaw(...)
// equals conjugateNonFinite(radical, ...). It is called only at package
// initialization, to fill the desinence tables below.
func nonFiniteDesinenceRaw(c *verbConjugation, form nonFiniteForm, p verbPerson) string {
	return themeVowelFor(c, form) + nonFiniteSuffix(form, p)
}

// buildNonFiniteDesinences precomputes conjugation c's non-finite desinence table.
func buildNonFiniteDesinences(c *verbConjugation) nonFiniteDesinences {
	d := nonFiniteDesinences{
		impersonalInf: nonFiniteDesinenceRaw(c, formImpersonalInfinitive, verbP1s),
		gerund:        nonFiniteDesinenceRaw(c, formGerund, verbP1s),
		participle:    nonFiniteDesinenceRaw(c, formParticiple, verbP1s),
	}
	for _, vp := range []verbPerson{verbP1s, verbP2s, verbP3s, verbP1p, verbP3p} {
		d.personalInf[vp] = nonFiniteDesinenceRaw(c, formPersonalInfinitive, vp)
	}
	return d
}

// The non-finite desinence tables of the three regular conjugations, precomputed at
// package initialization (mirroring the finite indicative and subjunctive tables).
var (
	nonFiniteDesinencesAr = buildNonFiniteDesinences(&conjAr)
	nonFiniteDesinencesEr = buildNonFiniteDesinences(&conjEr)
	nonFiniteDesinencesIr = buildNonFiniteDesinences(&conjIr)
)

// nonFiniteDesinencesFor returns conjugation c's non-finite desinence table,
// selected by its theme vowel (a, e or i), mirroring indicativeParadigmFor.
func nonFiniteDesinencesFor(c *verbConjugation) *nonFiniteDesinences {
	switch c.themeVowel {
	case "e":
		return &nonFiniteDesinencesEr
	case "i":
		return &nonFiniteDesinencesIr
	default: // "a"
		return &nonFiniteDesinencesAr
	}
}

// verbDesinence resolves every option that applies to the (already resolved) mood m
// and returns the desinence string appended to the bare radical. It draws from r
// for each Any option, applying the per-mood normalizations already implemented in
// the #28..#30 resolvers (resolveIndicativeTense for the indicative all-tense case,
// resolveSubjunctiveTense for the subjunctive N3 tense normalization,
// resolveVerbPerson for the N1 second-person-plural normalization, and
// resolveImperativePerson for the additional N2 first-person-singular
// normalization). Tense is ignored for the non-finite and imperative moods (N3).
// The returned string is a package-level table entry or a two-part concatenation of
// them, and by construction it makes radical + verbDesinence(...) equal the
// corresponding #28..#30 paradigm-exact assembler output.
func verbDesinence(r *rand.Rand, m Mood, c *verbConjugation, t Tense, p Person, n Number) string {
	switch m {
	case Indicative:
		rt := resolveIndicativeTense(r, t)
		vp := resolveVerbPerson(r, p, n)
		return indicativeSuffix(c, rt, vp)
	case Subjunctive:
		rt := resolveSubjunctiveTense(r, t)
		vp := resolveVerbPerson(r, p, n)
		return subjunctiveSuffix(c, rt, vp)
	case ImperativeAffirmative:
		vp := resolveImperativePerson(r, p, n)
		return imperativeSuffix(c, imperativeAff, vp)
	case ImperativeNegative:
		vp := resolveImperativePerson(r, p, n)
		return imperativeSuffix(c, imperativeNeg, vp)
	case Gerund:
		return nonFiniteDesinencesFor(c).gerund
	case Participle:
		return nonFiniteDesinencesFor(c).participle
	case Infinitive:
		form, vp := resolveInfinitiveForm(r, p, n)
		d := nonFiniteDesinencesFor(c)
		if form == formImpersonalInfinitive {
			return d.impersonalInf
		}
		return d.personalInf[vp]
	default: // unreachable: resolveMood returns a concrete mood
		return nonFiniteDesinencesFor(c).impersonalInf
	}
}

// ---------------------------------------------------------------------------
// Single-allocation assembly
// ---------------------------------------------------------------------------

// assembleVerbForm renders a verb form from a sampled radical stem and a resolved
// desinence in exactly one allocation: the leading syllables of the stem in full,
// then the final ("theme") syllable's onset only (the radical's terminal
// consonant), then the desinence. The final syllable's placeholder nucleus is
// discarded — the real thematic vowel, where a form keeps one, is part of the
// desinence — so the radical always ends in a consonant. The builder is pre-sized
// to the exact byte length, so the writes never reallocate and Builder.String
// returns the bytes without a second copy. The result is lowercase and valid UTF-8,
// because every syllable part and the desinence are lowercase and valid UTF-8.
func assembleVerbForm(stem syllabicStem, desinence string) string {
	syllables := stem.syllables
	last := len(syllables) - 1

	total := len(desinence) + len(syllables[last].onset)
	for i := 0; i < last; i++ {
		total += len(syllables[i].onset) + len(syllables[i].nucleus) + len(syllables[i].coda)
	}

	var b strings.Builder
	b.Grow(total)
	for i := 0; i < last; i++ {
		b.WriteString(syllables[i].onset)
		b.WriteString(syllables[i].nucleus)
		b.WriteString(syllables[i].coda)
	}
	b.WriteString(syllables[last].onset) // radical's terminal consonant; nucleus discarded
	b.WriteString(desinence)
	return b.String()
}

// ---------------------------------------------------------------------------
// Shared generation core
// ---------------------------------------------------------------------------

// verbCore is the single generation engine shared by the package-level verb
// functions and the *Generator verb methods. It draws only from r, so a seeded
// Generator is fully reproducible, and reuses a caller-provided syllable buffer
// (backed by a stack array) so that only the final assembly allocates.
//
// It resolves the mood (weighted, finite-indicative favored), picks a conjugation
// (weighted, first-conjugation dominant), and resolves every option that applies to
// the mood into a single desinence string (applying the specification's N1, N2 and
// N3 normalizations through the #28..#30 resolvers). It then sizes the leading stem
// for that desinence (normalizing the length category up to the minimum viable
// length), samples a hiatus-free radical whose assembled form falls in the category
// window, and appends the desinence in one allocation. The output is always
// lowercase, valid UTF-8, orthographically well-formed, and correctly accented.
func verbCore(r *rand.Rand, buf []syllable, m Mood, t Tense, p Person, n Number, l LengthTypeWords) string {
	rm := resolveMood(r, m)
	c := pickConjugation(r)
	desinence := verbDesinence(r, rm, c, t, p, n)
	desinenceRunes := utf8.RuneCountInString(desinence)

	// The radical's terminal consonant onset (one rune, the common case) plus the
	// desinence are the fixed affix cost on top of the bare leading stem, for both
	// length normalization and stem sizing.
	affixCost := 1 + desinenceRunes
	minViable := affixCost + nounStemFloor
	lo, hi := charRangeOf(normalizeLength(l, minViable))

	end := &verbRadicalEnding
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

		last := len(stem.syllables) - 1
		runes := stemRuneLen(stem.syllables[:last]) +
			utf8.RuneCountInString(stem.syllables[last].onset) + desinenceRunes
		if runes >= lo && runes <= hi {
			break
		}
	}

	return assembleVerbForm(stem, desinence)
}

// ---------------------------------------------------------------------------
// Public API: package-level functions
// ---------------------------------------------------------------------------

// VerbPT returns a random European-Portuguese (pt-PT) pseudo-verb form in a random
// paradigm slot: a random mood, tense, person, number and length. It is exact sugar
// for VerbPTOf(AnyMood, AnyTense, AnyPerson, AnyNumber, AnyLengthWord). The result
// is a lowercase, valid-UTF-8 word that follows the regular pt-PT verb conjugation
// and graphic accentuation; it is a plausible invented word, not guaranteed to be a
// real dictionary entry.
//
// VerbPT draws from the global, automatically seeded math/rand/v2 source and is safe
// for concurrent use by multiple goroutines. For reproducible output, use a seeded
// Generator ([New] or [NewSource]) and its VerbPT method.
func VerbPT() string {
	return VerbPTOf(AnyMood, AnyTense, AnyPerson, AnyNumber, AnyLengthWord)
}

// VerbPTOf returns a random pt-PT pseudo-verb produced in the paradigm slot
// addressed by mood m, tense t, person p and number n, within length category l.
// Each option's Any zero value ([AnyMood], [AnyTense], [AnyPerson], [AnyNumber],
// [AnyLengthWord]) selects a random valid value for the resolved mood.
//
// The mood governs which of the other options apply. The finite indicative and
// subjunctive use the tense, person and number; the imperative uses the person and
// number (it has no tense selection); the infinitive uses the person and number (to
// choose between the impersonal and the inflected personal infinitive); the gerund
// and the participle are invariable (person, number and tense are ignored).
//
// Grammatically impossible combinations never error or panic; each normalizes to
// the nearest valid form, per the deterministic matrix of
// specification/wordspt-open-classes.md: a second person plural request becomes
// third person plural (N1); a first person singular imperative becomes first person
// plural (N2); a subjunctive Preterite becomes Imperfect and a subjunctive
// Conditional becomes Future (N3); a tense on a non-finite or imperative mood is
// ignored (N3); and an unviable length category is normalized upward to the
// smallest category that can hold the resolved slot (N5), never downward. The result
// is always lowercase and valid UTF-8.
//
// VerbPTOf draws from the global, automatically seeded math/rand/v2 source and is
// safe for concurrent use by multiple goroutines.
func VerbPTOf(m Mood, t Tense, p Person, n Number, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return verbCore(wordsPTGlobalRand, buf[:0], m, t, p, n, l)
}

// ---------------------------------------------------------------------------
// Public API: *Generator methods
// ---------------------------------------------------------------------------

// VerbPT is the seeded-generator equivalent of [VerbPT]. A Generator created with
// the same seed produces the same sequence of verb forms.
func (g *Generator) VerbPT() string {
	return g.VerbPTOf(AnyMood, AnyTense, AnyPerson, AnyNumber, AnyLengthWord)
}

// VerbPTOf is the seeded-generator equivalent of [VerbPTOf]. A Generator created
// with the same seed produces the same sequence of verb forms for the same
// arguments.
func (g *Generator) VerbPTOf(m Mood, t Tense, p Person, n Number, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return verbCore(g.r, buf[:0], m, t, p, n, l)
}
