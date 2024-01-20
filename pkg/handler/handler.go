package handler

import (
	"HR/pkg/forecaster"
	"HR/pkg/models/user"
	"HR/pkg/repos"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Handler handles queries
// Stored implementors of UsersRepo and IForecaster interfaces
type Handler struct {
	UsersRepo  repos.UsersRepo
	Forecaster forecaster.IForecaster
}

// NewHandler - constructor
func NewHandler(usersRepo repos.UsersRepo, forecaster forecaster.IForecaster) *Handler {
	return &Handler{
		UsersRepo:  usersRepo,
		Forecaster: forecaster,
	}
}

// GetFilterPagination writes in response writer users from repository with other filters and pagination
// Writes 200 OK with successful handling
// Writes 400 BadRequest with incorrect limit or field names values
// Writes 500 InternalServerError with server errors
func (handler *Handler) GetFilterPagination(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()

	// Получаем и сразу удаляем из параметров запроса параметр limit
	limitStr := values.Get("limit")
	values.Del("limit")

	var err error
	var limit int
	if limitStr != "" {
		// Если значение параметра limit не указан или пустой, то значение переменной limit = 0
		// Иначе получаем из строки значение переменной limit
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, handler.jsonError(err), http.StatusBadRequest)
			return
		}
	}

	queryBuilder := new(strings.Builder)
	args := make([]any, 0)

	// Заполняем strings.Builder частью sql запроса, args значениями sql запроса
	err = handler.fillQueryOptions(updatesMap(values), queryBuilder, &args, " and %s=$%d", "")
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusBadRequest)
		return
	}

	// Выполняем получение данных из репозитория
	users, err := handler.UsersRepo.GetUsersByQuery(r.Context(), queryBuilder, args)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	enrichCh := make(chan user.EnrichedUsers)

	// Запускаем в отдельной горутине пагинацию и записываем список EnrichedUsers в канал enrichCh
	// И закрываем канал при окончании данных
	go handler.paginate(users, limit, enrichCh)

	// Читаем данные из канала и записываем списки в ResponseWriter
	for list := range enrichCh {
		err = handler.write(w, list)
		if err != nil {
			http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser deletes user from repository by id
// Writes 200 OK with successful handling
// Writes 400 BadRequest with incorrect user_id value
// Writes 500 InternalServerError with server errors
func (handler *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Получаем user_id по ключу из запроса
	id, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusBadRequest)
		return
	}

	// Выполняем удаление данных из репозитория
	err = handler.UsersRepo.DeleteByID(r.Context(), id)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UpdateUser updates user in repository by id and another field names
// Writes 200 OK with successful handling
// Writes 400 BadRequest with incorrect user_id or field names values
// Writes 500 InternalServerError with server errors
func (handler *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Получаем user_id по ключу из запроса
	id, err := strconv.Atoi(vars["user_id"])
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusBadRequest)
		return
	}

	values := r.URL.Query()

	queryBuilder := new(strings.Builder)
	args := []any{id}

	// Заполняем strings.Builder частью sql запроса, args значениями sql запроса
	err = handler.fillQueryOptions(updatesMap(values), queryBuilder, &args, " %s=$%d", " ,")
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusBadRequest)
		return
	}

	// Выполняем обновление данных в репозитории
	err = handler.UsersRepo.Update(r.Context(), queryBuilder, args)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddUser adds user in repository
// Writes 200 OK with successful handling
// Writes 400 BadRequest with incorrect body request
// Writes 500 InternalServerError with server errors
func (handler *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	buf, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Преобразуем json из тела запроса в сущность User
	u := new(user.User)
	err = json.Unmarshal(buf, u)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	errCh := make(chan error)
	enrichUsCh := make(chan *user.EnrichedUser)
	defer close(errCh)
	defer close(enrichUsCh)

	var enrUser *user.EnrichedUser
	ctx := r.Context()

	// Запускаем в отдельной горутине обогащение данных при помощи IForecaster
	go handler.Forecaster.ForecastUser(u, enrichUsCh, errCh)

	// Ждем один из 3 сценариев
	select {
	// Запись получила доп. данные и успешно вернулась
	case enrUser = <-enrichUsCh:
	// Произошла ошибка в горутине
	case err = <-errCh:
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	// Контекст запроса завершился
	case <-ctx.Done():
		http.Error(w, handler.jsonError(ctx.Err()), http.StatusInternalServerError)
		return
	}

	// Выполняем добавление данных в репозиторий и получаем вставленный id
	id, err := handler.UsersRepo.AddUser(r.Context(), enrUser)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	resp := response{
		InsertedID: id,
	}

	// Записываем в ResponseWriter вставленный id
	err = handler.write(w, resp)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
