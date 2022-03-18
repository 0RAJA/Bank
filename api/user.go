package api

import (
	"database/sql"
	db "github.com/0RAJA/Bank/db/sqlc"
	"github.com/0RAJA/Bank/pkg/app"
	"github.com/0RAJA/Bank/pkg/app/errcode"
	"github.com/0RAJA/Bank/pkg/bind"
	"github.com/0RAJA/Bank/pkg/utils"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserInfo struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserInfo(user db.User) UserInfo {
	return UserInfo{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	var param createUserRequest
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	hashPassword, err := utils.HashPassword(param.Password)
	if err != nil {
		errs := errcode.ServerErr.WithDetails(err.Error())
		settings.Logger.Info(errs.Error())
		response.ToErrorResponse(errs)
		return
	}
	arg := db.CreateUserParams{
		Username:       param.Username,
		HashedPassword: hashPassword,
		FullName:       param.FullName,
		Email:          param.Email,
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		settings.Logger.Info(err.Error())
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(err.Error()))
				return
			}
		}
		response.ToErrorResponse(errcode.ServerErr)
		return
	}
	token, err := s.maker.CreateToken(user.Username, settings.TokenSetting.Duration)
	if err != nil {
		response.ToErrorResponse(errcode.ServerErr)
		return
	}
	response.ToResponse(newUserInfoWithToken(token, newUserInfo(user)))
}

type getUserRequest struct {
	Username string `form:"username" binding:"required,alphanum"`
}

func (s *Server) GetUser(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := getUserRequest{}
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	user, err := s.store.GetUser(ctx, param.Username)
	if err != nil {
		settings.Logger.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			response.ToErrorResponse(errcode.NotFoundErr)
		} else {
			response.ToErrorResponse(errcode.ServerErr)
		}
		return
	}
	response.ToResponse(newUserInfo(user))
}

type loginsRequest struct {
	Username string `json:"username,omitempty" binding:"required,alphanum"`
	Password string `json:"password,omitempty" bindin:"required,min=6"`
}

type userInfoWithToken struct {
	Token string   `json:"token,omitempty"`
	User  UserInfo `json:"user"`
}

func newUserInfoWithToken(token string, user UserInfo) *userInfoWithToken {
	return &userInfoWithToken{Token: token, User: user}
}

func (s *Server) LoginUser(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := loginsRequest{}
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	user, err := s.store.GetUser(ctx, param.Username)
	if err != nil {
		settings.Logger.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			response.ToErrorResponse(errcode.LoginErr)
		} else {
			response.ToErrorResponse(errcode.ServerErr)
		}
		return
	}
	if err := utils.CheckPassword(param.Password, user.HashedPassword); err != nil {
		response.ToErrorResponse(errcode.LoginErr)
		return
	}
	token, err := s.maker.CreateToken(user.Username, settings.TokenSetting.Duration)
	if err != nil {
		response.ToErrorResponse(errcode.ServerErr)
		return
	}
	response.ToResponse(newUserInfoWithToken(token, newUserInfo(user)))
}
