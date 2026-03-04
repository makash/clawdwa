package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	GroupJID  string   `json:"group_jid"`
	GroupName string   `json:"group_name"`
	BotJID    string   `json:"bot_jid"`
	Prefixes  []string `json:"prefixes"`
	ClaudeBin string   `json:"claude_bin"`
}

func Dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "clawdwa")
}

func Path() string {
	return filepath.Join(Dir(), "config.json")
}

func WADBPath() string {
	return filepath.Join(Dir(), "wa.db")
}

func BotDBPath() string {
	return filepath.Join(Dir(), "bot.db")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(Path())
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(Dir(), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(Path(), data, 0600)
}

func Exists() bool {
	_, err := os.Stat(Path())
	return err == nil
}
