// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package repo

import (
  "github.com/light-speak/lighthouse/graphql/ast"
  "github.com/light-speak/lighthouse/graphql/model"
  "gorm.io/gorm"
  "user/models"
)

func Provide__User() map[string]*ast.Relation { return map[string]*ast.Relation{"created_at": {},"id": {},"name": {},"posts": {},"updated_at": {},}}
func Query__User(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.User{}).Scopes(scopes...)
}
func Fields__User(user *models.User, key string) (interface{}, error) {
  switch key {
    case "created_at": 
      return user.CreatedAt, nil
    case "id": 
      return user.Id, nil
    case "name": 
      return user.Name, nil
    case "posts": 
      return user.Posts, nil
    case "updated_at": 
      return user.UpdatedAt, nil
  }
  return nil, nil
} 
func First__User(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  selectColumns, selectRelations := model.GetSelectInfo(columns, Provide__User())
  user := &models.User{}
  err := Query__User().Scopes(scopes...).Select(selectColumns).First(user).Error
  if err != nil {
    return nil, err
  }
  res, err := model.StructToMap(user)
  if err != nil {
    return nil, err
  }
  for _, relation := range selectRelations {
    fieldValue, err := Fields__User(user, relation.Relation.ForeignKey)
    if err != nil {
      return nil, err
    }
    res, err = model.FetchRelation(res, relation, fieldValue)
    if err != nil {
      return nil, err
    }
  }
  return res, nil
}
func Provide__Post() map[string]*ast.Relation { return map[string]*ast.Relation{"content": {},"created_at": {},"deleted_at": {},"id": {},"title": {},"updated_at": {},"user": {Name: "user", RelationType: ast.RelationTypeBelongsTo, ForeignKey: "user_id", Reference: "id"},"user_id": {},}}
func Query__Post(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.Post{}).Scopes(scopes...)
}
func Fields__Post(post *models.Post, key string) (interface{}, error) {
  switch key {
    case "content": 
      return post.Content, nil
    case "created_at": 
      return post.CreatedAt, nil
    case "deleted_at": 
      return post.DeletedAt, nil
    case "id": 
      return post.Id, nil
    case "title": 
      return post.Title, nil
    case "updated_at": 
      return post.UpdatedAt, nil
    case "user": 
      return post.User, nil
    case "user_id": 
      return post.UserId, nil
  }
  return nil, nil
} 
func First__Post(columns map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  selectColumns, selectRelations := model.GetSelectInfo(columns, Provide__Post())
  post := &models.Post{}
  err := Query__Post().Scopes(scopes...).Select(selectColumns).First(post).Error
  if err != nil {
    return nil, err
  }
  res, err := model.StructToMap(post)
  if err != nil {
    return nil, err
  }
  for _, relation := range selectRelations {
    fieldValue, err := Fields__Post(post, relation.Relation.ForeignKey)
    if err != nil {
      return nil, err
    }
    res, err = model.FetchRelation(res, relation, fieldValue)
    if err != nil {
      return nil, err
    }
  }
  return res, nil
}


func init() {
  model.AddQuickFirst("User", First__User)
  model.AddQuickFirst("Post", First__Post)
}
