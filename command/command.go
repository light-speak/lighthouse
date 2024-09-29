package command

type CommandArgType int

const (
	String CommandArgType = iota
	Int
	Bool
)

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
