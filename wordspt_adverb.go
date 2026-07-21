package gengo

import (
	"math/rand/v2"
	"strings"
)

// This file implements the WordsPT pt-PT pseudo-adverb generators: the
// package-level AdverbPT and AdverbPTByLengthType, plus their mirrored
// *Generator methods. It reuses, without duplication, the adjective machinery
// built in wordspt_adjective.go (the feminine adjective ending inventory, the
// length-aware ending selection selectFittingEnding, and the combined
// leading-stem-plus-ending sampler sampleNounStem) and the shared length
// machinery in wordspt_inflection.go, and adds only what an adverb needs: the
// invariable -mente suffix and its single-allocation assembly from the BARE
// adjective stem.
//
// Formation (specification/wordspt-open-classes.md, "Adverb"):
//
// A productive pt-PT adverb is the FEMININE SINGULAR form of an adjective plus
// the suffix -mente, and it is INVARIABLE: it does not inflect for gender,
// number, or degree (rápido -> rápida -> rapidamente; fácil -> facilmente; the
// same -mente form serves every context). Because the suffix carries the word's
// primary stress, the base LOSES whatever graphic accent it carried in isolation
// (fácil -> facilmente, rápida -> rapidamente, só -> somente), and the whole
// adverb becomes a default paroxytone stressed on -men-te, which takes no graphic
// accent at all.
//
// Accent loss for free
//
// The adjective composition stores BARE syllables: Layer 3 (wordspt_accent.go)
// adds the acute or circumflex only during assembly (accentedWord). So building
// the adverb directly from the feminine adjective's BARE combined stem — never
// invoking Layer 3 — drops the base's graphic accent exactly as the rule
// requires, with no special-casing. The appended -mente is itself unaccented, so
// the assembled adverb carries no acute or circumflex anywhere. This mirrors how
// the superlative (wordspt_adjective.go) builds straight from the bare stem to
// shed the base accent. Nasal tildes on the bare syllables are part of the
// grapheme (nasality, not stress) and are preserved, as in any unstressed nasal.
//
// The junction between the feminine base and -mente never needs an orthographic
// adjustment: every feminine adjective ending terminates in a vowel (-a, -osa,
// -ica, -iva, -ente, -ante) or in the coda l (-al, -ável, -ível), and -mente
// begins with the consonant m, so no cedilla-before-front-vowel, hiatus, or
// coda/onset problem can arise. -mente is therefore a purely additive string
// suffix, validated (like the plural and the superlative) against the Sprint 7
// STRING oracles rather than the decomposition oracles.
//
// Invariability and length
//
// AdverbPT takes no gender, number, or degree option, because a -mente adverb is
// invariable; only length is selectable, through AdverbPTByLengthType. The
// shortest -mente adverb is the shortest feminine adjective base plus the five
// characters of "mente", which already exceeds the four-character ceiling of
// SmallLengthWord (and the eight-character ceiling of MediumLengthWords once the
// CV stem floor is counted), so an adverb requested in a small category always
// normalizes upward to the smallest viable category, per the specification's
// minimum-viable-length rule. Length is never normalized downward.
//
// Allocation. The syllable buffer is a fixed-size array on the caller's stack, so
// it does not allocate on the heap; a rejected length attempt is measured from
// the buffer with no allocation. An adverb therefore performs exactly one string
// allocation: adverbWord writes the bare stem and the -mente suffix into a single
// pre-sized builder.
//
// Randomness surfaces. The package-level AdverbPT/AdverbPTByLengthType draw from
// the global, automatically seeded math/rand/v2 source (through
// wordsPTGlobalRand); the (*Generator) methods draw from the generator's own
// seeded source. Both call the single shared core adverbCore with an injectable
// *rand.Rand, so there is no duplication and a seeded Generator is fully
// reproducible.
//
// Sources: Cunha, C. & Cintra, L., "Nova Gramática do Português Contemporâneo"
// (the productive -mente adverb, its formation from the feminine adjective, its
// invariability, and the loss of the base's graphic accent); Acordo Ortográfico
// da Língua Portuguesa (graphic accent placement and the paroxytone -mente
// ending taking no accent).

// ---------------------------------------------------------------------------
// -mente suffix
// ---------------------------------------------------------------------------

// adverbSuffix is the invariable adverbial suffix appended to the feminine
// singular adjective base to form a productive pt-PT adverb (belamente,
// facilmente, rapidamente). It carries the adverb's primary stress, which is why
// the base loses its own graphic accent.
const adverbSuffix = "mente"

// adverbSuffixRunes is the character (rune) length adverbSuffix adds to the base.
// "mente" is five ASCII characters, so its rune length equals its byte length,
// five. It sizes the leading stem and the minimum viable adverb length before any
// syllable is sampled. (The specification prose rounds this to "six characters";
// the literal suffix is "mente", five characters. The minimum-viable-length
// consequence — an adverb never fits SmallLengthWord — holds either way.)
const adverbSuffixRunes = 5

// ---------------------------------------------------------------------------
// -mente assembly
// ---------------------------------------------------------------------------

// adverbWord renders a productive -mente adverb from a bare combined adjective
// stem in a single allocation. It writes every syllable straight from the bare
// (unaccented) stem — so the adverb carries no base graphic accent — and appends
// the invariable -mente suffix. The stem's tonic index is intentionally ignored:
// the adverb's stress falls on -men-te, a default paroxytone that takes no graphic
// accent, so Layer 3 is never invoked. The result is lowercase and valid UTF-8,
// because every syllable part and the suffix are lowercase and valid UTF-8.
func adverbWord(stem syllabicStem) string {
	syllables := stem.syllables

	// Pre-size the builder to the exact byte length for a single allocation.
	total := len(adverbSuffix)
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
	b.WriteString(adverbSuffix)
	return b.String()
}

// ---------------------------------------------------------------------------
// Shared generation core
// ---------------------------------------------------------------------------

// adverbCore is the single generation engine shared by the package-level adverb
// functions and the *Generator adverb methods. It draws only from r, so a seeded
// Generator is fully reproducible, and reuses a caller-provided syllable buffer
// (backed by a stack array) so that only the final assembly allocates.
//
// A -mente adverb is invariable, so the core resolves no gender, number, or
// degree: it always builds the FEMININE SINGULAR positive adjective base. It
// selects a feminine adjective ending that fits the requested length, sizes the
// leading stem for that ending plus the -mente suffix (normalizing the length
// category up to the minimum viable length), then samples a combined stem whose
// assembled adverb falls in the category window. The output is always lowercase,
// valid UTF-8, orthographically well-formed, and unaccented (no acute or
// circumflex).
func adverbCore(r *rand.Rand, buf []syllable, l LengthTypeWords) string {
	_, ceiling := charRangeOf(l) // the requested category's character ceiling
	end := selectFittingEnding(r, feminineAdjectiveEndings, ceiling)

	// The -mente suffix is a fixed part of every adverb, so its characters count as
	// affix cost on top of the feminine adjective base for both length
	// normalization and stem sizing. This raises the minimum viable length above
	// every small category, so a small request normalizes upward.
	affixCost := end.affixRunes + adverbSuffixRunes
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

		runes := stemRuneLen(stem.syllables) + adverbSuffixRunes
		if runes >= lo && runes <= hi {
			break
		}
	}

	return adverbWord(stem)
}

// ---------------------------------------------------------------------------
// Public API: package-level functions
// ---------------------------------------------------------------------------

// AdverbPT returns a random European-Portuguese (pt-PT) pseudo-adverb: a
// productive -mente adverb of random length. It is exact sugar for
// AdverbPTByLengthType(AnyLengthWord). The result is a lowercase, valid-UTF-8
// word formed from the feminine singular of a generated adjective plus -mente; it
// is invariable and carries no graphic accent. It is a plausible invented word,
// not guaranteed to be a real dictionary entry.
//
// AdverbPT draws from the global, automatically seeded math/rand/v2 source and is
// safe for concurrent use by multiple goroutines. For reproducible output, use a
// seeded Generator ([New] or [NewSource]) and its AdverbPT method.
func AdverbPT() string {
	return AdverbPTByLengthType(AnyLengthWord)
}

// AdverbPTByLengthType returns a random pt-PT productive -mente adverb within
// length category l. [AnyLengthWord] (the zero value) selects a random valid
// length. A -mente adverb is invariable — it does not inflect for gender, number,
// or degree — so length is its only option. When l cannot hold any -mente adverb
// (its shortest form exceeds the category ceiling), l is normalized upward to the
// smallest category that can (length is never normalized downward); in particular
// an adverb never fits [SmallLengthWord] and normalizes up. The result is always
// lowercase and valid UTF-8; the function never returns an error and never panics.
//
// AdverbPTByLengthType draws from the global, automatically seeded math/rand/v2
// source and is safe for concurrent use by multiple goroutines.
func AdverbPTByLengthType(l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return adverbCore(wordsPTGlobalRand, buf[:0], l)
}

// ---------------------------------------------------------------------------
// Public API: *Generator methods
// ---------------------------------------------------------------------------

// AdverbPT is the seeded-generator equivalent of [AdverbPT]. A Generator created
// with the same seed produces the same sequence of adverbs.
func (g *Generator) AdverbPT() string {
	return g.AdverbPTByLengthType(AnyLengthWord)
}

// AdverbPTByLengthType is the seeded-generator equivalent of
// [AdverbPTByLengthType]. A Generator created with the same seed produces the
// same sequence of adverbs for the same length category.
func (g *Generator) AdverbPTByLengthType(l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return adverbCore(g.r, buf[:0], l)
}
