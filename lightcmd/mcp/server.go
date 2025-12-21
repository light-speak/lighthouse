package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// Server MCP 服务器
type Server struct {
	name    string
	version string

	tools     map[string]*ToolHandler
	resources map[string]*ResourceHandler

	mu     sync.RWMutex
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// ToolHandler 工具处理器
type ToolHandler struct {
	Tool    Tool
	Handler func(args map[string]interface{}) (*CallToolResult, error)
}

// ResourceHandler 资源处理器
type ResourceHandler struct {
	Resource Resource
	Handler  func() (*ResourceContent, error)
}

// NewServer 创建新的 MCP 服务器
func NewServer(name, version string) *Server {
	return &Server{
		name:      name,
		version:   version,
		tools:     make(map[string]*ToolHandler),
		resources: make(map[string]*ResourceHandler),
		stdin:     os.Stdin,
		stdout:    os.Stdout,
		stderr:    os.Stderr,
	}
}

// RegisterTool 注册工具
func (s *Server) RegisterTool(name, description string, inputSchema map[string]interface{}, handler func(args map[string]interface{}) (*CallToolResult, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tools[name] = &ToolHandler{
		Tool: Tool{
			Name:        name,
			Description: description,
			InputSchema: inputSchema,
		},
		Handler: handler,
	}
}

// RegisterResource 注册资源
func (s *Server) RegisterResource(uri, name, description, mimeType string, handler func() (*ResourceContent, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.resources[uri] = &ResourceHandler{
		Resource: Resource{
			URI:         uri,
			Name:        name,
			Description: description,
			MimeType:    mimeType,
		},
		Handler: handler,
	}
}

// Run 运行服务器
func (s *Server) Run() error {
	scanner := bufio.NewScanner(s.stdin)
	// 增大缓冲区以处理大型请求
	const maxScanTokenSize = 1024 * 1024 // 1MB
	buf := make([]byte, maxScanTokenSize)
	scanner.Buffer(buf, maxScanTokenSize)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var req Request
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			s.sendError(nil, ParseError, "Parse error")
			continue
		}

		resp := s.handleRequest(&req)
		if resp != nil {
			s.sendResponse(resp)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %w", err)
	}

	return nil
}

// handleRequest 处理请求
func (s *Server) handleRequest(req *Request) *Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "initialized":
		// 通知，无需响应
		return nil
	case "tools/list":
		return s.handleToolsList(req)
	case "tools/call":
		return s.handleToolsCall(req)
	case "resources/list":
		return s.handleResourcesList(req)
	case "resources/read":
		return s.handleResourcesRead(req)
	case "ping":
		return NewSuccessResponse(req.ID, map[string]interface{}{})
	default:
		return NewErrorResponse(req.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", req.Method))
	}
}

// handleInitialize 处理初始化请求
func (s *Server) handleInitialize(req *Request) *Response {
	result := InitializeResult{
		ProtocolVersion: MCPVersion,
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    s.name,
			Version: s.version,
		},
	}

	return NewSuccessResponse(req.ID, result)
}

// handleToolsList 处理工具列表请求
func (s *Server) handleToolsList(req *Request) *Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tools := make([]Tool, 0, len(s.tools))
	for _, th := range s.tools {
		tools = append(tools, th.Tool)
	}

	return NewSuccessResponse(req.ID, ToolsListResult{Tools: tools})
}

// handleToolsCall 处理工具调用请求
func (s *Server) handleToolsCall(req *Request) *Response {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, InvalidParams, "Invalid params")
	}

	s.mu.RLock()
	handler, ok := s.tools[params.Name]
	s.mu.RUnlock()

	if !ok {
		return NewErrorResponse(req.ID, InvalidParams, fmt.Sprintf("Tool not found: %s", params.Name))
	}

	result, err := handler.Handler(params.Arguments)
	if err != nil {
		return NewSuccessResponse(req.ID, CallToolResult{
			Content: []ContentBlock{NewTextContent(fmt.Sprintf("Error: %s", err.Error()))},
			IsError: true,
		})
	}

	return NewSuccessResponse(req.ID, result)
}

// handleResourcesList 处理资源列表请求
func (s *Server) handleResourcesList(req *Request) *Response {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resources := make([]Resource, 0, len(s.resources))
	for _, rh := range s.resources {
		resources = append(resources, rh.Resource)
	}

	return NewSuccessResponse(req.ID, ResourcesListResult{Resources: resources})
}

// handleResourcesRead 处理资源读取请求
func (s *Server) handleResourcesRead(req *Request) *Response {
	var params ReadResourceParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return NewErrorResponse(req.ID, InvalidParams, "Invalid params")
	}

	s.mu.RLock()
	handler, ok := s.resources[params.URI]
	s.mu.RUnlock()

	if !ok {
		return NewErrorResponse(req.ID, InvalidParams, fmt.Sprintf("Resource not found: %s", params.URI))
	}

	content, err := handler.Handler()
	if err != nil {
		return NewErrorResponse(req.ID, InternalError, err.Error())
	}

	return NewSuccessResponse(req.ID, ReadResourceResult{
		Contents: []ResourceContent{*content},
	})
}

// sendResponse 发送响应
func (s *Server) sendResponse(resp *Response) {
	data, err := json.Marshal(resp)
	if err != nil {
		s.logError("Failed to marshal response: %v", err)
		return
	}
	fmt.Fprintln(s.stdout, string(data))
}

// sendError 发送错误响应
func (s *Server) sendError(id json.RawMessage, code int, message string) {
	resp := NewErrorResponse(id, code, message)
	s.sendResponse(resp)
}

// logError 记录错误到 stderr
func (s *Server) logError(format string, args ...interface{}) {
	fmt.Fprintf(s.stderr, format+"\n", args...)
}
