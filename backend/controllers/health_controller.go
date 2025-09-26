package controllers

import "time"

type HealthController struct {
	BaseController
}

// @Title GetHealth
// @Description 鍋ュ悍妫€鏌ユ帴鍙?// @Success 200 {string} string "鏈嶅姟鐘舵€?
// @router /health [get]
func (c *HealthController) GetHealth() {
	c.Success(map[string]string{
		"status":    "ok",
		"service":   "backend",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	})
}
