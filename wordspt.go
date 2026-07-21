package gengo

// This file implements the WordsPT conformance validators: pure, allocation-light
// oracle functions that decide whether a decomposed pt-PT syllable, stem or word
// obeys the European-Portuguese phonotactic rules and the pt-PT graphic-accent
// rules. They serve as the empirical test oracle for the whole WordsPT feature:
// every later generation task is verified against them, and every regression is
// caught by them.
//
// The static data model these validators consult (inventories, weights,
// positional constraints and rune helpers) lives in wordspt_data.go. Linguistic
// sources are cited there and, per validator, below.
//
// Scope of this task: the phonotactic core and the conformance validators only.
// The validators operate on a decomposed representation (onset/nucleus/coda,
// plus a tonic-syllable index) because that is exactly what the generation
// layers produce and consume; it also avoids ambiguous re-syllabification of raw
// strings. The cedilla and diaeresis validators, which do not need
// syllabification, operate directly on the assembled rune string.

// syllable is a decomposed pt-PT syllable of the form (onset)(nucleus)(coda).
// The onset and coda are optional (empty when absent); the nucleus is
// mandatory. Every field holds lowercase pt-PT graphemes matching the inventory
// forms in wordspt_data.go (for example onset "tr", nucleus "ão", coda "s").
type syllable struct {
	onset   string
	nucleus string
	coda    string
}

// ---------------------------------------------------------------------------
// Phonotactic conformance
// ---------------------------------------------------------------------------

// onsetIsValid reports whether onset is a legal pt-PT onset. The empty string is
// valid, because the onset is optional. Any other value must be a member of the
// onset inventory (a single consonant, a digraph, or an obstruent+liquid
// cluster).
func onsetIsValid(onset string) bool {
	if onset == "" {
		return true
	}
	_, ok := onsetSet[onset]
	return ok
}

// nucleusIsValid reports whether nucleus is a legal pt-PT nucleus: an oral
// vowel, an oral diphthong, a nasal vowel, or a nasal diphthong from the nucleus
// inventory. The nucleus is mandatory, so the empty string is invalid.
func nucleusIsValid(nucleus string) bool {
	_, ok := nucleusSet[nucleus]
	return ok
}

// codaIsValid reports whether coda is a legal pt-PT coda. The empty string is
// valid, because the coda is optional. Any other value must be one of the seven
// permitted coda consonants s, r, l, m, n, z, x.
func codaIsValid(coda string) bool {
	if coda == "" {
		return true
	}
	_, ok := codaSet[coda]
	return ok
}

// syllableIsValid reports whether s is a well-formed pt-PT syllable. It checks
// that each part is drawn from the legal inventory and that the onset obeys its
// positional constraints:
//
//   - lh, nh and ç never appear word-initially (lh/nh occur only between
//     vowels; a word never begins with ç);
//   - qu and gu appear only before a front vowel (e or i);
//   - ç appears only before a back or central vowel (a, o or u), per the Acordo
//     Ortográfico cedilla rule.
//
// wordInitial must be true when s is the first syllable of a stem, so that the
// word-initial onset constraint can be applied. Sources: Cunha & Cintra; Mateus
// & d'Andrade; Acordo Ortográfico.
func syllableIsValid(s syllable, wordInitial bool) bool {
	if !nucleusIsValid(s.nucleus) {
		return false
	}
	if !codaIsValid(s.coda) {
		return false
	}
	if !onsetIsValid(s.onset) {
		return false
	}
	if s.onset == "" {
		return true
	}
	if wordInitial {
		if _, forbidden := onsetInitialForbidden[s.onset]; forbidden {
			return false
		}
	}
	if _, front := onsetFrontVowelOnly[s.onset]; front && !nucleusStartsFront(s.nucleus) {
		return false
	}
	if _, back := onsetBackVowelOnly[s.onset]; back && !nucleusStartsBack(s.nucleus) {
		return false
	}
	return true
}

// nucleusStartsFront reports whether the nucleus begins with a front vowel (e or
// i), which is the context that licenses the qu and gu onsets.
func nucleusStartsFront(nucleus string) bool {
	b, ok := vowelBase(firstRune(nucleus))
	return ok && (b == 'e' || b == 'i')
}

// nucleusStartsBack reports whether the nucleus begins with a back or central
// vowel (a, o or u), which is the context that licenses the ç onset.
func nucleusStartsBack(nucleus string) bool {
	b, ok := vowelBase(firstRune(nucleus))
	return ok && (b == 'a' || b == 'o' || b == 'u')
}

// transitionIsValid reports whether a coda may be immediately followed by the
// onset of the next syllable within a stem. It encodes the constraints that are
// orthographically certain in pt-PT:
//
//   - The nasal spelling rule (Acordo Ortográfico): a coda m occurs only before
//     p or b; a coda n never occurs before p or b (there it would be spelled m).
//   - The Maximum Onset Principle: a non-empty coda must be followed by a
//     consonant onset. A coda before a vowel-initial syllable would re-syllabify
//     as that syllable's onset, so such a transition is not well-formed.
//
// Other coda/onset pairs are permitted by this core (pt-PT is highly permissive
// at these boundaries). The exhaustive transition table is a flagged item for
// authoritative expansion; see the task report. nextOnset is the onset of the
// following syllable ("" when that syllable is vowel-initial). This function is
// only meaningful between two syllables, never for a word-final coda.
func transitionIsValid(coda, nextOnset string) bool {
	if coda == "" {
		return true
	}
	// Maximum Onset Principle: a coda must be followed by a consonant onset.
	if nextOnset == "" {
		return false
	}
	switch coda {
	case "m":
		// Coda m only before p or b.
		return nextOnset[0] == 'p' || nextOnset[0] == 'b'
	case "n":
		// Coda n never before p or b.
		return nextOnset[0] != 'p' && nextOnset[0] != 'b'
	}
	return true
}

// stemConformsToPhonotactics reports whether stem is a well-formed pt-PT
// syllable sequence: every syllable is valid in its position and every
// inter-syllable transition is legal. An empty stem is not conformant.
func stemConformsToPhonotactics(stem []syllable) bool {
	if len(stem) == 0 {
		return false
	}
	for i := range stem {
		if !syllableIsValid(stem[i], i == 0) {
			return false
		}
		if i < len(stem)-1 && !transitionIsValid(stem[i].coda, stem[i+1].onset) {
			return false
		}
	}
	return true
}

// ---------------------------------------------------------------------------
// Accentuation conformance
// ---------------------------------------------------------------------------

// stressClass classifies a word by the position of its tonic syllable.
type stressClass uint8

const (
	// oxytone (aguda): stress on the last syllable.
	oxytone stressClass = iota
	// paroxytone (grave): stress on the second-to-last syllable.
	paroxytone
	// proparoxytone (esdrúxula): stress on the third-to-last syllable or earlier.
	proparoxytone
)

// stressClassOf returns the stress class of a word with syllableCount syllables
// whose tonic syllable is at the zero-based index tonic. A stress position three
// or more syllables from the end is classified as proparoxytone.
func stressClassOf(syllableCount, tonic int) stressClass {
	switch syllableCount - 1 - tonic {
	case 0:
		return oxytone
	case 1:
		return paroxytone
	default:
		return proparoxytone
	}
}

// tonicRequiresGraphicAccent reports whether the tonic vowel of an oral-nucleus
// word must carry a graphic accent (acute or circumflex) under the pt-PT rules
// of the Acordo Ortográfico. It is defined for a word whose tonic nucleus is
// oral; nasal tonic nuclei are marked by the tilde and are handled by
// [stemConformsToAccentuation]. The rules encoded:
//
//   - Proparoxytones always require an accent.
//   - Oxytones require an accent when the word ends in a tonic a, e or o
//     (optionally followed by a coda s), or in the nasal sequence -em.
//   - Paroxytones require an accent when the word does NOT end in one of the
//     default paroxytone endings (see [isDefaultParoxytoneEnding]); notably,
//     paroxytones ending in a nasal vowel or diphthong (ã, ão) are accented on
//     the tonic vowel, as in órfã and órgão.
//
// Sources: Acordo Ortográfico da Língua Portuguesa, Bases IX–XI; Cunha & Cintra.
func tonicRequiresGraphicAccent(stem []syllable, tonic int) bool {
	n := len(stem)
	last := stem[n-1]
	switch stressClassOf(n, tonic) {
	case proparoxytone:
		return true
	case oxytone:
		base, ok := lastVowelBase(last.nucleus)
		if !ok {
			return false
		}
		if !nucleusIsNasal(last.nucleus) {
			if base == 'a' || base == 'e' || base == 'o' {
				if last.coda == "" || last.coda == "s" {
					return true
				}
			}
			// The nasal sequence -em (oral e nucleus plus m coda): também, porém.
			if base == 'e' && last.coda == "m" {
				return true
			}
		}
		return false
	default: // paroxytone
		return !isDefaultParoxytoneEnding(last)
	}
}

// isDefaultParoxytoneEnding reports whether the final syllable is one of the
// default paroxytone endings that carry no graphic accent: an oral a, e or o
// optionally closed by a coda s (plural) or m (nasal -am/-em), as in casa,
// casas, homem and falam. Every other ending is non-default and receives an
// accent on the tonic vowel; in particular, a nasal (tilde) ending is
// non-default, because paroxytones ending in ã or ão are accented (órfã, órgão).
func isDefaultParoxytoneEnding(last syllable) bool {
	if nucleusIsNasal(last.nucleus) {
		return false
	}
	base, ok := lastVowelBase(last.nucleus)
	if !ok {
		return false
	}
	if base == 'a' || base == 'e' || base == 'o' {
		switch last.coda {
		case "", "s", "m":
			return true
		}
	}
	return false
}

// stemConformsToAccentuation reports whether an accented stem (its nucleus
// graphemes already bearing any acute, circumflex or tilde) is orthographically
// correct for its stress position, given the zero-based tonic index. It checks:
//
//   - Placement: only the tonic oral nucleus may carry an acute or circumflex
//     accent; no other syllable may. A nasal tilde may appear on any syllable
//     (for example the unstressed ão of órgão).
//   - Nasal tonic: when the tonic nucleus is nasal, the tilde marks the stress,
//     so the nucleus must carry the tilde and must not also carry an acute or
//     circumflex accent.
//   - Oral tonic: the presence of an acute or circumflex accent on the tonic
//     nucleus must match exactly what [tonicRequiresGraphicAccent] demands.
//
// The choice between acute (open) and circumflex (closed) depends on the vowel
// quality recorded during generation; validating that quality is deferred to the
// generation tasks that carry the quality datum. Sources: Acordo Ortográfico,
// Bases IX–XIII; Cunha & Cintra.
func stemConformsToAccentuation(stem []syllable, tonic int) bool {
	n := len(stem)
	if n == 0 || tonic < 0 || tonic >= n {
		return false
	}
	// Placement: acute/circumflex only on the tonic syllable.
	for i := range stem {
		if i != tonic && nucleusHasAcuteOrCircumflex(stem[i].nucleus) {
			return false
		}
	}
	tonicNucleus := stem[tonic].nucleus
	if nucleusIsNasal(tonicNucleus) {
		return !nucleusHasAcuteOrCircumflex(tonicNucleus)
	}
	return tonicRequiresGraphicAccent(stem, tonic) == nucleusHasAcuteOrCircumflex(tonicNucleus)
}

// nucleusIsNasal reports whether the nucleus is a nasal vowel or nasal diphthong,
// detected by the presence of a tilde (ã or õ) in its spelling.
func nucleusIsNasal(nucleus string) bool {
	for _, r := range nucleus {
		if hasTilde(r) {
			return true
		}
	}
	return false
}

// nucleusHasAcuteOrCircumflex reports whether any vowel of the nucleus bears an
// acute or circumflex accent.
func nucleusHasAcuteOrCircumflex(nucleus string) bool {
	for _, r := range nucleus {
		if hasAcute(r) || hasCircumflex(r) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// String-level orthographic conformance (cedilla and diaeresis)
// ---------------------------------------------------------------------------

// wordConformsToCedilla reports whether every cedilla in word obeys the pt-PT
// cedilla rule of the Acordo Ortográfico: a ç must be immediately followed by a,
// o or u (any accented form counts by its base vowel) and must never appear
// word-initially. A ç before e or i, a word-initial ç, or a word-final ç is
// non-conformant. A word without any ç is trivially conformant. The scan is
// allocation-free (it ranges over the string's runes).
func wordConformsToCedilla(word string) bool {
	first := true
	pendingAfterCedilla := false
	for _, r := range word {
		if pendingAfterCedilla {
			base, ok := vowelBase(r)
			if !ok || (base != 'a' && base != 'o' && base != 'u') {
				return false
			}
			pendingAfterCedilla = false
		}
		if r == 'ç' {
			if first {
				return false
			}
			pendingAfterCedilla = true
		}
		first = false
	}
	// A ç with nothing after it (word-final) is non-conformant.
	return !pendingAfterCedilla
}

// wordHasForbiddenDiaeresis reports whether word contains a diaeresis (trema),
// which the Acordo Ortográfico abolished in pt-PT and which is therefore illegal
// in a generated word. The scan is allocation-free.
func wordHasForbiddenDiaeresis(word string) bool {
	for _, r := range word {
		if hasDiaeresis(r) {
			return true
		}
	}
	return false
}

// wordContainsGrave reports whether word contains a grave accent (à). In pt-PT
// the grave accent marks only the crasis contraction, which WordsPT does not
// generate in open-class words; an open-class output must therefore contain no
// grave accent. The scan is allocation-free.
func wordContainsGrave(word string) bool {
	for _, r := range word {
		if hasGrave(r) {
			return true
		}
	}
	return false
}
