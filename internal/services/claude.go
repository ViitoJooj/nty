package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ViitoJooj/nty/internal/config"
)

const (
	messagesURL        = "https://api.anthropic.com/v1/messages"
	model              = "claude-haiku-4-5"
	claudeCodeIdentity = "You are Claude Code, Anthropic's official CLI for Claude."
)

func complete(system, user string, maxTokens int) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", err
	}
	if cfg.ClaudeAccessToken == "" {
		return "", fmt.Errorf("não logado no Claude — rode: nty config --ai")
	}
	if cfg.ClaudeExpiresAt != 0 && time.Now().Unix() >= cfg.ClaudeExpiresAt {
		return "", fmt.Errorf("token Claude expirado — rode: nty config --ai")
	}

	body, err := json.Marshal(map[string]any{
		"model":      model,
		"max_tokens": maxTokens,
		"system":     claudeCodeIdentity + "\n\n" + system,
		"messages": []map[string]any{
			{"role": "user", "content": user},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, messagesURL, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+cfg.ClaudeAccessToken)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("anthropic-beta", "oauth-2025-04-20")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API %d: %s", resp.StatusCode, raw)
	}

	var parsed struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, c := range parsed.Content {
		if c.Type == "text" {
			sb.WriteString(c.Text)
		}
	}
	return strings.TrimSpace(sb.String()), nil
}

func CommitMessage(file, diff string) (string, error) {
	system := fmt.Sprintf("You write short, objective git commit messages in %s, using conventional commits style when applicable. Reply with ONLY the message, one line, no quotes and no explanation.", commitLanguage())
	user := fmt.Sprintf("File: %s\n\nDiff:\n%s", file, diff)
	return complete(system, user, 100)
}

// commitLanguage maps the saved lang setting to a language name for the prompt.
// Defaults to English when unset/unknown.
func commitLanguage() string {
	cfg, err := config.Load()
	if err == nil && cfg.Lang == "pt" {
		return "Portuguese"
	}
	return "English"
}

func BranchName(files []string) (string, error) {
	system := "Você sugere um nome de branch git curto em inglês, kebab-case, com prefixo feat/ fix/ ou chore/ quando fizer sentido. Responda APENAS com o nome, uma linha."
	user := "Arquivos alterados:\n" + strings.Join(files, "\n")
	return complete(system, user, 30)
}

// ClassifyArtifact decides whether the project ships a runnable binary or a
// library, from its manifests and file tree. Language-independent: the model
// generalizes instead of us coding one rule per ecosystem.
func ClassifyArtifact(manifests, fileTree string) (string, error) {
	system := "You classify a software project as either a runnable BINARY (CLI/app/service producing an executable) or a LIBRARY (imported by other code). Use the manifests and file tree. Reply with ONLY one word: binary or library."
	user := fmt.Sprintf("Manifests:\n%s\n\nFile tree:\n%s", manifests, fileTree)
	out, err := complete(system, user, 5)
	if err != nil {
		return "", err
	}
	out = strings.ToLower(strings.TrimSpace(out))
	if strings.Contains(out, "binary") {
		return "binary", nil
	}
	return "library", nil
}

// ReleaseNotes picks the next semver version (AI bump from the commits) and
// writes the release notes. Returns the first line (version) and the body.
func ReleaseNotes(lastTag, commits string) (version, notes string, err error) {
	system := fmt.Sprintf("You are a release manager. Given the last git tag and the commits since it, decide the next semantic version (bump major/minor/patch from the changes) and write clean release notes in %s. Respond EXACTLY in this format: first line is only the version like v1.2.3, then a blank line, then the release notes as markdown bullet points. No other text.", commitLanguage())
	user := fmt.Sprintf("Last tag: %s\n\nCommits since:\n%s", lastTag, commits)
	out, err := complete(system, user, 600)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(strings.TrimSpace(out), "\n", 2)
	version = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		notes = strings.TrimSpace(parts[1])
	}
	return version, notes, nil
}
