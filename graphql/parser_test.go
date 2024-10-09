package graphql

import (
	"fmt"
	"strings"
	"testing"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/parser"
	"github.com/light-speak/lighthouse/graphql/parser/lexer"
	"github.com/light-speak/lighthouse/graphql/validate"
	"github.com/light-speak/lighthouse/log"
)

func TestReadGraphQLFile(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	for {
		token := l.NextToken()
		log.Debug().Msgf("%+v", token.Value)
		if token.Type == lexer.EOF {
			break
		}
	}
}

func TestParseSchema(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser(l)
	nodes := p.ParseSchema()
	for _, node := range nodes {
		log.Debug().Msgf("Type: %s", node.GetNodeType())
	}
}

func TestValidate(t *testing.T) {
	l, err := parser.ReadGraphQLFile("demo.graphql")
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser(l)

	nodes := p.ParseSchema()
	for _, node := range nodes {
		err := validate.Validate(node, p)
		if err != nil {
			t.Fatal(err)
		}
	}

	schemaNodes := make([]ast.Node, 0, len(nodes))
	for _, node := range nodes {
		schemaNodes = append(schemaNodes, node)
	}
	schema := generateSchema(schemaNodes)
	log.Debug().Msgf("schema: %s", schema)
}

func generateObjectType(node ast.Node) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("type %s {\n", node.GetName()))

	for _, field := range node.GetFields() {
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

func generateSchema(nodes []ast.Node) string {
	var schemaBuilder strings.Builder
	for _, node := range nodes {
		switch node.GetNodeType() {
		case ast.NodeTypeType:
			// 生成 Object Type 定义
			schemaBuilder.WriteString(generateObjectType(node))
		case ast.NodeTypeInterface:
			// 生成 Interface 定义
			schemaBuilder.WriteString(generateInterfaceType(node))
		case ast.NodeTypeEnum:
			// 生成 Enum 定义
			schemaBuilder.WriteString(generateEnumType(node))
		case ast.NodeTypeInput:
			// 生成 Input 定义
			schemaBuilder.WriteString(generateInputType(node))
		case ast.NodeTypeUnion:
			// 生成 Union 定义
			schemaBuilder.WriteString(generateUnionType(node))
		}
		schemaBuilder.WriteString("\n")
	}

	return schemaBuilder.String()
}

// 生成 Enum 类型定义
func generateEnumType(node ast.Node) string {
	var builder strings.Builder
	enumNode := node.(*ast.EnumNode)
	builder.WriteString(fmt.Sprintf("enum %s {\n", enumNode.Name))

	for _, field := range enumNode.Values {
		log.Debug().Msgf("enum value: %+v", field.Value)
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
