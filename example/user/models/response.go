// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import  "github.com/light-speak/lighthouse/graphql/model"


type UserPaginateResponse struct {
  Data *[]*User `json:"data" `
  PaginateInfo *model.PaginateInfo `json:"paginate_info" `
}

type LoginResponse struct {
  Token string `json:"token" gorm:"type:varchar(255)" `
  Authorization string `json:"authorization" gorm:"type:varchar(255)" `
  User *User `json:"user" `
}

type Test struct {
  Test string `json:"test" gorm:"type:varchar(255)" `
}

type PostPaginateResponse struct {
  Data *[]*Post `json:"data" `
  PaginateInfo *model.PaginateInfo `json:"paginate_info" `
}
