// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import  "github.com/light-speak/lighthouse/graphql/model"


type User struct {
  model.Model
  Name string `json:"name" gorm:"index" `
  Posts []Post `json:"posts" `
}

func (*User) IsModel() bool { return true }
func (*User) IsHasName() bool { return true }
func (this *User) GetName() string { return this.Name }
func (*User) TableName() string { return "users" }
func (*User) TypeName() string { return "user" }

type Post struct {
  model.ModelSoftDelete
  UserId int64 `json:"user_id" `
  Title string `json:"title" gorm:"index" `
  Content string `json:"content" `
  User User `json:"user" `
}

func (*Post) IsModel() bool { return true }
func (*Post) TableName() string { return "posts" }
func (*Post) TypeName() string { return "post" }


func Migrate() error {
	return model.GetDB().AutoMigrate(
    &User{},
    &Post{},
  )
}