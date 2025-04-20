package repository

// конфиг для БД, подтягивается из apiserver.toml
type ConfigUser struct {
	DatabaseURL string `toml:"database_url_user"`
}

func NewConfigUser() *ConfigUser {
	return &ConfigUser{}
}
