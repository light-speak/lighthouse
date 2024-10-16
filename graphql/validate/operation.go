package validate

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/light-speak/lighthouse/errors"
// 	"github.com/light-speak/lighthouse/graphql/ast"
// 	"github.com/light-speak/lighthouse/log"
// 	"github.com/light-speak/lighthouse/utils"
// )

// var operationArgs = map[string]*ast.ArgumentNode{}

// func validateOperation(node ast.Node) error {

// 	operNode := node.(*ast.OperationNode)

// 	defFields, ok := p.TypeMap[utils.UcFirst(string(operNode.Type))]
// 	if !ok {
// 		return &errors.ValidateError{
// 			Node:    operNode,
// 			Message: "node is not an operation",
// 		}
// 	}

// 	err := validateOperationArgs(operNode.Args)

// 	if err != nil {
// 		return err
// 	}

// 	fields := node.GetFields()

// 	for _, field := range fields {
// 		err = validateFieldArgs(field)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	fields, _ = fillOperationFields(fields)

// 	var printChildren func(children []*ast.FieldNode, depth int)
// 	printChildren = func(children []*ast.FieldNode, depth int) {
// 		for _, f := range children {
// 			indent := strings.Repeat("  ", depth)
// 			log.Debug().Msgf("%s子字段名称: %+v", indent, f.GetName())
// 			if len(f.Children) > 0 {
// 				printChildren(f.Children, depth+1)
// 			}
// 		}
// 	}

// 	for _, field := range fields {
// 		printChildren(field.Children, 0)
// 	}

// 	for _, field := range fields {
// 		defField := findDefinedField(defFields.Fields, field.GetName())
// 		err = validateOperationField(field, defField)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func validateOperationArgs(args []*ast.ArgumentNode) error {
// 	for _, arg := range args {
// 		if arg.Type == nil && arg.Value.Type == nil {
// 			return &errors.ValidateError{
// 				Node:    arg,
// 				Message: fmt.Sprintf("argument %s type is invaild", arg.GetName()),
// 			}
// 		}
// 		operationArgs[arg.GetName()] = arg
// 	}
// 	return nil
// }

// func validateFieldArgs(field *ast.FieldNode) error {
// 	// 获取操作类型
// 	operationType := string(field.Parent.(*ast.OperationNode).Type)

// 	// 查找定义的字段
// 	defField := findDefinedField(p.TypeMap[utils.UcFirst(operationType)].Fields, field.GetName())
// 	if defField == nil {
// 		return &errors.ValidateError{
// 			Node:    field,
// 			Message: fmt.Sprintf("字段 %s 未在定义中找到", field.GetName()),
// 		}
// 	}

// 	// 创建参数映射
// 	defArgMap := createArgMap(defField.Args)
// 	argMap := make(map[string]*ast.ArgumentNode)
// 	for _, arg := range field.Args {
// 		argMap[arg.GetName()] = operationArgs[arg.Type.Name]
// 	}

// 	// 验证必需参数
// 	if err := validateRequiredArgs(defArgMap, argMap, field); err != nil {
// 		return err
// 	}

// 	// 验证列表参数
// 	if err := validateListArgs(defArgMap, argMap, field); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func findDefinedField(fields []*ast.FieldNode, name string) *ast.FieldNode {
// 	for _, f := range fields {
// 		if f.Name == name {
// 			return f
// 		}
// 	}
// 	return nil
// }

// func createArgMap(args []*ast.ArgumentNode) map[string]*ast.ArgumentNode {
// 	argMap := make(map[string]*ast.ArgumentNode)
// 	for _, arg := range args {
// 		argMap[arg.GetName()] = arg
// 	}
// 	return argMap
// }

// func createFieldMap(fields []*ast.FieldNode) map[string]*ast.FieldNode {
// 	fieldMap := make(map[string]*ast.FieldNode)
// 	for _, field := range fields {
// 		fieldMap[field.GetName()] = field
// 	}
// 	return fieldMap
// }

// func validateRequiredArgs(defArgMap, argMap map[string]*ast.ArgumentNode, field *ast.FieldNode) error {
// 	for name, defArg := range defArgMap {
// 		if defArg.Type.IsNonNull {
// 			if arg, ok := argMap[name]; ok {
// 				if defArg.Type.IsNonNull && !arg.Type.IsNonNull {
// 					log.Debug().Msgf("defArg: %+v", defArg.Type)
// 					log.Debug().Msgf("arg: %+v", arg.Type)
// 					return &errors.ValidateError{
// 						Node:    field,
// 						Message: fmt.Sprintf("参数 %s 是必需的", defArg.Name),
// 					}
// 				}
// 			} else {
// 				return &errors.ValidateError{
// 					Node:    field,
// 					Message: fmt.Sprintf("参数 %s 未找到", name),
// 				}
// 			}

// 		}
// 	}
// 	return nil
// }

// func validateListArgs(defArgMap, argMap map[string]*ast.ArgumentNode, field *ast.FieldNode) error {
// 	for name, defArg := range defArgMap {
// 		if defArg.Type.IsList {
// 			if arg, ok := argMap[name]; ok {
// 				if err := validateListArg(defArg.Type, arg.Type, field); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	}
// 	return nil
// }

// func validateListArg(defArg, arg *ast.FieldType, field *ast.FieldNode) error {
// 	operationArg, ok := operationArgs[arg.Name]
// 	if !ok {
// 		return &errors.ValidateError{
// 			Node:    field,
// 			Message: fmt.Sprintf("参数 %s 未找到", arg.Name),
// 		}
// 	}

// 	artType := operationArg.Type

// 	if defArg.IsList != artType.IsList {
// 		return &errors.ValidateError{
// 			Node:    field,
// 			Message: fmt.Sprintf("参数 %s 类型无效", arg.Name),
// 		}
// 	}

// 	for defArg.IsList {
// 		if !isValidListType(defArg, artType) {
// 			return &errors.ValidateError{
// 				Node:    field,
// 				Message: fmt.Sprintf("参数 %s 类型无效", arg.Name),
// 			}
// 		}
// 		defArg = defArg.ElemType
// 		artType = artType.ElemType
// 	}

// 	if !isValidType(defArg, artType) {
// 		return &errors.ValidateError{
// 			Node:    field,
// 			Message: fmt.Sprintf("参数 %s 类型无效", arg.Name),
// 		}
// 	}

// 	return nil
// }

// func isValidListType(defArg, artType *ast.FieldType) bool {
// 	return defArg.IsList == artType.IsList && defArg.IsNonNull == artType.IsNonNull
// }

// func isValidType(defArg, artType *ast.FieldType) bool {
// 	return defArg.Name == artType.Name && defArg.IsList == artType.IsList && defArg.IsNonNull == artType.IsNonNull
// }

// func fillOperationFields(fields []*ast.FieldNode) ([]*ast.FieldNode, error) {
// 	newFields := make(map[string]*ast.FieldNode)
// 	var res []*ast.FieldNode

// 	processField := func(f *ast.FieldNode) {
// 		if _, ok := newFields[f.GetName()]; !ok {
// 			res = append(res, f)
// 			newFields[f.GetName()] = f
// 		}
// 	}

// 	for _, field := range fields {
// 		if field.Type != nil && field.Type.TypeCategory == ast.NodeTypeFragment {
// 			fragment := getValueTypeFragment(field.GetName())
// 			if fragment == nil {
// 				return nil, &errors.ValidateError{
// 					Node:    field,
// 					Message: fmt.Sprintf("fragment %s not found", field.GetName()),
// 				}
// 			}
// 			for _, f := range fragment.Fields {
// 				processField(f)
// 			}
// 		} else {
// 			processField(field)
// 		}
// 	}

// 	for _, field := range res {
// 		if len(field.Children) > 0 {
// 			field.Children, _ = fillOperationFields(field.Children)
// 		}
// 	}

// 	return res, nil
// }

// // func isTypeInUnion(typeName string, unionNode *ast.UnionNode) (string, bool) {
// // 	for _, t := range unionNode.Types {
// // 		if t == typeName {
// // 			return t, true
// // 		}
// // 	}
// // 	return "", false
// // }

// // func validateUnionField(field *ast.FieldNode, defField *ast.FieldNode) error {
// // 	// 遍历前端请求的 union 字段
// // 	for _, unionField := range field.Children {
// // 		// 检查请求的类型是否存在于后端定义的 union 中
// // 		defUnion := defField.Type.Type.(*ast.UnionNode)
// // 		_, ok := isTypeInUnion(unionField.GetName(), defUnion)
// // 		if !ok {
// // 			return &errors.ValidateError{
// // 				Node:    unionField,
// // 				Message: fmt.Sprintf("类型 %s 不是 union %s 的有效成员", unionField.GetName(), defField.Type.Name),
// // 			}
// // 		}

// // 		// 解析 fragment 并添加字段
// // 		unionChildren, err := fillOperationFields(unionField.Children)
// // 		if err != nil {
// // 			return nil
// // 		}

// // 		unionField.Children = unionChildren

// // 	}
// // 	return nil
// // }

// func validateOperationField(field *ast.FieldNode, defField *ast.FieldNode) error {
// 	if defField.Type == nil {
// 		return &errors.ValidateError{
// 			Node:    field,
// 			Message: fmt.Sprintf("字段 %s 未在定义中找到", field.GetName()),
// 		}
// 	}

// 	if defField.Type.TypeCategory == ast.NodeTypeUnion {
// 		log.Debug().Msgf("defField: %+v", defField.Type.Name)
// 		log.Debug().Msgf("field: %+v", field)
// 	}

// 	// 拿到后端定义的字段 map -> fieldMap
// 	defFieldMap := createFieldMap(p.TypeMap[defField.Type.Name].Fields)
// 	// 遍历前端请求的字段，然后判断是否在 fieldMap 中，如果存在，则判断名称和类型是否一致
// 	for _, child := range field.Children {
// 		v, ok := defFieldMap[child.GetName()]
// 		if !ok {
// 			return &errors.ValidateError{
// 				Node:    field,
// 				Message: fmt.Sprintf("字段 %s 未在定义中找到", child.GetName()),
// 			}
// 		}
// 		// 附带参数
// 		if len(child.Args) > 0 || len(v.Args) > 0 {
// 			defArgMap := createArgMap(v.Args)
// 			argMap := make(map[string]*ast.ArgumentNode)
// 			for _, arg := range child.Args {
// 				argMap[arg.GetName()] = operationArgs[arg.Type.Name]
// 			}
// 			// 验证必需参数
// 			if err := validateRequiredArgs(defArgMap, argMap, child); err != nil {
// 				return err
// 			}
// 			// 验证列表参数
// 			if err := validateListArgs(defArgMap, argMap, child); err != nil {
// 				return err
// 			}
// 		}
// 		// 字段下还有子字段，则递归调用 validateOperationField 检查
// 		if len(child.Children) > 0 {
// 			return validateOperationField(child, v)
// 		}
// 	}

// 	return nil
// }
