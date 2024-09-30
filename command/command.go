package command

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

type Command interface {
	Name() string
	Usage() string
	Args() []*CommandArg
	Action() func(flagValues map[string]interface{}) error
}

type CommandArg struct {
	Name     string
	Type     CommandArgType
	Usage    string
	Required bool
	Default  interface{}
}

type CommandList interface {
	GetCommands() []Command
}
