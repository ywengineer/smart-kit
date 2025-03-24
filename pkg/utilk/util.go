package utilk

// DefaultIfEmpty if v is empty then return def
func DefaultIfEmpty[T string | []any | map[any]any](v T, def T) T {
	if v == nil || len(v) == 0 {
		return def
	}
	return v
}

// DefaultIfNil 函数接收一个指针和一个默认值
// 如果指针为 nil，则返回默认值；否则返回指针指向的值
func DefaultIfNil[T any](ptr *T, def T) T {
	if ptr == nil {
		return def
	}
	return *ptr
}
