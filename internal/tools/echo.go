package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/internal/utils"
)

// NewEchoTool creates a new echo tool for DiceDB
func NewEchoTool() mcp.Tool {
	return mcp.NewTool("echo",
		mcp.WithDescription("Echo a message through the DiceDB server"),
		utils.CommonURLParam(),
		mcp.WithString("message",
			mcp.Required(),
			mcp.Description("The message to echo"),
		),
	)
}

// HandleEchoTool handles the echo tool request
func HandleEchoTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, ok := request.GetArguments()["message"].(string)
	if !ok || message == "" {
		return nil, fmt.Errorf("missing or empty message parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "ECHO",
		Args: []string{message},
	})

	// Check if DiceDB returned an error
	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	return mcp.NewToolResultText(fmt.Sprintf("DiceDB echoed: %s", utils.FormatDiceDBResponse(resp))), nil
}
