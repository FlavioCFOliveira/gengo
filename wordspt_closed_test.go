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

// wantPronounsPT is the authoritative reference set of pt-PT pronouns across every
// subtype, written independently of pronounsPT so that any drift on either side
// (a typo, a missing form, a wrong diacritic, a lost vós form) is caught. It is
// duplicate-free: pronouns belonging to more than one subtype appear once.
var wantPronounsPT = map[string]struct{}{
	// Personal — subject.
	"eu": {}, "tu": {}, "ele": {}, "ela": {}, "nós": {}, "vós": {}, "eles": {}, "elas": {}, //nolint:misspell // real pt-PT pronouns; US misspell mistakes the plural form for the English "eels"
	// Personal — oblique atonic.
	"me": {}, "te": {}, "se": {}, "o": {}, "a": {}, "lhe": {}, "nos": {}, "vos": {}, "os": {}, "as": {}, "lhes": {},
	// Personal — oblique tonic.
	"mim": {}, "ti": {}, "si": {}, "comigo": {}, "contigo": {}, "consigo": {}, "connosco": {}, "convosco": {},
	// Possessive.
	"meu": {}, "minha": {}, "meus": {}, "minhas": {},
	"teu": {}, "tua": {}, "teus": {}, "tuas": {},
	"seu": {}, "sua": {}, "seus": {}, "suas": {},
	"nosso": {}, "nossa": {}, "nossos": {}, "nossas": {},
	"vosso": {}, "vossa": {}, "vossos": {}, "vossas": {},
	// Demonstrative.
	"este": {}, "esta": {}, "estes": {}, "estas": {}, "isto": {},
	"esse": {}, "essa": {}, "esses": {}, "essas": {}, "isso": {},
	"aquele": {}, "aquela": {}, "aqueles": {}, "aquelas": {}, "aquilo": {},
	"mesmo": {}, "mesma": {}, "mesmos": {}, "mesmas": {},
	"próprio": {}, "própria": {}, "próprios": {}, "próprias": {},
	"tal": {}, "tais": {},
	"semelhante": {}, "semelhantes": {},
	// Indefinite.
	"algum": {}, "alguma": {}, "alguns": {}, "algumas": {},
	"nenhum": {}, "nenhuma": {}, "nenhuns": {}, "nenhumas": {},
	"todo": {}, "toda": {}, "todos": {}, "todas": {},
	"outro": {}, "outra": {}, "outros": {}, "outras": {},
	"muito": {}, "muita": {}, "muitos": {}, "muitas": {},
	"pouco": {}, "pouca": {}, "poucos": {}, "poucas": {},
	"tanto": {}, "tanta": {}, "tantos": {}, "tantas": {},
	"quanto": {}, "quanta": {}, "quantos": {}, "quantas": {},
	"vário": {}, "vária": {}, "vários": {}, "várias": {},
	"certo": {}, "certa": {}, "certos": {}, "certas": {},
	"qualquer": {}, "quaisquer": {},
	"alguém": {}, "ninguém": {}, "tudo": {}, "nada": {}, "algo": {}, "cada": {}, "outrem": {},
	// Relative.
	"que": {}, "quem": {}, "qual": {}, "quais": {},
	"cujo": {}, "cuja": {}, "cujos": {}, "cujas": {},
	// Interrogative: every interrogative pronoun is also relative or indefinite
	// (quem/que/qual/quais and the quanto series), so no member is unique here.
}

// wantInterjectionsPT is the authoritative reference set of the curated pt-PT
// single-word interjections, written independently of interjectionsPT to guard
// against drift.
var wantInterjectionsPT = map[string]struct{}{
	"ah": {}, "oh": {}, "ó": {}, "olá": {}, "oi": {}, "ui": {}, "ai": {}, "eh": {},
	"hã": {}, "hum": {}, "uf": {}, "ufa": {}, "oxalá": {}, "tomara": {}, "viva": {},
	"bravo": {}, "olé": {}, "chiu": {}, "psiu": {}, "bolas": {}, "caramba": {},
	"credo": {}, "coitado": {}, "adeus": {}, "alto": {}, "avante": {}, "arre": {},
	"upa": {}, "eia": {}, "salve": {}, "ora": {}, "apre": {},
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
		{"pronoun", pronounsPT},
		{"interjection", interjectionsPT},
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
		{"pronoun", pronounsPT, wantPronounsPT},
		{"interjection", interjectionsPT, wantInterjectionsPT},
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
		{"pronoun", (*Generator).PronounPT, wantPronounsPT},
		{"interjection", (*Generator).InterjectionPT, wantInterjectionsPT},
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
		assertClosedClassMember(t, "pronoun", PronounPT(), wantPronounsPT)
		assertClosedClassMember(t, "interjection", InterjectionPT(), wantInterjectionsPT)
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
		{"pronoun", (*Generator).PronounPT},
		{"interjection", (*Generator).InterjectionPT},
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
		{"PronounPT", func() { _ = PronounPT() }},
		{"InterjectionPT", func() { _ = InterjectionPT() }},
		{"Generator.ArticlePT", func() { _ = g.ArticlePT() }},
		{"Generator.PrepositionPT", func() { _ = g.PrepositionPT() }},
		{"Generator.ConjunctionPT", func() { _ = g.ConjunctionPT() }},
		{"Generator.PronounPT", func() { _ = g.PronounPT() }},
		{"Generator.InterjectionPT", func() { _ = g.InterjectionPT() }},
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

func BenchmarkPronounPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = PronounPT()
	}
}

func BenchmarkInterjectionPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = InterjectionPT()
	}
}

func BenchmarkGeneratorPronounPT(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.PronounPT()
	}
}

// TestPronounPTSubtypeCoverage asserts, member by member, that the curated pronoun
// list spans every subtype required by the task and specification/wordspt-closed-
// classes.md: personal subject, personal atonic oblique, personal tonic oblique,
// possessive, demonstrative, indefinite, relative, and interrogative. It checks a
// representative sentinel from each subtype directly against the production list,
// so a whole subtype cannot silently disappear.
func TestPronounPTSubtypeCoverage(t *testing.T) {
	got := make(map[string]struct{}, len(pronounsPT))
	for _, w := range pronounsPT {
		got[w] = struct{}{}
	}
	subtypes := map[string][]string{
		"personal-subject":        {"eu", "tu", "ele", "ela", "nós", "vós", "eles", "elas"}, //nolint:misspell // real pt-PT pronouns; US misspell mistakes the plural form for the English "eels"
		"personal-atonic-oblique": {"me", "te", "se", "lhe", "nos", "vos", "lhes"},
		"personal-tonic-oblique":  {"mim", "ti", "si", "comigo", "contigo", "consigo", "connosco", "convosco"},
		"possessive":              {"meu", "minha", "teu", "tua", "seu", "sua", "nosso", "nossa", "vosso", "vossa"},
		"demonstrative":           {"este", "esta", "isto", "esse", "essa", "isso", "aquele", "aquilo", "mesmo", "próprio", "tal", "semelhante"},
		"indefinite":              {"algum", "nenhum", "todo", "outro", "muito", "pouco", "tanto", "vário", "certo", "qualquer", "alguém", "ninguém", "tudo", "nada", "algo", "cada", "outrem"},
		"relative":                {"que", "quem", "qual", "quais", "cujo", "cuja"},
		// Every interrogative pronoun is also relative or indefinite (there is no
		// interrogative-unique member); assert the genuine interrogative pronouns
		// are present rather than requiring a distinct-from-relative one.
		"interrogative": {"quem", "que", "qual", "quais", "quanto", "quanta", "quantos", "quantas"},
	}
	for subtype, members := range subtypes {
		for _, m := range members {
			if _, ok := got[m]; !ok {
				t.Errorf("pronoun subtype %q: required member %q is missing from pronounsPT", subtype, m)
			}
		}
	}
}

// TestPronounPTIncludesVosForms guards the task's explicit requirement that the
// archaic-but-real second-person-plural forms related to "vós" are present. These
// are easy to drop as "obsolete", so they are asserted individually.
func TestPronounPTIncludesVosForms(t *testing.T) {
	got := make(map[string]struct{}, len(pronounsPT))
	for _, w := range pronounsPT {
		got[w] = struct{}{}
	}
	vosForms := []string{"vós", "vos", "convosco", "vosso", "vossa", "vossos", "vossas"}
	for _, f := range vosForms {
		if _, ok := got[f]; !ok {
			t.Errorf("required vós-related form %q is missing from pronounsPT", f)
		}
	}
}
