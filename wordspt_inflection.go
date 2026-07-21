package gengo

import (
	"math"
	"math/rand/v2"
)

// This file establishes the public WordsPT flexion (inflection) enums shared by
// every open nominal class (noun, adjective, adverb) and the internal, reusable
// machinery those classes build on:
//
//   - the flexion enums Gender, Number and Degree, each a uint8 whose zero value
//     is "Any" ("choose at random among the valid options");
//   - the Any-resolution helpers that turn an Any trait into a concrete random
//     value using an injectable *rand.Rand source, so the package-level functions
//     and the *Generator methods share one core;
//   - the length machinery that maps a LengthTypeWords category to a target
//     syllable count and normalizes the category UP to a class-supplied minimum
//     viable length (never downward), as required by the WordsPT specification.
//
// The Portuguese plural-formation engine that the number inflection uses lives in
// wordspt_plural.go. The stem sampler (wordspt_stem.go) and the graphic
// accentuation (wordspt_accent.go) are the Sprint 7 foundation this machinery
// reuses. Sources for the linguistic rules are cited where they apply and, for
// the plural rules, in wordspt_plural.go (Cunha & Cintra; Acordo Ortográfico).
//
// Scope: this file introduces only the enums and the shared machinery. The public
// class generators (NounPT, AdjectivePT, AdverbPT) and their class endings are
// built by later tasks on top of these primitives.

// ---------------------------------------------------------------------------
// Public flexion enums
// ---------------------------------------------------------------------------

// Gender is the grammatical gender option of a generated nominal pt-PT word. Its
// zero value, [AnyGender], means "choose a valid gender at random".
type Gender uint8

const (
	// AnyGender is the zero value of [Gender] and selects a random valid gender.
	AnyGender Gender = iota
	// Masculine selects the masculine form of a gender-inflecting word.
	Masculine
	// Feminine selects the feminine form of a gender-inflecting word.
	Feminine
)

// Number is the grammatical number option of a generated nominal pt-PT word. Its
// zero value, [AnyNumber], means "choose a valid number at random".
type Number uint8

const (
	// AnyNumber is the zero value of [Number] and selects a random valid number.
	AnyNumber Number = iota
	// Singular selects the singular form.
	Singular
	// Plural selects the plural form, built by the Portuguese plural engine.
	Plural
)

// Degree is the degree option of a generated adjective. Its zero value,
// [AnyDegree], means "choose a valid degree at random". Only the positive and the
// synthetic absolute superlative (-íssimo) are supported; the comparative and the
// analytic superlatives are multi-word forms and are out of scope.
type Degree uint8

const (
	// AnyDegree is the zero value of [Degree] and selects a random valid degree.
	AnyDegree Degree = iota
	// Positive selects the base (positive) degree.
	Positive
	// Superlative selects the synthetic absolute superlative (-íssimo).
	Superlative
)

// ---------------------------------------------------------------------------
// Any-resolution helpers
// ---------------------------------------------------------------------------
//
// Each helper turns an unspecified trait (the Any zero value, or any value
// outside the defined concrete set) into a concrete, uniformly random valid
// value drawn from the injectable source r; a concrete value is returned
// unchanged. Resolving through the injected source is what lets the package-level
// functions (global source) and the *Generator methods (own source) share one
// generation core. The helpers never error and never panic: an out-of-range enum
// value is treated as Any, so an impossible request normalizes to a valid form.

// resolveGender resolves g to a concrete [Gender]. AnyGender (and any undefined
// value) becomes Masculine or Feminine with equal probability.
func resolveGender(r *rand.Rand, g Gender) Gender {
	switch g {
	case Masculine, Feminine:
		return g
	default:
		if r.Uint32()&1 == 0 {
			return Masculine
		}
		return Feminine
	}
}

// resolveNumber resolves n to a concrete [Number]. AnyNumber (and any undefined
// value) becomes Singular or Plural with equal probability.
func resolveNumber(r *rand.Rand, n Number) Number {
	switch n {
	case Singular, Plural:
		return n
	default:
		if r.Uint32()&1 == 0 {
			return Singular
		}
		return Plural
	}
}

// resolveDegree resolves d to a concrete [Degree]. AnyDegree (and any undefined
// value) becomes Positive or Superlative with equal probability.
func resolveDegree(r *rand.Rand, d Degree) Degree {
	switch d {
	case Positive, Superlative:
		return d
	default:
		if r.Uint32()&1 == 0 {
			return Positive
		}
		return Superlative
	}
}

// ---------------------------------------------------------------------------
// Length machinery: category window, normalization, and syllable mapping
// ---------------------------------------------------------------------------

// runesPerSyllable is the mean rune (character) length of one generated pt-PT
// syllable on the accented-stem path, measured empirically over the Sprint 7
// sampler: 200000 samples per syllable count (1..12) yield a stable mean of
// 2.47–2.51 runes/syllable. It is the inversion constant used to translate a
// target character window into a target syllable count. Because individual
// syllables range from 1 to 5 runes, this maps the center of a window; the
// caller enforces the exact window on the assembled word.
const runesPerSyllable = 2.48

// charRangeOf returns the inclusive character (rune) window of a length category.
// SmallLengthWord is 1..4, MediumLengthWords is 5..8 and BigLengthWords is 9..30;
// AnyLengthWord (the zero value) and any undefined value select the full 1..30
// range, matching [WordByLengthType].
func charRangeOf(l LengthTypeWords) (minChars, maxChars int) {
	switch l {
	case SmallLengthWord:
		return 1, 4
	case MediumLengthWords:
		return 5, 8
	case BigLengthWords:
		return 9, 30
	default:
		return 1, 30
	}
}

// normalizeLength returns the effective length category for a class whose
// shortest producible word is minViable characters. When the requested category's
// upper bound is below minViable, it is raised to the smallest category whose
// upper bound is at least minViable; the category is never lowered. AnyLengthWord
// (upper bound 30) is returned unchanged for any realistic minViable. This
// implements the specification's minimum-viable-length rule ("normalize the
// length upward ... length is never normalized downward").
func normalizeLength(l LengthTypeWords, minViable int) LengthTypeWords {
	_, maxChars := charRangeOf(l)
	if minViable <= maxChars {
		return l
	}
	switch {
	case minViable <= 4:
		return SmallLengthWord
	case minViable <= 8:
		return MediumLengthWords
	default:
		return BigLengthWords
	}
}

// syllablesForCharRange maps a character (rune) window to a syllable count. It
// draws a target character count uniformly in [minChars, maxChars] from r and
// inverts the measured [runesPerSyllable], returning a count of at least one. The
// uniform draw spreads generated lengths across the window; the returned count is
// a target, because a generated syllable spans 1..5 runes and the exact length is
// enforced by the caller on the assembled word.
func syllablesForCharRange(r *rand.Rand, minChars, maxChars int) int {
	if minChars < 1 {
		minChars = 1
	}
	if maxChars < minChars {
		maxChars = minChars
	}
	target := minChars + int(r.Uint32N(uint32(maxChars-minChars+1)))
	count := int(math.Round(float64(target) / runesPerSyllable))
	if count < 1 {
		count = 1
	}
	return count
}

// stemSyllablesFor sizes the syllabic stem for a nominal word of length category
// l. It is the shared length-mapping mechanism, parameterized by the class: the
// class supplies minViable (the character length of its shortest producible word)
// and affixRunes (the characters its fixed ending and inflection add on top of
// the bare stem). stemSyllablesFor normalizes l up to minViable, targets a word
// length within the resulting window (never below minViable), subtracts the affix
// cost, and inverts [runesPerSyllable] to a stem syllable count of at least one.
// It draws only from r, so it is reproducible through a seeded *Generator source.
func stemSyllablesFor(r *rand.Rand, l LengthTypeWords, minViable, affixRunes int) int {
	if minViable < 1 {
		minViable = 1
	}
	if affixRunes < 0 {
		affixRunes = 0
	}
	minChars, maxChars := charRangeOf(normalizeLength(l, minViable))
	if minChars < minViable {
		minChars = minViable
	}
	if maxChars < minChars {
		maxChars = minChars
	}
	// Size the stem: the window minus the fixed affix cost, keeping at least one
	// stem character so the stem is never empty.
	stemMin := minChars - affixRunes
	if stemMin < 1 {
		stemMin = 1
	}
	stemMax := maxChars - affixRunes
	if stemMax < stemMin {
		stemMax = stemMin
	}
	return syllablesForCharRange(r, stemMin, stemMax)
}
