package git

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func DetectRepo() (string, string, error) {
	out, err := exec.Command("git", "remote", "get-url", "origin").Output()
	if err != nil {
		return "", "", err
	}

	url := strings.TrimSpace(string(out))
	url = strings.Replace(url, "https://github.com/", "", 1)
	url = strings.Replace(url, "git@github.com:", "", 1)
	url = strings.Replace(url, ".git", "", 1)

	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid remote URL format")
	}

	return parts[0], parts[1], nil
}

func GetCommitActivity() ([]float64, []string, error) {
	activity := make([]float64, 30)
	labels := make([]string, 30)

	now := time.Now()
	for i := 0; i < 30; i++ {
		date := now.AddDate(0, 0, -i)
		labels[29-i] = date.Format("02 Jan")

		since := date.Format("2006-01-02 00:00:00")
		until := date.Format("2006-01-02 23:59:59")

		out, err := exec.Command("git", "rev-list", "--count", "--since", since, "--until", until, "HEAD").Output()
		if err != nil {
			activity[29-i] = 0
			continue
		}

		countStr := strings.TrimSpace(string(out))
		var count float64
		fmt.Sscanf(countStr, "%f", &count)
		activity[29-i] = count
	}

	return activity, labels, nil
}

type DirSize struct {
	Name string
	Size int64
}

func GetRepoSize() (int64, []DirSize, error) {
	var totalSize int64
	dirSizes := make(map[string]int64)

	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path == ".git" || path == "node_modules" || path == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		size := info.Size()
		totalSize += size

		parts := strings.Split(path, string(os.PathSeparator))
		if len(parts) > 1 {
			dirSizes[parts[0]] += size
		}

		return nil
	})

	if err != nil {
		return 0, nil, err
	}

	var sortedSizes []DirSize
	for name, size := range dirSizes {
		sortedSizes = append(sortedSizes, DirSize{Name: name, Size: size})
	}

	sort.Slice(sortedSizes, func(i, j int) bool {
		return sortedSizes[i].Size > sortedSizes[j].Size
	})

	limit := 5
	if len(sortedSizes) < limit {
		limit = len(sortedSizes)
	}

	return totalSize, sortedSizes[:limit], nil
}

func FormatSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

type CodeStats struct {
	Language string
	Files    int
	Lines    int
}

func GetCodeStats() (map[string]*CodeStats, int, error) {
	out, err := exec.Command("git", "ls-files").Output()
	if err != nil {
		return nil, 0, err
	}

	stats := make(map[string]*CodeStats)
	totalLines := 0

	files := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, file := range files {
		ext := filepath.Ext(file)
		lang := ""
		switch ext {
		case ".go":
			lang = "Go"
		case ".js", ".jsx":
			lang = "JavaScript"
		case ".ts", ".tsx":
			lang = "TypeScript"
		case ".md":
			lang = "Markdown"
		case ".html":
			lang = "HTML"
		case ".css":
			lang = "CSS"
		case ".json":
			lang = "JSON"
		case ".yml", ".yaml":
			lang = "YAML"
		default:
			lang = "Other"
		}

		if _, ok := stats[lang]; !ok {
			stats[lang] = &CodeStats{Language: lang}
		}
		stats[lang].Files++

		lines, err := countLines(file)
		if err == nil {
			stats[lang].Lines += lines
			totalLines += lines
		}
	}

	return stats, totalLines, nil
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines, scanner.Err()
}

func GetStaleBranches() ([]string, error) {
	out, err := exec.Command("git", "for-each-ref", "--sort=-committerdate", "--format=%(refname:short) (%(committerdate:relative))", "refs/heads/").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var stale []string
	for _, line := range lines {
		if strings.Contains(line, "weeks") || strings.Contains(line, "month") || strings.Contains(line, "year") {
			stale = append(stale, line)
		}
	}
	return stale, nil
}

func GetRepoHealth() (string, int, error) {
	stale, _ := GetStaleBranches()
	score := 100 - (len(stale) * 5)
	if score < 0 {
		score = 0
	}

	status := "Healthy"
	if score < 80 {
		status = "Attention Needed"
	} else if score < 50 {
		status = "Critical"
	}

	return status, score, nil
}