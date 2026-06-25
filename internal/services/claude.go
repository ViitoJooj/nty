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
	system := "Você escreve mensagens de commit curtas e objetivas em português, estilo conventional commits quando aplicável. Responda APENAS com a mensagem, uma linha, sem aspas e sem explicação."
	user := fmt.Sprintf("Arquivo: %s\n\nDiff:\n%s", file, diff)
	return complete(system, user, 100)
}

func BranchName(files []string) (string, error) {
	system := "Você sugere um nome de branch git curto em inglês, kebab-case, com prefixo feat/ fix/ ou chore/ quando fizer sentido. Responda APENAS com o nome, uma linha."
	user := "Arquivos alterados:\n" + strings.Join(files, "\n")
	return complete(system, user, 30)
}
