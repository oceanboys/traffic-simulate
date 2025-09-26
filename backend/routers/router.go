package routers

import (
	"backend/controllers"
	"github.com/beego/beego/v2/server/web"
)

func Init() {
	// 初始化控制器
	roadController := controllers.NewRoadController()
	gpsController := controllers.NewGPSController()
	healthController := &controllers.HealthController{}
	trafficController := controllers.NewTrafficController()

	// 健康检查
	web.Router("/api/health", healthController, "get:GetHealth")

	// 路段管理路由
	web.Router("/api/roads", roadController, "get:GetAllRoads")
	web.Router("/api/roads", roadController, "post:CreateRoad")
	web.Router("/api/roads/:id:int", roadController, "get:GetRoad")
	web.Router("/api/roads/:id:int", roadController, "put:UpdateRoad")
	web.Router("/api/roads/:id:int", roadController, "delete:DeleteRoad")

	// GPS数据路由
	web.Router("/api/gps", gpsController, "post:CreateGPSData")
	web.Router("/api/gps/road/:roadId:int", gpsController, "get:GetGPSDataByRoad")
	web.Router("/api/gps/vehicle/:vehicleId", gpsController, "get:GetGPSDataByVehicle")

	// 交通数据路由
	web.Router("/api/traffic/realtime", trafficController, "get:GetRealTimeTraffic")
	web.Router("/api/traffic/alerts", trafficController, "get:GetAlerts")
	web.Router("/api/traffic/stats", trafficController, "get:GetTrafficStats")
	web.Router("/api/traffic/congestion", trafficController, "get:GetCongestionData")
	web.Router("/api/traffic/flow", trafficController, "get:GetVehicleFlow")

	// 车辆管理路由
	web.Router("/api/vehicles", trafficController, "get:GetVehicles")
	web.Router("/api/vehicles", trafficController, "post:AddVehicle")
	web.Router("/api/vehicles/:id", trafficController, "delete:RemoveVehicle")

	// 模拟控制路由
	web.Router("/api/simulation/start", trafficController, "post:StartSimulation")
	web.Router("/api/simulation/stop", trafficController, "post:StopSimulation")
	web.Router("/api/simulation/status", trafficController, "get:GetSimulationStatus")
}
