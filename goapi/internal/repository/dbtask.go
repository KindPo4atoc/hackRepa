package repository

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type DataBaseTask struct {
	config         *ConfigTask
	db             *sql.DB
	dataRepository *DBTaskRepository
}

func NewTask(c *ConfigTask) *DataBaseTask {
	return &DataBaseTask{
		config: c,
	}
}

func (data *DataBaseTask) OpenNew(connect string) error {
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

func (data *DataBaseTask) Open() error {
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

func (db *DataBaseTask) Close() {
	db.db.Close()
}

func (data *DataBaseTask) Data() *DBTaskRepository {
	if data.dataRepository != nil {
		return data.dataRepository
	}

	data.dataRepository = &DBTaskRepository{
		store: data,
	}
	return data.dataRepository
}
