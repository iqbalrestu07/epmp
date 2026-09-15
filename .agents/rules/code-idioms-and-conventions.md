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

- ❌ Don't write "Java in Python" or "C in Go"  
- ❌ Don't force OOP patterns in functional languages  
- ❌ Don't avoid language features because they're "unfamiliar"  
- ✅ Learn and apply language-specific idioms

### Language-Specific Idioms

This file defines the universal principle. Language-specific idiom files provide concrete patterns and tooling choices.

| Language / Framework | Idiom File                         | When to Load                                    |
| -------------------- | ---------------------------------- | ----------------------------------------------- |
| **TypeScript / JS**  | @typescript-idioms-and-patterns.md | TypeScript or JavaScript projects               |
| **Directus**         | @directus-extension-patterns.md    | Directus endpoint extensions (this project)     |

> This project uses **plain ESM JavaScript** (not TypeScript). The TypeScript idioms file is retained for async/await and null-safety patterns that apply to JavaScript. For Directus-specific conventions, always defer to `directus-extension-patterns.md`.

### Related Principles
- Core Design Principles @core-design-principles.md
- Code Completion Mandate @code-completion-mandate.md
- Directus Extension Patterns @directus-extension-patterns.md
