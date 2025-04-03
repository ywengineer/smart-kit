// Code generated by hertz generator.

package passport

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/ywengineer/smart-kit/passport/internal"
	"github.com/ywengineer/smart-kit/passport/internal/middleware"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _bindMw() []app.HandlerFunc {
	return middleware.JwtWithValidate(middleware.IsUserMatch(internal.UserTypePlayer))
}

func _loginMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _registerMw() []app.HandlerFunc {
	return nil
}
