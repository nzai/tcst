package config

import (
	"errors"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config global config
type Config struct {
	AppID     string `toml:"app_id"`
	Bucket    string `toml:"bucket"`
	Region    string `toml:"region"`
	SecretID  string `toml:"secret_id"`
	SecretKey string `toml:"secret_key"`
}

// Valid validate config
func (s Config) Valid() error {
	if strings.TrimSpace(s.AppID) == "" {
		return errors.New("app_id undefined")
	}

	if strings.TrimSpace(s.Bucket) == "" {
		return errors.New("bucket undefined")
	}

	if strings.TrimSpace(s.Region) == "" {
		return errors.New("region undefined")
	}

	if strings.TrimSpace(s.SecretID) == "" {
		return errors.New("secret_id undefined")
	}

	if strings.TrimSpace(s.SecretKey) == "" {
		return errors.New("secret_key undefined")
	}

	return nil
}

var (
	currentConfig *Config
)

// Get get current config
func Get() *Config {
	return currentConfig
}

// Parse parse config from file
func Parse(filePath string) (*Config, error) {
	currentConfig = new(Config)
	_, err := toml.DecodeFile(filePath, currentConfig)
	if err != nil {
		return nil, err
	}

	return currentConfig, currentConfig.Valid()
}
