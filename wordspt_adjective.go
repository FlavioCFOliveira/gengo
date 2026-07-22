package gengo

import (
	"math/rand/v2"
	"strings"
)

// This file implements the WordsPT pt-PT pseudo-adjective generators: the
// package-level AdjectivePT and AdjectivePTOf, plus their mirrored *Generator
// methods. It reuses, without duplication, the noun composition pipeline built in
// wordspt_noun.go (the combined leading-stem-plus-ending sampler sampleNounStem,
// the length helpers, the shared global source), and adds only the two things an
// adjective needs that a noun does not:
//
//   1. an adjective ending inventory (wordspt_noun.go's nounEnding value type is
//      general enough to describe an adjective ending, so the endings below are
//      nounEnding values; no noun ending or noun behavior is changed); and
//   2. the synthetic absolute superlative -íssimo/-íssima, a string-level suffix
//      applied to the bare stem, with the pt-PT orthographic adjustments.
//
// Positive-degree adjectives are built exactly like nouns: (leading stem
// syllables) + (the adjective ending's syllables) assembled into one combined
// syllabic stem whose tonic is fixed by the ending's stress, then graphic
// accentuation (Layer 3) and, for the plural, the shared string pluralizer.
//
// Adjective ending inventory and stress (specification/wordspt-open-classes.md):
//
//	-o / -a       gender-inflecting, paroxytone, no graphic accent (belo/bela)
//	-oso / -osa   gender-inflecting, paroxytone, no graphic accent (famoso)
//	-ico / -ica   gender-inflecting, PROPAROXYTONE, accented antepenult (básico)
//	-ivo / -iva   gender-inflecting, paroxytone, no graphic accent (ativo)
//	-al           gender-invariable, oxytone, no graphic accent (legal)
//	-ável         gender-invariable, paroxytone, accented á (amável)
//	-ível         gender-invariable, paroxytone, accented í (possível)
//	-ente         gender-invariable, paroxytone, no graphic accent (presente)
//	-ante         gender-invariable, paroxytone, no graphic accent (elegante)
//
// -ico/-ica is the first proparoxytone ending in WordsPT (stressFromEnd 2): its
// tonic is the last leading-stem syllable, and the general accentuation layer
// (wordspt_accent.go) already places a proparoxytone accent correctly, so no new
// accentuation code is needed. Gender behavior follows the shared rule (Cunha &
// Cintra): the -o/-a, -ico/-ica and -ivo/-iva pairs inflect for gender, while
// -al, -ável, -ível, -ente and -ante are invariable and serve either gender with
// one form.
//
// The -íssimo/-íssima absolute superlative
//
// The synthetic absolute superlative is a string-level suffix, exactly like the
// plural: it is not re-syllabified, so it is validated against the Sprint 7
// STRING oracles (cedilla, diaeresis, grave, valid lowercase UTF-8), not the
// decomposition oracles. It moves the primary stress onto its own acute í, so the
// whole word becomes proparoxytone (X-í-ssi-mo) and the base loses whatever
// graphic accent it carried. This falls out for free here: the combined stem
// stores BARE syllables (Layer 3 adds the acute/circumflex only during assembly),
// so a superlative built directly from the bare syllables carries no base accent
// and exactly one acute — the í of -íssimo. Nasal tildes on the bare syllables are
// part of the grapheme (nasality, not stress) and are preserved, as in any
// unstressed nasal.
//
// Formation, from the bare combined stem (Acordo Ortográfico; Cunha & Cintra):
//
//   - Base whose last syllable is vowel-final (its coda is empty: -o, -oso, -ico,
//     -ivo, -ente, -ante): drop that final nucleus vowel, then keep the exposed
//     onset consonant, with three orthographic adjustments before the front í:
//     c -> qu (rico -> riquíssimo) and g -> gu (longo -> longuíssimo) preserve the
//     hard /k/,/g/; and ç -> c (a ç onset exposed by the thematic -o/-a, as in a
//     -ço base) keeps the soft /s/ while obeying the cedilla rule, since ç is never
//     written before e or i (the standard pt-PT alternation caçar -> cacei). Append
//     -íssim + the gender vowel (+ the plural s).
//   - Base whose last syllable is consonant-final (its coda is non-empty: -al,
//     -ável, -ível): keep the whole base and append -íssim + gender (+ plural)
//     directly (banal -> banalíssimo).
//
// The superlative inflects for gender and number through its own suffix vowel and
// plural s: -íssimo, -íssima, -íssimos, -íssimas. Because the final base vowel is
// dropped, the superlative is independent of the base's gender marker, so it is
// formed identically whichever gender realization the ending selected.
//
// Flagged simplifications (pseudo-word scope; recorded per the project rules):
//
//   - Some -vel/-il/-re adjectives take erudite superlatives (amável ->
//     amabilíssimo, fácil -> facílimo, célebre -> celebérrimo). Those -ílimo and
//     -érrimo variants are lexical, not productive; for an invented word the
//     REGULAR -íssimo is applied (amável -> amavelíssimo), as the task specifies.
//   - The thematic -o/-a ending draws the onset before its (back) thematic vowel,
//     so a c or g onset occurs (rico, longo). Its superlative drops the vowel and
//     exposes that onset before the front í, so BOTH orthographic hardenings are
//     reachable through generation: c -> qu (rico -> riquíssimo) and g -> gu (longo
//     -> longuíssimo). The proparoxytone -ico/-ica ending, whose final syllable has
//     a fixed c onset, additionally exercises c -> qu (básico -> basiquíssimo).
//
// Allocation. The positive singular assembles in exactly one allocation (Layer
// 3's builder). An ADDITIVE positive plural (the +s of the vowel- and
// diphthong-final endings -o, -oso, -ico, -ivo, -ente, -ante) is also a single
// allocation: the s is written into the same assembling builder (see
// accentedWordSuffixed and additivePluralSuffix). A SUBSTITUTIVE positive plural
// (the vowel+l endings -al, -ável, -ível -> -ais/-áveis/-íveis) adds one tail
// concatenation in the shared pluralizer. The superlative is also a single
// allocation: it is written straight from the bare syllables into one pre-sized
// builder, with no separate base assembly. The length-window search reuses the
// stack syllable buffer and measures the superlative length from a closed-form
// rune count, so it allocates nothing per attempt.
//
// Sources: Cunha, C. & Cintra, L., "Nova Gramática do Português Contemporâneo"
// (adjective formation, gender of suffixes, the -íssimo absolute superlative and
// its orthographic adjustments); Acordo Ortográfico da Língua Portuguesa (graphic
// accent, cedilla and the qu/gu spelling before front vowels). The sampling
// weights are ordinal estimates (flagged), consistent with the other WordsPT
// weights; they change only how often each valid ending appears.

// ---------------------------------------------------------------------------
// Adjective ending inventory
// ---------------------------------------------------------------------------

// The realized adjective endings. Each is a concrete gender realization (as for
// nouns): the masculine and feminine members of a gender-inflecting pair are two
// entries, and the invariable endings are single entries shared by both genders.
// They are nounEnding values because that type already describes exactly what an
// ending contributes — its trailing syllables, whether the first onset is drawn,
// its stress distance from the word end, and its rune/plural costs — with no
// noun-specific behavior beyond the unused fixedCaoPlural flag (false here).
var (
	// endAdjMascO is the masculine thematic -o (belo, alto, novo, rico, longo): the
	// most common adjective class in pt-PT. It is paroxytone, a single final
	// syllable with a sampled onset before the (back) thematic vowel o, plural +s.
	// It mirrors the noun thematic -o (endMascO) exactly, reusing the shared ending
	// machinery with no noun-specific behavior. Because its onset is sampled, a c or
	// g may precede the o (rico, longo), so its -íssimo superlative makes both the
	// c -> qu (riquíssimo) and the g -> gu (longuíssimo) hardening reachable through
	// generation.
	endAdjMascO = nounEnding{
		label: "-o", syllables: []syllable{{nucleus: "o"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 2, pluralRuneDelta: 1,
	}
	// endAdjFemA is the feminine counterpart -a (bela, alta, nova, rica, longa),
	// mirroring the noun thematic -a (endFemA).
	endAdjFemA = nounEnding{
		label: "-a", syllables: []syllable{{nucleus: "a"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 2, pluralRuneDelta: 1,
	}
	// endAdjMascOso is the masculine -oso (famoso): paroxytone, syllables o+so
	// with a sampled onset before the first (back) vowel, plural +s.
	endAdjMascOso = nounEnding{
		label: "-oso", syllables: []syllable{{nucleus: "o"}, {onset: "s", nucleus: "o"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjFemOsa is the feminine counterpart -osa (famosa).
	endAdjFemOsa = nounEnding{
		label: "-osa", syllables: []syllable{{nucleus: "o"}, {onset: "s", nucleus: "a"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjMascIco is the masculine -ico (básico): PROPAROXYTONE (stressFromEnd
	// 2), so the tonic is the last leading-stem syllable; syllables i+co with a
	// sampled onset before the first (front) vowel, plural +s.
	endAdjMascIco = nounEnding{
		label: "-ico", syllables: []syllable{{nucleus: "i"}, {onset: "c", nucleus: "o"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 2, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjFemIca is the feminine counterpart -ica (básica).
	endAdjFemIca = nounEnding{
		label: "-ica", syllables: []syllable{{nucleus: "i"}, {onset: "c", nucleus: "a"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 2, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjMascIvo is the masculine -ivo (ativo): paroxytone, syllables i+vo with
	// a sampled onset before the first (front) vowel, plural +s.
	endAdjMascIvo = nounEnding{
		label: "-ivo", syllables: []syllable{{nucleus: "i"}, {onset: "v", nucleus: "o"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjFemIva is the feminine counterpart -iva (ativa).
	endAdjFemIva = nounEnding{
		label: "-iva", syllables: []syllable{{nucleus: "i"}, {onset: "v", nucleus: "a"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 4, pluralRuneDelta: 1,
	}
	// endAdjAl is the gender-invariable -al (legal): oxytone, a single final
	// syllable (sampled onset + a + coda l), plural -ais (the vowel+l rule).
	endAdjAl = nounEnding{
		label: "-al", syllables: []syllable{{nucleus: "a", coda: "l"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 0, affixRunes: 3, pluralRuneDelta: 1,
	}
	// endAdjAvel is the gender-invariable -ável (amável): paroxytone whose tonic á
	// is graphically accented (a paroxytone ending in -l is non-default), syllables
	// a+vel with a sampled onset before the first (back) vowel, plural -áveis.
	endAdjAvel = nounEnding{
		label: "-ável", syllables: []syllable{{nucleus: "a"}, {onset: "v", nucleus: "e", coda: "l"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endAdjIvel is the gender-invariable -ível (possível): paroxytone whose tonic
	// í is graphically accented, syllables i+vel with a sampled onset before the
	// first (front) vowel, plural -íveis.
	endAdjIvel = nounEnding{
		label: "-ível", syllables: []syllable{{nucleus: "i"}, {onset: "v", nucleus: "e", coda: "l"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endAdjEnte is the gender-invariable -ente (presente): paroxytone with a
	// default (unaccented) ending, syllables en+te with a sampled onset before the
	// first (front) vowel; the coda n before the t is legal by construction, plural
	// +s.
	endAdjEnte = nounEnding{
		label: "-ente", syllables: []syllable{{nucleus: "e", coda: "n"}, {onset: "t", nucleus: "e"}},
		sampleOnset: true, frontOnset: true,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
	// endAdjAnte is the gender-invariable -ante (elegante): paroxytone with a
	// default (unaccented) ending, syllables an+te with a sampled onset before the
	// first (back) vowel, plural +s.
	endAdjAnte = nounEnding{
		label: "-ante", syllables: []syllable{{nucleus: "a", coda: "n"}, {onset: "t", nucleus: "e"}},
		sampleOnset: true, frontOnset: false,
		stressFromEnd: 1, affixRunes: 5, pluralRuneDelta: 1,
	}
)

// masculineAdjectiveEndings lists the endings a masculine request may realize: the
// masculine gender-inflecting forms and every gender-invariable ending. Weights
// are ordinal estimates (flagged): the thematic -o is the single most common
// adjective class, so it leads, while the productive derivational -oso, -al, -ico
// and -ente remain well represented (the thematic weight leads but does not
// overwhelm them).
var masculineAdjectiveEndings = []weightedNounEnding{
	{&endAdjMascO, 30},
	{&endAdjMascOso, 22},
	{&endAdjMascIco, 14},
	{&endAdjMascIvo, 12},
	{&endAdjAl, 18},
	{&endAdjAvel, 8},
	{&endAdjIvel, 6},
	{&endAdjEnte, 12},
	{&endAdjAnte, 8},
}

// feminineAdjectiveEndings lists the endings a feminine request may realize: the
// feminine gender-inflecting forms and every gender-invariable ending. Weights
// mirror masculineAdjectiveEndings (ordinal estimates, flagged).
var feminineAdjectiveEndings = []weightedNounEnding{
	{&endAdjFemA, 30},
	{&endAdjFemOsa, 22},
	{&endAdjFemIca, 14},
	{&endAdjFemIva, 12},
	{&endAdjAl, 18},
	{&endAdjAvel, 8},
	{&endAdjIvel, 6},
	{&endAdjEnte, 12},
	{&endAdjAnte, 8},
}

// adjectiveEndingsFor returns the ending list a concrete gender may realize. g is
// expected to be [Masculine] or [Feminine] (as returned by [resolveGender]); any
// other value is treated as masculine, so the function never panics.
func adjectiveEndingsFor(g Gender) []weightedNounEnding {
	if g == Feminine {
		return feminineAdjectiveEndings
	}
	return masculineAdjectiveEndings
}

// selectFittingEnding draws an ending from list that fits the length ceiling
// maxChars (its shortest positive word is at most maxChars characters), using the
// ordinal weights in a single bounded draw and a short linear scan, allocating
// nothing. When no ending fits (the requested category is shorter than the
// shortest ending of that gender), it returns the shortest ending; the caller then
// normalizes the length upward for it. It generalizes the length-aware selection
// pattern of [selectNounEnding] to any weighted ending list, so the adjective path
// selects endings the same way the noun path does, without altering the noun code.
func selectFittingEnding(r *rand.Rand, list []weightedNounEnding, maxChars int) *nounEnding {
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
// -íssimo/-íssima superlative machinery
// ---------------------------------------------------------------------------

// adjSuperlativeStem is the stressed core of the absolute superlative suffix,
// without its gender vowel. Its í carries the single acute of the whole word, so
// the superlative is proparoxytone (X-í-ssi-mo). The gender vowel (o or a) and, in
// the plural, a final s are appended after it.
const adjSuperlativeStem = "íssim"

// adjSuperlativeSuffixRunes is the rune (character) length that the superlative
// suffix adds for its stressed core plus its single gender vowel: len([íssim]) is
// five runes and the gender vowel is one, so -íssimo and -íssima are six runes
// each. The plural adds one more rune (s), counted separately.
const adjSuperlativeSuffixRunes = 6

// endingSuperlativeDropsVowel reports whether an ending's superlative drops the
// base's final nucleus vowel: it does when the ending's last syllable is
// vowel-final (its coda is empty, as in -oso, -ico, -ivo, -ente, -ante), and does
// not when the last syllable is consonant-final (-al, -ável, -ível), where the
// suffix attaches directly.
func endingSuperlativeDropsVowel(end *nounEnding) bool {
	return end.syllables[len(end.syllables)-1].coda == ""
}

// endingSuperlativeHardens reports whether dropping the final vowel exposes a c or
// g onset that must become qu or gu to preserve the hard /k/,/g/ before the front
// í of -íssimo (rico -> riquíssimo, longo -> longuíssimo). It inspects the ending's
// last syllable onset and is meaningful only when [endingSuperlativeDropsVowel] is
// true. It exists solely to account for the one extra RUNE that qu/gu adds when
// sizing the static minimum-viable length (the ç -> c hardening is length-neutral,
// one rune to one rune, so it is intentionally excluded). The exact per-word length
// after sampling, which does account for a sampled c/g onset, is computed by
// [superlativeStemRuneLen]. It reports false for a sampled-onset ending such as the
// thematic -o/-a, whose stored onset is empty because it is drawn at build time.
func endingSuperlativeHardens(end *nounEnding) bool {
	onset := end.syllables[len(end.syllables)-1].onset
	return onset == "c" || onset == "g"
}

// superlativeAffixRunes returns the rune cost the masculine-singular superlative
// adds on top of the bare leading stem for this ending: the ending's own rune cost,
// minus the dropped final vowel (and plus one when that exposes a c/g that becomes
// qu/gu), plus the -íssimo suffix. It sizes the leading stem and the minimum viable
// length for a superlative before any syllable is sampled.
func superlativeAffixRunes(end *nounEnding) int {
	a := end.affixRunes
	if endingSuperlativeDropsVowel(end) {
		a--
		if endingSuperlativeHardens(end) {
			a++
		}
	}
	return a + adjSuperlativeSuffixRunes
}

// superlativeStemRuneLen returns the exact character (rune) length of the
// superlative that a sampled combined stem assembles to for the given number,
// without assembling it. It mirrors [superlativeWord] on the actual stem: every
// syllable before the last contributes in full, the last syllable contributes its
// whole self when it is consonant-final or only its (possibly hardened) onset when
// it is vowel-final, and the -íssim suffix, the gender vowel and the plural s are
// added. Because it inspects the sampled last-syllable onset directly, it accounts
// for the c -> qu / g -> gu hardening even for the thematic -o/-a ending, whose
// onset is drawn at build time and is therefore invisible to a static estimate.
// The length window can thus be enforced from a closed-form count with no
// allocation.
func superlativeStemRuneLen(stem syllabicStem, plural bool) int {
	syllables := stem.syllables
	last := len(syllables) - 1
	n := stemRuneLen(syllables[:last])
	if syllables[last].coda == "" {
		// Vowel-final: only the exposed onset survives, after hardening. Every
		// hardened onset (qu, gu, c, or an unchanged onset) is ASCII, so its byte
		// length equals its rune length.
		n += hardenedOnsetLen(syllables[last].onset)
	} else {
		n += stemRuneLen(syllables[last : last+1])
	}
	n += adjSuperlativeSuffixRunes
	if plural {
		n++
	}
	return n
}

// superlativeWord renders the absolute superlative of a bare combined stem in a
// single allocation. It writes every base syllable straight from the bare
// (unaccented) stem — so the word carries no base graphic accent, only the acute í
// of the suffix — and treats the last syllable per the formation rule: a
// vowel-final last syllable contributes its onset only (with c -> qu, g -> gu,
// ç -> c), while a consonant-final one contributes in full. It then appends -íssim,
// the gender vowel (o for masculine, a for feminine) and, in the plural, a final s.
// g and n are the resolved gender and number.
func superlativeWord(stem syllabicStem, g Gender, n Number) string {
	syllables := stem.syllables
	last := len(syllables) - 1
	dropsVowel := syllables[last].coda == ""
	plural := n == Plural

	genderVowel := "o"
	if g == Feminine {
		genderVowel = "a"
	}

	// Pre-size the builder to the exact byte length for a single allocation.
	total := 0
	for i := 0; i < last; i++ {
		total += len(syllables[i].onset) + len(syllables[i].nucleus) + len(syllables[i].coda)
	}
	if dropsVowel {
		total += hardenedOnsetLen(syllables[last].onset)
	} else {
		total += len(syllables[last].onset) + len(syllables[last].nucleus) + len(syllables[last].coda)
	}
	total += len(adjSuperlativeStem) + len(genderVowel)
	if plural {
		total++
	}

	var b strings.Builder
	b.Grow(total)
	for i := 0; i < last; i++ {
		b.WriteString(syllables[i].onset)
		b.WriteString(syllables[i].nucleus)
		b.WriteString(syllables[i].coda)
	}
	if dropsVowel {
		writeHardenedOnset(&b, syllables[last].onset)
	} else {
		b.WriteString(syllables[last].onset)
		b.WriteString(syllables[last].nucleus)
		b.WriteString(syllables[last].coda)
	}
	b.WriteString(adjSuperlativeStem)
	b.WriteString(genderVowel)
	if plural {
		b.WriteByte('s')
	}
	return b.String()
}

// hardenedOnsetLen returns the byte length the onset occupies before the front í
// of -íssimo, applying the c -> qu, g -> gu and ç -> c hardening. A c or g becomes
// the two-byte qu or gu; the two-byte ç becomes the one-byte c (the cedilla rule
// forbids ç before a front vowel, and c already spells /s/ there); any other onset,
// being ASCII, keeps its own length. Every result is ASCII, so this byte length
// also equals the onset's rune length after hardening.
func hardenedOnsetLen(onset string) int {
	switch onset {
	case "c", "g":
		return 2
	case "ç":
		return 1
	default:
		return len(onset)
	}
}

// writeHardenedOnset writes the exposed final onset before -íssimo: c becomes qu
// and g becomes gu (keeping the hard /k/,/g/ sound), ç becomes c (keeping the soft
// /s/ sound while obeying the cedilla rule, which never writes ç before e or i),
// and any other onset is written unchanged.
func writeHardenedOnset(b *strings.Builder, onset string) {
	switch onset {
	case "c":
		b.WriteString("qu")
	case "g":
		b.WriteString("gu")
	case "ç":
		b.WriteString("c")
	default:
		b.WriteString(onset)
	}
}

// ---------------------------------------------------------------------------
// Shared generation core
// ---------------------------------------------------------------------------

// adjectiveCore is the single generation engine shared by the package-level
// adjective functions and the *Generator adjective methods. It draws only from r,
// so a seeded Generator is fully reproducible, and reuses a caller-provided
// syllable buffer (backed by a stack array) so that only the final assembly
// allocates.
//
// It resolves gender, number and degree (Any -> a random valid value), selects an
// ending that realizes the gender, sizes the leading stem for the requested degree
// and number (normalizing the length category up to the minimum viable length),
// then samples a combined stem whose assembled word falls in the category window.
// A positive adjective is the accented word (plus the shared pluralizer for a
// plural); a superlative is the single-allocation -íssimo form. The output is
// always lowercase, valid UTF-8, and orthographically well-formed.
func adjectiveCore(r *rand.Rand, buf []syllable, g Gender, n Number, d Degree, l LengthTypeWords) string {
	rg := resolveGender(r, g)
	rn := resolveNumber(r, n)
	rd := resolveDegree(r, d)

	_, ceiling := charRangeOf(l) // the requested category's character ceiling
	end := selectFittingEnding(r, adjectiveEndingsFor(rg), ceiling)

	// The degree and number both lengthen the word, so their added characters count
	// as fixed affix cost for length normalization and stem sizing: a superlative
	// uses the -íssimo affix, and a plural adds its own delta (the fixed +1 s of
	// -íssimos for a superlative, the ending's plural delta for a positive).
	affixCost := end.affixRunes
	if rd == Superlative {
		affixCost = superlativeAffixRunes(end)
	}
	if rn == Plural {
		if rd == Superlative {
			affixCost++ // -íssimo -> -íssimos
		} else {
			affixCost += end.pluralRuneDelta
		}
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

		var runes int
		if rd == Superlative {
			runes = superlativeStemRuneLen(stem, rn == Plural)
		} else {
			runes = stemRuneLen(stem.syllables)
			if rn == Plural {
				runes += end.pluralRuneDelta
			}
		}
		if runes >= lo && runes <= hi {
			break
		}
	}

	if rd == Superlative {
		return superlativeWord(stem, rg, rn)
	}
	if rn == Plural {
		// The vowel- and diphthong-final adjective endings (-o, -oso, -ico, -ivo,
		// -ente, -ante) take a regular +s, appended into the SAME builder that
		// assembles the word: one allocation. The substitutive vowel+l endings (-al,
		// -ável, -ível -> -ais/-áveis/-íveis) rewrite the word's tail, so they keep the
		// post-assembly string transform: two allocations. See [additivePluralSuffix].
		if suffix, additive := additivePluralSuffix(end); additive {
			return accentedWordSuffixed(stem, suffix)
		}
		return pluralize(r, accentedWord(stem), end.oxytone())
	}
	return accentedWord(stem)
}

// ---------------------------------------------------------------------------
// Public API: package-level functions
// ---------------------------------------------------------------------------

// AdjectivePT returns a random European-Portuguese (pt-PT) pseudo-adjective with a
// random gender, number, degree and length. It is exact sugar for
// AdjectivePTOf(AnyGender, AnyNumber, AnyDegree, AnyLengthWord). The result is a
// lowercase, valid-UTF-8 word that follows pt-PT adjective morphology and graphic
// accentuation; it is a plausible invented word, not guaranteed to be a real
// dictionary entry.
//
// AdjectivePT draws from the global, automatically seeded math/rand/v2 source and
// is safe for concurrent use by multiple goroutines. For reproducible output, use
// a seeded Generator ([New] or [NewSource]) and its AdjectivePT method.
func AdjectivePT() string {
	return AdjectivePTOf(AnyGender, AnyNumber, AnyDegree, AnyLengthWord)
}

// AdjectivePTOf returns a random pt-PT pseudo-adjective inflected for gender g,
// number n and degree d, within length category l. Each option's Any zero value
// ([AnyGender], [AnyNumber], [AnyDegree], [AnyLengthWord]) selects a random valid
// value. A gender-inflecting ending realizes the requested gender directly; a
// gender-invariable ending serves either request with one form. Degree [Positive]
// yields the base adjective; [Superlative] yields the synthetic absolute
// superlative -íssimo/-íssima, inflected for gender and number. When l cannot hold
// any adjective of the chosen ending and degree, l is normalized upward to the
// smallest category that can (length is never normalized downward); in particular
// a Superlative never fits SmallLengthWord and normalizes up. The result is always
// lowercase and valid UTF-8; the function never returns an error and never panics.
//
// AdjectivePTOf draws from the global, automatically seeded math/rand/v2 source and
// is safe for concurrent use by multiple goroutines.
func AdjectivePTOf(g Gender, n Number, d Degree, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return adjectiveCore(wordsPTGlobalRand, buf[:0], g, n, d, l)
}

// ---------------------------------------------------------------------------
// Public API: *Generator methods
// ---------------------------------------------------------------------------

// AdjectivePT is the seeded-generator equivalent of [AdjectivePT]. A Generator
// created with the same seed produces the same sequence of adjectives.
func (g *Generator) AdjectivePT() string {
	return g.AdjectivePTOf(AnyGender, AnyNumber, AnyDegree, AnyLengthWord)
}

// AdjectivePTOf is the seeded-generator equivalent of [AdjectivePTOf]. A Generator
// created with the same seed produces the same sequence of adjectives for the same
// arguments.
func (g *Generator) AdjectivePTOf(gender Gender, n Number, d Degree, l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return adjectiveCore(g.r, buf[:0], gender, n, d, l)
}
