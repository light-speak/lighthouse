package scalars

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

func MarshalDeletedAt(deletedAt *gorm.DeletedAt) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		json.NewEncoder(w).Encode(*deletedAt)
	})
}

func UnmarshalDeletedAt(v interface{}) (*gorm.DeletedAt, error) {
	switch v := v.(type) {
	case string:
		var deletedAt gorm.DeletedAt
		return &deletedAt, sonic.Unmarshal([]byte(v), &deletedAt)
	case []byte:
		var deletedAt gorm.DeletedAt
		return &deletedAt, sonic.Unmarshal(v, &deletedAt)
	case int:
		deletedAt := gorm.DeletedAt{Time: time.Unix(int64(v), 0)}
		return &deletedAt, nil
	case int64:
		deletedAt := gorm.DeletedAt{Time: time.Unix(v, 0)}
		return &deletedAt, nil
	case time.Time:
		deletedAt := gorm.DeletedAt{Time: v}
		return &deletedAt, nil
	default:
		return nil, fmt.Errorf("invalid type for gorm.DeletedAt: %T", v)
	}
}
