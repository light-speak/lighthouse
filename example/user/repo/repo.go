// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package repo

import (
  "github.com/light-speak/lighthouse/graphql/model"
  "github.com/light-speak/lighthouse/context"
  "user/models"
  "gorm.io/gorm"
  "github.com/light-speak/lighthouse/graphql/ast"
)

func Provide__User() map[string]*ast.Relation { return map[string]*ast.Relation{"created_at": {},"id": {},"myPosts": {Name: "post", RelationType: ast.RelationTypeHasMany, ForeignKey: "user_id", Reference: "id"},"name": {},"updated_at": {},}}
func Load__User(ctx *context.Context, key int64, field string) (map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "users", field).Load(key)
}
func LoadList__User(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "users", field).LoadList(key)
}
func Query__User(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.User{}).Scopes(scopes...)
}
func First__User(ctx *context.Context, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  var err error
  if data == nil {
    data = make(map[string]interface{})
    err = Query__User().Scopes(scopes...).First(data).Error
    if err != nil {
      return nil, err
    }
  }
  return data, nil
}
func List__User(ctx *context.Context, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
  var err error
  if datas == nil {
    datas = make([]map[string]interface{}, 0)
    err = Query__User().Scopes(scopes...).Find(&datas).Error
    if err != nil {
      return nil, err
    }
  }
  return datas, nil
}
func Count__User(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  var count int64
  err := Query__User().Scopes(scopes...).Count(&count).Error
  return count, err
}
func Provide__Post() map[string]*ast.Relation { return map[string]*ast.Relation{"BackId": {},"IsBool": {},"content": {},"created_at": {},"deleted_at": {},"enum": {},"id": {},"tagId": {},"title": {},"updated_at": {},"user": {Name: "user", RelationType: ast.RelationTypeBelongsTo, ForeignKey: "user_id", Reference: "id"},"userId": {},}}
func Load__Post(ctx *context.Context, key int64, field string) (map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "posts", field).Load(key)
}
func LoadList__Post(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "posts", field).LoadList(key)
}
func Query__Post(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.Post{}).Scopes(scopes...)
}
func First__Post(ctx *context.Context, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  var err error
  if data == nil {
    data = make(map[string]interface{})
    err = Query__Post().Scopes(scopes...).First(data).Error
    if err != nil {
      return nil, err
    }
  }
  return data, nil
}
func List__Post(ctx *context.Context, datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
  var err error
  if datas == nil {
    datas = make([]map[string]interface{}, 0)
    err = Query__Post().Scopes(scopes...).Find(&datas).Error
    if err != nil {
      return nil, err
    }
  }
  return datas, nil
}
func Count__Post(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  var count int64
  err := Query__Post().Scopes(scopes...).Count(&count).Error
  return count, err
}


func init() {
  model.AddQuickFirst("User", First__User)
  model.AddQuickList("User", List__User)
  model.AddQuickLoad("User", Load__User)
  model.AddQuickLoadList("User", LoadList__User)
  model.AddQuickCount("User", Count__User)
  model.AddQuickFirst("Post", First__Post)
  model.AddQuickList("Post", List__Post)
  model.AddQuickLoad("Post", Load__Post)
  model.AddQuickLoadList("Post", LoadList__Post)
  model.AddQuickCount("Post", Count__Post)
}
