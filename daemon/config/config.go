package config

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pelletier/go-toml"
)

type AccConfig struct {
	Name       string
	Cookie     string
	LastSignIn time.Time
}
type Config struct {
	Account []AccConfig
}

func (c Config) WriteConfig(writer io.Writer) error {
	encoder := toml.NewEncoder(writer)
	return encoder.Encode(c)
}
func (c Config) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	return c.WriteConfig(file)
}

func FromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error when reading config file: %w", err)
	}
	return ReadConfig(file)
}

func ReadConfig(reader io.Reader) (*Config, error) {
	var config Config
	decoder := toml.NewDecoder(reader)
	err := decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error decoding toml config: %w", err)
	}
	for i := range config.Account {
		config.Account[i].Cookie = strings.ReplaceAll(
			config.Account[i].Cookie, "\n", "")
	}
	return &config, nil
}
