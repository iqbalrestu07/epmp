package errs

import "errors"

// DomainError represents a business rule violation.
type DomainError struct {
	Code    string
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

// NewDomainError creates a typed domain error.
func NewDomainError(code, message string) *DomainError {
	return &DomainError{Code: code, Message: message}
}

// Common domain errors.
var (
	ErrNotFound      = NewDomainError("NOT_FOUND", "resource not found")
	ErrAlreadyExists = NewDomainError("ALREADY_EXISTS", "resource already exists")
	ErrConflict      = NewDomainError("CONFLICT", "resource conflict")
	ErrValidation    = NewDomainError("VALIDATION", "validation failed")
)

// IsDomainError checks if an error is a DomainError with the given code.
func IsDomainError(err error, code string) bool {
	var de *DomainError
	if errors.As(err, &de) {
		return de.Code == code
	}
	return false
}
