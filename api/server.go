package api

import (
	"github.com/0RAJA/Bank/configs"
	db "github.com/0RAJA/Bank/db/sqlc"
	"github.com/0RAJA/Bank/middles/auth"
	"github.com/0RAJA/Bank/middles/valid"
	"github.com/0RAJA/Bank/pkg/token"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)
import "github.com/gin-gonic/gin"

type Server struct {
	store      db.Store       //提供数据库交互
	Router     *gin.Engine    //提供路由，处理API请求
	serverConf configs.Server //
	maker      token.Maker
}

func NewServer(store *db.SqlStore) *Server {
	maker, err := token.NewPasetoMaker([]byte(settings.TokenSetting.Key))
	if err != nil {
		settings.Logger.Fatal(err.Error())
	}
	server := &Server{store: store, maker: maker}
	router := gin.Default()
	//注册自定校验
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", valid.ValidCurrency) //注册自定义标签
	}
	//TODO: register route
	auth := router.Group("/", auth.AuthMiddleware(server.maker))
	{
		auth.GET("/get_user", server.GetUser)
		auth.POST("/create_transfer", server.createTransfer)
		auth.GET("/get_transfer", server.getTransfer)
		auth.GET("/list_transfer", server.ListTransfers)
		auth.POST("/create_account", server.createAccount)
		auth.GET("/get_account", server.getAccount)
		auth.GET("/list_accounts", server.ListAccounts)
	}
	router.POST("/create_user", server.createUser)
	router.POST("/login", server.LoginUser)
	server.Router = router
	return server
}

func (s *Server) Start(address string) error {
	if err := s.Router.Run(address); err != nil {
		return err
	}
	return nil
}
