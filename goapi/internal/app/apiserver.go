package apiserver

import (
	"encoding/json"
	"fmt"
	"goapi/internal/entity"
	"goapi/internal/repository"
	"net/http"

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
		AllowedOrigins:   []string{"http://127.0.0.1:5500"}, // Разрешенные домены
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
	answer, err := s.db.Data().AddUsers(user.Login, string(hash))
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
	logrus.Info("Route /createDB: POST request")
	if err := s.configureDBTask(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	infoTables, err := s.dbTask.Data().CreateDBForTask(1)

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
	s.dbTask.Close()
	answer, err := s.db.Data().DestroyDBTask("task1")

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

// инициализация роутера
func (s *APIServer) configureRouter() {
	/// определение маршрутов
	s.router.HandleFunc("/addUser", s.handleAddUser).Methods("POST")
	s.router.HandleFunc("/validateUser", s.handleValidateUser).Methods("POST")
	s.router.HandleFunc("/createDB", s.handleCreateDB).Methods("GET")
	s.router.HandleFunc("/dropDB", s.handleDropDB).Methods("Get")
	s.router.HandleFunc("/executeCommand", s.handleExecuteCommand).Methods("POST")
}
