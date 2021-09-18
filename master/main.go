package main

import (
	"context"
	"dcs-sched/master/config"
	"dcs-sched/master/router"
	"flag"
	"fmt"
	"github.com/carter115/gslog"
	"os"
	"runtime"
)

func init() {
	// 加载配置
	configFiel := flag.String("config", "./config/config.yaml", "configure file")
	flag.Parse()

	if err := config.InitConfig(*configFiel); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	// 加载日志组件
	logConfig := gslog.Config{
		ProjectName: config.Config.Logger.Project,
		AppName:     config.Config.Logger.App,
		FileName:    config.Config.Logger.File,
		Level:       config.Config.Logger.Level,
		Outputs:     config.Config.Logger.Outputs,
		Hooks:       config.Config.Logger.Hooks,
		EsServer:    config.Config.Logger.EsServer,
		StashServer: config.Config.Logger.StashServer,
	}
	gslog.InitLogger(logConfig)

	runtime.GOMAXPROCS(runtime.NumCPU())
}

// @title 分布式任务调度器
// @version 0.1
func main() {
	ctx := context.Background()
	gslog.Infof(ctx, "load config: %+v", config.Config)

	if err := router.Server(); err != nil {
		gslog.Errorf(ctx, "server is shutdown: %+v", err)
	}
}
