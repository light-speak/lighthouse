package graphql

import (
	"fmt"
	"strings"

	"github.com/light-speak/lighthouse/graphql/ast"
)

func generateObjectType(node ast.Node) string {
	var builder strings.Builder

	typeNode := node.(*ast.TypeNode)
	// 生成类型声明，检查是否有实现的接口
	if len(typeNode.Implements) > 0 {
		builder.WriteString(fmt.Sprintf("type %s implements ", typeNode.Name))
		for i, iface := range typeNode.Implements {
			if i > 0 {
				builder.WriteString(" & ")
			}
			builder.WriteString(iface)
		}
		builder.WriteString(" {\n")
	} else {
		builder.WriteString(fmt.Sprintf("type %s {\n", node.GetName()))
	}

	for _, field := range typeNode.Fields {
		builder.WriteString(fmt.Sprintf("  %s", field.Name))

		// 检查该字段是否有参数
		if len(field.Args) > 0 {
			builder.WriteString("(")
			for i, arg := range field.Args {
				if i > 0 {
					builder.WriteString(", ")
				}
				// 生成参数类型
				builder.WriteString(fmt.Sprintf("%s: %s", arg.Name, generateFieldType(arg.Type)))
				if arg.Type.IsNonNull {
					builder.WriteString("!")
				}
			}
			builder.WriteString(")")
		}

		// 生成字段类型
		builder.WriteString(fmt.Sprintf(": %s", generateFieldType(field.Type)))
		if field.Type.IsNonNull {
			builder.WriteString("!")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")
	return builder.String()
}

func generateSchema(nodes map[string]ast.Node) string {
	var schemaBuilder strings.Builder

	var scalarNodes []ast.Node

	for _, node := range nodes {
		nextLine := true
		switch node.GetNodeType() {
		case ast.NodeTypeScalar:
			// 生成 Scalar 定义
			scalarNodes = append(scalarNodes, node)
			nextLine = false
		case ast.NodeTypeInterface:
			// 生成 Interface 定义
			schemaBuilder.WriteString(generateInterfaceType(node))
		case ast.NodeTypeType:
			// 生成 Object Type 定义
			schemaBuilder.WriteString(generateObjectType(node))
		case ast.NodeTypeEnum:
			// 生成 Enum 定义
			schemaBuilder.WriteString(generateEnumType(node))
		case ast.NodeTypeInput:
			// 生成 Input 定义
			schemaBuilder.WriteString(generateInputType(node))
		case ast.NodeTypeUnion:
			// 生成 Union 定义
			schemaBuilder.WriteString(generateUnionType(node))
		default:
			nextLine = false
		}
		if nextLine {
			schemaBuilder.WriteString("\n")
		}
	}

	schemaBuilder.WriteString("\n# ============== Scalar ==============\n")
	for _, scalarNode := range scalarNodes {
		schemaBuilder.WriteString(generateScalarType(scalarNode))
	}

	return schemaBuilder.String()
}

func generateScalarType(node ast.Node) string {
	var builder strings.Builder
	scalarNode := node.(*ast.ScalarNode)
	// not generate id , int , float , string , boolean
	if scalarNode.Name == "ID" || scalarNode.Name == "Int" || scalarNode.Name == "Float" || scalarNode.Name == "String" || scalarNode.Name == "Boolean" {
		return ""
	}
	builder.WriteString(fmt.Sprintf("scalar %s\n", scalarNode.Name))
	return builder.String()
}

// 生成 Enum 类型定义
func generateEnumType(node ast.Node) string {
	var builder strings.Builder
	enumNode := node.(*ast.EnumNode)
	builder.WriteString(fmt.Sprintf("enum %s {\n", enumNode.Name))

	for _, field := range enumNode.Values {
		builder.WriteString(fmt.Sprintf("  %s\n", field.Name))
	}

	builder.WriteString("}\n")
	return builder.String()
}

// 生成 Input 类型定义
func generateInputType(node ast.Node) string {
	var builder strings.Builder
	inputNode := node.(*ast.InputNode)
	builder.WriteString(fmt.Sprintf("input %s {\n", inputNode.Name))

	for _, field := range node.GetFields() {
		builder.WriteString(fmt.Sprintf("  %s: %s", field.Name, generateFieldType(field.Type)))
		if field.Type.IsNonNull {
			builder.WriteString("!")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")
	return builder.String()
}

// 生成 Union 类型定义
func generateUnionType(node ast.Node) string {
	var builder strings.Builder
	unionNode := node.(*ast.UnionNode)
	builder.WriteString(fmt.Sprintf("union %s = ", unionNode.Name))

	for i, t := range unionNode.Types {
		if i > 0 {
			builder.WriteString(" | ")
		}
		builder.WriteString(t)
	}

	builder.WriteString("\n")
	return builder.String()
}

// 生成 Interface 类型定义
func generateInterfaceType(node ast.Node) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("interface %s {\n", node.GetName()))

	for _, field := range node.GetFields() {
		builder.WriteString(fmt.Sprintf("  %s: %s", field.Name, generateFieldType(field.Type)))
		if field.Type.IsNonNull {
			builder.WriteString("!")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")
	return builder.String()
}

// 生成字段的类型定义
func generateFieldType(fieldType *ast.FieldType) string {
	if fieldType.IsList {
		return fmt.Sprintf("[%s]", generateFieldType(fieldType.ElemType))
	}
	return fieldType.Name
}
