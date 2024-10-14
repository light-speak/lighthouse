package validate

import (
	"fmt"
	"strings"

	"github.com/light-speak/lighthouse/errors"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	"github.com/light-speak/lighthouse/utils"
)

var operationArgs = map[string]*ast.ArgumentNode{}

func validateOperation(node ast.Node) error {

	operNode := node.(*ast.OperationNode)

	err := validateOperationArgs(operNode.Args)
	if err != nil {
		return err
	}

	_, ok := p.TypeMap[utils.UcFirst(string(operNode.Type))]
	if !ok {
		return &errors.ValidateError{
			Node:    operNode,
			Message: "node is not an operation",
		}
	}

	fields := node.GetFields()

	for _, field := range fields {
		err = validateFieldArgs(field)
		if err != nil {
			return err
		}
	}

	fields, _ = validateOperationFields(fields)

	var printChildren func(children []*ast.FieldNode, depth int)
	printChildren = func(children []*ast.FieldNode, depth int) {
		for _, f := range children {
			indent := strings.Repeat("  ", depth)
			log.Debug().Msgf("%s子字段名称: %+v", indent, f.GetName())
			if len(f.Children) > 0 {
				printChildren(f.Children, depth+1)
			}
		}
	}

	for _, field := range fields {
		printChildren(field.Children, 0)
	}

	return nil
}

func validateOperationArgs(args []*ast.ArgumentNode) error {
	for _, arg := range args {
		if arg.Type == nil && arg.Value.Type == nil {
			return &errors.ValidateError{
				Node:    arg,
				Message: fmt.Sprintf("argument %s type is invaild", arg.GetName()),
			}
		}
		operationArgs[arg.GetName()] = arg
	}
	return nil
}

func validateFieldArgs(field *ast.FieldNode) error {
	// 找到定义的接口
	t := string(field.Parent.(*ast.OperationNode).Type)
	defField := &ast.FieldNode{}
	for _, f := range p.TypeMap[utils.UcFirst(t)].Fields {
		if f.Name == field.GetName() {
			defField = f
			break
		}
	}

	// 获取后端定义的接口参数和类型
	defArgMap := make(map[string]*ast.ArgumentNode)
	for _, arg := range defField.Args {
		defArgMap[arg.Name] = arg
	}

	// 获取前端输入的参数和类型
	argMap := make(map[string]*ast.ArgumentNode)

	for _, arg := range field.Args {
		argMap[arg.GetName()] = arg
	}

	validateListArg := func(defArg *ast.ArgumentNode, arg *ast.ArgumentNode) error {
		// log.Debug().Msgf("defArg: %+v", defArg.Type)
		log.Debug().Msgf("arg: %+v", arg.Type.TypeCategory)
		if defArg.Type.IsList != arg.Type.IsList {
			return &errors.ValidateError{
				Node:    field,
				Message: fmt.Sprintf("argument %s type is invaild", arg.Name),
			}
		}
		if !defArg.Type.IsList {
			vaild := defArg.Type.Name == arg.Type.Name && defArg.Type.IsList == arg.Type.IsList && defArg.Type.IsNonNull == arg.Type.IsNonNull
			if !vaild {
				return &errors.ValidateError{
					Node:    field,
					Message: fmt.Sprintf("argument %s type is invaild", arg.Name),
				}
			}
		}
		return nil
	}

	for name, defArg := range defArgMap {
		// log.Debug().Msgf("defArg name: %+v", name)
		// log.Debug().Msgf("defArg: %+v", defArg.Value.Children[0])
		log.Debug().Msgf("arg: %+v", argMap[name].Type)
		// 必传参数
		if defArg.Type.IsNonNull {
			if _, ok := argMap[name]; !ok {
				return &errors.ValidateError{
					Node:    field,
					Message: fmt.Sprintf("argument %s is required", name),
				}
			}
		}

		// list 参数，递归检查内层元素是否一致
		if defArg.Type.IsList {
			err := validateListArg(defArg, argMap[name])
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func validateOperationFields(fields []*ast.FieldNode) ([]*ast.FieldNode, error) {
	newFields := make(map[string]*ast.FieldNode)
	var res []*ast.FieldNode

	var err error
	processField := func(f *ast.FieldNode) {
		if _, ok := newFields[f.GetName()]; !ok {
			res = append(res, f)
			newFields[f.GetName()] = f
		}
	}

	for _, field := range fields {
		if field.Type != nil && field.Type.TypeCategory == ast.NodeTypeFragment {
			fragment := getValueTypeFragment(field.GetName())
			if fragment == nil {
				return nil, &errors.ValidateError{
					Node:    field,
					Message: fmt.Sprintf("fragment %s not found", field.GetName()),
				}
			}
			for _, f := range fragment.Fields {
				processField(f)
			}
		} else {
			processField(field)
			if err != nil {
				return nil, err
			}
		}
	}

	for _, field := range res {
		if len(field.Children) > 0 {
			field.Children, _ = validateOperationFields(field.Children)
		}
	}

	return res, nil
}
