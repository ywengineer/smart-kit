package utilk

import (
	"github.com/bytedance/sonic"
)

func UnmarshalJSON[T interface{}](bytes []byte, v *T) (*T, error) {
	err := sonic.Unmarshal(bytes, v)
	if err != nil {
		return nil, err
	}
	return v, err
}
