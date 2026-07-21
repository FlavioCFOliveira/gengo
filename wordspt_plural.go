package gengo

import (
	"math/rand/v2"
	"unicode/utf8"
)

// This file implements the European-Portuguese (pt-PT) plural-formation engine
// used by the WordsPT number inflection. It operates at the STRING level: given a
// singular pt-PT word and whether that word is oxytone (stressed on its final
// syllable), it dispatches on the word ending and returns the plural. The rules
// are encoded as data-driven branches keyed on the final grapheme(s), sourced
// authoritatively from Cunha, C. & Cintra, L., "Nova Gramática do Português
// Contemporâneo" (plural formation) and the Acordo Ortográfico da Língua
// Portuguesa (the -ões/-ães/-ãos and -éis/-óis graphic forms).
//
// Rule table (singular ending -> plural), from Cunha & Cintra:
//
//   - a vowel or an oral/nasal diphthong, or a lone nasal vowel (-ã, -õ):  add -s
//     (casa->casas, café->cafés, maçã->maçãs, mãe->mães).
//   - -ão:  -ões (default/productive), -ães or -ãos, chosen by weight
//     (coração->corações, pão->pães, mão->mãos, leão->leões).
//   - -m:  -ns  (homem->homens).
//   - -r, -z, -n:  add -es  (flor->flores, luz->luzes, líquen->líquenes).
//   - -x:  invariable  (tórax->tórax).
//   - -s, oxytone (stressed final syllable):  add -es  (país->países).
//   - -s, not oxytone:  invariable  (lápis->lápis).
//   - -al/-el/-ol/-ul:  -ais/-éis/-óis/-uis  (animal->animais, papel->papéis,
//     anzol->anzóis, azul->azuis).  When -el is UNstressed the acute is not
//     added: -eis  (amável->amáveis, possível->possíveis).
//   - -il, oxytone:  -is  (funil->funis).
//   - -il, not oxytone:  -eis  (fácil->fáceis, útil->úteis).
//
// Scope and validation: the engine is a string transform, so it is validated
// against known plural forms and the STRING-LEVEL orthographic oracles of Sprint
// 7 (wordConformsToCedilla, wordHasForbiddenDiaeresis, wordContainsGrave, valid
// UTF-8). The phonotactic and syllabic-accentuation validators apply to the
// SINGULAR stem produced by Layer 1, before pluralization, not to the
// string-transformed plural: a plural such as "animais" is not re-syllabified
// here, so it is not offered to stemConformsToAccentuation. A handful of real
// pt-PT plurals carry lexical accent shifts that a pseudo-word generator neither
// needs nor sees (for example the closed-vowel oxytone inglês->ingleses, or the
// hiatus accent raiz->raízes); these are documented limitations, out of range for
// the generator, whose accentuation layer never emits a circumflex and whose
// output the string oracles still accept.
//
// Allocation: each transforming branch performs a single string concatenation
// (one allocation); the invariable branches (-x, unstressed -s) return the input
// unchanged (zero allocations). Because the plural touches only the final
// grapheme(s), a class can equivalently pluralize its (short) class ending before
// final assembly, keeping the whole word's build within the single-allocation
// budget.

// aoPlural identifies one of the three plural outcomes of a singular ending in
// -ão. In real pt-PT the outcome is lexical (memorized per word); for a generated
// pseudo-word it is chosen by weight, with -ões as the productive default.
type aoPlural uint8

const (
	// aoOes is the -ões plural (productive default): leão -> leões.
	aoOes aoPlural = iota
	// aoAes is the -ães plural: pão -> pães.
	aoAes
	// aoAos is the -ãos plural: mão -> mãos.
	aoAos
)

// The -ão plural weights are ordinal estimates (flagged, consistent with the
// other WordsPT sampling weights): -ões is by far the most productive outcome in
// pt-PT, while -ães and -ãos apply to smaller closed sets (Cunha & Cintra). The
// exact proportions do not affect correctness, only how often each valid form
// appears for a pseudo-word.
const (
	aoOesWeight   = 78
	aoAosWeight   = 12
	aoAesWeight   = 10
	aoTotalWeight = aoOesWeight + aoAosWeight + aoAesWeight
)

// sampleAoPlural draws a weighted -ão plural outcome from the injectable source
// r, with -ões dominant. It performs a single bounded draw and never rejects.
func sampleAoPlural(r *rand.Rand) aoPlural {
	switch v := r.Uint32N(aoTotalWeight); {
	case v < aoOesWeight:
		return aoOes
	case v < aoOesWeight+aoAosWeight:
		return aoAos
	default:
		return aoAes
	}
}

// pluralize returns the pt-PT plural of the singular word, given whether the
// singular is oxytone (stressed on its final syllable). It dispatches on the word
// ending per the rule table at the top of this file, drawing from r only for the
// weighted -ão outcome. It never errors and never panics: an unrecognized ending
// falls back to adding -s. The empty string pluralizes to the empty string.
func pluralize(r *rand.Rand, singular string, oxytone bool) string {
	if singular == "" {
		return ""
	}

	last, _ := utf8.DecodeLastRuneInString(singular)
	prefix := trimLastRunes(singular, 1) // singular without its final rune
	var prev rune
	if prefix != "" {
		prev, _ = utf8.DecodeLastRuneInString(prefix)
	}

	switch {
	case last == 'o' && prev == 'ã': // -ão
		return applyAoPlural(singular, sampleAoPlural(r))
	case last == 'm': // -m -> -ns
		return prefix + "ns"
	case last == 'l': // -al/-el/-ol/-ul/-il
		if base, ok := vowelBase(prev); ok {
			return pluralizeL(singular, base, oxytone)
		}
		return singular + "es" // defensive: -l not preceded by a vowel
	case last == 'r' || last == 'z' || last == 'n': // -r/-z/-n -> +es
		return singular + "es"
	case last == 'x': // -x invariable
		return singular
	case last == 's': // -s: oxytone +es, else invariable
		if oxytone {
			return singular + "es"
		}
		return singular
	default: // vowel or diphthong (oral or nasal) -> +s; safe fallback otherwise
		return singular + "s"
	}
}

// applyAoPlural replaces a singular -ão ending with the chosen plural outcome. It
// strips the two final runes ("ã" + "o") and appends the -ões, -ães or -ãos
// ending. It is the pure, deterministic transform behind the weighted
// [sampleAoPlural] selection, so each outcome can be tested against its canonical
// pt-PT example (leão->leões, pão->pães, mão->mãos).
func applyAoPlural(singular string, v aoPlural) string {
	stem := trimLastRunes(singular, 2) // drop the "ão"
	switch v {
	case aoAes:
		return stem + "ães"
	case aoAos:
		return stem + "ãos"
	default: // aoOes: the productive default
		return stem + "ões"
	}
}

// pluralizeL forms the plural of a word ending in a vowel plus -l, given the base
// vowel immediately before the final l and whether the word is oxytone. The final
// "vowel + l" is replaced by the corresponding diphthong-plus-s ending; the acute
// on -éis and -óis marks the open tonic vowel of an oxytone, and is not added when
// the ending is unstressed (amável->amáveis). The -il ending is special: an
// oxytone yields -is (funil->funis) and an unstressed one yields -eis
// (fácil->fáceis).
func pluralizeL(singular string, baseVowel rune, oxytone bool) string {
	stem := trimLastRunes(singular, 2) // drop the vowel + l
	switch baseVowel {
	case 'a':
		return stem + "ais"
	case 'e':
		if oxytone {
			return stem + "éis"
		}
		return stem + "eis"
	case 'i':
		if oxytone {
			return stem + "is"
		}
		return stem + "eis"
	case 'o':
		if oxytone {
			return stem + "óis"
		}
		return stem + "ois"
	case 'u':
		return stem + "uis"
	default:
		return singular + "es" // defensive: unreachable for a valid -l ending
	}
}

// trimLastRunes returns s without its last n runes, operating on the UTF-8 bytes
// with no allocation (it returns a sub-slice of s). It stops early if s becomes
// empty.
func trimLastRunes(s string, n int) string {
	for i := 0; i < n && s != ""; i++ {
		_, size := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-size]
	}
	return s
}
