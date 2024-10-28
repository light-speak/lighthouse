package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
)

func ValidateValue(field *ast.Field, value interface{}, isVariable bool) (interface{}, error) {
	realType := field.Type.GetRealType()
	if realType.Kind != ast.KindScalar {
		return nil, &errors.GraphQLError{
			Message:   fmt.Sprintf("field %s is not a scalar type", field.Name),
			Locations: []errors.GraphqlLocation{{Line: 1, Column: 1}},
		}
	}
	value, err := realType.TypeNode.(*ast.ScalarNode).ScalarType.Serialize(value)
	if err != nil {
		return nil, err
	}
	return value, nil
}
