package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/pottekkat/dicedb-mcp/internal/utils"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewHSetTool() mcp.Tool {
	return mcp.NewTool("hset",
		mcp.WithDescription("Set a value for a field in the string-string map held at key"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key holding string-string map"),
		),
		mcp.WithArray("pairs",
			mcp.Items(map[string]any{"type": "string"}),
			mcp.Required(),
			mcp.Description("The string string map to set in DiceDB"),
		),
	)
}

func HandleHSetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	kvs, ok := request.GetArguments()["pairs"].([]any)

	// Convert the keys to strings
	stringKvs := make([]string, len(kvs))
	for i, key := range kvs {
		stringKvs[i] = fmt.Sprintf("%v", key)
	}
	if len(kvs) == 0 || (len(kvs)&1) == 1 || !ok {
		return nil, fmt.Errorf("missing or empty pairs parameter")
	}

	client, err := utils.GetClientFromRequest(request)
	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "HSET",
		Args: append([]string{key}, stringKvs...),
	})

	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	value := utils.FormatDiceDBResponse(resp)

	var resultMessage string
	if strings.HasPrefix(value, "OK") {
		resultMessage = fmt.Sprintf("Key '%s' is set to string-string map '%v'", key, kvs)
	} else {
		resultMessage = fmt.Sprintf("Key '%s' was not set", key)
	}

	return mcp.NewToolResultText(resultMessage), nil
}
