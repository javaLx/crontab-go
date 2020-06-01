package main

import (
	"crontab-go/master"
	"flag"
	"log"
	"runtime"
)

var (
	configFile string
)

//  初始化命令行参数
func initArgs() {
	// master -config ./config.json -xxx 123 -yyy ddd
	flag.StringVar(&configFile, "config", "./config.json", "config配置文件地址")
	flag.Parse()
}

// 初始化线程数量
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)
	// 初始化终端参数
	initArgs()
	// 初始化线程数
	initEnv()
	// 加载配置信息
	if err = master.InitConfig(configFile); err != nil {
		log.Fatal("初始化配置文件出错", err.Error())
	}
	// 初始化服务发现模块
	if err = master.InitWorkerMgr(); err != nil {
		log.Fatal("初始化服务发现模块失败", err.Error())
	}
	// 初始化任务管理器
	if err = master.InitJobManager(); err != nil {
		log.Fatal("初始化任务管理器失败", err.Error())
	}
	// 启动Api HTTP服务
	if err = master.ApiServer().Run(); err != nil {
		log.Fatal("启动Api服务失败", err.Error())
	}

}
