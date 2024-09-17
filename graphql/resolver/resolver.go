package resolver

import (
	"context"
	"errors"
	"fmt"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/graphql/middleware"
	"gorm.io/gorm"
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

func ResolveData(ctx context.Context, db *gorm.DB, path string, data interface{}, option Option) error {
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
	if *option.Type == Query {
		for _, where := range lCtx.Wheres {
			if path == where.Path {
				tx = tx.Where(where.Query, where.Value)
				if err := tx.Error; err != nil {
					return err
				}
			}
		}
		switch option.QueryType {
		case &ListQuery:
			return resolveList(lCtx, tx, data)
		case &OneQuery:
			return resolveOne(lCtx, tx, data)
		case &CountQuery:
			return resolveCount(lCtx, tx, data)
		case &SumQuery:
			return resolveSum(lCtx, tx, data)
		default:
			return errors.New("error Query Type")
		}
	} else if *option.Type == Mutation {
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

func resolveList(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if lCtx.Paginate != nil {
		tx = tx.Limit(int(lCtx.Paginate.Size)).Offset((int(lCtx.Paginate.Page - 1)) * int(lCtx.Paginate.Size))
	}

	if err := tx.Find(data).Error; err != nil {
		return err
	}
	return nil
}

func resolveOne(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.First(data).Error; err != nil {
		return err
	}
	return nil
}

func resolveCreate(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.Create(data).Error; err != nil {
		return err

	}
	return nil
}

func resolveUpdate(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if err := tx.Model(data).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func resolveCount(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if intPtr, ok := data.(*int); ok {
		var count int64
		if err := tx.Count(&count).Error; err != nil {
			return err
		}
		*intPtr = int(count)
		return nil
	}

	return nil
}

func resolveSum(lCtx *graphql.Context, tx *gorm.DB, data interface{}) error {
	if intPtr, ok := data.(*int); ok {
		var sum int64
		if err := tx.Select(fmt.Sprintf("SUM(%s)", *lCtx.Column)).Scan(&sum).Error; err != nil {
			return err
		}
		*intPtr = int(sum)
		return nil
	}

	return nil
}
