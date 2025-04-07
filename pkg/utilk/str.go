package utilk

import "fmt"

// ToString 把任意基本类型转换为字符串
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}
	if v, ok := value.(string); ok {
		return v
	} else if v, ok := value.(*string); ok {
		return *v
	} else {
		return fmt.Sprintf("%v", value)
	}
}
