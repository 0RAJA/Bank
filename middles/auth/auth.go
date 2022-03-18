package auth

import (
	"github.com/0RAJA/Bank/pkg/app"
	"github.com/0RAJA/Bank/pkg/app/errcode"
	"github.com/0RAJA/Bank/pkg/token"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware(maker token.Maker) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		response := app.NewResponse(ctx)
		authorizationHeader := ctx.GetHeader(settings.TokenSetting.AuthorizationKey)
		if len(authorizationHeader) == 0 {
			response.ToErrorResponse(errcode.UnauthorizedAuthNotExistErr)
			ctx.Abort()
			return
		}
		fields := strings.SplitN(authorizationHeader, " ", 2)
		if len(fields) != 2 || strings.ToLower(fields[0]) != settings.TokenSetting.AuthorizationType {
			response.ToErrorResponse(errcode.UnauthorizedAuthNotExistErr)
			ctx.Abort()
			return
		}
		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken)
		if err != nil {
			response.ToErrorResponse(errcode.UnauthorizedAuthNotExistErr.WithDetails(err.Error()))
			ctx.Abort()
			return
		}
		ctx.Set(settings.TokenSetting.AuthorizationKey, payload)
		ctx.Next()
	}
}
