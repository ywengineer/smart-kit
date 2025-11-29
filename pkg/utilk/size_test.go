package utilk

import (
	"errors"
	"math"
	"testing"
)

func TestSize_MarshalJSON(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name     string
		size     Size
		expected string
	}{{
		name:     "0字节",
		size:     0,
		expected: "0B",
	}, {
		name:     "字节级别",
		size:     100,
		expected: "100.00B",
	}, {
		name:     "KB级别",
		size:     1024,
		expected: "1.00KB",
	}, {
		name:     "MB级别",
		size:     1024 * 1024,
		expected: "1.00MB",
	}, {
		name:     "GB级别",
		size:     1024 * 1024 * 1024,
		expected: "1.00GB",
	}, {
		name:     "TB级别",
		size:     1024 * 1024 * 1024 * 1024,
		expected: "1.00TB",
	}, {
		name:     "PB级别",
		size:     1024 * 1024 * 1024 * 1024 * 1024,
		expected: "1.00PB",
	}, {
		name:     "小数部分",
		size:     1536, // 1.5KB
		expected: "1.50KB",
	}, {
		name:     "接近边界值",
		size:     1023,
		expected: "1023.00B",
	}, {
		name:     "超出PB范围",
		size:     Size(math.MaxInt64),
		expected: "8192.00PB", // 最大值会显示为PB
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.size.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v, wantErr nil", err)
				return
			}
			if got := string(result); got != tt.expected {
				t.Errorf("MarshalJSON() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSize_UnmarshalJSON(t *testing.T) {
	// 成功的测试用例
	successTests := []struct {
		name     string
		input    string
		expected Size
	}{{
		name:     "无单位默认B",
		input:    "100",
		expected: 100,
	}, {
		name:     "B单位",
		input:    "100B",
		expected: 100,
	}, {
		name:     "小写b单位",
		input:    "100b",
		expected: 100,
	}, {
		name:     "KB单位",
		input:    "2KB",
		expected: 2 * 1024,
	}, {
		name:     "小写kb单位",
		input:    "2kb",
		expected: 2 * 1024,
	}, {
		name:     "简写K单位",
		input:    "3K",
		expected: 3 * 1024,
	}, {
		name:     "MB单位",
		input:    "4MB",
		expected: 4 * 1024 * 1024,
	}, {
		name:     "简写M单位",
		input:    "5M",
		expected: 5 * 1024 * 1024,
	}, {
		name:     "GB单位",
		input:    "6GB",
		expected: 6 * 1024 * 1024 * 1024,
	}, {
		name:     "简写G单位",
		input:    "7G",
		expected: 7 * 1024 * 1024 * 1024,
	}, {
		name:     "TB单位",
		input:    "8TB",
		expected: 8 * 1024 * 1024 * 1024 * 1024,
	}, {
		name:     "简写T单位",
		input:    "9T",
		expected: 9 * 1024 * 1024 * 1024 * 1024,
	}, {
		name:     "PB单位",
		input:    "10PB",
		expected: 10 * 1024 * 1024 * 1024 * 1024 * 1024,
	}, {
		name:     "简写P单位",
		input:    "11P",
		expected: 11 * 1024 * 1024 * 1024 * 1024 * 1024,
	}, {
		name:     "小数B",
		input:    "123.45B",
		expected: 123, // 小数部分会被截断
	}, {
		name:     "小数KB",
		input:    "1.5KB",
		expected: 1536, // 1.5 * 1024 = 1536
	}, {
		name:     "带空格",
		input:    " 100 MB ",
		expected: 100 * 1024 * 1024,
	}, {
		name:     "0值",
		input:    "0B",
		expected: 0,
	}}

	for _, tt := range successTests {
		t.Run(tt.name, func(t *testing.T) {
			var s Size
			err := s.UnmarshalJSON([]byte(tt.input))
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v, wantErr nil", err)
				return
			}
			if s != tt.expected {
				t.Errorf("UnmarshalJSON() = %v, want %v", s, tt.expected)
			}
		})
	}

	// 失败的测试用例
	failureTests := []struct {
		name        string
		input       string
		expectedErr error
	}{{
		name:        "空字符串",
		input:       "",
		expectedErr: errors.New("an empty string cannot be parsed"),
	}, {
		name:        "格式错误",
		input:       "abc",
		expectedErr: errors.New("format error, correct example: 10B, 20KB, 1.5MB, 50PB"),
	}, {
		name:        "负数",
		input:       "-10B",
		expectedErr: errors.New("file size cannot be negative"),
	}, {
		name:        "无效单位",
		input:       "100XB",
		expectedErr: errors.New("unsupported unit：XB"), // 注意这里可能需要调整错误消息
	}, {
		name:        "解析错误",
		input:       "100.B",
		expectedErr: nil, // 这个可能会失败，但错误信息是数字解析失败
	}}

	for _, tt := range failureTests {
		t.Run(tt.name, func(t *testing.T) {
			var s Size
			err := s.UnmarshalJSON([]byte(tt.input))
			if err == nil {
				t.Errorf("UnmarshalJSON() expected error, got nil")
				return
			}
			// 如果指定了预期的错误消息，则检查
			if tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestSize_Integration(t *testing.T) {
	// 集成测试：Marshal 后再 Unmarshal 应该得到原始值
	tests := []struct {
		name string
		size Size
	}{{
		name: "基本值",
		size: 1024 * 1024, // 1MB
	}, {
		name: "小值",
		size: 42,
	}, {
		name: "大值",
		size: 1024 * 1024 * 1024, // 1GB
	}, {
		name: "0值",
		size: 0,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			jsonData, err := tt.size.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			// Unmarshal
			var unmarshaledSize Size
			err = unmarshaledSize.UnmarshalJSON(jsonData)
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}

			// 检查是否相等
			// 注意：对于有小数的情况，可能会有精度损失，这里我们只测试整数情况
			if unmarshaledSize != tt.size {
				t.Errorf("Integration test failed: got %v, want %v", unmarshaledSize, tt.size)
			}
		})
	}
}

func TestSize_EdgeCases(t *testing.T) {
	// 测试边界情况
	tests := []struct {
		name string
		size Size
	}{{
		name: "接近0的值",
		size: 1,
	}, {
		name: "接近单位转换边界",
		size: 1023,
	}, {
		name: "接近单位转换边界+1",
		size: 1025,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 测试Marshal不会崩溃
			jsonData, err := tt.size.MarshalJSON()
			if err != nil {
				t.Errorf("MarshalJSON() error = %v", err)
				return
			}

			// 测试Unmarshal不会崩溃
			var s Size
			err = s.UnmarshalJSON(jsonData)
			if err != nil {
				t.Errorf("UnmarshalJSON() error = %v", err)
				return
			}
		})
	}
}
