package apps

import (
	"github.com/hertz-contrib/jwt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenToken(t *testing.T) {
	j := NewJwt(JwtConfig{
		Key:         "test",
		Realm:       "Smart-test",
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: "smart-id",
	}, nil)
	token, _, err := j.TokenGenerator(jwt.MapClaims{
		"id": 123456,
	})
	assert.Nil(t, err)
	jt, err := j.ParseTokenString(token)
	assert.Nil(t, err)
	assert.True(t, true)
	assert.Nil(t, jt.Claims.Valid())
	t.Logf("%+v", jt)
}
