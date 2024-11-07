// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import (
  "gorm.io/gorm"
  "github.com/light-speak/lighthouse/context"
  "github.com/light-speak/lighthouse/graphql/model"
)


func PostUserId1(ctx *context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
    // Func:UserId1 user code start. Do not remove this comment. 
		return db.Where("user_id = ?", 1)
    // Func:UserId1 user code end. Do not remove this comment. 
	}
}
func PostUserId2(ctx *context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
    // Func:UserId2 user code start. Do not remove this comment. 
		return db.Where("user_id = ?", 2)
    // Func:UserId2 user code end. Do not remove this comment. 
	}
}


func init() {
	model.AddScopes("PostUserId1", PostUserId1)
	model.AddScopes("PostUserId2", PostUserId2)
}