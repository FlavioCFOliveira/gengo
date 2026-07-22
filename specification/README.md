# gengo Functional Specification

This folder is the single source of truth for the functional specification of
the gengo project. Each file describes a distinct functional responsibility
area. This README is the index and the bridge to every specification file.

The specification serves two audiences at once: human stakeholders who need to
understand what the system does, and AI agents that use the specification as a
source of truth for code generation, review, and decision-making. Every file
follows a consistent structure and uses consistent terminology.

## How to read this specification

- Start from this index to find the feature you need.
- Each feature has a primary document and, where the feature is large, one or
  more companion documents that hold the detailed rules. The primary document
  links to its companions.

## Specification index

### WordsPT

The WordsPT feature generates plausible pseudo-words that follow the
morphological and phonotactic patterns of European Portuguese (pt-PT), and
selects real pt-PT single words for the closed word classes.

| File | Scope |
|------|-------|
| [wordspt.md](wordspt.md) | Primary WordsPT document: overview, definitions, scope, word class taxonomy, public API (enumerations and functions), cross-cutting semantics, constraints, acceptance criteria, and open items. |
| [wordspt-generation-model.md](wordspt-generation-model.md) | The three-layer generation model: syllabic phonotactics, morphology, and deterministic graphic accentuation. |
| [wordspt-open-classes.md](wordspt-open-classes.md) | The generated open classes (noun, adjective, verb, adverb): endings, inflection rules, the regular verb paradigm, normalization of impossible combinations, and minimum viable length per class. |
| [wordspt-closed-classes.md](wordspt-closed-classes.md) | The curated closed classes (article, pronoun, numeral, preposition, conjunction, interjection): membership rules and representative members. |

## Recorded decisions

Decisions that were raised as open items during specification have been resolved
by the user and are retained as a decision record in [wordspt.md](wordspt.md)
under "Open Items". They cover: the two randomness surfaces of WordsPT (the
package-level functions and the mirrored `*Generator` methods), the authoritative
sourcing of the class-frequency weight table (scheduled for the orchestration
phase), the verb normalization matrix, and the single-word resolution of the
negative imperative.
