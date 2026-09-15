---
trigger: always_on
---

## Documentation Principles

### Self-Documenting Code

**Clear naming reduces need for comments:**

- Code shows WHAT is happening  
- Comments explain WHY it's done this way

**When to comment:**

- Complex business logic deserves explanation  
- Non-obvious algorithms (explain approach)  
- Workarounds for bugs (link to issue tracker)  
- Performance optimizations (explain trade-offs)

### Documentation Levels

1. **Inline comments:** Explain WHY for complex code  
2. **Function/method docs:** API contract (parameters, returns, errors)  
3. **Module/package docs:** High-level purpose and usage  
4. **README:** Setup, usage, examples  
5. **Architecture docs:** System design, component interactions, use valid mermaid diagram

### PRD & TSD Update Mandate

**When implementing features that fulfill or change planned milestones, the corresponding PRD and TSD documents MUST be updated as part of the same work.**

**Triggers for update:**
- New collection or schema change → update TSD (scope, schema section) + PRD (milestone status, domain model if affected)
- Milestone status change (Planned → In Progress → Done) → update PRD milestone table
- New endpoint or API change → update TSD (API section)
- New FR (functional requirement) completed → update PRD (mark done) + TSD (implementation details)

**What to update:**
| Document | When | What |
|---|---|---|
| **PRD** | Milestone status changes | Update milestone table status, bump doc version |
| **PRD** | New collection/entity | Add to Domain Model table if it's a `ref_*` collection |
| **TSD** | Feature implemented | Update scope, mark items done, add implementation notes |
| **TSD** | Schema change | Update collection/field documentation |

> **Rule:** Never ship a feature that changes the schema or completes a milestone without updating PRD/TSD. Treat documentation as part of the definition of done.

### Related Principles
- Core Design Principles @core-design-principles.md
- Code Organization Principles @code-organization-principles.md
- Schema Documentation Mandate @schema-documentation-mandate.md
