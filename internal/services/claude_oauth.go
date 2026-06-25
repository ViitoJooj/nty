// Package oauth implements the Claude (Anthropic) subscription OAuth flow,
// the same PKCE flow Claude Code uses. No third-party app registration exists,
// so the public Claude Code client_id is reused for personal login.
package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	clientID     = "9d1c250a-e61b-44d9-88ed-5944d1962f5e"
	authorizeURL = "https://claude.ai/oauth/authorize"
	tokenURL     = "https://console.anthropic.com/v1/oauth/token"
	redirectURI  = "https://console.anthropic.com/oauth/code/callback"
	scopes       = "org:create_api_key user:profile user:inference"
)

// PKCE holds the verifier/state needed to start and finish a login.
type PKCE struct {
	Verifier string
	State    string
}

// Tokens is the result of a successful exchange.
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

// ExpiresAt is the unix second the access token stops being valid.
func (t Tokens) ExpiresAt() int64 {
	return time.Now().Unix() + int64(t.ExpiresIn)
}

func randURLSafe(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Start builds the authorize URL the user opens in a browser, and returns the
// PKCE state needed to later exchange the pasted code.
func Start() (loginURL string, p PKCE, err error) {
	verifier, err := randURLSafe(32)
	if err != nil {
		return "", PKCE{}, err
	}
	state, err := randURLSafe(32)
	if err != nil {
		return "", PKCE{}, err
	}

	sum := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(sum[:])

	q := url.Values{
		"code":                  {"true"},
		"client_id":             {clientID},
		"response_type":         {"code"},
		"redirect_uri":          {redirectURI},
		"scope":                 {scopes},
		"code_challenge":        {challenge},
		"code_challenge_method": {"S256"},
		"state":                 {state},
	}

	return authorizeURL + "?" + q.Encode(), PKCE{Verifier: verifier, State: state}, nil
}

// Exchange trades the pasted code (format "code#state") for tokens.
func Exchange(pasted string, p PKCE) (Tokens, error) {
	code, state, _ := strings.Cut(strings.TrimSpace(pasted), "#")
	if state == "" {
		// Some callbacks omit the #state suffix; fall back to our own.
		state = p.State
	}

	body, err := json.Marshal(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     clientID,
		"code":          code,
		"state":         state,
		"redirect_uri":  redirectURI,
		"code_verifier": p.Verifier,
	})
	if err != nil {
		return Tokens{}, err
	}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(string(body)))
	if err != nil {
		return Tokens{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Tokens{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(resp.Body)
		return Tokens{}, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, msg)
	}

	var tok Tokens
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return Tokens{}, err
	}
	if tok.AccessToken == "" {
		return Tokens{}, fmt.Errorf("token exchange returned no access_token")
	}
	return tok, nil
}
