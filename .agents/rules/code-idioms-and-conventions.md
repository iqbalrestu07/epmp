---
trigger: always_on
---

## Code Idioms and Conventions

### Universal Principle

**Write idiomatic code for the target language:**

- Code should look natural to developers familiar with that language
- Follow established community conventions, not personal preferences
- Use language built-ins and standard library effectively
- Apply language-appropriate patterns (don't force patterns from other languages)

### Idiomatic Code Characteristics

- Leverages language features (don't avoid features unnecessarily)
- Follows language naming conventions
- Uses appropriate error handling for language (exceptions vs Result types)
- Applies established community patterns

### Avoid Cross-Language Anti-Patterns

- ❌ Don't write "Java in Go" (no `IRepository`, no getters/setters for every field, no exceptions-as-panics)
- ❌ Don't write "Go in TypeScript" (no `err, result` tuples — use exceptions/`ApiError`)
- ❌ Don't force class hierarchies in React — components are functions + hooks
- ✅ Learn and apply language-specific idioms

### Language-Specific Idioms

This file defines the universal principle. Language-specific idiom files provide concrete patterns and tooling choices.

| Language / Framework                 | Idiom File                         | When to Load                                     |
| ------------------------------------ | ---------------------------------- | ------------------------------------------------ |
| **Go** (`backend/`, `tools/`)        | @go-idioms-and-patterns.md         | Any Go change                                    |
| **TypeScript / React** (`frontend/`) | @typescript-idioms-and-patterns.md | Any frontend change                              |
| **EPMP-specific structure**          | @epmp-module-patterns.md           | Module/feature creation or modification (always) |

> Naming conventions (packages, files, DTOs, events, errors, tests) are defined in `tools/epmp-ai/CONVENTIONS.md` and take precedence over generic community defaults when they differ.

### Related Principles

- Core Design Principles @core-design-principles.md
- Code Completion Mandate @code-completion-mandate.md
- EPMP Module Patterns @epmp-module-patterns.md
