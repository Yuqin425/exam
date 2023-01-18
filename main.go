package main

import (
	"LanShan/api"
	"LanShan/dao/mysql"
	"LanShan/logger"
	"LanShan/router"
	"LanShan/settings"
	"LanShan/utils/snowflake"
	"fmt"
)

func main() {
	// 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	// 雪花算法生成分布式ID
	if err := snowflake.Init(1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	if err := api.InitTrans("zh"); err != nil {
		fmt.Printf("init validator Trans failed,err:%v\n", err)
		return
	}
	// 注册路由
	r := router.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
