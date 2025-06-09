package tools

import (
	"context"
	"fmt"

	"github.com/pottekkat/dicedb-mcp/internal/utils"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewHGetTool() mcp.Tool {
	return mcp.NewTool("hget",
		mcp.WithDescription("Get a value of field present in the string-string map held at key"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key holding string-string map"),
		),
		mcp.WithString("field",
			mcp.Required(),
			mcp.Description("The field to retrieve from DiceDB"),
		),
	)
}

func HandleHGetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}
	field, ok := request.GetArguments()["field"].(string)
	if !ok || field == "" {
		return nil, fmt.Errorf("missing or empty field parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "HGET",
		Args: []string{key, field},
	})

	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	value := utils.FormatDiceDBResponse(resp)

	if value == "(nil)" {
		return mcp.NewToolResultText(fmt.Sprintf("Key '%s' not found or Field '%s' not found", key, field)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Value for key and field '%s', '%s', '%s'", key, field, value)), nil
}
