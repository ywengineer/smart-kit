package config

import (
	"context"
	"encoding/base64"
	"gitee.com/ywengineer/smart-kit/pkg/utilk"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type Auth struct {
	Realm   string            `json:"realm" yaml:"realm" redis:"realm"`
	UserKey string            `json:"userKey" yaml:"userKey" redis:"userKey"`
	Users   map[string]string `json:"users" yaml:"users" redis:"users"`
}

func (auth Auth) findUser(value string) (string, bool) {
	for user, password := range auth.Users {
		v := "Basic " + base64.StdEncoding.EncodeToString(utilk.S2b(user+":"+password))
		if v == value {
			return user, true
		}
	}
	return value, false
}

func BasicAuth() app.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote(p.Auth.Realm)
	return func(ctx context.Context, c *app.RequestContext) {
		// Search user in the slice of allowed credentials
		user, found := p.Auth.findUser(c.Request.Header.Get("Authorization"))
		if !found {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		c.Set(p.Auth.UserKey, user)
	}
}
