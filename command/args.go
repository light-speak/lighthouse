package command

import "fmt"

func GetArgs(argDefs []*CommandArg, flagValues map[string]interface{}) (map[string]interface{}, error) {
	args := make(map[string]interface{})

	for _, argDef := range argDefs {
		value, exists := flagValues[argDef.Name]
		if !exists || value == nil {
			if argDef.Required {
				return nil, fmt.Errorf("%s is required", argDef.Name)
			}
			// Set default value for optional args
			args[argDef.Name] = GetDefaultValue(argDef)
		} else {
			args[argDef.Name] = value
		}
	}

	return args, nil
}

// GetStringArg gets a string argument from the flag values
func GetStringArg(args map[string]interface{}, name string) (*string, error) {
	value, exists := args[name]
	if !exists || value == nil {
		return nil, fmt.Errorf("%s is required", name)
	}
	return value.(*string), nil
}

// GetIntArg gets an int argument from the flag values
func GetIntArg(args map[string]interface{}, name string) (*int, error) {
	value, exists := args[name]
	if !exists || value == nil {
		return nil, fmt.Errorf("%s is required", name)
	}
	return value.(*int), nil
}

// GetBoolArg gets a bool argument from the flag values
func GetBoolArg(args map[string]interface{}, name string) (*bool, error) {
	value, exists := args[name]
	if !exists || value == nil {
		return nil, fmt.Errorf("%s is required", name)
	}
	return value.(*bool), nil
}

func GetDefaultValue(argDef *CommandArg) interface{} {
	switch argDef.Type {
	case String:
		return argDef.Default.(string)
	case Int:
		return argDef.Default.(int)
	case Bool:
		return argDef.Default.(bool)
	default:
		return nil
	}
}
