package models


func (post *Post) Mapping() map[string]interface{} {
	return map[string]interface{}{
		"properties": map[string]interface{}{
			"title": map[string]interface{}{
				"type": "TEXT",
					"analyzer": "IK_MAX_WORD",
					"search_analyzer": "IK_SMART",
			},
		},
	}
}
