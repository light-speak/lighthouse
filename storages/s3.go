package storages

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Storage S3 存储实现
// 基于 minio-go v7 客户端
type S3Storage struct {
	client   *minio.Client
	endpoint string // 端点地址
	useSSL   bool   // 是否使用 SSL
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
		client:   client,
		endpoint: config.Endpoint,
		useSSL:   config.UseSSL,
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

// GetPresignedPutURL 获取预签名上传 URL
// 允许前端直接上传文件到 S3，无需经过后端服务器
func (s *S3Storage) GetPresignedPutURL(ctx context.Context, bucket string, key string, expiry time.Duration) (string, error) {
	// 生成预签名上传 URL
	url, err := s.client.PresignedPutObject(ctx, bucket, key, expiry)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}

// GetPublicURL 获取公开访问 URL
// 用于存储桶配置为公有读的场景
func (s *S3Storage) GetPublicURL(bucket string, key string) string {
	// 构建公开访问 URL
	if s.useSSL {
		return fmt.Sprintf("https://%s/%s/%s", s.endpoint, bucket, key)
	}
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, bucket, key)
}
