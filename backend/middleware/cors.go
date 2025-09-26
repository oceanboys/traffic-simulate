package middleware

import (
	"github.com/beego/beego/v2/server/web/context"
)

func CorsMiddleware(ctx *context.Context) {
	ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ctx.Output.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	ctx.Output.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	if ctx.Input.Method() == "OPTIONS" {
		ctx.Output.SetStatus(200)
		ctx.Output.Body([]byte(""))
		return
	}
}
