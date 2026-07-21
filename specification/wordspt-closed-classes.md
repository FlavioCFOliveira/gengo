# WordsPT Closed Classes

## Overview

This document defines the six closed classes that WordsPT does not generate but
selects from curated lists of real pt-PT single words: article, pronoun,
numeral, preposition, conjunction, and interjection.

This document is a companion to [wordspt.md](wordspt.md) and uses the
terminology defined there.

## General rules for closed classes

- **Real words only.** Every member of every closed-class list is a real pt-PT
  word. Unlike the open classes, closed-class output is never invented.
- **Single words only.** Every member is a single word. Multi-word locutions are
  excluded from every closed class, in line with [wordspt.md](wordspt.md). This
  excludes prepositional locutions, conjunctional locutions, and compound
  numbers, among others.
- **pt-PT spelling.** Every member uses European Portuguese spelling under the
  Acordo Ortográfico da Língua Portuguesa. Where pt-PT and Brazilian spellings
  differ, the pt-PT form is used (for example `dezasseis`, not `dezesseis`;
  `connosco`, not `conosco`).
- **Exhaustive and sourced.** Each list is exhaustive for its class and is
  curated from authoritative pt-PT references. The lists below give
  representative members to fix the intended scope; the definitive, exhaustive
  membership is finalized in the implementation data file and validated against
  the cited authorities. No member is guessed.
- **Length is ignored.** Closed-class functions ignore the length concept
  entirely, as stated in [wordspt.md](wordspt.md).
- **Lowercase.** Every member is output in lowercase, with correct accents.

## Article (Artigo)

Selected by `ArticlePT()`.

Subtypes and representative members:

- Definite: `o`, `a`, `os`, `as`.
- Indefinite: `um`, `uma`, `uns`, `umas`.
- Single-word contractions of a preposition with an article, for example `do`,
  `da`, `dos`, `das`, `no`, `na`, `nos`, `nas`, `pelo`, `pela`, `pelos`,
  `pelas`, `ao`, `aos`, `à`, `às`, `num`, `numa`, `dum`, `duma`.

The contraction members `à` and `às` carry the grave accent, which marks the
crasis contraction of the preposition `a` with the article `a` or `as`.

## Pronoun (Pronome)

Selected by `PronounPT()`.

The list is exhaustive across all pronoun subtypes and includes the
archaic-but-real forms related to `vós`. Subtypes and representative members:

- Personal, tonic: `eu`, `tu`, `ele`, `ela`, `nós`, `vós`, `eles`, `elas`,
  `mim`, `ti`, `si`, `connosco`, `convosco`, `comigo`, `contigo`, `consigo`.
- Personal, atonic oblique: `me`, `te`, `se`, `lhe`, `lhes`, `nos`, `vos`, `o`,
  `a`, `os`, `as`.
- Possessive: `meu`, `minha`, `teu`, `tua`, `seu`, `sua`, `nosso`, `nossa`,
  `vosso`, `vossa`, and their plurals.
- Demonstrative: `este`, `esta`, `esse`, `essa`, `aquele`, `aquela`, `isto`,
  `isso`, `aquilo`, and their inflected forms.
- Indefinite: `algum`, `alguma`, `nenhum`, `todo`, `outro`, `muito`, `pouco`,
  `tudo`, `nada`, `alguém`, `ninguém`, `cada`, `qualquer`.
- Relative: `que`, `quem`, `qual`, `quais`, `cujo`, `cuja`, `onde`, `quanto`.
- Interrogative: `quem`, `que`, `qual`, `quais`, `quanto`, `quanta`, `quantos`,
  `quantas`.

## Numeral (Numeral)

Selected by `NumeralPT()`.

Only single-word numerals are included. Compound numbers written with more than
one word (for example `vinte e um`) are excluded. Subtypes and representative
members:

- Cardinals: `zero`, `um`, `dois`, `três`, `quatro`, `cinco`, `seis`, `sete`,
  `oito`, `nove`, `dez`, `onze`, `doze`, `treze`, `catorze`, `quinze`,
  `dezasseis`, `dezassete`, `dezoito`, `dezanove`, `vinte`, `trinta`,
  `quarenta`, `cinquenta`, `sessenta`, `setenta`, `oitenta`, `noventa`, `cem`,
  `mil`, `milhão`, `bilião`.
- Ordinals: `primeiro`, `segundo`, `terceiro`, `quarto`, `quinto`, `sexto`,
  `sétimo`, `oitavo`, `nono`, `décimo`, `vigésimo`, `trigésimo`, `centésimo`,
  `milésimo`.

The pt-PT spellings `dezasseis`, `dezassete`, and `catorze` are used, not the
Brazilian `dezesseis`, `dezessete`, and `quatorze`.

## Preposition (Preposição)

Selected by `PrepositionPT()`.

The list is the set of essential (simple) pt-PT prepositions: `a`, `ante`,
`após`, `até`, `com`, `contra`, `de`, `desde`, `em`, `entre`, `para`,
`perante`, `por`, `sem`, `sob`, `sobre`, `trás`.

Prepositional locutions (for example `de acordo com`) are excluded.

## Conjunction (Conjunção)

Selected by `ConjunctionPT()`.

Only single-word conjunctions are included. Conjunctional locutions (for example
`à medida que`) are excluded. Subtypes and representative members:

- Coordinating: `e`, `nem`, `mas`, `ou`, `logo`, `pois`, `porém`, `contudo`,
  `todavia`.
- Subordinating: `que`, `se`, `porque`, `quando`, `enquanto`, `embora`,
  `conforme`, `como`, `conquanto`, `consoante`, `segundo`, `portanto`.

## Interjection (Interjeição)

Selected by `InterjectionPT()`.

Representative members: `ah`, `oh`, `ó`, `olá`, `ui`, `ai`, `oxalá`, `viva`,
`bravo`, `chiu`, `caramba`.

Interjectional locutions (for example `ai de mim`) are excluded.
