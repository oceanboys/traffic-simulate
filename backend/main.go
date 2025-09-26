package main

import (
	"backend/routers"
	"backend/utils"
	"fmt"
	"os"

	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 先初始化数据库
	if err := utils.InitDatabase(); err != nil {
		fmt.Printf("数据库初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 创建数据库表
	if err := utils.CreateTables(); err != nil {
		fmt.Printf("创建数据库表失败: %v\n", err)
		os.Exit(1)
	}

	// 手动初始化路由（在数据库初始化后）
	routers.Init()

	beego.Run()
}
