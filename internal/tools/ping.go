package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/internal/utils"
)

// NewPingTool creates a new ping tool for DiceDB
func NewPingTool() mcp.Tool {
	return mcp.NewTool("ping",
		mcp.WithDescription("Ping the DiceDB server to check connectivity"),
		// All tools have a url argument
		utils.CommonURLParam(),
	)
}

// HandlePingTool handles the ping tool request
func HandlePingTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	// Run the PING command on the DiceDB server
	resp := client.Fire(&wire.Command{Cmd: "PING"})

	// Check if DiceDB returned an error
	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	// Return the response to the MCP client
	return mcp.NewToolResultText(fmt.Sprintf("Response from DiceDB: %s", utils.FormatDiceDBResponse(resp))), nil
}
