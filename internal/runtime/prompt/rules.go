package prompt

import (
	"context"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func LoadRules(ctx context.Context, workingDir string, resolver IncludeResolver) (string, error) {
	ruleDirs := []string{
		filepath.Join(workingDir, ".lattice", "rules"),
		filepath.Join(workingDir, ".claude", "rules"),
	}
	var sections []string
	for _, dir := range ruleDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		names := make([]string, 0, len(files))
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".md") {
				names = append(names, file.Name())
			}
		}
		sort.Strings(names)
		for _, name := range names {
			path := filepath.Join(dir, name)
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			content := string(data)
			if resolver != nil {
				content, err = resolver.Resolve(ctx, filepath.Dir(path), content)
				if err != nil {
					return "", err
				}
			}
			sections = append(sections, "## "+path+"\n"+strings.TrimSpace(content))
		}
	}
	return strings.TrimSpace(strings.Join(sections, "\n\n")), nil
}
