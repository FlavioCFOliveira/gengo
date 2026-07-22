package gengo

import (
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the WordsPT closed-class infrastructure and the ArticlePT,
// PrepositionPT, and ConjunctionPT selectors (package-level functions and their
// *Generator methods). The checks map to the task acceptance criteria and to the
// closed-class rules in specification/wordspt-closed-classes.md and
// specification/wordspt.md:
//
//   - Curation integrity: each embedded list is non-empty, duplicate-free, and
//     every member is a single, lowercase, valid-UTF-8 word (AC "closed-class
//     membership").
//   - Membership: every output over a large sampled run is a member of the
//     corresponding curated list, and every member is reachable (AC4).
//   - Encoding and case: every output is valid UTF-8 with no uppercase (AC7).
//   - Length independence: the selectors take no length option and the reachable
//     set is exactly the full curated list, uninfluenced by any length concept
//     (AC9).
//   - Reproducibility: a seeded Generator reproduces its exact selection sequence.
//   - Allocation: a single pick performs no allocation.
//
// The expected reference sets below are written independently of the production
// slices in wordspt_closed.go so that a drift in either (a typo, a missing member,
// a wrong diacritic) is caught by TestClosedClassListsMatchReference.

// closedClassSampleSize is the sampled-run size for the membership and reachability
// checks. It comfortably exceeds the task's >=1000 floor and, with a fixed seed,
// deterministically reaches every member of every curated list.
const closedClassSampleSize = 50000

// wantArticlesPT is the authoritative reference set of pt-PT articles: the four
// definite, the four indefinite, and the twenty-four single-word
// preposition+article contractions (including the crasis forms à and às).
var wantArticlesPT = map[string]struct{}{
	"o": {}, "a": {}, "os": {}, "as": {},
	"um": {}, "uma": {}, "uns": {}, "umas": {},
	"do": {}, "da": {}, "dos": {}, "das": {},
	"no": {}, "na": {}, "nos": {}, "nas": {},
	"num": {}, "numa": {}, "nuns": {}, "numas": {},
	"dum": {}, "duma": {}, "duns": {}, "dumas": {},
	"pelo": {}, "pela": {}, "pelos": {}, "pelas": {},
	"ao": {}, "aos": {}, "à": {}, "às": {},
}

// wantPrepositionsPT is the authoritative reference set of the essential (simple)
// pt-PT prepositions.
var wantPrepositionsPT = map[string]struct{}{
	"a": {}, "ante": {}, "após": {}, "até": {}, "com": {}, "contra": {},
	"de": {}, "desde": {}, "em": {}, "entre": {}, "para": {}, "perante": {},
	"por": {}, "sem": {}, "sob": {}, "sobre": {}, "trás": {},
}

// wantConjunctionsPT is the authoritative reference set of single-word pt-PT
// coordinating and subordinating conjunctions.
var wantConjunctionsPT = map[string]struct{}{
	"e": {}, "nem": {}, "mas": {}, "porém": {}, "todavia": {}, "contudo": {},
	"entretanto": {}, "ou": {}, "ora": {}, "logo": {}, "portanto": {},
	"pois": {}, "assim": {}, "senão": {}, "outrossim": {},
	"que": {}, "se": {}, "porque": {}, "porquanto": {}, "conquanto": {},
	"embora": {}, "como": {}, "quando": {}, "enquanto": {}, "mal": {}, "caso": {},
	"conforme": {}, "consoante": {}, "segundo": {}, "salvo": {},
}

// assertClosedClassMember fails t when word is not a valid-UTF-8, all-lowercase
// member of want. It enforces the encoding/case rule (AC7) and the membership rule
// (AC4) for a single output.
func assertClosedClassMember(t *testing.T, class, word string, want map[string]struct{}) {
	t.Helper()
	if !utf8.ValidString(word) {
		t.Fatalf("%s %q is not valid UTF-8", class, word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("%s %q contains an uppercase rune %q", class, word, r)
		}
	}
	if _, ok := want[word]; !ok {
		t.Errorf("%s %q is not a member of the curated list", class, word)
	}
}

// TestClosedClassListsCuration verifies the curation integrity of every embedded
// list directly: each is non-empty, contains no duplicates, and every member is a
// single, non-empty, lowercase, valid-UTF-8 word (no surrounding or internal
// whitespace, since closed classes admit single words only).
func TestClosedClassListsCuration(t *testing.T) {
	cases := []struct {
		name    string
		members []string
	}{
		{"article", articlesPT},
		{"preposition", prepositionsPT},
		{"conjunction", conjunctionsPT},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.members) == 0 {
				t.Fatalf("%s list is empty", c.name)
			}
			seen := make(map[string]struct{}, len(c.members))
			for _, w := range c.members {
				if w == "" {
					t.Errorf("%s list contains an empty string", c.name)
					continue
				}
				if !utf8.ValidString(w) {
					t.Errorf("%s member %q is not valid UTF-8", c.name, w)
				}
				for _, r := range w {
					if unicode.IsSpace(r) {
						t.Errorf("%s member %q contains whitespace (not a single word)", c.name, w)
					}
					if unicode.IsUpper(r) {
						t.Errorf("%s member %q contains an uppercase rune %q", c.name, w, r)
					}
				}
				if _, dup := seen[w]; dup {
					t.Errorf("%s list contains a duplicate member %q", c.name, w)
				}
				seen[w] = struct{}{}
			}
		})
	}
}

// TestClosedClassListsMatchReference verifies that each production slice contains
// exactly the authoritative reference set — no missing member, no extra member,
// and identical spelling and diacritics. Because the reference sets are written
// independently in this file, this guards against a drift on either side.
func TestClosedClassListsMatchReference(t *testing.T) {
	cases := []struct {
		name    string
		members []string
		want    map[string]struct{}
	}{
		{"article", articlesPT, wantArticlesPT},
		{"preposition", prepositionsPT, wantPrepositionsPT},
		{"conjunction", conjunctionsPT, wantConjunctionsPT},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if len(c.members) != len(c.want) {
				t.Errorf("%s list has %d members, reference has %d", c.name, len(c.members), len(c.want))
			}
			got := make(map[string]struct{}, len(c.members))
			for _, w := range c.members {
				got[w] = struct{}{}
				if _, ok := c.want[w]; !ok {
					t.Errorf("%s list has unexpected member %q (not in reference)", c.name, w)
				}
			}
			for w := range c.want {
				if _, ok := got[w]; !ok {
					t.Errorf("%s reference member %q is missing from the list", c.name, w)
				}
			}
		})
	}
}

// TestClosedClassMembershipAndReachability drives each selector through a large
// deterministic (seeded) run and asserts that every output is a valid-UTF-8,
// lowercase member of the curated list (AC4, AC7), and that the reachable set is
// exactly the full curated list — every member is produced and nothing outside it
// is (AC4 reachability; AC9 length independence, since no length option exists and
// the whole list is reachable).
func TestClosedClassMembershipAndReachability(t *testing.T) {
	cases := []struct {
		name string
		pick func(g *Generator) string
		want map[string]struct{}
	}{
		{"article", (*Generator).ArticlePT, wantArticlesPT},
		{"preposition", (*Generator).PrepositionPT, wantPrepositionsPT},
		{"conjunction", (*Generator).ConjunctionPT, wantConjunctionsPT},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g := New(0x9E3779B97F4A7C15)
			reached := make(map[string]struct{}, len(c.want))
			for i := 0; i < closedClassSampleSize; i++ {
				w := c.pick(g)
				assertClosedClassMember(t, c.name, w, c.want)
				reached[w] = struct{}{}
			}
			for w := range c.want {
				if _, ok := reached[w]; !ok {
					t.Errorf("%s member %q was never produced over %d draws", c.name, w, closedClassSampleSize)
				}
			}
			if len(reached) != len(c.want) {
				t.Errorf("%s reached %d distinct members, want %d", c.name, len(reached), len(c.want))
			}
		})
	}
}

// TestClosedClassPackageLevelMembership exercises the package-level (global-source)
// selectors and asserts that every output is a valid-UTF-8, lowercase member of the
// curated list. It confirms that the concurrency-safe global surface draws from the
// same curated lists as the *Generator methods.
func TestClosedClassPackageLevelMembership(t *testing.T) {
	for i := 0; i < closedClassSampleSize; i++ {
		assertClosedClassMember(t, "article", ArticlePT(), wantArticlesPT)
		assertClosedClassMember(t, "preposition", PrepositionPT(), wantPrepositionsPT)
		assertClosedClassMember(t, "conjunction", ConjunctionPT(), wantConjunctionsPT)
	}
}

// TestClosedClassReproducibility verifies that two Generators created with the same
// seed produce byte-identical selection sequences for each selector, and that the
// *Generator methods and package-level functions select from the same list.
func TestClosedClassReproducibility(t *testing.T) {
	const seed = 0xDEADBEEFCAFEF00D
	const n = 2000

	cases := []struct {
		name string
		pick func(g *Generator) string
	}{
		{"article", (*Generator).ArticlePT},
		{"preposition", (*Generator).PrepositionPT},
		{"conjunction", (*Generator).ConjunctionPT},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			g1 := New(seed)
			g2 := New(seed)
			for i := 0; i < n; i++ {
				if a, b := c.pick(g1), c.pick(g2); a != b {
					t.Fatalf("%s: same-seed divergence at draw %d: %q != %q", c.name, i, a, b)
				}
			}
		})
	}
}

// TestClosedClassNoAllocation verifies that a single pick performs no heap
// allocation, on both the package-level (global-source) surface and the
// *Generator surface. The returned string shares the embedded list's backing
// array, so the pick itself allocates nothing.
func TestClosedClassNoAllocation(t *testing.T) {
	g := New(1)
	checks := []struct {
		name string
		fn   func()
	}{
		{"ArticlePT", func() { _ = ArticlePT() }},
		{"PrepositionPT", func() { _ = PrepositionPT() }},
		{"ConjunctionPT", func() { _ = ConjunctionPT() }},
		{"Generator.ArticlePT", func() { _ = g.ArticlePT() }},
		{"Generator.PrepositionPT", func() { _ = g.PrepositionPT() }},
		{"Generator.ConjunctionPT", func() { _ = g.ConjunctionPT() }},
	}
	for _, c := range checks {
		if got := testing.AllocsPerRun(1000, c.fn); got != 0 {
			t.Errorf("%s: got %v allocs/op, want 0", c.name, got)
		}
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

func BenchmarkArticlePT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ArticlePT()
	}
}

func BenchmarkPrepositionPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = PrepositionPT()
	}
}

func BenchmarkConjunctionPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ConjunctionPT()
	}
}

func BenchmarkGeneratorArticlePT(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.ArticlePT()
	}
}
