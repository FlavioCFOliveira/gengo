# WordsPT Open Classes

## Overview

This document defines the four open classes that WordsPT generates by
morphology: noun, adjective, verb, and adverb. For each class it defines the
class endings, the inflections, and the orthographic adjustments. It also
defines the normalization of impossible option combinations and the minimum
viable length per class.

This document is a companion to [wordspt.md](wordspt.md) and uses the
terminology defined there. The three-layer pipeline that produces these words is
defined in [wordspt-generation-model.md](wordspt-generation-model.md).

## Shared inflection rules for nouns and adjectives

Nouns and adjectives share the same gender and number inflection rules.

### Gender

- Endings that inflect for gender change the masculine marker to the feminine
  marker: `-o` becomes `-a` (`-eiro` becomes `-eira`, `-oso` becomes `-osa`,
  `-ico` becomes `-ica`, `-ivo` becomes `-iva`, `-or` becomes `-ora`).
- Endings that are invariable for gender keep the same form for masculine and
  feminine (`-ção`, `-dade`, `-agem`, `-mento`, `-al`, `-ável`, `-ível`,
  `-ente`, `-ante`, and the common-gender `-ista`).
- When a caller requests a specific gender, the generator selects a class ending
  whose grammatical gender is compatible with the request. For gender-inflecting
  endings it applies the requested gender form. For gender-fixed endings it uses
  them only when they are compatible with the requested gender. Common-gender
  endings serve either request with the same form.

### Number

The plural is formed from the singular by the following pt-PT rules. When an
ending admits more than one plural outcome, the generator selects one using
frequency weights sourced from authoritative references, and it uses the stress
position (which it controls) to resolve outcomes that depend on stress.

| Singular ending | Plural rule | Example pattern |
|-----------------|-------------|-----------------|
| vowel or oral diphthong | add `s` | `casa` to `casas` |
| `-ão` | `-ões`, `-ães`, or `-ãos` (weighted; `-ões` is the default) | `coração` to `corações`; `pão` to `pães`; `mão` to `mãos` |
| `-al`, `-el`, `-ol` | `-ais`, `-éis`, `-óis` | `animal` to `animais`; `papel` to `papéis`; `lençol` to `lençóis` |
| `-il` (stressed) | `-is` | `funil` to `funis` |
| `-il` (unstressed) | `-eis` | `fácil` to `fáceis` |
| `-m` | `-ns` | `homem` to `homens` |
| `-r`, `-z` | add `es` | `flor` to `flores`; `luz` to `luzes` |
| `-s` (stressed final syllable) | add `es` | `país` to `países` |
| `-s` (unstressed final syllable) | invariable | `lápis` to `lápis` |

The choice between `-is` and `-eis` for `-il`, and between adding `es` and
staying invariable for `-s`, is decided by the stress position of the generated
word, which is known from Layer 1.

## Noun (Substantivo)

### Endings

A generated noun uses one of the following typical pt-PT noun endings:
`-o`, `-a`, `-ção`, `-dade`, `-mento`, `-agem`, `-eiro`, `-eira`, `-or`,
`-ista`.

### Inherent gender of endings

- Gender-inflecting: `-o` / `-a`, `-eiro` / `-eira`, `-or` / `-ora`.
- Feminine-fixed: `-ção`, `-dade`, `-agem`.
- Masculine-fixed: `-mento`.
- Common gender: `-ista` (the same form is masculine or feminine).

### Inflection

A generated noun inflects for gender and number using the shared rules above.
The function `NounPTOf(g Gender, n Number, l LengthTypeWords)` resolves `g` and
`n` (with `Any` meaning a random valid choice) and produces the corresponding
form within length category `l`.

## Adjective (Adjetivo)

### Endings

A generated adjective uses one of the following typical pt-PT adjective endings.
The thematic `-o` / `-a` ending is the primary class, and the remaining endings
are the common derivational classes.

- **Thematic `-o` / `-a`**: the most common adjective class in European
  Portuguese (for example `belo` / `bela`, `alto` / `alta`, `novo` / `nova`,
  `rico` / `rica`, `longo` / `longa`). Its masculine form ends in `-o` and its
  feminine form ends in `-a`. It inflects for gender, it is paroxytone, and it
  forms the regular plural by adding `s`.
- **Derivational endings**: `-oso` / `-osa`, `-ável`, `-ível`, `-al`,
  `-ico` / `-ica`, `-ente`, `-ante`, `-ivo` / `-iva`. The `-ante` ending is
  common gender: it is invariable for gender, it is paroxytone, and it forms the
  regular plural by adding `s` (for example `elegante` / `elegante`,
  `elegantes`).

The thematic `-o` / `-a` class is what makes the synthetic superlative `g` to
`gu` hardening reachable: a thematic base ending in `-go` produces forms such as
`longo` to `longuíssimo`. It complements the `c` to `qu` hardening already
reachable through the `-ico` / `-ica` ending (`rico` to `riquíssimo`). Both
adjustments are defined under Orthographic adjustments for the superlative
below.

### Inflection

A generated adjective inflects for gender and number using the shared rules
above, and it supports two degrees through the `Degree` option:

- `Positive`: the base adjective, inflected for gender and number.
- `Superlative`: the synthetic absolute superlative, formed with the suffix
  `-íssimo` / `-íssima`, itself inflected for gender and number
  (`-íssimo`, `-íssima`, `-íssimos`, `-íssimas`).

No other degree is supported. The comparative, the relative superlative, and
the analytic absolute superlative are multi-word forms and are out of scope, as
stated in [wordspt.md](wordspt.md).

### Orthographic adjustments for the superlative

The superlative applies pt-PT orthographic adjustments so that the sound of the
base is preserved before the `-íssimo` suffix:

- A base ending in `-co` changes `c` to `qu` before the suffix, so that the hard
  sound is preserved: `rico` becomes `riquíssimo`.
- A base ending in `-go` changes `g` to `gu` before the suffix, so that the hard
  sound is preserved: `longo` becomes `longuíssimo`.
- The base loses its final vowel before the suffix (`rico` to `riqu-` plus
  `-íssimo`).

The stress of a superlative falls on the suffix (`-íssimo` is proparoxytone and
carries the acute accent), and the accentuation layer marks it accordingly.

## Verb (Verbo)

### Conjugation classes

Every generated verb is regular. There are three regular conjugation classes,
identified by the infinitive ending:

- First conjugation: infinitive in `-ar`, theme vowel `a`. Model verb: `falar`.
- Second conjugation: infinitive in `-er`, theme vowel `e`. Model verb:
  `comer`.
- Third conjugation: infinitive in `-ir`, theme vowel `i`. Model verb:
  `partir`.

Because generated verbs are always regular, the regular paradigm of the chosen
conjugation applies uniformly. The exact ending tables for each conjugation are
embedded as static data sourced from a reference grammar. The first-conjugation
model is given in full below; the second and third conjugations follow the same
paradigm with their own theme vowels.

### Persons

WordsPT supports the persons of current pt-PT usage, without the archaic second
person plural ("vós"):

- Singular: first (`eu`), second (`tu`), third (`ele` / `ela` / `você`).
- Plural: first (`nós`), third (`eles` / `elas` / `vocês`).

There is no second person plural form. The consequences for impossible requests
are defined under Normalization below.

### Paradigm to generate

The following slots are generated. The simple pluperfect indicative is out of
scope.

- Non-finite: impersonal infinitive, inflected personal infinitive, gerund,
  participle.
- Indicative: present, imperfect, preterite (perfect past), future,
  conditional.
- Subjunctive (Conjuntivo): present, imperfect, future.
- Imperative: affirmative and negative.

### First-conjugation model paradigm (falar)

Non-finite forms:

| Form | Value |
|------|-------|
| Impersonal infinitive | `falar` |
| Personal infinitive (eu) | `falar` |
| Personal infinitive (tu) | `falares` |
| Personal infinitive (ele) | `falar` |
| Personal infinitive (nós) | `falarmos` |
| Personal infinitive (eles) | `falarem` |
| Gerund | `falando` |
| Participle | `falado` |

Indicative:

| Person | Present | Imperfect | Preterite | Future | Conditional |
|--------|---------|-----------|-----------|--------|-------------|
| eu (1sg) | `falo` | `falava` | `falei` | `falarei` | `falaria` |
| tu (2sg) | `falas` | `falavas` | `falaste` | `falarás` | `falarias` |
| ele (3sg) | `fala` | `falava` | `falou` | `falará` | `falaria` |
| nós (1pl) | `falamos` | `falávamos` | `falámos` | `falaremos` | `falaríamos` |
| eles (3pl) | `falam` | `falavam` | `falaram` | `falarão` | `falariam` |

The present first person plural `falamos` and the preterite first person plural
`falámos` differ only by the acute accent, and both are produced as shown, in
line with pt-PT usage.

Subjunctive (Conjuntivo):

| Person | Present | Imperfect | Future |
|--------|---------|-----------|--------|
| eu (1sg) | `fale` | `falasse` | `falar` |
| tu (2sg) | `fales` | `falasses` | `falares` |
| ele (3sg) | `fale` | `falasse` | `falar` |
| nós (1pl) | `falemos` | `falássemos` | `falarmos` |
| eles (3pl) | `falem` | `falassem` | `falarem` |

Imperative:

| Person | Affirmative | Negative (verb form) |
|--------|-------------|----------------------|
| eu (1sg) | not applicable | not applicable |
| tu (2sg) | `fala` | `fales` |
| ele (3sg) | `fale` | `fale` |
| nós (1pl) | `falemos` | `falemos` |
| eles (3pl) | `falem` | `falem` |

The imperative has no first person singular. The negative imperative in
Portuguese is periphrastic: it is formed with the particle `não` followed by the
subjunctive form (for example `não fales`), and is therefore a multi-word
construction. Because WordsPT is a single-word generator, the `ImperativeNegative`
mood produces only the verb form used in negative commands (which coincides with
the present subjunctive form for the imperative persons); the negating particle
`não` is not part of the output. This is a confirmed decision, recorded under
Open Items in [wordspt.md](wordspt.md#open-items).

### Selecting the infinitive form

The `Mood` value `Infinitive` covers both the impersonal infinitive and the
inflected personal infinitive:

- When both `Person` and `Number` are `Any`, the generator produces either the
  impersonal infinitive or a random inflected personal infinitive.
- When `Person` or `Number` is specified, the generator produces the inflected
  personal infinitive for the resolved person and number. For the first and
  third person singular, this form is identical to the impersonal infinitive.

### Non-finite inflection

- The gerund is invariable; `Person` and `Number` are ignored for it.
- The verb participle is produced as the invariable masculine singular form (for
  example `falado`); `Person` and `Number` are ignored for it. The adjectival
  inflection of participles is not part of the verb generator; an inflected
  participle used as a modifier is an adjective, produced by the adjective
  functions.

### Verb options

The function `VerbPTOf(m Mood, t Tense, p Person, n Number, l LengthTypeWords)`
resolves each option (with `Any` meaning a random valid choice) and produces the
addressed paradigm slot within length category `l`. Options that do not apply to
the resolved mood are ignored, as defined above and under Normalization.

## Adverb (Advérbio)

### Formation

A generated adverb is a productive `-mente` adverb. It is formed by appending
`-mente` to the feminine singular form of a generated adjective.

### Orthographic rule

When the base adjective carries a graphic accent, that accent is removed in the
adverb, because the primary stress moves to the `-mente` suffix:

- `fácil` gives `facilmente`.
- `rápido` gives the feminine `rápida`, which gives `rapidamente`.

### Invariability

A `-mente` adverb is invariable. It does not inflect for gender, number, or
degree. The function `AdverbPTByLengthType(l LengthTypeWords)` produces such an
adverb within length category `l`, subject to the minimum viable length rule
below.

## Normalization of impossible combinations

WordsPT never returns an error and never panics. When a requested combination is
grammatically impossible, WordsPT normalizes it to the nearest valid form and
produces that form. The rules below are deterministic and are confirmed
decisions, recorded under Open Items in [wordspt.md](wordspt.md#open-items). Each
rule states the linguistic justification for its target.

First, every `Any` option is resolved to a random valid value: `Mood` to a
random valid mood; `Tense` to a random valid tense for the resolved mood;
`Person` and `Number` to a random valid person and number for the resolved mood.
Then the following rules apply.

- **N1. Second person plural.** A requested `Person` of `Second` with a `Number`
  of `Plural`, in any mood, normalizes to `Third` person `Plural`. Justification:
  current pt-PT expresses second person plural address with `vocês`, which
  governs third person plural agreement; this is the nearest existing form.
- **N2. First person singular imperative.** A requested `Person` of `First` with
  a `Number` of `Singular`, in `ImperativeAffirmative` or `ImperativeNegative`,
  normalizes to `First` person `Plural`. Justification: the imperative has no
  first person singular; its first person exists only in the plural (the "let
  us" form). The mood and the person are preserved, and only the number is
  relaxed.
- **N3. Tense not available in the mood.**
  - For the non-finite moods (`Infinitive`, `Gerund`, `Participle`) and for the
    imperative moods (`ImperativeAffirmative`, `ImperativeNegative`), there is no
    tense selection, so a requested `Tense` is ignored. This is not an error.
  - For `Subjunctive`, only `Present`, `Imperfect`, and `Future` exist. A
    requested `Preterite` normalizes to `Imperfect` (the nearest past), and a
    requested `Conditional` normalizes to `Future` (the nearest prospective
    form). Justification: the mood is preserved and the tense maps to the nearest
    existing subjunctive tense by temporal value.
  - For `Indicative`, all five tenses exist, so no normalization is needed.
- **N4. Options that do not apply.** `Person` and `Number` on the gerund and on
  the participle are ignored, as defined under Non-finite inflection. `Gender`
  is not an option of the verb functions. `Degree` is an option of the adjective
  functions only. `Person`, `Tense`, and `Mood` are options of the verb
  functions only.
- **N5. Unviable length.** A length category that cannot contain a valid word of
  the resolved class and slot normalizes upward, as defined below.

## Minimum viable length per class

Length is measured in characters (runes), as defined in [wordspt.md](wordspt.md).
Each class, and each verb paradigm slot, has a minimum viable character count
equal to the shortest stem the phonotactic layer can build plus the shortest
applicable ending and inflection. The concrete per-slot thresholds are derived
from the embedded syllable and ending data.

The normalization rule is: when the maximum of the requested length category is
below the minimum viable character count of the resolved class and slot, WordsPT
normalizes the length upward to the smallest category whose maximum is at least
that minimum. Length is never normalized downward.

Consequences that hold regardless of the concrete data:

- **Adverb.** A `-mente` adverb is the feminine singular adjective base plus the
  six characters `mente`. Its shortest form exceeds the four-character ceiling of
  `SmallLengthWord`, so a `-mente` adverb requested in `SmallLengthWord` always
  normalizes upward to the smallest viable category.
- **Superlative adjective.** The `-íssimo` suffix adds several characters to the
  base, so a `Superlative` adjective requested in `SmallLengthWord` normalizes
  upward.
- **Long verb slots.** Slots such as the first person plural conditional (for
  example `falaríamos`) exceed the ceilings of the smaller categories, so a
  request that pairs such a slot with a small category normalizes upward.

Closed-class functions ignore length entirely and are never subject to this
rule.
