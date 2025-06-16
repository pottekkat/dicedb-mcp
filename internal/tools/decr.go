package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/internal/utils"
)

// NewDecrTool creates a new DECR tool for DiceDB
func NewDecrTool() mcp.Tool {
	return mcp.NewTool("decr",
		mcp.WithDescription("Decrement the integer value of a key by one"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key whose value to decrement"),
		),
	)
}

// HandleDecrTool handles the DECR tool request
func HandleDecrTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "DECR",
		Args: []string{key},
	})

	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	value := utils.FormatDiceDBResponse(resp)
	return mcp.NewToolResultText(fmt.Sprintf("Decremented key '%s', new value: %s", key, value)), nil
}
