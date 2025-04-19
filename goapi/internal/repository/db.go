package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBase struct {
	config         *Config
	db             *sql.DB
	dataRepository *DBRepository
}

func New(config *Config) *DataBase {
	return &DataBase{
		config: config,
	}
}

func (data *DataBase) Open() error {
	db, err := sql.Open("postgres", data.config.DatabaseURL)
	fmt.Println(data.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}
	fmt.Println("ok")
	data.db = db

	return nil
}

func (db *DataBase) Close() {
	db.db.Close()
}

func (data *DataBase) Data() *DBRepository {
	if data.dataRepository != nil {
		return data.dataRepository
	}

	data.dataRepository = &DBRepository{
		store: data,
	}
	return data.dataRepository
}
