package gengo

import "math/rand/v2"

// This file implements the top-level WordsPT orchestration generators: the
// package-level WordPT, WordPTByLengthType and WordsPT, plus their mirrored
// *Generator methods. They sit above the ten per-class generators built by the
// earlier tasks and return a random pt-PT word of a random class:
//
//	WordPT()                       -> one random word of a random class
//	WordPTByLengthType(l)          -> as WordPT, with l applied to the open classes
//	WordsPT(n)                     -> n random words (empty, non-nil slice for n<=0)
//
// Dispatch model
//
// WordPT draws a word class from a WEIGHTED distribution (resolveWordClass) and
// delegates to that class's generation CORE with the injected source, never to the
// package-level wrapper. Delegating to the core (nounCore, adjectiveCore, verbCore,
// adverbCore, pickClosedClass) avoids a double dispatch: the wrapper would allocate
// a fresh stack buffer and re-enter through the global source, whereas the core
// takes the already-injected *rand.Rand and the already-allocated shared syllable
// buffer directly. The open-class cores request every inflection option as its Any
// value, so a random gender, number, degree and verb slot are drawn per word.
//
// Class distribution (favoring the open/content classes)
//
// The specification's "natural distribution" approximates the relative frequency of
// each word class in European-Portuguese running text. In running text the most
// frequent TOKENS are function words (articles, prepositions, conjunctions, and a
// few pronouns such as "de", "a", "o", "que", "e"); a handful of those short words
// account for a large share of every text. That token frequency is the wrong target
// for a random-WORD generator: the six closed classes are small, fixed inventories
// (a few dozen words each), so weighting them by their running-text token share
// would make WordPT return the same handful of function words most of the time and
// almost never exercise the productive, generated open classes that are the point of
// the feature.
//
// The distribution below therefore favors the open/content classes as the bulk and
// keeps the six closed classes at a small combined share. This is grounded in the
// LEXICAL (type) frequency of Portuguese rather than its token frequency: the
// content classes vastly outnumber the function classes in the lexicon, with nouns
// the largest class, verbs and adjectives comparably large, and adverbs (the
// productive -mente class) smaller; the closed classes are a fixed, small tail. All
// ten classes remain reachable (every weight is at least one). The concrete weights
// are ordinal estimates (flagged), consistent with every other WordsPT weight table
// and with how WordLengthRatio encodes the length distribution; they change only how
// often each class appears, never the conformance of any output.
//
// Length
//
// Length applies only to the open classes; the closed classes ignore it, because
// their output is a fixed real word (specification: "Length applies only to open
// classes"). WordPTByLengthType applies a CONCRETE requested category to the chosen
// open class directly. When the request is AnyLengthWord (the WordPT default), the
// open-class length is drawn from WordLengthRatio, the same natural length
// distribution the existing Words generator uses, so WordsPT genuinely mirrors
// Words: each open-class word's length category follows WordLengthRatio, and the
// open core then applies its own within-category length skew and any upward
// minimum-viable-length normalization. Closed classes never consult the length.
//
// Randomness surfaces. The package-level WordPT/WordPTByLengthType/WordsPT draw from
// the global, automatically seeded math/rand/v2 source (through wordsPTGlobalRand);
// the (*Generator) methods draw from the generator's own seeded source. Both call
// the single shared core wordPTCore with an injectable *rand.Rand, so there is no
// duplication of dispatch logic and a seeded Generator is fully reproducible.
//
// Allocation. An open-class word is one string allocation (the class core's final
// assembly), and the orchestration adds none: it reuses a single stack syllable
// buffer across the whole WordsPT slice and passes it straight to the open core. A
// closed-class word allocates nothing (the selector returns a string that shares the
// curated list's backing array). So WordPT performs at most one allocation per word.

// ---------------------------------------------------------------------------
// Word-class distribution
// ---------------------------------------------------------------------------

// wordClass is an internal identifier for one of the ten Portuguese word classes
// WordPT can produce. It is unexported because the class is an implementation
// detail of the orchestration: callers select a class only indirectly, through the
// weighted distribution.
type wordClass uint8

const (
	classNoun         wordClass = iota // open
	classAdjective                     // open
	classVerb                          // open
	classAdverb                        // open
	classArticle                       // closed
	classPronoun                       // closed
	classNumeral                       // closed
	classPreposition                   // closed
	classConjunction                   // closed
	classInterjection                  // closed
)

// weightedWordClass pairs a word class with its ordinal sampling weight.
type weightedWordClass struct {
	class  wordClass
	weight uint32
}

// wordClasses is the weighted class distribution WordPT draws from. The weights sum
// to 100, so each reads directly as a percentage: the four open classes take 88% of
// the mass (noun 35, verb 22, adjective 21, adverb 10) and the six closed classes
// share the remaining 12% (pronoun 3, numeral 2, preposition 2, conjunction 2,
// article 2, interjection 1). The open bulk and the small closed tail encode the
// design decision documented in this file's header: favor the productive content
// classes, which dominate the lexicon, over the small fixed function-word
// inventories, which dominate running-text token counts. Every weight is at least
// one, so every class is reachable. The values are ordinal estimates (flagged).
var wordClasses = []weightedWordClass{
	{classNoun, 35},
	{classAdjective, 21},
	{classVerb, 22},
	{classAdverb, 10},
	{classArticle, 2},
	{classPronoun, 3},
	{classNumeral, 2},
	{classPreposition, 2},
	{classConjunction, 2},
	{classInterjection, 1},
}

// wordClassesTotal is the total class weight, the exclusive upper bound for the
// single bounded draw in [resolveWordClass]. It is computed once, at package
// initialization, mirroring verbMoodsTotal.
var wordClassesTotal = func() uint32 {
	var total uint32
	for i := range wordClasses {
		total += wordClasses[i].weight
	}
	return total
}()

// resolveWordClass draws a word class from the weighted distribution using the
// injectable source r: a single bounded selection and a short linear scan, with no
// allocation and no rejection. It mirrors resolveMood.
func resolveWordClass(r *rand.Rand) wordClass {
	v := r.Uint32N(wordClassesTotal)
	var acc uint32
	for i := range wordClasses {
		acc += wordClasses[i].weight
		if v < acc {
			return wordClasses[i].class
		}
	}
	return classNoun // unreachable: v < wordClassesTotal always matches an entry
}

// openClassLength resolves the length category passed to an open-class core. A
// concrete category ([SmallLengthWord], [MediumLengthWords] or [BigLengthWords]) is
// honored as requested. AnyLengthWord (and any undefined value), the WordPT default,
// is resolved to a concrete category through WordLengthRatio, so the open-class
// words of WordPT and WordsPT follow the same natural length distribution the
// existing Words generator uses (specification: Natural distribution). It draws from
// r only when it must resolve an unspecified length, so a concrete request consumes
// no randomness.
func openClassLength(r *rand.Rand, l LengthTypeWords) LengthTypeWords {
	switch l {
	case SmallLengthWord, MediumLengthWords, BigLengthWords:
		return l
	default: // AnyLengthWord or any undefined value
		switch WordLengthRatio[r.IntN(len(WordLengthRatio))] {
		case '1':
			return SmallLengthWord
		case '2':
			return MediumLengthWords
		default:
			return BigLengthWords
		}
	}
}

// ---------------------------------------------------------------------------
// Shared orchestration core
// ---------------------------------------------------------------------------

// wordPTCore is the single dispatch engine shared by the package-level WordPT
// functions and the *Generator WordPT methods. It draws only from r, so a seeded
// Generator is fully reproducible, and threads a caller-provided syllable buffer
// (backed by a stack array) into whichever open-class core it selects, so an
// open-class word costs exactly one string allocation and a closed-class word costs
// none.
//
// It draws a class from the weighted distribution and delegates to that class's
// core with the injected source. Each open class requests every inflection option
// as its Any value and receives the length resolved by [openClassLength]; each
// closed class is a uniform pick over its curated list and ignores the length. The
// output is always a lowercase, valid-UTF-8 pt-PT word.
func wordPTCore(r *rand.Rand, buf []syllable, l LengthTypeWords) string {
	switch resolveWordClass(r) {
	case classNoun:
		return nounCore(r, buf, AnyGender, AnyNumber, openClassLength(r, l))
	case classAdjective:
		return adjectiveCore(r, buf, AnyGender, AnyNumber, AnyDegree, openClassLength(r, l))
	case classVerb:
		return verbCore(r, buf, AnyMood, AnyTense, AnyPerson, AnyNumber, openClassLength(r, l))
	case classAdverb:
		return adverbCore(r, buf, openClassLength(r, l))
	case classArticle:
		return pickClosedClass(r, articlesPT)
	case classPronoun:
		return pickClosedClass(r, pronounsPT)
	case classNumeral:
		return pickClosedClass(r, numeralsPT)
	case classPreposition:
		return pickClosedClass(r, prepositionsPT)
	case classConjunction:
		return pickClosedClass(r, conjunctionsPT)
	case classInterjection:
		return pickClosedClass(r, interjectionsPT)
	default: // unreachable: resolveWordClass returns one of the ten classes
		return nounCore(r, buf, AnyGender, AnyNumber, openClassLength(r, l))
	}
}

// ---------------------------------------------------------------------------
// Public API: package-level functions
// ---------------------------------------------------------------------------

// WordPT returns a random European-Portuguese (pt-PT) word of a random class,
// following the natural class distribution described in this package. It is exact
// sugar for WordPTByLengthType(AnyLengthWord). The result is a lowercase, valid-UTF-8
// pt-PT word: a generated pseudo-word for the open classes (noun, adjective, verb,
// adverb) or a real curated word for the closed classes (article, pronoun, numeral,
// preposition, conjunction, interjection).
//
// WordPT draws from the global, automatically seeded math/rand/v2 source and is safe
// for concurrent use by multiple goroutines. For reproducible output, use a seeded
// Generator ([New] or [NewSource]) and its WordPT method.
func WordPT() string {
	return WordPTByLengthType(AnyLengthWord)
}

// WordPTByLengthType returns a random pt-PT word of a random class, as [WordPT], with
// the length category l applied to the open classes. A concrete category
// ([SmallLengthWord], [MediumLengthWords] or [BigLengthWords]) constrains a generated
// open-class word to that character window, normalizing upward when the category
// cannot hold the chosen class (length is never normalized downward). [AnyLengthWord]
// (the zero value) selects the natural length distribution. The closed classes ignore
// l entirely, because their output is a fixed real word. The result is always
// lowercase and valid UTF-8; the function never returns an error and never panics.
//
// WordPTByLengthType draws from the global, automatically seeded math/rand/v2 source
// and is safe for concurrent use by multiple goroutines.
func WordPTByLengthType(l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return wordPTCore(wordsPTGlobalRand, buf[:0], l)
}

// WordsPT returns a slice of n random pt-PT words following the natural class
// distribution, mirroring the semantics of [Words]. It returns a non-nil, empty
// slice when n is less than or equal to zero, and a slice of exactly n words
// otherwise. Each word is drawn independently, exactly as [WordPT] draws one.
//
// WordsPT draws from the global, automatically seeded math/rand/v2 source and is safe
// for concurrent use by multiple goroutines. For reproducible output, use a seeded
// Generator ([New] or [NewSource]) and its WordsPT method.
func WordsPT(n int) []string {
	if n <= 0 {
		return []string{}
	}
	result := make([]string, n)
	var buf [nounMaxSyllables]syllable
	for i := 0; i < n; i++ {
		result[i] = wordPTCore(wordsPTGlobalRand, buf[:0], AnyLengthWord)
	}
	return result
}

// ---------------------------------------------------------------------------
// Public API: *Generator methods
// ---------------------------------------------------------------------------

// WordPT is the seeded-generator equivalent of [WordPT]. A Generator created with the
// same seed produces the same sequence of words.
func (g *Generator) WordPT() string {
	return g.WordPTByLengthType(AnyLengthWord)
}

// WordPTByLengthType is the seeded-generator equivalent of [WordPTByLengthType]. A
// Generator created with the same seed produces the same sequence of words for the
// same length category.
func (g *Generator) WordPTByLengthType(l LengthTypeWords) string {
	var buf [nounMaxSyllables]syllable
	return wordPTCore(g.r, buf[:0], l)
}

// WordsPT is the seeded-generator equivalent of [WordsPT]. A Generator created with
// the same seed produces the same sequence of word slices.
func (g *Generator) WordsPT(n int) []string {
	if n <= 0 {
		return []string{}
	}
	result := make([]string, n)
	var buf [nounMaxSyllables]syllable
	for i := 0; i < n; i++ {
		result[i] = wordPTCore(g.r, buf[:0], AnyLengthWord)
	}
	return result
}
