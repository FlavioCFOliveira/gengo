package gengo

import (
	"strings"
	"unicode/utf8"
)

// This file implements WordsPT Layer 3: deterministic European-Portuguese
// (pt-PT) graphic accentuation over a generated syllabic stem, plus the
// single-allocation accented assembler. It builds on the phonotactic sampler in
// wordspt_stem.go (which fixes the syllable sequence and the tonic index) and
// satisfies the accentuation oracle in wordspt.go (stemConformsToAccentuation)
// exactly.
//
// Correct-by-construction accentuation
//
// Because Layer 1 controls the tonic syllable, the accent is derived from the
// pt-PT rules, never guessed and never found by trial-and-check:
//
//   - A nasal tonic nucleus (one bearing a tilde: ã, õ, ão, ãe, õe) is already
//     marked for stress by the tilde, so no acute or circumflex is added.
//   - An oral tonic nucleus receives an acute or circumflex accent if and only
//     if tonicRequiresGraphicAccent (the Acordo Ortográfico rule set encoded in
//     wordspt.go) demands one; otherwise it stays bare.
//   - The accent is placed on the FIRST vowel of the tonic nucleus. For a
//     monophthong that is the vowel itself; for a pt-PT falling diphthong
//     (vowel + glide, the only oral diphthongs in the inventory) it is the
//     nuclear vowel, which is the first element (papéis -> éi, herói -> ói).
//
// No other (non-tonic) syllable ever receives an acute or circumflex accent,
// which is exactly the placement invariant the oracle checks. The nasal tilde,
// which the oracle permits on any syllable, is produced by the sampler and left
// untouched here (for example the unstressed ão of órgão).
//
// Accent shape (acute vs circumflex)
//
// The oracle constrains only the PRESENCE of an acute-or-circumflex accent on
// the tonic oral nucleus, not its quality, so the acute/circumflex choice is a
// documented generation decision, grounded in pt-PT orthography:
//
//   - i and u take only the acute (í, ú); pt-PT has no î or û.
//   - e and o take the acute (é, ó). pt-PT favors the OPEN tonic vowel far more
//     than pt-BR (também, porém, herói, avó, café, género, António, fenómeno),
//     and, decisively, the -em/-ens oxytone ending is written with the acute
//     (também, porém, parabéns); a blanket "closed before a nasal" rule would
//     wrongly yield *tambêm. The acute is therefore always correct-by-oracle and
//     never produces a known-wrong pt-PT spelling.
//   - a takes the circumflex (â) when it is a lone tonic vowel immediately
//     closed by a nasal consonant (a coda m or n, or a following syllable whose
//     onset begins with a nasal m, n or the digraph nh), and the acute (á)
//     otherwise. â before a nasal is the one circumflex context that pt-PT
//     applies consistently and without a stressed-oral counterexample (câmara,
//     lâmpada, âmbar, âncora, âmago), whereas a stressed oral á never precedes a
//     nasal in pt-PT (the -am ending is unstressed: falam, cantam).
//
// Consequences of this rule, recorded for transparency (see the task report):
// ê and ô never appear, because the closed quality of a tonic e or o in pt-PT is
// lexically determined and cannot be inferred from spelling context for an
// invented word; the open acute is the safe, always-valid default. This choice
// affects only vowel quality, which the oracle does not validate; every output
// remains orthographically well-formed pt-PT.
//
// Determinism and allocation
//
// The accent is a pure function of the stem and its tonic index: it uses no
// randomness, so it needs no *rand.Rand source and is trivially reproducible.
// The accented word is produced in a SINGLE allocation: accentuateStem computes
// a tiny plan without allocating, and assembleAccentedStem renders the word in
// one pre-sized pass, substituting the tonic nucleus's first vowel in place, so
// the returned string is the only allocation on the path.
//
// Sources: Acordo Ortográfico da Língua Portuguesa (1990), Bases IX–XIII (acute,
// circumflex and tilde); Cunha, C. & Cintra, L., "Nova Gramática do Português
// Contemporâneo" (stress classes, vowel quality and the pt-PT open-vowel
// preference); Mateus, M. H. & d'Andrade, E., "The Phonology of Portuguese"
// (the falling-diphthong nucleus).

// accentPlan is the deterministic, allocation-free description of the single
// graphic accent a stem's tonic nucleus must carry. tonic is the index of the
// syllable whose nucleus receives the accent, or -1 when the stem takes no
// acute or circumflex accent at all (a nasal tonic, or an oral tonic that the
// rules leave bare). When tonic is non-negative, accented is the accented vowel
// rune that replaces the first vowel of that nucleus during assembly.
type accentPlan struct {
	tonic    int
	accented rune
}

// accentuateStem derives the pt-PT graphic accent for stem from its tonic index,
// returning an [accentPlan]. It allocates nothing and draws no randomness: the
// accent is a deterministic function of the stem, correct by construction from
// the rules in [tonicRequiresGraphicAccent]. It never mutates stem; the plan is
// applied later, during assembly, by [assembleAccentedStem].
func accentuateStem(stem syllabicStem) accentPlan {
	syllables := stem.syllables
	tonic := stem.tonic
	if len(syllables) == 0 || tonic < 0 || tonic >= len(syllables) {
		return accentPlan{tonic: -1}
	}

	nucleus := syllables[tonic].nucleus

	// A nasal tonic nucleus is marked for stress by its tilde; add nothing.
	if nucleusIsNasal(nucleus) {
		return accentPlan{tonic: -1}
	}

	// An oral tonic nucleus is accented only when the rules demand it.
	if !tonicRequiresGraphicAccent(syllables, tonic) {
		return accentPlan{tonic: -1}
	}

	// Place the accent on the first vowel of the nucleus (the nuclear vowel;
	// in a falling diphthong the second element is a glide).
	base, _ := vowelBase(firstRune(nucleus)) // always a vowel by inventory construction
	accented := accentedVowel(base, tonicTakesCircumflex(syllables, tonic))
	return accentPlan{tonic: tonic, accented: accented}
}

// tonicTakesCircumflex reports whether the tonic vowel closes to a circumflex
// under the documented pt-PT accent-shape rule: only a lone tonic a immediately
// followed by a nasal consonant (its own coda m or n, or a following syllable
// whose onset begins with a nasal m, n or nh) closes to â. Every other oral
// tonic vowel takes the acute. See the accent-shape note at the top of the file.
func tonicTakesCircumflex(syllables []syllable, tonic int) bool {
	s := syllables[tonic]
	// Only a lone a (a single-vowel oral nucleus) closes before a nasal; a
	// diphthong's a is followed by its glide, not by the nasal.
	if s.nucleus != "a" {
		return false
	}
	if s.coda == "m" || s.coda == "n" {
		return true
	}
	if s.coda == "" && tonic+1 < len(syllables) {
		next := syllables[tonic+1].onset
		if next != "" && (next[0] == 'm' || next[0] == 'n') {
			return true
		}
	}
	return false
}

// accentedVowel returns the accented pt-PT form of the base oral vowel. With
// circumflex true it returns the closed form (â, ê, ô) for a, e and o; i and u,
// which have no circumflex in pt-PT, always return their acute form. With
// circumflex false it returns the open acute form (á, é, í, ó, ú). base must be
// one of the five base vowels produced by [vowelBase].
func accentedVowel(base rune, circumflex bool) rune {
	if circumflex {
		switch base {
		case 'a':
			return 'â'
		case 'e':
			return 'ê'
		case 'o':
			return 'ô'
		}
		// i and u fall through to the acute: pt-PT has no î or û.
	}
	switch base {
	case 'a':
		return 'á'
	case 'e':
		return 'é'
	case 'i':
		return 'í'
	case 'o':
		return 'ó'
	case 'u':
		return 'ú'
	}
	return base // unreachable for a valid base vowel
}

// assembleAccentedStem renders syllables to a single lowercase, valid-UTF-8
// pt-PT word in exactly one allocation, applying plan as it writes. It pre-sizes
// a strings.Builder to the exact byte length of the output — the concatenated
// syllable parts, adjusted for the extra bytes the accented rune adds over the
// plain first vowel it replaces — so the writes never reallocate and
// Builder.String returns the bytes without a second copy. When plan.tonic is
// negative it produces the bare (unaccented) assembly, identical to
// [assembleStem]. Every part and every substituted rune is lowercase and valid
// UTF-8, so the result is lowercase, valid UTF-8.
func assembleAccentedStem(syllables []syllable, plan accentPlan) string {
	return assembleAccentedStemSuffixed(syllables, plan, "")
}

// assembleAccentedStemSuffixed is [assembleAccentedStem] with an optional
// trailing suffix written into the SAME pre-sized builder, so an additive pt-PT
// plural marker (a bare "s" or "es"; see [additivePluralSuffix]) costs no extra
// allocation over the singular assembly. The suffix is appended verbatim after
// the last syllable's coda, which is exactly where the string-level [pluralize]
// would concatenate it for an additive ending, so the result is byte-identical
// to pluralizing the assembled word. An empty suffix reproduces
// [assembleAccentedStem]'s output exactly (WriteString of "" writes nothing and
// adds no allocation).
func assembleAccentedStemSuffixed(syllables []syllable, plan accentPlan, suffix string) string {
	total := len(suffix)
	for i := range syllables {
		total += len(syllables[i].onset) + len(syllables[i].nucleus) + len(syllables[i].coda)
	}
	if plan.tonic >= 0 {
		old, _ := utf8.DecodeRuneInString(syllables[plan.tonic].nucleus)
		total += utf8.RuneLen(plan.accented) - utf8.RuneLen(old)
	}

	var b strings.Builder
	b.Grow(total)
	for i := range syllables {
		b.WriteString(syllables[i].onset)
		if i == plan.tonic {
			nucleus := syllables[i].nucleus
			_, size := utf8.DecodeRuneInString(nucleus)
			b.WriteRune(plan.accented)
			b.WriteString(nucleus[size:])
		} else {
			b.WriteString(syllables[i].nucleus)
		}
		b.WriteString(syllables[i].coda)
	}
	b.WriteString(suffix)
	return b.String()
}

// accentedWord is the composed Layer 3 entry point: it accentuates stem and
// renders it to the final accented, lowercase, valid-UTF-8 pt-PT word. Given a
// stem sampled into a reused buffer (zero-allocation), the returned string is
// the single allocation of the whole build path.
func accentedWord(stem syllabicStem) string {
	return assembleAccentedStem(stem.syllables, accentuateStem(stem))
}

// accentedWordSuffixed is [accentedWord] with an additive pt-PT plural suffix
// appended in the SAME single allocation. It is the one-allocation plural path
// for endings whose plural is a pure suffix append (the additive endings
// classified by [additivePluralSuffix]): the stem is accentuated and assembled
// exactly as [accentedWord] does, with suffix written into the same builder. The
// output equals accentedWord(stem)+suffix, which for an additive ending is
// exactly pluralize(accentedWord(stem)); passing an empty suffix reproduces
// [accentedWord].
func accentedWordSuffixed(stem syllabicStem, suffix string) string {
	return assembleAccentedStemSuffixed(stem.syllables, accentuateStem(stem), suffix)
}
