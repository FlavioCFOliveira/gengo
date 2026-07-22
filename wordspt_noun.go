package gengo

import (
	"math/rand/v2"
	"unicode/utf8"
)

// This file implements the FIRST public WordsPT generators: the pt-PT
// pseudo-noun functions NounPT and NounPTOf, plus their mirrored *Generator
// methods. It composes the Sprint 7 / task-#22 pipeline into a complete word:
//
//   1. Layer 1 (wordspt_stem.go): a correct-by-construction syllabic stem.
//   2. Layer 2 (this file + wordspt_inflection.go + wordspt_plural.go): the noun
//      ending inventory, gender realization, and Portuguese plural formation.
//   3. Layer 3 (wordspt_accent.go): deterministic pt-PT graphic accentuation.
//
// A noun is built as (leading stem syllables) + (the class ending's syllables),
// assembled into ONE combined syllabic stem whose tonic index is fixed by the
// ending's stress pattern, so the shared accentuation layer places the accent
// correctly on the COMBINED word (never on the bare stem in isolation). The
// leading stem and the ending are joined by re-deriving the leading stem's final
// coda against the ending's first onset, exactly as the Layer 1 sampler joins two
// internal syllables (transitionIsValid); this keeps the whole word
// correct-by-construction, with no generate-and-reject phonotactic loop.
//
// Randomness surfaces. The package-level NounPT/NounPTOf draw from the global,
// automatically seeded math/rand/v2 source (through wordsPTGlobalRand, below);
// the (*Generator).NounPT/NounPTOf methods draw from the generator's own seeded
// source. Both call the single shared core nounCore with an injectable
// *rand.Rand, so there is no duplication of generation logic and a seeded
// Generator is fully reproducible.
//
// Allocation. The syllable buffer is a fixed-size array on the caller's stack, so
// it does not allocate on the heap; a rejected length attempt is measured from
// the buffer with no allocation. A singular noun therefore performs exactly one
// string allocation (the final assembly). An ADDITIVE plural — the regular +s (or
// +es for the -or ending), which merely appends a suffix — is also a single
// allocation: the suffix is written into the same assembling builder (see
// accentedWordSuffixed and additivePluralSuffix). A SUBSTITUTIVE plural performs
// one additional tail concatenation in the shared pluralizer (wordspt_plural.go),
// because a pt-PT plural such as -agens (-m -> -ns), -ção -> -ções or -al -> -ais
// rewrites the word's tail — carrying, for -agens, a coda cluster ("ns") that
// lies outside the single-coda syllable inventory and cannot be represented as a
// syllable — so it is a string transform applied after assembly.
//
// Sources for the ending inventory, the inherent gender of each ending and the
// stress each ending carries: Cunha, C. & Cintra, L., "Nova Gramática do
// Português Contemporâneo" (noun formation, gender of derivational suffixes,
// stress classes); Acordo Ortográfico da Língua Portuguesa (graphic accent,
// cedilla and nasal spelling). The relative sampling weights are ordinal
// estimates (flagged), consistent with the other WordsPT weights; they change
// only how often each valid ending appears, never the conformance of an output.

// ---------------------------------------------------------------------------
// Shared global randomness source for the package-level WordsPT functions
// ---------------------------------------------------------------------------

// globalSource is a math/rand/v2 Source whose Uint64 delegates to the package's
// global, automatically seeded generator (the top-level math/rand/v2 functions).
// Those functions are safe for concurrent use, and a [rand.Rand] holds no mutable
// state of its own (it is a thin wrapper over its Source), so a *rand.Rand built
// over globalSource is likewise safe for concurrent use. This is what lets the
// package-level WordsPT functions feed the shared, *rand.Rand-based generation
// core while keeping the same concurrency guarantees as the rest of gengo.
type globalSource struct{}

// Uint64 returns the next value from the global math/rand/v2 generator.
func (globalSource) Uint64() uint64 { return rand.Uint64() }

// wordsPTGlobalRand adapts the global math/rand/v2 source to the *rand.Rand type
// that the shared WordsPT core requires. It is the package-level counterpart of a
// Generator's own r, and is safe for concurrent use (see [globalSource]).
var wordsPTGlobalRand = rand.New(globalSource{})

// ---------------------------------------------------------------------------
// Noun ending inventory
// ---------------------------------------------------------------------------

// nounEnding describes one realized pt-PT noun ending: the trailing syllables it
// contributes to a word, its stress, and how it is pluralized. Each value is a
// concrete GENDER realization (for example the masculine -o and the feminine -a
// are two entries), so selecting an ending for a requested gender is a lookup in
// the matching list, with no run-time gender transformation.
//
// The syllables field holds the ending's fixed final syllables. When sampleOnset
// is true, the onset of the FIRST of those syllables is drawn at build time from
// the onset inventory the first nucleus licenses (frontOnset selects the front-
// vs back-vowel inventory), so a thematic ending such as -o realizes as a real
// final syllable (for example -to in "gato") rather than a bare vowel.
type nounEnding struct {
	// label is the human-readable ending name (for example "-ção"); it is used in
	// documentation and tests, never in generation.
	label string
	// syllables are the ending's fixed trailing syllables, in order. They already
	// satisfy the phonotactic validators; the internal and word-final positions
	// are legal by construction.
	syllables []syllable
	// sampleOnset requests a drawn onset on syllables[0] (a thematic or
	// derivational ending whose leading consonant varies), instead of the fixed
	// onset stored in syllables[0].
	sampleOnset bool
	// frontOnset selects, when sampleOnset is true, the onset inventory licensed
	// before the first nucleus: the front-vowel inventory (excludes ç) when the
	// first nucleus is front (e or i), the back-vowel inventory (excludes qu, gu)
	// otherwise.
	frontOnset bool
	// stressFromEnd is the tonic syllable's distance from the last syllable of the
	// whole word: 0 oxytone (last), 1 paroxytone (penult). Every noun ending here
	// is oxytone or paroxytone; none forms a proparoxytone.
	stressFromEnd int
	// affixRunes is the ending region's character (rune) cost, used to map a
	// length category to a leading-stem size. It counts a sampled onset as one
	// rune (the common case), so it is an estimate for length targeting, not an
	// exact width; the exact word length is enforced on the assembled word.
	affixRunes int
	// pluralRuneDelta is the number of characters the pt-PT plural adds to the
	// singular for this ending (for example +1 for a -s plural, +2 for the -es
	// plural of -or). It lets the length window be enforced on a plural word
	// without assembling it first.
	pluralRuneDelta int
	// fixedCaoPlural marks the -ção ending, whose plural is deterministically
	// -ções (the productive Latin -tionem outcome), rather than a weighted choice
	// among the -ão plural forms used for arbitrary words.
	fixedCaoPlural bool
	// onsetInvOverride optionally replaces, when sampleOnset is true, the
	// front/back onset inventory that syllables[0]'s drawn onset is taken from. It
	// is nil for every noun and adjective ending (they keep the default front/back
	// inventory selected by frontOnset), and is set only by the verb radical ending
	// (wordspt_verb_public.go), which needs an onset inventory that excludes the
	// frontness-sensitive onsets so a stripped radical stays orthographically valid
	// before any desinence vowel.
	onsetInvOverride *weightedInventory
}

// oxytone reports whether the ending stresses the last syllable of the word,
// which is the flag the pt-PT pluralizer needs to resolve stress-dependent
// plural outcomes.
func (e *nounEnding) oxytone() bool { return e.stressFromEnd == 0 }

// The realized noun endings. Each variable is a single gender realization; the
// masculine/feminine/common grouping below decides which requests can select it.
// Gender behavior (Cunha & Cintra): -o/-a, -eiro/-eira and -or/-ora inflect for
// gender; -ção, -dade and -agem are feminine; -mento is masculine; -ista is
// common (one form for either gender).
var (
	// endMascO is the masculine thematic ending -o (for example "gato"):
	// paroxytone, a single final syllable (sampled onset + o), plural +s.
	endMascO = nounEnding{
		label: "-o", syllables: []syllable{{nucleus: "o"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 2, pluralRuneDelta: 1,
	}
	// endFemA is the feminine thematic ending -a (for example "gata"):
	// paroxytone, a single final syllable (sampled onset + a), plural +s.
	endFemA = nounEnding{
		label: "-a", syllables: []syllable{{nucleus: "a"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 2, pluralRuneDelta: 1,
	}
	// endMascEiro is the masculine agent/tool ending -eiro (for example
	// "barbeiro"): paroxytone on the ei diphthong, syllables ei+ro with a sampled
	// onset on ei, plural +s.
	endMascEiro = nounEnding{
		label: "-eiro", syllables: []syllable{{nucleus: "ei"}, {onset: "r", nucleus: "o"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endFemEira is the feminine counterpart -eira (for example "barbeira").
	endFemEira = nounEnding{
		label: "-eira", syllables: []syllable{{nucleus: "ei"}, {onset: "r", nucleus: "a"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endMascOr is the masculine agent ending -or (for example "professor"):
	// oxytone, a single final syllable (sampled onset + o + coda r), plural -es.
	endMascOr = nounEnding{
		label: "-or", syllables: []syllable{{nucleus: "o", coda: "r"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 0, affixRunes: 3, pluralRuneDelta: 2,
	}
	// endFemOra is the feminine counterpart -ora (for example "professora"):
	// paroxytone, syllables o+ra with a sampled onset on o, plural +s.
	endFemOra = nounEnding{
		label: "-ora", syllables: []syllable{{nucleus: "o"}, {onset: "r", nucleus: "a"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endMascMento is the masculine result ending -mento (for example
	// "movimento"): paroxytone, syllables men+to, plural +s.
	endMascMento = nounEnding{
		label: "-mento", syllables: []syllable{{onset: "m", nucleus: "e", coda: "n"}, {onset: "t", nucleus: "o"}},
		sampleOnset:   false,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endFemCao is the feminine action ending -ção (for example "coração"):
	// oxytone on the nasal diphthong ão (marked by its tilde, no acute), plural
	// -ções by construction.
	endFemCao = nounEnding{
		label: "-ção", syllables: []syllable{{onset: "ç", nucleus: "ão"}},
		sampleOnset:   false,
		stressFromEnd: 0, affixRunes: 3, pluralRuneDelta: 1, fixedCaoPlural: true,
	}
	// endFemDade is the feminine quality ending -dade (for example "cidade"):
	// paroxytone, syllables da+de, plural +s.
	endFemDade = nounEnding{
		label: "-dade", syllables: []syllable{{onset: "d", nucleus: "a"}, {onset: "d", nucleus: "e"}},
		sampleOnset:   false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endFemAgem is the feminine action/collective ending -agem (for example
	// "viagem"): paroxytone, syllables a+gem, plural -agens (the -m -> -ns rule).
	endFemAgem = nounEnding{
		label: "-agem", syllables: []syllable{{nucleus: "a"}, {onset: "g", nucleus: "e", coda: "m"}},
		sampleOnset:   false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endComumIsta is the common-gender agent ending -ista (for example
	// "artista"): the same form is masculine or feminine; paroxytone, syllables
	// is+ta, plural +s.
	endComumIsta = nounEnding{
		label: "-ista", syllables: []syllable{{nucleus: "i", coda: "s"}, {onset: "t", nucleus: "a"}},
		sampleOnset:   false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
)

// weightedNounEnding pairs a realized ending with its ordinal sampling weight.
type weightedNounEnding struct {
	ending *nounEnding
	weight uint32
}

// masculineNounEndings lists the endings a masculine request may realize: the
// masculine gender-inflecting forms, the masculine-fixed -mento, and the
// common-gender -ista. Weights are ordinal estimates (flagged): the thematic -o
// is by far the most productive masculine ending.
var masculineNounEndings = []weightedNounEnding{
	{&endMascO, 50},
	{&endMascOr, 12},
	{&endMascMento, 15},
	{&endMascEiro, 8},
	{&endComumIsta, 8},
}

// feminineNounEndings lists the endings a feminine request may realize: the
// feminine gender-inflecting forms, the feminine-fixed -ção, -dade and -agem, and
// the common-gender -ista. Weights are ordinal estimates (flagged): the thematic
// -a and the productive -ção dominate.
var feminineNounEndings = []weightedNounEnding{
	{&endFemA, 40},
	{&endFemOra, 8},
	{&endFemEira, 6},
	{&endFemCao, 18},
	{&endFemDade, 12},
	{&endFemAgem, 8},
	{&endComumIsta, 8},
}

// nounStemFloor is the shortest realistic leading stem, in characters: one
// consonant-plus-vowel (CV) syllable. A bare vowel-only stem is phonotactically
// legal but degenerate, so the minimum viable length of a noun with a given
// ending is the ending's rune cost plus this floor. Using the CV floor (rather
// than a single vowel) also keeps every length category comfortably wider than an
// ending's shortest word, which is what makes the length window reliably
// reachable without a tight lower boundary.
const nounStemFloor = 2

// endingMinViable is the minimum viable character length of a noun formed with
// this ending: the ending's rune cost plus the shortest realistic stem.
func (e *nounEnding) endingMinViable() int { return e.affixRunes + nounStemFloor }

// selectNounEnding draws a noun ending that both realizes the resolved gender g
// and fits the length ceiling maxChars (its minimum viable word is at most
// maxChars characters), using the ordinal weights in a single bounded draw and a
// short linear scan (no allocation). Length-aware selection is what lets a
// concrete length category be honored strictly: a Small request, for example,
// only ever selects the short thematic endings that a Small word can contain,
// instead of a long ending that would overflow the category.
//
// When no gender-compatible ending fits maxChars (the requested category is
// shorter than the shortest ending of that gender), it falls back to the shortest
// gender-compatible ending; the caller then normalizes the length upward for that
// ending. For nouns this fallback never triggers, because the thematic -o and -a
// fit even the smallest category. g is expected to be a concrete gender
// ([Masculine] or [Feminine]) as returned by [resolveGender].
func selectNounEnding(r *rand.Rand, g Gender, maxChars int) *nounEnding {
	list := masculineNounEndings
	if g == Feminine {
		list = feminineNounEndings
	}

	// First pass: total the weight of the endings that fit, and track the
	// shortest ending as a fallback, in one scan.
	var fitTotal uint32
	shortest := list[0].ending
	for i := range list {
		e := list[i].ending
		if e.endingMinViable() <= maxChars {
			fitTotal += list[i].weight
		}
		if e.endingMinViable() < shortest.endingMinViable() {
			shortest = e
		}
	}
	if fitTotal == 0 {
		return shortest
	}

	// Second pass: weighted draw among the fitting endings only.
	v := r.Uint32N(fitTotal)
	var acc uint32
	for i := range list {
		e := list[i].ending
		if e.endingMinViable() > maxChars {
			continue
		}
		acc += list[i].weight
		if v < acc {
			return e
		}
	}
	return shortest // unreachable: v < fitTotal always matches a fitting ending
}

// ---------------------------------------------------------------------------
// Combined stem assembly (leading stem + ending)
// ---------------------------------------------------------------------------

// nounMaxSyllables bounds the syllable buffer sized on the caller's stack. A Big
// word tops out near 30 characters, which is at most ~12 syllables at the
// measured ~2.61 runes per syllable; 24 leaves generous headroom for the ending
// and keeps the array off the heap.
const nounMaxSyllables = 24

// nounMaxLengthAttempts bounds the length-window resampling. Each rejected
// attempt only resamples into the reused buffer and measures runes, allocating
// nothing, so the cap costs nothing against the single-allocation budget; it
// exists only to guarantee termination. It is set well above the number of
// attempts the tightest category (a Small word, whose four-character window
// admits only the shortest stems) needs in practice, so that the probability of
// exhausting it and returning a marginally out-of-window word is negligible.
const nounMaxLengthAttempts = 64

// sampleNounStem builds the combined syllabic stem of a noun: leadingCount
// stem syllables followed by the ending's syllables, written into buf (reused
// across length attempts, so it does not allocate). It mirrors the two-pass order
// of [sampleSyllabicStem] — onsets and nuclei first, then codas conditioned on
// the following onset — and extends it across the leading/ending junction: the
// last leading coda is drawn against the ending's first onset, so the join is
// phonotactically legal by construction. The returned stem's tonic index is fixed
// by the ending's stress (oxytone or paroxytone), so the whole word carries the
// correct stress before accentuation.
func sampleNounStem(r *rand.Rand, buf []syllable, leadingCount int, end *nounEnding) syllabicStem {
	if leadingCount < 1 {
		leadingCount = 1
	}
	total := leadingCount + len(end.syllables)
	if cap(buf) < total {
		buf = make([]syllable, total)
	} else {
		buf = buf[:total]
	}

	// Pass 1: leading onsets and nuclei. Each depends only on the syllable's
	// position, so the nucleus is drawn from the sub-inventory its onset licenses,
	// exactly as in sampleSyllabicStem.
	for i := 0; i < leadingCount; i++ {
		onsetInv := &onsetMedialInv
		if i == 0 {
			onsetInv = &onsetInitialInv
		}
		onset := sampleForm(r, onsetInv)

		nucInv := &nuclei
		if _, front := onsetFrontVowelOnly[onset]; front {
			nucInv = &nucleiFrontInv
		} else if _, back := onsetBackVowelOnly[onset]; back {
			nucInv = &nucleiBackInv
		}

		buf[i].onset = onset
		buf[i].nucleus = sampleForm(r, nucInv)
		buf[i].coda = "" // fixed in pass 2
	}

	// Place the ending's fixed syllables after the leading stem, then, when the
	// ending's first onset is sampled, draw it from the inventory its first
	// nucleus licenses (medial position, so ç, lh and nh are allowed).
	copy(buf[leadingCount:], end.syllables)
	firstOnset := end.syllables[0].onset
	if end.sampleOnset {
		onsetInv := &nounOnsetBeforeBackInv
		if end.frontOnset {
			onsetInv = &nounOnsetBeforeFrontInv
		}
		if end.onsetInvOverride != nil {
			onsetInv = end.onsetInvOverride
		}
		firstOnset = sampleForm(r, onsetInv)
		buf[leadingCount].onset = firstOnset
	}

	// Pass 2: leading codas, each conditioned on the following onset (the next
	// leading syllable, or the ending's first onset for the last leading
	// syllable), mirroring transitionIsValid. A vowel-initial follower forces an
	// empty coda (Maximum Onset Principle).
	for i := 0; i < leadingCount; i++ {
		nextOnset := firstOnset
		if i < leadingCount-1 {
			nextOnset = buf[i+1].onset
		}
		switch {
		case nextOnset == "":
			buf[i].coda = ""
		case nextOnset[0] == 'p' || nextOnset[0] == 'b':
			buf[i].coda = sampleForm(r, &codaBeforePBInv)
		default:
			buf[i].coda = sampleForm(r, &codaBeforeOtherInv)
		}
	}

	return syllabicStem{syllables: buf, tonic: total - 1 - end.stressFromEnd}
}

// nounOnsetBeforeFrontInv samples a consonant onset legal before a front vowel (e
// or i): every non-empty onset except ç, which requires a back or central vowel.
// It is used for the drawn onset of front-nucleus endings (-eiro, -eira).
var nounOnsetBeforeFrontInv = deriveInventory(0, onsetBeforeFront, onsetSingles, onsetClusters)

// nounOnsetBeforeBackInv samples a consonant onset legal before a back or central
// vowel (a, o or u): every non-empty onset except qu and gu, which require a front
// vowel. It is used for the drawn onset of back-nucleus endings (-o, -a, -or,
// -ora).
var nounOnsetBeforeBackInv = deriveInventory(0, onsetBeforeBack, onsetSingles, onsetClusters)

// onsetBeforeFront reports whether an onset may precede a front vowel: any onset
// except the cedilla ç (which is licensed only before a back or central vowel).
func onsetBeforeFront(form string) bool { return form != "ç" }

// onsetBeforeBack reports whether an onset may precede a back or central vowel:
// any onset except qu and gu (which are licensed only before a front vowel).
func onsetBeforeBack(form string) bool {
	_, front := onsetFrontVowelOnly[form]
	return !front
}

// stemRuneLen returns the character (rune) length of the word a stem assembles
// to, without assembling it. Graphic accentuation replaces a nucleus vowel with
// an accented vowel of the same rune count, so the accented word has the same
// length as this sum; the length window can therefore be checked before the
// single assembling allocation.
func stemRuneLen(syllables []syllable) int {
	n := 0
	for i := range syllables {
		n += utf8.RuneCountInString(syllables[i].onset) +
			utf8.RuneCountInString(syllables[i].nucleus) +
			utf8.RuneCountInString(syllables[i].coda)
	}
	return n
}

// pluralizeNoun returns the pt-PT plural of a singular noun. It reuses the shared
// string-level pluralizer for every ending except -ção, whose plural is the fixed
// productive outcome -ções (applied through the same [applyAoPlural] transform),
// rather than the weighted -ão choice the generic pluralizer makes for arbitrary
// words. The oxytone flag the generic pluralizer needs comes from the ending's
// stress.
func pluralizeNoun(r *rand.Rand, singular string, end *nounEnding) string {
	if end.fixedCaoPlural {
		return applyAoPlural(singular, aoOes)
	}
	return pluralize(r, singular, end.oxytone())
}

// ---------------------------------------------------------------------------
// Shared generation core
// ---------------------------------------------------------------------------

// nounCore is the single generation engine shared by the package-level noun
// functions and the *Generator noun methods. It draws only from r, so a seeded
// Generator is fully reproducible. buf is a caller-provided syllable buffer
// (backed by a stack array), reused across length attempts so that only the final
// assembly allocates.
//
// It resolves gender and number (Any -> a random valid value), selects an ending
// that realizes the gender, normalizes the length category up to the ending's
// minimum viable length, then samples a combined stem whose assembled length
// falls in the category window. A singular noun is the single-allocation
// assembled word; a plural noun applies the shared pluralizer to that word. The
// output is always lowercase, valid UTF-8, and orthographically well-formed.
func nounCore(r *rand.Rand, buf []syllable, g Gender, n Number, l LengthTypeWords) string {
	rg := resolveGender(r, g)
	rn := resolveNumber(r, n)

	_, ceiling := charRangeOf(l) // the requested category's character ceiling
	end := selectNounEnding(r, rg, ceiling)

	// The plural inflection lengthens the word, so its added characters count as
	// part of the fixed affix cost for both length normalization and stem sizing.
	// This raises the minimum viable length of this class-and-slot (a Small plural,
	// which cannot fit four characters comfortably, normalizes upward to a Medium
	// word) and makes the stem sampler target a stem short enough that the whole
	// PLURAL word lands in the window, not the singular. It is the specification's
	// per-slot minimum-viable-length rule, with the plural counted as inflection.
	affixCost := end.affixRunes
	if rn == Plural {
		affixCost += end.pluralRuneDelta
	}
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

		runes := stemRuneLen(stem.syllables)
		if rn == Plural {
			runes += end.pluralRuneDelta
		}
		if runes >= lo && runes <= hi {
			break
		}
	}

	if rn == Plural {
		// Additive plurals (a regular +s, or +es for -r/-z/-n/oxytone -s) append a
		// bare suffix onto the singular, so they are written into the SAME builder
		// that assembles the word: one allocation. Substitutive plurals (-ão, -m and
		// the vowel+l endings, plus the fixed -ção) rewrite the word's tail, so they
		// keep the post-assembly string transform: two allocations. See
		// [additivePluralSuffix].
		if suffix, additive := additivePluralSuffix(end); additive {
			return accentedWordSuffixed(stem, suffix)
		}
		return pluralizeNoun(r, accentedWord(stem), end)
	}
	return accentedWord(stem)
}

// ---------------------------------------------------------------------------
// Public API: package-level functions
// ---------------------------------------------------------------------------

// NounPT returns a random European-Portuguese (pt-PT) pseudo-noun with a random
// gender, number and length. It is exact sugar for NounPTOf(AnyGender, AnyNumber,
// AnyLengthWord). The result is a lowercase, valid-UTF-8 word that follows pt-PT
// noun morphology and graphic accentuation; it is a plausible invented word, not
// guaranteed to be a real dictionary entry.
//
// NounPT draws from the global, automatically seeded math/rand/v2 source and is
// safe for concurrent use by multiple goroutines. For reproducible output, use a
// seeded Generator ([New] or [NewSource]) and its NounPT method.
func NounPT() string {
	return NounPTOf(AnyGender, AnyNumber, AnyLengthWord)
}

// NounPTOf returns a random pt-PT pseudo-noun inflected for gender g and number n,
// within length category l. Each option's Any zero value ([AnyGender],
// [AnyNumber], [AnyLengthWord]) selects a random valid value. A gender-inflecting
// ending realizes the requested gender directly; a gender-fixed or common-gender
// ending is chosen only when it is compatible with the request. When l cannot
// hold any noun of the chosen ending, l is normalized upward to the smallest
// category that can (length is never normalized downward). The result is always
// lowercase and valid UTF-8; the function never returns an error and never
// panics.
//
// NounPTOf draws from the global, automatically seeded math/rand/v2 source and is
// safe for concurrent use by multiple goroutines.
func NounPTOf(g Gender, n Number, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return nounCore(wordsPTGlobalRand, buf[:0], g, n, l)
}

// ---------------------------------------------------------------------------
// Public API: *Generator methods
// ---------------------------------------------------------------------------

// NounPT is the seeded-generator equivalent of [NounPT]. A Generator created with
// the same seed produces the same sequence of nouns.
func (g *Generator) NounPT() string {
	return g.NounPTOf(AnyGender, AnyNumber, AnyLengthWord)
}

// NounPTOf is the seeded-generator equivalent of [NounPTOf]. A Generator created
// with the same seed produces the same sequence of nouns for the same arguments.
func (g *Generator) NounPTOf(gender Gender, n Number, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return nounCore(g.r, buf[:0], gender, n, l)
}
