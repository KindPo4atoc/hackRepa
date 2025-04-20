package repository

import (
	"fmt"
	"goapi/internal/entity"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// структура для взаимодействия с бд
type DBTaskRepository struct {
	store          *DataBaseTask
	answerTaskPath string
	dbName         string
}

func CheckCommand(cmd string) bool {
	var dropCmd = []string{"insert", "delete", "update", "alter", "drop"}
	cmd = strings.ToLower(cmd)
	var tmp []string
	tmp = strings.Split(cmd, " ")
	for i := 0; i < len(dropCmd); i++ {
		check := false
		tmpCmd := dropCmd[i]
		for j := 0; j < len(tmp); j++ {
			if tmpCmd == tmp[j] {
				check = true
				return check
			}
		}
	}
	return false
}
func (r *DBTaskRepository) ExecuteCommand(cmd string) (entity.ResultExecute, error) {
	fmt.Println(r.answerTaskPath)
	if CheckCommand(cmd) {
		return entity.ResultExecute{Status: "Запрещенная команда"}, nil
	}
	var result entity.ResultExecute
	content, err := os.ReadFile(r.answerTaskPath)

	if err != nil {
		return entity.ResultExecute{Status: "Not read answer file"}, err
	}
	cmdForAnswer := string(content)
	fmt.Println(cmd)
	var count int
	rows, err := r.store.db.Query(cmd)

	if err != nil {
		return entity.ResultExecute{Status: err.Error()}, nil
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return entity.ResultExecute{Status: "Ошибка при чтении колонок"}, err
	}
	result.Columns = columns
	var dataQuery []string
	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			return entity.ResultExecute{Status: "Ошибка чтения данных"}, err
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
	rowsAnswer, err := r.store.db.Query(cmdForAnswer)
	if err != nil {
		return entity.ResultExecute{Status: "Ошибка при запросе"}, err
	}

	defer rowsAnswer.Close()
	columns, err = rowsAnswer.Columns()
	if err != nil {
		return entity.ResultExecute{Status: "Ошибка при чтении колонок"}, err
	}
	var dataAnswer []string
	valuesAnswer := make([]interface{}, len(columns))
	pointersAnswer := make([]interface{}, len(columns))
	for i := range valuesAnswer {
		pointersAnswer[i] = &valuesAnswer[i]
	}
	for rowsAnswer.Next() {
		err := rowsAnswer.Scan(pointersAnswer...)
		if err != nil {
			return entity.ResultExecute{Status: "Ошибка чтении данных"}, err
		}
		dbValAnswer := fmt.Sprintf("%v", valuesAnswer)
		dataAnswer = append(dataAnswer, dbValAnswer)
	}
	fmt.Println(dataAnswer)
	for i := 0; i < len(dataAnswer); i++ {
		tmp := dataAnswer[i]
		for j := 0; j < len(dataQuery); j++ {
			if tmp == dataQuery[j] {
				count++
			}
		}
	}
	if count == len(dataAnswer) && len(dataAnswer) == len(dataQuery) {
		result.Status = "200 OK"
		return result, nil
	} else {
		result.Status = "Wrong answer"
		return result, nil
	}
}

/*
	func (r *DBTaskRepository) CreateDBForTask(taskNumber int) (entity.ContextTables, error) {
		pathToTables := "./internal/databases"
		var infoTables entity.ContextTables
		dbName := "task" + strconv.Itoa(taskNumber)
		var pathFiles []string
		entries, err := os.ReadDir(pathToTables)
		if err != nil {
			return entity.ContextTables{}, err
		}

		for _, entry := range entries {
			tmp := strings.Split(entry.Name(), ".")
			if tmp[0] == strconv.Itoa(taskNumber) {
				pathToTables = pathToTables + "/" + entry.Name()
			}
		}
		files, err := os.ReadDir(pathToTables)
		if err != nil {
			return entity.ContextTables{}, err
		}
		for _, file := range files {
			tmp := pathToTables + "/" + file.Name()
			pathFiles = append(pathFiles, tmp)
		}
		fmt.Println("ok")
		r.dbName = dbName
		_, err = r.store.db.Exec("CREATE DATABASE " + dbName)
		if err != nil {
			r.store.Close()
			connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", dbName)
			r.store.OpenNew(connect)
			logrus.Info("Database ready")
			return entity.ContextTables{}, nil
		}

		r.store.Close()
		connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", dbName)
		r.store.OpenNew(connect)
		curDir, err := os.Executable()
		if err != nil {
			return entity.ContextTables{}, err
		}
		fmt.Println(curDir)
		for _, pathFile := range pathFiles {
			var infoTable entity.InfoTable
			fmt.Println(pathFile)
			tmp := strings.Split(pathFile, "/")
			if tmp[len(tmp)-1][len(tmp[len(tmp)-1])-4] != '.' {
				r.answerTaskPath = pathFile
				continue
			}
			strCreateTable := "CREATE TABLE " + strings.Split(tmp[len(tmp)-1], ".")[0] + " ("
			content, err := os.ReadFile(pathFile)
			if err != nil {
				return entity.ContextTables{}, err
			}
			infoTable.TableName = strings.Split(tmp[len(tmp)-1], ".")[0]
			strContent := strings.Split(string(content), "\n")
			row := strings.Split(strContent[0], ",")
			infoTable.TableColumns = row
			rowData := strings.Split(strContent[1], ",")
			for i := 0; i < len(row); i++ {
				if strings.HasPrefix(row[i], "id") {
					strCreateTable = strCreateTable + row[i] + " serial primary key not null"
					infoTable.TableColumnsTypes = append(infoTable.TableColumnsTypes, "PK")
					continue
				} else {
					//как определить что это вторичный ключ?
					typeColumn := CheckType(rowData[i])
					strCreateTable = strCreateTable + ", " + row[i] + " " + typeColumn
					infoTable.TableColumnsTypes = append(infoTable.TableColumnsTypes, typeColumn)
				}
			}
			strCreateTable = strCreateTable + ");"
			fmt.Println(strCreateTable)

			_, err = r.store.db.Exec(strCreateTable)
			if err != nil {
				return entity.ContextTables{}, err
			}
			infoTables.Tables = append(infoTables.Tables, infoTable)
			fullPathToCsv := GetPathCsv(curDir, pathFile)
			cmdCopy := fmt.Sprintf("COPY %s FROM '%s' DELIMITER ',' CSV HEADER;", strings.Split(tmp[len(tmp)-1], ".")[0], fullPathToCsv)
			fmt.Println(cmdCopy)
			_, err = r.store.db.Exec(cmdCopy)
			if err != nil {
				return entity.ContextTables{}, err
			}

		}
		return infoTables, nil
	}
*/
func GetPathCsv(currDir, pathToCsv string) string {
	currDirArr := strings.Split(currDir, "/")
	pathToCsvArr := strings.Split(pathToCsv, "/")
	var result string
	for i := 1; i < len(currDirArr); i++ {
		if i == len(currDirArr)-2 {
			result = result + "/" + currDirArr[i]
			for j := 1; j < len(pathToCsvArr); j++ {
				result = result + "/" + pathToCsvArr[j]
			}
			break
		} else {
			result = result + "/" + currDirArr[i]
		}
	}
	return result
}

func (r *DBTaskRepository) CreateDb(taskNumber int, path string) (entity.Answer, error) {
	r.dbName = "task" + strconv.Itoa(taskNumber)
	var pathFiles []string
	files, err := os.ReadDir(path)
	if err != nil {
		return entity.Answer{}, err
	}
	for _, file := range files {
		tmp := path + "/" + file.Name()
		if tmp[len(tmp)-4] != '.' {
			r.answerTaskPath = tmp
		}
		pathFiles = append(pathFiles, tmp)
	}
	fmt.Println(pathFiles)
	curDir, err := os.Executable()
	fmt.Println(curDir)
	if err != nil {
		return entity.Answer{}, err
	}

	_, err = r.store.db.Exec("CREATE DATABASE " + r.dbName)
	if err != nil {
		r.store.Close()
		connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", r.dbName)
		r.store.OpenNew(connect)
		logrus.Info("Database ready")
		return entity.Answer{Status: "200 OK"}, nil
	}
	r.store.Close()
	connect := fmt.Sprintf("postgres://kp:admin@localhost/%s?sslmode=disable", r.dbName)
	r.store.OpenNew(connect)
	for _, pathFile := range pathFiles {
		tmp := strings.Split(pathFile, "/")
		if tmp[len(tmp)-1] == "answer" {
			continue
		}
		strCreateTable := "CREATE TABLE " + strings.Split(tmp[len(tmp)-1], ".")[0] + " ("
		fmt.Println(strCreateTable)
		content, err := os.ReadFile(pathFile)
		if err != nil {
			return entity.Answer{}, err
		}
		strContent := strings.Split(string(content), "\n")
		fmt.Println(len(strContent))
		fmt.Println(strContent)
		row := strings.Split(strContent[0], ",")
		rowData := strings.Split(strContent[1], ",")
		fmt.Println(row)
		fmt.Println(rowData)
		for i := 0; i < len(row); i++ {
			if strings.Contains(row[i], "id") {
				strCreateTable = strCreateTable + row[i] + " serial primary key not null"
				fmt.Println(strCreateTable)
			} else {
				//как определить что это вторичный ключ?
				typeColumn := CheckType(rowData[i])
				strCreateTable = strCreateTable + ", " + row[i] + " " + typeColumn
				fmt.Println(strCreateTable)
			}
		}
		strCreateTable = strCreateTable + ");"
		fmt.Println("SSSSSS")
		fmt.Println(strCreateTable)

		_, err = r.store.db.Exec(strCreateTable)
		if err != nil {
			return entity.Answer{}, err
		}
		fullPathToCsv := GetPathCsv(curDir, pathFile)
		cmdCopy := fmt.Sprintf("COPY %s FROM '%s' DELIMITER ',' CSV HEADER;", strings.Split(tmp[len(tmp)-1], ".")[0], fullPathToCsv)
		fmt.Println(cmdCopy)
		_, err = r.store.db.Exec(cmdCopy)
		if err != nil {
			return entity.Answer{}, err
		}
	}

	return entity.Answer{Status: "200 OK"}, nil
}
func (r *DBTaskRepository) GetActivityUsers() (int, error) {
	var countUsers int
	count, err := r.store.db.Query("SELECT COUNT(*) FROM pg_stat_activity WHERE datname = current_database();")
	if err != nil {
		return -1, err
	}
	defer count.Close()
	for count.Next() {
		err := count.Scan(&countUsers)
		if err != nil {
			return -1, err
		}
	}
	if countUsers > 1 {
		return -1, nil
	} else {
		return 1, nil
	}

}
func (r *DBTaskRepository) GetInfoTables() (entity.ContextTables, error) {
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

func (r *DBTaskRepository) GetDbName() string {
	return r.dbName
}
func CheckType(data string) string {
	reInteger := regexp.MustCompile(`^[+-]?(0|[1-9]\d*)$`)
	reReal := regexp.MustCompile(`^[+-]?(?:0|[1-9]\d*)(?:[.,]\d+)?(?:[eE][+-]?\d+)?$|^[+-]?[.,]\d+(?:[eE][+-]?\d+)?$`)
	reDate := regexp.MustCompile(`^([0-9]{4})-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)

	if reInteger.MatchString(data) {
		return "integer"
	}
	if reReal.MatchString(data) {
		return "real"
	}
	if reDate.MatchString(data) {
		return "date"
	}
	return "text"
}
