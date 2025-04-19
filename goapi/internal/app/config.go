package apiserver

import "goapi/internal/repository"

// конфиг для подтягивания данных из apiserver.toml
type Config struct {
	BindAddr     string                 `toml:"bind_addr"`
	LogLevel     string                 `toml:"log_level"`
	DBConfig     *repository.Config     `toml:"database"`
	DBTaskConfig *repository.ConfigTask `toml:"database_task"`
}

func NewConfig() *Config {
	return &Config{
		DBConfig: repository.NewConfig(),
	}
}
