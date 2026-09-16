---
trigger: model_decision
description: When implementing logging, working with loggers, or setting up observability for operations (API handlers, database queries, background jobs, external API calls)
---

## Logging and Observability Principles

> **⚠️ Prerequisite:** All operations MUST be logged per Logging and Observability Mandate @logging-and-observability-mandate.md. This guide provides implementation patterns only.

### Logging Standards

#### Log Levels (Standard Priority)

Use consistent log levels across all services:

| Level     | When to Use                             | Examples                                                 |
| --------- | --------------------------------------- | -------------------------------------------------------- |
| **TRACE** | Extremely detailed diagnostic info      | Function entry/exit, variable states (dev only)          |
| **DEBUG** | Detailed flow for debugging             | Query execution, cache hits/misses, state transitions    |
| **INFO**  | General informational messages          | Request started, task created, user logged in            |
| **WARN**  | Potentially harmful situations          | Deprecated API usage, fallback triggered, retry attempt  |
| **ERROR** | Error events that allow app to continue | Request failed, external API timeout, validation failure |
| **FATAL** | Severe errors causing shutdown          | Database unreachable, critical config missing            |

#### Logging Rules

**1. Every request/operation must log:**

```

// Start of operation
log.Info("creating task",
"correlationId", correlationID,
"userId", userID,
"title", task.Title,
)

// Success
log.Info("task created successfully",
"correlationId", correlationID,
"taskId", task.ID,
"duration", time.Since(start),
)

// Error
log.Error("failed to create task",
"correlationId", correlationID,
"error", err,
"userId", userID,
)

```

**2. Always include context:**

- `correlationId`: Trace requests across services (UUID)
- `userId`: Who triggered the action
- `duration`: Operation timing (milliseconds)
- `error`: Error details (if failed)

**3. Structured logging only** (no string formatting):

```

// ✅ Structured
log.Info("user login", "userId", userID, "ip", clientIP)

// ❌ String formatting
log.Info(fmt.Sprintf("User %s logged in from %s", userID, clientIP))

```

**4. Security - Never log:**

- Passwords or password hashes
- API keys or tokens
- Credit card numbers
- PII in production logs (email/phone only if necessary and sanitized)
- Full request/response bodies (unless DEBUG level in non-prod)

**5. Performance - Never log in hot paths:**

- Inside tight loops
- Per-item processing in batch operations (use summary instead)
- Synchronous logging in latency-critical paths

**Best Practice:** "Use logger middleware redaction rather than manual string manipulation."

#### Language-Specific Implementations

##### Go (zerolog — `backend/internal/pkg/logger`)

```go
// Logger is created once in cmd/server/bootstrap.go via logger.New() and injected
// into modules through module.go — never create a new zerolog.Logger inside a module.
log := logger.New() // APP_ENV=production → JSON; otherwise coloured console with caller

log.Info().
    Str("request_id", reqID).
    Str("org_id", orgID).
    Str("module", "reservation").
    Msg("reservation created")

log.Error().
    Err(err).
    Str("request_id", reqID).
    Str("org_id", orgID).
    Msg("failed to create reservation")
```

Field names are `snake_case`. Never log JWTs, passwords, NIK/identity numbers, or phone numbers (WhatsApp).

##### Frontend (React)

There is no logging library in `frontend/`. `console.error` is allowed only inside `components/ErrorBoundary.tsx` and the error path of `services/api.ts`; components and hooks must not log — surface errors through React Query state / UI instead.

##### Log Patterns by Operation Type

##### API Request/Response

```

// Request received
log.Info("request received",
  "method", r.Method,
  "path", r.URL.Path,
  "correlationId", correlationID,
  "userId", userID,
)

// Request completed
log.Info("request completed",
  "correlationId", correlationID,
  "status", statusCode,
  "duration", duration,
)

```

##### Database Operations

```

// Query start (DEBUG level)
log.Debug("executing query",
  "correlationId", correlationID,
  "query", "SELECT * FROM tasks WHERE user_id = $1",
)

// Query success (DEBUG level)
log.Debug("query completed",
  "correlationId", correlationID,
  "rowsReturned", len(results),
  "duration", duration,
)

// Query error (ERROR level)
log.Error("query failed",
  "correlationId", correlationID,
  "error", err,
  "query", "SELECT * FROM tasks WHERE user_id = $1",
)

```

##### External API Calls

```

// Call start
log.Info("calling external API",
  "correlationId", correlationID,
  "service", "email-provider",
  "endpoint", "/send",
)

// Retry (WARN level)
log.Warn("retrying external API call",
  "correlationId", correlationID,
  "service", "email-provider",
  "attempt", retryCount,
  "error", err,
)

// Circuit breaker open (WARN level)
log.Warn("circuit breaker opened",
  "correlationId", correlationID,
  "service", "email-provider",
  "failureCount", failures,
)

```

##### Error Scenarios

```

// Recoverable error (ERROR level)
log.Error("validation failed",
  "correlationId", correlationID,
  "userId", userID,
  "error", "invalid email format",
  "input", sanitizedInput, // Sanitized!
)

// Fatal error (FATAL level)
log.Fatal("critical dependency unavailable",
  "error", err,
  "dependency", "database",
  "action", "shutting down",
)

```

#### Environment-Specific Configuration

Controlled by `APP_ENV` (see `backend/internal/pkg/logger/logger.go`): `production` emits JSON to stdout; anything else emits human-readable console output with caller info. Do not add per-environment log configuration elsewhere.

#### Testing Logs

**Unit tests:** Capture and assert on log output.

#### Checklist for Every Feature

- [ ] All public operations log INFO on start
- [ ] All operations log INFO/ERROR on complete/failure
- [ ] All logs include correlationId
- [ ] No sensitive data in logs
- [ ] Structured logging (key-value pairs)
- [ ] Appropriate log level used
- [ ] Error logs include error details
- [ ] Performance-critical paths use DEBUG level

### Related Principles

- Logging and Observability Mandate @logging-and-observability-mandate.md
- Error Handling Principles @error-handling-principles.md
- Security Mandate @security-mandate.md
- Security Principles @security-principles.md
- API Design Principles @api-design-principles.md
