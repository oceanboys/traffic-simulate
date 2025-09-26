package controllers

import (
	"backend/models"
	"backend/services"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

// RoadController 路段控制器
type RoadController struct {
	web.Controller
	roadService *services.RoadService
}

func NewRoadController() *RoadController {
	return &RoadController{
		roadService: services.NewRoadService(),
	}
}

// GetAllRoads 获取所有路段
func (c *RoadController) GetAllRoads() {
	roads, err := c.roadService.GetAllRoads()
	if err != nil {
		c.CustomAbort(500, "Failed to get roads: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    roads,
	}
	c.ServeJSON()
}

// GetRoad 获取单个路段
func (c *RoadController) GetRoad() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.CustomAbort(400, "Invalid road ID")
		return
	}

	road, err := c.roadService.GetRoadByID(uint(id))
	if err != nil {
		c.CustomAbort(404, "Road not found")
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"data":    road,
	}
	c.ServeJSON()
}

// CreateRoad 创建路段
func (c *RoadController) CreateRoad() {
	var road models.RoadSegment
	if err := c.ParseForm(&road); err != nil {
		c.CustomAbort(400, "Invalid form data")
		return
	}

	if road.Name == "" {
		c.CustomAbort(400, "Road name cannot be empty")
		return
	}

	if err := c.roadService.CreateRoad(&road); err != nil {
		c.CustomAbort(500, "Failed to create road: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Road created successfully",
		"data":    road,
	}
	c.ServeJSON()
}

// UpdateRoad 更新路段
func (c *RoadController) UpdateRoad() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.CustomAbort(400, "Invalid road ID")
		return
	}

	var road models.RoadSegment
	if err := c.ParseForm(&road); err != nil {
		c.CustomAbort(400, "Invalid form data")
		return
	}

	road.ID = uint(id)
	if err := c.roadService.UpdateRoad(&road); err != nil {
		c.CustomAbort(500, "Failed to update road: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Road updated successfully",
		"data":    road,
	}
	c.ServeJSON()
}

// DeleteRoad 删除路段
func (c *RoadController) DeleteRoad() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.CustomAbort(400, "Invalid road ID")
		return
	}

	if err := c.roadService.DeleteRoad(uint(id)); err != nil {
		c.CustomAbort(500, "Failed to delete road: "+err.Error())
		return
	}

	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Road deleted successfully",
	}
	c.ServeJSON()
}
