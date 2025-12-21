# 文件存储

Lighthouse 提供统一的文件存储接口，支持 S3/MinIO 和腾讯云 COS。

## 环境变量

### S3/MinIO

```bash
STORAGE_DRIVER=s3
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_DEFAULT_BUCKET=uploads
S3_USE_SSL=false
```

### 腾讯云 COS

```bash
STORAGE_DRIVER=cos
COS_SECRET_ID=xxx
COS_SECRET_KEY=xxx
COS_BUCKET=bucket-appid
COS_REGION=ap-guangzhou
```

## 初始化存储

```go
import "github.com/light-speak/lighthouse/storages"

storage, err := storages.NewStorage()
if err != nil {
    log.Fatal(err)
}
```

## 存储接口

```go
type Storage interface {
    // 上传文件
    Put(ctx context.Context, key string, reader io.Reader) error

    // 获取预签名下载 URL
    GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

    // 获取预签名上传 URL（客户端直传）
    GetPresignedPutURL(ctx context.Context, key string, expiry time.Duration) (string, error)

    // 获取公开访问 URL
    GetPublicURL(key string) string
}
```

## 在 Resolver 中使用

### 配置 Resolver

```go
// resolver/resolver.go
type Resolver struct {
    LDB     *databases.LightDatabase
    Storage storages.Storage
}
```

### 上传文件

```go
func (r *mutationResolver) UploadAvatar(ctx context.Context, file graphql.Upload) (string, error) {
    // 生成唯一的文件路径
    key := fmt.Sprintf("avatars/%d/%s", time.Now().Unix(), file.Filename)

    // 上传到存储
    err := r.Storage.Put(ctx, key, file.File)
    if err != nil {
        return "", lighterr.NewOperationFailedError("上传失败", err)
    }

    // 返回存储路径（不是 URL）
    return key, nil
}
```

### URL 转换

在字段 Resolver 中将存储路径转换为 URL：

```go
func (r *userResolver) Avatar(ctx context.Context, obj *models.User) (*string, error) {
    if obj.Avatar == nil {
        return nil, nil
    }
    url := r.Storage.GetPublicURL(*obj.Avatar)
    return &url, nil
}
```

### 预签名上传（客户端直传）

```go
func (r *mutationResolver) GetUploadURL(ctx context.Context, filename string) (*UploadURL, error) {
    // 生成唯一的文件路径
    key := fmt.Sprintf("uploads/%d/%s", time.Now().Unix(), filename)

    // 获取预签名上传 URL
    url, err := r.Storage.GetPresignedPutURL(ctx, key, 15*time.Minute)
    if err != nil {
        return nil, err
    }

    return &UploadURL{
        URL: url,
        Key: key,
    }, nil
}
```

客户端使用预签名 URL 直接上传：

```javascript
const { url, key } = await getUploadURL(filename);

await fetch(url, {
  method: 'PUT',
  body: file,
  headers: {
    'Content-Type': file.type,
  },
});

// 上传完成后，将 key 保存到数据库
await updateAvatar(key);
```

### 预签名下载

```go
func (r *queryResolver) GetDownloadURL(ctx context.Context, key string) (string, error) {
    return r.Storage.GetPresignedURL(ctx, key, 1*time.Hour)
}
```

## MinIO 部署

### Docker

```bash
docker run -p 9000:9000 -p 9001:9001 \
  -e MINIO_ROOT_USER=minioadmin \
  -e MINIO_ROOT_PASSWORD=minioadmin \
  minio/minio server /data --console-address ":9001"
```

### Docker Compose

```yaml
services:
  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    volumes:
      - minio_data:/data

volumes:
  minio_data:
```

访问 [http://localhost:9001](http://localhost:9001) 进入 MinIO 控制台。

## 最佳实践

### 文件路径规范

```go
// 按类型和时间组织
key := fmt.Sprintf("%s/%s/%s", fileType, time.Now().Format("2006/01/02"), filename)

// 示例
// avatars/2024/01/15/user-123.jpg
// documents/2024/01/15/report.pdf
```

### 文件名安全

```go
import "path/filepath"

// 清理文件名
safeFilename := filepath.Base(filename)

// 生成唯一文件名
uniqueFilename := fmt.Sprintf("%s_%s", uuid.New().String(), safeFilename)
```

### 文件类型验证

```go
allowedTypes := map[string]bool{
    "image/jpeg": true,
    "image/png":  true,
    "image/gif":  true,
}

if !allowedTypes[file.ContentType] {
    return "", errors.New("不支持的文件类型")
}
```
