package main

import (
	"context"
	"flag"
	"go_webapp/global"
	"go_webapp/internal/dao/mysql"
	"go_webapp/internal/dao/redis"
	"go_webapp/internal/routers"
	"go_webapp/pkg/logger"
	"go_webapp/pkg/settings"
	"go_webapp/pkg/snowflake"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	//命令行参数
	port    string
	runMode string
	config  string
)

func init() {
	err := setupSetting()
	if err != nil {
		global.Logger.Fatal("init.setupSetting failed", zap.Error(err))
	}

	err = setupLogger()
	if err != nil {
		global.Logger.Fatal("init.setupLogger failed", zap.Error(err))
	}

	err = setupSnowflake()
	if err != nil {
		global.Logger.Fatal("init.setupSnowflake failed", zap.Error(err))
	}
	err = setupMySql()
	if err != nil {
		global.Logger.Fatal("init:MySql init failed", zap.Error(err))
	}

	err = setupRedis()
	if err != nil {
		global.Logger.Fatal("init:Redis init failed", zap.Error(err))
	}
}

// @title 抖声
// @version 1.0
// @description 青训营项目 · 组名 - 大师我悟了
func main() {
	global.Logger.Debug("一切运行正常")

	//启动服务
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           ":" + global.ServerSetting.Port,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	//优雅重启
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("s.ListenAndServe err: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

	defer redis.Close()
}

//
//  setupFlag
//  @Description: 载入命令行参数
//  @return error
//
func setupFlag() error {
	flag.StringVar(&port, "port", "", "启动端口")
	flag.StringVar(&runMode, "mode", "", "启动模式")
	flag.StringVar(&config, "config", "configs/", "指定要使用的配置文件路径")
	//flag.BoolVar(&isVersion, "version", false, "编译信息")
	flag.Parse()

	return nil
}

//
//  setupSetting
//  @Description: 载入各部分的设置，绑定到结构体
//  @return error
//
func setupSetting() error {
	setting, err := settings.NewSetting(strings.Split(config, ",")...)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Log", &global.LogSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("MySql", &global.MySqlSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Redis", &global.RedisSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	if port != "" {
		global.ServerSetting.Port = port
	}
	if runMode != "" {
		global.ServerSetting.RunMode = runMode
	}
	return nil
}

func setupLogger() error {
	var err error
	global.Logger, err = logger.InitLogger(global.LogSetting, global.ServerSetting.RunMode)
	global.AccessLogger, err = logger.InitAccessLogger(global.LogSetting, global.ServerSetting.RunMode)
	return err
}

func setupSnowflake() error {
	return snowflake.Init(global.ServerSetting.StartTime, global.ServerSetting.MachineId)
}

func setupMySql() error {
	var err error
	global.DBEngine, err = mysql.Init(global.MySqlSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupRedis() error {
	return redis.Init(global.RedisSetting)
}
