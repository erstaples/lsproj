package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"

	"github.com/estaples/lsproj/internal/action"
	"github.com/estaples/lsproj/internal/cache"
	"github.com/estaples/lsproj/internal/config"
	"github.com/estaples/lsproj/internal/history"
	"github.com/estaples/lsproj/internal/scanner"
	"github.com/estaples/lsproj/internal/tui"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `lsproj - fuzzy find and open your git projects

Recursively scans root directories for git repos, then presents a
fuzzy-searchable TUI to select a project and run an action on it.

Usage:
  lsproj [flags]

Examples:
  lsproj                          Launch with config file defaults
  lsproj --root ~/code            Scan a specific root directory
  lsproj --extra ~/code/nke       Include a non-git directory
  lsproj --refresh                Force rescan and reset the cache
  lsproj --cache-ttl 24h          Override cache TTL for this run

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Config:
  Default config path: %s

  roots:      list of directories to recursively scan for git repos
  extras:     list of directories to include regardless of .git presence
  cache_ttl:  how long scan results are cached (e.g. 5m, 1h, 24h)
`, config.DefaultPath())
	}

	configPath := flag.String("config", config.DefaultPath(), "path to config file")
	roots := flag.StringArray("root", nil, "root directory to scan (can be specified multiple times)")
	extras := flag.StringArray("extra", nil, "extra directory to include (can be specified multiple times)")
	cacheTTL := flag.String("cache-ttl", "", "cache TTL duration (e.g. 1h, 30m)")
	refresh := flag.Bool("refresh", false, "force rescan and reset the cache")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// CLI flags override config
	if len(*roots) > 0 {
		cfg.Roots = *roots
	}
	if len(*extras) > 0 {
		cfg.Extras = *extras
	}
	if *cacheTTL != "" {
		cfg.CacheTTLStr = *cacheTTL
	}

	expandedRoots := cfg.ExpandRoots()
	expandedExtras := cfg.ExpandExtras()
	if len(expandedRoots) == 0 && len(expandedExtras) == 0 {
		fmt.Fprintln(os.Stderr, "No root directories configured.")
		fmt.Fprintln(os.Stderr, "Use --root <path> or create a config file at:", config.DefaultPath())
		os.Exit(1)
	}

	// Try cache first unless --refresh
	var projects []scanner.Project
	if !*refresh {
		projects = cache.Load(expandedRoots, cfg.CacheTTL())
	}
	if projects == nil {
		projects = scanner.Scan(expandedRoots, expandedExtras)
		// Save to cache (best effort)
		_ = cache.Save(expandedRoots, projects)
	}

	if len(projects) == 0 {
		fmt.Fprintln(os.Stderr, "No git projects found in configured roots.")
		os.Exit(1)
	}

	counts := history.Load()
	counts.Sort(projects)

	result, err := tui.Run(projects)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if result == nil {
		os.Exit(0)
	}

	// Record selection (best effort)
	_ = counts.Record(result.Project.Path)

	if action.IsCDAction(result.Action) {
		resultFile := os.Getenv("LSPROJ_RESULT_FILE")
		if resultFile == "" {
			// No wrapper — just print the path so it's still useful
			fmt.Println(result.Project.Path)
			return
		}
		if err := os.WriteFile(resultFile, []byte(result.Project.Path), 0o600); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing result file: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if err := result.Action.Run(result.Project.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error running action: %v\n", err)
		os.Exit(1)
	}
}
