package excute

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model"
)

func ValidateValue(field *ast.Field, value interface{}, isVariable bool) (interface{}, errors.GraphqlErrorInterface) {
	realType := field.Type.GetRealType()
	var v interface{}
	var err errors.GraphqlErrorInterface
	switch realType.Kind {
	case ast.KindScalar:
		v, err = realType.TypeNode.(*ast.ScalarNode).ScalarType.Serialize(value, field.GetLocation())
		if err != nil {
			return nil, &errors.GraphQLError{
				Message:   err.Error(),
				Locations: []*errors.GraphqlLocation{field.GetLocation()},
			}
		}
	case ast.KindEnum:
		v = value.(model.EnumInterface).ToString()
	}
	return v, nil
}
