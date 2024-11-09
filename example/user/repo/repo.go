// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package repo

import (
  "sync"
  "user/models"
  "github.com/light-speak/lighthouse/context"
  "github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/utils"
	"gorm.io/gorm"  
)

// Generic loader function
func loadEntity[T any](ctx *context.Context, key int64, table string, field string, filters ...*model.Filter) (*sync.Map, error) {
  data, err := model.GetLoader[int64](model.GetDB(), table, field, filters).Load(key)
  if err != nil {
    return nil, err
  }
  return utils.MapToSyncMap(data), nil
}

// Generic list loader function  
func loadEntityList[T any](ctx *context.Context, key int64, table string, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  datas, err := model.GetLoader[int64](model.GetDB(), table, field, filters).LoadList(key)
  if err != nil {
    return nil, err
  }
  return utils.MapSliceToSyncMapSlice(datas), nil
}

// Generic query function
func queryEntity[T any](m interface{}, scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return model.GetDB().Model(m).Scopes(scopes...)
}

// Generic first function
func firstEntity[T any](ctx *context.Context, data *sync.Map, enumFieldsFn func(string) func(interface{}) interface{}, 
  model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  
  var err error
  var mu sync.Mutex
  
  if data == nil {
    mapData := make(map[string]interface{})
    err = queryEntity[T](model).Scopes(scopes...).First(&mapData).Error
    if err != nil {
      return nil, err
    }
    data = utils.MapToSyncMap(mapData)
  }

  result := &sync.Map{}
  data.Range(func(key, value interface{}) bool {
    k := key.(string)
    if fn := enumFieldsFn(k); fn != nil {
      mu.Lock()
      result.Store(k, fn(value))
      mu.Unlock()
    } else {
      result.Store(k, value)
    }
    return true
  })
  return result, nil
}

// Generic list function
func listEntity[T any](ctx *context.Context, datas []*sync.Map, enumFieldsFn func(string) func(interface{}) interface{}, model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  if datas == nil {
    mapDatas := make([]map[string]interface{}, 0)
    err := queryEntity[T](model).Scopes(scopes...).Find(&mapDatas).Error
    if err != nil {
      return nil, err
    }
    datas = utils.MapSliceToSyncMapSlice(mapDatas)
  }

  var mu sync.Mutex
  results := make([]*sync.Map, len(datas))
  
  for i, data := range datas {
    result := &sync.Map{}
    data.Range(func(key, value interface{}) bool {
      k := key.(string)
      if fn := enumFieldsFn(k); fn != nil {
        mu.Lock()
        result.Store(k, fn(value))
        mu.Unlock()
      } else {
        result.Store(k, value)
      }
      return true
    })
    results[i] = result
  }
  
  return results, nil
}

// Generic count function
func countEntity[T any](model interface{}, scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  var count int64
  err := queryEntity[T](model).Scopes(scopes...).Count(&count).Error
  return count, err
}

// UserTag functions
func Load__UserTag(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.UserTag](ctx, key, "user_tags", field, filters...)
}

func LoadList__UserTag(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.UserTag](ctx, key, "user_tags", field, filters...)
}

func Query__UserTag(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.UserTag](&models.UserTag{}, scopes...)
}

func First__UserTag(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.UserTag](ctx, data, models.UserTagEnumFields, &models.UserTag{}, scopes...)
}

func List__UserTag(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.UserTag](ctx, datas, models.UserTagEnumFields, &models.UserTag{}, scopes...)
}

func Count__UserTag(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.UserTag](&models.UserTag{}, scopes...)
}
// Wallet functions
func Load__Wallet(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.Wallet](ctx, key, "wallets", field, filters...)
}

func LoadList__Wallet(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.Wallet](ctx, key, "wallets", field, filters...)
}

func Query__Wallet(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.Wallet](&models.Wallet{}, scopes...)
}

func First__Wallet(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.Wallet](ctx, data, models.WalletEnumFields, &models.Wallet{}, scopes...)
}

func List__Wallet(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.Wallet](ctx, datas, models.WalletEnumFields, &models.Wallet{}, scopes...)
}

func Count__Wallet(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.Wallet](&models.Wallet{}, scopes...)
}
// Comment functions
func Load__Comment(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.Comment](ctx, key, "comments", field, filters...)
}

func LoadList__Comment(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.Comment](ctx, key, "comments", field, filters...)
}

func Query__Comment(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.Comment](&models.Comment{}, scopes...)
}

func First__Comment(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.Comment](ctx, data, models.CommentEnumFields, &models.Comment{}, scopes...)
}

func List__Comment(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.Comment](ctx, datas, models.CommentEnumFields, &models.Comment{}, scopes...)
}

func Count__Comment(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.Comment](&models.Comment{}, scopes...)
}
// Tag functions
func Load__Tag(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.Tag](ctx, key, "tags", field, filters...)
}

func LoadList__Tag(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.Tag](ctx, key, "tags", field, filters...)
}

func Query__Tag(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.Tag](&models.Tag{}, scopes...)
}

func First__Tag(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.Tag](ctx, data, models.TagEnumFields, &models.Tag{}, scopes...)
}

func List__Tag(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.Tag](ctx, datas, models.TagEnumFields, &models.Tag{}, scopes...)
}

func Count__Tag(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.Tag](&models.Tag{}, scopes...)
}
// User functions
func Load__User(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.User](ctx, key, "users", field, filters...)
}

func LoadList__User(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.User](ctx, key, "users", field, filters...)
}

func Query__User(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.User](&models.User{}, scopes...)
}

func First__User(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.User](ctx, data, models.UserEnumFields, &models.User{}, scopes...)
}

func List__User(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.User](ctx, datas, models.UserEnumFields, &models.User{}, scopes...)
}

func Count__User(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.User](&models.User{}, scopes...)
}
// Article functions
func Load__Article(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.Article](ctx, key, "articles", field, filters...)
}

func LoadList__Article(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.Article](ctx, key, "articles", field, filters...)
}

func Query__Article(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.Article](&models.Article{}, scopes...)
}

func First__Article(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.Article](ctx, data, models.ArticleEnumFields, &models.Article{}, scopes...)
}

func List__Article(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.Article](ctx, datas, models.ArticleEnumFields, &models.Article{}, scopes...)
}

func Count__Article(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.Article](&models.Article{}, scopes...)
}
// Post functions
func Load__Post(ctx *context.Context, key int64, field string, filters ...*model.Filter) (*sync.Map, error) {
  return loadEntity[models.Post](ctx, key, "posts", field, filters...)
}

func LoadList__Post(ctx *context.Context, key int64, field string, filters ...*model.Filter) ([]*sync.Map, error) {
  return loadEntityList[models.Post](ctx, key, "posts", field, filters...)
}

func Query__Post(scopes ...func(db *gorm.DB) *gorm.DB) *gorm.DB {
  return queryEntity[models.Post](&models.Post{}, scopes...)
}

func First__Post(ctx *context.Context, data *sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) (*sync.Map, error) {
  return firstEntity[models.Post](ctx, data, models.PostEnumFields, &models.Post{}, scopes...)
}

func List__Post(ctx *context.Context, datas []*sync.Map, scopes ...func(db *gorm.DB) *gorm.DB) ([]*sync.Map, error) {
  return listEntity[models.Post](ctx, datas, models.PostEnumFields, &models.Post{}, scopes...)
}

func Count__Post(scopes ...func(db *gorm.DB) *gorm.DB) (int64, error) {
  return countEntity[models.Post](&models.Post{}, scopes...)
}


func init() {
  model.AddQuickFirst("UserTag", First__UserTag)
  model.AddQuickList("UserTag", List__UserTag)
  model.AddQuickLoad("UserTag", Load__UserTag)
  model.AddQuickLoadList("UserTag", LoadList__UserTag)
  model.AddQuickCount("UserTag", Count__UserTag)
  model.AddQuickFirst("Wallet", First__Wallet)
  model.AddQuickList("Wallet", List__Wallet)
  model.AddQuickLoad("Wallet", Load__Wallet)
  model.AddQuickLoadList("Wallet", LoadList__Wallet)
  model.AddQuickCount("Wallet", Count__Wallet)
  model.AddQuickFirst("Comment", First__Comment)
  model.AddQuickList("Comment", List__Comment)
  model.AddQuickLoad("Comment", Load__Comment)
  model.AddQuickLoadList("Comment", LoadList__Comment)
  model.AddQuickCount("Comment", Count__Comment)
  model.AddQuickFirst("Tag", First__Tag)
  model.AddQuickList("Tag", List__Tag)
  model.AddQuickLoad("Tag", Load__Tag)
  model.AddQuickLoadList("Tag", LoadList__Tag)
  model.AddQuickCount("Tag", Count__Tag)
  model.AddQuickFirst("User", First__User)
  model.AddQuickList("User", List__User)
  model.AddQuickLoad("User", Load__User)
  model.AddQuickLoadList("User", LoadList__User)
  model.AddQuickCount("User", Count__User)
  model.AddQuickFirst("Article", First__Article)
  model.AddQuickList("Article", List__Article)
  model.AddQuickLoad("Article", Load__Article)
  model.AddQuickLoadList("Article", LoadList__Article)
  model.AddQuickCount("Article", Count__Article)
  model.AddQuickFirst("Post", First__Post)
  model.AddQuickList("Post", List__Post)
  model.AddQuickLoad("Post", Load__Post)
  model.AddQuickLoadList("Post", LoadList__Post)
  model.AddQuickCount("Post", Count__Post)
}
