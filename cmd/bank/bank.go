package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github.com/0RAJA/Bank/api"
	db "github.com/0RAJA/Bank/db/sqlc"
	"github.com/0RAJA/Bank/pkg/app"
	"github.com/0RAJA/Bank/pkg/logger"
	"github.com/0RAJA/Bank/pkg/setting"
	"github.com/0RAJA/Bank/settings"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	address  = "127.0.0.1:8080"
)

func main() {
	conn, err := sql.Open(settings.PostgresSetting.DBDriver, fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%s", settings.PostgresSetting.UserName, settings.PostgresSetting.Password, settings.PostgresSetting.Address, settings.PostgresSetting.DBName, settings.PostgresSetting.Sslmode))
	if err != nil {
		log.Fatalln(err)
	}
	store := db.NewStore(conn)
	gin.SetMode(settings.ServerSetting.RunMode)
	server := api.NewServer(store)
	s := &http.Server{
		Addr:           settings.ServerSetting.Address,
		Handler:        server.Router,
		ReadTimeout:    settings.ServerSetting.ReadTimeout,
		WriteTimeout:   settings.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()
	gracefulExit(s) //优雅退出
	settings.Logger.Info("OVER")
}
func init() {
	if err := SetupSetting(); err != nil { //加载配置文件
		panic("init setting failed:" + err.Error())
	}
	initLog() //加载日志
	//初始化分页器
	app.Init(settings.PagelinesSetting.DefaultPage, settings.PagelinesSetting.DefaultPageSize, settings.PagelinesSetting.PageKey, settings.PagelinesSetting.PageSizeKey)
}

func initLog() {
	logger.Init(&logger.InitStruct{
		LogSavePath:   settings.LogSetting.LogSavePath,
		LogFileExt:    settings.LogSetting.LogFileExt,
		MaxSize:       settings.LogSetting.MaxSize,
		MaxBackups:    settings.LogSetting.MaxBackups,
		MaxAge:        settings.LogSetting.MaxAge,
		Compress:      settings.LogSetting.Compress,
		LowLevelFile:  settings.LogSetting.LowLevelFile,
		HighLevelFile: settings.LogSetting.HighLevelFile,
	})
	settings.Logger = logger.NewLogger(settings.LogSetting.Level)
}

var (
	configPaths string
	configName  string
	configType  string
)

func setupFlag() {
	//命令行参数绑定
	flag.StringVar(&configName, "name", "config", "配置文件名")
	flag.StringVar(&configType, "type", "yml", "配置文件类型")
	flag.StringVar(&configPaths, "path", "configs/", "指定要使用的配置文件路径")
	flag.Parse()
}

//优雅关机
func gracefulExit(s *http.Server) {
	//退出通知
	quit := make(chan os.Signal)
	//等待退出通知
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	settings.Logger.Info("ShutDown Server...")
	//给几秒完成剩余任务
	ctx, cancel := context.WithTimeout(context.Background(), settings.ServerSetting.DefaultContextTimeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil { //优雅退出
		settings.Logger.Info("Server forced to ShutDown,Err:" + err.Error())
	}
}

func SetupSetting() (err error) {
	setupFlag()
	newSetting, err := setting.NewSetting(configName, configType, strings.Split(configPaths, ",")...)
	if err != nil {
		return err
	}
	readSetting := func(k string, v interface{}) error {
		if err != nil {
			return err
		}
		return newSetting.ReadSection(k, v)
	}
	err = readSetting("Server", settings.ServerSetting)
	err = newSetting.ReadSection("App", settings.AppSetting)
	err = newSetting.ReadSection("Log", settings.LogSetting)
	err = newSetting.ReadSection("Postgres", settings.PostgresSetting)
	err = newSetting.ReadSection("Email", settings.EmailSetting)
	err = newSetting.ReadSection("Pagelines", settings.PagelinesSetting)
	err = newSetting.ReadSection("token", settings.TokenSetting)
	if err != nil {
		return err
	}
	return nil
}
