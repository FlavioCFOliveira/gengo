package gengo

import (
	"math/rand/v2"
	"testing"
	"unicode/utf8"
)

// The tests in this file are the acceptance-criteria regression guards for the
// stem-naturalness refinement (Sprint 8): the reduced word-internal vowel hiatus
// (wordspt_stem.go, the small medial empty-onset weight) and the tightened,
// downward-skewed effective word length (wordspt_inflection.go). They assert
// absolute bounds the refined sampler meets by a wide margin and that the previous
// (Sprint 7) sampler could not, so a revert of either refinement fails them.
//
// Measured before/after (200000-word samples), recorded here so the thresholds
// are auditable:
//
//	metric                              before (Sprint 7)   after (this task)
//	NounPT words with a >=3-vowel run   22.47%              3.92%
//	NounPT mean of per-word max run     2.093               1.579
//	NounPT max vowel run observed       9                   5
//	bare-stem words with a >=3-vowel run  (proportional)    3.48%
//	Big/Singular mean word length       19.16 chars         12.95 chars
//	Big/Singular words <= 14 chars      ~27%                ~74%

// isVowel reports whether r is a pt-PT vowel rune (any accented form).
func isVowel(r rune) bool {
	_, ok := vowelBase(r)
	return ok
}

// longestVowelRun returns the length of the longest run of consecutive vowel
// runes in s. A word-internal vowel hiatus (a coda-less syllable followed by a
// vowel-initial one) shows up here as a run of three or more vowels once a
// diphthong meets the junction; a value of two or less is the natural case (a
// single vowel, or one pt-PT diphthong).
func longestVowelRun(s string) int {
	best, cur := 0, 0
	for _, r := range s {
		if isVowel(r) {
			cur++
			if cur > best {
				best = cur
			}
		} else {
			cur = 0
		}
	}
	return best
}

// TestStemHiatusReduced asserts, over a large sample of bare stems drawn straight
// from the refined sampler, that the word-internal vowel hiatus is sharply
// reduced: few stems carry a run of three or more vowels and the mean longest run
// stays well below the Sprint 7 level. Bare stems isolate the sampler itself,
// without the noun endings, so this is the tightest guard on the onset
// conditioning.
func TestStemHiatusReduced(t *testing.T) {
	r := rand.New(rand.NewPCG(0xC0FFEE, 0xBEEF))
	var dst []syllable

	const loop = 2000 // 2000 * 8 counts = 16000 stems, well above the >=1000 floor
	var samples, sumMax, runsGE3 int
	for i := 0; i < loop; i++ {
		for count := 1; count <= 8; count++ {
			st := sampleSyllabicStem(r, dst, count, i%count)
			dst = st.syllables
			run := longestVowelRun(accentedWord(st))
			samples++
			sumMax += run
			if run >= 3 {
				runsGE3++
			}
		}
	}

	meanMax := float64(sumMax) / float64(samples)
	fracGE3 := float64(runsGE3) / float64(samples)

	// The Sprint 7 sampler produced a >=3-vowel run in a large minority of words
	// (proportional to the 22% seen on nouns) and a mean longest run above two; the
	// refined sampler stays far below both, with generous margins against seed noise.
	if fracGE3 >= 0.08 {
		t.Errorf("bare-stem >=3-vowel runs = %.3f%% of %d samples, want < 8%% (Sprint 7 was far higher)",
			100*fracGE3, samples)
	}
	if meanMax >= 1.9 {
		t.Errorf("bare-stem mean longest vowel run = %.4f, want < 1.9 (Sprint 7 was ~2.0+)", meanMax)
	}
	t.Logf("bare-stem hiatus over %d samples: mean longest run=%.4f, >=3-run fraction=%.3f%%",
		samples, meanMax, 100*fracGE3)
}

// TestNounPTHiatusReduced asserts the same hiatus reduction on the public NounPT
// surface (random gender, number and length), the words a user actually sees.
func TestNounPTHiatusReduced(t *testing.T) {
	g := New(123456789)

	const n = 40000
	var sumMax, runsGE3, maxObserved int
	for i := 0; i < n; i++ {
		run := longestVowelRun(g.NounPT())
		sumMax += run
		if run >= 3 {
			runsGE3++
		}
		if run > maxObserved {
			maxObserved = run
		}
	}

	meanMax := float64(sumMax) / float64(n)
	fracGE3 := float64(runsGE3) / float64(n)

	// Sprint 7: 22.47% of words had a >=3-vowel run, mean 2.093, max 9.
	if fracGE3 >= 0.08 {
		t.Errorf("NounPT >=3-vowel runs = %.3f%% of %d words, want < 8%% (Sprint 7 was 22.47%%)",
			100*fracGE3, n)
	}
	if meanMax >= 1.8 {
		t.Errorf("NounPT mean longest vowel run = %.4f, want < 1.8 (Sprint 7 was 2.093)", meanMax)
	}
	t.Logf("NounPT hiatus over %d words: mean longest run=%.4f, >=3-run fraction=%.3f%%, max run=%d",
		n, meanMax, 100*fracGE3, maxObserved)
}

// TestBigNounLengthRealistic asserts that the Big category no longer yields the
// implausible ~19-character, uniformly-spread words of Sprint 7: the mean word
// length now lands in the realistic ~12-13 character range, and the clear
// majority of Big words fall in the shorter part of the window, all while every
// word stays strictly inside the 9..30 window (the strict window is also enforced
// by TestNounPTLengthHonored).
func TestBigNounLengthRealistic(t *testing.T) {
	g := New(0xB19B19)

	const n = 40000
	var sum, shortHalf int
	for i := 0; i < n; i++ {
		w := g.NounPTOf(AnyGender, Singular, BigLengthWords)
		l := utf8.RuneCountInString(w)
		if l < 9 || l > 30 {
			t.Fatalf("Big noun %q has %d runes, outside the 9..30 window", w, l)
		}
		sum += l
		if l <= 14 {
			shortHalf++
		}
	}

	mean := float64(sum) / float64(n)
	shortFrac := float64(shortHalf) / float64(n)

	// A uniform draw over 9..30 has mean 19.5 (Sprint 7 measured 19.16); the skew
	// must pull the mean down into the realistic range without collapsing to the
	// floor.
	if mean < 9.5 || mean > 15.0 {
		t.Errorf("Big/Singular mean length = %.3f, want in [9.5,15.0] (Sprint 7 was 19.16)", mean)
	}
	// Under the downward skew most Big words are short; a uniform draw would put
	// only ~27%% at or below 14 characters.
	if shortFrac <= 0.5 {
		t.Errorf("Big/Singular words <= 14 chars = %.2f%%, want > 50%% (Sprint 7 was ~27%%)",
			100*shortFrac)
	}
	t.Logf("Big/Singular over %d words: mean length=%.3f, fraction <=14 chars=%.2f%%",
		n, mean, 100*shortFrac)
}

// TestMedialOnsetPrefersConsonant is a structural guard on the mechanism: the
// medial onset inventory must still offer the empty (vowel-initial) onset (so a
// little real hiatus like país remains and the sampler stays rejection-free), but
// with a weight far below the word-initial one, so medial vowel onsets are rare.
func TestMedialOnsetPrefersConsonant(t *testing.T) {
	if !inventoryContains(onsetMedialInv, "") {
		t.Fatal("onsetMedialInv must still offer the empty onset (some real hiatus remains)")
	}
	if emptyOnsetWeightMedial >= emptyOnsetWeightInitial {
		t.Errorf("medial empty-onset weight %d must be far below the initial weight %d",
			emptyOnsetWeightMedial, emptyOnsetWeightInitial)
	}

	// The empty onset's share of the medial inventory must be a small minority,
	// which is what makes medial hiatus rare.
	var emptyWeight, total uint32
	for i := range onsetMedialInv.forms {
		total += onsetMedialInv.forms[i].weight
		if onsetMedialInv.forms[i].form == "" {
			emptyWeight = onsetMedialInv.forms[i].weight
		}
	}
	if share := float64(emptyWeight) / float64(total); share >= 0.05 {
		t.Errorf("medial empty-onset share = %.3f, want < 0.05 (a small residual hiatus)", share)
	}
}
