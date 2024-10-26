package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
)

func ExecuteQuery(query string, variables map[string]any) (interface{}, error) {
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
	err := qp.Validate(p.NodeStore)
	if err != nil {
		return nil, &errors.GraphQLError{
			Message:   err.Error(),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
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
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("operation type %s not supported", qp.OperationType),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}

	for _, field := range qp.Fields {
		quickRes, isQuick, err := QuickExecute(field)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   err.Error(),
				Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
			}
		}
		if quickRes != nil {
			if err != nil {
				return nil, &errors.GraphQLError{
					Message:   err.Error(),
					Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
				}
			}
			res[field.Name] = quickRes
			continue
		}
		if isQuick {
			if field.Type.Kind == ast.KindNonNull {
				return nil, &errors.GraphQLError{
					Message:   "field is not nullable",
					Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
				}
			}
			res[field.Name] = nil
			continue
		}
		queryFunc, ok := funMap[field.Name]
		if !ok {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("query %s not found", field.Name),
				Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
			}
		}
		res[field.Name], err = queryFunc(qp, field)
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   err.Error(),
				Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
			}
		}
	}
	return res, nil
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
