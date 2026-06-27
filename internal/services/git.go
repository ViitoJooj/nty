package services

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out, errBuf strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %v: %s", strings.Join(args, " "), err, strings.TrimSpace(errBuf.String()))
	}
	return out.String(), nil
}

func ChangedFiles() ([]string, error) {
	out, err := run("status", "--porcelain")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		path := line[3:]
		if i := strings.Index(path, " -> "); i >= 0 {
			path = path[i+4:]
		}
		files = append(files, strings.Trim(strings.TrimSpace(path), "\""))
	}
	return files, nil
}

func CreateBranch(name string) error {
	_, err := run("checkout", "-b", name)
	return err
}

func CurrentBranch() (string, error) {
	out, err := run("rev-parse", "--abbrev-ref", "HEAD")
	return strings.TrimSpace(out), err
}

func IsDefaultBranch(name string) bool {
	switch name {
	case "main", "master", "develop", "HEAD":
		return true
	}
	return false
}

func CleanMergedBranches() ([]string, error) {
	if _, err := run("fetch", "--prune"); err != nil {
		return nil, err
	}

	out, err := run("for-each-ref", "--format=%(refname:short) %(upstream:track)", "refs/heads")
	if err != nil {
		return nil, err
	}

	var deleted []string
	for _, line := range strings.Split(out, "\n") {
		name, track, _ := strings.Cut(strings.TrimSpace(line), " ")
		if name == "" || IsDefaultBranch(name) {
			continue
		}
		if !strings.Contains(track, "[gone]") {
			continue
		}
		if _, err := run("branch", "-d", name); err != nil {
			continue
		}
		deleted = append(deleted, name)
	}
	return deleted, nil
}

func StageAndDiff(file string) (string, error) {
	if _, err := run("add", "--", file); err != nil {
		return "", err
	}
	return run("diff", "--cached", "--", file)
}

func Commit(message string) error {
	_, err := run("commit", "-m", message)
	return err
}

func Push(branch string) error {
	_, err := run("push", "-u", "origin", branch)
	return err
}

func LastTag() string {
	out, err := run("describe", "--tags", "--abbrev=0")
	if err != nil {
		return "v0.0.0"
	}
	return strings.TrimSpace(out)
}

func CommitsSince(tag string) string {
	out, err := run("log", tag+"..HEAD", "--pretty=format:%s")
	if err != nil {
		out, _ = run("log", "--pretty=format:%s")
	}
	return strings.TrimSpace(out)
}

func PushTags() error {
	_, err := run("push", "--tags")
	return err
}

func SanitizeBranch(name string) string {
	name = strings.TrimSpace(name)
	if i := strings.IndexByte(name, '\n'); i >= 0 {
		name = name[:i]
	}
	name = strings.ToLower(name)

	var b strings.Builder
	prevDash := false
	for _, r := range name {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'),
			r == '/' || r == '-' || r == '_':
			b.WriteRune(r)
			prevDash = false

		case r == ' ':
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-/")
	if out == "" {
		out = "nty-changes"
	}
	return out
}
