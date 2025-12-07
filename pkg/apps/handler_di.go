package apps

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/samber/do/v2"
	dohttp "github.com/samber/do/v2/http"
)

// DiMonitor
// Use integrates the do library's web-based debugging interface with a Hertz router.
// This function sets up HTTP routes for the debugging UI, allowing you to inspect
// your DI container through a web browser.
//
// Parameters:
//   - router: The Hertz router group to add the debugging routes to
//   - injector: The injector instance to debug
//
// The function sets up the following routes:
//   - GET /: The main debugging interface home page
//   - GET /scope: Scope tree visualization with optional scope_id parameter
//   - GET /service: Service inspection with optional scope_id and service_name parameters
//
// Example:
//
//	router := hertz.New()
//	api := router.Group("/api")
//	debug := router.Group("/debug/di")
//
//	// Add the debugging interface
//	dohertz.Use(debug, injector)
//
//	// Your application routes
//	api.GET("/users", userHandler)
//
// The debugging interface will be available at /debug/di and provides:
//   - Visual representation of scope hierarchy
//   - Service dependency graphs
//   - Service inspection and debugging tools
//   - Navigation between different views
//
// Security:
// Do not expose this group publicly in production. Protect it with authentication
// (e.g., Basic Auth) and/or network restrictions, since it can leak internals
// about your application's DI graph. Attach auth middleware to the router group
// before calling Use.
func DiMonitor(router *route.RouterGroup, injector do.Injector) {
	basePathDo := router.BasePath()

	router.GET("", func(c context.Context, ctx *app.RequestContext) {
		output := fmt.Sprintf(
			`<!DOCTYPE html>
					<html>
						<head>
							<title>Dependency injection UI - samber/do</title>
						</head>
						<body>
							<h1>Welcome to do UI ✌️</h1>
							
							<ul>
								<li><a href="%s/debug/pprof">Profile</a></li>
								<li><a href="%s/scope">Inspect scopes</a></li>
								<li><a href="%s/service">Inspect services</a></li>
							</ul>
						</body>
					</html>`,
			basePathDo, basePathDo, basePathDo)
		response(ctx, output, nil)
	})

	router.GET("/scope", func(c context.Context, ctx *app.RequestContext) {
		scopeID := ctx.Query("scope_id")
		if scopeID == "" {
			url := fmt.Sprintf("%s/scope?scope_id=%s", basePathDo, injector.ID())
			ctx.Redirect(consts.StatusTemporaryRedirect, []byte(url))
			return
		}

		output, err := dohttp.ScopeTreeHTML(basePathDo, injector, scopeID)
		response(ctx, output, err)
	})

	router.GET("/service", func(c context.Context, ctx *app.RequestContext) {
		scopeID := ctx.Query("scope_id")
		serviceName := ctx.Query("service_name")

		if scopeID == "" || serviceName == "" {
			output, err := dohttp.ServiceListHTML(basePathDo, injector)
			response(ctx, output, err)
			return
		}

		output, err := dohttp.ServiceHTML(basePathDo, injector, scopeID, serviceName)
		response(ctx, output, err)
	})
}

func response(c *app.RequestContext, data string, err error) {
	c.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	if err != nil {
		c.Response.AppendBodyString(err.Error())
	} else {
		c.Response.AppendBodyString(data)
	}
}
