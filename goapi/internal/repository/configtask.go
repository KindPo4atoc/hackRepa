package repository

type ConfigTask struct {
	DatabaseURL string `toml:"database_url_task"`
}

func NewConfigTask() *ConfigTask {
	return &ConfigTask{}
}
