package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/epmp/sdk/fe-codegen/config"
	"github.com/epmp/sdk/fe-codegen/internal/filesystem"
	"github.com/epmp/sdk/fe-codegen/internal/generator"
)

func main() {
	configPath := flag.String("config", "", "path to the module config YAML file")
	outputRoot := flag.String("output", "", "output root directory (overrides config)")
	baseURL := flag.String("base-url", "", "API base URL (overrides config)")
	dryRun := flag.Bool("dry-run", false, "render templates without writing files")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "error: --config is required")
		flag.Usage()
		os.Exit(1)
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err, "path", *configPath)
		os.Exit(1)
	}

	if errs := config.Validate(cfg); len(errs) > 0 {
		slog.Error("config validation failed", "path", *configPath)
		for _, e := range errs {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		os.Exit(1)
	}

	out := *outputRoot
	if out == "" {
		out = cfg.Frontend.OutputRoot
	}

	if *dryRun {
		cfg.DryRun = true
	}

	absOut, err := filepath.Abs(out)
	if err != nil {
		slog.Error("cannot resolve output path", "error", err)
		os.Exit(1)
	}

	fs := filesystem.New()
	gen := generator.New(generator.DefaultRenderer(), fs)

	slog.Info("generating frontend module",
		"domain", cfg.Domain.Name,
		"package", cfg.Domain.Package,
		"output", absOut,
		"dry_run", cfg.DryRun,
	)

	req := &generator.GenerateRequest{
		Config:     cfg,
		OutputRoot: absOut,
		BaseURL:    *baseURL,
	}

	if cfg.DryRun {
		if err := gen.DryRun(req); err != nil {
			slog.Error("dry run failed", "error", err)
			os.Exit(1)
		}
		slog.Info("dry run completed — no files written")
		return
	}

	if err := gen.Generate(req); err != nil {
		slog.Error("generation failed", "error", err)
		os.Exit(1)
	}

	slog.Info("frontend module generated successfully",
		"domain", cfg.Domain.Name,
		"output", absOut,
	)
}
