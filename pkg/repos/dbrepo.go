package repos

import (
	"HR/pkg/models/user"
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// DatabaseUsersRepo - implementor of UsersRepo interface
type DatabaseUsersRepo struct {
	DB *sql.DB
}

// NewDatabaseUsersRepo - constructor
func NewDatabaseUsersRepo(db *sql.DB) *DatabaseUsersRepo {
	return &DatabaseUsersRepo{
		DB: db,
	}
}

// GetUsersByQuery gets users from db with filters in queryBuilder and args
func (repo *DatabaseUsersRepo) GetUsersByQuery(ctx context.Context, queryBuilder *strings.Builder, args []any) (*user.EnrichedUsers, error) {
	log.Printf("getting users from db with args %v", args)

	// Выполняем запрос к бд
	// Если фильтр пустой, то вернутся все записи
	// Иначе queryBuilder запишет названия колонок и placeholder`ы для значений, а args передаст значения для запроса
	rows, err := repo.DB.QueryContext(ctx,
		fmt.Sprintf("select id, name, surname, patronymic, age, gender, nation from users where 1=1 %s", queryBuilder.String()),
		args...)
	if err != nil {
		return nil, err
	}

	users := make(user.EnrichedUsers, 0)
	// Читаем результат, пока он не достигнет конца
	for rows.Next() {
		u := new(user.EnrichedUser)

		// Записываем в поля переменной u значения из результата
		err = rows.Scan(&u.ID, &u.Name, &u.Surname, &u.Patronymic, &u.Age, &u.Gender, &u.Nation)
		if err != nil {
			log.Printf("user %v, err '%s'", *u, err.Error())
			continue
		}

		users.Append(u)
	}
	rows.Close()

	return &users, nil
}

// DeleteByID deletes user from db by id
func (repo *DatabaseUsersRepo) DeleteByID(ctx context.Context, id int) error {
	log.Printf("deleting from db user with id '%d'", id)

	// Выполняем запрос к бд для удаления записи по id
	_, err := repo.DB.ExecContext(ctx, "delete from users where id = $1", id)

	return err
}

// Update updates user in db by id = args[0]
func (repo *DatabaseUsersRepo) Update(ctx context.Context, queryBuilder *strings.Builder, args []any) error {
	log.Printf("updating users with args %v", args)

	// Выполняем запрос к бд
	// Первым значением в списке args должен быть id записи, оставшиеся - значения полей для обновления
	_, err := repo.DB.ExecContext(ctx, fmt.Sprintf("update users set%s where id=$1", queryBuilder.String()), args...)

	return err
}

// AddUser adds user in db
func (repo *DatabaseUsersRepo) AddUser(ctx context.Context, user *user.EnrichedUser) (int, error) {
	log.Printf("adding in db of user %v", user)

	var insertedID int

	// Выполняем запрос к бд и возвращаем вставленный id
	err := repo.DB.QueryRowContext(ctx, "insert into users (name, surname, patronymic, age, gender, nation) "+
		"values ($1, $2, $3, $4, $5, $6) returning id;",
		user.Name, user.Surname, user.Patronymic, user.Age, user.Gender, user.Nation).Scan(&insertedID)

	return insertedID, err
}
