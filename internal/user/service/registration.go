package service

import (
	"boilerplate/internal/user/db"
	"context"
	application "github.com/debugger84/modulus-application"
	"github.com/gofrs/uuid"
	"time"
)

const emailExists application.ErrorIdentifier = "emailExists"

type RegisterRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Registration struct {
	finder *db.UserFinder
	saver  *db.UserSaver
	logger application.Logger
}

func NewRegistration(finder *db.UserFinder, saver *db.UserSaver, logger application.Logger) *Registration {
	return &Registration{finder: finder, saver: saver, logger: logger}
}

// Register returns emailExists error
func (r Registration) Register(ctx context.Context, request RegisterRequest) (*db.User, error) {
	if r.emailExist(ctx, request.Email) {
		return nil, application.NewCommonError(emailExists, "not unique email")
	}
	id, _ := uuid.NewV6()
	user := db.User{
		ID:           id,
		Name:         request.Name,
		Email:        request.Email,
		RegisteredAt: time.Now(),
		Settings:     nil,
		Contacts:     []string{"Contact 1"},
	}

	user.Settings = &db.Settings{Incognito: true}
	err := r.saver.Create(ctx, user)
	if err != nil {
		r.logger.Error(ctx, err.Error())
		return nil, err
	}

	return &user, nil
}

func (r Registration) emailExist(ctx context.Context, email string) bool {
	query := r.finder.CreateQuery(ctx)
	query.Email(email)
	user := r.finder.OneByQuery(query)

	return user != nil
}
