package utils

import (
	"backend/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func InitDatabase() error {
	// 读取配置
	dbHost, _ := beego.AppConfig.String("db.host")
	dbPort, _ := beego.AppConfig.String("db.port")
	dbUser, _ := beego.AppConfig.String("db.user")
	dbPassword, _ := beego.AppConfig.String("db.password")
	dbName, _ := beego.AppConfig.String("db.name")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 注册数据库
	err := orm.RegisterDataBase("default", "mysql", dataSource)
	if err != nil {
		logs.Error("数据库连接失败: ", err)
		return err
	}

	// 注册模型
	orm.RegisterModel(new(models.RoadSegment), new(models.GPSData), new(models.TrafficAlert), new(models.Vehicle))

	// 自动建表（开发环境）
	runMode, _ := beego.AppConfig.String("runmode")
	if runMode == "dev" {
		orm.RunSyncdb("default", false, true)
	}

	logs.Info("数据库初始化成功")
	return nil
}
