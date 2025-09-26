package controllers

import "github.com/beego/beego/v2/server/web"

type BaseController struct {
	web.Controller
}

// JsonResponse 缁熶竴JSON鍝嶅簲鏍煎紡
func (c *BaseController) JsonResponse(code int, message string, data interface{}) {
	response := map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
	c.Data["json"] = response
	c.ServeJSON()
}

// Success 鎴愬姛鍝嶅簲
func (c *BaseController) Success(data interface{}) {
	c.JsonResponse(200, "success", data)
}

// Error 閿欒鍝嶅簲
func (c *BaseController) Error(message string, code int) {
	if code == 0 {
		code = 500
	}
	c.JsonResponse(code, message, nil)
}
