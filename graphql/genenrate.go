package graphql

import (
	"os"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/graphql/model/generate"
	"github.com/light-speak/lighthouse/log"
)

func Generate() error {
	err := LoadSchema()
	if err != nil {
		log.Error().Msgf("Failed to load schema: %v", err)
		return err
	}
	p := GetParser()
	nodes := p.NodeStore.Nodes

	typeNodes := []*ast.ObjectNode{}
	responseNodes := []*ast.ObjectNode{}

	for _, node := range nodes {
		switch node.GetKind() {
		case ast.KindObject:
			objectNode, _ := node.(*ast.ObjectNode)
			if !objectNode.IsModel {
				responseNodes = append(responseNodes, objectNode)
			} else {
				typeNodes = append(typeNodes, objectNode)
			}
		}
	}

	currentPath, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := generate.GenObject(typeNodes, currentPath); err != nil {
		return err
	}
	// if err := generate.GenResponse(responseNodes, currentPath); err != nil {
	// 	return err
	// }
	// if err := generate.GenInterface(p.NodeStore.Interfaces, currentPath); err != nil {
	// 	return err
	// }
	// if err := generate.GenInput(p.NodeStore.Inputs, currentPath); err != nil {
	// 	return err
	// }

	return nil
}
