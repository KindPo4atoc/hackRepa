package apiserver

import (
	"encoding/json"
	"fmt"
	"goapi/internal/entity"
	"goapi/internal/repository"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	logrus "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Инициализация структуры
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	db     *repository.DataBase
	dbTask *repository.DataBaseTask
	dbUser *repository.DataBaseUser
}

// конструктор
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

// Метод запуска API. Инициализирует все поля структуры
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	if err := s.configureDB(); err != nil {
		return fmt.Errorf("Ошибка инициализации основной БД")
	}

	s.logger.Info("Starting api server")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5500", "http://192.168.1.4:3000", "http://localhost:3000"}, // Разрешенные домены
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Разрешить куки и заголовки авторизации
		Debug:            true, // Логирование (опционально)
	})

	handler := c.Handler(s.router)
	return http.ListenAndServe(s.config.BindAddr, handler)
}

// инициализация логгера
func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

// инициализация бд
func (s *APIServer) configureDB() error {

	database := repository.New(s.config.DBConfig)
	if err := database.Open(); err != nil {
		return err
	}

	s.db = database
	s.db.Data()

	return nil
}

func (s *APIServer) configureDBTask() error {
	database := repository.NewTask(s.config.DBTaskConfig)
	if err := database.Open(); err != nil {
		return err
	}

	s.dbTask = database
	s.dbTask.Data()

	return nil
}

func (s *APIServer) configureDBUser() error {
	database := repository.NewUser(s.config.DBUserConfig)
	if err := database.Open(); err != nil {
		return err
	}

	s.dbUser = database
	s.dbUser.Data()

	return nil
}

/*
	пример построения post запроса

	func (s *APIServer) handlePredictModel(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Route /predict: POST request")
		var dataForPredict entity.UserData
		json.NewDecoder(r.Body).Decode(&dataForPredict)
		var tmp []float64
		tmp = append(tmp, float64(dataForPredict.IncomeAnnum))
		tmp = append(tmp, float64(dataForPredict.LoanAmount))
		tmp = append(tmp, float64(dataForPredict.LoanTerm))
		tmp = append(tmp, float64(dataForPredict.CibilScore))
		var data [][]float64
		data = append(data, tmp)
		predict, distance, dataLDA, err := s.model.Predict(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		s.predict = lda.NewPredict(predict, distance, dataLDA)
		w.Header().Set("Content-type", "application/json")
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		err = encoder.Encode(s.predict)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

// примеры построения Get запроса

	func (s *APIServer) handleGetConvData(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Route /getConvData: GET request")
		w.Header().Set("Content-type", "application/json")

		encoder := json.NewEncoder(w)
		encoder.SetIndent("", " ")
		err := encoder.Encode(s.model)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
*/
func (s *APIServer) handleValidateUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /existUser: POST request")
	var user entity.UserData
	json.NewDecoder(r.Body).Decode(&user)

	answer, err := s.db.Data().ValidateUser(user.Login, user.PasswordHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *APIServer) handleAddUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /addUser: POST request")
	var user entity.UserData
	json.NewDecoder(r.Body).Decode(&user)

	cost := 14
	hash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), cost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	answer, err := s.db.Data().AddUsers(user.Login, user.Email, string(hash))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleCreateDB(w http.ResponseWriter, r *http.Request) {
	// переделать в POST запрос create -> нужно чтобы передавали в EP для create номер задачи
	var task entity.Task
	fmt.Println()
	json.NewDecoder(r.Body).Decode(&task)
	fmt.Println(task)
	logrus.Info("Route /createDB: POST request")
	if err := s.configureDBTask(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	path, err := s.db.Data().GetPathCsv(task.IdTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	infoTables, err := s.dbTask.Data().CreateDb(task.IdTask, path)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(infoTables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleDropDB(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /dropDB: POST request")
	dbName := s.dbTask.Data().GetDbName()
	countActiveUsers, err := s.dbTask.Data().GetActivityUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var answer entity.Answer
	s.dbTask.Close()
	if countActiveUsers == 1 {
		answer, err = s.db.Data().DestroyDBTask(dbName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		answer = entity.Answer{Status: "200 OK"}
	}

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleExecuteCommand(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /executeCommand: POST request")
	var command entity.Command
	json.NewDecoder(r.Body).Decode(&command)
	fmt.Println(command.Cmd)
	answer, err := s.dbTask.Data().ExecuteCommand(command.Cmd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /getAllTasks: GET request")
	tasks, err := s.db.Data().GetAllTasks()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetTask(w http.ResponseWriter, r *http.Request) {
	idTask := mux.Vars(r)["item"]
	logrus.Info("Route /getTask/{item}: GET request")
	task, err := s.db.Data().GetTask(idTask)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetTasksByLevel(w http.ResponseWriter, r *http.Request) {
	level := mux.Vars(r)["item"]
	logrus.Info("Route /getTasksByLevel/{item}: GET request")
	task, err := s.db.Data().GetTasksByLevel(level)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetInfoTables(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /getInfoTable: GET request")
	infoTable, err := s.dbTask.Data().GetInfoTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(infoTable)

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(infoTable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetTasksToolCheck(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /getTasksToCheck: GET request")
	tasks, err := s.db.Data().GetTasksStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleAddTask(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /addTask: POST request")
	var taskForAdd entity.AddTask
	json.NewDecoder(r.Body).Decode(&taskForAdd)
	fmt.Println(taskForAdd)

	answer := s.db.Data().AddTask(taskForAdd)
	if answer.Status != "200 OK" {
		http.Error(w, answer.Status, http.StatusInternalServerError)
	}
	fmt.Println("ok")
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) handleChangeStatus(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /changeStatus/{id}/{status}/{level}: POST request")
	idTask, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	status, err := strconv.Atoi(mux.Vars(r)["status"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	level := mux.Vars(r)["level"]
	answer := s.db.Data().ChangeStatus(idTask, status, level)
	if answer.Status != "200 OK" {
		http.Error(w, answer.Status, http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleCreateDBUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /createDBUser: POST request")
	var user entity.UserData
	json.NewDecoder(r.Body).Decode(&user)
	fmt.Println(user)

	if err := s.configureDBUser(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	answer := s.dbUser.Data().CreateDb(user.Login)
	if answer.Status != "200 OK" {
		http.Error(w, answer.Status, http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func (s *APIServer) handleExecuteCommandUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /executeCommandUser: POST request")
	var command entity.Command
	json.NewDecoder(r.Body).Decode(&command)
	fmt.Println(command.Cmd)
	answer := s.dbUser.Data().ExecuteCommand(command.Cmd)
	if answer.Status != "200 OK" {
		http.Error(w, answer.Status, http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleDropDBUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /dropDBUser: GET request")
	dbName := s.dbUser.Data().GetDbName()
	var answer entity.Answer
	s.dbUser.Close()
	answer, err := s.db.Data().DestroyDBTask(dbName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(answer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleGetInfoTablesUser(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /getInfoTablesUser: GET request")
	infoTable, err := s.dbUser.Data().GetInfoTables()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Println(infoTable)

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(infoTable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *APIServer) handleDownloadFiles(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Route /downloadTaskTable/{id}: GET request")
	idTask, err := strconv.Atoi(mux.Vars(r)["id"])
	var task entity.AddTask
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	path, err := s.db.Data().GetPathCsv(idTask)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	files, err := os.ReadDir(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	for _, file := range files {
		pathToFile := path + file.Name()
		task.FilesName = append(task.FilesName, file.Name())
		content, err := os.ReadFile(pathToFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		dataFile := string(content)
		task.Contents = append(task.Contents, dataFile)
	}
	fmt.Println(task)

	w.Header().Set("Content-type", "application/json")

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err = encoder.Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// инициализация роутера
func (s *APIServer) configureRouter() {
	//проверить id задач
	/// определение маршрутов
	s.router.HandleFunc("/addUser", s.handleAddUser).Methods("POST")
	s.router.HandleFunc("/validateUser", s.handleValidateUser).Methods("POST")
	s.router.HandleFunc("/createDB", s.handleCreateDB).Methods("POST")
	s.router.HandleFunc("/getInfoTable", s.handleGetInfoTables).Methods("GET")
	s.router.HandleFunc("/dropDB", s.handleDropDB).Methods("GET")
	s.router.HandleFunc("/executeCommand", s.handleExecuteCommand).Methods("POST")
	s.router.HandleFunc("/getAllTasks", s.handleGetAllTasks).Methods("GET")
	s.router.HandleFunc("/getTasksToCheck", s.handleGetTasksToolCheck).Methods("GET")
	s.router.HandleFunc("/getTask/{item}", s.handleGetTask).Methods("GET")
	s.router.HandleFunc("/getTasksByLevel/{item}", s.handleGetTasksByLevel).Methods("GET")
	s.router.HandleFunc("/addTask", s.handleAddTask).Methods("POST")
	s.router.HandleFunc("/changeStatus/{id}/{status}/{level}", s.handleChangeStatus).Methods("GET")
	s.router.HandleFunc("/createDBUser", s.handleCreateDBUser).Methods("POST")
	s.router.HandleFunc("/dropDBUser", s.handleDropDBUser).Methods("GET")
	s.router.HandleFunc("/getInfoTablesUser", s.handleGetInfoTablesUser).Methods("GET")
	s.router.HandleFunc("/executeCommandUser", s.handleExecuteCommandUser).Methods("POST")
	s.router.HandleFunc("/downloadTaskTables/{id}", s.handleDownloadFiles).Methods("GET")
}

/*
todo ->песочницу
	 ->добавить везде cookie
	 ->одобрение/отклонение задачи -> бэк готов (/changeStatus)
	 ->собрать сайт
	 ->добавление автора к addTask
*/
