package utilk

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Size int

// 定义单位序列（默认小写+完整单位）
var units = []string{"B", "KB", "MB", "GB", "TB", "PB"}

// 定义单位对应的 1024 次方（B=1024^0, KB=1024^1, ..., PB=1024^5）
var unitMap = map[string]int{
	"B":  0, // Bytes
	"KB": 1,
	"MB": 2,
	"GB": 3,
	"TB": 4,
	"PB": 5,
}

func (s *Size) MarshalJSON() ([]byte, error) {
	// 处理 0 字节
	if *s == 0 {
		return []byte("0B"), nil
	}
	// 计算单位索引（log2(1024) ≈ 10，每 10 位代表一个单位）
	base := float64(1024)
	exp := int(math.Log(float64(*s)) / math.Log(base))
	if exp >= len(units) {
		exp = len(units) - 1 // 超出最大单位（PB）则用 PB 表示
	}
	// 计算格式化后的值（保留指定小数位数）
	value := float64(*s) / math.Pow(base, float64(exp))
	// 格式化并返回结果（去除末尾多余的 .0）
	result := fmt.Sprintf("%.2f%s", value, units[exp])
	return []byte(result), nil
}

// 正则表达式：匹配 数字（整数/小数）+ 单位（可选，默认B）
// 支持格式：10B、20KB、1.5MB、1024gb、50PB 等
var sizeRegex = regexp.MustCompile(`^(\d+(\.\d+)?)\s*([BbKkMmGgTtPp]?[Bb]?)$`)

func (s *Size) UnmarshalJSON(b []byte) error {
	// 去除字符串前后空格
	sizeStr := strings.TrimSpace(string(b))
	if sizeStr == "" {
		return errors.New("an empty string cannot be parsed")
	}
	matches := sizeRegex.FindStringSubmatch(sizeStr)
	if len(matches) == 0 {
		return errors.New("format error, correct example: 10B, 20KB, 1.5MB, 50PB")
	}
	// 提取数字部分和单位部分
	numStr := matches[1]
	unitStr := strings.ToUpper(matches[3])
	// 解析数字（支持整数和小数）
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return errors.New("number parsing failed: " + err.Error())
	}
	if num < 0 {
		return errors.New("file size cannot be negative")
	}
	// 处理单位（兼容简写，如 K→KB、M→MB，无单位默认B）
	unitStr = strings.ToUpper(unitStr)
	switch unitStr {
	case "":
		unitStr = "B" // 无单位默认是字节（如 "100" → 100B）
	case "K":
		unitStr = "KB"
	case "M":
		unitStr = "MB"
	case "G":
		unitStr = "GB"
	case "T":
		unitStr = "TB"
	case "P":
		unitStr = "PB"
	}
	// 检查单位是否合法
	exp, ok := unitMap[unitStr]
	if !ok {
		return fmt.Errorf("unsupported unit：%s", matches[3])
	}
	// 计算字节数：num * (1024^exp)
	// 1024^5 = 1125899906842624（PB级别），乘以 100 后仍在 int64 范围内（int64最大值：9223372036854775807）
	multiplier := math.Pow(1024, float64(exp))
	bytes := num * multiplier
	// 检查是否超出 int64 范围（避免溢出）
	if bytes > math.MaxInt64 || bytes < math.MinInt64 {
		return fmt.Errorf("converted byte size exceeds int64 range (max %.0f PB)", math.MaxInt64/math.Pow(1024, 5))
	}
	// 转换为整数（小数部分直接舍弃，如 1.9KB → 1*1024=1024B）
	*s = Size(int64(bytes))
	return nil
}

func ParseSize(fmt string) (Size, error) {
	var size Size
	err := (&size).UnmarshalJSON([]byte(fmt))
	if err != nil {
		return 0, err
	}
	return size, nil
}
