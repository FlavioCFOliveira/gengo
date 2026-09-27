# Test Report - gengo

**Date:** 2026-09-27
**Package:** github.com/FlavioCFOliveira/gengo
**Version:** v0.2.1
**Platform:** linux/amd64 (AMD Ryzen 9 5900HX, 16 threads) · Go go1.27.1
**Test command:** `go test -race -count=1 -cover -v .`
**Benchmark command:** see [Performance Tests (Benchmarks)](#performance-tests-benchmarks)

---

## Executive Summary

| Metric | Value |
|--------|-------|
| **Unit Tests** | 238/238 passed (100%) |
| **Subtests** | 178/178 passed (100%) |
| **Examples** | 7/7 passed (100%) |
| **Race Detector** | No data races reported |
| **Statement Coverage** | 96.5% |
| **Benchmarks** | 36 executed |
| **Total Test Time** | 85.955s (with `-race`) |
| **Total Benchmark Time** | 233.688s |

---

## Unit Tests

All 238 unit tests, 178 subtests, and 7 examples passed, with no failures and no skipped tests. Durations were measured with the race detector enabled.

### Detailed Results

| Test | Status | Duration |
|------|--------|----------|
| `TestBool` | PASS | 0.00s |
| `TestDate` | PASS | 0.00s |
| `TestUnixDate` | PASS | 0.00s |
| `TestDateBetween` | PASS | 0.00s |
| `TestDateBetweenGranularity` | PASS | 0.00s |
| `TestGeneratorReproducible` | PASS | 0.02s |
| `TestNewSourceReproducible` | PASS | 0.00s |
| `TestGeneratorRanges` | PASS | 0.00s |
| `TestGeneratorIntSwapBranches` | PASS | 0.00s |
| `TestGeneratorUintSwapAndFullRange` | PASS | 0.00s |
| `TestGeneratorFloatBranches` | PASS | 0.00s |
| `TestGeneratorStringEdges` | PASS | 0.14s |
| `TestGeneratorDateAndWordEdges` | PASS | 0.00s |
| `TestGeneratorStringAllocs` | PASS | 0.18s |
| `TestGeneratorStringGolden` | PASS | 0.01s |
| `TestInt8` | PASS | 0.00s |
| `TestInt8Between` | PASS | 0.00s |
| `TestInt16` | PASS | 0.00s |
| `TestInt16Between` | PASS | 0.00s |
| `TestInt32` | PASS | 0.00s |
| `TestInt32Between` | PASS | 0.00s |
| `TestInt` | PASS | 0.00s |
| `TestIntBetween` | PASS | 0.00s |
| `TestInt64` | PASS | 0.00s |
| `TestInt64Between` | PASS | 0.00s |
| `TestUint8` | PASS | 0.00s |
| `TestUint8Between` | PASS | 0.00s |
| `TestByte` | PASS | 0.00s |
| `TestUint16` | PASS | 0.00s |
| `TestUint16Between` | PASS | 0.00s |
| `TestUint32` | PASS | 0.00s |
| `TestUint32Between` | PASS | 0.00s |
| `TestUint64` | PASS | 0.00s |
| `TestUint64Between` | PASS | 0.00s |
| `TestFloat32` | PASS | 0.00s |
| `TestFloat32Between` | PASS | 0.00s |
| `TestFloat64` | PASS | 0.00s |
| `TestFloat64Between` | PASS | 0.00s |
| `TestComplex64` | PASS | 0.00s |
| `TestComplex64Between` | PASS | 0.00s |
| `TestComplex128` | PASS | 0.00s |
| `TestComplex128Between` | PASS | 0.00s |
| `TestString` | PASS | 0.27s |
| `TestStringByteOriented` | PASS | 0.01s |
| `TestStringAllChars` | PASS | 0.01s |
| `TestStringAlphanumeric` | PASS | 0.01s |
| `TestStringAlphabetic` | PASS | 0.01s |
| `TestStringAlphabeticUppercase` | PASS | 0.01s |
| `TestStringAlphabeticLowercase` | PASS | 0.01s |
| `TestStringNumeric` | PASS | 0.01s |
| `TestStringHexadecimal` | PASS | 0.01s |
| `TestStringSymbols` | PASS | 0.01s |
| `TestStringBetween` | PASS | 0.01s |
| `TestStringAllocs` | PASS | 0.18s |
| `TestStringSingleCharLong` | PASS | 0.00s |
| `TestWord` | PASS | 0.00s |
| `TestWordByLengthType` | PASS | 0.01s |
| `TestWords` | PASS | 0.00s |
| `TestAccentuateStemConformsToAccentuation` | PASS | 0.28s |
| `TestAccentuateStemKnownCases` | PASS | 0.00s |
| `TestAccentuateStemDeterministic` | PASS | 0.02s |
| `TestAccentedWordNoRandomDraws` | PASS | 0.00s |
| `TestAssembleAccentedStemSingleAllocation` | PASS | 0.00s |
| `TestAccentedWordAllocationBudget` | PASS | 0.01s |
| `TestAccentuateStemEmptyIsSafe` | PASS | 0.00s |
| `TestAdjectivePTOrthographicConformance` | PASS | 4.93s |
| `TestAdjectivePTGenderHonored` | PASS | 1.88s |
| `TestAdjectivePTNumberHonored` | PASS | 1.88s |
| `TestAdjectivePTDegreeHonored` | PASS | 0.97s |
| `TestAdjectivePTSuperlativeInflection` | PASS | 1.73s |
| `TestAdjectivePTSuperlativeOrthography` | PASS | 0.76s |
| `TestAdjectivePTThematicOAProduced` | PASS | 0.96s |
| `TestAdjectivePTSuperlativeHardeningViaGeneration` | PASS | 0.00s |
| `TestAdjectiveSuperlativeHardeningRule` | PASS | 0.00s |
| `TestAdjectivePTLengthHonored` | PASS | 4.73s |
| `TestAdjectivePTAnyGenderVaried` | PASS | 0.48s |
| `TestAdjectivePTAnyNumberVaried` | PASS | 0.50s |
| `TestAdjectivePTAnyDegreeVaried` | PASS | 0.49s |
| `TestAdjectivePTAnyLengthVaried` | PASS | 0.48s |
| `TestAdjectivePTEndingVariety` | PASS | 0.97s |
| `TestAdjectivePTReproducible` | PASS | 1.47s |
| `TestAdjectiveEndingLengthAwareSelection` | PASS | 0.09s |
| `TestAdjectiveStemConformsToOracles` | PASS | 0.86s |
| `TestAdjectivePTConcurrentSafe` | PASS | 0.60s |
| `TestAdverbPTOrthographicConformance` | PASS | 2.07s |
| `TestAdverbPTEndsInMente` | PASS | 1.43s |
| `TestAdverbPTNoGraphicAccent` | PASS | 1.86s |
| `TestAdverbPTInvariable` | PASS | 0.38s |
| `TestAdverbWordFormation` | PASS | 0.00s |
| `TestAdverbPTLengthHonored` | PASS | 1.44s |
| `TestAdverbPTSmallNormalizesUp` | PASS | 0.34s |
| `TestAdverbPTEndingVariety` | PASS | 0.38s |
| `TestAdverbEndingLengthAwareSelection` | PASS | 0.06s |
| `TestAdverbPTReproducible` | PASS | 1.04s |
| `TestAdverbPTConcurrentSafe` | PASS | 0.59s |
| `TestClosedClassListsCuration` | PASS | 0.00s |
| `TestClosedClassListsMatchReference` | PASS | 0.00s |
| `TestClosedClassMembershipAndReachability` | PASS | 0.30s |
| `TestClosedClassPackageLevelMembership` | PASS | 0.34s |
| `TestClosedClassReproducibility` | PASS | 0.01s |
| `TestClosedClassNoAllocation` | PASS | 0.01s |
| `TestPronounPTSubtypeCoverage` | PASS | 0.00s |
| `TestPronounPTIncludesVosForms` | PASS | 0.00s |
| `TestNumeralPTSubtypeCoverage` | PASS | 0.00s |
| `TestNumeralPTPortugueseSpelling` | PASS | 0.00s |
| `TestFlexionEnumZeroValues` | PASS | 0.00s |
| `TestResolveConcretePassthrough` | PASS | 0.00s |
| `TestResolveAnyUniformAndValid` | PASS | 0.24s |
| `TestResolveReproducible` | PASS | 0.00s |
| `TestCharRangeOf` | PASS | 0.00s |
| `TestNormalizeLength` | PASS | 0.00s |
| `TestSyllablesForCharRange` | PASS | 0.31s |
| `TestStemSyllablesForNormalizesAndSizes` | PASS | 0.54s |
| `TestStemSyllablesForLandsInWindow` | PASS | 1.83s |
| `TestStemHiatusReduced` | PASS | 0.26s |
| `TestNounPTHiatusReduced` | PASS | 0.94s |
| `TestBigNounLengthRealistic` | PASS | 0.99s |
| `TestMedialOnsetPrefersConsonant` | PASS | 0.00s |
| `TestNounPTOrthographicConformance` | PASS | 4.86s |
| `TestNounPTGenderHonored` | PASS | 0.82s |
| `TestNounPTNumberHonored` | PASS | 0.85s |
| `TestNounPTLengthHonored` | PASS | 2.67s |
| `TestNounPTAnyGenderVaried` | PASS | 0.41s |
| `TestNounPTAnyNumberVaried` | PASS | 0.43s |
| `TestNounPTAnyLengthVaried` | PASS | 0.41s |
| `TestNounPTReproducible` | PASS | 1.17s |
| `TestNounEndingLengthAwareSelection` | PASS | 0.06s |
| `TestNounPTConcurrentSafe` | PASS | 0.44s |
| `TestNounStemConformsToOracles` | PASS | 0.85s |
| `TestAdditivePluralClassificationMatchesPluralize` | PASS | 1.39s |
| `TestAccentedWordSuffixedMatchesPluralize` | PASS | 1.55s |
| `TestPluralBuildAllocationBudget` | PASS | 0.03s |
| `TestAdditivePluralSuffixClassification` | PASS | 0.00s |
| `TestPluralizeRuleTable` | PASS | 0.00s |
| `TestApplyAoPlural` | PASS | 0.00s |
| `TestPluralizeAoIsWeightedAndReachable` | PASS | 0.10s |
| `TestSampleAoPluralDistribution` | PASS | 0.08s |
| `TestPluralizeEmpty` | PASS | 0.00s |
| `TestPluralizePreservesStringOrthography` | PASS | 0.23s |
| `TestPluralizeReproducible` | PASS | 0.03s |
| `TestPluralizeAllocations` | PASS | 0.00s |
| `TestAllWeightedInventoriesListed` | PASS | 0.00s |
| `TestInventoryLookupTables` | PASS | 0.01s |
| `TestOnsetNucleusLicenseMatchesSourceSets` | PASS | 0.00s |
| `TestSampleOnsetNucleusMatchesSetLookup` | PASS | 0.72s |
| `TestSkewThresholdsMatchFormula` | PASS | 0.20s |
| `TestSkewTableCoversEveryLengthWindow` | PASS | 0.00s |
| `TestSkewedCharTargetWideSpanUsesFormula` | PASS | 0.00s |
| `TestSkewTableIsGenerated` | PASS | 0.00s |
| `TestSampleSyllabicStemConformsToPhonotactics` | PASS | 0.20s |
| `TestSampleSyllabicStemHonoursTonic` | PASS | 0.00s |
| `TestSampleSyllabicStemReproducible` | PASS | 0.02s |
| `TestSampleSyllabicStemBufferReuse` | PASS | 0.00s |
| `TestDerivedInventoriesNoRejection` | PASS | 0.00s |
| `TestOnsetInventoryMembership` | PASS | 0.00s |
| `TestNucleiInventoryPartition` | PASS | 0.00s |
| `TestCodaInventoryMembership` | PASS | 0.00s |
| `TestAssembleStemSingleAllocation` | PASS | 0.00s |
| `TestAssembleStemContent` | PASS | 0.00s |
| `TestFullStemAllocationBudget` | PASS | 0.01s |
| `TestNoRejectionDeterministicDraws` | PASS | 0.00s |
| `TestOnsetIsValid` | PASS | 0.00s |
| `TestNucleusIsValid` | PASS | 0.00s |
| `TestCodaIsValid` | PASS | 0.00s |
| `TestSyllableIsValid` | PASS | 0.00s |
| `TestTransitionIsValid` | PASS | 0.00s |
| `TestStemConformsToPhonotactics` | PASS | 0.00s |
| `TestStressClassOf` | PASS | 0.00s |
| `TestStemConformsToAccentuation` | PASS | 0.00s |
| `TestWordConformsToCedilla` | PASS | 0.00s |
| `TestWordHasForbiddenDiaeresis` | PASS | 0.00s |
| `TestWordContainsGrave` | PASS | 0.00s |
| `TestSampleInventoriesNoRejection` | PASS | 0.00s |
| `TestInventoryFormsAreValidUTF8` | PASS | 0.00s |
| `TestVerbPTOrthographicConformance` | PASS | 5.43s |
| `TestVerbDesinenceMatchesAssemblers` | PASS | 0.00s |
| `TestVerbPTSlotHonored` | PASS | 0.60s |
| `TestVerbPTPresentTenseHonored` | PASS | 0.52s |
| `TestVerbPTAnyMoodVaried` | PASS | 0.12s |
| `TestResolveMoodPassthrough` | PASS | 0.00s |
| `TestVerbPTIndicativeTenseVaried` | PASS | 0.26s |
| `TestVerbPTPersonNumberVaried` | PASS | 0.26s |
| `TestResolveInfinitiveForm` | PASS | 0.13s |
| `TestVerbPTInfinitiveConcretePersonIsPersonal` | PASS | 0.51s |
| `TestVerbPTPresentNeedsNoAccent` | PASS | 0.79s |
| `TestVerbPTNormalizationN1` | PASS | 0.82s |
| `TestVerbPTNormalizationN2` | PASS | 0.20s |
| `TestVerbPTNormalizationN3Subjunctive` | PASS | 0.46s |
| `TestVerbPTNormalizationN3TenseIgnored` | PASS | 1.21s |
| `TestVerbPTLengthHonored` | PASS | 2.66s |
| `TestVerbPTAnyLengthVaried` | PASS | 0.26s |
| `TestVerbPTReproducible` | PASS | 0.96s |
| `TestVerbPTSingleAllocation` | PASS | 0.10s |
| `TestVerbPTConcurrentSafe` | PASS | 0.27s |
| `TestVerbRadicalOnsetInventoryExcludesSensitive` | PASS | 0.00s |
| `TestNonFiniteExactForms` | PASS | 0.00s |
| `TestPersonalInfinitiveIgnoresPersonWhereApplicable` | PASS | 0.00s |
| `TestNonFiniteVerbOrthographicConformance` | PASS | 0.75s |
| `TestNonFiniteVerbLengthHonored` | PASS | 0.62s |
| `TestNonFiniteVerbReproducible` | PASS | 0.24s |
| `TestPickConjugationWeighted` | PASS | 0.03s |
| `TestResolveVerbPerson` | PASS | 0.02s |
| `TestVerbRadicalStemConformsToOracle` | PASS | 0.23s |
| `TestConjugateNonFiniteSingleAllocation` | PASS | 0.00s |
| `TestIndicativeExactForms` | PASS | 0.00s |
| `TestIndicativePreteriteVsPresentFirstPlural` | PASS | 0.00s |
| `TestConjugateIndicativeSingleAllocation` | PASS | 0.03s |
| `TestIndicativeOrthographicConformance` | PASS | 0.00s |
| `TestResolveIndicativeTense` | PASS | 0.01s |
| `TestIndicativeVerbFormConformance` | PASS | 0.31s |
| `TestIndicativeVerbFormReproducible` | PASS | 0.05s |
| `TestSubjunctiveExactForms` | PASS | 0.00s |
| `TestSubjunctiveImperfectFirstPlural` | PASS | 0.00s |
| `TestSubjunctiveFutureEqualsPersonalInfinitive` | PASS | 0.00s |
| `TestImperativeExactForms` | PASS | 0.00s |
| `TestImperativeDerivationRule` | PASS | 0.00s |
| `TestImperativeNormalizesFirstSingular` | PASS | 0.00s |
| `TestConjugateSubjunctiveSingleAllocation` | PASS | 0.02s |
| `TestConjugateImperativeSingleAllocation` | PASS | 0.01s |
| `TestSubjunctiveOrthographicConformance` | PASS | 0.00s |
| `TestImperativeOrthographicConformance` | PASS | 0.00s |
| `TestResolveSubjunctiveTense` | PASS | 0.01s |
| `TestResolveImperativePerson` | PASS | 0.04s |
| `TestSubjunctiveVerbFormConformance` | PASS | 0.32s |
| `TestImperativeVerbFormConformance` | PASS | 0.09s |
| `TestSubjunctiveVerbFormReproducible` | PASS | 0.05s |
| `TestImperativeVerbFormReproducible` | PASS | 0.01s |
| `TestResolveWordClassReachesAll` | PASS | 0.28s |
| `TestResolveWordClassDistribution` | PASS | 0.28s |
| `TestWordPTOutputValidity` | PASS | 2.35s |
| `TestWordPTClassMix` | PASS | 0.84s |
| `TestWordPTByLengthTypeBigOpenWindow` | PASS | 0.94s |
| `TestWordPTClosedIgnoreLength` | PASS | 0.00s |
| `TestWordsPTCount` | PASS | 0.08s |
| `TestWordsPTEmptyForNonPositive` | PASS | 0.00s |
| `TestWordPTReproducible` | PASS | 2.13s |
| `TestWordPTConcurrentSafe` | PASS | 1.40s |
| `TestWordPTAllocationBudget` | PASS | 0.22s |

### Subtests
- `TestWordByLengthType/SmallLengthWord` - PASS
- `TestWordByLengthType/MediumLengthWords` - PASS
- `TestWordByLengthType/BigLengthWords` - PASS
- `TestWordByLengthType/Default` - PASS
- `TestWords/PositiveLength` - PASS
- `TestWords/ZeroLength` - PASS
- `TestWords/NegativeLength` - PASS
- `TestAccentuateStemKnownCases/oxytone-a` - PASS
- `TestAccentuateStemKnownCases/oxytone-e` - PASS
- `TestAccentuateStemKnownCases/oxytone-o` - PASS
- `TestAccentuateStemKnownCases/oxytone-as` - PASS
- `TestAccentuateStemKnownCases/oxytone-em-acute` - PASS
- `TestAccentuateStemKnownCases/oxytone-im-none` - PASS
- `TestAccentuateStemKnownCases/oxytone-u-none` - PASS
- `TestAccentuateStemKnownCases/proparoxytone-a` - PASS
- `TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-onset` - PASS
- `TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-coda` - PASS
- `TestAccentuateStemKnownCases/proparoxytone-e-stays-acute-before-nasal` - PASS
- `TestAccentuateStemKnownCases/paroxytone-r-nondefault` - PASS
- `TestAccentuateStemKnownCases/paroxytone-diphthong-nondefault` - PASS
- `TestAccentuateStemKnownCases/paroxytone-default-a-none` - PASS
- `TestAccentuateStemKnownCases/paroxytone-default-as-none` - PASS
- `TestAccentuateStemKnownCases/paroxytone-nasal-ending-acute-tonic` - PASS
- `TestAccentuateStemKnownCases/nasal-vowel-tonic` - PASS
- `TestAccentuateStemKnownCases/nasal-diphthong-tonic-cedilla` - PASS
- `TestClosedClassListsCuration/article` - PASS
- `TestClosedClassListsCuration/preposition` - PASS
- `TestClosedClassListsCuration/conjunction` - PASS
- `TestClosedClassListsCuration/pronoun` - PASS
- `TestClosedClassListsCuration/interjection` - PASS
- `TestClosedClassListsCuration/numeral` - PASS
- `TestClosedClassListsMatchReference/article` - PASS
- `TestClosedClassListsMatchReference/preposition` - PASS
- `TestClosedClassListsMatchReference/conjunction` - PASS
- `TestClosedClassListsMatchReference/pronoun` - PASS
- `TestClosedClassListsMatchReference/interjection` - PASS
- `TestClosedClassListsMatchReference/numeral` - PASS
- `TestClosedClassMembershipAndReachability/article` - PASS
- `TestClosedClassMembershipAndReachability/preposition` - PASS
- `TestClosedClassMembershipAndReachability/conjunction` - PASS
- `TestClosedClassMembershipAndReachability/pronoun` - PASS
- `TestClosedClassMembershipAndReachability/interjection` - PASS
- `TestClosedClassMembershipAndReachability/numeral` - PASS
- `TestClosedClassReproducibility/article` - PASS
- `TestClosedClassReproducibility/preposition` - PASS
- `TestClosedClassReproducibility/conjunction` - PASS
- `TestClosedClassReproducibility/pronoun` - PASS
- `TestClosedClassReproducibility/interjection` - PASS
- `TestClosedClassReproducibility/numeral` - PASS
- `TestPluralizeRuleTable/casa` - PASS
- `TestPluralizeRuleTable/café` - PASS
- `TestPluralizeRuleTable/maçã` - PASS
- `TestPluralizeRuleTable/mãe` - PASS
- `TestPluralizeRuleTable/herói` - PASS
- `TestPluralizeRuleTable/gato` - PASS
- `TestPluralizeRuleTable/homem` - PASS
- `TestPluralizeRuleTable/jardim` - PASS
- `TestPluralizeRuleTable/bom` - PASS
- `TestPluralizeRuleTable/flor` - PASS
- `TestPluralizeRuleTable/amor` - PASS
- `TestPluralizeRuleTable/luz` - PASS
- `TestPluralizeRuleTable/rapaz` - PASS
- `TestPluralizeRuleTable/liquen` - PASS
- `TestPluralizeRuleTable/tórax` - PASS
- `TestPluralizeRuleTable/látex` - PASS
- `TestPluralizeRuleTable/país` - PASS
- `TestPluralizeRuleTable/lápis` - PASS
- `TestPluralizeRuleTable/vírus` - PASS
- `TestPluralizeRuleTable/animal` - PASS
- `TestPluralizeRuleTable/papel` - PASS
- `TestPluralizeRuleTable/anzol` - PASS
- `TestPluralizeRuleTable/azul` - PASS
- `TestPluralizeRuleTable/amável` - PASS
- `TestPluralizeRuleTable/possível` - PASS
- `TestPluralizeRuleTable/funil` - PASS
- `TestPluralizeRuleTable/barril` - PASS
- `TestPluralizeRuleTable/fácil` - PASS
- `TestPluralizeRuleTable/útil` - PASS
- `TestPluralizeRuleTable/difícil` - PASS
- `TestInventoryLookupTables/onsetSingles` - PASS
- `TestInventoryLookupTables/onsetClusters` - PASS
- `TestInventoryLookupTables/nuclei` - PASS
- `TestInventoryLookupTables/codas` - PASS
- `TestInventoryLookupTables/onsetInitialInv` - PASS
- `TestInventoryLookupTables/onsetMedialInv` - PASS
- `TestInventoryLookupTables/nucleiFrontInv` - PASS
- `TestInventoryLookupTables/nucleiBackInv` - PASS
- `TestInventoryLookupTables/codaBeforePBInv` - PASS
- `TestInventoryLookupTables/codaBeforeOtherInv` - PASS
- `TestInventoryLookupTables/codaFinalInv` - PASS
- `TestInventoryLookupTables/nounOnsetBeforeFrontInv` - PASS
- `TestInventoryLookupTables/nounOnsetBeforeBackInv` - PASS
- `TestInventoryLookupTables/verbRadicalOnsetInv` - PASS
- `TestDerivedInventoriesNoRejection/onsetInitialInv` - PASS
- `TestDerivedInventoriesNoRejection/onsetMedialInv` - PASS
- `TestDerivedInventoriesNoRejection/nucleiFrontInv` - PASS
- `TestDerivedInventoriesNoRejection/nucleiBackInv` - PASS
- `TestDerivedInventoriesNoRejection/codaBeforePBInv` - PASS
- `TestDerivedInventoriesNoRejection/codaBeforeOtherInv` - PASS
- `TestDerivedInventoriesNoRejection/codaFinalInv` - PASS
- `TestSyllableIsValid/simple_open` - PASS
- `TestSyllableIsValid/onset_cluster` - PASS
- `TestSyllableIsValid/closed_with_coda` - PASS
- `TestSyllableIsValid/no_onset` - PASS
- `TestSyllableIsValid/nasal_diphthong_nucleus` - PASS
- `TestSyllableIsValid/lh_word-initial_rejected` - PASS
- `TestSyllableIsValid/lh_medial_accepted` - PASS
- `TestSyllableIsValid/nh_word-initial_rejected` - PASS
- `TestSyllableIsValid/nh_medial_accepted` - PASS
- `TestSyllableIsValid/ç_word-initial_rejected` - PASS
- `TestSyllableIsValid/ç_before_back_vowel_accepted` - PASS
- `TestSyllableIsValid/ç_before_front_vowel_rejected` - PASS
- `TestSyllableIsValid/ç_before_i_rejected` - PASS
- `TestSyllableIsValid/qu_before_front_vowel_accepted` - PASS
- `TestSyllableIsValid/qu_before_i_accepted` - PASS
- `TestSyllableIsValid/qu_before_back_vowel_rejected` - PASS
- `TestSyllableIsValid/gu_before_front_vowel_accepted` - PASS
- `TestSyllableIsValid/gu_before_back_vowel_rejected` - PASS
- `TestSyllableIsValid/bad_onset_cluster` - PASS
- `TestSyllableIsValid/illegal_coda` - PASS
- `TestSyllableIsValid/invalid_nucleus` - PASS
- `TestTransitionIsValid/empty_coda` - PASS
- `TestTransitionIsValid/m_before_p` - PASS
- `TestTransitionIsValid/m_before_b` - PASS
- `TestTransitionIsValid/m_before_t_rejected` - PASS
- `TestTransitionIsValid/m_before_vowel_rejected` - PASS
- `TestTransitionIsValid/n_before_t` - PASS
- `TestTransitionIsValid/n_before_p_rejected` - PASS
- `TestTransitionIsValid/n_before_b_rejected` - PASS
- `TestTransitionIsValid/s_before_consonant` - PASS
- `TestTransitionIsValid/s_before_cluster` - PASS
- `TestTransitionIsValid/r_before_consonant` - PASS
- `TestTransitionIsValid/coda_before_vowel_rejected_(MOP)` - PASS
- `TestTransitionIsValid/l_before_vowel_rejected_(MOP)` - PASS
- `TestStemConformsToPhonotactics/accept/prato` - PASS
- `TestStemConformsToPhonotactics/accept/campo_(m_before_p)` - PASS
- `TestStemConformsToPhonotactics/accept/canto_(n_before_t)` - PASS
- `TestStemConformsToPhonotactics/accept/olho_(lh_medial)` - PASS
- `TestStemConformsToPhonotactics/accept/casaco` - PASS
- `TestStemConformsToPhonotactics/accept/coração_(ç_medial_+_nasal)` - PASS
- `TestStemConformsToPhonotactics/reject/empty_stem` - PASS
- `TestStemConformsToPhonotactics/reject/m_before_t` - PASS
- `TestStemConformsToPhonotactics/reject/n_before_p` - PASS
- `TestStemConformsToPhonotactics/reject/coda_before_vowel_(MOP)` - PASS
- `TestStemConformsToPhonotactics/reject/bad_onset_cluster` - PASS
- `TestStemConformsToPhonotactics/reject/illegal_coda` - PASS
- `TestStemConformsToPhonotactics/reject/lh_word-initial` - PASS
- `TestStemConformsToPhonotactics/reject/ç_word-initial` - PASS
- `TestStemConformsToPhonotactics/reject/invalid_nucleus` - PASS
- `TestStemConformsToAccentuation/accept/médico` - PASS
- `TestStemConformsToAccentuation/accept/câmara` - PASS
- `TestStemConformsToAccentuation/accept/café` - PASS
- `TestStemConformsToAccentuation/accept/sofá` - PASS
- `TestStemConformsToAccentuation/accept/também_(-em)` - PASS
- `TestStemConformsToAccentuation/accept/animal` - PASS
- `TestStemConformsToAccentuation/accept/feliz` - PASS
- `TestStemConformsToAccentuation/accept/casa` - PASS
- `TestStemConformsToAccentuation/accept/casas` - PASS
- `TestStemConformsToAccentuation/accept/homem_(-em_default)` - PASS
- `TestStemConformsToAccentuation/accept/fácil_(-l)` - PASS
- `TestStemConformsToAccentuation/accept/lápis_(i+s)` - PASS
- `TestStemConformsToAccentuation/accept/órfã_(nasal_end)` - PASS
- `TestStemConformsToAccentuation/accept/irmã` - PASS
- `TestStemConformsToAccentuation/accept/coração` - PASS
- `TestStemConformsToAccentuation/reject/medico_(no_accent,_proparoxytone)` - PASS
- `TestStemConformsToAccentuation/reject/cafe_(no_accent,_oxytone)` - PASS
- `TestStemConformsToAccentuation/reject/tambem_(no_accent,_-em_oxytone)` - PASS
- `TestStemConformsToAccentuation/reject/facil_(no_accent,_paroxytone_-l)` - PASS
- `TestStemConformsToAccentuation/reject/cása_(spurious_on_tonic_paroxytone)` - PASS
- `TestStemConformsToAccentuation/reject/anìmal-like_spurious` - PASS
- `TestStemConformsToAccentuation/reject/accent_on_non-tonic` - PASS
- `TestStemConformsToAccentuation/reject/acute_on_non-tonic_with_nasal_tonic` - PASS
- `TestStemConformsToAccentuation/reject/tonic_out_of_range` - PASS
- `TestStemConformsToAccentuation/reject/empty_stem` - PASS
- `TestSampleInventoriesNoRejection/onsetSingles` - PASS
- `TestSampleInventoriesNoRejection/onsetClusters` - PASS
- `TestSampleInventoriesNoRejection/nuclei` - PASS
- `TestSampleInventoriesNoRejection/codas` - PASS

### Examples
- `ExampleNounPTOf` - PASS
- `ExampleAdjectivePTOf` - PASS
- `ExampleVerbPTOf` - PASS
- `ExampleAdverbPT` - PASS
- `ExamplePrepositionPT` - PASS
- `ExampleWordsPT` - PASS
- `ExampleGenerator` - PASS

---

## Performance Tests (Benchmarks)

These benchmarks are the same 36 benchmarks as in the previous edition of this report, run on the same machine. `BENCHMARKS.md` covers the WordsPT and concurrent (`BenchmarkParallel*`) benchmarks.

Command:

```bash
go test -benchmem -run='^$' -bench '^Benchmark(Bool|Date|UnixDate|DateBetween|Int8|Int8Between|Int16|Int16Between|Int32|Int32Between|Int|IntBetween|Int64|Int64Between|Uint8|Uint8Between|Byte|Uint16|Uint16Between|Uint32|Uint32Between|Uint64|Uint64Between|Float32|Float32Between|Float64|Float64Between|Complex64|Complex64Between|Complex128|Complex128Between|String|StringNumeric|Word|WordByLengthType|Words)$' -benchtime 5s .
```

### Summary by Category

#### Booleans and Dates
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkBool` | 1,000,000,000 | 4.763 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkDate` | 813,372,056 | 7.472 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUnixDate` | 872,620,770 | 7.134 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkDateBetween` | 571,210,183 | 10.53 ns/op | 0 B/op, 0 allocs/op |

#### Signed Integers (Int8/16/32/64)
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkInt8` | 1,000,000,000 | 4.713 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt8Between` | 837,411,355 | 7.431 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt16` | 1,000,000,000 | 4.724 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt16Between` | 808,625,276 | 7.383 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt32` | 1,000,000,000 | 4.731 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt32Between` | 809,719,130 | 7.433 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt` | 1,000,000,000 | 4.742 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkIntBetween` | 781,577,872 | 7.620 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt64` | 1,000,000,000 | 4.724 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkInt64Between` | 792,635,524 | 7.630 ns/op | 0 B/op, 0 allocs/op |

#### Unsigned Integers (Uint8/16/32/64) and Byte
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkUint8` | 1,000,000,000 | 4.731 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint8Between` | 820,339,515 | 7.410 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkByte` | 1,000,000,000 | 4.734 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint16` | 1,000,000,000 | 4.716 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint16Between` | 830,340,750 | 7.418 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint32` | 1,000,000,000 | 4.719 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint32Between` | 818,196,580 | 7.431 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint64` | 1,000,000,000 | 4.722 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkUint64Between` | 787,919,929 | 7.664 ns/op | 0 B/op, 0 allocs/op |

#### Floating Point and Complex Numbers
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkFloat32` | 1,000,000,000 | 6.000 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat32Between` | 960,039,884 | 6.387 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat64` | 1,000,000,000 | 6.112 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkFloat64Between` | 914,943,406 | 6.549 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex64` | 593,035,335 | 10.23 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex64Between` | 445,344,910 | 13.54 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex128` | 573,194,046 | 10.48 ns/op | 0 B/op, 0 allocs/op |
| `BenchmarkComplex128Between` | 431,893,446 | 13.82 ns/op | 0 B/op, 0 allocs/op |

#### Strings and Words
| Benchmark | Operations | Time/op | Allocations |
|-----------|------------|---------|-------------|
| `BenchmarkString` | 215,655,939 | 27.77 ns/op | 8 B/op, 1 allocs/op |
| `BenchmarkStringNumeric` | 95,627,598 | 60.93 ns/op | 8 B/op, 1 allocs/op |
| `BenchmarkWord` | 70,306,759 | 84.56 ns/op | 18 B/op, 1 allocs/op |
| `BenchmarkWordByLengthType` | 121,923,086 | 49.38 ns/op | 7 B/op, 1 allocs/op |
| `BenchmarkWords` | 7,855,239 | 769.1 ns/op | 258 B/op, 10 allocs/op |

---

## Performance Analysis

### Highlights

1. **Zero Allocations in Primitive Types**: All number, boolean, and date generation functions perform no memory allocations (0 allocs/op), adhering to the library's performance philosophy.

2. **Fast Execution**: Integer, float, and boolean functions execute in ~4.7–7.7 nanoseconds per operation, achieving hundreds of millions to over a billion operations in 5 seconds of benchmarking.

3. **Fastest Functions** (top 5):
   - `BenchmarkInt8`: 4.713 ns/op
   - `BenchmarkUint16`: 4.716 ns/op
   - `BenchmarkUint32`: 4.719 ns/op
   - `BenchmarkUint64`: 4.722 ns/op
   - `BenchmarkInt16`: 4.724 ns/op

4. **Complex Numbers Slower**: `BenchmarkComplex64` and `BenchmarkComplex128` are naturally slower (~10.4 ns/op), and their `Between` variants slower still (~13.7 ns/op), as they generate two parts (real and imaginary).

5. **Strings with Controlled Allocation**: String and word functions (`BenchmarkString`, `BenchmarkStringNumeric`, `BenchmarkWord`, `BenchmarkWordByLengthType`, `BenchmarkWords`) are the only ones that allocate (1 alloc/op for single results; `BenchmarkWords` allocates once per generated word), as they must allocate memory for the returned string(s).

### Aggregated Metrics

| Category | Best Time | Worst Time | Approximate Average |
|----------|-----------|------------|---------------------|
| Booleans & Dates | 4.763 ns/op | 10.53 ns/op | ~7.5 ns/op |
| Integers | 4.713 ns/op | 7.664 ns/op | ~6.0 ns/op |
| Floats & Complex | 6.000 ns/op | 13.82 ns/op | ~9.1 ns/op |
| Strings | 27.77 ns/op | 60.93 ns/op | ~44.4 ns/op |

---

## Evidence (Raw Output)

### Unit Test Output
```
=== RUN   TestBool
--- PASS: TestBool (0.00s)
=== RUN   TestDate
--- PASS: TestDate (0.00s)
=== RUN   TestUnixDate
--- PASS: TestUnixDate (0.00s)
=== RUN   TestDateBetween
--- PASS: TestDateBetween (0.00s)
=== RUN   TestDateBetweenGranularity
--- PASS: TestDateBetweenGranularity (0.00s)
=== RUN   TestGeneratorReproducible
--- PASS: TestGeneratorReproducible (0.02s)
=== RUN   TestNewSourceReproducible
--- PASS: TestNewSourceReproducible (0.00s)
=== RUN   TestGeneratorRanges
--- PASS: TestGeneratorRanges (0.00s)
=== RUN   TestGeneratorIntSwapBranches
--- PASS: TestGeneratorIntSwapBranches (0.00s)
=== RUN   TestGeneratorUintSwapAndFullRange
--- PASS: TestGeneratorUintSwapAndFullRange (0.00s)
=== RUN   TestGeneratorFloatBranches
--- PASS: TestGeneratorFloatBranches (0.00s)
=== RUN   TestGeneratorStringEdges
--- PASS: TestGeneratorStringEdges (0.14s)
=== RUN   TestGeneratorDateAndWordEdges
--- PASS: TestGeneratorDateAndWordEdges (0.00s)
=== RUN   TestGeneratorStringAllocs
--- PASS: TestGeneratorStringAllocs (0.18s)
=== RUN   TestGeneratorStringGolden
--- PASS: TestGeneratorStringGolden (0.01s)
=== RUN   TestInt8
--- PASS: TestInt8 (0.00s)
=== RUN   TestInt8Between
--- PASS: TestInt8Between (0.00s)
=== RUN   TestInt16
--- PASS: TestInt16 (0.00s)
=== RUN   TestInt16Between
--- PASS: TestInt16Between (0.00s)
=== RUN   TestInt32
--- PASS: TestInt32 (0.00s)
=== RUN   TestInt32Between
--- PASS: TestInt32Between (0.00s)
=== RUN   TestInt
--- PASS: TestInt (0.00s)
=== RUN   TestIntBetween
--- PASS: TestIntBetween (0.00s)
=== RUN   TestInt64
--- PASS: TestInt64 (0.00s)
=== RUN   TestInt64Between
--- PASS: TestInt64Between (0.00s)
=== RUN   TestUint8
--- PASS: TestUint8 (0.00s)
=== RUN   TestUint8Between
--- PASS: TestUint8Between (0.00s)
=== RUN   TestByte
--- PASS: TestByte (0.00s)
=== RUN   TestUint16
--- PASS: TestUint16 (0.00s)
=== RUN   TestUint16Between
--- PASS: TestUint16Between (0.00s)
=== RUN   TestUint32
--- PASS: TestUint32 (0.00s)
=== RUN   TestUint32Between
--- PASS: TestUint32Between (0.00s)
=== RUN   TestUint64
--- PASS: TestUint64 (0.00s)
=== RUN   TestUint64Between
--- PASS: TestUint64Between (0.00s)
=== RUN   TestFloat32
--- PASS: TestFloat32 (0.00s)
=== RUN   TestFloat32Between
--- PASS: TestFloat32Between (0.00s)
=== RUN   TestFloat64
--- PASS: TestFloat64 (0.00s)
=== RUN   TestFloat64Between
--- PASS: TestFloat64Between (0.00s)
=== RUN   TestComplex64
--- PASS: TestComplex64 (0.00s)
=== RUN   TestComplex64Between
--- PASS: TestComplex64Between (0.00s)
=== RUN   TestComplex128
--- PASS: TestComplex128 (0.00s)
=== RUN   TestComplex128Between
--- PASS: TestComplex128Between (0.00s)
=== RUN   TestString
--- PASS: TestString (0.27s)
=== RUN   TestStringByteOriented
--- PASS: TestStringByteOriented (0.01s)
=== RUN   TestStringAllChars
--- PASS: TestStringAllChars (0.01s)
=== RUN   TestStringAlphanumeric
--- PASS: TestStringAlphanumeric (0.01s)
=== RUN   TestStringAlphabetic
--- PASS: TestStringAlphabetic (0.01s)
=== RUN   TestStringAlphabeticUppercase
--- PASS: TestStringAlphabeticUppercase (0.01s)
=== RUN   TestStringAlphabeticLowercase
--- PASS: TestStringAlphabeticLowercase (0.01s)
=== RUN   TestStringNumeric
--- PASS: TestStringNumeric (0.01s)
=== RUN   TestStringHexadecimal
--- PASS: TestStringHexadecimal (0.01s)
=== RUN   TestStringSymbols
--- PASS: TestStringSymbols (0.01s)
=== RUN   TestStringBetween
--- PASS: TestStringBetween (0.01s)
=== RUN   TestStringAllocs
--- PASS: TestStringAllocs (0.18s)
=== RUN   TestStringSingleCharLong
--- PASS: TestStringSingleCharLong (0.00s)
=== RUN   TestWord
--- PASS: TestWord (0.00s)
=== RUN   TestWordByLengthType
=== RUN   TestWordByLengthType/SmallLengthWord
=== RUN   TestWordByLengthType/MediumLengthWords
=== RUN   TestWordByLengthType/BigLengthWords
=== RUN   TestWordByLengthType/Default
--- PASS: TestWordByLengthType (0.01s)
    --- PASS: TestWordByLengthType/SmallLengthWord (0.00s)
    --- PASS: TestWordByLengthType/MediumLengthWords (0.00s)
    --- PASS: TestWordByLengthType/BigLengthWords (0.00s)
    --- PASS: TestWordByLengthType/Default (0.00s)
=== RUN   TestWords
=== RUN   TestWords/PositiveLength
=== RUN   TestWords/ZeroLength
=== RUN   TestWords/NegativeLength
--- PASS: TestWords (0.00s)
    --- PASS: TestWords/PositiveLength (0.00s)
    --- PASS: TestWords/ZeroLength (0.00s)
    --- PASS: TestWords/NegativeLength (0.00s)
=== RUN   TestAccentuateStemConformsToAccentuation
--- PASS: TestAccentuateStemConformsToAccentuation (0.28s)
=== RUN   TestAccentuateStemKnownCases
=== RUN   TestAccentuateStemKnownCases/oxytone-a
=== RUN   TestAccentuateStemKnownCases/oxytone-e
=== RUN   TestAccentuateStemKnownCases/oxytone-o
=== RUN   TestAccentuateStemKnownCases/oxytone-as
=== RUN   TestAccentuateStemKnownCases/oxytone-em-acute
=== RUN   TestAccentuateStemKnownCases/oxytone-im-none
=== RUN   TestAccentuateStemKnownCases/oxytone-u-none
=== RUN   TestAccentuateStemKnownCases/proparoxytone-a
=== RUN   TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-onset
=== RUN   TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-coda
=== RUN   TestAccentuateStemKnownCases/proparoxytone-e-stays-acute-before-nasal
=== RUN   TestAccentuateStemKnownCases/paroxytone-r-nondefault
=== RUN   TestAccentuateStemKnownCases/paroxytone-diphthong-nondefault
=== RUN   TestAccentuateStemKnownCases/paroxytone-default-a-none
=== RUN   TestAccentuateStemKnownCases/paroxytone-default-as-none
=== RUN   TestAccentuateStemKnownCases/paroxytone-nasal-ending-acute-tonic
=== RUN   TestAccentuateStemKnownCases/nasal-vowel-tonic
=== RUN   TestAccentuateStemKnownCases/nasal-diphthong-tonic-cedilla
--- PASS: TestAccentuateStemKnownCases (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-a (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-e (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-o (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-as (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-em-acute (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-im-none (0.00s)
    --- PASS: TestAccentuateStemKnownCases/oxytone-u-none (0.00s)
    --- PASS: TestAccentuateStemKnownCases/proparoxytone-a (0.00s)
    --- PASS: TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-onset (0.00s)
    --- PASS: TestAccentuateStemKnownCases/proparoxytone-a-circumflex-nasal-coda (0.00s)
    --- PASS: TestAccentuateStemKnownCases/proparoxytone-e-stays-acute-before-nasal (0.00s)
    --- PASS: TestAccentuateStemKnownCases/paroxytone-r-nondefault (0.00s)
    --- PASS: TestAccentuateStemKnownCases/paroxytone-diphthong-nondefault (0.00s)
    --- PASS: TestAccentuateStemKnownCases/paroxytone-default-a-none (0.00s)
    --- PASS: TestAccentuateStemKnownCases/paroxytone-default-as-none (0.00s)
    --- PASS: TestAccentuateStemKnownCases/paroxytone-nasal-ending-acute-tonic (0.00s)
    --- PASS: TestAccentuateStemKnownCases/nasal-vowel-tonic (0.00s)
    --- PASS: TestAccentuateStemKnownCases/nasal-diphthong-tonic-cedilla (0.00s)
=== RUN   TestAccentuateStemDeterministic
--- PASS: TestAccentuateStemDeterministic (0.02s)
=== RUN   TestAccentedWordNoRandomDraws
--- PASS: TestAccentedWordNoRandomDraws (0.00s)
=== RUN   TestAssembleAccentedStemSingleAllocation
--- PASS: TestAssembleAccentedStemSingleAllocation (0.00s)
=== RUN   TestAccentedWordAllocationBudget
--- PASS: TestAccentedWordAllocationBudget (0.01s)
=== RUN   TestAccentuateStemEmptyIsSafe
--- PASS: TestAccentuateStemEmptyIsSafe (0.00s)
=== RUN   TestAdjectivePTOrthographicConformance
--- PASS: TestAdjectivePTOrthographicConformance (4.93s)
=== RUN   TestAdjectivePTGenderHonored
--- PASS: TestAdjectivePTGenderHonored (1.88s)
=== RUN   TestAdjectivePTNumberHonored
--- PASS: TestAdjectivePTNumberHonored (1.88s)
=== RUN   TestAdjectivePTDegreeHonored
--- PASS: TestAdjectivePTDegreeHonored (0.97s)
=== RUN   TestAdjectivePTSuperlativeInflection
--- PASS: TestAdjectivePTSuperlativeInflection (1.73s)
=== RUN   TestAdjectivePTSuperlativeOrthography
--- PASS: TestAdjectivePTSuperlativeOrthography (0.76s)
=== RUN   TestAdjectivePTThematicOAProduced
--- PASS: TestAdjectivePTThematicOAProduced (0.96s)
=== RUN   TestAdjectivePTSuperlativeHardeningViaGeneration
--- PASS: TestAdjectivePTSuperlativeHardeningViaGeneration (0.00s)
=== RUN   TestAdjectiveSuperlativeHardeningRule
--- PASS: TestAdjectiveSuperlativeHardeningRule (0.00s)
=== RUN   TestAdjectivePTLengthHonored
--- PASS: TestAdjectivePTLengthHonored (4.73s)
=== RUN   TestAdjectivePTAnyGenderVaried
--- PASS: TestAdjectivePTAnyGenderVaried (0.48s)
=== RUN   TestAdjectivePTAnyNumberVaried
--- PASS: TestAdjectivePTAnyNumberVaried (0.50s)
=== RUN   TestAdjectivePTAnyDegreeVaried
--- PASS: TestAdjectivePTAnyDegreeVaried (0.49s)
=== RUN   TestAdjectivePTAnyLengthVaried
--- PASS: TestAdjectivePTAnyLengthVaried (0.48s)
=== RUN   TestAdjectivePTEndingVariety
--- PASS: TestAdjectivePTEndingVariety (0.97s)
=== RUN   TestAdjectivePTReproducible
--- PASS: TestAdjectivePTReproducible (1.47s)
=== RUN   TestAdjectiveEndingLengthAwareSelection
--- PASS: TestAdjectiveEndingLengthAwareSelection (0.09s)
=== RUN   TestAdjectiveStemConformsToOracles
--- PASS: TestAdjectiveStemConformsToOracles (0.86s)
=== RUN   TestAdjectivePTConcurrentSafe
--- PASS: TestAdjectivePTConcurrentSafe (0.60s)
=== RUN   TestAdverbPTOrthographicConformance
--- PASS: TestAdverbPTOrthographicConformance (2.07s)
=== RUN   TestAdverbPTEndsInMente
--- PASS: TestAdverbPTEndsInMente (1.43s)
=== RUN   TestAdverbPTNoGraphicAccent
--- PASS: TestAdverbPTNoGraphicAccent (1.86s)
=== RUN   TestAdverbPTInvariable
--- PASS: TestAdverbPTInvariable (0.38s)
=== RUN   TestAdverbWordFormation
--- PASS: TestAdverbWordFormation (0.00s)
=== RUN   TestAdverbPTLengthHonored
--- PASS: TestAdverbPTLengthHonored (1.44s)
=== RUN   TestAdverbPTSmallNormalizesUp
--- PASS: TestAdverbPTSmallNormalizesUp (0.34s)
=== RUN   TestAdverbPTEndingVariety
--- PASS: TestAdverbPTEndingVariety (0.38s)
=== RUN   TestAdverbEndingLengthAwareSelection
--- PASS: TestAdverbEndingLengthAwareSelection (0.06s)
=== RUN   TestAdverbPTReproducible
--- PASS: TestAdverbPTReproducible (1.04s)
=== RUN   TestAdverbPTConcurrentSafe
--- PASS: TestAdverbPTConcurrentSafe (0.59s)
=== RUN   TestClosedClassListsCuration
=== RUN   TestClosedClassListsCuration/article
=== RUN   TestClosedClassListsCuration/preposition
=== RUN   TestClosedClassListsCuration/conjunction
=== RUN   TestClosedClassListsCuration/pronoun
=== RUN   TestClosedClassListsCuration/interjection
=== RUN   TestClosedClassListsCuration/numeral
--- PASS: TestClosedClassListsCuration (0.00s)
    --- PASS: TestClosedClassListsCuration/article (0.00s)
    --- PASS: TestClosedClassListsCuration/preposition (0.00s)
    --- PASS: TestClosedClassListsCuration/conjunction (0.00s)
    --- PASS: TestClosedClassListsCuration/pronoun (0.00s)
    --- PASS: TestClosedClassListsCuration/interjection (0.00s)
    --- PASS: TestClosedClassListsCuration/numeral (0.00s)
=== RUN   TestClosedClassListsMatchReference
=== RUN   TestClosedClassListsMatchReference/article
=== RUN   TestClosedClassListsMatchReference/preposition
=== RUN   TestClosedClassListsMatchReference/conjunction
=== RUN   TestClosedClassListsMatchReference/pronoun
=== RUN   TestClosedClassListsMatchReference/interjection
=== RUN   TestClosedClassListsMatchReference/numeral
--- PASS: TestClosedClassListsMatchReference (0.00s)
    --- PASS: TestClosedClassListsMatchReference/article (0.00s)
    --- PASS: TestClosedClassListsMatchReference/preposition (0.00s)
    --- PASS: TestClosedClassListsMatchReference/conjunction (0.00s)
    --- PASS: TestClosedClassListsMatchReference/pronoun (0.00s)
    --- PASS: TestClosedClassListsMatchReference/interjection (0.00s)
    --- PASS: TestClosedClassListsMatchReference/numeral (0.00s)
=== RUN   TestClosedClassMembershipAndReachability
=== RUN   TestClosedClassMembershipAndReachability/article
=== RUN   TestClosedClassMembershipAndReachability/preposition
=== RUN   TestClosedClassMembershipAndReachability/conjunction
=== RUN   TestClosedClassMembershipAndReachability/pronoun
=== RUN   TestClosedClassMembershipAndReachability/interjection
=== RUN   TestClosedClassMembershipAndReachability/numeral
--- PASS: TestClosedClassMembershipAndReachability (0.30s)
    --- PASS: TestClosedClassMembershipAndReachability/article (0.05s)
    --- PASS: TestClosedClassMembershipAndReachability/preposition (0.05s)
    --- PASS: TestClosedClassMembershipAndReachability/conjunction (0.05s)
    --- PASS: TestClosedClassMembershipAndReachability/pronoun (0.05s)
    --- PASS: TestClosedClassMembershipAndReachability/interjection (0.05s)
    --- PASS: TestClosedClassMembershipAndReachability/numeral (0.05s)
=== RUN   TestClosedClassPackageLevelMembership
--- PASS: TestClosedClassPackageLevelMembership (0.34s)
=== RUN   TestClosedClassReproducibility
=== RUN   TestClosedClassReproducibility/article
=== RUN   TestClosedClassReproducibility/preposition
=== RUN   TestClosedClassReproducibility/conjunction
=== RUN   TestClosedClassReproducibility/pronoun
=== RUN   TestClosedClassReproducibility/interjection
=== RUN   TestClosedClassReproducibility/numeral
--- PASS: TestClosedClassReproducibility (0.01s)
    --- PASS: TestClosedClassReproducibility/article (0.00s)
    --- PASS: TestClosedClassReproducibility/preposition (0.00s)
    --- PASS: TestClosedClassReproducibility/conjunction (0.00s)
    --- PASS: TestClosedClassReproducibility/pronoun (0.00s)
    --- PASS: TestClosedClassReproducibility/interjection (0.00s)
    --- PASS: TestClosedClassReproducibility/numeral (0.00s)
=== RUN   TestClosedClassNoAllocation
--- PASS: TestClosedClassNoAllocation (0.01s)
=== RUN   TestPronounPTSubtypeCoverage
--- PASS: TestPronounPTSubtypeCoverage (0.00s)
=== RUN   TestPronounPTIncludesVosForms
--- PASS: TestPronounPTIncludesVosForms (0.00s)
=== RUN   TestNumeralPTSubtypeCoverage
--- PASS: TestNumeralPTSubtypeCoverage (0.00s)
=== RUN   TestNumeralPTPortugueseSpelling
--- PASS: TestNumeralPTPortugueseSpelling (0.00s)
=== RUN   TestFlexionEnumZeroValues
--- PASS: TestFlexionEnumZeroValues (0.00s)
=== RUN   TestResolveConcretePassthrough
--- PASS: TestResolveConcretePassthrough (0.00s)
=== RUN   TestResolveAnyUniformAndValid
--- PASS: TestResolveAnyUniformAndValid (0.24s)
=== RUN   TestResolveReproducible
--- PASS: TestResolveReproducible (0.00s)
=== RUN   TestCharRangeOf
--- PASS: TestCharRangeOf (0.00s)
=== RUN   TestNormalizeLength
--- PASS: TestNormalizeLength (0.00s)
=== RUN   TestSyllablesForCharRange
--- PASS: TestSyllablesForCharRange (0.31s)
=== RUN   TestStemSyllablesForNormalizesAndSizes
--- PASS: TestStemSyllablesForNormalizesAndSizes (0.54s)
=== RUN   TestStemSyllablesForLandsInWindow
--- PASS: TestStemSyllablesForLandsInWindow (1.83s)
=== RUN   TestStemHiatusReduced
    wordspt_naturalness_test.go:91: bare-stem hiatus over 16000 samples: mean longest run=1.6220, >=3-run fraction=3.281%
--- PASS: TestStemHiatusReduced (0.26s)
=== RUN   TestNounPTHiatusReduced
    wordspt_naturalness_test.go:124: NounPT hiatus over 40000 words: mean longest run=1.5828, >=3-run fraction=4.145%, max run=5
--- PASS: TestNounPTHiatusReduced (0.94s)
=== RUN   TestBigNounLengthRealistic
    wordspt_naturalness_test.go:166: Big/Singular over 40000 words: mean length=12.978, fraction <=14 chars=73.58%
--- PASS: TestBigNounLengthRealistic (0.99s)
=== RUN   TestMedialOnsetPrefersConsonant
--- PASS: TestMedialOnsetPrefersConsonant (0.00s)
=== RUN   TestNounPTOrthographicConformance
--- PASS: TestNounPTOrthographicConformance (4.86s)
=== RUN   TestNounPTGenderHonored
--- PASS: TestNounPTGenderHonored (0.82s)
=== RUN   TestNounPTNumberHonored
--- PASS: TestNounPTNumberHonored (0.85s)
=== RUN   TestNounPTLengthHonored
--- PASS: TestNounPTLengthHonored (2.67s)
=== RUN   TestNounPTAnyGenderVaried
--- PASS: TestNounPTAnyGenderVaried (0.41s)
=== RUN   TestNounPTAnyNumberVaried
--- PASS: TestNounPTAnyNumberVaried (0.43s)
=== RUN   TestNounPTAnyLengthVaried
--- PASS: TestNounPTAnyLengthVaried (0.41s)
=== RUN   TestNounPTReproducible
--- PASS: TestNounPTReproducible (1.17s)
=== RUN   TestNounEndingLengthAwareSelection
--- PASS: TestNounEndingLengthAwareSelection (0.06s)
=== RUN   TestNounPTConcurrentSafe
--- PASS: TestNounPTConcurrentSafe (0.44s)
=== RUN   TestNounStemConformsToOracles
--- PASS: TestNounStemConformsToOracles (0.85s)
=== RUN   TestAdditivePluralClassificationMatchesPluralize
--- PASS: TestAdditivePluralClassificationMatchesPluralize (1.39s)
=== RUN   TestAccentedWordSuffixedMatchesPluralize
--- PASS: TestAccentedWordSuffixedMatchesPluralize (1.55s)
=== RUN   TestPluralBuildAllocationBudget
--- PASS: TestPluralBuildAllocationBudget (0.03s)
=== RUN   TestAdditivePluralSuffixClassification
--- PASS: TestAdditivePluralSuffixClassification (0.00s)
=== RUN   TestPluralizeRuleTable
=== RUN   TestPluralizeRuleTable/casa
=== RUN   TestPluralizeRuleTable/café
=== RUN   TestPluralizeRuleTable/maçã
=== RUN   TestPluralizeRuleTable/mãe
=== RUN   TestPluralizeRuleTable/herói
=== RUN   TestPluralizeRuleTable/gato
=== RUN   TestPluralizeRuleTable/homem
=== RUN   TestPluralizeRuleTable/jardim
=== RUN   TestPluralizeRuleTable/bom
=== RUN   TestPluralizeRuleTable/flor
=== RUN   TestPluralizeRuleTable/amor
=== RUN   TestPluralizeRuleTable/luz
=== RUN   TestPluralizeRuleTable/rapaz
=== RUN   TestPluralizeRuleTable/liquen
=== RUN   TestPluralizeRuleTable/tórax
=== RUN   TestPluralizeRuleTable/látex
=== RUN   TestPluralizeRuleTable/país
=== RUN   TestPluralizeRuleTable/lápis
=== RUN   TestPluralizeRuleTable/vírus
=== RUN   TestPluralizeRuleTable/animal
=== RUN   TestPluralizeRuleTable/papel
=== RUN   TestPluralizeRuleTable/anzol
=== RUN   TestPluralizeRuleTable/azul
=== RUN   TestPluralizeRuleTable/amável
=== RUN   TestPluralizeRuleTable/possível
=== RUN   TestPluralizeRuleTable/funil
=== RUN   TestPluralizeRuleTable/barril
=== RUN   TestPluralizeRuleTable/fácil
=== RUN   TestPluralizeRuleTable/útil
=== RUN   TestPluralizeRuleTable/difícil
--- PASS: TestPluralizeRuleTable (0.00s)
    --- PASS: TestPluralizeRuleTable/casa (0.00s)
    --- PASS: TestPluralizeRuleTable/café (0.00s)
    --- PASS: TestPluralizeRuleTable/maçã (0.00s)
    --- PASS: TestPluralizeRuleTable/mãe (0.00s)
    --- PASS: TestPluralizeRuleTable/herói (0.00s)
    --- PASS: TestPluralizeRuleTable/gato (0.00s)
    --- PASS: TestPluralizeRuleTable/homem (0.00s)
    --- PASS: TestPluralizeRuleTable/jardim (0.00s)
    --- PASS: TestPluralizeRuleTable/bom (0.00s)
    --- PASS: TestPluralizeRuleTable/flor (0.00s)
    --- PASS: TestPluralizeRuleTable/amor (0.00s)
    --- PASS: TestPluralizeRuleTable/luz (0.00s)
    --- PASS: TestPluralizeRuleTable/rapaz (0.00s)
    --- PASS: TestPluralizeRuleTable/liquen (0.00s)
    --- PASS: TestPluralizeRuleTable/tórax (0.00s)
    --- PASS: TestPluralizeRuleTable/látex (0.00s)
    --- PASS: TestPluralizeRuleTable/país (0.00s)
    --- PASS: TestPluralizeRuleTable/lápis (0.00s)
    --- PASS: TestPluralizeRuleTable/vírus (0.00s)
    --- PASS: TestPluralizeRuleTable/animal (0.00s)
    --- PASS: TestPluralizeRuleTable/papel (0.00s)
    --- PASS: TestPluralizeRuleTable/anzol (0.00s)
    --- PASS: TestPluralizeRuleTable/azul (0.00s)
    --- PASS: TestPluralizeRuleTable/amável (0.00s)
    --- PASS: TestPluralizeRuleTable/possível (0.00s)
    --- PASS: TestPluralizeRuleTable/funil (0.00s)
    --- PASS: TestPluralizeRuleTable/barril (0.00s)
    --- PASS: TestPluralizeRuleTable/fácil (0.00s)
    --- PASS: TestPluralizeRuleTable/útil (0.00s)
    --- PASS: TestPluralizeRuleTable/difícil (0.00s)
=== RUN   TestApplyAoPlural
--- PASS: TestApplyAoPlural (0.00s)
=== RUN   TestPluralizeAoIsWeightedAndReachable
--- PASS: TestPluralizeAoIsWeightedAndReachable (0.10s)
=== RUN   TestSampleAoPluralDistribution
--- PASS: TestSampleAoPluralDistribution (0.08s)
=== RUN   TestPluralizeEmpty
--- PASS: TestPluralizeEmpty (0.00s)
=== RUN   TestPluralizePreservesStringOrthography
--- PASS: TestPluralizePreservesStringOrthography (0.23s)
=== RUN   TestPluralizeReproducible
--- PASS: TestPluralizeReproducible (0.03s)
=== RUN   TestPluralizeAllocations
--- PASS: TestPluralizeAllocations (0.00s)
=== RUN   TestAllWeightedInventoriesListed
--- PASS: TestAllWeightedInventoriesListed (0.00s)
=== RUN   TestInventoryLookupTables
=== RUN   TestInventoryLookupTables/onsetSingles
=== RUN   TestInventoryLookupTables/onsetClusters
=== RUN   TestInventoryLookupTables/nuclei
=== RUN   TestInventoryLookupTables/codas
=== RUN   TestInventoryLookupTables/onsetInitialInv
=== RUN   TestInventoryLookupTables/onsetMedialInv
=== RUN   TestInventoryLookupTables/nucleiFrontInv
=== RUN   TestInventoryLookupTables/nucleiBackInv
=== RUN   TestInventoryLookupTables/codaBeforePBInv
=== RUN   TestInventoryLookupTables/codaBeforeOtherInv
=== RUN   TestInventoryLookupTables/codaFinalInv
=== RUN   TestInventoryLookupTables/nounOnsetBeforeFrontInv
=== RUN   TestInventoryLookupTables/nounOnsetBeforeBackInv
=== RUN   TestInventoryLookupTables/verbRadicalOnsetInv
--- PASS: TestInventoryLookupTables (0.01s)
    --- PASS: TestInventoryLookupTables/onsetSingles (0.00s)
    --- PASS: TestInventoryLookupTables/onsetClusters (0.00s)
    --- PASS: TestInventoryLookupTables/nuclei (0.00s)
    --- PASS: TestInventoryLookupTables/codas (0.00s)
    --- PASS: TestInventoryLookupTables/onsetInitialInv (0.00s)
    --- PASS: TestInventoryLookupTables/onsetMedialInv (0.00s)
    --- PASS: TestInventoryLookupTables/nucleiFrontInv (0.00s)
    --- PASS: TestInventoryLookupTables/nucleiBackInv (0.00s)
    --- PASS: TestInventoryLookupTables/codaBeforePBInv (0.00s)
    --- PASS: TestInventoryLookupTables/codaBeforeOtherInv (0.00s)
    --- PASS: TestInventoryLookupTables/codaFinalInv (0.00s)
    --- PASS: TestInventoryLookupTables/nounOnsetBeforeFrontInv (0.00s)
    --- PASS: TestInventoryLookupTables/nounOnsetBeforeBackInv (0.00s)
    --- PASS: TestInventoryLookupTables/verbRadicalOnsetInv (0.00s)
=== RUN   TestOnsetNucleusLicenseMatchesSourceSets
--- PASS: TestOnsetNucleusLicenseMatchesSourceSets (0.00s)
=== RUN   TestSampleOnsetNucleusMatchesSetLookup
--- PASS: TestSampleOnsetNucleusMatchesSetLookup (0.72s)
=== RUN   TestSkewThresholdsMatchFormula
--- PASS: TestSkewThresholdsMatchFormula (0.20s)
=== RUN   TestSkewTableCoversEveryLengthWindow
--- PASS: TestSkewTableCoversEveryLengthWindow (0.00s)
=== RUN   TestSkewedCharTargetWideSpanUsesFormula
--- PASS: TestSkewedCharTargetWideSpanUsesFormula (0.00s)
=== RUN   TestSkewTableIsGenerated
--- PASS: TestSkewTableIsGenerated (0.00s)
=== RUN   TestSampleSyllabicStemConformsToPhonotactics
--- PASS: TestSampleSyllabicStemConformsToPhonotactics (0.20s)
=== RUN   TestSampleSyllabicStemHonoursTonic
--- PASS: TestSampleSyllabicStemHonoursTonic (0.00s)
=== RUN   TestSampleSyllabicStemReproducible
--- PASS: TestSampleSyllabicStemReproducible (0.02s)
=== RUN   TestSampleSyllabicStemBufferReuse
--- PASS: TestSampleSyllabicStemBufferReuse (0.00s)
=== RUN   TestDerivedInventoriesNoRejection
=== RUN   TestDerivedInventoriesNoRejection/onsetInitialInv
=== RUN   TestDerivedInventoriesNoRejection/onsetMedialInv
=== RUN   TestDerivedInventoriesNoRejection/nucleiFrontInv
=== RUN   TestDerivedInventoriesNoRejection/nucleiBackInv
=== RUN   TestDerivedInventoriesNoRejection/codaBeforePBInv
=== RUN   TestDerivedInventoriesNoRejection/codaBeforeOtherInv
=== RUN   TestDerivedInventoriesNoRejection/codaFinalInv
--- PASS: TestDerivedInventoriesNoRejection (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/onsetInitialInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/onsetMedialInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/nucleiFrontInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/nucleiBackInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/codaBeforePBInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/codaBeforeOtherInv (0.00s)
    --- PASS: TestDerivedInventoriesNoRejection/codaFinalInv (0.00s)
=== RUN   TestOnsetInventoryMembership
--- PASS: TestOnsetInventoryMembership (0.00s)
=== RUN   TestNucleiInventoryPartition
--- PASS: TestNucleiInventoryPartition (0.00s)
=== RUN   TestCodaInventoryMembership
--- PASS: TestCodaInventoryMembership (0.00s)
=== RUN   TestAssembleStemSingleAllocation
--- PASS: TestAssembleStemSingleAllocation (0.00s)
=== RUN   TestAssembleStemContent
--- PASS: TestAssembleStemContent (0.00s)
=== RUN   TestFullStemAllocationBudget
--- PASS: TestFullStemAllocationBudget (0.01s)
=== RUN   TestNoRejectionDeterministicDraws
--- PASS: TestNoRejectionDeterministicDraws (0.00s)
=== RUN   TestOnsetIsValid
--- PASS: TestOnsetIsValid (0.00s)
=== RUN   TestNucleusIsValid
--- PASS: TestNucleusIsValid (0.00s)
=== RUN   TestCodaIsValid
--- PASS: TestCodaIsValid (0.00s)
=== RUN   TestSyllableIsValid
=== RUN   TestSyllableIsValid/simple_open
=== RUN   TestSyllableIsValid/onset_cluster
=== RUN   TestSyllableIsValid/closed_with_coda
=== RUN   TestSyllableIsValid/no_onset
=== RUN   TestSyllableIsValid/nasal_diphthong_nucleus
=== RUN   TestSyllableIsValid/lh_word-initial_rejected
=== RUN   TestSyllableIsValid/lh_medial_accepted
=== RUN   TestSyllableIsValid/nh_word-initial_rejected
=== RUN   TestSyllableIsValid/nh_medial_accepted
=== RUN   TestSyllableIsValid/ç_word-initial_rejected
=== RUN   TestSyllableIsValid/ç_before_back_vowel_accepted
=== RUN   TestSyllableIsValid/ç_before_front_vowel_rejected
=== RUN   TestSyllableIsValid/ç_before_i_rejected
=== RUN   TestSyllableIsValid/qu_before_front_vowel_accepted
=== RUN   TestSyllableIsValid/qu_before_i_accepted
=== RUN   TestSyllableIsValid/qu_before_back_vowel_rejected
=== RUN   TestSyllableIsValid/gu_before_front_vowel_accepted
=== RUN   TestSyllableIsValid/gu_before_back_vowel_rejected
=== RUN   TestSyllableIsValid/bad_onset_cluster
=== RUN   TestSyllableIsValid/illegal_coda
=== RUN   TestSyllableIsValid/invalid_nucleus
--- PASS: TestSyllableIsValid (0.00s)
    --- PASS: TestSyllableIsValid/simple_open (0.00s)
    --- PASS: TestSyllableIsValid/onset_cluster (0.00s)
    --- PASS: TestSyllableIsValid/closed_with_coda (0.00s)
    --- PASS: TestSyllableIsValid/no_onset (0.00s)
    --- PASS: TestSyllableIsValid/nasal_diphthong_nucleus (0.00s)
    --- PASS: TestSyllableIsValid/lh_word-initial_rejected (0.00s)
    --- PASS: TestSyllableIsValid/lh_medial_accepted (0.00s)
    --- PASS: TestSyllableIsValid/nh_word-initial_rejected (0.00s)
    --- PASS: TestSyllableIsValid/nh_medial_accepted (0.00s)
    --- PASS: TestSyllableIsValid/ç_word-initial_rejected (0.00s)
    --- PASS: TestSyllableIsValid/ç_before_back_vowel_accepted (0.00s)
    --- PASS: TestSyllableIsValid/ç_before_front_vowel_rejected (0.00s)
    --- PASS: TestSyllableIsValid/ç_before_i_rejected (0.00s)
    --- PASS: TestSyllableIsValid/qu_before_front_vowel_accepted (0.00s)
    --- PASS: TestSyllableIsValid/qu_before_i_accepted (0.00s)
    --- PASS: TestSyllableIsValid/qu_before_back_vowel_rejected (0.00s)
    --- PASS: TestSyllableIsValid/gu_before_front_vowel_accepted (0.00s)
    --- PASS: TestSyllableIsValid/gu_before_back_vowel_rejected (0.00s)
    --- PASS: TestSyllableIsValid/bad_onset_cluster (0.00s)
    --- PASS: TestSyllableIsValid/illegal_coda (0.00s)
    --- PASS: TestSyllableIsValid/invalid_nucleus (0.00s)
=== RUN   TestTransitionIsValid
=== RUN   TestTransitionIsValid/empty_coda
=== RUN   TestTransitionIsValid/m_before_p
=== RUN   TestTransitionIsValid/m_before_b
=== RUN   TestTransitionIsValid/m_before_t_rejected
=== RUN   TestTransitionIsValid/m_before_vowel_rejected
=== RUN   TestTransitionIsValid/n_before_t
=== RUN   TestTransitionIsValid/n_before_p_rejected
=== RUN   TestTransitionIsValid/n_before_b_rejected
=== RUN   TestTransitionIsValid/s_before_consonant
=== RUN   TestTransitionIsValid/s_before_cluster
=== RUN   TestTransitionIsValid/r_before_consonant
=== RUN   TestTransitionIsValid/coda_before_vowel_rejected_(MOP)
=== RUN   TestTransitionIsValid/l_before_vowel_rejected_(MOP)
--- PASS: TestTransitionIsValid (0.00s)
    --- PASS: TestTransitionIsValid/empty_coda (0.00s)
    --- PASS: TestTransitionIsValid/m_before_p (0.00s)
    --- PASS: TestTransitionIsValid/m_before_b (0.00s)
    --- PASS: TestTransitionIsValid/m_before_t_rejected (0.00s)
    --- PASS: TestTransitionIsValid/m_before_vowel_rejected (0.00s)
    --- PASS: TestTransitionIsValid/n_before_t (0.00s)
    --- PASS: TestTransitionIsValid/n_before_p_rejected (0.00s)
    --- PASS: TestTransitionIsValid/n_before_b_rejected (0.00s)
    --- PASS: TestTransitionIsValid/s_before_consonant (0.00s)
    --- PASS: TestTransitionIsValid/s_before_cluster (0.00s)
    --- PASS: TestTransitionIsValid/r_before_consonant (0.00s)
    --- PASS: TestTransitionIsValid/coda_before_vowel_rejected_(MOP) (0.00s)
    --- PASS: TestTransitionIsValid/l_before_vowel_rejected_(MOP) (0.00s)
=== RUN   TestStemConformsToPhonotactics
=== RUN   TestStemConformsToPhonotactics/accept/prato
=== RUN   TestStemConformsToPhonotactics/accept/campo_(m_before_p)
=== RUN   TestStemConformsToPhonotactics/accept/canto_(n_before_t)
=== RUN   TestStemConformsToPhonotactics/accept/olho_(lh_medial)
=== RUN   TestStemConformsToPhonotactics/accept/casaco
=== RUN   TestStemConformsToPhonotactics/accept/coração_(ç_medial_+_nasal)
=== RUN   TestStemConformsToPhonotactics/reject/empty_stem
=== RUN   TestStemConformsToPhonotactics/reject/m_before_t
=== RUN   TestStemConformsToPhonotactics/reject/n_before_p
=== RUN   TestStemConformsToPhonotactics/reject/coda_before_vowel_(MOP)
=== RUN   TestStemConformsToPhonotactics/reject/bad_onset_cluster
=== RUN   TestStemConformsToPhonotactics/reject/illegal_coda
=== RUN   TestStemConformsToPhonotactics/reject/lh_word-initial
=== RUN   TestStemConformsToPhonotactics/reject/ç_word-initial
=== RUN   TestStemConformsToPhonotactics/reject/invalid_nucleus
--- PASS: TestStemConformsToPhonotactics (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/prato (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/campo_(m_before_p) (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/canto_(n_before_t) (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/olho_(lh_medial) (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/casaco (0.00s)
    --- PASS: TestStemConformsToPhonotactics/accept/coração_(ç_medial_+_nasal) (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/empty_stem (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/m_before_t (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/n_before_p (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/coda_before_vowel_(MOP) (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/bad_onset_cluster (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/illegal_coda (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/lh_word-initial (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/ç_word-initial (0.00s)
    --- PASS: TestStemConformsToPhonotactics/reject/invalid_nucleus (0.00s)
=== RUN   TestStressClassOf
--- PASS: TestStressClassOf (0.00s)
=== RUN   TestStemConformsToAccentuation
=== RUN   TestStemConformsToAccentuation/accept/médico
=== RUN   TestStemConformsToAccentuation/accept/câmara
=== RUN   TestStemConformsToAccentuation/accept/café
=== RUN   TestStemConformsToAccentuation/accept/sofá
=== RUN   TestStemConformsToAccentuation/accept/também_(-em)
=== RUN   TestStemConformsToAccentuation/accept/animal
=== RUN   TestStemConformsToAccentuation/accept/feliz
=== RUN   TestStemConformsToAccentuation/accept/casa
=== RUN   TestStemConformsToAccentuation/accept/casas
=== RUN   TestStemConformsToAccentuation/accept/homem_(-em_default)
=== RUN   TestStemConformsToAccentuation/accept/fácil_(-l)
=== RUN   TestStemConformsToAccentuation/accept/lápis_(i+s)
=== RUN   TestStemConformsToAccentuation/accept/órfã_(nasal_end)
=== RUN   TestStemConformsToAccentuation/accept/irmã
=== RUN   TestStemConformsToAccentuation/accept/coração
=== RUN   TestStemConformsToAccentuation/reject/medico_(no_accent,_proparoxytone)
=== RUN   TestStemConformsToAccentuation/reject/cafe_(no_accent,_oxytone)
=== RUN   TestStemConformsToAccentuation/reject/tambem_(no_accent,_-em_oxytone)
=== RUN   TestStemConformsToAccentuation/reject/facil_(no_accent,_paroxytone_-l)
=== RUN   TestStemConformsToAccentuation/reject/cása_(spurious_on_tonic_paroxytone)
=== RUN   TestStemConformsToAccentuation/reject/anìmal-like_spurious
=== RUN   TestStemConformsToAccentuation/reject/accent_on_non-tonic
=== RUN   TestStemConformsToAccentuation/reject/acute_on_non-tonic_with_nasal_tonic
=== RUN   TestStemConformsToAccentuation/reject/tonic_out_of_range
=== RUN   TestStemConformsToAccentuation/reject/empty_stem
--- PASS: TestStemConformsToAccentuation (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/médico (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/câmara (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/café (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/sofá (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/também_(-em) (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/animal (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/feliz (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/casa (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/casas (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/homem_(-em_default) (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/fácil_(-l) (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/lápis_(i+s) (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/órfã_(nasal_end) (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/irmã (0.00s)
    --- PASS: TestStemConformsToAccentuation/accept/coração (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/medico_(no_accent,_proparoxytone) (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/cafe_(no_accent,_oxytone) (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/tambem_(no_accent,_-em_oxytone) (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/facil_(no_accent,_paroxytone_-l) (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/cása_(spurious_on_tonic_paroxytone) (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/anìmal-like_spurious (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/accent_on_non-tonic (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/acute_on_non-tonic_with_nasal_tonic (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/tonic_out_of_range (0.00s)
    --- PASS: TestStemConformsToAccentuation/reject/empty_stem (0.00s)
=== RUN   TestWordConformsToCedilla
--- PASS: TestWordConformsToCedilla (0.00s)
=== RUN   TestWordHasForbiddenDiaeresis
--- PASS: TestWordHasForbiddenDiaeresis (0.00s)
=== RUN   TestWordContainsGrave
--- PASS: TestWordContainsGrave (0.00s)
=== RUN   TestSampleInventoriesNoRejection
=== RUN   TestSampleInventoriesNoRejection/onsetSingles
=== RUN   TestSampleInventoriesNoRejection/onsetClusters
=== RUN   TestSampleInventoriesNoRejection/nuclei
=== RUN   TestSampleInventoriesNoRejection/codas
--- PASS: TestSampleInventoriesNoRejection (0.00s)
    --- PASS: TestSampleInventoriesNoRejection/onsetSingles (0.00s)
    --- PASS: TestSampleInventoriesNoRejection/onsetClusters (0.00s)
    --- PASS: TestSampleInventoriesNoRejection/nuclei (0.00s)
    --- PASS: TestSampleInventoriesNoRejection/codas (0.00s)
=== RUN   TestInventoryFormsAreValidUTF8
--- PASS: TestInventoryFormsAreValidUTF8 (0.00s)
=== RUN   TestVerbPTOrthographicConformance
--- PASS: TestVerbPTOrthographicConformance (5.43s)
=== RUN   TestVerbDesinenceMatchesAssemblers
--- PASS: TestVerbDesinenceMatchesAssemblers (0.00s)
=== RUN   TestVerbPTSlotHonored
--- PASS: TestVerbPTSlotHonored (0.60s)
=== RUN   TestVerbPTPresentTenseHonored
--- PASS: TestVerbPTPresentTenseHonored (0.52s)
=== RUN   TestVerbPTAnyMoodVaried
--- PASS: TestVerbPTAnyMoodVaried (0.12s)
=== RUN   TestResolveMoodPassthrough
--- PASS: TestResolveMoodPassthrough (0.00s)
=== RUN   TestVerbPTIndicativeTenseVaried
--- PASS: TestVerbPTIndicativeTenseVaried (0.26s)
=== RUN   TestVerbPTPersonNumberVaried
--- PASS: TestVerbPTPersonNumberVaried (0.26s)
=== RUN   TestResolveInfinitiveForm
--- PASS: TestResolveInfinitiveForm (0.13s)
=== RUN   TestVerbPTInfinitiveConcretePersonIsPersonal
--- PASS: TestVerbPTInfinitiveConcretePersonIsPersonal (0.51s)
=== RUN   TestVerbPTPresentNeedsNoAccent
--- PASS: TestVerbPTPresentNeedsNoAccent (0.79s)
=== RUN   TestVerbPTNormalizationN1
--- PASS: TestVerbPTNormalizationN1 (0.82s)
=== RUN   TestVerbPTNormalizationN2
--- PASS: TestVerbPTNormalizationN2 (0.20s)
=== RUN   TestVerbPTNormalizationN3Subjunctive
--- PASS: TestVerbPTNormalizationN3Subjunctive (0.46s)
=== RUN   TestVerbPTNormalizationN3TenseIgnored
--- PASS: TestVerbPTNormalizationN3TenseIgnored (1.21s)
=== RUN   TestVerbPTLengthHonored
--- PASS: TestVerbPTLengthHonored (2.66s)
=== RUN   TestVerbPTAnyLengthVaried
--- PASS: TestVerbPTAnyLengthVaried (0.26s)
=== RUN   TestVerbPTReproducible
--- PASS: TestVerbPTReproducible (0.96s)
=== RUN   TestVerbPTSingleAllocation
--- PASS: TestVerbPTSingleAllocation (0.10s)
=== RUN   TestVerbPTConcurrentSafe
--- PASS: TestVerbPTConcurrentSafe (0.27s)
=== RUN   TestVerbRadicalOnsetInventoryExcludesSensitive
--- PASS: TestVerbRadicalOnsetInventoryExcludesSensitive (0.00s)
=== RUN   TestNonFiniteExactForms
--- PASS: TestNonFiniteExactForms (0.00s)
=== RUN   TestPersonalInfinitiveIgnoresPersonWhereApplicable
--- PASS: TestPersonalInfinitiveIgnoresPersonWhereApplicable (0.00s)
=== RUN   TestNonFiniteVerbOrthographicConformance
--- PASS: TestNonFiniteVerbOrthographicConformance (0.75s)
=== RUN   TestNonFiniteVerbLengthHonored
--- PASS: TestNonFiniteVerbLengthHonored (0.62s)
=== RUN   TestNonFiniteVerbReproducible
--- PASS: TestNonFiniteVerbReproducible (0.24s)
=== RUN   TestPickConjugationWeighted
--- PASS: TestPickConjugationWeighted (0.03s)
=== RUN   TestResolveVerbPerson
--- PASS: TestResolveVerbPerson (0.02s)
=== RUN   TestVerbRadicalStemConformsToOracle
--- PASS: TestVerbRadicalStemConformsToOracle (0.23s)
=== RUN   TestConjugateNonFiniteSingleAllocation
--- PASS: TestConjugateNonFiniteSingleAllocation (0.00s)
=== RUN   TestIndicativeExactForms
--- PASS: TestIndicativeExactForms (0.00s)
=== RUN   TestIndicativePreteriteVsPresentFirstPlural
--- PASS: TestIndicativePreteriteVsPresentFirstPlural (0.00s)
=== RUN   TestConjugateIndicativeSingleAllocation
--- PASS: TestConjugateIndicativeSingleAllocation (0.03s)
=== RUN   TestIndicativeOrthographicConformance
--- PASS: TestIndicativeOrthographicConformance (0.00s)
=== RUN   TestResolveIndicativeTense
--- PASS: TestResolveIndicativeTense (0.01s)
=== RUN   TestIndicativeVerbFormConformance
--- PASS: TestIndicativeVerbFormConformance (0.31s)
=== RUN   TestIndicativeVerbFormReproducible
--- PASS: TestIndicativeVerbFormReproducible (0.05s)
=== RUN   TestSubjunctiveExactForms
--- PASS: TestSubjunctiveExactForms (0.00s)
=== RUN   TestSubjunctiveImperfectFirstPlural
--- PASS: TestSubjunctiveImperfectFirstPlural (0.00s)
=== RUN   TestSubjunctiveFutureEqualsPersonalInfinitive
--- PASS: TestSubjunctiveFutureEqualsPersonalInfinitive (0.00s)
=== RUN   TestImperativeExactForms
--- PASS: TestImperativeExactForms (0.00s)
=== RUN   TestImperativeDerivationRule
--- PASS: TestImperativeDerivationRule (0.00s)
=== RUN   TestImperativeNormalizesFirstSingular
--- PASS: TestImperativeNormalizesFirstSingular (0.00s)
=== RUN   TestConjugateSubjunctiveSingleAllocation
--- PASS: TestConjugateSubjunctiveSingleAllocation (0.02s)
=== RUN   TestConjugateImperativeSingleAllocation
--- PASS: TestConjugateImperativeSingleAllocation (0.01s)
=== RUN   TestSubjunctiveOrthographicConformance
--- PASS: TestSubjunctiveOrthographicConformance (0.00s)
=== RUN   TestImperativeOrthographicConformance
--- PASS: TestImperativeOrthographicConformance (0.00s)
=== RUN   TestResolveSubjunctiveTense
--- PASS: TestResolveSubjunctiveTense (0.01s)
=== RUN   TestResolveImperativePerson
--- PASS: TestResolveImperativePerson (0.04s)
=== RUN   TestSubjunctiveVerbFormConformance
--- PASS: TestSubjunctiveVerbFormConformance (0.32s)
=== RUN   TestImperativeVerbFormConformance
--- PASS: TestImperativeVerbFormConformance (0.09s)
=== RUN   TestSubjunctiveVerbFormReproducible
--- PASS: TestSubjunctiveVerbFormReproducible (0.05s)
=== RUN   TestImperativeVerbFormReproducible
--- PASS: TestImperativeVerbFormReproducible (0.01s)
=== RUN   TestResolveWordClassReachesAll
--- PASS: TestResolveWordClassReachesAll (0.28s)
=== RUN   TestResolveWordClassDistribution
--- PASS: TestResolveWordClassDistribution (0.28s)
=== RUN   TestWordPTOutputValidity
--- PASS: TestWordPTOutputValidity (2.35s)
=== RUN   TestWordPTClassMix
--- PASS: TestWordPTClassMix (0.84s)
=== RUN   TestWordPTByLengthTypeBigOpenWindow
--- PASS: TestWordPTByLengthTypeBigOpenWindow (0.94s)
=== RUN   TestWordPTClosedIgnoreLength
--- PASS: TestWordPTClosedIgnoreLength (0.00s)
=== RUN   TestWordsPTCount
--- PASS: TestWordsPTCount (0.08s)
=== RUN   TestWordsPTEmptyForNonPositive
--- PASS: TestWordsPTEmptyForNonPositive (0.00s)
=== RUN   TestWordPTReproducible
--- PASS: TestWordPTReproducible (2.13s)
=== RUN   TestWordPTConcurrentSafe
--- PASS: TestWordPTConcurrentSafe (1.40s)
=== RUN   TestWordPTAllocationBudget
--- PASS: TestWordPTAllocationBudget (0.22s)
=== RUN   ExampleNounPTOf
--- PASS: ExampleNounPTOf (0.00s)
=== RUN   ExampleAdjectivePTOf
--- PASS: ExampleAdjectivePTOf (0.00s)
=== RUN   ExampleVerbPTOf
--- PASS: ExampleVerbPTOf (0.00s)
=== RUN   ExampleAdverbPT
--- PASS: ExampleAdverbPT (0.00s)
=== RUN   ExamplePrepositionPT
--- PASS: ExamplePrepositionPT (0.00s)
=== RUN   ExampleWordsPT
--- PASS: ExampleWordsPT (0.00s)
=== RUN   ExampleGenerator
--- PASS: ExampleGenerator (0.00s)
PASS
coverage: 96.5% of statements
ok  	github.com/FlavioCFOliveira/gengo	85.955s	coverage: 96.5% of statements
```

### Benchmark Output
```
goos: linux
goarch: amd64
pkg: github.com/FlavioCFOliveira/gengo
cpu: AMD Ryzen 9 5900HX with Radeon Graphics        
BenchmarkBool-16                 	1000000000	         4.763 ns/op	       0 B/op	       0 allocs/op
BenchmarkDate-16                 	813372056	         7.472 ns/op	       0 B/op	       0 allocs/op
BenchmarkUnixDate-16             	872620770	         7.134 ns/op	       0 B/op	       0 allocs/op
BenchmarkDateBetween-16          	571210183	        10.53 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt8-16                 	1000000000	         4.713 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt8Between-16          	837411355	         7.431 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt16-16                	1000000000	         4.724 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt16Between-16         	808625276	         7.383 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt32-16                	1000000000	         4.731 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt32Between-16         	809719130	         7.433 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt-16                  	1000000000	         4.742 ns/op	       0 B/op	       0 allocs/op
BenchmarkIntBetween-16           	781577872	         7.620 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt64-16                	1000000000	         4.724 ns/op	       0 B/op	       0 allocs/op
BenchmarkInt64Between-16         	792635524	         7.630 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint8-16                	1000000000	         4.731 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint8Between-16         	820339515	         7.410 ns/op	       0 B/op	       0 allocs/op
BenchmarkByte-16                 	1000000000	         4.734 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16-16               	1000000000	         4.716 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint16Between-16        	830340750	         7.418 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32-16               	1000000000	         4.719 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint32Between-16        	818196580	         7.431 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint64-16               	1000000000	         4.722 ns/op	       0 B/op	       0 allocs/op
BenchmarkUint64Between-16        	787919929	         7.664 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat32-16              	1000000000	         6.000 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat32Between-16       	960039884	         6.387 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat64-16              	1000000000	         6.112 ns/op	       0 B/op	       0 allocs/op
BenchmarkFloat64Between-16       	914943406	         6.549 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex64-16            	593035335	        10.23 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex64Between-16     	445344910	        13.54 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex128-16           	573194046	        10.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkComplex128Between-16    	431893446	        13.82 ns/op	       0 B/op	       0 allocs/op
BenchmarkString-16               	215655939	        27.77 ns/op	       8 B/op	       1 allocs/op
BenchmarkStringNumeric-16        	95627598	        60.93 ns/op	       8 B/op	       1 allocs/op
BenchmarkWord-16                 	70306759	        84.56 ns/op	      18 B/op	       1 allocs/op
BenchmarkWordByLengthType-16     	121923086	        49.38 ns/op	       7 B/op	       1 allocs/op
BenchmarkWords-16                	 7855239	       769.1 ns/op	     258 B/op	      10 allocs/op
PASS
ok  	github.com/FlavioCFOliveira/gengo	233.688s
```

---

## Conclusion

The **gengo** library demonstrates high test coverage and optimized performance:

- **100% tests passing** (238/238 tests, 178 subtests, 7 examples), with no data races under `-race`
- **96.5% statement coverage**
- **Zero memory allocations** in all primitive functions
- **Execution times in the nanosecond range** for numeric types
- **Consistency** between functions with and without custom ranges

The library achieves its design goals of simplicity and performance, prioritizing execution speed and memory efficiency as documented in the project philosophy.

---

# WordsPT Allocation & Benchmark Audit (Task #37)

> **Historical section.** The audit below was recorded on 2026-07-22 for
> Task #37, in the v0.2.0 cycle. Its figures predate the v0.2.1 WordsPT sampling-core
> optimization; see the "WordsPT sampling core" section of
> [BENCHMARKS.md](BENCHMARKS.md) for the current figures.

**Date:** 2026-07-22
**Package:** github.com/FlavioCFOliveira/gengo
**Platform:** linux/amd64 (AMD Ryzen 9 5900HX, 16 threads) · Go go1.26.5
**Command:** `go test -run=^$ -bench 'PT' -benchmem -benchtime 2s -count 1 .`

This section covers the pt-PT word generators (the `WordsPT` feature) that were added after the original report above. It records the full allocation/benchmark surface of every public `WordsPT` function and documents the Task #37 optimization: building an **additive** noun/adjective plural in a single allocation.

## 1. Optimization: single-allocation additive plural

The noun and adjective plural path used to be two allocations: assemble the singular word (one allocation), then apply a tail string transform (`pluralize`) to it (a second allocation). Task #37 splits the plural into two build strategies decided up front from each ending's own metadata (`additivePluralSuffix`):

- **Additive plural** — a pure suffix append: a regular `+s` for vowel/diphthong endings, or `+es` for `-r`/`-z`/`-n`/oxytone-`-s` endings. The suffix is now written into the **same** `strings.Builder` that assembles the word (`accentedWordSuffixed`), so the plural is **one allocation**. Reachable additive endings: nouns `-o`, `-a`, `-eiro`, `-eira`, `-or` (the only `+es` case), `-ora`, `-mento`, `-dade`, `-ista`; adjectives `-o`, `-a`, `-oso`, `-osa`, `-ico`, `-ica`, `-ivo`, `-iva`, `-ente`, `-ante`.
- **Substitutive plural** — a tail rewrite that cannot be a suffix append: `-ão→-ões/-ães/-ãos` (incl. the fixed `-ção→-ções`), `-m→-ns` (`-agem→-agens`), and the vowel+`l` endings `-al/-ável/-ível→-ais/-áveis/-íveis`. These **keep the post-assembly `pluralize` transform** and stay **two allocations**, because they modify graphemes inside the word (e.g. `-agens` carries the `ns` coda cluster that no single-coda syllable can represent).

The produced plural **string is unchanged** — only how it is built changed. The additive classifier draws no randomness, and the substitutive endings it defers (fixed `-ção`, `-m`, `-l`) draw none either, so the random stream is byte-for-byte identical: same-seed reproducibility is preserved. This equivalence is locked by regression tests in `wordspt_plural_optimization_test.go` (`TestAdditivePluralClassificationMatchesPluralize`, `TestAccentedWordSuffixedMatchesPluralize`, `TestPluralBuildAllocationBudget`, `TestAdditivePluralSuffixClassification`).

### Before / After (plural benchmarks, `-benchtime 2s`)

| Benchmark (Generator surface) | Allocs before | Allocs after | B/op before | B/op after | ns/op before | ns/op after |
|-------------------------------|:-------------:|:------------:|:-----------:|:----------:|:------------:|:-----------:|
| `BenchmarkNounPTOfPlural` | 2 | **1** | 25 | **15** | 451.3 | **404.4** |
| `BenchmarkAdjectivePTOfPositivePlural` | 2 | **1** | 27 | **17** | 477.9 | **450.1** |

`testing`'s reported `allocs/op` is the integer-truncated **average** over the benchmark's random mix of endings (both benchmarks force `Plural` but pick the ending at random). The average now truncates to 1 because additive endings dominate (~87% of noun plurals, ~75% of adjective plurals). The exact per-path costs are pinned by `TestPluralBuildAllocationBudget`: an additive build is **exactly 1** allocation and a substitutive build is **exactly 2**.

## 2. Full public-function benchmark surface

`ns/op`, `B/op` and `allocs/op` for every public `WordsPT` function, on the package-level (global source) and `*Generator` surfaces. `*Of…` and `Generator…` rows use a seeded `New(1)` generator.

### Open classes (noun, adjective, verb, adverb)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkNounPT` | package | 457.4 | 13 | 1 |
| `BenchmarkNounPTOfSingular` | Generator | 405.4 | 12 | 1 |
| `BenchmarkNounPTOfPlural` | Generator | 405.3 | 15 | 1 (avg; additive 1 / substitutive 2) |
| `BenchmarkAdjectivePT` | package | 491.4 | 17 | 1 |
| `BenchmarkAdjectivePTOfPositiveSingular` | Generator | 430.6 | 13 | 1 |
| `BenchmarkAdjectivePTOfPositivePlural` | Generator | 448.2 | 17 | 1 (avg; additive 1 / substitutive 2) |
| `BenchmarkAdjectivePTOfSuperlativeSingular` | Generator | 398.7 | 19 | 1 |
| `BenchmarkVerbPT` | package | 444.1 | 13 | 1 |
| `BenchmarkVerbPTOfPresentIndicative` | Generator | 381.0 | 12 | 1 |
| `BenchmarkVerbPTOfImperfectSubjunctive` | Generator | 363.5 | 18 | 1 |
| `BenchmarkVerbPTOfInfinitive` | Generator | 371.6 | 12 | 1 |
| `BenchmarkAdverbPT` | package | 405.6 | 18 | 1 |
| `BenchmarkAdverbPTByLengthType` | Generator (Big) | 385.6 | 18 | 1 |

### Closed classes (curated selectors)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkArticlePT` | package | 8.567 | 0 | 0 |
| `BenchmarkGeneratorArticlePT` | Generator | 4.584 | 0 | 0 |
| `BenchmarkPrepositionPT` | package | 8.430 | 0 | 0 |
| `BenchmarkConjunctionPT` | package | 8.922 | 0 | 0 |
| `BenchmarkPronounPT` | package | 8.487 | 0 | 0 |
| `BenchmarkGeneratorPronounPT` | Generator | 4.579 | 0 | 0 |
| `BenchmarkInterjectionPT` | package | 8.582 | 0 | 0 |
| `BenchmarkNumeralPT` | package | 8.240 | 0 | 0 |
| `BenchmarkGeneratorNumeralPT` | Generator | 4.740 | 0 | 0 |

### Orchestration (`WordPT` / `WordsPT`)

| Benchmark | Surface | ns/op | B/op | allocs/op |
|-----------|---------|------:|-----:|:---------:|
| `BenchmarkWordPT` | package | 454.2 | 11 | 0 (avg; see note) |
| `BenchmarkGeneratorWordPT` | Generator | 414.8 | 11 | 0 (avg; see note) |
| `BenchmarkWordsPT` | Generator (16 words) | 6646 | 436 | 15 (≈0.94/word) |

## 3. Allocation classification (0 / 1 / 2 allocations, and why)

- **0 allocations — the six closed-class selectors** (`ArticlePT`, `PrepositionPT`, `ConjunctionPT`, `PronounPT`, `InterjectionPT`, `NumeralPT`). They return a string that already lives in a package-level `[]string`; the selector only indexes into the shared backing array, copying no bytes (asserted at 0 allocs by `TestClosedClassNoAllocation`). The `*Generator` variants are ~2× faster than the package-level ones (~4.6 ns vs ~8.5 ns) because the package-level path goes through the `globalSource` wrapper over the concurrency-safe global `math/rand/v2` generator, while a `*Generator` draws from its own unsynchronized PCG source.
- **1 allocation — every open-class singular, the additive plural, the superlative, verbs, and adverbs.** The syllable buffer is a stack array (reused across length attempts with no heap cost), so the only allocation is the final assembled string. The additive plural (Task #37) keeps this to one allocation by appending its suffix into the assembling builder.
- **2 allocations — the substitutive noun/adjective plural only.** The singular is assembled (one allocation) and then its tail is rewritten by `pluralize`/`pluralizeNoun` (a second allocation). This is inherent: `-ões`, `-ns`, `-ais` and friends change graphemes inside the word, so they cannot be produced by appending to the singular.
- **`WordPT` averages 0 allocations (truncated).** `WordPT` dispatches across the open and closed classes by type; the frequent zero-allocation closed-class words pull the truncated per-word average below one (`TestWordPTAllocationBudget` asserts the average stays below two). `WordsPT(16)` reports 15 allocs for 16 words (≈0.94/word) for the same reason, reusing one syllable buffer across the whole slice.

## 4. Gate results (Task #37)

| Gate | Result |
|------|--------|
| `gofmt -l .` | clean (no files) |
| `go vet ./...` | clean |
| `golangci-lint run` | 0 issues |
| `go test -race -count=1 .` | ok (full suite, no regression) |
| `go mod tidy -diff` | clean |
| `govulncheck ./...` | No vulnerabilities found |
