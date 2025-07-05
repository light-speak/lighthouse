package utils

import (
	"fmt"
	"regexp"
	"time"
)

// IsValidPhoneNumber 校验中国手机号
// 中国手机号的正则表达式：以1开头，第二位是3-9，后面跟9个数字，总共11位
func IsValidPhoneNumber(phone string) bool {
	var re = regexp.MustCompile(`^1[3-9]\d{9}$`)
	return re.MatchString(phone)
}

// IsValidIP 校验IPv4地址
func IsValidIP(ip string) bool {
	var re = regexp.MustCompile(`^((25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)\.){3}(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)$`)
	return re.MatchString(ip)
}

// IsValidIPv6 校验IPv6地址
func IsValidIPv6(ip string) bool {
	//TODO: 实现 IPv6 地址的校验
	return true
}

// IsValidIDCard 验证身份证号码是否有效
// 支持15位老身份证、18位新身份证和外国人永久居留证
// 15位身份证格式: PPPPPPYYMMDDXXX
// 18位身份证格式: PPPPPPYYYYMMDDXXXC
// 外国人永久居留证格式: 前3位为字母，后13位为数字
// P: 省份和城市代码, Y: 年份, M: 月份, D: 日期, X: 顺序码, C: 校验码
func IsValidIDCard(id string) bool {
	// 检查外国人永久居留证
	if len(id) == 16 {
		var re = regexp.MustCompile(`^[A-Z]{3}\d{13}$`)
		return re.MatchString(id)
	}

	// 检查15位老身份证
	if len(id) == 15 {
		var re = regexp.MustCompile(`^[1-9]\d{5}\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}$`)
		if !re.MatchString(id) {
			return false
		}
		// 校验日期合法性
		year := "19" + id[6:8]
		month := id[8:10]
		day := id[10:12]
		return IsValidDate(fmt.Sprintf("%s-%s-%s", year, month, day))
	}

	// 检查18位新身份证
	if len(id) == 18 {
		var re = regexp.MustCompile(`^[1-9]\d{5}(19|20)\d{2}(0[1-9]|1[0-2])(0[1-9]|[12]\d|3[01])\d{3}[\dX]$`)
		if !re.MatchString(id) {
			return false
		}

		if id[17] != 'X' && (id[17] < '0' || id[17] > '9') {
			return false
		}
		// 校验日期合法性
		year := id[6:10]
		month := id[10:12]
		day := id[12:14]
		if !IsValidDate(fmt.Sprintf("%s-%s-%s", year, month, day)) {
			return false
		}

		weights := []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		checkCode := []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
		sum := 0
		for i := 0; i < 17; i++ {
			sum += int(id[i]-'0') * weights[i]
		}
		return id[17] == checkCode[sum%11]
	}

	return false
}

// IsValidBankCard checks if Chinese bank card number is valid (16-19 digits)
func IsValidBankCard(card string) bool {
	var re = regexp.MustCompile(`^[1-9]\d{15,18}$`)
	if !re.MatchString(card) {
		return false
	}

	// Luhn 算法校验验证最后一位校验码
	sum := 0
	alt := false
	for i := len(card) - 1; i >= 0; i-- {
		n := int(card[i] - '0')
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	return sum%10 == 0
}

// IsValidPostCode 邮编校验 仅针对中国大陆邮编
func IsValidPostCode(code string) bool {
	var re = regexp.MustCompile(`^[1-9]\d{5}$`)
	return re.MatchString(code)
}

// IsValidDate 日期校验 格式: YYYY-MM-DD
func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// IsValidTime 时间校验 格式: HH:MM:SS
func IsValidTime(timeStr string) bool {
	_, err := time.Parse("15:04:05", timeStr)
	return err == nil
}

// IsValidDateTime 日期时间校验 格式: YYYY-MM-DD HH:MM:SS
func IsValidDateTime(datetime string) bool {
	_, err := time.Parse("2006-01-02 15:04:05", datetime)
	return err == nil
}

// IsValidAmount 金额校验 格式: 123.45
func IsValidAmount(amount string) bool {
	var re = regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	return re.MatchString(amount)
}
