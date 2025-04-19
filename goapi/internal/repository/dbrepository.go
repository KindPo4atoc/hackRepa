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
func (r *DBRepository) GetAllTasks() (entity.ContextTask, error) {
	var tasks entity.ContextTask
	rows, err := r.store.db.Query(
		"select id_task, heading, description, s.login, task_level " +
			"from tasks as t left join users as s on t.fk_user = s.id where t.status = 1;",
	)
	if err != nil {
		return entity.ContextTask{}, err
	}

	defer rows.Close()

	for rows.Next() {
		var dataRow entity.Task
		err := rows.Scan(
			&dataRow.IdTask,
			&dataRow.Header,
			&dataRow.DescriptionTask,
			&dataRow.LoginAuthor,
			&dataRow.LevelTask,
		)
		if err != nil {
			return entity.ContextTask{}, err
		}
		tasks.Data = append(tasks.Data, dataRow)
	}
	if len(tasks.Data) > 0 {
		return tasks, nil
	} else {
		return entity.ContextTask{}, nil
	}
}
func (r *DBRepository) GetTask(idTask string) (entity.Task, error) {
	var task entity.Task
	rows, err := r.store.db.Query(
		"select id_task, heading, description, s.login, task_level "+
			"from tasks as t left join users as s on t.fk_user = s.id where t.status = 1 and t.id_task = $1;",
		idTask,
	)
	if err != nil {
		return entity.Task{}, err
	}

	defer rows.Close()
	i := 0
	for rows.Next() {
		err := rows.Scan(
			&task.IdTask,
			&task.Header,
			&task.DescriptionTask,
			&task.LoginAuthor,
			&task.LevelTask,
		)
		if err != nil {
			return entity.Task{}, err
		}
		i++
	}
	if i > 0 {
		return task, nil
	} else {
		return entity.Task{}, nil
	}
}
func (r *DBRepository) GetTasksByLevel(level string) (entity.ContextTask, error) {
	var tasks entity.ContextTask
	rows, err := r.store.db.Query(
		"select id_task, heading, description, s.login, task_level "+
			"from tasks as t left join users as s on t.fk_user = s.id where t.status = 1 and t.task_level = $1;",
		level,
	)
	if err != nil {
		return entity.ContextTask{}, err
	}

	defer rows.Close()
	for rows.Next() {
		var dataRow entity.Task
		err := rows.Scan(
			&dataRow.IdTask,
			&dataRow.Header,
			&dataRow.DescriptionTask,
			&dataRow.LoginAuthor,
			&dataRow.LevelTask,
		)
		if err != nil {
			return entity.ContextTask{}, err
		}
		tasks.Data = append(tasks.Data, dataRow)
	}
	if len(tasks.Data) > 0 {
		return tasks, nil
	} else {
		return entity.ContextTask{}, nil
	}
}
