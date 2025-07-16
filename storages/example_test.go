package storages_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/light-speak/lighthouse/storages"
)

// Example_Upload 展示如何在 GraphQL resolver 中实现文件上传
func Example_Upload() {
	// 在 GraphQL mutation 中获取上传凭证
	type GetUploadURLInput struct {
		Filename string `json:"filename"`
		Type     string `json:"type"` // 文件类型：avatar, post, document
	}

	type GetUploadURLResponse struct {
		UploadURL string `json:"uploadUrl"`
		PublicURL string `json:"publicUrl"`
		ExpiresAt int64  `json:"expiresAt"`
	}

	// Mutation resolver
	getUploadURL := func(ctx context.Context, input GetUploadURLInput) (*GetUploadURLResponse, error) {
		// 根据文件类型决定存储路径
		var prefix string
		switch input.Type {
		case "avatar":
			prefix = "avatars"
		case "post":
			prefix = "posts"
		case "document":
			prefix = "documents"
		default:
			prefix = "misc"
		}

		// 生成文件路径
		fileKey := storages.GenerateFileKey(prefix, input.Filename)

		// 获取上传凭证
		uploadURL, publicURL, err := storages.GetUploadURL(ctx, storages.UploadConfig{
			Key:    fileKey,
			Expiry: 15 * time.Minute,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get upload URL: %w", err)
		}

		return &GetUploadURLResponse{
			UploadURL: uploadURL,
			PublicURL: publicURL,
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		}, nil
	}

	// 使用示例
	ctx := context.Background()
	resp, err := getUploadURL(ctx, GetUploadURLInput{
		Filename: "profile.jpg",
		Type:     "avatar",
	})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Upload URL: %s\n", resp.UploadURL)
	fmt.Printf("Public URL: %s\n", resp.PublicURL)
	fmt.Printf("Expires at: %d\n", resp.ExpiresAt)
}

// Example_DirectUpload 展示如何直接通过后端上传文件
func Example_DirectUpload() {
	ctx := context.Background()

	// 获取存储实例
	storage, err := storages.GetStorage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 假设有一个文件 reader
	var fileReader *strings.Reader // 实际使用时替换为真实的文件 reader

	// 生成文件路径
	fileKey := storages.GenerateFileKey("uploads", "document.pdf")

	// 上传文件
	err = storage.Put(ctx, storages.GetDefaultBucket(), fileKey, fileReader)
	if err != nil {
		fmt.Printf("Upload error: %v\n", err)
		return
	}

	// 获取公开访问 URL
	publicURL := storage.GetPublicURL(storages.GetDefaultBucket(), fileKey)
	fmt.Printf("File uploaded successfully: %s\n", publicURL)
}
