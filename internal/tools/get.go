package tools

import (
	"context"
	"fmt"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/internal/utils"
)

// NewGetTool creates a new GET tool for DiceDB
func NewGetTool() mcp.Tool {
	return mcp.NewTool("get",
		mcp.WithDescription("Get a value from DiceDB by key"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key to retrieve from DiceDB"),
		),
	)
}

// HandleGetTool handles the GET tool request
func HandleGetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "GET",
		Args: []string{key},
	})

	// Check if DiceDB returned an error
	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	value := utils.FormatDiceDBResponse(resp)

	// If value is "(nil)", the key doesn't exist
	if value == "(nil)" {
		return mcp.NewToolResultText(fmt.Sprintf("Key '%s' not found", key)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Value for key '%s': %s", key, value)), nil
}
