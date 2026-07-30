package config

import (
	"fmt"
	"strings"
	"unicode"
)

// Validate returns all validation errors for a GeneratorConfig.
// Returns nil slice if config is valid.
func Validate(cfg *GeneratorConfig) []error {
	var errs []error

	if cfg.Version == "" {
		errs = append(errs, fmt.Errorf("version is required"))
	}
	if cfg.Frontend.OutputRoot == "" {
		errs = append(errs, fmt.Errorf("frontend.output_root is required"))
	}
	if !isPascalCase(cfg.Domain.Name) {
		errs = append(errs, fmt.Errorf(
			"domain.name %q must be PascalCase (e.g. Property)", cfg.Domain.Name))
	}
	if cfg.Domain.Package == "" {
		errs = append(errs, fmt.Errorf("domain.package is required"))
	}
	if cfg.Domain.Package != strings.ToLower(cfg.Domain.Package) {
		errs = append(errs, fmt.Errorf(
			"domain.package %q must be lowercase", cfg.Domain.Package))
	}
	if cfg.Domain.REST.BasePath != "" && !strings.HasPrefix(cfg.Domain.REST.BasePath, "/api/") {
		errs = append(errs, fmt.Errorf(
			"domain.rest.base_path must start with /api/ (got %q)", cfg.Domain.REST.BasePath))
	}
	for _, f := range cfg.Domain.Fields {
		if f.Type == FieldTypeEnum && len(f.EnumValues) < 2 {
			errs = append(errs, fmt.Errorf(
				"field %q of type enum must have at least 2 enum_values", f.Name))
		}
	}
	return errs
}

func isPascalCase(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	return unicode.IsUpper(runes[0])
}
