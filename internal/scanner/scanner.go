package scanner

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Project struct {
	Name string
	Path string
}

func Scan(roots []string, extras []string) []Project {
	var projects []Project
	seen := make(map[string]bool)

	// Add explicitly listed extra directories first
	for _, extra := range extras {
		info, err := os.Stat(extra)
		if err != nil || !info.IsDir() {
			continue
		}
		if !seen[extra] {
			seen[extra] = true
			projects = append(projects, Project{
				Name: filepath.Base(extra),
				Path: extra,
			})
		}
	}

	for _, root := range roots {
		filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return filepath.SkipDir
			}
			if !d.IsDir() {
				return nil
			}
			// Skip hidden directories (except .git which we're looking for)
			name := d.Name()
			if name != "." && strings.HasPrefix(name, ".") && name != ".git" {
				return filepath.SkipDir
			}
			if name == ".git" {
				projectPath := filepath.Dir(path)
				if !seen[projectPath] {
					seen[projectPath] = true
					displayName := relativeDisplay(roots, projectPath)
					projects = append(projects, Project{
						Name: displayName,
						Path: projectPath,
					})
				}
				return filepath.SkipDir
			}
			// Skip common non-project directories
			switch name {
			case "node_modules", "vendor", "__pycache__", ".cache":
				return filepath.SkipDir
			}
			return nil
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
	return projects
}

func relativeDisplay(roots []string, projectPath string) string {
	for _, root := range roots {
		if rel, err := filepath.Rel(root, projectPath); err == nil && !strings.HasPrefix(rel, "..") {
			return rel
		}
	}
	return filepath.Base(projectPath)
}
