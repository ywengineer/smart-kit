package validator

import (
	"errors"
	"net"
	"strings"
)

// IPType 定义 IP 类型枚举
type IPType int

const (
	IPTypeUnknown   IPType = iota // 未知类型
	IPTypeIPv4                    // 纯 IPv4 地址
	IPTypeIPv6                    // 纯 IPv6 地址
	IPTypeIPv4CIDR                // IPv4 CIDR 段（如 192.168.1.0/24）
	IPTypeIPv6CIDR                // IPv6 CIDR 段（如 2001:0db8::/32）
	IPTypeIPv4Range               // IPv4 范围段（如 192.168.1.1-192.168.1.255）
)

var ErrEmptyString = errors.New("empty string")
var ErrMaskIpv4 = errors.New("IPv4 CIDR 掩码范围非法（需 0-32）")
var ErrMaskIpv6 = errors.New("IPv6 CIDR 掩码范围非法（需 0-128）")
var ErrIllegalInput = errors.New("illegal input")
var ErrIpv4Range = errors.New("IPv4 范围格式非法")

// IsValidIPOrIPRange 校验输入是否为合法 IP 地址/IP 段
// 返回：IP 类型、无错误则返回 nil
func IsValidIPOrIPRange(input string) (IPType, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return IPTypeUnknown, ErrEmptyString
	}

	// 1. 先校验是否为纯 IP 地址（IPv4/IPv6）
	if ip := net.ParseIP(input); ip != nil {
		if ip.To4() != nil {
			return IPTypeIPv4, nil // 纯 IPv4 地址
		}
		return IPTypeIPv6, nil // 纯 IPv6 地址
	}

	// 2. 校验是否为 CIDR 格式（IPv4/IPv6）
	if _, ipNet, err := net.ParseCIDR(input); err == nil {
		if ipNet.IP.To4() != nil {
			// 额外校验 IPv4 CIDR 的掩码是否合法（0-32）
			maskSize, _ := ipNet.Mask.Size()
			if maskSize < 0 || maskSize > 32 {
				return IPTypeUnknown, ErrMaskIpv4
			}
			return IPTypeIPv4CIDR, nil
		}
		// 额外校验 IPv6 CIDR 的掩码是否合法（0-128）
		maskSize, _ := ipNet.Mask.Size()
		if maskSize < 0 || maskSize > 128 {
			return IPTypeUnknown, ErrMaskIpv6
		}
		return IPTypeIPv6CIDR, nil
	}

	// 3. 校验是否为 IPv4 范围格式（如 192.168.1.1-192.168.1.255）
	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return IPTypeUnknown, ErrIpv4Range //errors.New("IP 范围格式非法（需仅包含一个 '-' 分隔符）")
		}
		startIP := net.ParseIP(strings.TrimSpace(parts[0]))
		endIP := net.ParseIP(strings.TrimSpace(parts[1]))
		if startIP == nil || endIP == nil {
			return IPTypeUnknown, ErrIpv4Range //errors.New("IP 范围中包含非法 IP 地址")
		}
		if startIP.To4() == nil || endIP.To4() == nil {
			return IPTypeUnknown, ErrIpv4Range // errors.New("仅支持 IPv4 地址范围格式")
		}
		// 校验起始 IP ≤ 结束 IP（转换为 uint32 比较）
		startInt := IpToUint32(startIP.To4())
		endInt := IpToUint32(endIP.To4())
		if startInt > endInt {
			return IPTypeUnknown, ErrIpv4Range //errors.New("IP 范围非法（起始 IP 大于结束 IP）")
		}
		return IPTypeIPv4Range, nil
	}

	// 4. 所有格式均不匹配
	return IPTypeUnknown, ErrIllegalInput //errors.New("输入内容非合法 IP 地址/IP 段")
}

// IpToUint32 ipToUint32 将 IPv4 地址转换为 uint32（用于比较大小）
func IpToUint32(ip net.IP) uint32 {
	if len(ip) != 4 {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// 辅助函数：获取 IP 类型的字符串描述
func (t IPType) String() string {
	switch t {
	case IPTypeIPv4:
		return "纯 IPv4 地址"
	case IPTypeIPv6:
		return "纯 IPv6 地址"
	case IPTypeIPv4CIDR:
		return "IPv4 CIDR 段"
	case IPTypeIPv6CIDR:
		return "IPv6 CIDR 段"
	case IPTypeIPv4Range:
		return "IPv4 范围段"
	default:
		return "未知类型"
	}
}
