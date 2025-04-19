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
}

func (r *DBTaskRepository) ExecuteCommand(cmd string) (entity.Answer, error) {
	fmt.Println(r.answerTaskPath)

	content, err := os.ReadFile(r.answerTaskPath)

	if err != nil {
		return entity.Answer{Status: "Not read answer file"}, err
	}
	answer := strings.Split(string(content), ",")
	fmt.Println(cmd)
	var count int
	rows, err := r.store.db.Query(cmd)

	if err != nil {
		return entity.Answer{}, err
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return entity.Answer{Status: "Ошибка при чтении колонок"}, err
	}

	values := make([]interface{}, len(columns))
	pointers := make([]interface{}, len(columns))
	for i := range values {
		pointers[i] = &values[i]
	}
	for rows.Next() {
		err := rows.Scan(pointers...)
		if err != nil {
			return entity.Answer{}, err
		}
		dbValue := fmt.Sprintf("%v", values[0])
		for i := 0; i < len(answer); i++ {
			if dbValue == answer[i] {
				count++
			}
		}
	}
	fmt.Printf("len(answer) - %d\ncount - %d", len(answer), count)
	if count == len(answer) {
		return entity.Answer{Status: "200 OK"}, err
	} else {
		return entity.Answer{Status: "Wrong answer"}, nil
	}
}
func (r *DBTaskRepository) CreateDBForTask(taskNumber int) (entity.ContextTables, error) {
	pathToTables := "./internal/databases"
	var infoTables entity.ContextTables
	dbName := "task"
	var pathFiles []string
	entries, err := os.ReadDir(pathToTables)
	if err != nil {
		return entity.ContextTables{}, err
	}

	for _, entry := range entries {
		tmp := strings.Split(entry.Name(), ".")
		if tmp[0] == strconv.Itoa(taskNumber) {
			dbName = dbName + tmp[0]
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
