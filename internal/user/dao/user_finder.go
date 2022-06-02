package dao

import (
	userQuery "boilerplate/internal/user/dao/query"
	"boilerplate/internal/user/dto"
	"context"
	"gorm.io/gorm"
)

type UserFinder struct {
	db *gorm.DB
}

func NewUserFinder(db *gorm.DB) *UserFinder {
	return &UserFinder{db: db}
}

func (f *UserFinder) One(ctx context.Context, id string) *dto.User {
	query := f.CreateQuery(ctx)
	query.Id(id)

	return f.OneByQuery(query)
}

func (f *UserFinder) OneByQuery(query *userQuery.UserQuery) *dto.User {
	var user *dto.User
	query.Build().Limit(1).Scan(&user)

	return user
}

func (f *UserFinder) ListByQuery(query *userQuery.UserQuery, count int) []*dto.User {
	var users []*dto.User
	query.Build().Limit(count).Scan(&users)

	return users
}

func (f *UserFinder) CreateQuery(ctx context.Context) *userQuery.UserQuery {
	return userQuery.NewUserQuery(ctx, f.db)
}
