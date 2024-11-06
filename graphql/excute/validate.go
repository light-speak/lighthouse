package excute

import (
	"fmt"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
)

func ValidateValue(field *ast.Field, value interface{}, isVariable bool) (interface{}, errors.GraphqlErrorInterface) {
	realType := field.Type.GetRealType()
	var v interface{}
	var err errors.GraphqlErrorInterface
	if field.Name == "__typename" {
		return value, nil
	}

	switch realType.Kind {
	case ast.KindScalar:
		if scalarNode, ok := realType.TypeNode.(*ast.ScalarNode); ok {
			v, err = scalarNode.ScalarType.Serialize(value, field.GetLocation())
			if err != nil {
				return nil, &errors.GraphQLError{
					Message:   err.Error(),
					Locations: []*errors.GraphqlLocation{field.GetLocation()},
				}
			}
		} else {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("scalar type is not a scalar node, field: %s", field.Name),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}
	case ast.KindEnum:
		if e, ok := value.(model.EnumInterface); ok {
			v = e.ToString()
		} else {
			return nil, &errors.GraphQLError{
				Message:   fmt.Sprintf("enum value type not supported, field: %s， got %v , type: %T", field.Name, value, value),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}
	}
	return v, nil
}
