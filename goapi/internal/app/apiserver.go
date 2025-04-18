package apiserver

import (
	"encoding/json"
	"goapi/internal/entity"
	"goapi/internal/repository"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	logrus "github.com/sirupsen/logrus"
)

// Инициализация структуры
type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	db     *repository.DataBase
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
		return nil
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
func (s *APIServer) handleTestGet(w http.ResponseWriter, r *http.Request) {
	var user entity.UserData
	var users entity.ContextData

	user = entity.UserData{0, 7, 123131, 131231, 600, "rejected"}
	users.Data = append(users.Data, user)

	w.Header().Set("Content-type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", " ")
	err := encoder.Encode(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

/* пример построения post запроса
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
// инициализация роутера
func (s *APIServer) configureRouter() {
	/// определение маршрутов
	s.router.HandleFunc("/test", s.handleTestGet).Methods("GET")
}
