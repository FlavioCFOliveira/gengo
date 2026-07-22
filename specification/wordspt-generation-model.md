# WordsPT Generation Model

## Overview

This document defines the three-layer model that WordsPT uses to build every
open-class pseudo-word. The layers run in a fixed order: phonotactics, then
morphology, then graphic accentuation. Closed-class words bypass this model
entirely and are selected from the curated lists defined in
[wordspt-closed-classes.md](wordspt-closed-classes.md).

This document is a companion to [wordspt.md](wordspt.md) and uses the
terminology defined there.

## Layer 1: Syllabic phonotactics

The first layer builds a stem as an ordered sequence of syllables. Each syllable
has the structure `(onset)(nucleus)(coda)`, where the onset and the coda are
optional and the nucleus is mandatory.

### Syllable inventories

The generator draws each syllable part from a real pt-PT inventory. Each
inventory carries frequency weights so that common sounds appear more often than
rare ones. The concrete members and weights are embedded as static tables,
sourced from authoritative pt-PT references; the categories are fixed as
follows.

- **Onsets**:
  - Single consonants: the standard pt-PT consonant onsets, including the
    digraphs `ch`, `lh`, and `nh`, and the sequences `qu` and `gu` before `e`
    and `i`.
    - Constraint: `lh` and `nh` never appear in word-initial position; they
      occur only between vowels.
  - Consonant clusters: an obstruent followed by a liquid, drawn from the set
    `pr, br, tr, dr, cr, gr, fr, vr, pl, bl, cl, gl, fl`.
- **Nuclei**:
  - Oral vowels: `a, e, i, o, u`, including their open and closed qualities.
  - Oral diphthongs: for example `ai, ei, oi, ui, au, eu, iu, ou`.
  - Nasal vowels: for example `am/an, em/en, im/in, om/on, um/un`, and the
    tilde vowels `a` with tilde and `o` with tilde.
  - Nasal diphthongs: `ao` with tilde, `ae` with tilde, `oe` with tilde
    (written respectively as the standard pt-PT forms). These are the nasal
    diphthongs used, among others, in noun and verb endings.
- **Codas**: restricted to the consonants `s, r, l, m, n, z, x`. No other coda
  is permitted.

### Positional and transition constraints

- **Word position**: onset clusters are permitted in any syllable that begins
  with a consonant cluster; the coda inventory above applies to any syllable,
  subject to the word-final restriction below.
- **Word-final coda**: a stem that will receive a vowel-final class ending does
  not need a coda in its final syllable; when a stem is used with a
  consonant-final class ending (for example `-or`), the ending supplies the
  final coda.
- **Nasal codas**: `m` and `n` in coda position mark nasalization of the
  preceding nucleus and follow pt-PT spelling (`m` before `p` and `b`, and at
  word end; `n` otherwise).
- **Inter-syllable transitions**: only transitions that produce pronounceable
  pt-PT sequences are permitted. A coda and the following onset must form a
  valid pt-PT consonant sequence. The full transition table is embedded as
  static data, sourced from authoritative references.

### Output of Layer 1

Layer 1 outputs a stem together with the index of its tonic syllable. The tonic
syllable is chosen by the generator, which is what allows Layer 3 to apply
graphic accentuation deterministically.

## Layer 2: Morphology

The second layer turns a stem into a word of a specific class by applying the
class ending and then the requested inflection.

The ordered steps are:

1. Attach the class ending to the stem (for example a noun ending or an
   adjective ending).
2. Apply the requested inflection (gender, number, degree, or the verb paradigm
   slot).
3. Apply any orthographic adjustment that the inflection requires (for example
   the consonant change in the absolute superlative `rico` to `riquíssimo`).

The class endings, the inflection rules, and the orthographic adjustments are
defined in full in [wordspt-open-classes.md](wordspt-open-classes.md). Layer 2
may adjust the tonic syllable index when an inflection moves the stress (for
example the `-íssimo` superlative and the `-mente` adverb).

## Layer 3: Graphic accentuation

Because Layer 1 fixes the tonic syllable and Layer 2 tracks any stress shift,
the generator knows the exact stress position of every word. Layer 3 therefore
applies pt-PT graphic accentuation deterministically, in full compliance with
the Acordo Ortográfico da Língua Portuguesa.

### Stress classification

Every word is classified by the position of its tonic syllable:

- **Oxytone** (aguda): stress on the last syllable.
- **Paroxytone** (grave): stress on the second-to-last syllable.
- **Proparoxytone** (esdrúxula): stress on the third-to-last syllable.

### Accentuation rules

The generator applies the following rules. The precise trigger sets are embedded
as static data sourced from the Acordo Ortográfico da Língua Portuguesa; the
rules below define the required behaviour.

- **Proparoxytones**: always receive a graphic accent on the tonic vowel.
- **Oxytones**: receive a graphic accent when the word ends in a tonic `a`, `e`,
  or `o` (optionally followed by `s`), or in `em` or `ens`.
- **Paroxytones**: receive a graphic accent when the word ends in a sound that
  is not among the default paroxytone endings (that is, they are accented when
  ending in `l`, `n`, `r`, `x`, `i`, `is`, `us`, `um`, `uns`, a diphthong, or
  the other pt-PT triggers), following the Acordo Ortográfico.
- **Accent shape**: the acute accent marks an open tonic vowel, and the
  circumflex accent marks a closed tonic vowel. The choice between acute and
  circumflex is determined by the vowel quality recorded during generation.
- **Nasal tilde**: the tilde marks the nasal vowels and the nasal diphthongs
  (`a` with tilde, `o` with tilde, and the diphthongs `ao`, `ae`, `oe` with
  tilde). The tilde carries the stress mark where applicable, so no additional
  acute or circumflex is added to a nasal vowel that the tilde already marks.
- **Cedilla**: the cedilla is placed under `c` (producing `ç`) only before `a`,
  `o`, or `u`, and never in word-initial position. Before `e` or `i`, the `s`
  sound is written with `c` or `s`, never with a cedilla.
- **Grave accent**: the grave accent occurs only in the crasis contraction
  (`a` plus `a` written as `à`). WordsPT does not generate crasis in open-class
  words; the grave accent therefore appears only in the relevant closed-class
  contractions listed in
  [wordspt-closed-classes.md](wordspt-closed-classes.md).
- **Diaeresis**: the diaeresis is not used, in line with the Acordo
  Ortográfico, which abolished it in pt-PT.

### Determinism

Given a stem, its tonic syllable index, its class ending, and its inflection,
Layer 3 produces exactly one correctly accented spelling. No accentuation
decision is left to chance, and no accentuation is guessed. The correctness of
this layer is validated by acceptance criterion 10 in [wordspt.md](wordspt.md).
