// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package resolver

import (
  "fmt"
  "github.com/light-speak/lighthouse/context"
  "user/models"
)


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