package repository

import (
	"fmt"
	"goapi/internal/entity"

	"golang.org/x/crypto/bcrypt"
)

// структура для взаимодействия с бд
type DBRepository struct {
	store *DataBase
}

func (r *DBRepository) ValidateUser(login, pass string) (entity.Answer, error) {
	var dataRow entity.UserData
	fmt.Println(pass)
	rows, err := r.store.db.Query(
		"SELECT passhash "+
			"FROM users where login = $1;",
		login,
	)

	if err != nil {
		return entity.Answer{}, err
	}
	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(
			&dataRow.PasswordHash,
		)
		if err != nil {
			return entity.Answer{}, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(dataRow.PasswordHash), []byte(pass))
	if err != nil {
		return entity.Answer{Status: err.Error()}, err
	} else {
		return entity.Answer{Status: "200 OK"}, nil
	}
}

func (r *DBRepository) ExistUser(login string) (entity.Answer, error) {
	var users entity.ContextData
	rows, err := r.store.db.Query(
		"SELECT login "+
			"FROM users where login = $1;",
		login,
	)

	if err != nil {
		return entity.Answer{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var dataRow entity.UserData
		err := rows.Scan(
			&dataRow.Login,
		)
		if err != nil {
			return entity.Answer{}, err
		}
		users.Data = append(users.Data, dataRow)
	}
	if len(users.Data) > 0 {
		return entity.Answer{Status: "200 OK"}, nil
	} else {
		return entity.Answer{Status: "404 Not found"}, nil
	}
}

func (r *DBRepository) AddUsers(login, passHash string) (entity.Answer, error) {
	answer, err := r.ExistUser(login)
	if err != nil || answer.Status == "200 OK" {
		return entity.Answer{Status: "User exist"}, err
	}

	_, err = r.store.db.Exec("Insert Into users(login, passhash) values($1, $2)", login, passHash)

	if err != nil {
		return entity.Answer{Status: err.Error()}, err
	}

	return entity.Answer{Status: "200 OK"}, nil
}

func (r *DBRepository) DestroyDBTask(dbName string) (entity.Answer, error) {
	cmd := fmt.Sprintf("drop database %s;", dbName)
	_, err := r.store.db.Exec(cmd)

	if err != nil {
		return entity.Answer{Status: err.Error()}, err
	} else {
		return entity.Answer{Status: "200 OK"}, nil
	}

}
