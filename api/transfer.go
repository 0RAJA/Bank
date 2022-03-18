package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/0RAJA/Bank/db/sqlc"
	"github.com/0RAJA/Bank/pkg/app"
	"github.com/0RAJA/Bank/pkg/app/errcode"
	"github.com/0RAJA/Bank/pkg/bind"
	"github.com/0RAJA/Bank/pkg/token"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

//判断是否存在account以及对应的currency是否相对应
func (s *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (*db.Account, bool) {
	response := app.NewResponse(ctx)
	account, err := s.store.GetAccountForUpdate(ctx, accountID)
	if err != nil {
		settings.Logger.Info(err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			response.ToErrorResponse(errcode.NotFoundErr)
		} else {
			if pqErr, ok := err.(*pq.Error); ok {
				switch pqErr.Code.Name() {
				case "foreign_key_violation":
					response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(err.Error()))
					return nil, false
				}
			}
			response.ToErrorResponse(errcode.ServerErr)
		}
		return nil, false
	}
	if account.Currency != currency {
		err := errcode.InvalidParamsErr.WithDetails(fmt.Sprintf("account [%d] currency mismatch %s vs %s", account.ID, account.Currency, currency))
		settings.Logger.Info(err.Error())
		response.ToErrorResponse(err)
		return nil, false
	}
	return &account, true
}

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	var param createTransferRequest
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		settings.Logger.Info(bind.FormatBindErr(errs))
		response.ToErrorResponse(errcode.InvalidParamsErr.WithDetails(errs.Errors()...))
		return
	}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	account1, ok := s.validAccount(ctx, param.FromAccountID, param.Currency)
	if !ok {
		return
	}
	if payload.UserName != account1.Owner {
		response.ToErrorResponse(errcode.InsufficientPermissionsErr)
		return
	}
	_, ok = s.validAccount(ctx, param.ToAccountID, param.Currency)
	if !ok {
		return
	}
	account, err := s.store.TransferTx(ctx, db.TransferTxParams{
		FromAccountID: param.FromAccountID,
		ToAccountID:   param.ToAccountID,
		Amount:        param.Amount,
	})
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

type getTransferParams struct {
	ID int64 `form:"id" binding:"required,min=1"`
}

func (s *Server) getTransfer(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := getTransferParams{}
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		err := errcode.InvalidParamsErr.WithDetails(errs.Errors()...)
		settings.Logger.Info(err.Error())
		response.ToErrorResponse(err)
		return
	}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	transfer, err := s.store.GetTransfer(ctx, db.GetTransferParams{
		ID:       param.ID,
		Username: payload.UserName,
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
	response.ToResponse(transfer)
}

type ListTransfersParams struct {
	FromAccountID int64 `form:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `form:"to_account_id" binding:"required,min=1"`
	Page          int   `form:"page" binding:"min=1"`
	PageSize      int   `form:"page_size" binding:"min=1"`
}

func (s *Server) ListTransfers(ctx *gin.Context) {
	response := app.NewResponse(ctx)
	param := ListTransfersParams{Page: app.GetPage(ctx), PageSize: app.GetPageSize(ctx)}
	valid, errs := bind.BindAndValid(ctx, &param)
	if !valid {
		err := errcode.InvalidParamsErr.WithDetails(errs.Errors()...)
		settings.Logger.Info(err.Error())
		response.ToErrorResponse(err)
		return
	}
	payload := ctx.MustGet(settings.TokenSetting.AuthorizationKey).(*token.Payload)
	accounts, err := s.store.ListTransfers(ctx, db.ListTransfersParams{
		Username:      payload.UserName,
		FromAccountID: param.FromAccountID,
		ToAccountID:   param.ToAccountID,
		Limit:         int32(param.PageSize),
		Offset:        int32(app.GetPageOffset(param.Page, param.PageSize)),
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
