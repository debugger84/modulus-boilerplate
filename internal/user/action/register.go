package action

import (
	"boilerplate/internal/user/dto"
	"boilerplate/internal/user/service"
	"context"
	application "github.com/debugger84/modulus-application"
	"net/http"
)

type RegisterRequest struct {
	Name  string `json:"name"  validate:"required,min=3,max=50,alphaunicode"`
	Email string `json:"email"  validate:"required,email,max=150"`
}
type RegisterResponse struct {
	Id string `json:"id"`
}

type RegisterAction struct {
	runner       *application.ActionRunner
	registration *service.Registration
}

func NewRegisterAction(runner *application.ActionRunner, registration *service.Registration) *RegisterAction {
	return &RegisterAction{runner: runner, registration: registration}
}

func (a *RegisterAction) Handle(w http.ResponseWriter, r *http.Request) {
	a.runner.Run(w, r, func(ctx context.Context, request any) application.ActionResponse {
		return a.process(ctx, request.(*RegisterRequest))
	}, &RegisterRequest{})
}

func (a *RegisterAction) process(ctx context.Context, request *RegisterRequest) application.ActionResponse {
	user := dto.User{
		Name:  request.Name,
		Email: request.Email,
	}
	result, err := a.registration.Register(ctx, user)
	if err != nil {
		return application.NewUnprocessableEntityResponse(ctx, err)
	}
	return application.NewSuccessCreationResponse(
		RegisterResponse{
			Id: result.Id,
		},
	)
}
