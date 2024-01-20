package repos

import (
	"HR/pkg/models/user"
	"context"
	"strings"
)

// UsersRepo - interface for getting, deleting, updating and adding users in repository
type UsersRepo interface {
	GetUsersByQuery(ctx context.Context, queryBuilder *strings.Builder, args []any) (*user.EnrichedUsers, error)
	DeleteByID(ctx context.Context, id int) error
	Update(ctx context.Context, queryBuilder *strings.Builder, args []any) error
	AddUser(ctx context.Context, user *user.EnrichedUser) (int, error)
}
