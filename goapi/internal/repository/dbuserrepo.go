package repository

import (
	"fmt"
	"goapi/internal/entity"
	"strings"

	"github.com/sirupsen/logrus"
)

// структура для взаимодействия с бд
type DBUserRepository struct {
	store  *DataBaseUser
	dbName string
}

func (r *DBUserRepository) CreateDb(login string) entity.Answer {
	r.dbName = "userdb" + login
	fmt.Println(r.dbName)
	_, err := r.store.db.Exec("CREATE DATABASE " + r.dbName + ";")
	if err != nil {
		r.store.Close()
		connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", r.dbName)
		r.store.OpenNew(connect)
		logrus.Info("Database ready")
		return entity.Answer{Status: "200 OK"}
	}
	r.store.Close()
	connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", r.dbName)
	r.store.OpenNew(connect)
	logrus.Info("Database ready")
	return entity.Answer{Status: "200 OK"}
}

func (r *DBUserRepository) ExecuteCommand(cmd string) entity.ResultExecute {
	typeRequest := CheckCommandForExecute(cmd)
	var result entity.ResultExecute
	if typeRequest == "query" {
		data, err := r.store.db.Query(cmd)
		if err != nil {
			return entity.ResultExecute{Status: err.Error()}
		}
		defer data.Close()
		columns, err := data.Columns()
		if err != nil {
			return entity.ResultExecute{Status: "Ошибка при чтении колонок"}
		}
		result.Columns = columns
		var dataQuery []string
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		for data.Next() {
			err := data.Scan(pointers...)
			if err != nil {
				return entity.ResultExecute{Status: "Ошибка чтения данных"}
			}
			dbVal := fmt.Sprintf("%v", values)
			dataQuery = append(dataQuery, dbVal)
		}
		var queryData [][]string
		for i := 0; i < len(dataQuery); i++ {
			data := strings.Replace(strings.Replace(dataQuery[i], "[", "", -1), "]", "", -1)
			tmp := strings.Split(data, " ")
			queryData = append(queryData, tmp)
		}
		result.Data = queryData
		result.Status = "200 OK"
	} else {
		_, err := r.store.db.Exec(cmd)
		if err != nil {
			return entity.ResultExecute{Status: err.Error()}
		}
		result.Status = "200 OK"
	}
	return result
}
func CheckCommandForExecute(cmd string) string {
	var dropCmd = []string{"insert", "delete", "update", "alter", "drop"}
	cmd = strings.ToLower(cmd)
	var tmp []string
	check := "query"
	tmp = strings.Split(cmd, " ")
	for i := 0; i < len(dropCmd); i++ {
		tmpCmd := dropCmd[i]
		for j := 0; j < len(tmp); j++ {
			if tmpCmd == tmp[j] {
				check = "exec"
				return check
			}
		}
	}
	return check
}

func (r *DBUserRepository) GetDbName() string {
	return r.dbName
}
func (r *DBUserRepository) GetInfoTables() (entity.ContextTables, error) {
	var infoTables entity.ContextTables
	nameTable, err := r.store.db.Query("SELECT table_name " +
		"FROM information_schema.tables " +
		"WHERE table_schema = 'public' " +
		"ORDER BY table_name",
	)
	if err != nil {
		fmt.Println("ssss")
		return entity.ContextTables{}, err
	}
	defer nameTable.Close()
	fmt.Println("ok")
	for nameTable.Next() {
		var infoTable entity.InfoTable
		err := nameTable.Scan(&infoTable.TableName)
		if err != nil {
			fmt.Println("pl")
			return entity.ContextTables{}, err
		}
		infoTables.Tables = append(infoTables.Tables, infoTable)
	}
	fmt.Println("ok")
	for i := 0; i < len(infoTables.Tables); i++ {

		cmd := fmt.Sprintf("SELECT "+
			"column_name as name, "+
			"CASE "+
			"WHEN data_type = 'timestamp without time zone' THEN 'DATETIME' "+
			"ELSE upper(data_type) "+
			"END as type, "+
			"CASE WHEN is_nullable = 'NO' THEN 'NOT NULL' ELSE '' END as nullable, "+
			"EXISTS ( "+
			"SELECT 1 "+
			"FROM information_schema.key_column_usage kcu "+
			"JOIN information_schema.table_constraints tc "+
			"ON kcu.constraint_name = tc.constraint_name "+
			"WHERE tc.constraint_type = 'PRIMARY KEY' "+
			"AND kcu.table_name = '%s' "+
			"AND kcu.column_name = information_schema.columns.column_name "+
			") as is_pk, "+
			"EXISTS ( "+
			"SELECT 1 "+
			"FROM information_schema.key_column_usage kcu "+
			"JOIN information_schema.table_constraints tc "+
			"ON kcu.constraint_name = tc.constraint_name "+
			"WHERE tc.constraint_type = 'FOREIGN KEY' "+
			"AND kcu.table_name = '%s' "+
			"AND kcu.column_name = information_schema.columns.column_name "+
			") as is_fk "+
			"FROM information_schema.columns "+
			"WHERE table_name = '%s' "+
			"ORDER BY ordinal_position ",
			infoTables.Tables[i].TableName,
			infoTables.Tables[i].TableName,
			infoTables.Tables[i].TableName,
		)

		typeColumns, err := r.store.db.Query(cmd)
		if err != nil {
			fmt.Println("ne ok")
			return entity.ContextTables{}, err
		}
		defer typeColumns.Close()

		for typeColumns.Next() {
			var data entity.TypeColumn
			var notNull string
			err := typeColumns.Scan(
				&data.NameColumn,
				&data.TypeColumn,
				&notNull,
				&data.Is_PK,
				&data.Is_FK,
			)
			if err != nil {
				fmt.Println("Sssssss")
				return entity.ContextTables{}, err
			}
			var tmp []string
			if data.Is_PK {
				tmp = append(tmp, "PK")
			}
			if data.Is_FK {
				tmp = append(tmp, "FK")
			}
			tmp = append(tmp, data.TypeColumn)
			infoTables.Tables[i].TableColumns = append(infoTables.Tables[i].TableColumns, data.NameColumn)
			infoTables.Tables[i].TableColumnsTypes = append(infoTables.Tables[i].TableColumnsTypes, tmp)
		}
	}
	return infoTables, nil
}
