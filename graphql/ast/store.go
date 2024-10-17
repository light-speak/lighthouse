package ast

type NodeStore struct {
	// Scalars is a map of all scalars
	Scalars map[string]*ScalarNode

	// Interfaces is a map of all interfaces
	Interfaces map[string]*InterfaceNode

	// Objects is a map of all objects
	Objects map[string]*ObjectNode

	// Unions is a map of all unions
	Unions map[string]*UnionNode

	// Enums is a map of all enums
	Enums map[string]*EnumNode

	// Inputs is a map of all inputs
	Inputs map[string]*InputObjectNode

	// Directives is a map of all directives
	Directives map[string]*DirectiveDefinition

	// ScalarTypes is a map of all scalar types
	ScalarTypes map[string]ScalarType

	// Names is a map of all names , it can be a type, enum, interface, input, scalar, union, directive, extend
	Names map[string]any

	// Nodes is a map of all nodes
	Nodes map[string]Node
}

func (s *NodeStore) InitStore() {
	s.Names = make(map[string]any)
	s.Objects = make(map[string]*ObjectNode)
	s.Unions = make(map[string]*UnionNode)
	s.Enums = make(map[string]*EnumNode)
	s.Interfaces = make(map[string]*InterfaceNode)
	s.Inputs = make(map[string]*InputObjectNode)
	s.Scalars = make(map[string]*ScalarNode)
	s.Directives = make(map[string]*DirectiveDefinition)
	s.Nodes = make(map[string]Node)
}

func (s *NodeStore) AddObject(name string, node *ObjectNode) {
	s.Objects[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddScalar(name string, node *ScalarNode) {
	s.Scalars[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddInterface(name string, node *InterfaceNode) {
	s.Interfaces[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddUnion(name string, node *UnionNode) {
	s.Unions[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddEnum(name string, node *EnumNode) {
	s.Enums[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddInput(name string, node *InputObjectNode) {
	s.Inputs[name] = node
	s.Nodes[name] = node
	s.Names[name] = node
}

func (s *NodeStore) AddDirective(name string, node *DirectiveDefinition) {
	s.Directives[name] = node
	s.Names[name] = node
}

