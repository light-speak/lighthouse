// Code generated by github.com/light-speak/lighthouse, DO NOT EDIT.
package resolver

import (
  "github.com/light-speak/lighthouse/context"
  "fmt"
  "user/models"
  "github.com/light-speak/lighthouse/graphql/excute"
)

func init() {
  excute.AddResolver("getPost", func(ctx *context.Context, args map[string]any) (interface{}, error) {
    fuck, ok := args["fuck"].(string)
    if !ok {
      return nil, fmt.Errorf("argument: 'fuck' is not a string, got %T", args["fuck"])
    }
    return GetPostResolver(ctx, fuck)
  })
  excute.AddResolver("testPostEnum", func(ctx *context.Context, args map[string]any) (interface{}, error) {
    enum, ok := models.TestEnumMap[args["enum"].(string)]
    if !ok {
      return nil, fmt.Errorf("argument: 'enum' is not a models.TestEnum, got %T", args["enum"])
    }
    return TestPostEnumResolver(ctx, enum)
  })
  excute.AddResolver("testPostInput", func(ctx *context.Context, args map[string]any) (interface{}, error) {
    input, ok := args["input"].(models.TestInput)
    if !ok {
      return nil, fmt.Errorf("argument: 'input' is not a models.TestInput, got %T", args["input"])
    }
    return TestPostInputResolver(ctx, input)
  })
}
