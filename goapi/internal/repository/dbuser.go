package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataBaseUser struct {
	config         *ConfigUser
	db             *sql.DB
	dataRepository *DBUserRepository
}

func NewUser(c *ConfigUser) *DataBaseUser {
	return &DataBaseUser{
		config: c,
	}
}

func (data *DataBaseUser) OpenNew(connect string) error {
	db, err := sql.Open("postgres", connect)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	data.db = db

	return nil
}

func (data *DataBaseUser) Open() error {
	db, err := sql.Open("postgres", data.config.DatabaseURL)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	data.db = db

	return nil
}

func (db *DataBaseUser) Close() {
	db.db.Close()
}

func (data *DataBaseUser) Data() *DBUserRepository {
	if data.dataRepository != nil {
		return data.dataRepository
	}

	data.dataRepository = &DBUserRepository{
		store: data,
	}
	return data.dataRepository
}
