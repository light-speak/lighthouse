// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package models

import "github.com/light-speak/lighthouse/graphql/model"

type Test struct {
	Test string `json:"test" gorm:"type:varchar(255)" `
}

type LoginResponse struct {
	User          *User  `json:"user" `
	Token         string `json:"token" gorm:"type:varchar(255)" `
	Authorization string `json:"authorization" gorm:"type:varchar(255)" `
}

type PostPaginateResponse struct {
	Data         *[]*Post            `json:"data" `
	PaginateInfo *model.PaginateInfo `json:"paginate_info" `
}

type UserPaginateResponse struct {
	Data         *[]*User            `json:"data" `
	PaginateInfo *model.PaginateInfo `json:"paginate_info" `
}
