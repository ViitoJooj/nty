package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var semverRe = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// manifestFiles are the files we feed to the AI to detect language + bin/lib.
var manifestFiles = []string{
	"go.mod", "Cargo.toml", "package.json", "pyproject.toml",
	"setup.py", "pom.xml", "build.gradle", "CMakeLists.txt",
}

// RunRelease builds a release: AI-decided version + notes, optional artifact
// build, and a GitHub release. Dry-run unless publish is true.
func RunRelease(publish bool) error {
	lastTag := LastTag()
	commits := CommitsSince(lastTag)
	if commits == "" {
		fmt.Printf("nada novo desde %s.\n", lastTag)
		return nil
	}

	manifests, tree := projectContext()
	kind, err := ClassifyArtifact(manifests, tree)
	if err != nil {
		return fmt.Errorf("classificar projeto: %w", err)
	}

	version, notes, err := ReleaseNotes(lastTag, commits)
	if err != nil {
		return fmt.Errorf("gerar release notes: %w", err)
	}
	if !semverRe.MatchString(version) {
		return fmt.Errorf("versão inválida sugerida pela IA: %q", version)
	}

	artifact := ""
	if kind == "binary" {
		artifact, err = buildArtifact(version)
		if err != nil {
			fmt.Println("aviso: falha ao buildar artefato, seguindo sem ele:", err)
		}
	}

	// Plano.
	fmt.Printf("última tag : %s\n", lastTag)
	fmt.Printf("nova versão: %s\n", version)
	fmt.Printf("tipo       : %s\n", kind)
	if artifact != "" {
		fmt.Printf("artefato   : %s\n", artifact)
	}
	fmt.Printf("\n%s\n\n", notes)

	if !publish {
		fmt.Println("dry-run. rode com --publish para criar a release no GitHub.")
		return nil
	}

	return ghRelease(version, notes, artifact)
}

// projectContext reads known manifests and a trimmed file tree.
func projectContext() (manifests, tree string) {
	var mb strings.Builder
	for _, f := range manifestFiles {
		if data, err := os.ReadFile(f); err == nil {
			fmt.Fprintf(&mb, "=== %s ===\n%s\n", f, data)
		}
	}
	out, _ := run("ls-files")
	lines := strings.Split(out, "\n")
	if len(lines) > 100 {
		lines = lines[:100] // ponytail: first 100 paths is enough signal for classification
	}
	return mb.String(), strings.Join(lines, "\n")
}

// buildArtifact compiles a release binary for known languages. Uses a fixed
// command table, never AI-generated shell.
// ponytail: Go only for now; add cargo/npm rows when a non-Go project needs it.
func buildArtifact(version string) (string, error) {
	if _, err := os.Stat("go.mod"); err != nil {
		return "", fmt.Errorf("build automático só suportado para Go por enquanto")
	}

	name := filepath.Base(mustWd())
	out := filepath.Join("dist", name+"-"+version+binExt())
	if err := os.MkdirAll("dist", 0o755); err != nil {
		return "", err
	}

	cmd := exec.Command("go", "build", "-o", out, ".")
	if combined, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v: %s", err, strings.TrimSpace(string(combined)))
	}
	return out, nil
}

// ghRelease creates the GitHub release (and its tag) via the gh CLI.
func ghRelease(version, notes, artifact string) error {
	args := []string{"release", "create", version, "--title", version, "--notes", notes}
	if artifact != "" {
		args = append(args, artifact)
	}
	cmd := exec.Command("gh", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gh release: %v: %s", err, strings.TrimSpace(string(out)))
	}
	fmt.Printf("release %s criada.\n%s", version, out)
	return nil
}

func mustWd() string {
	wd, _ := os.Getwd()
	return wd
}

func binExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}
