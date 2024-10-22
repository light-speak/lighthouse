package graphql

func Generate() error {

	// p := GetParser()
	// nodes := p.NodeStore.Nodes

	// typeNodes := []*ast.ObjectNode{}
	// responseNodes := []*ast.ObjectNode{}

	// for _, node := range nodes {
	// 	switch node.GetKind() {
	// 	case ast.KindObject:
	// 		typeNode, _ := node.(*ast.ObjectNode)
	// 		if typeNode.IsResponse {
	// 			responseNodes = append(responseNodes, typeNode)
	// 		} else {
	// 			typeNodes = append(typeNodes, typeNode)
	// 		}
	// 	}
	// }

	// if err := generate.GenType(typeNodes, currentPath); err != nil {
	// 	return err
	// }
	// if err := generate.GenResponse(responseNodes, currentPath); err != nil {
	// 	return err
	// }
	// if err := generate.GenInterface(p.InterfaceMap, currentPath); err != nil {
	// 	return err
	// }
	// if err := generate.GenInput(p.InputMap, currentPath); err != nil {
	// 	return err
	// }

	// schema := generateSchema(nodes)
	// options := &template.Options{
	// 	Path:         currentPath,
	// 	Template:     schema,
	// 	FileName:     "schema",
	// 	FileExt:      "graphql",
	// 	Editable:     false,
	// 	SkipIfExists: false,
	// }
	// template.Render(options)

	return nil
}
