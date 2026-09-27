package gengo

import (
	"io/fs"
	"math"
	"math/rand/v2"
	"os"
	"slices"
	"strings"
	"testing"
)

// This file guards the precomputed sampling tables of the WordsPT core: the
// per-inventory value-to-index lookup tables and onset nucleus licenses
// (wordspt_data.go), and the generated word-length threshold table
// (wordspt_skew_table.go). Each table replaces an equivalent computation, so the
// tests prove that the table and the computation it replaces agree.

// allWeightedInventories lists every package-level weighted inventory.
// TestAllWeightedInventoriesListed proves that the list is complete.
var allWeightedInventories = []struct {
	name string
	inv  *weightedInventory
}{
	{"onsetSingles", &onsetSingles},
	{"onsetClusters", &onsetClusters},
	{"nuclei", &nuclei},
	{"codas", &codas},
	{"onsetInitialInv", &onsetInitialInv},
	{"onsetMedialInv", &onsetMedialInv},
	{"nucleiFrontInv", &nucleiFrontInv},
	{"nucleiBackInv", &nucleiBackInv},
	{"codaBeforePBInv", &codaBeforePBInv},
	{"codaBeforeOtherInv", &codaBeforeOtherInv},
	{"codaFinalInv", &codaFinalInv},
	{"nounOnsetBeforeFrontInv", &nounOnsetBeforeFrontInv},
	{"nounOnsetBeforeBackInv", &nounOnsetBeforeBackInv},
	{"verbRadicalOnsetInv", &verbRadicalOnsetInv},
}

// TestAllWeightedInventoriesListed proves that allWeightedInventories names
// exactly the package-level variables initialized by newWeightedInventory or
// deriveInventory in the package's source files, so every table test below
// reaches every inventory, including any inventory added later.
func TestAllWeightedInventoriesListed(t *testing.T) {
	declared := declaredWeightedInventories(t)
	listed := make([]string, 0, len(allWeightedInventories))
	for _, it := range allWeightedInventories {
		listed = append(listed, it.name)
	}
	slices.Sort(declared)
	slices.Sort(listed)
	if !slices.Equal(declared, listed) {
		t.Fatalf("declared inventories %v, listed %v", declared, listed)
	}
}

// declaredWeightedInventories scans the package's gofmt-formatted non-test
// source files and returns the names of the package-level variables initialized
// by newWeightedInventory or deriveInventory. It recognizes both
// "var name = call(" and the "name = call(" line of a grouped var block. It
// deliberately scans lines instead of importing go/parser: linking the go/*
// packages into the test binary shifts the address of runtime code and
// measurably perturbs unrelated nanosecond-scale benchmarks.
func declaredWeightedInventories(t *testing.T) []string {
	t.Helper()
	fsys := os.DirFS(".")
	names, err := fs.Glob(fsys, "*.go")
	if err != nil {
		t.Fatal(err)
	}
	var declared []string
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := fs.ReadFile(fsys, name)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(string(src), "package gengo\n") && !strings.Contains(string(src), "\npackage gengo\n") {
			continue // the go:generate program (package main)
		}
		for line := range strings.Lines(string(src)) {
			line = strings.TrimPrefix(strings.TrimLeft(line, "\t"), "var ")
			ident, call, ok := strings.Cut(line, " = ")
			if ok && !strings.ContainsAny(ident, " ,\t") &&
				(strings.HasPrefix(call, "newWeightedInventory(") || strings.HasPrefix(call, "deriveInventory(")) {
				declared = append(declared, ident)
			}
		}
	}
	return declared
}

// TestInventoryLookupTables proves, for every inventory and every value r in
// [0,total), that the lookup table returns the index a binary search over the
// independently computed cumulative weights returns (the smallest i with
// cumulative[i] > r), that every form is reachable, and that the form count fits
// the uint8 table entries.
func TestInventoryLookupTables(t *testing.T) {
	for _, it := range allWeightedInventories {
		t.Run(it.name, func(t *testing.T) {
			if n := len(it.inv.forms); n > maxInventoryForms {
				t.Fatalf("%d forms exceed the uint8 lookup-table limit %d", n, maxInventoryForms)
			}
			verifyInventorySampling(t, it.name, it.inv)
			cumulative := cumulativeWeights(t, it.name, it.inv)
			for r := uint32(0); r < it.inv.total; r++ {
				want, _ := slices.BinarySearch(cumulative, r+1)
				if got := it.inv.sampleIndex(r); got != want {
					t.Fatalf("sampleIndex(%d) = %d, binary search over cumulative weights = %d", r, got, want)
				}
			}
		})
	}
}

// TestOnsetNucleusLicenseMatchesSourceSets proves that the nucleus license
// recorded for every form of every inventory is exactly the one the source sets
// onsetFrontVowelOnly and onsetBackVowelOnly define, that the source sets hold
// exactly qu, gu and ç, and that each license selects the matching nucleus
// inventory.
func TestOnsetNucleusLicenseMatchesSourceSets(t *testing.T) {
	for _, it := range allWeightedInventories {
		if len(it.inv.license) != len(it.inv.forms) {
			t.Fatalf("%s: %d licenses for %d forms", it.name, len(it.inv.license), len(it.inv.forms))
		}
		for i := range it.inv.forms {
			form := it.inv.forms[i].form
			_, front := onsetFrontVowelOnly[form]
			_, back := onsetBackVowelOnly[form]
			want := licenseAnyNucleus
			switch {
			case front && back:
				t.Fatalf("%s: onset %q is in both the front-only and back-only sets", it.name, form)
			case front:
				want = licenseFrontNucleus
			case back:
				want = licenseBackNucleus
			}
			if got := it.inv.license[i]; got != want {
				t.Errorf("%s: license of %q = %d, want %d", it.name, form, got, want)
			}
		}
	}

	constrained := make([]string, 0, len(onsetFrontVowelOnly)+len(onsetBackVowelOnly))
	for form := range onsetFrontVowelOnly {
		constrained = append(constrained, form)
	}
	for form := range onsetBackVowelOnly {
		constrained = append(constrained, form)
	}
	slices.Sort(constrained)
	if want := []string{"gu", "qu", "ç"}; !slices.Equal(constrained, want) {
		t.Errorf("constrained onsets = %q, want %q", constrained, want)
	}

	if licensedNuclei[licenseAnyNucleus] != &nuclei ||
		licensedNuclei[licenseFrontNucleus] != &nucleiFrontInv ||
		licensedNuclei[licenseBackNucleus] != &nucleiBackInv {
		t.Error("licensedNuclei does not map each license to its nucleus inventory")
	}
}

// TestSampleOnsetNucleusMatchesSetLookup proves, draw for draw, that
// sampleOnsetNucleus produces the same onset and nucleus, consuming the same
// random values, as drawing the onset and then selecting the nucleus inventory by
// looking the onset up in onsetFrontVowelOnly and onsetBackVowelOnly.
func TestSampleOnsetNucleusMatchesSetLookup(t *testing.T) {
	for _, onsetInv := range []*weightedInventory{&onsetInitialInv, &onsetMedialInv} {
		got := rand.New(rand.NewPCG(20260927, 41))
		ref := rand.New(rand.NewPCG(20260927, 41))
		for i := 0; i < 200000; i++ {
			onset, nucleus := sampleOnsetNucleus(got, onsetInv)

			wantOnset := sampleForm(ref, onsetInv)
			nucInv := &nuclei
			if _, front := onsetFrontVowelOnly[wantOnset]; front {
				nucInv = &nucleiFrontInv
			} else if _, back := onsetBackVowelOnly[wantOnset]; back {
				nucInv = &nucleiBackInv
			}
			wantNucleus := sampleForm(ref, nucInv)

			if onset != wantOnset || nucleus != wantNucleus {
				t.Fatalf("draw %d: got %q+%q, want %q+%q", i, onset, nucleus, wantOnset, wantNucleus)
			}
		}
	}
}

// fixedSource is a rand.Source that always returns the same value. A *rand.Rand
// over fixedSource(uint64(x) << 32) returns x from every Uint32 call, which lets a
// test feed a chosen 32-bit draw to skewedCharTarget.
type fixedSource uint64

func (s fixedSource) Uint64() uint64 { return uint64(s) }

// skewOffsetForDraw returns the offset skewedCharTarget selects above minChars 0
// for the given span when the random source yields the 32-bit draw x.
func skewOffsetForDraw(span int, x uint32) int {
	return skewedCharTarget(rand.New(fixedSource(uint64(x)<<32)), 0, span)
}

// TestSkewThresholdsMatchFormula proves that the generated threshold table
// reproduces skewOffsetFormula exactly for every span in [1, skewTableMaxSpan].
//
// The table path returns the number of thresholds of the span that the draw
// reaches, a step function that is constant between consecutive thresholds. The
// formula is non-decreasing in the draw (u grows with x, 1-u*denom shrinks, and
// the logarithm divided by the negative ln d grows), so it too can change value
// only at a threshold. The test therefore checks both functions for equality at
// both sides of every threshold (t-1 and t) and at the two ends of the draw range
// (0 and MaxUint32), and checks at every threshold that the formula steps up
// (formula(t-1) < k <= formula(t)), which pins each threshold to the smallest
// draw that reaches its offset. A deterministic sample of further draws per span
// checks the equality between thresholds as well.
func TestSkewThresholdsMatchFormula(t *testing.T) {
	if want := skewTableMaxSpan * (skewTableMaxSpan + 1) / 2; len(skewThresholds) != want {
		t.Fatalf("len(skewThresholds) = %d, want %d", len(skewThresholds), want)
	}
	check := func(span int, x uint32) {
		t.Helper()
		if got, want := skewOffsetForDraw(span, x), skewOffsetFormula(span, x); got != want {
			t.Fatalf("span %d, draw %d: table offset %d, formula offset %d", span, x, got, want)
		}
	}
	sample := rand.New(rand.NewPCG(41, 20260927))
	for span := 1; span <= skewTableMaxSpan; span++ {
		row := skewThresholds[span*(span-1)/2 : span*(span+1)/2]
		check(span, 0)
		check(span, math.MaxUint32)
		if got := skewOffsetFormula(span, 0); got != 0 {
			t.Errorf("span %d: formula(0) = %d, want 0", span, got)
		}
		if got := skewOffsetFormula(span, math.MaxUint32); got != span {
			t.Errorf("span %d: formula(MaxUint32) = %d, want %d", span, got, span)
		}
		for i, thr := range row {
			k := i + 1
			if i > 0 && thr < row[i-1] {
				t.Fatalf("span %d: threshold %d (%d) < threshold %d (%d)", span, k, thr, i, row[i-1])
			}
			if got := skewOffsetFormula(span, thr); got < k {
				t.Fatalf("span %d: formula(%d) = %d, want >= %d", span, thr, got, k)
			}
			check(span, thr)
			if thr == 0 {
				continue
			}
			if got := skewOffsetFormula(span, thr-1); got >= k {
				t.Fatalf("span %d: formula(%d) = %d, want < %d", span, thr-1, got, k)
			}
			if skewOffsetFormula(span, thr-1) > skewOffsetFormula(span, thr) {
				t.Fatalf("span %d: formula decreases between %d and %d", span, thr-1, thr)
			}
			check(span, thr-1)
		}
		for range 1 << 12 {
			check(span, sample.Uint32())
		}
	}
}

// TestSkewTableCoversEveryLengthWindow guards that every character window of
// charRangeOf has a span of at most skewTableMaxSpan, so the length machinery
// never reaches the floating-point fallback of skewedCharTarget.
func TestSkewTableCoversEveryLengthWindow(t *testing.T) {
	for _, l := range []LengthTypeWords{AnyLengthWord, SmallLengthWord, MediumLengthWords, BigLengthWords, 9} {
		minChars, maxChars := charRangeOf(l)
		if span := maxChars - minChars; span > skewTableMaxSpan {
			t.Errorf("charRangeOf(%d) span %d exceeds skewTableMaxSpan %d", l, span, skewTableMaxSpan)
		}
	}
}

// TestSkewedCharTargetWideSpanUsesFormula checks the fallback for spans above
// skewTableMaxSpan: the target is minChars plus the closed-form offset, and stays
// inside the window.
func TestSkewedCharTargetWideSpanUsesFormula(t *testing.T) {
	const minChars = 3
	draws := []uint32{0, 1, 1 << 31, math.MaxUint32 - 1, math.MaxUint32}
	for _, span := range []int{skewTableMaxSpan + 1, 64, 1000} {
		for _, x := range draws {
			got := skewedCharTarget(rand.New(fixedSource(uint64(x)<<32)), minChars, minChars+span)
			if want := minChars + skewOffsetFormula(span, x); got != want {
				t.Errorf("span %d, draw %d: got %d, want %d", span, x, got, want)
			}
			if got < minChars || got > minChars+span {
				t.Errorf("span %d, draw %d: target %d outside [%d,%d]", span, x, got, minChars, minChars+span)
			}
		}
	}
}

// TestSkewTableIsGenerated guards that wordspt_skew_table.go carries the standard
// generated-code header, so tools and reviewers treat it as generated output.
func TestSkewTableIsGenerated(t *testing.T) {
	src, err := os.ReadFile("wordspt_skew_table.go")
	if err != nil {
		t.Fatal(err)
	}
	const header = "// Code generated by go run wordspt_skew_gen.go; DO NOT EDIT.\n"
	if !strings.HasPrefix(string(src), header) {
		t.Errorf("wordspt_skew_table.go does not start with %q", header)
	}
}
