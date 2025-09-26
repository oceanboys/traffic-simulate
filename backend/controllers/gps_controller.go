package controllers

import (
	"backend/models"
	"backend/services"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

// GPSController GPS控制器
type GPSController struct {
	web.Controller
	gpsService *services.GPSService
}

func NewGPSController() *GPSController {
	return &GPSController{
		gpsService: services.NewGPSService(),
	}
}

// CreateGPSData 创建GPS数据
func (c *GPSController) CreateGPSData() {
	var gpsData models.GPSData
	if err := c.ParseForm(&gpsData); err != nil {
		c.CustomAbort(400, "Invalid form data")
		return
	}

	if gpsData.VehicleID == "" {
		c.CustomAbort(400, "Vehicle ID cannot be empty")
		return
	}
	if gpsData.Longitude == 0 || gpsData.Latitude == 0 {
		c.CustomAbort(400, "Longitude and latitude cannot be zero")
		return
	}

	if err := c.gpsService.CreateGPSData(&gpsData); err != nil {
		c.CustomAbort(500, "Failed to create GPS data: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "GPS data created successfully",
	}
	c.ServeJSON()
}

// GetGPSDataByRoad 根据路段获取GPS数据
func (c *GPSController) GetGPSDataByRoad() {
	roadIdStr := c.Ctx.Input.Param(":roadId")
	roadId, err := strconv.ParseUint(roadIdStr, 10, 32)
	if err != nil {
		c.CustomAbort(400, "Invalid road ID")
		return
	}

	minutes := 30 // 默认获取最近30分钟的数据
	if c.GetString("minutes") != "" {
		if m, err := strconv.Atoi(c.GetString("minutes")); err == nil {
			minutes = m
		}
	}

	gpsData, err := c.gpsService.GetGPSDataByRoad(uint(roadId), minutes)
	if err != nil {
		c.CustomAbort(500, "Failed to get GPS data: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    gpsData,
	}
	c.ServeJSON()
}

// GetGPSDataByVehicle 根据车辆获取GPS数据
func (c *GPSController) GetGPSDataByVehicle() {
	vehicleId := c.Ctx.Input.Param(":vehicleId")
	if vehicleId == "" {
		c.CustomAbort(400, "Vehicle ID cannot be empty")
		return
	}

	limit := 100 // 默认限制100条
	if c.GetString("limit") != "" {
		if l, err := strconv.Atoi(c.GetString("limit")); err == nil {
			limit = l
		}
	}

	gpsData, err := c.gpsService.GetGPSDataByVehicle(vehicleId, limit)
	if err != nil {
		c.CustomAbort(500, "Failed to get GPS data: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    gpsData,
	}
	c.ServeJSON()
}
