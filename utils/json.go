package utils

import "encoding/json"

func ParseJSONNumberToI64(num json.Number) (int, error) {
	val, err := num.Int64()
	if err != nil {
		return 0, err
	}
	result := int(val)
	return result, nil
}
