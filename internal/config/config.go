package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// configDir returns ~/.nty (e.g. C:\Users\X\.nty on Windows).
func configDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".nty"), nil
}

type Config struct {
	Lang                  string
	CreateBranchForDeploy bool
	ClaudeAccessToken     string
	ClaudeRefreshToken    string
	ClaudeExpiresAt       int64
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if dir, err := configDir(); err == nil {
		viper.AddConfigPath(dir)
	}

	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	config := &Config{
		Lang:                  viper.GetString("lang"),
		CreateBranchForDeploy: viper.GetBool("create_branch_for_deploy"),
		ClaudeAccessToken:     viper.GetString("claude_access_token"),
		ClaudeRefreshToken:    viper.GetString("claude_refresh_token"),
		ClaudeExpiresAt:       viper.GetInt64("claude_expires_at"),
	}

	return config, nil
}

func SaveClaudeAuth(accessToken, refreshToken string, expiresAt int64) error {
	viper.Set("claude_access_token", accessToken)
	viper.Set("claude_refresh_token", refreshToken)
	viper.Set("claude_expires_at", expiresAt)
	return write()
}

func write() error {
	dir, err := configDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	viper.SetConfigFile(filepath.Join(dir, "config.yaml"))
	return viper.WriteConfig()
}

func SaveGithub(user, email string) error {
	viper.Set("github_user", user)
	viper.Set("github_email", email)
	return write()
}
