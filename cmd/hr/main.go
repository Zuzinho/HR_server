package main

import (
	"HR/pkg/env"
	"HR/pkg/forecaster"
	"HR/pkg/handler"
	"HR/pkg/middleware"
	"HR/pkg/repos"
	"database/sql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// Main - Главная функция
func main() {
	// Инициализируем переменные, получаемые из конфигурационных данных
	linksConfig := env.MustLinksConfig()
	port := env.MustPort()
	maxOpenConns := env.MustMaxOpenConns()
	db, err := sql.Open("postgres", env.MustDBConnString())
	if err != nil {
		log.Fatal(err)
	}

	// Ставим максимальное количество подключений к бд
	db.SetMaxOpenConns(maxOpenConns)

	// Проверяем подключение к бд
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("pinged connection to db")

	// Инициализируем класс-реализатор интерфейса UsersRepo
	usersRepo := repos.NewDatabaseUsersRepo(db)

	// Инициализируем класс-реализатор интерфейса IForecaster
	fc := forecaster.NewForecaster(linksConfig)

	// Инициализируем класс-обработчик запросов
	userHandler := handler.NewHandler(usersRepo, fc)

	router := mux.NewRouter()

	// Обрабатываем запросы по шаблону
	// Адрес запроса - /users
	// Возможные параметры - limit (максимальное получаемое количество записей за раз)
	// name, surname, patronymic, age (число), gender (male или female) и nation
	router.HandleFunc("/users", userHandler.GetFilterPagination).Methods("GET")
	// Адрес запроса - /users/{id записи (число)}
	router.HandleFunc("/users/{user_id:[0-9]+}", userHandler.DeleteUser).Methods("DELETE")
	// Адрес запроса - /users/{id записи (число)}
	// Возможные параметры - name, surname, patronymic, age (число), gender (male или female), nation
	router.HandleFunc("/users/{user_id:[0-9]+}", userHandler.UpdateUser).Methods("UPDATE")
	// Адрес запроса - /users
	// Тело запроса - {
	// 	name string
	//  surname string
	//  patronymic string (необязательно)
	// }
	router.HandleFunc("/users", userHandler.AddUser).Methods("POST")

	// Добавляем слой логгирования
	api := middleware.Logging(router)
	// Добавляем слой восстановления от паники
	api = middleware.RecoverPanic(api)

	log.Printf("starting on port %s", port)

	// Запускаем сервер
	log.Println(http.ListenAndServe(":"+port, api))
}
