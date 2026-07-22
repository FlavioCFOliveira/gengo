# Knowledge Graph Model — gengo

Canonical description of the `gengo` project knowledge graph (a Label-Property
Graph managed through `rmp graph`, roadmap `gengo`). This file is the source of
truth for the graph's shape: its node labels, edge types, and properties. Keep it
in conformance with the live graph.

## Provenance (every node and edge)

- `gitCommit` — full commit hash when the element was last confirmed.
- `gitDate` — that commit's ISO date (`YYYY-MM-DD`).

## Node labels

| Label | Identity property | Other properties |
|-------|-------------------|------------------|
| `Package` | `name` | `module`, `path` |
| `File` | `path` | `role` (`source` \| `test`) |
| `Feature` | `name` | `status` (`released` \| `in-progress` \| `specified`), `description` |
| `Component` | `name` | `kind`, `description` |
| `Type` | `name` | `kind` (`struct` \| `enum` \| `interface`), `exported` (bool) |
| `Function` | `name` | `signature`, `exported` (bool) |
| `Spec` | `path` | `title` |
| `Sprint` | `id` | `title`, `status` |
| `Task` | `id` | `title`, `type`, `status` |
| `Commit` | `hash` | `type` (conventional-commit type), `message` |

## Edge types

| Edge | From → To | Meaning |
|------|-----------|---------|
| `IN_PACKAGE` | `File` → `Package` | file belongs to the Go package |
| `IMPLEMENTED_IN` | `Feature` → `File` | feature's code lives (in part) in this file |
| `HAS_COMPONENT` | `Feature` → `Component` | feature is decomposed into this component |
| `SPECIFIED_BY` | `Feature` → `Spec` | feature is specified by this document |
| `DEFINED_IN` | `Component` \| `Type` \| `Function` → `File` | symbol/component defined in this file |
| `PART_OF` | `Task` → `Sprint` | task belongs to the sprint |
| `DELIVERS` | `Task` → `Component` \| `Feature` | task delivers this work |
| `COMMITTED_IN` | `Task` → `Commit` | task's work recorded in this commit |
| `TOUCHES` | `Commit` → `File` | commit adds/modifies this file |

## Notes

- Sprints/Tasks mirror the `rmp` roadmap (owned by the roadmap-manager); the graph
  holds them to link planning to code.
- The graph is enriched incrementally: per-`Function` granularity is added on demand
  (fallback reads), not exhaustively up front.
- Bootstrapped 2026-07-21 (commit `f211c44`), focused scope: package skeleton +
  WordsPT feature (Sprint 7, task #18).
