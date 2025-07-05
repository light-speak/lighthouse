package utils

import (
	"math/rand"
	"time"
)

// 生成指定长度的随机数字码
func RandomCode(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// 生成指定长度的随机字符串,包含大小写字母和数字
func RandomString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// 生成指定范围内的随机整数
func RandomInt(min, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max-min+1) + min
}

// 生成指定范围内的随机浮点数
func RandomFloat(min, max float64) float64 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + r.Float64()*(max-min)
}

// 生成随机布尔值
func RandomBool() bool {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(2) == 1
}

// 生成指定长度的随机小写字母字符串
func RandomLowerString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// 生成指定长度的随机大写字母字符串
func RandomUpperString(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}

// 从切片中随机选择一个元素
func RandomChoice[T any](slice []T) T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return slice[r.Intn(len(slice))]
}

// 生成指定范围内的随机时间
func RandomTime(start, end time.Time) time.Time {
	min := start.Unix()
	max := end.Unix()
	delta := max - min

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sec := r.Int63n(delta) + min
	return time.Unix(sec, 0)
}
