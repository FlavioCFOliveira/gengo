package gengo

import (
	"math/rand/v2"
	"strconv"
	"testing"
	"unicode/utf8"
)

// The tests in this file exercise the WordsPT Layer 3 accentuation application
// and the single-allocation accented assembler (wordspt_accent.go). They prove,
// over a large sampled run, that every accented stem is correct by construction
// against the accentuation oracle in wordspt.go (stemConformsToAccentuation) and
// the string-level orthographic oracles (cedilla, diaeresis, grave); that the
// output is always valid lowercase UTF-8 with exactly one graphic accent where
// the rules demand one; that the accent is fully deterministic (no rejection,
// no random draws); and that the accented word is produced in a single
// allocation, preserving the Layer 1 allocation budget.

// applyAccentPlan returns a copy of src with plan applied to the tonic nucleus,
// so the accented syllables can be fed to the [stemConformsToAccentuation]
// oracle, which inspects the nucleus graphemes. It is a test-only helper and
// allocates freely; production assembly substitutes in place without allocating.
func applyAccentPlan(src []syllable, plan accentPlan) []syllable {
	out := make([]syllable, len(src))
	copy(out, src)
	if plan.tonic >= 0 {
		nucleus := src[plan.tonic].nucleus
		_, size := utf8.DecodeRuneInString(nucleus)
		out[plan.tonic].nucleus = string(plan.accented) + nucleus[size:]
	}
	return out
}

// countGraphicAccents returns how many runes of s bear an acute or circumflex
// accent. WordsPT places exactly one such accent per word that requires it, and
// none otherwise.
func countGraphicAccents(s string) int {
	n := 0
	for _, r := range s {
		if hasAcute(r) || hasCircumflex(r) {
			n++
		}
	}
	return n
}

// TestAccentuateStemConformsToAccentuation is the acceptance-criterion test:
// over at least 1000 iterations, across a range of syllable counts and tonic
// positions, every generated stem is accentuated, and the accented result
// satisfies the accentuation oracle and every string-level orthographic oracle.
// It also verifies that the assembled string carries exactly the graphic accent
// the plan describes (one when required, none otherwise), tying the produced
// string to the deterministic plan.
func TestAccentuateStemConformsToAccentuation(t *testing.T) {
	const loop = 1000
	r := newTestRand()

	var dst []syllable
	for i := 0; i < loop; i++ {
		for count := 1; count <= 8; count++ {
			tonic := i % count // exercise every tonic position over the run

			st := sampleSyllabicStem(r, dst, count, tonic)
			dst = st.syllables // reuse the buffer across iterations

			plan := accentuateStem(st)

			// The accented syllables must satisfy the accentuation oracle.
			accented := applyAccentPlan(st.syllables, plan)
			if !stemConformsToAccentuation(accented, st.tonic) {
				t.Fatalf("accented stem violates accentuation oracle: tonic=%d %+v",
					st.tonic, accented)
			}

			// A nasal tonic must never gain an acute or circumflex accent.
			if nucleusIsNasal(st.syllables[st.tonic].nucleus) && plan.tonic >= 0 {
				t.Fatalf("nasal tonic nucleus %q was given a graphic accent",
					st.syllables[st.tonic].nucleus)
			}

			word := accentedWord(st)

			if !utf8.ValidString(word) {
				t.Fatalf("accented word is not valid UTF-8: %q", word)
			}
			if word == "" {
				t.Fatalf("accented word is empty for count %d", count)
			}
			if hasUpper(word) {
				t.Fatalf("accented word contains an uppercase letter: %q", word)
			}

			// Exactly one graphic accent iff the rules require one.
			wantAccents := 0
			if plan.tonic >= 0 {
				wantAccents = 1
			}
			if got := countGraphicAccents(word); got != wantAccents {
				t.Fatalf("accented word %q has %d graphic accents, want %d (tonic=%d)",
					word, got, wantAccents, st.tonic)
			}

			// The string-level orthographic oracles must hold after accentuation.
			if !wordConformsToCedilla(word) {
				t.Fatalf("accented word violates the cedilla rule: %q", word)
			}
			if wordHasForbiddenDiaeresis(word) {
				t.Fatalf("accented word contains a forbidden diaeresis: %q", word)
			}
			if wordContainsGrave(word) {
				t.Fatalf("accented word contains a grave accent: %q", word)
			}
		}
	}
}

// TestAccentuateStemKnownCases pins each stress-class rule and the documented
// acute/circumflex choice to concrete, hand-built stems with a known correct
// pt-PT spelling. Each stem is phonotactically valid, so the cases are honest.
func TestAccentuateStemKnownCases(t *testing.T) {
	tests := []struct {
		name  string
		syls  []syllable
		tonic int
		want  string
	}{
		// Oxytones accented on final tonic a/e/o (optionally + s), and -em.
		{"oxytone-a", []syllable{syl("s", "o", ""), syl("f", "a", "")}, 1, "sofá"},
		{"oxytone-e", []syllable{syl("c", "a", ""), syl("f", "e", "")}, 1, "café"},
		{"oxytone-o", []syllable{syl("", "a", ""), syl("v", "o", "")}, 1, "avó"},
		{"oxytone-as", []syllable{syl("", "a", ""), syl("tr", "a", "s")}, 1, "atrás"},
		{"oxytone-em-acute", []syllable{syl("", "a", ""), syl("l", "e", "m")}, 1, "além"},
		// Oxytone endings that take NO accent.
		{"oxytone-im-none", []syllable{syl("", "a", "s"), syl("s", "i", "m")}, 1, "assim"},
		{"oxytone-u-none", []syllable{syl("", "a", ""), syl("z", "u", "l")}, 1, "azul"},
		// Proparoxytones always accented; â only for a before a nasal.
		{"proparoxytone-a", []syllable{syl("s", "a", ""), syl("b", "a", ""), syl("d", "o", "")}, 0, "sábado"},
		{"proparoxytone-a-circumflex-nasal-onset", []syllable{syl("c", "a", ""), syl("m", "a", ""), syl("r", "a", "")}, 0, "câmara"},
		{"proparoxytone-a-circumflex-nasal-coda", []syllable{syl("", "a", "n"), syl("c", "o", ""), syl("r", "a", "")}, 0, "âncora"},
		{"proparoxytone-e-stays-acute-before-nasal", []syllable{syl("g", "e", ""), syl("n", "e", ""), syl("r", "o", "")}, 0, "género"},
		// Paroxytones: accented on non-default endings, bare on default endings.
		{"paroxytone-r-nondefault", []syllable{syl("c", "a", ""), syl("d", "a", ""), syl("v", "e", "r")}, 1, "cadáver"},
		{"paroxytone-diphthong-nondefault", []syllable{syl("j", "o", ""), syl("qu", "ei", "")}, 0, "jóquei"},
		{"paroxytone-default-a-none", []syllable{syl("c", "a", ""), syl("s", "a", "")}, 0, "casa"},
		{"paroxytone-default-as-none", []syllable{syl("c", "a", ""), syl("s", "a", "s")}, 0, "casas"},
		{"paroxytone-nasal-ending-acute-tonic", []syllable{syl("", "o", "r"), syl("g", "ão", "")}, 0, "órgão"},
		// Nasal tonic nuclei carry the tilde only, never an added accent.
		{"nasal-vowel-tonic", []syllable{syl("", "i", "r"), syl("m", "ã", "")}, 1, "irmã"},
		{"nasal-diphthong-tonic-cedilla", []syllable{syl("c", "o", ""), syl("r", "a", ""), syl("ç", "ão", "")}, 2, "coração"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if !stemConformsToPhonotactics(tc.syls) {
				t.Fatalf("test stem is not phonotactically valid: %+v", tc.syls)
			}
			st := syllabicStem{syllables: tc.syls, tonic: tc.tonic}
			if got := accentedWord(st); got != tc.want {
				t.Errorf("accentedWord = %q, want %q", got, tc.want)
			}
			// The produced spelling must also satisfy the accentuation oracle.
			plan := accentuateStem(st)
			if !stemConformsToAccentuation(applyAccentPlan(tc.syls, plan), tc.tonic) {
				t.Errorf("known case %q fails the accentuation oracle", tc.want)
			}
		})
	}
}

// TestAccentuateStemDeterministic confirms the accent is a pure function of the
// stem: repeated calls yield an identical plan and word, with no dependence on
// any random source. This is the deterministic, correct-by-construction property
// the task requires.
func TestAccentuateStemDeterministic(t *testing.T) {
	r := newTestRand()
	var dst []syllable
	for i := 0; i < 1000; i++ {
		count := 1 + i%6
		st := sampleSyllabicStem(r, dst, count, i%count)
		dst = st.syllables

		p1 := accentuateStem(st)
		p2 := accentuateStem(st)
		if p1 != p2 {
			t.Fatalf("accentuateStem not deterministic: %+v != %+v", p1, p2)
		}
		if w1, w2 := accentedWord(st), accentedWord(st); w1 != w2 {
			t.Fatalf("accentedWord not deterministic: %q != %q", w1, w2)
		}
	}
}

// TestAccentedWordNoRandomDraws proves the full accented-word path consumes
// exactly three random draws per syllable — the sampler's fixed budget — which
// shows the accentuation stage adds no draw and performs no rejection loop.
func TestAccentedWordNoRandomDraws(t *testing.T) {
	for count := 1; count <= 8; count++ {
		got := countUint32Draws(func(r *rand.Rand) {
			_ = accentedWord(sampleSyllabicStem(r, nil, count, 0))
		})
		want := uint64(3 * count)
		if got != want {
			t.Errorf("count=%d: consumed %d draws, want %d (accentuation adds none)", count, got, want)
		}
	}
}

// TestAssembleAccentedStemSingleAllocation asserts the accented assembler
// performs exactly one allocation for the returned string, for both an accented
// and an unaccented plan.
func TestAssembleAccentedStemSingleAllocation(t *testing.T) {
	r := newTestRand()

	// A stem that requires an accent (forced proparoxytone).
	st := sampleSyllabicStem(r, nil, 3, 0)
	plan := accentuateStem(st)
	if plan.tonic < 0 {
		t.Fatalf("expected an accented plan for a proparoxytone, got none")
	}
	if got := testing.AllocsPerRun(1000, func() { _ = assembleAccentedStem(st.syllables, plan) }); got != 1 {
		t.Errorf("assembleAccentedStem (accented) allocations/op = %v, want 1", got)
	}

	// An explicitly unaccented plan must also assemble in one allocation.
	none := accentPlan{tonic: -1}
	if got := testing.AllocsPerRun(1000, func() { _ = assembleAccentedStem(st.syllables, none) }); got != 1 {
		t.Errorf("assembleAccentedStem (unaccented) allocations/op = %v, want 1", got)
	}
}

// TestAccentedWordAllocationBudget demonstrates the end-to-end Layer 3 budget:
// with a reused syllable buffer, sampling and accentuation allocate nothing, and
// the accented string is the single allocation per generated word.
func TestAccentedWordAllocationBudget(t *testing.T) {
	r := newTestRand()
	dst := make([]syllable, 0, 16) // pre-sized reusable buffer
	got := testing.AllocsPerRun(1000, func() {
		st := sampleSyllabicStem(r, dst, 4, 1)
		dst = st.syllables
		_ = accentedWord(st)
	})
	if got != 1 {
		t.Errorf("full accented-word allocations/op = %v, want 1 (the accented string)", got)
	}
}

// TestAccentuateStemEmptyIsSafe verifies the defensive guards: an empty or
// out-of-range stem yields a no-accent plan rather than panicking.
func TestAccentuateStemEmptyIsSafe(t *testing.T) {
	if p := accentuateStem(syllabicStem{syllables: nil, tonic: 0}); p.tonic != -1 {
		t.Errorf("empty stem: got plan %+v, want no accent", p)
	}
	st := syllabicStem{syllables: []syllable{syl("c", "a", "")}, tonic: 5}
	if p := accentuateStem(st); p.tonic != -1 {
		t.Errorf("out-of-range tonic: got plan %+v, want no accent", p)
	}
}

// BenchmarkAccentuateStem measures deriving the accent plan for a fixed stem; it
// allocates nothing (the plan is a value) and draws no randomness.
func BenchmarkAccentuateStem(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	st := sampleSyllabicStem(r, nil, 4, 0)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = accentuateStem(st)
	}
}

// BenchmarkAssembleAccentedStem measures assembly of a fixed accented stem; it
// performs exactly one allocation per operation (the returned string).
func BenchmarkAssembleAccentedStem(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	st := sampleSyllabicStem(r, nil, 4, 0)
	plan := accentuateStem(st)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = assembleAccentedStem(st.syllables, plan)
	}
}

// BenchmarkBuildAccentedWord measures the full Layer 3 path (sample into a reused
// buffer, accentuate, assemble); it performs exactly one allocation per
// operation: the accented string.
func BenchmarkBuildAccentedWord(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	var dst []syllable
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		st := sampleSyllabicStem(r, dst, 4, 1)
		dst = st.syllables
		_ = accentedWord(st)
	}
}

// BenchmarkBuildAccentedWordByCount measures the full accented path across
// syllable counts, so the per-call work is visibly linear in the syllable count
// (deterministic, no rejection).
func BenchmarkBuildAccentedWordByCount(b *testing.B) {
	for _, count := range []int{1, 2, 3, 4, 6, 8} {
		b.Run("syllables="+strconv.Itoa(count), func(b *testing.B) {
			r := rand.New(rand.NewPCG(1, 2))
			var dst []syllable
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				st := sampleSyllabicStem(r, dst, count, 0)
				dst = st.syllables
				_ = accentedWord(st)
			}
		})
	}
}
