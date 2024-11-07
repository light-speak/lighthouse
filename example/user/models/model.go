// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import  "github.com/light-speak/lighthouse/graphql/model"


type Wallet struct {
  model.Model
  UserId int64 `json:"user_id" `
  User *User `json:"user" gorm:"-" `
  Balance int64 `json:"balance" `
}

func (*Wallet) IsModel() bool { return true }
func (*Wallet) TableName() string { return "wallets" }
func (*Wallet) TypeName() string { return "wallet" }
func WalletEnumFields(key string) func(interface{}) interface{} {
  return nil
}

type Post struct {
  model.ModelSoftDelete
  UserId int64 `json:"user_id" gorm:"index" `
  IsBool bool `json:"is_bool" gorm:"default:false" `
  User *User `json:"user" gorm:"-" `
  Enum TestEnum `json:"enum" `
  BackId int64 `json:"back_id" `
  Comments *[]Comment `json:"comments" gorm:"-" `
  FuckingAttr string `json:"fucking_attr" gorm:"-" `
  Title string `json:"title" gorm:"index;type:varchar(255)" `
  Content string `json:"content" gorm:"type:varchar(255)" `
  TagId int64 `json:"tag_id" `
}

func (*Post) IsModel() bool { return true }
func (*Post) TableName() string { return "posts" }
func (*Post) TypeName() string { return "post" }
func PostEnumFields(key string) func(interface{}) interface{} {
  switch key {
  case "enum":
    return func(value interface{}) interface{} {
      switch v := value.(type) {
      case int64:
        return TestEnum(v)
      case int8:
        return TestEnum(v)
      default:
        return v
      }
    }
  }
  return nil
}

type Comment struct {
  model.Model
  Content string `json:"content" gorm:"type:varchar(255)" `
  CommentableId int64 `json:"commentable_id" gorm:"index:commentable" `
  CommentableType string `json:"commentable_type" gorm:"index:commentable;type:varchar(255)" `
  Commentable interface{} `json:"commentable" gorm:"-" `
}

func (*Comment) IsModel() bool { return true }
func (*Comment) TableName() string { return "comments" }
func (*Comment) TypeName() string { return "comment" }
func CommentEnumFields(key string) func(interface{}) interface{} {
  return nil
}

type User struct {
  model.Model
  Name string `json:"name" gorm:"index;type:varchar(255)" `
  MyPosts *[]Post `json:"my_posts" gorm:"-" `
  Wallet *Wallet `json:"wallet" gorm:"-" `
  Tags *[]Tag `json:"tags" gorm:"-" `
}

func (*User) IsModel() bool { return true }
func (*User) IsHasName() bool { return true }
func (this *User) GetName() string { return this.Name }
func (*User) TableName() string { return "users" }
func (*User) TypeName() string { return "user" }
func UserEnumFields(key string) func(interface{}) interface{} {
  return nil
}

type UserTag struct {
  model.Model
  User *User `json:"user" gorm:"-" `
  Tag *Tag `json:"tag" gorm:"-" `
  UserId int64 `json:"user_id" gorm:"uniqueIndex:user_tag_index" `
  TagId int64 `json:"tag_id" gorm:"uniqueIndex:user_tag_index" `
}

func (*UserTag) IsModel() bool { return true }
func (*UserTag) TableName() string { return "user_tags" }
func (*UserTag) TypeName() string { return "user_tag" }
func UserTagEnumFields(key string) func(interface{}) interface{} {
  return nil
}

type Article struct {
  model.Model
  Name string `json:"name" gorm:"type:varchar(255)" `
  Content string `gorm:"type:varchar(255)" json:"content" `
}

func (*Article) IsModel() bool { return true }
func (*Article) TableName() string { return "articles" }
func (*Article) TypeName() string { return "article" }
func ArticleEnumFields(key string) func(interface{}) interface{} {
  return nil
}

type Tag struct {
  model.Model
  Name string `json:"name" gorm:"type:varchar(255)" `
}

func (*Tag) IsModel() bool { return true }
func (*Tag) TableName() string { return "tags" }
func (*Tag) TypeName() string { return "tag" }
func TagEnumFields(key string) func(interface{}) interface{} {
  return nil
}


func Migrate() error {
	return model.GetDB().AutoMigrate(
    &Wallet{},
    &Post{},
    &Comment{},
    &User{},
    &UserTag{},
    &Article{},
    &Tag{},
  )
}