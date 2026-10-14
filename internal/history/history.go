package history

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/estaples/lsproj/internal/scanner"
)

// Counts maps project path -> access count.
type Counts map[string]int

func Path() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "lsproj", "history.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "lsproj", "history.json")
}

func Load() Counts {
	data, err := os.ReadFile(Path())
	if err != nil {
		return make(Counts)
	}
	var c Counts
	if err := json.Unmarshal(data, &c); err != nil {
		return make(Counts)
	}
	return c
}

func (c Counts) Record(projectPath string) error {
	c[projectPath]++
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	dir := filepath.Dir(Path())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0o644)
}

// Sort reorders projects by descending access count, with alphabetical
// as a tiebreaker for projects with the same count.
func (c Counts) Sort(projects []scanner.Project) {
	sort.SliceStable(projects, func(i, j int) bool {
		ci, cj := c[projects[i].Path], c[projects[j].Path]
		if ci != cj {
			return ci > cj
		}
		return projects[i].Name < projects[j].Name
	})
}
