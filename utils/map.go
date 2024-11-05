package utils

import "sync"

func MapToSyncMap(data map[string]interface{}) *sync.Map {
	syncMap := &sync.Map{}
	for k, v := range data {
		syncMap.Store(k, v)
	}
	return syncMap
}

func MapSliceToSyncMapSlice(datas []map[string]interface{}) []*sync.Map {
	result := make([]*sync.Map, len(datas))
	for i, data := range datas {
		result[i] = MapToSyncMap(data)
	}
	return result
}
