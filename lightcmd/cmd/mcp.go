package cmd

import (
	"github.com/light-speak/lighthouse/lightcmd/mcp"
)

// MCPCommand MCP 服务器命令
type MCPCommand struct{}

func (c *MCPCommand) Name() string {
	return "mcp"
}

func (c *MCPCommand) Usage() string {
	return "Start the MCP (Model Context Protocol) server for AI integration"
}

func (c *MCPCommand) Args() []*CommandArg {
	return []*CommandArg{}
}

func (c *MCPCommand) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		return mcp.StartServer()
	}
}

func (c *MCPCommand) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&MCPCommand{})
}
