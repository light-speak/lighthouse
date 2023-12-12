package resolver

import (
	"context"
	"errors"
	"github.com/light-speak/lighthouse"
	"github.com/light-speak/lighthouse/middleware"
	"gorm.io/gorm"
	"time"
)

type QueryType int

var (
	ListQuery  QueryType = 1
	OneQuery   QueryType = 2
	CountQuery QueryType = 3
	SumQuery   QueryType = 4
)

type Type int

var (
	Query    Type = 1
	Mutation Type = 2
)

type MutationType int

var (
	CreateMutation MutationType = 1
	UpdateMutation MutationType = 2
)

type Option struct {
	*Type
	*QueryType
	*MutationType
}

func ResolveData(ctx context.Context, db *gorm.DB, data interface{}, option Option) error {
	lCtx := middleware.GetContext(ctx)
	tx := db
	if option.Type == nil {
		option.Type = &Query
	}
	if option.QueryType == nil {
		option.QueryType = &ListQuery
	}
	if option.MutationType == nil {
		option.MutationType = &CreateMutation
	}

	if option.Type == &Query {
		for _, where := range lCtx.Wheres {
			tx = tx.Where(where.Query, where.Value)
			if err := tx.Error; err != nil {
				return err
			}
		}
		switch option.QueryType {
		case &ListQuery:
			return resolveList(lCtx, tx, data)
		case &OneQuery:
			return resolveOne(lCtx, tx, data)
		default:
			return errors.New("error Query Type")
		}
	} else if option.Type == &Mutation {
		switch option.MutationType {
		case &CreateMutation:
			return resolveCreate(lCtx, tx, data)
		case &UpdateMutation:
			return resolveUpdate(lCtx, tx, data)
		default:
			return errors.New("error Mutation Type")
		}
	} else {
		return errors.New("error Type")
	}
}

func resolveList(lCtx *lighthouse.Context, tx *gorm.DB, data interface{}) error {
	if lCtx.Paginate != nil {
		tx = tx.Limit(lCtx.Paginate.Size).Offset((lCtx.Paginate.Page - 1) * lCtx.Paginate.Size)
	}
	if err := tx.Find(&data).Error; err != nil {
		return err
	}
	return nil
}

func resolveOne(lCtx *lighthouse.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.First(&data).Error; err != nil {
		return err
	}
	return nil
}

func resolveCreate(lCtx *lighthouse.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.Create(data).Error; err != nil {
		return err

	}
	return nil
}

func resolveUpdate(lCtx *lighthouse.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.Model(data).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func setCreatedAndUpdatedTime(isCreate bool, data *map[string]interface{}) {
	if isCreate {
		(*data)["CreatedAt"] = time.Now()
	}
	(*data)["UpdatedAt"] = time.Now()
}
