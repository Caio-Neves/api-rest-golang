package config

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Port                   int    `toml:"Port"`
	Logs                   string `toml:"logs"`
	BaseUrl                string `toml:"baseUrl"`
	SqlServerDatabase      SqlServerDBConfig
	PostgresServerDatabase PostgresSqlDBConfig
}

func ReadConfigFile(path string) (*Config, error) {
	fileConfig, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	cfg, err := io.ReadAll(fileConfig)
	if err != nil {
		return nil, err
	}
	config := Config{}
	_, err = toml.Decode(string(cfg), &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
