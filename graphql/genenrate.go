package graphql

import (
	"errors"
	"os"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model/generate"
	"github.com/light-speak/lighthouse/log"
)

func Generate() error {
	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	err = LoadSchema()
	if err != nil {
		log.Error().Msgf("Failed to load schema: %v", err)
		return err
	}
	p := GetParser()
	nodes := p.NodeStore.Nodes

	typeNodes := []*ast.ObjectNode{}
	resNodes := []*ast.ObjectNode{}

	for _, node := range nodes {
		if isInternalType(node.GetName()) {
			continue
		}
		switch node.GetKind() {
		case ast.KindObject:
			objectNode, _ := node.(*ast.ObjectNode)
			if objectNode.IsModel {
				typeNodes = append(typeNodes, objectNode)
			} else {
				resNodes = append(resNodes, objectNode)
			}
		}
	}

	interfaces := []ast.Node{}
	for _, interfaceNode := range p.NodeStore.Interfaces {
		interfaces = append(interfaces, interfaceNode)
	}
	for _, union := range p.NodeStore.Unions {
		interfaces = append(interfaces, union)
	}

	if err := generate.GenObject(typeNodes, currentPath); err != nil {
		return err
	}
	if err := generate.GenInterface(interfaces, currentPath); err != nil {
		return err
	}
	if err := generate.GenInput(p.NodeStore.Inputs, currentPath); err != nil {
		return err
	}
	if err := generate.GenRepo(typeNodes, currentPath); err != nil {
		return err
	}
	if err := generate.GenEnum(p.NodeStore.Enums, currentPath); err != nil {
		return err
	}
	if err := generate.GenResponse(resNodes, currentPath); err != nil {
		return err
	}

	operationNodes := []*ast.ObjectNode{}

	if _, ok := p.NodeStore.Objects["Query"]; ok {
		operationNodes = append(operationNodes, p.NodeStore.Objects["Query"])
		if err := generate.GenOperationResolver(p.NodeStore.Objects["Query"], currentPath, "query"); err != nil {
			return err
		}
	}
	if _, ok := p.NodeStore.Objects["Mutation"]; ok {
		operationNodes = append(operationNodes, p.NodeStore.Objects["Mutation"])
		if err := generate.GenOperationResolver(p.NodeStore.Objects["Mutation"], currentPath, "mutation"); err != nil {
			return err
		}
	}
	attrNodes := map[string][]*ast.Field{}
	fieldNameMap := make(map[string]string) // field name -> type name mapping

	searchableFields := map[*ast.ObjectNode][]*generate.SearchableField{}
	for _, node := range p.NodeStore.Objects {
		if node.IsModel {
			for _, field := range node.Fields {
				if field.IsSearchable {
					for _, directive := range field.Directives {
						if directive.Name == "searchable" {
							indexAnalyzer := directive.Definition.Args["indexAnalyzer"].GetDefaultValue()
							searchAnalyzer := directive.Definition.Args["searchAnalyzer"].GetDefaultValue()

							if directive.GetArg("indexAnalyzer") != nil {
								indexAnalyzer = directive.GetArg("indexAnalyzer").Value.(*string)
							}
							if directive.GetArg("searchAnalyzer") != nil {
								searchAnalyzer = directive.GetArg("searchAnalyzer").Value.(*string)
							}

							searchableFields[node] = append(searchableFields[node], &generate.SearchableField{
								Field:          field,
								Type:           directive.GetArg("type").Value.(string),
								IndexAnalyzer:  indexAnalyzer,
								SearchAnalyzer: searchAnalyzer,
							})
						}
					}
				}
				if field.IsAttr {
					if existingType, exists := fieldNameMap[field.Name]; exists {
						return errors.New("duplicate attr field name '" + field.Name + "' found in types '" + existingType + "' and '" + node.GetName() + "'")
					}
					fieldNameMap[field.Name] = node.GetName()

					if _, ok := attrNodes[node.GetName()]; !ok {
						attrNodes[node.GetName()] = []*ast.Field{}
					}
					attrNodes[node.GetName()] = append(attrNodes[node.GetName()], field)
				}
			}
		}
	}

	if err := generate.GenAttr(attrNodes, currentPath); err != nil {
		return err
	}

	if err := generate.GenSearchable(searchableFields, currentPath); err != nil {
		return err
	}

	if err := generate.GenOperationResolverGen(operationNodes, currentPath); err != nil {
		return err
	}

	return nil
}
