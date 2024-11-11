package manor

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/light-speak/lighthouse/graphql/ast"
	"github.com/light-speak/lighthouse/log"
)

type ManorStatus int32

const (
	ManorStatusReady   ManorStatus = 0
	ManorStatusMerging ManorStatus = 1
	ManorStatusError   ManorStatus = 2
)

var (
	mergedSchema *ast.NodeStore
	mergedMutex  sync.RWMutex
	manorStatus  ManorStatus
)

var enumMap = map[string]struct{}{
	"SortOrder":          {},
	"SearchableType":     {},
	"SearchableAnalyzer": {},
}

// initMergedSchema 初始化一个新的合并 schema
func initMergedSchema() *ast.NodeStore {
	return &ast.NodeStore{
		Objects:    make(map[string]*ast.ObjectNode),
		Interfaces: make(map[string]*ast.InterfaceNode),
		Unions:     make(map[string]*ast.UnionNode),
		Enums:      make(map[string]*ast.EnumNode),
		Inputs:     make(map[string]*ast.InputObjectNode),
		Scalars:    make(map[string]*ast.ScalarNode),
	}
}

func MergeSchemas(services map[string]*ServiceInfo) (*ast.NodeStore, error) {
	// Initialize a new schema for merging
	tempSchema := initMergedSchema()

	var errors []error
	for _, service := range services {
		err := mergeSchemaIntoTemp(service.Value, tempSchema)
		if err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("schema merge completed with %d errors", len(errors))
	}

	// Update the global mergedSchema
	mergedMutex.Lock()
	mergedSchema = tempSchema
	mergedMutex.Unlock()

	return mergedSchema, nil
}

func mergeSchemaIntoTemp(service *RegisterValue, tempSchema *ast.NodeStore) error {
	nodeStore := &ast.NodeStore{
		Objects:    make(map[string]*ast.ObjectNode),
		Interfaces: make(map[string]*ast.InterfaceNode),
		Unions:     make(map[string]*ast.UnionNode),
		Enums:      make(map[string]*ast.EnumNode),
		Inputs:     make(map[string]*ast.InputObjectNode),
		Scalars:    make(map[string]*ast.ScalarNode),
	}

	if err := sonic.UnmarshalString(service.Schema, nodeStore); err != nil {
		return err
	}

	var errors []error

	for _, obj := range nodeStore.Objects {
		if curObj, ok := tempSchema.Objects[obj.Name]; ok {
			if obj.IsExtend || obj.Name == "Query" || obj.Name == "Mutation" || obj.Name == "Subscription" {
				for _, field := range obj.Fields {
					if _, ok := curObj.Fields[field.Name]; !ok {
						field.ServiceInfo = &ast.ServiceInfo{
							Service: service.Address,
						}
						curObj.Fields[field.Name] = field
					}
				}
			} else {
				if curObj.MainService == "" {
					for _, field := range obj.Fields {
						field.ServiceInfo = &ast.ServiceInfo{
							Service: service.Address,
						}
						curObj.Fields[field.Name] = field
					}
					curObj.MainService = service.Address
					continue
				}
				if strings.HasPrefix(obj.Name, "__") {
					continue
				}
				for _, field := range obj.Fields {
					if _, ok := curObj.Fields[field.Name]; ok {
						continue
					} else {
						errors = append(errors, fmt.Errorf("field %s cannot define the same object without @extend in object %s", field.Name, obj.Name))
					}
				}
			}
		} else {
			for _, field := range obj.Fields {
				field.ServiceInfo = &ast.ServiceInfo{
					Service: service.Address,
				}
			}
			if !obj.IsExtend {
				obj.MainService = service.Address
			}
			tempSchema.Objects[obj.Name] = obj
		}
	}

	for name, intf := range nodeStore.Interfaces {
		if _, exists := tempSchema.Interfaces[name]; exists {
			if strings.HasPrefix(name, "__") {
				continue
			}
			errors = append(errors, fmt.Errorf("interface %s is already defined", name))
			continue
		}
		intf.MainService = service.Address
		tempSchema.Interfaces[name] = intf
	}

	for name, enum := range nodeStore.Enums {
		if _, exists := tempSchema.Enums[name]; exists {
			if _, ok := enumMap[name]; ok {
				continue
			}
			if strings.HasPrefix(name, "__") {
				continue
			}
			errors = append(errors, fmt.Errorf("enum %s is already defined", name))
			continue
		}
		enum.MainService = service.Address
		tempSchema.Enums[name] = enum
	}

	for name, input := range nodeStore.Inputs {
		if _, exists := tempSchema.Inputs[name]; exists {
			errors = append(errors, fmt.Errorf("input object %s is already defined", name))
			continue
		}
		input.MainService = service.Address
		tempSchema.Inputs[name] = input
	}

	for name, scalar := range nodeStore.Scalars {
		if _, exists := tempSchema.Scalars[name]; !exists {
			scalar.MainService = service.Address
			tempSchema.Scalars[name] = scalar
		}
	}

	for name, union := range nodeStore.Unions {
		if existingUnion, exists := tempSchema.Unions[name]; exists {
			for typeName, possibleType := range union.PossibleTypes {
				if _, hasType := existingUnion.PossibleTypes[typeName]; !hasType {
					existingUnion.PossibleTypes[typeName] = possibleType
				}
			}
			continue
		}
		union.MainService = service.Address
		tempSchema.Unions[name] = union
	}

	for _, err := range errors {
		log.Error().Msgf("Error merging schema: %v", err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("schema merge completed with %d errors", len(errors))
	}

	return nil
}

func GetMergedSchema() *ast.NodeStore {
	mergedMutex.RLock()
	defer mergedMutex.RUnlock()
	return mergedSchema
}

func init() {
	mergedSchema = initMergedSchema()
	manorStatus = ManorStatusMerging
}
