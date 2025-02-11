package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logs                   string `toml:"logs"`
	SqlServerDatabase      SqlServerDBConfig
	PostgresServerDatabase PostgresSqlDBConfig
}

func ReadConfigFile(path string) (Config, error) {
	fileConfig, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	cfg, err := io.ReadAll(fileConfig)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	_, err = toml.Decode(string(cfg), &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
