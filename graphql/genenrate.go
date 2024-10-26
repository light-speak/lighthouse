package graphql

import (
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

	for _, node := range nodes {
		if isInternalType(node.GetName()) {
			continue
		}
		switch node.GetKind() {
		case ast.KindObject:
			objectNode, _ := node.(*ast.ObjectNode)
			if objectNode.IsModel {
				typeNodes = append(typeNodes, objectNode)
			}
		}
	}

	if err := generate.GenObject(typeNodes, currentPath); err != nil {
		return err
	}
	if err := generate.GenInterface(p.NodeStore.Interfaces, currentPath); err != nil {
		return err
	}
	if err := generate.GenInput(p.NodeStore.Inputs, currentPath); err != nil {
		return err
	}
	if err := generate.GenRepo(typeNodes, currentPath); err != nil {
		return err
	}

	return nil
}
