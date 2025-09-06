package config

import (
	"context"
	"net/http"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
)

type Auth struct {
	Realm   string            `json:"realm" yaml:"realm" redis:"realm"`
	UserKey string            `json:"userKey" yaml:"userKey" redis:"userKey"`
	Users   map[string]string `json:"users" yaml:"users" redis:"users"`
}

func BasicAuth() app.HandlerFunc {
	realm := "Basic realm=" + strconv.Quote(p.Auth.Realm)
	return func(ctx context.Context, c *app.RequestContext) {
		// Search user in the slice of allowed credentials
		user, found := p.Auth.Users[c.Request.Header.Get("Authorization")]
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
