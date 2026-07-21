package gengo

import (
	"math/rand/v2"
	"testing"
	"unicode/utf8"
)

// The tests in this file exercise the pt-PT plural-formation engine
// (wordspt_plural.go). They cover the rule table exhaustively with real pt-PT
// words, verify the three -ão outcomes and their weighting, prove that a
// pluralized string keeps its Sprint 7 string-level orthographic conformance over
// a large sample, and confirm reproducibility and the allocation budget.

// TestPluralizeRuleTable covers the plural-rule dispatch exhaustively with real
// pt-PT singular/plural pairs, one (or more) per rule branch. The -ão branch is
// covered separately by TestApplyAoPlural because its outcome is weighted.
func TestPluralizeRuleTable(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	tests := []struct {
		name     string
		singular string
		oxytone  bool
		want     string
		rule     string
	}{
		// Vowel or diphthong (oral/nasal) -> +s.
		{"casa", "casa", false, "casas", "vowel +s"},
		{"café", "café", true, "cafés", "vowel +s"},
		{"maçã", "maçã", true, "maçãs", "nasal vowel +s"},
		{"mãe", "mãe", true, "mães", "nasal diphthong +s"},
		{"herói", "herói", true, "heróis", "oral diphthong +s"},
		{"gato", "gato", false, "gatos", "vowel +s"},
		// -m -> -ns.
		{"homem", "homem", false, "homens", "-m -> -ns"},
		{"jardim", "jardim", true, "jardins", "-m -> -ns"},
		{"bom", "bom", true, "bons", "-m -> -ns"},
		// -r/-z/-n -> +es.
		{"flor", "flor", true, "flores", "-r +es"},
		{"amor", "amor", true, "amores", "-r +es"},
		{"luz", "luz", true, "luzes", "-z +es"},
		{"rapaz", "rapaz", true, "rapazes", "-z +es"},
		{"liquen", "líquen", false, "líquenes", "-n +es"},
		// -x invariable.
		{"tórax", "tórax", false, "tórax", "-x invariable"},
		{"látex", "látex", false, "látex", "-x invariable"},
		// -s: oxytone +es, paroxytone invariable.
		{"país", "país", true, "países", "-s oxytone +es"},
		{"lápis", "lápis", false, "lápis", "-s paroxytone invariable"},
		{"vírus", "vírus", false, "vírus", "-s paroxytone invariable"},
		// -al/-el/-ol/-ul -> -ais/-éis/-óis/-uis (all oxytone).
		{"animal", "animal", true, "animais", "-al -> -ais"},
		{"papel", "papel", true, "papéis", "-el oxytone -> -éis"},
		{"anzol", "anzol", true, "anzóis", "-ol oxytone -> -óis"},
		{"azul", "azul", true, "azuis", "-ul -> -uis"},
		// -el/-ável/-ível unstressed -> -eis (no added acute).
		{"amável", "amável", false, "amáveis", "-el paroxytone -> -eis"},
		{"possível", "possível", false, "possíveis", "-el paroxytone -> -eis"},
		// -il: oxytone -> -is, paroxytone -> -eis.
		{"funil", "funil", true, "funis", "-il oxytone -> -is"},
		{"barril", "barril", true, "barris", "-il oxytone -> -is"},
		{"fácil", "fácil", false, "fáceis", "-il paroxytone -> -eis"},
		{"útil", "útil", false, "úteis", "-il paroxytone -> -eis"},
		{"difícil", "difícil", false, "difíceis", "-il paroxytone -> -eis"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := pluralize(r, tc.singular, tc.oxytone)
			if got != tc.want {
				t.Errorf("pluralize(%q, oxytone=%v) = %q, want %q [%s]",
					tc.singular, tc.oxytone, got, tc.want, tc.rule)
			}
		})
	}
}

// TestApplyAoPlural verifies each of the three -ão plural outcomes against its
// canonical pt-PT example. These are the pure transforms the weighted selector
// dispatches to, so the rule table is covered exhaustively including the -ão row.
func TestApplyAoPlural(t *testing.T) {
	tests := []struct {
		singular string
		variant  aoPlural
		want     string
	}{
		{"leão", aoOes, "leões"},
		{"coração", aoOes, "corações"},
		{"pão", aoAes, "pães"},
		{"cão", aoAes, "cães"},
		{"mão", aoAos, "mãos"},
		{"irmão", aoAos, "irmãos"},
	}
	for _, tc := range tests {
		if got := applyAoPlural(tc.singular, tc.variant); got != tc.want {
			t.Errorf("applyAoPlural(%q, %d) = %q, want %q", tc.singular, tc.variant, got, tc.want)
		}
	}
}

// TestPluralizeAoIsWeightedAndReachable proves that the top-level dispatch of an
// -ão word yields only the three valid outcomes, that every outcome is reachable,
// and that -ões is the dominant (productive default) outcome.
func TestPluralizeAoIsWeightedAndReachable(t *testing.T) {
	r := rand.New(rand.NewPCG(0xA0, 0xB0))
	counts := map[string]int{}
	const n = 60000
	valid := map[string]bool{"corações": true, "coraçães": true, "coraçãos": true}
	for i := 0; i < n; i++ {
		got := pluralize(r, "coração", true)
		if !valid[got] {
			t.Fatalf("pluralize(%q) produced unexpected form %q", "coração", got)
		}
		counts[got]++
	}
	for form := range valid {
		if counts[form] == 0 {
			t.Errorf("-ão outcome %q was never produced over %d draws", form, n)
		}
	}
	if counts["corações"] <= counts["coraçães"] || counts["corações"] <= counts["coraçãos"] {
		t.Errorf("-ões should dominate: ões=%d ães=%d ãos=%d",
			counts["corações"], counts["coraçães"], counts["coraçãos"])
	}
}

// TestSampleAoPluralDistribution checks the weighted selector directly: all three
// outcomes appear and their observed frequencies follow the documented weights
// (-ões dominant), within a generous tolerance.
func TestSampleAoPluralDistribution(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 11))
	counts := [3]int{}
	const n = 300000
	for i := 0; i < n; i++ {
		counts[sampleAoPlural(r)]++
	}
	for v, c := range counts {
		if c == 0 {
			t.Fatalf("aoPlural outcome %d never sampled", v)
		}
	}
	if counts[aoOes] <= counts[aoAes] || counts[aoOes] <= counts[aoAos] {
		t.Errorf("-ões must dominate: got ões=%d ães=%d ãos=%d", counts[aoOes], counts[aoAes], counts[aoAos])
	}
	// Observed -ões fraction should be near the configured 78% (±3%).
	frac := float64(counts[aoOes]) / float64(n)
	if frac < 0.75 || frac > 0.81 {
		t.Errorf("-ões fraction = %.3f, want ~0.78", frac)
	}
}

// TestPluralizeEmpty confirms the empty-string contract.
func TestPluralizeEmpty(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	if got := pluralize(r, "", false); got != "" {
		t.Errorf("pluralize(\"\") = %q, want empty", got)
	}
}

// TestPluralizePreservesStringOrthography proves that pluralizing a generated,
// accented pt-PT word preserves the Sprint 7 STRING-LEVEL orthographic oracles
// (valid lowercase UTF-8, cedilla rule, no abolished diaeresis, no grave accent)
// over a large sample across syllable counts and stress positions. Per the engine
// design, the syllabic-accentuation validator is NOT applied to the transformed
// plural; the string oracles are the correct oracle for a string transform.
func TestPluralizePreservesStringOrthography(t *testing.T) {
	const loop = 2000
	r := rand.New(rand.NewPCG(0xC0FFEE, 0xBEEF))
	var dst []syllable
	for i := 0; i < loop; i++ {
		for count := 1; count <= 6; count++ {
			tonic := i % count
			st := sampleSyllabicStem(r, dst, count, tonic)
			dst = st.syllables
			singular := accentedWord(st)
			oxytone := tonic == count-1

			plural := pluralize(r, singular, oxytone)

			if plural == "" {
				t.Fatalf("plural of %q is empty", singular)
			}
			if !utf8.ValidString(plural) {
				t.Fatalf("plural %q of %q is not valid UTF-8", plural, singular)
			}
			if hasUpper(plural) {
				t.Fatalf("plural %q of %q contains an uppercase letter", plural, singular)
			}
			if !wordConformsToCedilla(plural) {
				t.Fatalf("plural %q of %q violates the cedilla rule", plural, singular)
			}
			if wordHasForbiddenDiaeresis(plural) {
				t.Fatalf("plural %q of %q contains a forbidden diaeresis", plural, singular)
			}
			if wordContainsGrave(plural) {
				t.Fatalf("plural %q of %q contains a grave accent", plural, singular)
			}
		}
	}
}

// TestPluralizeReproducible confirms that two identically seeded sources produce
// byte-for-byte identical plurals, so the *Generator surface (own seeded source)
// is reproducible through the plural engine, including its weighted -ão branch.
func TestPluralizeReproducible(t *testing.T) {
	a := New(20260721)
	b := New(20260721)
	var da, db []syllable
	for i := 0; i < 1000; i++ {
		count := 1 + i%6
		tonic := i % count
		sa := sampleSyllabicStem(a.r, da, count, tonic)
		sb := sampleSyllabicStem(b.r, db, count, tonic)
		da, db = sa.syllables, sb.syllables
		wa := pluralize(a.r, accentedWord(sa), tonic == count-1)
		wb := pluralize(b.r, accentedWord(sb), tonic == count-1)
		if wa != wb {
			t.Fatalf("identical seeds diverged at i=%d: %q != %q", i, wa, wb)
		}
	}
}

// TestPluralizeAllocations documents the allocation budget of the engine: a
// transforming ending performs exactly one allocation (the new string) and an
// invariable ending performs none (it returns the input unchanged).
func TestPluralizeAllocations(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	if got := testing.AllocsPerRun(1000, func() { _ = pluralize(r, "casa", false) }); got != 1 {
		t.Errorf("pluralize transforming allocs/op = %v, want 1", got)
	}
	if got := testing.AllocsPerRun(1000, func() { _ = pluralize(r, "tórax", false) }); got != 0 {
		t.Errorf("pluralize invariable (-x) allocs/op = %v, want 0", got)
	}
	if got := testing.AllocsPerRun(1000, func() { _ = pluralize(r, "lápis", false) }); got != 0 {
		t.Errorf("pluralize invariable (-s paroxytone) allocs/op = %v, want 0", got)
	}
}

// BenchmarkPluralize measures a transforming plural (the +s branch), which is one
// allocation per call.
func BenchmarkPluralize(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pluralize(r, "casa", false)
	}
}

// BenchmarkPluralizeL measures the stress-conditioned -l branch (papel -> papéis).
func BenchmarkPluralizeL(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pluralize(r, "papel", true)
	}
}

// BenchmarkPluralizeAo measures the weighted -ão branch (a draw plus one
// allocation).
func BenchmarkPluralizeAo(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pluralize(r, "coração", true)
	}
}

// BenchmarkPluralizeInvariable measures an invariable plural (-x), which performs
// no allocation.
func BenchmarkPluralizeInvariable(b *testing.B) {
	r := rand.New(rand.NewPCG(1, 2))
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = pluralize(r, "tórax", false)
	}
}
