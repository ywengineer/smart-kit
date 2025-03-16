// Code generated by hertz generator. DO NOT EDIT.

package router

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	mgr "github.com/ywengineer/smart-kit/passport/biz/router/mgr"
	mgr_pst "github.com/ywengineer/smart-kit/passport/biz/router/mgr/pst"
	passport "github.com/ywengineer/smart-kit/passport/biz/router/passport"
)

// GeneratedRegister registers routers generated by IDL.
func GeneratedRegister(r *server.Hertz) {
	//INSERT_POINT: DO NOT DELETE THIS LINE!
	mgr_pst.Register(r)

	mgr.Register(r)

	passport.Register(r)
}
