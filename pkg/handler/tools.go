package handler

import (
	"HR/pkg/models/user"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
)

type (
	// updatesMap - user type for using url.Query map
	updatesMap map[string][]string
	// response - struct for storing response with inserted id
	response struct {
		InsertedID int `json:"inserted_id"`
	}
)

// write marshals data in json and writes it to ResponseWriter
func (handler *Handler) write(w http.ResponseWriter, data any) error {
	log.Printf("writing to response writer %v", data)

	buf, err := json.Marshal(data)
	if err != nil {
		http.Error(w, handler.jsonError(err), http.StatusInternalServerError)
		return nil
	}

	_, err = w.Write(buf)
	return err
}

// jsonError logs error and marshals error to json
func (handler *Handler) jsonError(err error) string {
	log.Println(err)

	return fmt.Sprintf("{\"error\": \"%s\"}", err.Error())
}

// paginate paginates users list to several lists with max length limit and writes these lists to enrichCh
func (handler *Handler) paginate(users *user.EnrichedUsers, limit int, enrichCh chan user.EnrichedUsers) {
	defer close(enrichCh)

	// Если limit равен 0, то записываем все данные без пагинации
	if limit == 0 {
		log.Println("writing all users")

		enrichCh <- *users
		return
	}

	sub := float64(len(*users)) / float64(limit)

	for i := 0; i < int(math.Ceil(sub)); i++ {
		from := i * limit
		to := (i + 1) * limit
		if to > len(*users) {
			to = len(*users)
		}

		log.Printf("writing users[%d, %d)", from, to)

		enrichCh <- (*users)[from:to]
	}
}

// fillQueryOptions fills queryBuilder part of sql query by pattern and delimiter and args values from updatesMap
// pattern must contain 2 placeholders for field name and it`s number
// Example: updatesMap {"name": "Nikita", "age": "20"}, pattern " %s=$%d" ; queryBuilder {" name=$1 %delimiter% age=$2"}, args {"Nikita", 20}
func (handler *Handler) fillQueryOptions(updatesMap updatesMap, queryBuilder *strings.Builder, args *[]any, pattern, delimiter string) error {
	firstVal := len(*args) + 1

	count := firstVal
	for k, v := range updatesMap {
		// Если у параметра k нет значений, то идем на следующую итерацию
		if len(v[0]) == 0 {
			log.Printf("value of '%s' parametr - empty set", k)
			continue
		}

		// Преобразуем k в тип FieldName и проверяем корректность значения
		fieldName, err := user.NewFieldName(k)
		if err != nil {
			return err
		}

		if count > firstVal {
			// Если не первая итерация, то дописываем разделитель
			queryBuilder.WriteString(delimiter)
		}

		// Вписываем в queryBuilder pattern и передаем ему параметры fieldName и номер параметра в queryBuilder
		queryBuilder.WriteString(fmt.Sprintf(pattern, fieldName, count))
		val, err := fieldName.ConvertQueryParam(v[0])
		if err != nil {
			return err
		}

		*args = append(*args, val)
		count++
	}

	return nil
}
