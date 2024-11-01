// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package repo

import (
  "github.com/light-speak/lighthouse/graphql/model"
  "github.com/light-speak/lighthouse/context"
  "test/models"
  "github.com/light-speak/lighthouse/graphql/ast"
  "sync"
  "gorm.io/gorm"
)

func Provide__Test() map[string]*ast.Relation { return map[string]*ast.Relation{"created_at": {},"email": {},"id": {},"name": {},"updated_at": {},}}
func Load__Test(ctx *context.Context, key int64, field string) (map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "tests", field).Load(key)
}
func LoadList__Test(ctx *context.Context, key int64, field string) ([]map[string]interface{}, error) {
  return model.GetLoader[int64](model.GetDB(), "tests", field).LoadList(key)
}
func Query__Test(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(&models.Test{}).Scopes(scopes...)
}
func First__Test(ctx *context.Context, columns map[string]interface{}, data map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (map[string]interface{}, error) {
  var err error
  selectColumns, selectRelations := model.GetSelectInfo(columns, Provide__Test())
  if data == nil {
    data = make(map[string]interface{})
    err = Query__Test().Scopes(scopes...).Select(selectColumns).First(data).Error
    if err != nil {
      return nil, err
    }
  }
  var wg sync.WaitGroup
  errChan := make(chan error, len(selectRelations))
  var mu sync.Mutex
  
  for key, relation := range selectRelations {
    wg.Add(1)
    go func(data map[string]interface{}, relation *model.SelectRelation)  {
      defer wg.Done()
      cData, err := model.FetchRelation(ctx, data, relation)
      if err != nil {
        errChan <- err
      }
      mu.Lock()
      defer mu.Unlock()
      data[key] = cData
    }(data, relation) 
  }
  wg.Wait()
  close(errChan)
  for err := range errChan {
    return nil, err
  }
  return data, nil
}
func List__Test(ctx *context.Context, columns map[string]interface{},datas []map[string]interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]map[string]interface{}, error) {
  var err error
  selectColumns, selectRelations := model.GetSelectInfo(columns, Provide__Test())
  if datas == nil {
    datas = make([]map[string]interface{}, 0)
    err = Query__Test().Scopes(scopes...).Select(selectColumns).Find(&datas).Error
    if err != nil {
      return nil, err
    }
  }
  var wg sync.WaitGroup
  errChan := make(chan error, len(datas)*len(selectRelations))
  var mu sync.Mutex
  
  for _, data := range datas {
    for key, relation := range selectRelations {
      wg.Add(1)
      go func(data map[string]interface{}, relation *model.SelectRelation)  {
        defer wg.Done()
        cData, err := model.FetchRelation(ctx, data, relation)
        if err != nil {
          errChan <- err
        }
        mu.Lock()
        defer mu.Unlock()
        data[key] = cData
      }(data, relation) 
    }
  }
  wg.Wait()
  close(errChan)
  for err := range errChan {
    return nil, err
  }
  return datas, nil
}


func init() {
  model.AddQuickFirst("Test", First__Test)
  model.AddQuickList("Test", List__Test)
  model.AddQuickLoad("Test", Load__Test)
  model.AddQuickLoadList("Test", LoadList__Test)
}
