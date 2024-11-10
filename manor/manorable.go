package manor

import (
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/env"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
	etcd "go.etcd.io/etcd/client/v3"
)

func Manorable(nodeStore *ast.NodeStore) {
	client, err := etcd.New(etcd.Config{
		Endpoints:   env.LighthouseConfig.Etcd.Endpoints,
		Username:    env.LighthouseConfig.Etcd.Username,
		Password:    env.LighthouseConfig.Etcd.Password,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Error().Msgf("Failed to connect to etcd: %v", err)
		return
	}
	defer client.Close()

	register(nodeStore)
}

func register(nodeStore *ast.NodeStore) error {
	// Filter out nodes starting with __
	filteredStore := &ast.NodeStore{
		Objects:    make(map[string]*ast.ObjectNode),
		Interfaces: make(map[string]*ast.InterfaceNode),
		Unions:     make(map[string]*ast.UnionNode),
		Enums:      make(map[string]*ast.EnumNode),
		Inputs:     make(map[string]*ast.InputObjectNode),
		Scalars:    make(map[string]*ast.ScalarNode),
	}

	// Copy objects excluding __ prefixed ones
	for name, obj := range nodeStore.Objects {
		if !strings.HasPrefix(name, "__") {
			for _, field := range obj.Fields {
				if !strings.HasPrefix(field.Name, "__") || field.Type.Name == "__typename" {
					filteredStore.Objects[name] = obj
				}
			}
		}
	}

	// Copy other node types excluding __ prefixed ones
	for name, node := range nodeStore.Interfaces {
		if !strings.HasPrefix(name, "__") {
			filteredStore.Interfaces[name] = node
		}
	}

	for name, node := range nodeStore.Unions {
		if !strings.HasPrefix(name, "__") {
			filteredStore.Unions[name] = node
		}
	}

	for name, node := range nodeStore.Enums {
		if !strings.HasPrefix(name, "__") {
			filteredStore.Enums[name] = node
		}
	}

	for name, node := range nodeStore.Inputs {
		if !strings.HasPrefix(name, "__") {
			filteredStore.Inputs[name] = node
		}
	}

	for name, node := range nodeStore.Scalars {
		if !strings.HasPrefix(name, "__") {
			filteredStore.Scalars[name] = node
		}
	}

	jsonData, err := sonic.Marshal(filteredStore)
	if err != nil {
		return err
	}
	log.Info().Msgf("%v", string(jsonData))

	return nil
}
