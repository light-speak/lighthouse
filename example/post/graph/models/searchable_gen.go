package models

import (
	"encoding/json"
	"github.com/light-speak/lighthouse/log"
)


// Mapping 方法将生成映射信息的 JSON 字符串
func (post *Post) Mapping() string {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart",
				},
				"content": map[string]interface{}{
					"type": "text",
					"analyzer": "ik_max_word",
					"search_analyzer": "ik_smart",
				},
			},
		},
	}

	// 将映射数据转换为 JSON 字符串
	jsonBytes, err := json.Marshal(mapping)
	if err != nil {
		log.Error("Error converting to JSON: %v", err)
		return "{}" 
	}
	return string(jsonBytes)
}

// SearchId 方法返回模型的 ID
func (post *Post) SearchId() int64 {
	return post.GetID()
}

// IndexName 方法返回索引名称
func (post *Post) IndexName() string {
	return "posts"
}

// GetIndexData 方法将返回模型数据的 JSON 字符串
func (post *Post) GetIndexData() string {
	data := map[string]interface{}{
		"title": post.Title,
		"content": post.Content,
	}

	// 将数据转换为 JSON 字符串
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Error("Error converting to JSON: %v", err)
		return "{}" 
	}
	return string(jsonBytes)
}