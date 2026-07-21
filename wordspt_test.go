package gengo

import (
	"testing"
	"unicode/utf8"
)

// The tests in this file exercise the WordsPT phonotactic and accentuation
// conformance validators (the test oracle) with curated well-formed and
// ill-formed pt-PT data, and they verify that the weighted inventories support
// correct-by-construction sampling with no rejection loop.

// syl is a compact constructor for a decomposed syllable, used by the tables.
func syl(onset, nucleus, coda string) syllable {
	return syllable{onset: onset, nucleus: nucleus, coda: coda}
}

func TestOnsetIsValid(t *testing.T) {
	valid := []string{
		"", // optional onset
		"p", "b", "t", "d", "c", "ç", "g", "f", "v", "s", "z", "j", "l", "m", "n", "r", "x",
		"ch", "lh", "nh", "qu", "gu",
		"pr", "br", "tr", "dr", "cr", "gr", "fr", "vr", "pl", "bl", "cl", "gl", "fl",
	}
	for _, o := range valid {
		if !onsetIsValid(o) {
			t.Errorf("onsetIsValid(%q) = false, want true", o)
		}
	}
	invalid := []string{
		"tl", "dl", "sr", "vl", "zr", // clusters not in the pt-PT inventory
		"k", "w", "y", // letters absent from the pt-PT onset inventory
		"rr", "ss", // medial digraphs, never onsets
		"pt", "gn", "ps", // non-pt-PT onset clusters
		"prr", "abc",
	}
	for _, o := range invalid {
		if onsetIsValid(o) {
			t.Errorf("onsetIsValid(%q) = true, want false", o)
		}
	}
}

func TestNucleusIsValid(t *testing.T) {
	valid := []string{
		"a", "e", "i", "o", "u",
		"ei", "ai", "ou", "au", "oi", "eu", "iu", "ui",
		"ã", "õ", "ão", "ãe", "õe",
	}
	for _, n := range valid {
		if !nucleusIsValid(n) {
			t.Errorf("nucleusIsValid(%q) = false, want true", n)
		}
	}
	invalid := []string{
		"",               // the nucleus is mandatory
		"ae", "ao", "oe", // oral versions of the nasal diphthongs are not nuclei here
		"b", "s", "x", // consonants
		"ia", "ua", // rising sequences not in the falling-diphthong inventory
		"aa", "xyz",
	}
	for _, n := range invalid {
		if nucleusIsValid(n) {
			t.Errorf("nucleusIsValid(%q) = true, want false", n)
		}
	}
}

func TestCodaIsValid(t *testing.T) {
	valid := []string{"", "s", "r", "l", "m", "n", "z", "x"}
	for _, c := range valid {
		if !codaIsValid(c) {
			t.Errorf("codaIsValid(%q) = false, want true", c)
		}
	}
	invalid := []string{"p", "b", "t", "d", "c", "g", "f", "v", "j", "k", "ss", "rs"}
	for _, c := range invalid {
		if codaIsValid(c) {
			t.Errorf("codaIsValid(%q) = true, want false", c)
		}
	}
}

func TestSyllableIsValid(t *testing.T) {
	tests := []struct {
		name        string
		s           syllable
		wordInitial bool
		want        bool
	}{
		{"simple open", syl("t", "a", ""), true, true},
		{"onset cluster", syl("pr", "a", ""), true, true},
		{"closed with coda", syl("c", "a", "m"), true, true},
		{"no onset", syl("", "a", ""), true, true},
		{"nasal diphthong nucleus", syl("ç", "ão", ""), false, true},

		{"lh word-initial rejected", syl("lh", "o", ""), true, false},
		{"lh medial accepted", syl("lh", "o", ""), false, true},
		{"nh word-initial rejected", syl("nh", "a", ""), true, false},
		{"nh medial accepted", syl("nh", "a", ""), false, true},

		{"ç word-initial rejected", syl("ç", "a", ""), true, false},
		{"ç before back vowel accepted", syl("ç", "a", ""), false, true},
		{"ç before front vowel rejected", syl("ç", "e", ""), false, false},
		{"ç before i rejected", syl("ç", "i", ""), false, false},

		{"qu before front vowel accepted", syl("qu", "e", ""), true, true},
		{"qu before i accepted", syl("qu", "i", ""), false, true},
		{"qu before back vowel rejected", syl("qu", "a", ""), true, false},
		{"gu before front vowel accepted", syl("gu", "e", ""), false, true},
		{"gu before back vowel rejected", syl("gu", "o", ""), false, false},

		{"bad onset cluster", syl("tl", "a", ""), true, false},
		{"illegal coda", syl("c", "a", "p"), true, false},
		{"invalid nucleus", syl("c", "b", ""), true, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := syllableIsValid(tc.s, tc.wordInitial); got != tc.want {
				t.Errorf("syllableIsValid(%+v, %v) = %v, want %v", tc.s, tc.wordInitial, got, tc.want)
			}
		})
	}
}

func TestTransitionIsValid(t *testing.T) {
	tests := []struct {
		name      string
		coda      string
		nextOnset string
		want      bool
	}{
		{"empty coda", "", "t", true},
		{"m before p", "m", "p", true},
		{"m before b", "m", "b", true},
		{"m before t rejected", "m", "t", false},
		{"m before vowel rejected", "m", "", false},
		{"n before t", "n", "t", true},
		{"n before p rejected", "n", "p", false},
		{"n before b rejected", "n", "b", false},
		{"s before consonant", "s", "t", true},
		{"s before cluster", "s", "tr", true},
		{"r before consonant", "r", "m", true},
		{"coda before vowel rejected (MOP)", "s", "", false},
		{"l before vowel rejected (MOP)", "l", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := transitionIsValid(tc.coda, tc.nextOnset); got != tc.want {
				t.Errorf("transitionIsValid(%q, %q) = %v, want %v", tc.coda, tc.nextOnset, got, tc.want)
			}
		})
	}
}

func TestStemConformsToPhonotactics(t *testing.T) {
	conformant := []struct {
		name string
		stem []syllable
	}{
		{"prato", []syllable{syl("pr", "a", ""), syl("t", "o", "")}},
		{"campo (m before p)", []syllable{syl("c", "a", "m"), syl("p", "o", "")}},
		{"canto (n before t)", []syllable{syl("c", "a", "n"), syl("t", "o", "")}},
		{"olho (lh medial)", []syllable{syl("", "o", ""), syl("lh", "o", "")}},
		{"casaco", []syllable{syl("c", "a", ""), syl("s", "a", ""), syl("c", "o", "")}},
		{"coração (ç medial + nasal)", []syllable{syl("c", "o", ""), syl("r", "a", ""), syl("ç", "ão", "")}},
	}
	for _, tc := range conformant {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			if !stemConformsToPhonotactics(tc.stem) {
				t.Errorf("stemConformsToPhonotactics(%q) = false, want true", tc.name)
			}
		})
	}

	nonConformant := []struct {
		name string
		stem []syllable
	}{
		{"empty stem", []syllable{}},
		{"m before t", []syllable{syl("c", "a", "m"), syl("t", "o", "")}},
		{"n before p", []syllable{syl("c", "a", "n"), syl("p", "o", "")}},
		{"coda before vowel (MOP)", []syllable{syl("c", "a", "s"), syl("", "a", "")}},
		{"bad onset cluster", []syllable{syl("tl", "a", ""), syl("t", "o", "")}},
		{"illegal coda", []syllable{syl("c", "a", "p"), syl("t", "o", "")}},
		{"lh word-initial", []syllable{syl("lh", "a", ""), syl("t", "o", "")}},
		{"ç word-initial", []syllable{syl("ç", "a", ""), syl("t", "o", "")}},
		{"invalid nucleus", []syllable{syl("c", "xyz", "")}},
	}
	for _, tc := range nonConformant {
		t.Run("reject/"+tc.name, func(t *testing.T) {
			if stemConformsToPhonotactics(tc.stem) {
				t.Errorf("stemConformsToPhonotactics(%q) = true, want false", tc.name)
			}
		})
	}
}

func TestStressClassOf(t *testing.T) {
	tests := []struct {
		count, tonic int
		want         stressClass
	}{
		{1, 0, oxytone},
		{2, 1, oxytone},
		{2, 0, paroxytone},
		{3, 2, oxytone},
		{3, 1, paroxytone},
		{3, 0, proparoxytone},
		{4, 1, proparoxytone},
	}
	for _, tc := range tests {
		if got := stressClassOf(tc.count, tc.tonic); got != tc.want {
			t.Errorf("stressClassOf(%d, %d) = %d, want %d", tc.count, tc.tonic, got, tc.want)
		}
	}
}

func TestStemConformsToAccentuation(t *testing.T) {
	conformant := []struct {
		name  string
		stem  []syllable
		tonic int
	}{
		// Proparoxytones: always accented.
		{"médico", []syllable{syl("m", "é", ""), syl("d", "i", ""), syl("c", "o", "")}, 0},
		{"câmara", []syllable{syl("c", "â", ""), syl("m", "a", ""), syl("r", "a", "")}, 0},
		// Oxytones ending in tonic a/e/o (+ optional s) or -em.
		{"café", []syllable{syl("c", "a", ""), syl("f", "é", "")}, 1},
		{"sofá", []syllable{syl("s", "o", ""), syl("f", "á", "")}, 1},
		{"também (-em)", []syllable{syl("t", "a", "m"), syl("b", "é", "m")}, 1},
		// Oxytones NOT requiring an accent (default: consonant coda / i / u).
		{"animal", []syllable{syl("", "a", ""), syl("n", "i", ""), syl("m", "a", "l")}, 2},
		{"feliz", []syllable{syl("f", "e", ""), syl("l", "i", "z")}, 1},
		// Paroxytones: default endings unaccented.
		{"casa", []syllable{syl("c", "a", ""), syl("s", "a", "")}, 0},
		{"casas", []syllable{syl("c", "a", ""), syl("s", "a", "s")}, 0},
		{"homem (-em default)", []syllable{syl("", "o", ""), syl("m", "e", "m")}, 0},
		// Paroxytones requiring an accent (non-default endings).
		{"fácil (-l)", []syllable{syl("f", "á", ""), syl("c", "i", "l")}, 0},
		{"lápis (i+s)", []syllable{syl("l", "á", ""), syl("p", "i", "s")}, 0},
		{"órfã (nasal end)", []syllable{syl("", "ó", "r"), syl("f", "ã", "")}, 0},
		// Nasal tonic nuclei: marked by the tilde alone.
		{"irmã", []syllable{syl("", "i", "r"), syl("m", "ã", "")}, 1},
		{"coração", []syllable{syl("c", "o", ""), syl("r", "a", ""), syl("ç", "ão", "")}, 2},
	}
	for _, tc := range conformant {
		t.Run("accept/"+tc.name, func(t *testing.T) {
			if !stemConformsToAccentuation(tc.stem, tc.tonic) {
				t.Errorf("stemConformsToAccentuation(%q, %d) = false, want true", tc.name, tc.tonic)
			}
		})
	}

	nonConformant := []struct {
		name  string
		stem  []syllable
		tonic int
	}{
		// Missing required accents.
		{"medico (no accent, proparoxytone)", []syllable{syl("m", "e", ""), syl("d", "i", ""), syl("c", "o", "")}, 0},
		{"cafe (no accent, oxytone)", []syllable{syl("c", "a", ""), syl("f", "e", "")}, 1},
		{"tambem (no accent, -em oxytone)", []syllable{syl("t", "a", "m"), syl("b", "e", "m")}, 1},
		{"facil (no accent, paroxytone -l)", []syllable{syl("f", "a", ""), syl("c", "i", "l")}, 0},
		// Spurious accents where none is required.
		{"cása (spurious on tonic paroxytone)", []syllable{syl("c", "á", ""), syl("s", "a", "")}, 0},
		{"anìmal-like spurious", []syllable{syl("", "a", ""), syl("n", "í", ""), syl("m", "a", "l")}, 2},
		// Accent on the wrong (non-tonic) syllable.
		{"accent on non-tonic", []syllable{syl("c", "á", ""), syl("s", "a", "")}, 1},
		// Nasal tonic must not carry acute/circumflex (constructed illegal grapheme aside,
		// an oral acute on a non-tonic together with a nasal tonic is rejected on placement).
		{"acute on non-tonic with nasal tonic", []syllable{syl("", "í", "r"), syl("m", "ã", "")}, 1},
		// Out-of-range tonic index.
		{"tonic out of range", []syllable{syl("c", "a", ""), syl("s", "a", "")}, 5},
		{"empty stem", []syllable{}, 0},
	}
	for _, tc := range nonConformant {
		t.Run("reject/"+tc.name, func(t *testing.T) {
			if stemConformsToAccentuation(tc.stem, tc.tonic) {
				t.Errorf("stemConformsToAccentuation(%q, %d) = true, want false", tc.name, tc.tonic)
			}
		})
	}
}

func TestWordConformsToCedilla(t *testing.T) {
	conformant := []string{
		"casa",    // no cedilla
		"caça",    // ç before a
		"moço",    // ç before o
		"açúcar",  // ç before ú (base u)
		"maçã",    // ç before ã (base a)
		"coração", // ç before ão
		"caçula",  // ç before u
	}
	for _, w := range conformant {
		if !wordConformsToCedilla(w) {
			t.Errorf("wordConformsToCedilla(%q) = false, want true", w)
		}
	}
	nonConformant := []string{
		"çasa", // word-initial cedilla
		"çe",   // word-initial cedilla before e
		"laçe", // ç before e
		"laçi", // ç before i
		"laçé", // ç before é (base e)
		"laç",  // word-final cedilla
		"laçt", // ç before a consonant
	}
	for _, w := range nonConformant {
		if wordConformsToCedilla(w) {
			t.Errorf("wordConformsToCedilla(%q) = true, want false", w)
		}
	}
}

func TestWordHasForbiddenDiaeresis(t *testing.T) {
	clean := []string{"pinguim", "aguentar", "cinquenta", "casa"}
	for _, w := range clean {
		if wordHasForbiddenDiaeresis(w) {
			t.Errorf("wordHasForbiddenDiaeresis(%q) = true, want false", w)
		}
	}
	withTrema := []string{"pingüim", "agüentar", "cinqüenta"} // pre-1990 spellings
	for _, w := range withTrema {
		if !wordHasForbiddenDiaeresis(w) {
			t.Errorf("wordHasForbiddenDiaeresis(%q) = false, want true", w)
		}
	}
}

func TestWordContainsGrave(t *testing.T) {
	without := []string{"casa", "café", "coração", "animal"}
	for _, w := range without {
		if wordContainsGrave(w) {
			t.Errorf("wordContainsGrave(%q) = true, want false", w)
		}
	}
	with := []string{"à", "às", "àquele", "àquela"}
	for _, w := range with {
		if !wordContainsGrave(w) {
			t.Errorf("wordContainsGrave(%q) = false, want true", w)
		}
	}
}

// TestSampleInventoriesNoRejection verifies that every weighted inventory
// supports correct-by-construction sampling: for every r in [0,total) the index
// returned partitions the cumulative table correctly, and every form is
// reachable. No value of r is ever rejected. This is the smoke test required by
// the acceptance criteria.
func TestSampleInventoriesNoRejection(t *testing.T) {
	inventories := []struct {
		name string
		inv  *weightedInventory
	}{
		{"onsetSingles", &onsetSingles},
		{"onsetClusters", &onsetClusters},
		{"nuclei", &nuclei},
		{"codas", &codas},
	}
	for _, it := range inventories {
		t.Run(it.name, func(t *testing.T) {
			verifyInventorySampling(t, it.name, it.inv)
		})
	}
}

// verifyCumulativeTable checks that the cumulative-weight table is a strictly
// increasing prefix sum of positive weights matching the recorded total.
func verifyCumulativeTable(t *testing.T, name string, inv *weightedInventory) {
	t.Helper()
	if len(inv.forms) == 0 {
		t.Fatalf("%s is empty", name)
	}
	var sum uint32
	for i := range inv.forms {
		if inv.forms[i].weight == 0 {
			t.Errorf("%s[%d] %q has zero weight", name, i, inv.forms[i].form)
		}
		sum += inv.forms[i].weight
		if inv.cumulative[i] != sum {
			t.Errorf("%s cumulative[%d] = %d, want %d", name, i, inv.cumulative[i], sum)
		}
	}
	if sum != inv.total {
		t.Fatalf("%s total = %d, want %d", name, inv.total, sum)
	}
}

// verifySampleIndex checks, for a single r, that sampleIndex returns the index
// whose cumulative range contains r, proving the mapping never rejects.
func verifySampleIndex(t *testing.T, name string, inv *weightedInventory, r uint32) {
	t.Helper()
	idx := inv.sampleIndex(r)
	if idx < 0 || idx >= len(inv.forms) {
		t.Fatalf("%s sampleIndex(%d) = %d out of range [0,%d)", name, r, idx, len(inv.forms))
	}
	if inv.cumulative[idx] <= r {
		t.Fatalf("%s sampleIndex(%d) = %d but cumulative[%d]=%d <= r", name, r, idx, idx, inv.cumulative[idx])
	}
	if idx > 0 && inv.cumulative[idx-1] > r {
		t.Fatalf("%s sampleIndex(%d) = %d but cumulative[%d]=%d > r", name, r, idx, idx-1, inv.cumulative[idx-1])
	}
}

// verifyInventorySampling runs the full no-rejection and reachability checks for
// one inventory.
func verifyInventorySampling(t *testing.T, name string, inv *weightedInventory) {
	t.Helper()
	if inv.total == 0 {
		t.Fatalf("%s total weight is 0", name)
	}
	verifyCumulativeTable(t, name, inv)

	// Full sweep: every r in [0,total) maps to a correct index with no rejection.
	for r := uint32(0); r < inv.total; r++ {
		verifySampleIndex(t, name, inv, r)
	}

	// Reachability: every index is the selection for at least one r.
	for i := range inv.forms {
		r := inv.cumulative[i] - 1
		if got := inv.sampleIndex(r); got != i {
			t.Errorf("%s form %q (index %d) unreachable: sampleIndex(%d) = %d", name, inv.forms[i].form, i, r, got)
		}
	}
}

// TestInventoryFormsAreValidUTF8 guards that every embedded form is valid UTF-8,
// so the rune-oriented generation layers can rely on it.
func TestInventoryFormsAreValidUTF8(t *testing.T) {
	all := []weightedInventory{onsetSingles, onsetClusters, nuclei, codas}
	for _, inv := range all {
		for i := range inv.forms {
			if !utf8.ValidString(inv.forms[i].form) {
				t.Errorf("form %q is not valid UTF-8", inv.forms[i].form)
			}
		}
	}
}
