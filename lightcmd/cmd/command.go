package cmd

type CommandArgType int

const (
	String CommandArgType = iota
	Int
	Bool
)

var commandArgTypeNames = map[CommandArgType]string{
	String: "String",
	Int:    "Int",
	Bool:   "Bool",
}

func GetTypeName(argType CommandArgType) string {
	return commandArgTypeNames[argType]
}

type CommandInterface interface {
	Name() string
	Usage() string
	Args() []*CommandArg
	Action() func(flagValues map[string]interface{}) error
	OnExit() func()
}

type CommandArg struct {
	Name     string
	Type     CommandArgType
	Usage    string
	Required bool
	Default  interface{}
}

type CommandListInterface interface {
	GetCommands() []CommandInterface
}
