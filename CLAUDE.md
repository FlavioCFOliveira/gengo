# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Purpose

**gengo** is a Go library designed to provide a simple and practical way to generate random data structures. It offers an intuitive, straightforward API for creating random values of various types (strings, numbers, dates, booleans, words) without complex configuration or boilerplate code.

### Design Philosophy

- **Simplicity First**: One function call to get random data - no initialization, no configuration
- **Practical Defaults**: Sensible defaults that cover 90% of use cases out of the box
- **Intuitive Naming**: Functions follow the pattern `Type()` for full range, `TypeBetween(min, max)` for custom ranges
- **Zero Dependencies**: Only uses Go standard library - no external dependencies to manage

### Performance Philosophy

- **Execution Speed**: Implementations prioritize the fastest execution path, leveraging hardware capabilities where possible
- **Minimal Allocations**: Reduce memory allocations to the absolute minimum - prefer stack over heap
- **Memory Efficiency**: Keep memory footprint low; avoid unnecessary buffering or temporary structures
- **No Bloat**: Avoid abstractions that add overhead; prefer direct, efficient implementations

## Collaboration Rules (CRITICAL)

### Decision Authority

**You are NOT authorized to make decisions on your own.** Whenever instructions are insufficient, unclear, non-specific, or non-concrete, or whenever contradictions or ambiguities exist, you MUST ALWAYS ASK the user how to proceed.

When asking the user:
- Provide multiple options (a, b, c, ...) and clearly indicate which option is your recommendation.
- When multiple clarifications are needed, ask the questions sequentially (one at a time), not all at once.

**Boundary between acting and asking:** Obvious, low-risk corrections (for example, a pre-existing bug with an unambiguous fix) proceed immediately; any decision that changes the scope, expected behavior, architecture, or requirements requires asking the user first.

### Never Guess

All interactions in this project must be based EXCLUSIVELY on verified knowledge. Never guess the intended answer. When the information you have is insufficient, look for answers in official or authoritative sources: specifications, RFCs, papers, books, or reference authors in the field. Use the Knowledge Graph as the primary source of information — both to consult it and to record the relationships you discover.

### Measure to Decide

Whenever you need to evaluate performance, completeness, or correctness, you MUST ALWAYS collect empirical evidence from the project to determine the answer. Decisions must ALWAYS be made empirically — never by assumption.

## Documentation Standards

### Language

All project documentation (including this `CLAUDE.md`) MUST be written in flawless English — free of orthographic, grammatical, and syntactic errors. Use clear, simple, and unambiguous technical language aimed at human readers.

### Accuracy

Documentation must be precise and faithful to the code. It must never describe behavior that diverges from the actual implementation.

## Development Workflow

The workflow MUST always follow these phases, in this exact order:

1. **Specify** — Define the requirement, scope, and acceptance criteria.
2. **Implement** — Write the code.
3. **Test** — Validate correctness and performance with evidence.
4. **Document** — Update documentation to reflect the change.

### Self-Contained Development Policy

Every development cycle MUST be self-contained. NEVER deliver only part of a task — each cycle must produce a complete, usable result. When new needs are discovered during a task (that were not previously anticipated), they must be resolved within the same development cycle, as immediately as possible: add new tasks and develop them right away.

- All code and development must, as a rule, be **full-fledged** (complete and ready to use).
- NEVER create skipped tests (for example, using `t.Skip`).
- Whenever you find a pre-existing bug, fix it immediately, and then resume the work you were doing when you found it.

### Regression Prevention

Whenever you identify a bug, create the regression tests needed to guarantee that the same bug cannot reappear as a result of future development.

### Production-Grade Orientation

**ALL** the actions you perform — development, fixes, evaluations, analyses, audits, and any others — must be treated with production-level rigor. Apply both maximum knowledge and maximum effort to ensure that every deliverable is ready for production use.

## Decision Framework

To decide the expected result of the project — whether in evaluations and audits or during code implementation — follow this priority order: **correct → safe → fast.**

1. **Is it correct?** Does the result comply with the objective, the project specification, and the applicable authoritative sources (RFCs, standards, etc.)?
2. **Is it safe?** Does the decision or task avoid introducing any characteristic or behavior that compromises the safe use of the deliverable?
3. **Is it fast?** Is it as fast as it can be without compromising correctness or safety? What can be done to maximize the deliverable's performance?

If these criteria conflict, or you have difficulty following them, immediately ask the user how to proceed, presenting the possible options.

## Roadmap

**Name:** gengo

## Planning and Task Execution

- For any operation involving Tasks or Sprints, use the `roadmap-manager` skill.
- Use the `rmp` CLI (the roadmap management tool available on the system) to plan and coordinate task execution.
- Treat `rmp` as the **single source of truth** for planning and execution of tasks in this project. No other mechanism may be used for this purpose.
- Use the Knowledge Graph to understand the project, its components, and the relationships between them, so you can more easily identify the scope and impact of each task.

### Planning

- Carefully evaluate the scope proposed by the user and determine first whether it justifies multiple development phases. Each phase must correspond to a solid deliverable.
- Every task must have a clear and objective definition of:
  - Goals
  - Functional requirements
  - Technical requirements
  - Acceptance criteria that confirm the task can be closed
- When a task is completed, it must be closed with a short summary describing what was done.
- Phases are represented as **Sprints** in `rmp` and serve to group tasks.
- When the work requires multiple phases, planning MUST be done in two distinct steps:
  1. First, define which phases (sprints) are needed and the scope (objective) of each.
  2. Only then, sprint by sprint, define the tasks of each sprint.

  In both steps, use `rmp` as the single source of truth.
- Use the Knowledge Graph to identify the highest-gain/highest-impact tasks, the foundational tasks, and the tasks that unblock other tasks or features, in order to optimize the execution order.
- **Prioritization:** by default, always work from the highest-gain/impact tasks toward the least essential ones. Foundational tasks and tasks that unblock others are always prioritized.
- When a task is too large to be executed in a single pass by an AI agent (such as Claude Code), subdivide it into parts, respecting the principles already defined (in particular, the self-contained-task principle).

### Task Execution

Task execution is the natural continuation of planning. Always use `rmp` and follow this sequence:

1. Check whether there is an open, unfinished task to continue.
2. Identify the next task.
3. Understand the goal of the task starting (description, functional and technical requirements).
4. Determine the most appropriate subagent and delegate its execution.
5. Always validate the acceptance criteria before closing a task.
6. Close the task with a short summary of what was done.
7. After closing the task, and before moving to the next, create a git commit following best practices, explaining what was done.
8. Update the Knowledge Graph.

Execution notes:

- Whenever possible, adapt the model and its reasoning-effort level to the requirements of each individual operation within a task.
- Task and sprint execution is **sequential**.
- Evaluations and audits may run in parallel, but parallel execution must **ALWAYS be authorized by the user**.
  - Even when authorized, **NEVER run more than 2 (two) simultaneously**. Plan all the evaluations/audits that are needed, but execute at most 2 at a time: as one finishes, start the next, always keeping the limit of 2 in parallel.

## Knowledge Graph

- Manage the Knowledge Graph (KG) with the help of the `knowledge-authority` skill.
- Use the "Graph" features of `rmp` (Groadmap) to create, maintain (update), and query a project knowledge graph.
- This graph **MUST CONTAIN EVERYTHING** useful to know about the project. Examples:
  - which features exist, and where they are specified and implemented;
  - which tests exist and what they test;
  - which components exist, how they relate, and the dependencies between them;
  - in which git commit each feature was specified, implemented, and tested;
  - the `rmp` tasks and their links to components.
- The graph **MUST ALWAYS BE UPDATED on every git commit**, recording the changes to the graph objects. Each node and edge update must identify the corresponding commit and its date.
- **This graph is the absolute truth about the project.** Keep it as up to date as possible so that, before having to read files, you can query the graph and obtain what you need.
- Create the nodes and edges that make the most sense for the project. Use the graph together with tasks and sprints to coordinate the work.

## Working Team (Subagents)

- You have at your disposal a team composed of all available subagents (global, user, or project).
- Use them collaboratively and complementarily, so that each task is completed with maximum confidence, effectiveness, and assertiveness.
- Each subagent must contribute proactively with its specialty.

## Build and Test Commands

```bash
# Run all tests
go test .

# Run a single test by name
go test -run TestInt8 .

# Run benchmarks
make bench
# Or directly:
go test -benchmem -run=^$ -bench ^Benchmark -benchtime 5s .

# Format code
go fmt .
# Or:
make format

# Lint code
make check          # Runs gofmt -l and golangci-lint

# Security check
make vulncheck      # Runs govulncheck ./...

# Install dev dependencies
make dev-install    # Installs govulncheck and golangci-lint
```

## Architecture

**gengo** is a Go module for generating random data. It uses `math/rand/v2` (Go 1.22+) which has automatic seed initialization.

### Code Organization

The package follows a flat file structure where each file handles a specific data type category:

- **init.go** - Package initialization with date range constants (`dtMin`, `dtMax`, `dtUnixMin`, `dtUnixMax`) and error variables (`ErrorNilFunc`, `ErrorFuncNilResult`)
- **strings.go** - Random string generation with predefined character set constants (`AllChars`, `Alphanumeric`, `Alphabetic`, `Numeric`, `Hexadecimal`, `Symbols`)
- **numbers.go** - Random number generation for all numeric types (int8/16/32/64, uint8/16/32/64, float32/64, complex64/128). Each type has a `Type()` function for full range and `TypeBetween(min, max)` for custom ranges
- **dates.go** - Random `time.Time` generation with `Date()`, `UnixDate()`, and `DateBetween(start, end)`
- **bool.go** - Boolean generation via `Bool()`
- **words.go** - Word generation with length categories (Small/Medium/Big) via `LengthTypeWords` enum

### Key Patterns

1. **Range Handling**: Functions accepting min/max parameters swap values if `min > max` to ensure valid ranges
2. **Overflow Prevention**: Integer functions use wider types for intermediate calculations (e.g., `int16` for `int8` math)
3. **String Generation**: `String(length, sourceChars)` is the core function; other string functions are convenience wrappers using predefined constants
4. **Testing**: Tests use a `loop` variable (1000 iterations) to validate randomness through repeated execution rather than asserting specific values

### Dependencies

- No external runtime dependencies
- Uses only Go standard library (`math/rand/v2`, `math`, `time`, `errors`)
- Dev tools: `golangci-lint`, `govulncheck`

### CI/CD

GitHub Actions workflow (`.github/workflows/go.yml`) runs on every push:
- Dependency check (`go mod tidy -diff`)
- Vulnerability scan (`govulncheck`)
- Linting (`golangci-lint`)
- Format check (`gofmt`)
- Tests (`go test -v`)

## Git Workflow

**IMPORTANT:** All commits must exclusively be authored by `flaviocfo <flaviocfo@gmail.com>`. Do not include Co-Authored-By tags or any other author attribution in commits. The user will handle commit authorship.

## Commit Policy

**All commits must be authored exclusively by flaviocfo.** When creating commits, I will use your configured git identity. Do not add Co-Authored-By tags or other attribution that would change the commit author from flaviocfo.
