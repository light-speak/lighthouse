package resolver

import (
	"context"
	"github.com/light-speak/lighthouse/middleware"
	"gorm.io/gorm"
)

func ResolverData(ctx context.Context, db *gorm.DB, data interface{}, isList bool) error {
	lCtx := middleware.GetContext(ctx)
	tx := db
	for _, where := range lCtx.Wheres {
		tx = tx.Where(where.Query, where.Value)
		if err := tx.Error; err != nil {
			return err
		}
	}
	if isList {
		if lCtx.Paginate != nil {
			tx = tx.Limit(lCtx.Paginate.Size).Offset((lCtx.Paginate.Page - 1) * lCtx.Paginate.Size)
		}
		if err := tx.Find(&data).Error; err != nil {
			return err
		}
	} else {
		if err := tx.First(&data).Error; err != nil {
			return err
		}
	}
	return nil
}
