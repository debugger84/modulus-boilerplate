package resolver

import (
	"boilerplate/internal/graph/model"
	"boilerplate/internal/user/storage"
	"boilerplate/internal/user/transformer"
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
)

type QueryResolver struct {
	finder *storage.Queries
}

func NewQueryResolver(finder *storage.Queries) *QueryResolver {
	return &QueryResolver{finder: finder}
}

func (r *QueryResolver) User(ctx context.Context, reqId string) (*model.User, error) {
	id, _ := uuid.FromString(reqId)
	user, err := r.finder.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		} else {
			return nil, errors.New("db error")
		}
	}
	return transformer.TransformUser(user), nil
}

func (r *QueryResolver) Users(ctx context.Context, first int, after *string) (*model.UserList, error) {
	var users []storage.User
	var err error
	var cursor *transformer.UsersListCursor
	if after != nil && *after != "" {
		cursor = transformer.NewUsersListCursorFromString(*after)
	}
	if cursor == nil {
		users, err = r.finder.GetUsersFirstPage(ctx, int32(first+1))
		if err != nil {
			return nil, errors.New("db error")
		}
	} else {
		id, _ := uuid.FromString(cursor.Id)
		users, err = r.finder.GetUsersAfterCursor(
			ctx, storage.GetUsersAfterCursorParams{
				RegisteredAt:   cursor.RegisteredAt,
				RegisteredAt_2: cursor.RegisteredAt,
				ID:             id,
				Limit:          int32(first + 1),
			},
		)
		if err != nil {
			return nil, errors.New("db error")
		}
	}

	return transformer.TransformUserList(
		users, first, func(user storage.User) string {
			return transformer.NewUsersListCursor(user).ToString()
		},
	), nil
}
