package validator

import (
	"errors"
	"github.com/samber/lo"
)

func Every(args ...interface{}) error {
	if len(args) <= 1 {
		return nil
	}
	keys := args[1:]
	if m, ok := args[0].([]string); ok {
		for _, v := range keys {
			if !lo.Contains(m, v.(string)) {
				return errors.New("not contains: " + v.(string))
			}
		}
		return nil
	} else if a, ok := args[0].(map[string]string); ok {
		for _, k := range keys {
			if v, ok := a[k.(string)]; !ok || len(v) == 0 {
				return errors.New("not contains: " + k.(string))
			}
		}
		return nil
	} else {
		return errors.New("only support []string or map[string]string")
	}
}
