// Code generated by github.com/light-speak/lighthouse, YOU CAN FUCKING EDIT BY YOURSELF.
package resolver

import (
  "github.com/light-speak/lighthouse/graphql"
  "github.com/light-speak/lighthouse/context"
  "github.com/light-speak/lighthouse/graphql/model"
  "github.com/light-speak/lighthouse/graphql/excute"
  "fmt"
  "github.com/light-speak/lighthouse/resolve"
  "user/models"
)

func init() {
  excute.AddResolver("getPost", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    pfuck, e := graphql.Parser.NodeStore.Scalars["String"].ScalarType.ParseValue(args["fuck"], nil)
    if e != nil {
      return nil, e
    }
    fuck, ok := pfuck.(string)
    if !ok {
      return nil, fmt.Errorf("argument: 'fuck' is not a string, got %T", args["fuck"])
    }
    res, err := r.GetPostResolver(ctx, fuck)
    if res == nil {
      return nil, err
    }
    return model.StructToMap(res)
  })
  excute.AddResolver("getPostIds", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    res, err := r.GetPostIdsResolver(ctx)
    if res == nil {
      return nil, err
    }
    return res, nil
  })
  excute.AddResolver("getPosts", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    pfuck, e := graphql.Parser.NodeStore.Scalars["String"].ScalarType.ParseValue(args["fuck"], nil)
    if e != nil {
      return nil, e
    }
    fuck, ok := pfuck.(string)
    if !ok {
      return nil, fmt.Errorf("argument: 'fuck' is not a string, got %T", args["fuck"])
    }
    list, err := r.GetPostsResolver(ctx, fuck)
    if list == nil {
      return nil, err
    }
    res := []map[string]interface{}{}
    for _, item := range list {
      itemMap, err := model.StructToMap(item)
      if err != nil {
        return nil, err
      }
      res = append(res, itemMap)
    }
    return res, nil
  })
  excute.AddResolver("testNullableEnum", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    enumValue, ok := models.TestEnumMap[args["enum"].(string)]
    if !ok {
      return nil, fmt.Errorf("argument: 'enum' is not a models.TestEnum, got %T", args["enum"])
    }
    enum := &enumValue
    res, err := r.TestNullableEnumResolver(ctx, enum)
    return res, err
  })
  excute.AddResolver("testPostEnum", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    enumValue, ok := models.TestEnumMap[args["enum"].(string)]
    if !ok {
      return nil, fmt.Errorf("argument: 'enum' is not a models.TestEnum, got %T", args["enum"])
    }
    enum := &enumValue
    res, err := r.TestPostEnumResolver(ctx, enum)
    return res, err
  })
  excute.AddResolver("testPostId", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    pid, e := graphql.Parser.NodeStore.Scalars["ID"].ScalarType.ParseValue(args["id"], nil)
    if e != nil {
      return nil, e
    }
    id, ok := pid.(int64)
    if !ok {
      return nil, fmt.Errorf("argument: 'id' is not a int64, got %T", args["id"])
    }
    res, err := r.TestPostIdResolver(ctx, id)
    if res == nil {
      return nil, err
    }
    return model.StructToMap(res)
  })
  excute.AddResolver("testPostInput", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    input, err := models.MapToTestInput(args["input"].(map[string]interface{}))
    if err != nil {
      return nil, fmt.Errorf("argument: 'input' can not convert to models.TestInput, got %T", args["input"])
    }
    res, err := r.TestPostInputResolver(ctx, input)
    return res, err
  })
  excute.AddResolver("testPostInt", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    pid, e := graphql.Parser.NodeStore.Scalars["Boolean"].ScalarType.ParseValue(args["id"], nil)
    if e != nil {
      return nil, e
    }
    id, ok := pid.(bool)
    if !ok {
      return nil, fmt.Errorf("argument: 'id' is not a bool, got %T", args["id"])
    }
    res, err := r.TestPostIntResolver(ctx, id)
    if res == nil {
      return nil, err
    }
    return model.StructToMap(res)
  })
  excute.AddResolver("createPost", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    input, err := models.MapToTestInput(args["input"].(map[string]interface{}))
    if err != nil {
      return nil, fmt.Errorf("argument: 'input' can not convert to models.TestInput, got %T", args["input"])
    }
    res, err := r.CreatePostResolver(ctx, input)
    if res == nil {
      return nil, err
    }
    return model.StructToMap(res)
  })
  excute.AddResolver("login", func(ctx *context.Context, args map[string]any, resolve resolve.Resolve) (interface{}, error) {
    r := resolve.(*Resolver)
    pname, e := graphql.Parser.NodeStore.Scalars["String"].ScalarType.ParseValue(args["name"], nil)
    if e != nil {
      return nil, e
    }
    name, ok := pname.(string)
    if !ok {
      return nil, fmt.Errorf("argument: 'name' is not a string, got %T", args["name"])
    }
    res, err := r.LoginResolver(ctx, name)
    if res == nil {
      return nil, err
    }
    return model.TypeToMap(res)
  })
}
