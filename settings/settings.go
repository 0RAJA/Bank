package settings

import (
	"github.com/0RAJA/Bank/configs"
	"github.com/0RAJA/Bank/pkg/logger"
)

var (
	ServerSetting    = new(configs.Server)
	AppSetting       = new(configs.App)
	LogSetting       = new(configs.Log)
	PostgresSetting  = new(configs.Postgres)
	EmailSetting     = new(configs.Email)
	PagelinesSetting = new(configs.Pagelines)
	TokenSetting     = new(configs.Token)
)

var (
	Logger = new(logger.Log)
)
