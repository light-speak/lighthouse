package storages

import (
	"context"
	"errors"
	"io"
	"time"
)

// COSStorage 腾讯云 COS 存储实现
// 当前为空实现，待后续开发
type COSStorage struct {
	// TODO: 添加 COS 客户端配置
}

// COSConfig COS 配置参数
type COSConfig struct {
	SecretID  string // 密钥 ID
	SecretKey string // 密钥
	Region    string // 地域
	// TODO: 添加更多配置参数
}

// NewCOSStorage 创建新的 COS 存储实例
// config: COS 配置参数
func NewCOSStorage(config COSConfig) (*COSStorage, error) {
	// TODO: 实现 COS 客户端初始化
	return &COSStorage{}, nil
}

// Put 上传文件到 COS 存储桶 (空实现)
func (c *COSStorage) Put(ctx context.Context, key string, reader io.Reader) error {
	// TODO: 实现 COS 文件上传
	return errors.New("COS storage not implemented yet")
}

// GetPresignedURL 获取预签名下载 URL (空实现)
func (c *COSStorage) GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// TODO: 实现 COS 预签名 URL 生成
	return "", errors.New("COS storage not implemented yet")
}

// GetPresignedPutURL 获取预签名上传 URL (空实现)
func (c *COSStorage) GetPresignedPutURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// TODO: 实现 COS 预签名上传 URL 生成
	return "", errors.New("COS storage not implemented yet")
}

// GetPublicURL 获取公开访问 URL (空实现)
func (c *COSStorage) GetPublicURL(key string) string {
	// TODO: 实现 COS 公开 URL 生成
	// 格式: https://{bucket}.cos.{region}.myqcloud.com/{key}
	return ""
}
