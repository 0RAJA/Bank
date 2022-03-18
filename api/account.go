package api

import (
	"database/sql"
	"errors"
	db "github.com/0RAJA/Bank/db/sqlc"
	"github.com/0RAJA/Bank/pkg/app"
	"github.com/0RAJA/Bank/pkg/app/errcode"
	"github.com/0RAJA/Bank/pkg/bind"
	"github.com/0RAJA/Bank/pkg/token"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,oneof=USD EUR RMB"`
}

func (s *Server) createAccount(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	var param createAccountRequest
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    payload.UserName,
		Currency: param.Currency,
	}
	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		settings.Logger.Info(err.Error())
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(err.Error()))
				return
			}
		}
		response.ToErrorResponse(errcode.ServerErr)
		return
	}
	response.ToResponse(account)
}

type getAccountParams struct {
	ID int64 `json:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := getAccountParams{ID: app.StrTo(ctx.Query("id")).MustInt64()}
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	account, err := s.store.GetAccountForUpdate(ctx, param.ID)
	if err != nil {
		settings.Logger.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			response.ToErrorResponse(errcode.NotFoundErr)
		} else {
			response.ToErrorResponse(errcode.ServerErr)
		}
		return
	}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	if account.Owner != payload.UserName {
		response.ToErrorResponse(errcode.InsufficientPermissionsErr)
		return
	}
	response.ToResponse(account)
}

type ListAccountsParams struct {
	Page     int `form:"page" binding:"min=1"`
	PageSize int `form:"page_size" binding:"min=1"`
}

func (s *Server) ListAccounts(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := ListAccountsParams{Page: app.GetPage(ctx), PageSize: app.GetPageSize(ctx)}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	accounts, err := s.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  int32(param.PageSize),
		Offset: int32(app.GetPageOffset(param.Page, param.PageSize)),
		Owner:  payload.UserName,
	})
	if err != nil {
		settings.Logger.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			response.ToErrorResponse(errcode.NotFoundErr)
		} else {
			response.ToErrorResponse(errcode.ServerErr)
		}
		return
	}
	response.ToResponseList(accounts, len(accounts))
}
