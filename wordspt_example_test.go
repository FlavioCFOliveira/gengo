package gengo_test

import (
	"fmt"

	"github.com/FlavioCFOliveira/gengo"
)

// The examples below use a seeded Generator (gengo.New) so their output is
// reproducible and can be checked. In real code, call the package-level
// functions (gengo.NounPT, gengo.VerbPT, and so on) to draw from the global,
// automatically seeded source instead.

// ExampleNounPTOf generates a pt-PT pseudo-noun inflected for a specific gender,
// number, and length category. Here the request is a feminine plural noun in the
// medium length category (5 to 8 characters).
func ExampleNounPTOf() {
	g := gengo.New(1)
	fmt.Println(g.NounPTOf(gengo.Feminine, gengo.Plural, gengo.MediumLengthWords))
	// Output: gruistas
}

// ExampleAdjectivePTOf generates a pt-PT pseudo-adjective in the synthetic
// absolute superlative degree. The masculine singular superlative always ends in
// the -íssimo suffix.
func ExampleAdjectivePTOf() {
	g := gengo.New(7)
	fmt.Println(g.AdjectivePTOf(gengo.Masculine, gengo.Singular, gengo.Superlative, gengo.AnyLengthWord))
	// Output: racantíssimo
}

// ExampleVerbPTOf generates a pt-PT pseudo-verb in a specific paradigm slot: the
// present indicative, first person singular.
func ExampleVerbPTOf() {
	g := gengo.New(3)
	fmt.Println(g.VerbPTOf(gengo.Indicative, gengo.Present, gengo.First, gengo.Singular, gengo.AnyLengthWord))
	// Output: onalgro
}

// ExampleAdverbPT generates a pt-PT pseudo-adverb. Every generated adverb is a
// productive -mente form and is invariable.
func ExampleAdverbPT() {
	g := gengo.New(6)
	fmt.Println(g.AdverbPT())
	// Output: dilrantemente
}

// ExamplePrepositionPT selects a real pt-PT preposition. Closed-class selectors
// return genuine, curated words rather than generated pseudo-words.
func ExamplePrepositionPT() {
	g := gengo.New(5)
	fmt.Println(g.PrepositionPT())
	// Output: perante
}

// ExampleWordsPT generates a slice of random pt-PT words following the natural
// distribution across word classes. Open-class words are generated pseudo-words,
// while closed-class words (such as the numeral "quatro" below) are real words.
func ExampleWordsPT() {
	g := gengo.New(6)
	fmt.Println(g.WordsPT(5))
	// Output: [quatro escas sesmando murnalmente seblares]
}

// ExampleGenerator shows that two generators created with the same seed produce
// the same sequence of words, which is what makes seeded output reproducible.
func ExampleGenerator() {
	a := gengo.New(7)
	b := gengo.New(7)
	fmt.Println(a.WordsPT(3))
	fmt.Println(b.WordsPT(3))
	// Output:
	// [a sedoraríamos teichãospores]
	// [a sedoraríamos teichãospores]
}
