package controllers

import (
	"backend/models"
	"backend/services"
	"encoding/json"
	"github.com/beego/beego/v2/server/web"
	"time"
)

// TrafficController 交通控制器
type TrafficController struct {
	web.Controller
	trafficService *services.TrafficService
}

// NewTrafficController 创建交通控制器
func NewTrafficController() *TrafficController {
	return &TrafficController{
		trafficService: services.NewTrafficService(),
	}
}

// GetRealTimeTraffic 获取实时交通数据
// @Title GetRealTimeTraffic
// @Description 获取实时交通数据
// @Success 200 {object} services.TrafficSummary
// @router /realtime [get]
func (c *TrafficController) GetRealTimeTraffic() {
	summary := c.trafficService.GetRealTimeSummary()
	c.Data["json"] = summary
	c.ServeJSON()
}

// GetAlerts 获取告警列表
// @Title GetAlerts
// @Description 获取告警列表
// @Success 200 {array} models.TrafficAlert
// @router /alerts [get]
func (c *TrafficController) GetAlerts() {
	alerts := c.trafficService.GetRecentAlerts(50) // 最近50条告警
	c.Data["json"] = alerts
	c.ServeJSON()
}

// GetTrafficStats 获取交通统计
// @Title GetTrafficStats
// @Description 获取交通统计数据
// @Success 200 {object} services.TrafficStats
// @router /stats [get]
func (c *TrafficController) GetTrafficStats() {
	stats := c.trafficService.GetTrafficStats()
	c.Data["json"] = stats
	c.ServeJSON()
}

// GetCongestionData 获取拥堵数据
// @Title GetCongestionData
// @Description 获取拥堵路段数据
// @Success 200 {array} services.CongestionData
// @router /congestion [get]
func (c *TrafficController) GetCongestionData() {
	congestion := c.trafficService.GetCongestionData()
	c.Data["json"] = congestion
	c.ServeJSON()
}

// GetVehicleFlow 获取车流数据
// @Title GetVehicleFlow
// @Description 获取车流统计数据
// @Success 200 {object} services.VehicleFlow
// @router /flow [get]
func (c *TrafficController) GetVehicleFlow() {
	flow := c.trafficService.GetVehicleFlow()
	c.Data["json"] = flow
	c.ServeJSON()
}

// GetVehicles 获取车辆列表
// @Title GetVehicles
// @Description 获取实时车辆列表
// @Success 200 {array} models.Vehicle
// @router /vehicles [get]
func (c *TrafficController) GetVehicles() {
	vehicles := c.trafficService.GetVehicles()
	c.Data["json"] = vehicles
	c.ServeJSON()
}

// AddVehicle 添加车辆
// @Title AddVehicle
// @Description 添加新车辆到模拟系统
// @Success 200 {object} map[string]interface{}
// @router /vehicles [post]
func (c *TrafficController) AddVehicle() {
	var vehicle models.Vehicle
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &vehicle); err != nil {
		c.CustomAbort(400, "Invalid request body")
		return
	}

	// 设置默认值
	vehicle.ID = uint(time.Now().Unix())
	vehicle.CreatedAt = time.Now()
	vehicle.Status = "normal"

	// 添加到服务
	result := c.trafficService.AddVehicle(vehicle)

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "车辆添加成功",
		"data":    result,
	}
	c.ServeJSON()
}

// RemoveVehicle 移除车辆
// @Title RemoveVehicle
// @Description 从模拟系统移除车辆
// @Success 200 {object} map[string]interface{}
// @router /vehicles/:id [delete]
func (c *TrafficController) RemoveVehicle() {
	vehicleID := c.GetString(":id")
	if vehicleID == "" {
		c.CustomAbort(400, "Vehicle ID is required")
		return
	}

	success := c.trafficService.RemoveVehicle(vehicleID)

	c.Data["json"] = map[string]interface{}{
		"success": success,
		"message": "车辆移除成功",
	}
	c.ServeJSON()
}

// StartSimulation 开始模拟
// @Title StartSimulation
// @Description 开始交通流模拟
// @Success 200 {object} map[string]interface{}
// @router /simulation/start [post]
func (c *TrafficController) StartSimulation() {
	c.trafficService.StartSimulation()

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "模拟已开始",
	}
	c.ServeJSON()
}

// StopSimulation 停止模拟
// @Title StopSimulation
// @Description 停止交通流模拟
// @Success 200 {object} map[string]interface{}
// @router /simulation/stop [post]
func (c *TrafficController) StopSimulation() {
	c.trafficService.StopSimulation()

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "模拟已停止",
	}
	c.ServeJSON()
}

// GetSimulationStatus 获取模拟状态
// @Title GetSimulationStatus
// @Description 获取当前模拟状态
// @Success 200 {object} map[string]interface{}
// @router /simulation/status [get]
func (c *TrafficController) GetSimulationStatus() {
	status := c.trafficService.GetSimulationStatus()

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    status,
	}
	c.ServeJSON()
}
