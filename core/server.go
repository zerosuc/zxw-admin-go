package core

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"time"

	"server/global"
	"server/initialize"
)

type Server struct {
}

func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "k8s-admin",
		Long: "The k8s-admin server controller is a daemon that embeds the core control loops.",
		Run: func(cmd *cobra.Command, args []string) {
			s := Server{}
			// 加载配置所有的配置文件
			s.Init()
			s.RunServer()
		},
	}
	return cmd
}

func (s *Server) Init() {
	// 所有全局变量在最初就初始化了; 这样可以满足并发安全；如果考虑多副本可以使用once
	global.Viper = Viper() // 初始化viper
	global.Log = Zap()     // 初始化zap日志
	zap.ReplaceGlobals(global.Log)
	global.Db = initialize.Gorm() // gorm连接数据库
	initialize.Redis()            // 初始化redis

	if global.Db == nil {
		global.Log.Error("mysql连接失败，退出程序")
		os.Exit(127)
	} else {
		initialize.RegisterTables(global.Db) // 初始化表
		initialize.MigrateData(global.Db)
		//前面库和表已经初始化好了才可以
		global.Cron = initialize.InitCron() // 初始化cron
		initialize.CheckCron()              // start cron entry, if exists
		// 程序结束前关闭数据库链接
		global.Db.DB()
		//defer db.Close()
	}
}
func (s *Server) RunServer() {
	addr := fmt.Sprintf("%s:%d", global.Config.System.Host, global.Config.System.Port)
	router := initialize.Routers()
	srv := http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	global.Log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error("Server Shutdown", zap.Error(err))
	}
	global.Log.Info("Server exiting")
}
