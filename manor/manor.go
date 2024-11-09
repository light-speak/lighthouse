package manor

import (
	"sync"
	"time"

	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/manor/kitex_gen/manor/rpc"
)

var serviceStore = &ServiceStore{}

type ServiceHost struct {
	// host address of the service
	Address string
	// whether the service is alive
	IsAlive bool
	// last heartbeat time
	LastHeartbeat int64
}

type ServiceStore struct {
	// store is the ast store of the service
	store *ast.NodeStore
	// services list of the service
	// serviceName -> []ServiceHost
	services sync.Map
}

func convertRPCToAST(rpcStore *rpc.NodeStore) error {
	if serviceStore.store == nil {
		store := &ast.NodeStore{}
		store.InitStore()
		serviceStore.store = store
	}

	for name, node := range rpcStore.Scalars {
		serviceStore.store.Scalars[name] = &ast.ScalarNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				Kind:        ast.KindScalar,
				IsMain:      node.IsMain,
			},
		}
	}

	for name, node := range rpcStore.Objects {
		fields := convertRPCFieldsToAST(node.Fields)
		serviceStore.store.Objects[name] = &ast.ObjectNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				Fields:      fields,
				Kind:        ast.KindObject,
				IsMain:      node.IsMain,
			},
			InterfaceNames: node.InterfaceNames,
			IsModel:        node.IsModel,
			Scopes:         node.Scopes,
			Table:          node.Table,
		}
	}

	for name, node := range rpcStore.Interfaces {
		fields := convertRPCFieldsToAST(node.Fields)
		serviceStore.store.Interfaces[name] = &ast.InterfaceNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				Fields:      fields,
				Kind:        ast.KindInterface,
				IsMain:      node.IsMain,
			},
		}
	}

	for name, node := range rpcStore.Unions {
		serviceStore.store.Unions[name] = &ast.UnionNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				Kind:        ast.KindUnion,
				IsMain:      node.IsMain,
			},
			TypeNames: node.TypeNames,
		}
	}

	for name, node := range rpcStore.Enums {
		enumValues := make(map[string]*ast.EnumValue)
		for vname, value := range node.EnumValues {
			enumValues[vname] = &ast.EnumValue{
				Name:              value.Name,
				Description:       value.Description,
				IsDeprecated:      value.IsDeprecated,
				DeprecationReason: value.DeprecationReason,
			}
		}

		serviceStore.store.Enums[name] = &ast.EnumNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				EnumValues:  enumValues,
				Kind:        ast.KindEnum,
				IsMain:      node.IsMain,
			},
		}
	}

	for name, node := range rpcStore.Inputs {
		fields := convertRPCFieldsToAST(node.Fields)
		serviceStore.store.Inputs[name] = &ast.InputObjectNode{
			BaseNode: ast.BaseNode{
				Name:        node.Name,
				Description: node.Description,
				Fields:      fields,
				Kind:        ast.KindInputObject,
				IsMain:      node.IsMain,
			},
		}
	}

	return nil
}

func convertRPCFieldsToAST(rpcFields map[string]*rpc.FieldNode) map[string]*ast.Field {
	fields := make(map[string]*ast.Field)
	for name, field := range rpcFields {
		fields[name] = &ast.Field{
			Name:              field.Name,
			Description:       field.Description,
			Args:              convertRPCArgumentsToAST(field.Args_),
			Type:              convertRPCTypeRefToAST(field.Type),
			IsDeprecated:      field.IsDeprecated,
			DeprecationReason: field.DeprecationReason,
		}
	}
	return fields
}

func convertRPCArgumentsToAST(rpcArgs map[string]*rpc.ArgumentNode) map[string]*ast.Argument {
	args := make(map[string]*ast.Argument)
	for name, arg := range rpcArgs {
		var defaultValue interface{}
		if arg.DefaultValue != nil {
			defaultValue = *arg.DefaultValue
		}
		args[name] = &ast.Argument{
			Name:         arg.Name,
			Type:         convertRPCTypeRefToAST(arg.Type),
			DefaultValue: defaultValue,
		}
	}
	return args
}

func convertRPCTypeRefToAST(rpcType *rpc.TypeRef) *ast.TypeRef {
	if rpcType == nil {
		return nil
	}
	return &ast.TypeRef{
		Kind:   ast.Kind(rpcType.Kind),
		Name:   rpcType.Name,
		OfType: convertRPCTypeRefToAST(rpcType.OfType),
	}
}

func RegisterService(serviceName string, serviceAddr string, store *rpc.NodeStore) error {
	// Try to load existing hosts first
	_, loaded := serviceStore.services.Load(serviceName)

	// Only convert RPC store to AST store if this is first time registration
	if !loaded {
		err := convertRPCToAST(store)
		if err != nil {
			return err
		}
	}

	// Use LoadOrStore to handle concurrent registrations
	for {
		// Try to load existing hosts or store new empty slice atomically
		value, loaded := serviceStore.services.LoadOrStore(serviceName, []*ServiceHost{})
		if !loaded {
			// First registration for this service
			hosts := []*ServiceHost{{
				Address:       serviceAddr,
				IsAlive:       true,
				LastHeartbeat: time.Now().Unix(),
			}}
			serviceStore.services.Store(serviceName, hosts)
			return nil
		}

		// Get current hosts
		currentHosts := value.([]*ServiceHost)

		// Check if service already registered
		exists := false
		for _, host := range currentHosts {
			if host.Address == serviceAddr {
				host.IsAlive = true
				host.LastHeartbeat = time.Now().Unix()
				exists = true
				break
			}
		}

		if !exists {
			// Create new hosts slice with additional service
			newHosts := append(currentHosts, &ServiceHost{
				Address:       serviceAddr,
				IsAlive:       true,
				LastHeartbeat: time.Now().Unix(),
			})

			// Try to update atomically
			if serviceStore.services.CompareAndSwap(serviceName, currentHosts, newHosts) {
				return nil
			}
			// If update failed, retry the whole process
			continue
		}
		return nil
	}
}
