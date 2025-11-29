package apps

import (
	"context"

	"gitee.com/ywengineer/smart-kit/pkg/logk"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/swagger"
	"github.com/swaggo/swag"
)

func RedocHandler(h *server.Hertz) app.HandlerFunc {
	h.GET("/swagger/doc.json", func() func(c context.Context, ctx *app.RequestContext) {
		config := swagger.Config{
			URL:                      "doc.json",
			DocExpansion:             "list",
			InstanceName:             swag.Name,
			Title:                    "Swagger UI",
			DefaultModelsExpandDepth: 1,
			DeepLinking:              true,
			PersistAuthorization:     false,
			Oauth2DefaultClientID:    "",
		}
		r, err := swag.ReadDoc(config.InstanceName)
		if err != nil {
			logk.Errorf("load swagger doc failed: %v", err)
			return func(c context.Context, ctx *app.RequestContext) {
				ctx.AbortWithMsg(err.Error(), consts.StatusInternalServerError)
			}
		}
		var doc map[string]interface{}
		if err := sonic.UnmarshalString(r, &doc); err != nil {
			logk.Errorf("unmarshal swagger doc failed: %v", err)
			return func(c context.Context, ctx *app.RequestContext) {
				ctx.AbortWithMsg(err.Error(), consts.StatusInternalServerError)
			}
		}
		return func(c context.Context, ctx *app.RequestContext) {
			ctx.JSON(consts.StatusOK, doc)
		}
	}())
	//
	return func(c context.Context, ctx *app.RequestContext) {
		ctx.Header("Content-Type", "text/html; charset=utf-8")
		ctx.Response.AppendBodyString(`
		<!DOCTYPE html>
		<html>
		  <head>
			<title>SEO API Documentation</title>
			<!-- needed for adaptive design -->
			<meta charset="utf-8"/>
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
			<!--
			Redoc doesn't change outer page styles
			-->
			<style>
			  body {
				margin: 0;
				padding: 0;
			  }
			</style>
		  </head>
		  <body>
			<redoc spec-url='/seo-app/swagger/doc.json'></redoc>
			<script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
		  </body>
		</html>
		`)
	}
}
