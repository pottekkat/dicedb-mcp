package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/pottekkat/dicedb-mcp/internal/tools"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"DiceDB MCP",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	// Create and add the ping tool
	pingTool := tools.NewPingTool()
	s.AddTool(pingTool, tools.HandlePingTool)

	echoTool := tools.NewEchoTool()
	s.AddTool(echoTool, tools.HandleEchoTool)

	// Create and add the get tool
	getTool := tools.NewGetTool()
	s.AddTool(getTool, tools.HandleGetTool)

	// Create and add the set tool
	setTool := tools.NewSetTool()
	s.AddTool(setTool, tools.HandleSetTool)

	// Create and add the del tool
	delTool := tools.NewDelTool()
	s.AddTool(delTool, tools.HandleDelTool)

	// Create and add the incr tool
	incrTool := tools.NewIncrTool()
	s.AddTool(incrTool, tools.HandleIncrTool)

	hgetTool := tools.NewHGetTool()
	s.AddTool(hgetTool, tools.HandleHGetTool)

	hgetall := tools.NewHGetAllTool()
	s.AddTool(hgetall, tools.HandleHGetAllTool)

	hsetTool := tools.NewHSetTool()
	s.AddTool(hsetTool, tools.HandleHSetTool)

	// Start the server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error starting server: %v\n", err)
	}
}
