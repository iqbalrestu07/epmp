---
trigger: model_decision
description: When writing TypeScript code, reviewing TypeScript idioms, or working on any TypeScript project (backend, frontend, or CLI)
---

## TypeScript Idioms and Patterns

### Core Philosophy

TypeScript's type system is your documentation, your test, and your specification — all at once. Make the type system encode the invariants of your domain so that invalid states are unrepresentable. Lean into the compiler.

> **Scope:** This file covers TypeScript-specific _type system and language idioms_. For quality commands, see `code-completion-mandate.md`. For logging, see `logging-and-observability-principles.md`. For EPMP feature layout (api/hooks/types/schema/pages), TanStack Query, and the `@/services/api` wrapper, see `epmp-module-patterns.md`.
>
> **This project:** `frontend/` is React 18 + Vite + TypeScript 5 with `npm`. Path alias `@/` → `src/`. Typecheck is `tsc -b` (run by `npm run build`).

---

### Strict Mode — Non-Negotiable

**Always enable strict mode** in `tsconfig.json`:

```json
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true
  }
}
```

These flags catch the majority of runtime errors at compile time. Never disable them on a per-file basis without a `// STRICT-DISABLE:` comment explaining the rationale.

---

### Type System Idioms

1. **`unknown` over `any` — always**

   ```typescript
   // ✅ Forces narrowing before use
   function parse(data: unknown): User {
     if (!isUser(data)) throw new Error("Invalid user shape");
     return data;
   }

   // ❌ Disables the type checker entirely
   function parse(data: any): User {
     return data;
   }
   ```

2. **Use `readonly` to enforce immutability at compile time**

   ```typescript
   interface TaskState {
       readonly id: string;
       readonly items: readonly Task[];
   }

   // For function params that should not be mutated
   function process(tasks: readonly Task[]): Summary { ... }
   ```

3. **Discriminated unions for type-safe state machines**

   ```typescript
   type AsyncState<T> =
     | { status: "idle" }
     | { status: "loading" }
     | { status: "success"; data: T }
     | { status: "error"; error: Error };

   // Exhaustive handling — compiler catches missing cases
   function render(state: AsyncState<User>): string {
     switch (state.status) {
       case "idle":
         return "Waiting...";
       case "loading":
         return "Loading...";
       case "success":
         return state.data.name; // data is typed
       case "error":
         return state.error.message;
     }
   }
   ```

4. **Const assertions for literal types**

   ```typescript
   const ROLES = ["admin", "editor", "viewer"] as const;
   type Role = (typeof ROLES)[number]; // 'admin' | 'editor' | 'viewer'
   ```

5. **Type narrowing — use type guards instead of `as` casts**

   ```typescript
   // ✅ Type guard — safe narrowing
   function isError(value: unknown): value is Error {
     return value instanceof Error;
   }

   // ❌ Type assertion — bypasses type checker
   const err = value as Error;
   ```

6. **Never use non-null assertion `!` in production code**

   ```typescript
   // ❌ Hides a potential null/undefined bug
   const name = user!.profile!.name;

   // ✅ Explicit handling
   const name = user?.profile?.name ?? "Anonymous";
   ```

7. **`satisfies` operator for type-checked object literals (TS 4.9+)**
   ```typescript
   // ✅ satisfies: compile-checked against interface, type stays as literal
   const config = {
     endpoint: "/api/tasks",
     retries: 3,
   } satisfies ApiConfig;
   // config.retries is typed as `3` (literal), not `number` — narrower and safer
   ```

---

### Null Safety

1. **Prefer `??` (nullish coalescing) over `||` for default values**

   ```typescript
   // ✅ Only falls back for null/undefined
   const count = input.count ?? 0;

   // ❌ Also falls back for 0, '', false
   const count = input.count || 0;
   ```

2. **Use optional chaining `?.` for safe navigation**

   ```typescript
   const city = user?.address?.city;
   ```

3. **Distinguish `undefined` (absence) from `null` (explicit empty)**
   - Use `undefined` for optional fields
   - Use `null` only when you need to represent "intentionally empty" on the wire (JSON APIs)

---

### Async/Await

> This section covers TypeScript-specific async idioms.

1. **Always `await` or handle returned Promises — no floating promises**

   ```typescript
   // ❌ Fire-and-forget — errors are silently swallowed
   sendEmail(user);

   // ✅ Awaited
   await sendEmail(user);

   // ✅ Intentionally fire-and-forget needs explicit void annotation
   void sendEmail(user); // still logs errors internally
   ```

2. **Use `Promise.all` for concurrent independent operations**

   ```typescript
   // ✅ Concurrent — total time = max(individual times)
   const [user, tasks] = await Promise.all([getUser(id), getTasks(id)]);

   // ❌ Sequential — total time = sum of individual times
   const user = await getUser(id);
   const tasks = await getTasks(id);
   ```

3. **Use `Promise.allSettled` when partial failure is acceptable**

   ```typescript
   const results = await Promise.allSettled(notifications.map(send));
   const failed = results.filter((r) => r.status === "rejected");
   ```

4. **Never mix `async/await` with raw `.then()/.catch()` chains in the same function**

---

### Runtime Validation at Boundaries

**All data crossing a system boundary must be validated at runtime**, not just typed.

```typescript
import { z } from "zod";

// Define schema as the single source of truth
const CreateTaskSchema = z.object({
  title: z.string().min(1).max(200),
  priority: z.enum(["low", "medium", "high"]),
  dueDate: z.string().datetime().optional(),
});

// Infer the TypeScript type from the schema — no duplication
type CreateTaskRequest = z.infer<typeof CreateTaskSchema>;

// Validate at the API boundary
function parseCreateTask(body: unknown): CreateTaskRequest {
  return CreateTaskSchema.parse(body); // throws ZodError on invalid input
}
```

- Use `zod` for runtime schema validation at API ingress and external API egress
- Never use TypeScript's `as` operator as a substitute for runtime validation
- Validate on ingress; trust validated types thereafter

---

### Centralized HTTP Client

All requests go through `frontend/src/services/api.ts` (`BASE_URL=/api/v1`, JWT + `X-Organization-ID` headers, refresh-token handling). Never call `fetch`/`axios` directly from a component, hook, or feature `api/` file — import the wrapper. See `epmp-module-patterns.md` for the feature `api/` → `hooks/` layering.

---

### Module and Export Patterns

1. **Prefer named exports over default exports**

   ```typescript
   // ✅ Named — explicit, refactor-safe, IDE-friendly
   export function createTask() { ... }
   export type { Task };

   // ❌ Default — ambiguous import names, harder to auto-import
   export default function createTask() { ... }
   ```

2. **Avoid barrel re-exports that create circular dependency risk**
   - Use feature `index.ts` files only for the public API of a feature, not as catch-all re-exports

3. **Import type separately to avoid bundling runtime artifacts**
   ```typescript
   import type { Task } from "./types";
   ```

---

### Testing

> Test naming, file conventions, and pyramid proportions are defined in `testing-strategy.md`. This section covers TypeScript-specific tooling.

1. **Type your test doubles against the real interface** — never use `as any` in test doubles

   > `frontend/` has no unit-test runner configured today (only Playwright E2E). Adding one is an architectural decision — record an ADR first. Until then, keep logic in pure functions and write hand-written fakes that satisfy the interface:

   ```typescript
   const fakeApi: ReservationApi = {
     create: async (input) => ({ id: "r-1", ...input }),
     getById: async () => null,
   };
   ```

2. **Assert on error types, not just error messages**

   ```typescript
   await expect(service.create(invalid)).rejects.toThrow(ZodError);
   ```

3. **Use `satisfies` operator in tests for type-checked fixtures**
   ```typescript
   const fixture = {
     id: "abc",
     title: "Test task",
   } satisfies Task;
   ```

---

### Formatting and Static Analysis

| Tool                           | Purpose                                  | Notes                                                                                                                  |
| ------------------------------ | ---------------------------------------- | ---------------------------------------------------------------------------------------------------------------------- |
| `tsc -b`                       | Full type checking (via `npm run build`) | Must pass zero errors                                                                                                  |
| `eslint .`                     | Lint rules + style (`npm run lint`)      | No config committed yet — adding one needs `@typescript-eslint/recommended-type-checked` + `eslint-plugin-react-hooks` |
| `prettier`                     | Canonical formatting                     | Not configured yet; match existing file style (2 spaces, double quotes, semicolons)                                    |
| `npm audit --audit-level=high` | Dependency CVE scanning                  | Run in CI; fail on high severity                                                                                       |

See `code-completion-mandate.md` for the exact commands to run before committing.

---

### Related Principles

- Code Idioms and Conventions @code-idioms-and-conventions.md
- Testing Strategy @testing-strategy.md
- Error Handling Principles @error-handling-principles.md
- Security Principles @security-principles.md
- Dependency Management Principles @dependency-management-principles.md
