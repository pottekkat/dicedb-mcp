package tools

import (
	"context"
	"fmt"

	"github.com/pottekkat/dicedb-mcp/internal/utils"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewHGetAllTool() mcp.Tool {
	return mcp.NewTool("hgetall",
		mcp.WithDescription("Get all fields and values in the string-string map held at key"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key holding string-string map"),
		),
	)
}

func HandleHGetAllTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "HGETALL",
		Args: []string{key},
	})

	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	result := utils.FormatDiceDBResponse(resp)

	if result == "(nil)" {
		return mcp.NewToolResultText(fmt.Sprintf("Key '%s' not found or no fields present", key)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("All fields and values for key '%s': %s", key, result)), nil
}
