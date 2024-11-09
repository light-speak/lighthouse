package manor

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/manor/kitex_gen/manor/rpc"
	"github.com/light-speak/lighthouse/manor/kitex_gen/manor/rpc/manor"
)

var (
	manorClient manor.Client
	serviceAddr string
	serviceName string
)

// InitManorClient initializes the manor client
func InitManorClient(addr string, name string) error {
	c, err := manor.NewClient("manor.rpc", client.WithHostPorts(addr))
	if err != nil {
		return fmt.Errorf("failed to create manor client: %w", err)
	}

	manorClient = c
	serviceAddr = addr
	serviceName = name
	return nil
}

// convertTypeRef 转换类型引用
func convertTypeRef(t *ast.TypeRef) *rpc.TypeRef {
	if t == nil {
		return nil
	}
	return &rpc.TypeRef{
		Kind:   string(t.Kind),
		Name:   t.Name,
		OfType: convertTypeRef(t.OfType),
	}
}

// convertArgument 转换参数
func convertArgument(a *ast.Argument) *rpc.ArgumentNode {
	return &rpc.ArgumentNode{
		Name: a.Name,
		Type: convertTypeRef(a.Type),
	}
}

// convertField 转换字段
func convertField(f *ast.Field) *rpc.FieldNode {
	args := make(map[string]*rpc.ArgumentNode)
	for name, arg := range f.Args {
		args[name] = convertArgument(arg)
	}

	return &rpc.FieldNode{
		Name:              f.Name,
		Description:       f.Description,
		Args_:             args,
		Type:              convertTypeRef(f.Type),
		IsDeprecated:      f.IsDeprecated,
		DeprecationReason: f.DeprecationReason,
	}
}

// convertStore 转换整个存储
func convertStore(store *ast.NodeStore) *rpc.NodeStore {
	rpcStore := &rpc.NodeStore{
		Scalars:    make(map[string]*rpc.ScalarNode),
		Interfaces: make(map[string]*rpc.InterfaceNode),
		Objects:    make(map[string]*rpc.ObjectNode),
		Unions:     make(map[string]*rpc.UnionNode),
		Enums:      make(map[string]*rpc.EnumNode),
		Inputs:     make(map[string]*rpc.InputObjectNode),
	}

	// Convert Scalars
	for name, node := range store.Scalars {
		rpcStore.Scalars[name] = &rpc.ScalarNode{
			Name:        node.Name,
			Description: node.Description,
			IsMain:      node.IsMain,
		}
	}

	// Convert Objects
	for name, node := range store.Objects {
		fields := make(map[string]*rpc.FieldNode)
		for fname, field := range node.Fields {
			fields[fname] = convertField(field)
		}

		rpcStore.Objects[name] = &rpc.ObjectNode{
			Name:           node.Name,
			Description:    node.Description,
			Fields:         fields,
			InterfaceNames: node.InterfaceNames,
			IsModel:        node.IsModel,
			Scopes:         node.Scopes,
			Table:          node.Table,
			IsMain:         node.IsMain,
		}
	}

	// Convert Interfaces
	for name, node := range store.Interfaces {
		fields := make(map[string]*rpc.FieldNode)
		for fname, field := range node.Fields {
			fields[fname] = convertField(field)
		}

		rpcStore.Interfaces[name] = &rpc.InterfaceNode{
			Name:        node.Name,
			Description: node.Description,
			Fields:      fields,
			IsMain:      node.IsMain,
		}
	}

	// Convert Unions
	for name, node := range store.Unions {
		rpcStore.Unions[name] = &rpc.UnionNode{
			Name:        node.Name,
			Description: node.Description,
			TypeNames:   node.TypeNames,
			IsMain:      node.IsMain,
		}
	}

	// Convert Enums
	for name, node := range store.Enums {
		values := make(map[string]*rpc.EnumValueNode)
		for vname, value := range node.EnumValues {
			values[vname] = &rpc.EnumValueNode{
				Name:              value.Name,
				Description:       value.Description,
				IsDeprecated:      value.IsDeprecated,
				DeprecationReason: value.DeprecationReason,
			}
		}

		rpcStore.Enums[name] = &rpc.EnumNode{
			Name:        node.Name,
			Description: node.Description,
			EnumValues:  values,
			IsMain:      node.IsMain,
		}
	}

	// Convert Input Objects
	for name, node := range store.Inputs {
		fields := make(map[string]*rpc.FieldNode)
		for fname, field := range node.Fields {
			fields[fname] = convertField(field)
		}

		rpcStore.Inputs[name] = &rpc.InputObjectNode{
			Name:        node.Name,
			Description: node.Description,
			Fields:      fields,
			IsMain:      node.IsMain,
		}
	}

	return rpcStore
}

// Register 注册服务
func Register(store *ast.NodeStore) error {
	if manorClient == nil {
		return fmt.Errorf("manor client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &rpc.RegisterRequest{
		ServiceName: serviceName,
		ServiceAddr: serviceAddr,
		Store:       convertStore(store),
	}

	resp, err := manorClient.Register(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("register failed: %s", resp.Message)
	}

	return nil
}
