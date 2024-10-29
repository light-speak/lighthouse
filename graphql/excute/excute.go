package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/context"
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/log"
)

func ExecuteQuery(ctx *context.Context, query string, variables map[string]any) interface{} {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(errors.GraphqlErrorInterface); ok {
				ctx.Errors = append(ctx.Errors, err.GraphqlError())
				return
			}
			panic(r)
		}
	}()
	p := graphql.GetParser()
	qp := p.NewQueryParser(lexer.NewLexer([]*lexer.Content{
		{
			Content: query,
		},
	}))
	qp.ParseSchema()
	qp.Variables = make(map[string]any)
	for k, v := range variables {
		qp.Variables[fmt.Sprintf("$%s", k)] = v
	}
	e := qp.Validate(p.NodeStore)
	if e != nil {
		ctx.Errors = append(ctx.Errors, e)
		panic(e)
	}
	res := make(map[string]interface{})

	var funMap map[string]func(qp *parser.QueryParser, field *ast.Field) (interface{}, error)
	switch qp.OperationType {
	case "Mutation":
		funMap = mutationMap
	case "Subscription":
		funMap = subscriptionMap
	case "Query":
		funMap = queryMap
	default:
		e := &errors.ParserError{
			Message:   fmt.Sprintf("operation type %s not supported", qp.OperationType),
			Locations: &errors.GraphqlLocation{Line: 1, Column: 1},
		}
		ctx.Errors = append(ctx.Errors, e)
		panic(e)
	}

	for _, field := range qp.Fields {
		quickRes, isQuick, err := QuickExecute(ctx, field)
		if err != nil {
			ctx.Errors = append(ctx.Errors, err)
			continue
		}
		if quickRes != nil {
			res[field.Name] = quickRes
			continue
		}
		if isQuick {
			if field.Type.Kind == ast.KindNonNull {
				e := &errors.GraphQLError{
					Message:   "field is not nullable",
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
				ctx.Errors = append(ctx.Errors, e)
				continue
			}
			res[field.Name] = nil
			continue
		}
		queryFunc, ok := funMap[field.Name]
		if ok {
			r, e := queryFunc(qp, field)
			if e != nil {
				ee := &errors.GraphQLError{
					Message:   e.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
				ctx.Errors = append(ctx.Errors, ee)
			}
			res[field.Name] = r
			continue
		}

		resolverFunc, ok := resolverMap[field.Name]
		if ok {
			log.Debug().Msgf("resolverFunc Args: %+v", field.Args)
			r, e := resolverFunc(ctx, nil)
			if e != nil {
				ee := &errors.GraphQLError{
					Message:   e.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
				ctx.Errors = append(ctx.Errors, ee)
				continue
			}
			res[field.Name] = r
			continue
		}

		ee := &errors.GraphQLError{
			Message:   fmt.Sprintf("query %s not found", field.Name),
			Locations: []*errors.GraphqlLocation{field.GetLocation()},
		}
		ctx.Errors = append(ctx.Errors, ee)
		continue
	}
	return res
}

var queryMap = map[string]func(qp *parser.QueryParser, field *ast.Field) (interface{}, error){
	"__schema": graphql.ResolveSchemaFields,
	"__type":   graphql.ResolveTypeByName,
}

var mutationMap = make(map[string]func(qp *parser.QueryParser, field *ast.Field) (interface{}, error))

var subscriptionMap = make(map[string]func(qp *parser.QueryParser, field *ast.Field) (interface{}, error))

func AddQuery(name string, fn func(qp *parser.QueryParser, field *ast.Field) (interface{}, error)) {
	queryMap[name] = fn
}

func AddMutation(name string, fn func(qp *parser.QueryParser, field *ast.Field) (interface{}, error)) {
	mutationMap[name] = fn
}

func AddSubscription(name string, fn func(qp *parser.QueryParser, field *ast.Field) (interface{}, error)) {
	subscriptionMap[name] = fn
}

var resolverMap = make(map[string]func(ctx *context.Context, args map[string]any) (interface{}, error))

func AddResolver(name string, fn func(ctx *context.Context, args map[string]any) (interface{}, error)) {
	resolverMap[name] = fn
}
