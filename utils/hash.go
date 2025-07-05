package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt 对字符串加密
func HashAndSalt(value string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func IsHashAndSalt(value string, plainLength int) bool {
	// 检查长度和前缀是否符合 bcrypt 格式
	if len(value) != 60 || !strings.HasPrefix(value, "$2") {
		return false
	}

	// 模拟生成符合长度的明文并进行 bcrypt 验证
	dummyPlain := strings.Repeat("x", plainLength) // 生成指定长度的字符串
	err := bcrypt.CompareHashAndPassword([]byte(value), []byte(dummyPlain))
	return err == bcrypt.ErrMismatchedHashAndPassword // 验证哈希格式有效
}

// ComparePasswords 比较哈希密码和输入的密码
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // 如果密码匹配，返回 true
}

// MD5Hash 使用MD5算法对字符串进行哈希
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA1Hash 使用SHA1算法对字符串进行哈希
func SHA1Hash(text string) string {
	hash := sha1.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA256Hash 使用SHA256算法对字符串进行哈希
func SHA256Hash(text string) string {
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}

// SHA512Hash 使用SHA512算法对字符串进行哈希
func SHA512Hash(text string) string {
	hash := sha512.Sum512([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GenerateSalt 生成指定长度的随机盐值
func GenerateSalt(length int) string {
	return RandomString(length)
}

// HashWithSalt 将密码与盐值组合后进行哈希
func HashWithSalt(password string, salt string) string {
	combined := password + salt
	return SHA256Hash(combined)
}

// Base64Encode 将字符串进行Base64编码
func Base64Encode(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

// Base64Decode 将Base64编码的字符串解码
func Base64Decode(encodedText string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// VerifyHash 验证字符串与其哈希值是否匹配（支持多种哈希算法）
func VerifyHash(text string, hash string, algorithm string) bool {
	var computedHash string
	switch algorithm {
	case "md5":
		computedHash = MD5Hash(text)
	case "sha1":
		computedHash = SHA1Hash(text)
	case "sha256":
		computedHash = SHA256Hash(text)
	case "sha512":
		computedHash = SHA512Hash(text)
	default:
		return false
	}
	return computedHash == hash
}

// HashWithPepper 将密码与pepper(固定盐值)组合后进行哈希
func HashWithPepper(password string, pepper string) string {
	combined := password + pepper
	return SHA256Hash(combined)
}

// DoubleHash 对字符串进行双重哈希，增加安全性
func DoubleHash(text string) string {
	firstHash := SHA256Hash(text)
	return SHA256Hash(firstHash)
}
