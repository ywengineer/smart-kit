package oauths

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

type Oauth map[string]map[string]string

const oauthType = "type"

func (o Oauth) Get(id string) (AuthFacade, error) {
	if af, ok := facadeMap[id]; ok {
		return af, nil
	} else if authProp, ok := o[id]; !ok {
		return nil, errors.New("oauth facade not found: " + id)
	} else if ot, ok := authProp[oauthType]; !ok {
		return nil, errors.New("missing oauth facade type: " + id)
	} else {
		ot = strings.ToLower(ot)
		switch ot {
		case "wx":
			af = NewWxAuth(authProp["app-id"], authProp["app-secret"])
		case "qq":
			af = NewQQAuth(authProp["app-id"], authProp["app-secret"], authProp["redirect-url"])
		case "smart":
			af = &anoAuth{}
		case "steam":
			af = NewSteamWebAuth(authProp["app-id"], authProp["app-secret"])
		case "google":
			af = NewGoogleAuth(authProp["app-id"])
		case "apple":
			af = NewAppleAuth(authProp["app-id"])
		default:
			return nil, errors.New(fmt.Sprintf("unsupported oauth type [%s] for auth facade [%s] ", ot, id))
		}
		var err error
		if af, err = af.Validate(id); err == nil {
			facadeMap[id] = af
		}
		return af, err
	}
}
