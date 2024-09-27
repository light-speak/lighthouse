package command

type CommandArgType int

const (
	String CommandArgType = iota
	Int
	Bool
)

type Command interface {
	Action() func(flagValues map[string]interface{}) error
	Name() string
	Usage() string
	Args() []*CommandArg
}

type CommandArg struct {
	Name     string
	Type     CommandArgType
	Usage    string
	Required bool
}

type CommandList interface {
	GetCommands() []Command
}
