package storages

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage S3 存储实现
// 基于 minio-go v7 客户端
type S3Storage struct {
	client *minio.Client
}

// S3Config S3 配置参数
type S3Config struct {
	Endpoint        string // 端点地址
	AccessKeyID     string // 访问密钥 ID
	SecretAccessKey string // 访问密钥
	UseSSL          bool   // 是否使用 SSL
}

// NewS3Storage 创建新的 S3 存储实例
// config: S3 配置参数
func NewS3Storage(config S3Config) (*S3Storage, error) {
	// 创建 minio 客户端
	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKeyID, config.SecretAccessKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	return &S3Storage{
		client: client,
	}, nil
}

// Put 上传文件到 S3 存储桶
func (s *S3Storage) Put(ctx context.Context, bucket string, key string, reader io.Reader) error {
	// 使用 PutObject 上传文件
	_, err := s.client.PutObject(ctx, bucket, key, reader, -1, minio.PutObjectOptions{})
	return err
}

// GetPresignedURL 获取预签名下载 URL
func (s *S3Storage) GetPresignedURL(ctx context.Context, bucket string, key string, expiry time.Duration) (string, error) {
	// 生成预签名 URL
	url, err := s.client.PresignedGetObject(ctx, bucket, key, expiry, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
