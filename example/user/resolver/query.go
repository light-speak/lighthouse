// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package resolver

import (
  "fmt"
  "github.com/light-speak/lighthouse/log"
  "github.com/light-speak/lighthouse/context"
  "user/models"
  "github.com/light-speak/lighthouse/graphql/model"
)


// GetPosts <nil>
// 
// Parameters:
// - fuck: <nil>
// 
// Returns:
// 
// - []Post
func (r *Resolver) GetPostsResolver(ctx *context.Context,fuck string) ([]*models.Post, error) {
	// Func:GetPosts user code start. Do not remove this comment.
	posts := []*models.Post{}
	db := model.GetDB()
	db.Find(&posts)
	return posts, nil
	// Func:GetPosts user code end. Do not remove this comment. 
}
// TestNullableEnum <nil>
// 
// Parameters:
// - enum: <nil>
// 
// Returns:
// 
// - string
func (r *Resolver) TestNullableEnumResolver(ctx *context.Context,enum *models.TestEnum) (string, error) {
	// Func:TestNullableEnum user code start. Do not remove this comment.
	panic("not implement")
	// Func:TestNullableEnum user code end. Do not remove this comment. 
}
// TestPostInput <nil>
// 
// Parameters:
// - input: <nil>
// 
// Returns:
// 
// - string
func (r *Resolver) TestPostInputResolver(ctx *context.Context,input *models.TestInput) (string, error) {
	// Func:TestPostInput user code start. Do not remove this comment.
	res := fmt.Sprintf("input: %+v", input)
	return res, nil
	// Func:TestPostInput user code end. Do not remove this comment. 
}
// GetPostIds <nil>
// 
// Parameters:
// 
// Returns:
// 
// - []int64
func (r *Resolver) GetPostIdsResolver(ctx *context.Context) ([]int64, error) {
	// Func:GetPostIds user code start. Do not remove this comment.
	return []int64{1, 2, 3}, nil
	// Func:GetPostIds user code end. Do not remove this comment. 
}
// TestPostEnum <nil>
// 
// Parameters:
// - enum: <nil>
// 
// Returns:
// 
// - string
func (r *Resolver) TestPostEnumResolver(ctx *context.Context,enum *models.TestEnum) (string, error) {
	// Func:TestPostEnum user code start. Do not remove this comment.
	log.Debug().Msgf("enum: %+v", enum)
	res := fmt.Sprintf("啥也不是！：%v", *enum == models.TestEnumA)
	return res, nil
	// Func:TestPostEnum user code end. Do not remove this comment. 
}
// GetPost <nil>
// 
// Parameters:
// - fuck: <nil>
// 
// Returns:
// 
// - Post
func (r *Resolver) GetPostResolver(ctx *context.Context,fuck string) (*models.Post, error) {
	// Func:GetPost user code start. Do not remove this comment.
	log.Debug().Msg("GetPostResolver")
	db := model.GetDB()
	post := &models.Post{}
	db.Where("id = ?", fuck).First(post)
	return post, nil
	// Func:GetPost user code end. Do not remove this comment. 
}
// TestPostId <nil>
// 
// Parameters:
// - id: <nil>
// 
// Returns:
// 
// - *Post
func (r *Resolver) TestPostIdResolver(ctx *context.Context,id int64) (*models.Post, error) {
	// Func:TestPostId user code start. Do not remove this comment.
	log.Debug().Msgf("id: %d", id)
	return nil, nil
	// Func:TestPostId user code end. Do not remove this comment. 
}
// TestPostInt <nil>
// 
// Parameters:
// - id: <nil>
// 
// Returns:
// 
// - *Post
func (r *Resolver) TestPostIntResolver(ctx *context.Context,id bool) (*models.Post, error) {
	// Func:TestPostInt user code start. Do not remove this comment.
	return nil, nil
	// Func:TestPostInt user code end. Do not remove this comment. 
}