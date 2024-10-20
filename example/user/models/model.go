// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import (
	"github.com/light-speak/lighthouse/graphql/model"
)


type User struct {
  model.Model
  Name string `json:"name" gorm:"index" `
  Posts PostPaginateResponse `json:"posts" `
}

func (*User) IsModel() bool { return true }
func (*User) IsHasName() bool { return true }
func (this *User) GetName() string { return this.Name }


type Post struct {
  model.ModelSoftDelete
  Title string `json:"title" gorm:"index" `
  Content string `json:"content" `
  UserId int64 `json:"userId" gorm:"index" `
  User User `json:"user" `
}

func (*Post) IsModel() bool { return true }


func init() {
	model.GetDB().AutoMigrate(
    &User{},
    &Post{},
  )
}
