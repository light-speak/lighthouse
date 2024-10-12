package validate

import (
	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

func validateOperation(node ast.Node) error {

	operNode := node.(*ast.OperationNode)

	typeMap, ok := p.TypeMap[utils.UcFirst(string(operNode.Type))]
	if !ok {
		return &errors.ValidateError{
			Node:    operNode,
			Message: "node is not an operation",
		}
	}

	// 步骤1：获取操作的字段
	fields := node.GetFields()

	// 步骤2：验证字段是否存在 验证字段返回类型是否有效，并验证子字段是否匹配返回类型
	for _, field := range fields {
		if fieldNode, err := validateOperationFieldType(field, typeMap); fieldNode != nil {
			if err := validateFieldReturnType(field, fieldNode); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// 步骤4：递归验证子字段
	// for _, field := range fields {
	// 	if err := validateChildFields(field); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func validateOperationFieldType(field *ast.FieldNode, typeMap *ast.TypeNode) (*ast.FieldNode, error) {

	// 获取 typeMap 中的字段
	typeMapFields := typeMap.GetFields()

	// 检查 field 是否存在于 typeMap 中
	for _, typeMapField := range typeMapFields {
		if typeMapField.GetName() == field.GetName() {
			return typeMapField, nil
		}
	}

	return nil, &errors.ValidateError{
		Node:    field,
		Message: "字段在操作类型中不存在",
	}

}

func validateFieldReturnType(field *ast.FieldNode, typeField *ast.FieldNode) error {
	log.Debug().Msgf("field: %+v", field)
	log.Debug().Msgf("typeField: %+v", typeField.Type)
	// log.Debug().Msgf("32131231232131231231231231312312312 %+v", getValueTypeType(typeField.Type.Name))
	return nil
}

func validateChildFields(field *ast.FieldNode) error {
	for _, childField := range field.Children {
		if err := validateOperation(childField); err != nil {
			return err
		}
	}
	return nil
}
