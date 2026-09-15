package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

// maxAuditBodySize caps stored request bodies to keep the audit table lean.
const maxAuditBodySize = 8 * 1024

// sensitiveKeys are stripped from stored request bodies.
var sensitiveKeys = []string{"password", "token", "secret", "authorization", "credential", "api_key", "apikey"}

// AuditLog returns a middleware that records every mutating request
// (POST/PUT/PATCH/DELETE) on authenticated routes into the audit_logs table.
// Audit failures are logged but never fail the request.
func AuditLog(db *pgxpool.Pool, log zerolog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			method := c.Request().Method
			if method != "POST" && method != "PUT" && method != "PATCH" && method != "DELETE" {
				return next(c)
			}

			var body []byte
			if c.Request().Body != nil {
				body, _ = io.ReadAll(io.LimitReader(c.Request().Body, maxAuditBodySize+1))
				c.Request().Body = io.NopCloser(bytes.NewBuffer(body))
			}

			err := next(c)

			module, entityID := auditModuleAndEntity(c.Request().URL.Path)
			record := map[string]interface{}{
				"orgID":     GetOrgID(c),
				"userID":    GetUserID(c),
				"userEmail": GetUserEmail(c),
				"action":    auditAction(method),
				"module":    module,
				"entityID":  entityID,
				"method":    method,
				"path":      c.Request().URL.Path,
				"status":    c.Response().Status,
				"ip":        c.RealIP(),
				"ua":        c.Request().UserAgent(),
				"body":      redactAuditBody(body),
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if _, qerr := db.Exec(ctx, `
				INSERT INTO audit_logs
					(organization_id, user_id, user_email, action, module, entity_id, method, path, status_code, ip_address, user_agent, request_body)
				VALUES
					(NULLIF($1,'')::uuid, NULLIF($2,'')::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
				record["orgID"], record["userID"], record["userEmail"], record["action"], record["module"],
				record["entityID"], record["method"], record["path"], record["status"],
				record["ip"], record["ua"], record["body"],
			); qerr != nil {
				log.Error().Err(qerr).Str("path", c.Request().URL.Path).Msg("audit log insert failed")
			}

			return err
		}
	}
}

func auditAction(method string) string {
	switch method {
	case "POST":
		return "CREATE"
	case "PUT", "PATCH":
		return "UPDATE"
	case "DELETE":
		return "DELETE"
	}
	return method
}

// auditModuleAndEntity extracts module and entity id from /api/v1/<module>/<id>/...
func auditModuleAndEntity(path string) (module, entityID string) {
	p := strings.TrimPrefix(path, "/api/v1/")
	segs := strings.Split(strings.Trim(p, "/"), "/")
	if len(segs) > 0 {
		module = segs[0]
	}
	if len(segs) > 1 {
		entityID = segs[1]
	}
	return module, entityID
}

// redactAuditBody parses the request body as JSON and removes sensitive fields.
// Non-JSON bodies are stored as a truncated string.
func redactAuditBody(body []byte) interface{} {
	if len(body) == 0 {
		return nil
	}
	if len(body) > maxAuditBodySize {
		body = body[:maxAuditBodySize]
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return map[string]interface{}{"_raw": string(body)}
	}
	for k := range parsed {
		lk := strings.ToLower(k)
		for _, sk := range sensitiveKeys {
			if strings.Contains(lk, sk) {
				parsed[k] = "[REDACTED]"
				break
			}
		}
	}
	return parsed
}
