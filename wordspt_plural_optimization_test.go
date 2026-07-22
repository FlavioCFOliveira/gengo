package gengo

import (
	"math/rand/v2"
	"strings"
	"testing"
)

// This file guards the single-allocation additive-plural optimization (task #37).
// The noun and adjective cores now build an ADDITIVE plural (a regular +s, or +es
// for -r/-z/-n/oxytone -s) by appending the suffix into the SAME builder that
// assembles the word (accentedWordSuffixed), so the whole plural is one
// allocation, while SUBSTITUTIVE plurals (-ão, -m, vowel+l, and the fixed -ção)
// keep the post-assembly string transform of pluralize/pluralizeNoun (two
// allocations). These tests prove the optimization changes only HOW the plural is
// built, never WHICH string it produces, and that its allocation budget is exactly
// as documented.

// allPluralEndings returns every realized noun and adjective ending, paired with
// the pluralizer the core uses for it: pluralizeNoun for noun endings (so the
// fixed -ção outcome is honored) and the generic pluralize for adjective endings.
// It is the exhaustive set the optimization must preserve.
type pluralEndingCase struct {
	end       *nounEnding
	pluralize func(r *rand.Rand, singular string, end *nounEnding) string
}

func adjectivePluralizeRef(r *rand.Rand, singular string, end *nounEnding) string {
	return pluralize(r, singular, end.oxytone())
}

func allPluralEndings() []pluralEndingCase {
	nouns := []*nounEnding{
		&endMascO, &endFemA, &endMascEiro, &endFemEira, &endMascOr,
		&endFemOra, &endMascMento, &endFemCao, &endFemDade, &endFemAgem, &endComumIsta,
	}
	adj := []*nounEnding{
		&endAdjMascO, &endAdjFemA, &endAdjMascOso, &endAdjFemOsa,
		&endAdjMascIco, &endAdjFemIca, &endAdjMascIvo, &endAdjFemIva,
		&endAdjAl, &endAdjAvel, &endAdjIvel, &endAdjEnte, &endAdjAnte,
	}
	cases := make([]pluralEndingCase, 0, len(nouns)+len(adj))
	for _, e := range nouns {
		cases = append(cases, pluralEndingCase{e, pluralizeNoun})
	}
	for _, e := range adj {
		cases = append(cases, pluralEndingCase{e, adjectivePluralizeRef})
	}
	return cases
}

// TestAdditivePluralClassificationMatchesPluralize cross-checks the up-front
// classifier additivePluralSuffix against the ACTUAL behavior of the string-level
// pluralizer, independently of the classifier itself: a plural is additive if and
// only if the pluralizer left the singular as a prefix and merely appended to it,
// and in that case the appended bytes are exactly the classifier's suffix. It runs
// over every ending and many sampled stems, so a future ending misclassified as
// additive (or a wrong suffix) is caught before it can corrupt an output string.
func TestAdditivePluralClassificationMatchesPluralize(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5EED, 0xF00D))
	var buf []syllable
	for _, c := range allPluralEndings() {
		suffix, additive := additivePluralSuffix(c.end)
		for i := 0; i < 4000; i++ {
			leading := 1 + int(r.Uint32N(6))
			stem := sampleNounStem(r, buf, leading, c.end)
			buf = stem.syllables

			singular := accentedWord(stem)
			ref := c.pluralize(r, singular, c.end)

			// The pluralizer performed a pure suffix append exactly when the singular
			// survives as a prefix of its plural; that is the ground-truth definition of
			// "additive" the classifier must agree with.
			isAppend := strings.HasPrefix(ref, singular)
			if additive != isAppend {
				t.Fatalf("ending %s: additivePluralSuffix additive=%v but pluralize append=%v (singular %q -> plural %q)",
					c.end.label, additive, isAppend, singular, ref)
			}
			if additive {
				if got := ref[len(singular):]; got != suffix {
					t.Fatalf("ending %s: classifier suffix %q, pluralize appended %q (singular %q -> plural %q)",
						c.end.label, suffix, got, singular, ref)
				}
			}
		}
	}
}

// TestAccentedWordSuffixedMatchesPluralize proves the load-bearing equivalence:
// for every ADDITIVE ending, the one-allocation inline build
// accentedWordSuffixed(stem, suffix) is byte-identical to the previous two-step
// build pluralize(accentedWord(stem)). If this holds, the optimization cannot have
// changed any produced plural string.
func TestAccentedWordSuffixedMatchesPluralize(t *testing.T) {
	r := rand.New(rand.NewPCG(0xBEEF, 0xCAFE))
	var buf []syllable
	for _, c := range allPluralEndings() {
		suffix, additive := additivePluralSuffix(c.end)
		if !additive {
			continue // substitutive endings keep the unchanged pluralize path
		}
		for i := 0; i < 4000; i++ {
			leading := 1 + int(r.Uint32N(6))
			stem := sampleNounStem(r, buf, leading, c.end)
			buf = stem.syllables

			singular := accentedWord(stem)
			want := c.pluralize(r, singular, c.end) // the pre-optimization result
			got := accentedWordSuffixed(stem, suffix)
			if got != want {
				t.Fatalf("ending %s: accentedWordSuffixed = %q, want (pluralize) %q",
					c.end.label, got, want)
			}
		}
	}
}

// TestPluralBuildAllocationBudget pins the allocation cost of each build path: the
// additive inline build is exactly one allocation, while the substitutive
// assemble-then-transform build is exactly two. It measures fixed, pre-sampled
// stems so testing.AllocsPerRun sees only the build cost, not the sampler.
func TestPluralBuildAllocationBudget(t *testing.T) {
	r := rand.New(rand.NewPCG(0x1234, 0x5678))

	// Additive endings: +s (thematic -o) and +es (agentive -or) must each be one
	// allocation via the inline suffix build.
	stemS := sampleNounStem(r, nil, 3, &endMascO)
	if got := testing.AllocsPerRun(1000, func() { _ = accentedWordSuffixed(stemS, "s") }); got != 1 {
		t.Errorf("additive +s inline build: allocs/op = %v, want 1", got)
	}
	stemES := sampleNounStem(r, nil, 3, &endMascOr)
	if got := testing.AllocsPerRun(1000, func() { _ = accentedWordSuffixed(stemES, "es") }); got != 1 {
		t.Errorf("additive +es inline build: allocs/op = %v, want 1", got)
	}

	// Substitutive endings: the -m -> -ns and the fixed -ção -> -ções plurals keep
	// the assemble-then-transform build, which is one allocation for the assembly
	// plus one for the tail rewrite.
	stemAgem := sampleNounStem(r, nil, 3, &endFemAgem)
	if got := testing.AllocsPerRun(1000, func() { _ = pluralizeNoun(r, accentedWord(stemAgem), &endFemAgem) }); got != 2 {
		t.Errorf("substitutive -agens build: allocs/op = %v, want 2", got)
	}
	stemCao := sampleNounStem(r, nil, 3, &endFemCao)
	if got := testing.AllocsPerRun(1000, func() { _ = pluralizeNoun(r, accentedWord(stemCao), &endFemCao) }); got != 2 {
		t.Errorf("substitutive -ções build: allocs/op = %v, want 2", got)
	}

	// Substitutive vowel+l adjective ending (-al -> -ais) is likewise two.
	stemAl := sampleNounStem(r, nil, 3, &endAdjAl)
	if got := testing.AllocsPerRun(1000, func() { _ = pluralize(r, accentedWord(stemAl), endAdjAl.oxytone()) }); got != 2 {
		t.Errorf("substitutive -ais build: allocs/op = %v, want 2", got)
	}
}

// TestAdditivePluralSuffixClassification records the expected additive/substitutive
// verdict and suffix of every ending, so an accidental change to an ending's plural
// metadata (its final coda, nasal nucleus, or the fixed-ção flag) is caught
// directly at the classifier.
func TestAdditivePluralSuffixClassification(t *testing.T) {
	type want struct {
		suffix   string
		additive bool
	}
	cases := map[*nounEnding]want{
		// Nouns.
		&endMascO: {"s", true}, &endFemA: {"s", true},
		&endMascEiro: {"s", true}, &endFemEira: {"s", true},
		&endMascOr: {"es", true}, &endFemOra: {"s", true},
		&endMascMento: {"s", true}, &endComumIsta: {"s", true},
		&endFemDade: {"s", true},
		&endFemCao:  {"", false}, // -ção -> -ções
		&endFemAgem: {"", false}, // -agem -> -agens
		// Adjectives.
		&endAdjMascO: {"s", true}, &endAdjFemA: {"s", true},
		&endAdjMascOso: {"s", true}, &endAdjFemOsa: {"s", true},
		&endAdjMascIco: {"s", true}, &endAdjFemIca: {"s", true},
		&endAdjMascIvo: {"s", true}, &endAdjFemIva: {"s", true},
		&endAdjEnte: {"s", true}, &endAdjAnte: {"s", true},
		&endAdjAl:   {"", false}, // -al -> -ais
		&endAdjAvel: {"", false}, // -ável -> -áveis
		&endAdjIvel: {"", false}, // -ível -> -íveis
	}
	for end, w := range cases {
		suffix, additive := additivePluralSuffix(end)
		if suffix != w.suffix || additive != w.additive {
			t.Errorf("ending %s: additivePluralSuffix = (%q, %v), want (%q, %v)",
				end.label, suffix, additive, w.suffix, w.additive)
		}
	}
}
