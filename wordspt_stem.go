package gengo

import (
	"math/rand/v2"
	"strings"
)

// This file implements WordsPT Layer 1 generation: the correct-by-construction
// syllabic stem sampler and the single-allocation stem assembler. It builds on
// the static phonotactic data model in wordspt_data.go (the weighted
// inventories, positional constraints and rune helpers) and produces stems that
// always satisfy the conformance validators in wordspt.go (stemConformsToPhonotactics
// is the oracle). Accentuation (Layer 3) and morphology (Layer 2) are applied by
// later tasks; the stem produced here is the bare, unaccented syllable sequence.
//
// Correct-by-construction sampling
//
// Every draw is constrained BEFORE it is taken, so no candidate is ever rejected
// and there is no generate-and-reject loop. Work per call is deterministic:
// exactly three weighted draws per syllable (onset, nucleus, coda). Correctness
// maps one-to-one onto the branches of the validators:
//
//   - Word-initial onset (lh, nh, ç forbidden): the first syllable draws its
//     onset from onsetInitialInv, which omits those forms. Later syllables use
//     onsetMedialInv, where they are allowed.
//   - Onset/nucleus agreement (qu, gu need a front vowel; ç needs a back or
//     central vowel): after the onset is fixed, the nucleus is drawn from the
//     matching sub-inventory (nucleiFrontInv, nucleiBackInv, or the full nuclei).
//   - Inter-syllable transition (nasal m/n spelling and the Maximum Onset
//     Principle): the coda of each non-final syllable is drawn AFTER the next
//     syllable's onset is fixed, from a coda inventory conditioned on that onset,
//     mirroring transitionIsValid exactly. When the next syllable is
//     vowel-initial the coda is forced empty (MOP).
//
// The onsets and nuclei of a syllable depend only on that syllable's position,
// never on its neighbors, so they are sampled in a first pass; the codas, which
// depend on the following onset, are sampled in a second pass over the fixed
// onsets. This two-pass order is what removes every backtrack.
//
// Sampling weights for the derived inventories reuse the weights of the base
// inventories in wordspt_data.go (single source of truth). The only new weights
// are the empty-onset weights (word-initial and medial) and the empty-coda weight
// below; like the base weights, their exact values are ordinal estimates
// (flagged), while their ordering encodes authoritative pt-PT facts: the language
// strongly prefers consonant onsets and open (CV) syllables (Mateus, M. H. &
// d'Andrade, E., "The Phonology of Portuguese", Oxford, 2000; Cunha & Cintra), and
// it strongly avoids word-internal vowel hiatus (a coda-less syllable directly
// followed by a vowel-initial one). The values change only how often onsetless or
// closed syllables appear; they never affect the conformance validators.
//
// Hiatus and the two empty-onset weights
//
// A medial syllable with an absent (empty) onset is exactly what produces a
// word-internal vowel-vowel junction: because the Maximum Onset Principle forces
// the preceding syllable's coda empty before a vowel-initial onset (see the pass-2
// coda draw below), an empty medial onset always leaves the previous syllable
// vowel-final, so the two nuclei meet with no consonant between them. Runs of
// three or four vowels (a diphthong meeting a vowel-initial syllable) follow the
// same way. pt-PT largely avoids such medial hiatus, so the medial empty-onset
// weight is kept far below the word-initial one: word-initial vowel onsets stay
// natural (amor, ave, ilha), while medial vowel onsets — the unnatural junctions —
// become rare. A small residual is intentional: real pt-PT medial hiatus exists
// (país, saída, poesia), so the medial weight is small but non-zero, which also
// keeps the empty onset reachable in the medial inventory. Separating the two
// weights reduces hiatus while every draw stays a single, rejection-free weighted
// selection.

// emptyOnsetWeightInitial is the sampling weight of the absent (vowel-initial)
// onset at word-initial position, relative to the consonant-onset weights in
// wordspt_data.go. It is an ordinal estimate (flagged), kept below the consonant
// weights so that vowel-initial words stay a minority, matching the pt-PT
// preference for consonant onsets while allowing natural vowel-initial words
// (amor, ave, ilha).
const emptyOnsetWeightInitial = 90

// emptyOnsetWeightMedial is the sampling weight of the absent (vowel-initial)
// onset at a non-initial (medial) position. It is an ordinal estimate (flagged),
// kept far below emptyOnsetWeightInitial because a medial empty onset is precisely
// what creates the word-internal vowel hiatus that pt-PT avoids (see above). It is
// small but non-zero: a little real medial hiatus (país, saída) remains, and the
// empty onset stays a member of the medial onset inventory.
const emptyOnsetWeightMedial = 12

// emptyCodaWeight is the sampling weight of the absent (open-syllable) coda,
// relative to the coda weights in wordspt_data.go. It is an ordinal estimate
// (flagged), kept well above the coda weights so that most syllables are open
// (CV), matching the pt-PT preference for open syllables.
const emptyCodaWeight = 520

// ---------------------------------------------------------------------------
// Derived sampling inventories (built once, at package initialization)
// ---------------------------------------------------------------------------

// onsetInitialInv samples a word-initial onset: the empty (absent) onset or any
// legal onset except the ones forbidden word-initially (lh, nh, ç). It uses the
// larger word-initial empty-onset weight, because vowel-initial words are natural.
var onsetInitialInv = deriveInventory(emptyOnsetWeightInitial, onsetAllowedInitial, onsetSingles, onsetClusters)

// onsetMedialInv samples a non-initial onset: the empty onset or any legal onset,
// including lh, nh and ç, which are legal only between vowels. It uses the small
// medial empty-onset weight, so a medial vowel-initial syllable — the source of
// word-internal vowel hiatus — is rare, sharply reducing unnatural vowel runs.
var onsetMedialInv = deriveInventory(emptyOnsetWeightMedial, keepAllForms, onsetSingles, onsetClusters)

// nucleiFrontInv samples a nucleus that begins with a front vowel (e or i). It is
// the only nucleus set licensed after the qu and gu onsets.
var nucleiFrontInv = deriveInventory(0, nucleusStartsFront, nuclei)

// nucleiBackInv samples a nucleus that begins with a back or central vowel (a, o
// or u). It is the only nucleus set licensed after the ç onset.
var nucleiBackInv = deriveInventory(0, nucleusStartsBack, nuclei)

// codaBeforePBInv samples the coda of a syllable whose following onset begins
// with p or b: the empty coda or any coda except n (the nasal coda before p or b
// is spelled m). It mirrors the m/n branch of transitionIsValid.
var codaBeforePBInv = deriveInventory(emptyCodaWeight, codaAllowedBeforePB, codas)

// codaBeforeOtherInv samples the coda of a syllable whose following onset does
// not begin with p or b: the empty coda or any coda except m (a coda m occurs
// only before p or b). It mirrors the m/n branch of transitionIsValid.
var codaBeforeOtherInv = deriveInventory(emptyCodaWeight, codaAllowedBeforeOther, codas)

// codaFinalInv samples a word-final coda: the empty coda or any coda except n.
// The nasal coda at word end is spelled m in pt-PT (m before p, b and at word
// end; n otherwise), per the Acordo Ortográfico spelling convention.
var codaFinalInv = deriveInventory(emptyCodaWeight, codaAllowedFinal, codas)

// deriveInventory builds a sampling inventory from one or more base inventories,
// keeping only the forms for which keep returns true and, when emptyWeight is
// positive, prepending the empty form (an absent onset or coda) with that weight.
// It runs once, at package initialization, so every sampling weight keeps its
// single source of truth in the base inventories of wordspt_data.go.
func deriveInventory(emptyWeight uint32, keep func(form string) bool, bases ...weightedInventory) weightedInventory {
	forms := make([]weightedForm, 0, 8)
	if emptyWeight > 0 {
		forms = append(forms, weightedForm{form: "", weight: emptyWeight})
	}
	for i := range bases {
		for j := range bases[i].forms {
			f := bases[i].forms[j]
			if keep(f.form) {
				forms = append(forms, f)
			}
		}
	}
	return newWeightedInventory(forms)
}

// keepAllForms is the identity predicate for deriveInventory; it keeps every form.
func keepAllForms(string) bool { return true }

// onsetAllowedInitial reports whether an onset may begin a word: every onset
// except those in onsetInitialForbidden (lh, nh, ç).
func onsetAllowedInitial(form string) bool {
	_, forbidden := onsetInitialForbidden[form]
	return !forbidden
}

// codaAllowedBeforePB reports whether a coda may precede an onset that begins
// with p or b: any coda except n (the nasal coda there is spelled m).
func codaAllowedBeforePB(form string) bool { return form != "n" }

// codaAllowedBeforeOther reports whether a coda may precede an onset that does
// not begin with p or b: any coda except m (a coda m occurs only before p or b).
func codaAllowedBeforeOther(form string) bool { return form != "m" }

// codaAllowedFinal reports whether a coda may end a word: any coda except n,
// because the word-final nasal coda is spelled m in pt-PT.
func codaAllowedFinal(form string) bool { return form != "n" }

// ---------------------------------------------------------------------------
// Stem sampling
// ---------------------------------------------------------------------------

// syllabicStem is the internal output of Layer 1: an ordered sequence of
// syllables together with the zero-based index of its tonic (stressed) syllable.
// The tonic index is carried unchanged for Layer 3 (task #20), which uses it to
// place graphic accents deterministically; it does not influence the phonotactic
// sampling.
type syllabicStem struct {
	syllables []syllable
	tonic     int
}

// sampleForm draws one form from inv using the injectable source r. It performs a
// single bounded draw in [0,total) followed by the inventory's O(log n) cumulative
// lookup, so it never rejects a draw. inv must be non-empty (total > 0), which
// every inventory in this package guarantees.
func sampleForm(r *rand.Rand, inv *weightedInventory) string {
	return inv.forms[inv.sampleIndex(r.Uint32N(inv.total))].form
}

// sampleSyllabicStem builds a correct-by-construction pt-PT stem of syllableCount
// syllables, drawing from the injectable source r, and marks tonic as the
// stressed syllable. It writes into dst, reusing dst's backing array when it has
// enough capacity so that repeated calls allocate nothing; pass nil on the first
// call and the returned slice on later calls to keep per-call allocation at the
// single final string produced by assembleStem.
//
// syllableCount is honored exactly; a value below 1 is clamped to 1, because a
// stem has at least one syllable. tonic is clamped into [0,syllableCount) so the
// result always carries a valid tonic index. The returned stem always satisfies
// stemConformsToPhonotactics and assembles to valid, lowercase UTF-8.
func sampleSyllabicStem(r *rand.Rand, dst []syllable, syllableCount, tonic int) syllabicStem {
	if syllableCount < 1 {
		syllableCount = 1
	}
	if tonic < 0 {
		tonic = 0
	} else if tonic >= syllableCount {
		tonic = syllableCount - 1
	}

	if cap(dst) < syllableCount {
		dst = make([]syllable, syllableCount)
	} else {
		dst = dst[:syllableCount]
	}

	// Pass 1: onset and nucleus. Both depend only on the syllable's position,
	// never on a neighbor, so they can be fixed independently. The nucleus is
	// drawn from the sub-inventory that the chosen onset licenses, so the
	// onset/nucleus agreement holds by construction.
	for i := 0; i < syllableCount; i++ {
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

		dst[i].onset = onset
		dst[i].nucleus = sampleForm(r, nucInv)
	}

	// Pass 2: coda. Each non-final coda is fixed only after the next syllable's
	// onset is known, from a coda inventory conditioned on that onset, mirroring
	// transitionIsValid exactly; a vowel-initial next syllable forces an empty
	// coda (Maximum Onset Principle). The final coda is drawn from codaFinalInv.
	for i := 0; i < syllableCount; i++ {
		if i == syllableCount-1 {
			dst[i].coda = sampleForm(r, &codaFinalInv)
			continue
		}
		nextOnset := dst[i+1].onset
		switch {
		case nextOnset == "":
			dst[i].coda = ""
		case nextOnset[0] == 'p' || nextOnset[0] == 'b':
			dst[i].coda = sampleForm(r, &codaBeforePBInv)
		default:
			dst[i].coda = sampleForm(r, &codaBeforeOtherInv)
		}
	}

	return syllabicStem{syllables: dst, tonic: tonic}
}

// ---------------------------------------------------------------------------
// Stem assembly
// ---------------------------------------------------------------------------

// assembleStem renders a syllable sequence to a single lowercase pt-PT string in
// exactly one allocation. It pre-sizes a strings.Builder to the exact byte length
// of the concatenated syllable parts, so the writes never reallocate and
// Builder.String returns the assembled bytes without a second copy. Every part is
// already lowercase and valid UTF-8, so the result is lowercase, valid UTF-8.
func assembleStem(syllables []syllable) string {
	total := 0
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
	return b.String()
}
