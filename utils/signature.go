package utils

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"time"
)

func GenerateSignature(length int, prefix string) (string, error) {
	// 数字字符集
	const digits = "0123456789"
	randomPart := make([]byte, length-14) // 减去时间戳长度

	for i := 0; i < len(randomPart); i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		randomPart[i] = digits[n.Int64()]
	}

	// 获取当前时间戳（精确到微秒）
	timestamp := time.Now().UnixMicro()
	timestampStr := strconv.FormatInt(timestamp, 10)

	// 拼接时间戳和随机部分
	return prefix + timestampStr + string(randomPart), nil
}
