package validate

import "github.com/light-speak/lighthouse/graphql/ast"

func validateOperation(node ast.Node) error {
	//TODO: validate operation
	// Step1 use node.GetFields() to get the fields of the operation
	// Step2 validate each field Type is valid
	//       type should be in Query or Mutation or Subscription
	// Step3 validate the field return type is valid
	//       validate the children of the field is matching the return type
	// Step4 recursion...
	return nil
}
