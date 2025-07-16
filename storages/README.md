# Storages 模块

本模块提供统一的存储接口，支持多种存储后端。

## 支持的存储类型

- **S3**: 兼容 S3 协议的存储服务（包括 AWS S3、MinIO、阿里云 OSS 等）
- **COS**: 腾讯云对象存储（预留实现）

## 环境变量配置

在项目根目录的 `.env` 文件中添加以下配置：

### 通用配置

```env
# 存储驱动类型：s3 或 cos
STORAGE_DRIVER=s3
```

### S3 存储配置

```env
# S3 端点地址
S3_ENDPOINT=localhost:9000
# 访问密钥 ID
S3_ACCESS_KEY=your-access-key
# 访问密钥
S3_SECRET_KEY=your-secret-key
# 是否使用 SSL
S3_USE_SSL=false
# 默认存储桶
S3_DEFAULT_BUCKET=default
```

### COS 存储配置

```env
# 密钥 ID
COS_SECRET_ID=your-secret-id
# 密钥
COS_SECRET_KEY=your-secret-key
# 地域
COS_REGION=ap-beijing
# 默认存储桶
COS_DEFAULT_BUCKET=default
```

## 使用方法

### 获取存储实例

```go
import "github.com/light-speak/lighthouse/storages"

// 获取存储实例
storage, err := storages.GetStorage()
if err != nil {
    log.Fatal(err)
}
```

### 上传文件

```go
// 使用默认存储桶
bucket := storages.GetDefaultBucket()

// 上传文件
err := storage.Put(ctx, bucket, "path/to/file.jpg", fileReader)
if err != nil {
    log.Printf("上传失败: %v", err)
}
```

### 获取预签名 URL

```go
// 生成有效期为 1 小时的预签名 URL
url, err := storage.GetPresignedURL(ctx, bucket, "path/to/file.jpg", time.Hour)
if err != nil {
    log.Printf("生成预签名 URL 失败: %v", err)
}
```

## 扩展存储类型

要添加新的存储类型，需要：

1. 在 `env.go` 中添加新的驱动常量
2. 创建新的存储实现文件（如 `oss.go`）
3. 实现 `Storage` 接口
4. 在 `env.go` 的 `initStorage()` 函数中添加初始化逻辑

## 注意事项

- 如果存储配置不完整，系统会跳过初始化并记录警告日志
- 使用前请确保相关的存储服务已经正确配置和启动
- 对于生产环境，建议使用 SSL 连接 