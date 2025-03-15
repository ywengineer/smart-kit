// Code generated by hertz generator. DO NOT EDIT.

package mgr

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	mgr "github.com/ywengineer/smart-kit/passport/biz/handler/mgr"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_mgr := root.Group("/mgr", _mgrMw()...)
		_mgr.POST("/sign", append(_signMw(), mgr.Sign)...)
		{
			_white_list := _mgr.Group("/white-list", _white_listMw()...)
			_white_list.GET("/add", append(_addMw(), mgr.Add)...)
			_white_list.GET("/page", append(_pageMw(), mgr.Page)...)
			_white_list.GET("/rm", append(_removeMw(), mgr.Remove)...)
		}
	}
}
