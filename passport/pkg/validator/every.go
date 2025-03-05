package validator

import (
	"fmt"
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
				return fmt.Errorf("not contains: %v", v)
			}
		}
		return nil
	} else if a, ok := args[0].(map[string]string); ok {
		for _, k := range keys {
			if v, ok := a[k.(string)]; !ok || len(v) == 0 {
				return fmt.Errorf("not contains: %v", k)
			}
		}
		return nil
	} else {
		return fmt.Errorf("only support []string or map[string]string")
	}
}
