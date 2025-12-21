package mcp

// Version MCP 服务器版本
const Version = "1.0.0"

// StartServer 启动 MCP 服务器
func StartServer() error {
	server := NewServer("lighthouse", Version)

	// 注册所有工具
	RegisterTools(server)

	// 注册所有资源
	RegisterResources(server)

	// 运行服务器
	return server.Run()
}
