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

	if err := generate.GenOperationResolverGen(operationNodes, currentPath); err != nil {
		return err
	}

	return nil
}
