package gengo

import (
	"math/rand/v2"
	"testing"
	"unicode"
	"unicode/utf8"
)

// This file verifies the top-level WordsPT orchestration generators (WordPT,
// WordPTByLengthType, WordsPT, and their *Generator methods). The checks map to the
// task acceptance criteria and the specification: every class is reachable, the open
// content classes dominate the mix, every output is an orthographically well-formed
// (or curated) lowercase valid-UTF-8 word, the length category is applied to the
// open classes and ignored by the closed classes, WordsPT honors the Words-style
// slice contract (empty non-nil slice for n<=0, exactly n otherwise), a seeded
// Generator is reproducible, and an open-class word costs at most one allocation.

// wordPTSampleSize is the sampled-run size for the conformance tests. It exceeds the
// task's >=1000 floor comfortably while keeping the suite fast.
const wordPTSampleSize = 40000

// closedClassSets pairs each closed class with its curated list, for building the
// union of all closed-class members and for per-class reachability checks.
var closedClassSets = []struct {
	name    string
	members []string
}{
	{"article", articlesPT},
	{"pronoun", pronounsPT},
	{"numeral", numeralsPT},
	{"preposition", prepositionsPT},
	{"conjunction", conjunctionsPT},
	{"interjection", interjectionsPT},
}

// allClosedMembers returns the set union of every curated closed-class list. A
// WordPT output that is a member of this set is a curated real word (always valid);
// an output outside it is a generated open-class pseudo-word, checked against the
// strict open-class orthography oracle.
func allClosedMembers() map[string]struct{} {
	set := make(map[string]struct{}, 256)
	for _, c := range closedClassSets {
		for _, w := range c.members {
			set[w] = struct{}{}
		}
	}
	return set
}

// assertWordPTValid fails t when word is not a valid WordPT output. Every output,
// regardless of class, must be a non-empty, valid-UTF-8, all-lowercase string. A
// curated closed-class member is valid by construction (it may legitimately carry a
// grave accent, as in the article "à", which the open-class oracle forbids), so only
// a generated open-class word (one outside the closed set) is held to the strict
// Sprint 7 open-class oracle: the cedilla rule, no forbidden diaeresis, and no grave
// accent.
func assertWordPTValid(t *testing.T, word string, closed map[string]struct{}) {
	t.Helper()
	if word == "" {
		t.Fatalf("WordPT produced an empty string")
	}
	if !utf8.ValidString(word) {
		t.Fatalf("WordPT %q is not valid UTF-8", word)
	}
	for _, r := range word {
		if unicode.IsUpper(r) {
			t.Errorf("WordPT %q contains an uppercase rune %q", word, r)
		}
	}
	if _, ok := closed[word]; ok {
		return // curated closed-class member: valid by construction
	}
	if !wordConformsToCedilla(word) {
		t.Errorf("open-class WordPT %q violates the pt-PT cedilla rule", word)
	}
	if wordHasForbiddenDiaeresis(word) {
		t.Errorf("open-class WordPT %q contains a forbidden diaeresis", word)
	}
	if wordContainsGrave(word) {
		t.Errorf("open-class WordPT %q contains a grave accent", word)
	}
}

// TestResolveWordClassReachesAll verifies that every one of the ten word classes is
// reachable through the weighted distribution: over a large deterministic run, each
// class is drawn at least once. This is the authoritative reachability proof (the
// class is chosen internally, so it cannot be read back unambiguously from an
// output surface where closed lists overlap and open classes are indistinguishable).
func TestResolveWordClassReachesAll(t *testing.T) {
	r := rand.New(rand.NewPCG(0x5EED, 0xC1A55))
	counts := make(map[wordClass]int, len(wordClasses))
	for i := 0; i < 500000; i++ {
		counts[resolveWordClass(r)]++
	}
	for _, wc := range wordClasses {
		if counts[wc.class] == 0 {
			t.Errorf("word class %d was never drawn over the run", wc.class)
		}
	}
}

// TestResolveWordClassDistribution verifies that the drawn distribution matches the
// weighted table: the four open classes dominate (about 88% of the mass), the noun
// is the single most frequent class, and each class lands within a small tolerance of
// its configured share. The run is deterministic (fixed seed), so the check never
// flakes.
func TestResolveWordClassDistribution(t *testing.T) {
	const n = 500000
	r := rand.New(rand.NewPCG(0xA11CE, 0xB0B))
	counts := make(map[wordClass]int, len(wordClasses))
	for i := 0; i < n; i++ {
		counts[resolveWordClass(r)]++
	}

	// Expected share per class, from the weight table (weights sum to 100).
	expected := map[wordClass]float64{
		classNoun: 0.35, classAdjective: 0.21, classVerb: 0.22, classAdverb: 0.10,
		classArticle: 0.02, classPronoun: 0.03, classNumeral: 0.02,
		classPreposition: 0.02, classConjunction: 0.02, classInterjection: 0.01,
	}
	for wc, want := range expected {
		got := float64(counts[wc]) / float64(n)
		if got < want-0.02 || got > want+0.02 {
			t.Errorf("class %d share = %.4f, want %.2f +/- 0.02", wc, got, want)
		}
	}

	open := counts[classNoun] + counts[classAdjective] + counts[classVerb] + counts[classAdverb]
	openFrac := float64(open) / float64(n)
	if openFrac < 0.85 || openFrac > 0.91 {
		t.Errorf("open-class fraction = %.4f, want ~0.88 (open must dominate)", openFrac)
	}

	// The noun must be the single largest class.
	for wc, c := range counts {
		if wc != classNoun && c >= counts[classNoun] {
			t.Errorf("class %d count %d >= noun count %d; noun must be the largest", wc, c, counts[classNoun])
		}
	}
}

// TestWordPTOutputValidity samples the package-level WordPT and the *Generator
// WordPT and asserts every output is a valid WordPT word: non-empty, valid UTF-8,
// lowercase, and (for the generated open classes) orthographically conformant.
func TestWordPTOutputValidity(t *testing.T) {
	closed := allClosedMembers()
	for i := 0; i < wordPTSampleSize; i++ {
		assertWordPTValid(t, WordPT(), closed)
	}
	g := New(0xB0B)
	for i := 0; i < wordPTSampleSize; i++ {
		assertWordPTValid(t, g.WordPT(), closed)
	}
}

// TestWordPTClassMix verifies at the output surface that WordPT mixes generated
// open-class pseudo-words with curated closed-class words: over a sampled run at
// least one output falls outside every closed list (an open pseudo-word), and each
// of the six closed lists contributes at least one member. Combined with
// TestResolveWordClassReachesAll, this confirms both the open and the closed classes
// surface through the public API.
func TestWordPTClassMix(t *testing.T) {
	closed := allClosedMembers()
	perList := make([]map[string]struct{}, len(closedClassSets))
	for i, c := range closedClassSets {
		set := make(map[string]struct{}, len(c.members))
		for _, w := range c.members {
			set[w] = struct{}{}
		}
		perList[i] = set
	}

	g := New(0x1234567)
	sawOpen := false
	sawListMember := make([]bool, len(closedClassSets))
	for i := 0; i < wordPTSampleSize; i++ {
		w := g.WordPT()
		if _, isClosed := closed[w]; !isClosed {
			sawOpen = true
			continue
		}
		for j := range perList {
			if _, ok := perList[j][w]; ok {
				sawListMember[j] = true
			}
		}
	}
	if !sawOpen {
		t.Error("no open-class pseudo-word appeared over the run")
	}
	for i, c := range closedClassSets {
		if !sawListMember[i] {
			t.Errorf("no %s member appeared over the run", c.name)
		}
	}
}

// TestWordPTByLengthTypeBigOpenWindow verifies that a concrete length category is
// applied to the open classes: for a BigLengthWords request, every generated
// open-class output (an output outside every closed list) falls within the Big
// character window of 9..30 runes. Big is viable for every open class and never
// normalizes, so the window is exact. Closed-class outputs are skipped, since they
// ignore length by design.
func TestWordPTByLengthTypeBigOpenWindow(t *testing.T) {
	closed := allClosedMembers()
	g := New(0xB16)
	for i := 0; i < wordPTSampleSize; i++ {
		w := g.WordPTByLengthType(BigLengthWords)
		if _, isClosed := closed[w]; isClosed {
			continue // closed classes ignore length
		}
		n := utf8.RuneCountInString(w)
		if n < 9 || n > 30 {
			t.Fatalf("Big open-class word %q has %d runes, outside window [9,30]", w, n)
		}
	}
}

// TestWordPTClosedIgnoreLength verifies that closed classes ignore the requested
// length: over a SmallLengthWord run, at least one curated closed-class word longer
// than the four-character Small ceiling still appears, proving that the length
// category never constrains a closed-class output (specification AC9).
func TestWordPTClosedIgnoreLength(t *testing.T) {
	closed := allClosedMembers()
	g := New(0x5117)
	sawLongClosed := false
	for i := 0; i < wordPTSampleSize; i++ {
		w := g.WordPTByLengthType(SmallLengthWord)
		if _, isClosed := closed[w]; isClosed && utf8.RuneCountInString(w) > 4 {
			sawLongClosed = true
			break
		}
	}
	if !sawLongClosed {
		t.Error("no closed-class word longer than the Small ceiling appeared under a Small request")
	}
}

// TestWordsPTCount verifies the WordsPT slice contract for positive n: the result
// has exactly n elements, and every element is a valid word. Both the package-level
// function and the *Generator method are checked.
func TestWordsPTCount(t *testing.T) {
	closed := allClosedMembers()
	for _, n := range []int{1, 2, 5, 50, 1000} {
		got := WordsPT(n)
		if len(got) != n {
			t.Fatalf("WordsPT(%d) returned %d words, want %d", n, len(got), n)
		}
		for _, w := range got {
			assertWordPTValid(t, w, closed)
		}

		g := New(uint64(n))
		gg := g.WordsPT(n)
		if len(gg) != n {
			t.Fatalf("Generator.WordsPT(%d) returned %d words, want %d", n, len(gg), n)
		}
		for _, w := range gg {
			assertWordPTValid(t, w, closed)
		}
	}
}

// TestWordsPTEmptyForNonPositive verifies the WordsPT empty-slice contract: for n<=0
// the result is a non-nil slice of length zero, matching the existing Words
// semantics, on both surfaces.
func TestWordsPTEmptyForNonPositive(t *testing.T) {
	g := New(1)
	for _, n := range []int{0, -1, -100} {
		if got := WordsPT(n); got == nil || len(got) != 0 {
			t.Errorf("WordsPT(%d) = %#v, want non-nil empty slice", n, got)
		}
		if got := g.WordsPT(n); got == nil || len(got) != 0 {
			t.Errorf("Generator.WordsPT(%d) = %#v, want non-nil empty slice", n, got)
		}
	}
}

// TestWordPTReproducible verifies that a seeded Generator is deterministic: two
// generators with the same seed produce identical WordPT, WordPTByLengthType, and
// WordsPT sequences.
func TestWordPTReproducible(t *testing.T) {
	a := New(0xC0FFEE)
	b := New(0xC0FFEE)
	for i := 0; i < wordPTSampleSize; i++ {
		if x, y := a.WordPT(), b.WordPT(); x != y {
			t.Fatalf("WordPT diverged at %d: %q != %q", i, x, y)
		}
	}

	lengths := []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords}
	c := New(42)
	d := New(42)
	for _, le := range lengths {
		for i := 0; i < 2000; i++ {
			if x, y := c.WordPTByLengthType(le), d.WordPTByLengthType(le); x != y {
				t.Fatalf("WordPTByLengthType(%d) diverged at %d: %q != %q", le, i, x, y)
			}
		}
	}

	e := New(7)
	f := New(7)
	for i := 0; i < 200; i++ {
		xs, ys := e.WordsPT(20), f.WordsPT(20)
		if len(xs) != len(ys) {
			t.Fatalf("WordsPT length diverged at %d: %d != %d", i, len(xs), len(ys))
		}
		for j := range xs {
			if xs[j] != ys[j] {
				t.Fatalf("WordsPT diverged at call %d index %d: %q != %q", i, j, xs[j], ys[j])
			}
		}
	}
}

// TestWordPTConcurrentSafe exercises the package-level WordPT and WordsPT from many
// goroutines. Run with -race, it demonstrates that the package-level surface is safe
// for concurrent use, as documented: its shared global source delegates each draw to
// the concurrency-safe global math/rand/v2 generator.
func TestWordPTConcurrentSafe(t *testing.T) {
	const goroutines = 16
	const perGoroutine = 3000
	done := make(chan struct{}, goroutines)
	for gr := 0; gr < goroutines; gr++ {
		go func() {
			defer func() { done <- struct{}{} }()
			for i := 0; i < perGoroutine; i++ {
				if WordPT() == "" {
					t.Error("WordPT returned an empty string")
					return
				}
				for _, w := range WordsPT(4) {
					if w == "" {
						t.Error("WordsPT returned an empty string")
						return
					}
				}
			}
		}()
	}
	for gr := 0; gr < goroutines; gr++ {
		<-done
	}
}

// TestWordPTAllocationBudget verifies that the orchestration adds no allocation of
// its own: a WordPT word costs only what the delegated class core costs. A curated
// closed-class word is zero allocations (the selector shares the list's backing
// array); a generated open-class word is one allocation for the singular, verb,
// adverb and superlative paths, and two for a positive noun or adjective plural (the
// documented tail concatenation of the shared pluralizer). testing.AllocsPerRun
// truncates the per-run average to an integer, so this asserts the per-word average
// stays below two; a regression that made the class dispatch, the shared syllable
// buffer, or the length resolution allocate would push it to two or more.
func TestWordPTAllocationBudget(t *testing.T) {
	if got := testing.AllocsPerRun(5000, func() { _ = WordPT() }); got > 1 {
		t.Errorf("WordPT: got %v allocs/op (truncated), want <= 1 (average below 2)", got)
	}
	g := New(1)
	if got := testing.AllocsPerRun(5000, func() { _ = g.WordPT() }); got > 1 {
		t.Errorf("Generator.WordPT: got %v allocs/op (truncated), want <= 1 (average below 2)", got)
	}
}

// ---------------------------------------------------------------------------
// Benchmarks
// ---------------------------------------------------------------------------

// BenchmarkWordPT measures the default package-level orchestration path over the
// global source. Its reported allocs/op reflects the open/closed mix: a closed-class
// word is zero allocations, an open-class word is one (two for a positive noun or
// adjective plural, whose shared pluralizer adds a tail concatenation).
func BenchmarkWordPT(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = WordPT()
	}
}

// BenchmarkGeneratorWordPT measures the seeded-generator orchestration path.
func BenchmarkGeneratorWordPT(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.WordPT()
	}
}

// BenchmarkWordsPT measures producing a slice of words, reusing one syllable buffer
// across the whole slice.
func BenchmarkWordsPT(b *testing.B) {
	g := New(1)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = g.WordsPT(16)
	}
}
