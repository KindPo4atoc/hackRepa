package repository

import (
	"fmt"
	"goapi/internal/entity"
	"os"
	"strconv"
	"strings"

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
		"SELECT passhash,role "+
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
			&dataRow.Role,
		)
		if err != nil {
			return entity.Answer{}, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(dataRow.PasswordHash), []byte(pass))
	if err != nil {
		return entity.Answer{Status: err.Error()}, err
	} else {
		return entity.Answer{Status: strconv.Itoa(dataRow.Role)}, nil
	}
}

func (r *DBRepository) ExistUser(login, email string) (entity.Answer, error) {
	var users entity.ContextData
	rows, err := r.store.db.Query(
		"SELECT login "+
			"FROM users where login = $1 or email = $2;",
		login,
		email,
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

func (r *DBRepository) AddUsers(login, email, passHash string) (entity.Answer, error) {
	answer, err := r.ExistUser(login, email)
	if err != nil || answer.Status == "200 OK" {
		return entity.Answer{Status: "User exist"}, err
	}

	_, err = r.store.db.Exec("Insert Into users(login, passhash,email) values($1, $2, $3)", login, passHash, email)

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
func (r *DBRepository) GetTasksStatus() (entity.ContextTask, error) {
	var tasks entity.ContextTask
	rows, err := r.store.db.Query(
		"select id_task, heading, description, s.login " +
			"from tasks as t left join users as s on t.fk_user = s.id where t.status = 0;",
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

func (r *DBRepository) GetPathCsv(taskNumber int) (string, error) {
	var tablePath string
	pathToTables, err := r.store.db.Query("select  path_to_bd from tasks where id_task = $1", taskNumber)

	if err != nil {
		return "", err
	}
	defer pathToTables.Close()

	for pathToTables.Next() {
		var path entity.Task
		err := pathToTables.Scan(&path.PathTotables)
		if err != nil {
			return "", err
		}
		tablePath = fmt.Sprintf("%v", path.PathTotables)
	}
	fmt.Println(tablePath)
	return tablePath, nil
}
func (r *DBRepository) AddTask(task entity.AddTask) entity.Answer {
	//var idUser int64
	var idNewEntry int64
	/*id, err := r.store.db.Query("select id_user from users where login = $1", task.Author)
	if err != nil {
		return entity.Answer{Status: err.Error()}
	}
	defer id.Close()

	for id.Next() {
		err := id.Scan(&idUser)
		if err != nil {
			return entity.Answer{Status: err.Error()}
		}
	}*/
	cmd := fmt.Sprintf("insert into tasks(heading, description, fk_user, path_to_bd) values('%s','%s',%d, '') returning id_task;", task.Header, task.Description, 1) //idUser)
	fmt.Println(cmd)
	err := r.store.db.QueryRow(cmd).Scan(&idNewEntry)
	if err != nil {
		fmt.Println(err.Error())
		return entity.Answer{Status: err.Error()}
	}

	fmt.Printf("id new entry %d", idNewEntry)
	pathToDB, err := CreateNewPathToDB(idNewEntry, task.FilesName, task.Contents, strings.Replace(task.Header, " ", "", -1))

	if err != nil {
		return entity.Answer{Status: err.Error()}
	}
	cmd = fmt.Sprintf("update tasks set path_to_bd = '%s' where id_task = %d;", pathToDB, idNewEntry)
	fmt.Println(pathToDB)
	_, err = r.store.db.Exec(cmd)
	if err != nil {
		return entity.Answer{Status: err.Error()}
	}
	return entity.Answer{Status: "200 OK"}
}

func CreateNewPathToDB(id int64, filenames, contents []string, header string) (string, error) {
	pathToDB := "./internal/databases/" + strconv.Itoa(int(id)) + "." + header

	err := os.MkdirAll(pathToDB, 0755)
	if err != nil {
		return "", err
	}

	for i := 0; i < len(filenames); i++ {
		pathToFile := pathToDB + "/" + filenames[i]
		file, err := os.Create(pathToFile)
		if err != nil {
			return "", err
		}
		defer file.Close()

		_, err = file.WriteString(contents[i])
		if err != nil {
			return "", nil
		}
	}
	return pathToDB, nil
}

func (r *DBRepository) ChangeStatus(id, status int, level string) entity.Answer {
	if status == 1 {
		_, err := r.store.db.Exec("update tasks set status = $1, task_level = $2 where id_task = $3", status, level, id)

		if err != nil {
			return entity.Answer{Status: err.Error()}
		}

		return entity.Answer{Status: "200 OK"}
	} else {
		var pathToCsv string
		err := r.store.db.QueryRow("select path_to_bd from tasks where id_task = $1;", id).Scan(&pathToCsv)
		if err != nil {
			return entity.Answer{Status: err.Error()}
		}
		err = os.RemoveAll(pathToCsv)
		if err != nil {
			return entity.Answer{Status: err.Error()}
		}

		_, err = r.store.db.Exec("delete from tasks where id_task = $1", id)

		if err != nil {
			return entity.Answer{Status: err.Error()}
		}

		return entity.Answer{Status: "200 OK"}
	}
}
