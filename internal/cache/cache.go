package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/estaples/lsproj/internal/scanner"
)

type entry struct {
	Projects  []scanner.Project `json:"projects"`
	Roots     []string          `json:"roots"`
	UpdatedAt time.Time         `json:"updated_at"`
}

func Path() string {
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "lsproj", "cache.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", "lsproj", "cache.json")
}

// Load returns cached projects if the cache exists, is not expired, and
// was built from the same set of roots. Returns nil if the cache is stale.
func Load(roots []string, ttl time.Duration) []scanner.Project {
	data, err := os.ReadFile(Path())
	if err != nil {
		return nil
	}
	var e entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil
	}
	if time.Since(e.UpdatedAt) > ttl {
		return nil
	}
	if !slicesEqual(e.Roots, roots) {
		return nil
	}
	return e.Projects
}

// Save writes the project list to the cache file.
func Save(roots []string, projects []scanner.Project) error {
	e := entry{
		Projects:  projects,
		Roots:     roots,
		UpdatedAt: time.Now(),
	}
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	dir := filepath.Dir(Path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0o644)
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
