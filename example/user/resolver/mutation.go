// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package resolver

import (
	"fmt"
	"user/models"

	"github.com/light-speak/lighthouse/auth"
	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/graphql/model"
	"github.com/light-speak/lighthouse/log"
)

func (r *Resolver) LoginResolver(ctx *context.Context, name string) (*models.LoginResponse, error) {
	// Func:Login user code start. Do not remove this comment.
	user := &models.User{}
	db := model.GetDB()
	if err := db.Model(&models.User{}).Where("id = ?", 2).First(user).Error; err != nil {
		return nil, err
	}
	token, err := auth.GetToken(user.Id)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("currentUser: %v", ctx.UserId)
	return &models.LoginResponse{
		User:          user,
		Token:         token,
		Authorization: fmt.Sprintf("Bearer %s", token),
	}, nil
	// Func:Login user code end. Do not remove this comment.
}
func (r *Resolver) CreatePostResolver(ctx *context.Context, input *models.TestInput) (*models.Post, error) {
	// Func:CreatePost user code start. Do not remove this comment.
	panic("not implement")
	// Func:CreatePost user code end. Do not remove this comment.
}
