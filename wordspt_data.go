package gengo

// This file holds the static, embedded European-Portuguese (pt-PT) phonotactic
// data model for the WordsPT feature: the syllable-part inventories with their
// sampling weights, the positional constraints, and the rune-level helpers that
// classify pt-PT vowels and diacritics.
//
// The data model is deliberately organized for fast, correct-by-construction
// sampling in a later task: every inventory is a [weightedInventory] that
// precomputes a cumulative-weight (prefix-sum) table, so a downstream sampler
// draws a single random value and binary-searches the table with no
// generate-and-reject loop. This file introduces no run-time allocation on the
// sampling or validation paths and adds no external dependency.
//
// Linguistic sources (never guessed; see the task report for flagged items):
//   - Acordo Ortográfico da Língua Portuguesa (1990): the graphic-accent,
//     cedilla, tilde and diaeresis rules (Bases IX–XIV) and the m/n nasal
//     spelling convention.
//   - Cunha, C. & Cintra, L., "Nova Gramática do Português Contemporâneo":
//     syllable structure, plural formation and accentuation.
//   - Mateus, M. H. & d'Andrade, E., "The Phonology of Portuguese" (Oxford,
//     2000): the syllable template (onset)(nucleus)(coda), the obstruent+liquid
//     onset clusters, the coda inventory and the Maximum Onset Principle.
//
// Frequency weights: the single-consonant onset and oral-vowel nucleus weights
// are derived, ordinally, from published per-mille letter/vowel frequencies of
// Portuguese running text (used here only to make common sounds appear more
// often than rare ones). The cluster, digraph, diphthong and nasal weights are
// ordinal estimates. The exact weights are flagged for authoritative review,
// consistent with Open Item 2 of specification/wordspt.md; they do not affect
// the conformance validators, only the relative sampling frequency.

// weightedForm is one member of a phonotactic inventory together with its
// sampling weight. A larger weight means the form is drawn more often. Every
// weight must be strictly positive so that every form stays reachable.
type weightedForm struct {
	form   string
	weight uint32
}

// weightedInventory is a sampling-ready inventory of weighted forms. It stores a
// precomputed cumulative-weight table so that a caller can select a form in
// O(log n) with no rejection loop: draw r in [0,total) and call [sampleIndex].
type weightedInventory struct {
	forms      []weightedForm
	cumulative []uint32 // cumulative[i] == sum of weights[0..i]; strictly increasing
	total      uint32   // total == cumulative[len-1]; the exclusive upper bound for r
}

// newWeightedInventory builds a [weightedInventory] from forms, precomputing the
// cumulative-weight table and the total weight. It is called once, at package
// initialization, for each static inventory.
func newWeightedInventory(forms []weightedForm) weightedInventory {
	cumulative := make([]uint32, len(forms))
	var total uint32
	for i := range forms {
		total += forms[i].weight
		cumulative[i] = total
	}
	return weightedInventory{forms: forms, cumulative: cumulative, total: total}
}

// sampleIndex maps r, a value in [0,total), to the index of the selected form.
// It performs a branch-lean binary search over the cumulative-weight table and
// never rejects: every r in range yields exactly one valid index, and every
// form with a positive weight is reachable. The result is the smallest index i
// for which cumulative[i] > r.
func (w *weightedInventory) sampleIndex(r uint32) int {
	lo, hi := 0, len(w.cumulative)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if w.cumulative[mid] <= r {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

// ---------------------------------------------------------------------------
// Onset inventory
// ---------------------------------------------------------------------------

// onsetSingles lists the single-consonant and digraph onsets of pt-PT. It
// includes the digraphs ch, lh and nh and the front-vowel sequences qu and gu.
// Positional constraints (lh/nh never word-initial; qu/gu only before e or i;
// ç only before a, o or u) are enforced by the validators, not by this list, so
// that the list stays a plain weighted inventory for sampling.
//
// Weights: consonant weights are ordinal, derived from Portuguese letter
// frequency; the digraph, qu/gu and ç weights are ordinal estimates (flagged).
var onsetSingles = newWeightedInventory([]weightedForm{
	{"p", 25}, {"b", 10}, {"t", 43}, {"d", 50},
	{"c", 39}, {"ç", 6}, {"g", 13}, {"f", 10},
	{"v", 17}, {"s", 78}, {"z", 5}, {"j", 4},
	{"l", 28}, {"m", 47}, {"n", 50}, {"r", 65},
	{"x", 2}, {"ch", 8}, {"lh", 6}, {"nh", 6},
	{"qu", 12}, {"gu", 4},
})

// onsetClusters lists the tautosyllabic onset clusters of pt-PT: an obstruent
// followed by a liquid (r or l). This is exactly the set permitted by the
// specification and by Mateus & d'Andrade; no other onset cluster is legal.
//
// Weights: ordinal estimates (flagged). Liquid-r clusters are generally more
// frequent than liquid-l clusters; vr is the rarest.
var onsetClusters = newWeightedInventory([]weightedForm{
	{"pr", 14}, {"br", 10}, {"tr", 14}, {"dr", 8},
	{"cr", 10}, {"gr", 12}, {"fr", 8}, {"vr", 2},
	{"pl", 10}, {"bl", 6}, {"cl", 8}, {"gl", 4},
	{"fl", 6},
})

// ---------------------------------------------------------------------------
// Nucleus inventory
// ---------------------------------------------------------------------------

// nuclei lists the syllable nuclei of pt-PT: the five oral vowels, the common
// oral (falling) diphthongs, the tilde nasal vowels, and the three nasal
// diphthongs. Nasalisation produced by a nasal coda (for example the -am/-en
// sequences) is modeled as an oral nucleus plus an m or n coda, not as a
// nucleus entry here.
//
// Weights: oral-vowel weights are ordinal, derived from Portuguese vowel
// frequency; diphthong and nasal weights are ordinal estimates (flagged). The
// standalone nasal vowel õ is rare outside the õe diphthong and carries a low
// weight (flagged).
var nuclei = newWeightedInventory([]weightedForm{
	// Oral vowels.
	{"a", 146}, {"e", 126}, {"o", 107}, {"i", 62}, {"u", 46},
	// Oral (falling) diphthongs.
	{"ei", 20}, {"ai", 12}, {"ou", 12}, {"au", 8},
	{"oi", 6}, {"eu", 6}, {"iu", 3}, {"ui", 3},
	// Nasal vowels (tilde).
	{"ã", 8}, {"õ", 2},
	// Nasal diphthongs (tilde).
	{"ão", 24}, {"ãe", 4}, {"õe", 4},
})

// ---------------------------------------------------------------------------
// Coda inventory
// ---------------------------------------------------------------------------

// codas lists the permitted syllable codas of pt-PT. The coda position is
// restricted to the seven consonants s, r, l, m, n, z and x; no other coda is
// legal. The nasal codas m and n additionally obey the m/n spelling rule
// enforced by [transitionIsValid].
//
// Weights: ordinal estimates (flagged). The sibilant s (plurals, verb endings)
// and r (infinitives, the -or ending) dominate; x is the rarest coda.
var codas = newWeightedInventory([]weightedForm{
	{"s", 80}, {"r", 50}, {"m", 40}, {"l", 20},
	{"n", 10}, {"z", 8}, {"x", 2},
})

// ---------------------------------------------------------------------------
// Membership sets (derived from the inventories, for O(1) validation lookups)
// ---------------------------------------------------------------------------

// onsetSet holds every legal onset form (singles, digraphs and clusters) for
// O(1) membership tests by the validators.
var onsetSet = mergeFormSet(onsetSingles, onsetClusters)

// nucleusSet holds every legal nucleus form for O(1) membership tests.
var nucleusSet = formSet(nuclei)

// codaSet holds every legal coda form for O(1) membership tests.
var codaSet = formSet(codas)

// formSet returns a set of the forms of one inventory.
func formSet(inv weightedInventory) map[string]struct{} {
	set := make(map[string]struct{}, len(inv.forms))
	for i := range inv.forms {
		set[inv.forms[i].form] = struct{}{}
	}
	return set
}

// mergeFormSet returns a set of the forms of several inventories combined.
func mergeFormSet(invs ...weightedInventory) map[string]struct{} {
	set := make(map[string]struct{})
	for i := range invs {
		for j := range invs[i].forms {
			set[invs[i].forms[j].form] = struct{}{}
		}
	}
	return set
}

// ---------------------------------------------------------------------------
// Positional onset constraints
// ---------------------------------------------------------------------------

// onsetInitialForbidden holds the onsets that may not appear word-initially. The
// palatal digraphs lh and nh occur only between vowels in pt-PT (Cunha &
// Cintra; Mateus & d'Andrade), and the cedilla onset ç never begins a word
// (Acordo Ortográfico cedilla rule).
var onsetInitialForbidden = map[string]struct{}{
	"lh": {},
	"nh": {},
	"ç":  {},
}

// onsetFrontVowelOnly holds the onsets that are legal only before a front vowel
// (e or i): the sequences qu and gu, whose /k/ and /ɡ/ readings the tables model
// before e and i (specification, Layer 1).
var onsetFrontVowelOnly = map[string]struct{}{
	"qu": {},
	"gu": {},
}

// onsetBackVowelOnly holds the onsets that are legal only before a back or
// central vowel (a, o or u): the cedilla onset ç, per the Acordo Ortográfico
// cedilla rule (ç only before a, o or u, never before e or i).
var onsetBackVowelOnly = map[string]struct{}{
	"ç": {},
}

// ---------------------------------------------------------------------------
// Rune-level pt-PT vowel and diacritic classification
// ---------------------------------------------------------------------------

// vowelBase returns the base oral vowel of a pt-PT vowel rune, ignoring any
// graphic accent, and reports whether r is a vowel at all. For example both
// 'á' and 'â' map to 'a'. Consonants and any non-vowel rune report false.
func vowelBase(r rune) (rune, bool) {
	switch r {
	case 'a', 'á', 'à', 'â', 'ã':
		return 'a', true
	case 'e', 'é', 'ê':
		return 'e', true
	case 'i', 'í':
		return 'i', true
	case 'o', 'ó', 'ô', 'õ':
		return 'o', true
	case 'u', 'ú':
		return 'u', true
	}
	return 0, false
}

// hasAcute reports whether r is a vowel bearing the acute accent (open quality).
func hasAcute(r rune) bool {
	switch r {
	case 'á', 'é', 'í', 'ó', 'ú':
		return true
	}
	return false
}

// hasCircumflex reports whether r is a vowel bearing the circumflex accent
// (closed quality).
func hasCircumflex(r rune) bool {
	switch r {
	case 'â', 'ê', 'ô':
		return true
	}
	return false
}

// hasTilde reports whether r is a vowel bearing the nasal tilde.
func hasTilde(r rune) bool {
	return r == 'ã' || r == 'õ'
}

// hasGrave reports whether r is a vowel bearing the grave accent. In pt-PT the
// grave accent occurs only in the crasis contraction 'à', which WordsPT does not
// generate in open-class words.
func hasGrave(r rune) bool {
	return r == 'à'
}

// hasDiaeresis reports whether r bears a diaeresis (trema). The Acordo
// Ortográfico abolished the diaeresis in pt-PT, so any such rune is illegal in a
// generated word.
func hasDiaeresis(r rune) bool {
	switch r {
	case 'ä', 'ë', 'ï', 'ö', 'ü':
		return true
	}
	return false
}

// firstRune returns the first rune of s, or 0 if s is empty.
func firstRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}

// lastVowelBase returns the base vowel of the last vowel rune in s and reports
// whether s contains any vowel. It is used to inspect a nucleus ending (for
// example the 'i' of the diphthong 'ei').
func lastVowelBase(s string) (rune, bool) {
	var base rune
	var found bool
	for _, r := range s {
		if b, ok := vowelBase(r); ok {
			base, found = b, true
		}
	}
	return base, found
}
