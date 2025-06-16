package tools

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dicedb/dicedb-go/wire"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/pottekkat/dicedb-mcp/internal/utils"
)

// NewSetTool creates a new SET tool for DiceDB
func NewSetTool() mcp.Tool {
	return mcp.NewTool("set",
		mcp.WithDescription("Set a key-value pair in DiceDB"),
		utils.CommonURLParam(),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("The key to set"),
		),
		mcp.WithString("value",
			mcp.Required(),
			mcp.Description("The value to set"),
		),
		mcp.WithNumber("ex",
			mcp.Description("Set the expiration time in seconds"),
		),
		mcp.WithNumber("px",
			mcp.Description("Set the expiration time in milliseconds"),
		),
		mcp.WithNumber("exat",
			mcp.Description("Set the expiration time in seconds since epoch"),
		),
		mcp.WithNumber("pxat",
			mcp.Description("Set the expiration time in milliseconds since epoch"),
		),
		mcp.WithBoolean("xx",
			mcp.Description("Only set the key if it already exists"),
		),
		mcp.WithBoolean("nx",
			mcp.Description("Only set the key if it does not already exist"),
		),
		mcp.WithBoolean("keepttl",
			mcp.Description("Keep the existing TTL of the key"),
		),
		mcp.WithBoolean("get",
			mcp.Description("Return the value of the key after setting it"),
		),
	)
}

// HandleSetTool handles the SET tool request
func HandleSetTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, ok := request.GetArguments()["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("missing or empty key parameter")
	}

	value, ok := request.GetArguments()["value"].(string)
	if !ok {
		return nil, fmt.Errorf("missing value parameter")
	}

	args := []string{key, value}

	if ex, ok := request.GetArguments()["ex"].(float64); ok {
		args = append(args, "EX", strconv.FormatInt(int64(ex), 10))
	}

	if px, ok := request.GetArguments()["px"].(float64); ok {
		args = append(args, "PX", strconv.FormatInt(int64(px), 10))
	}

	if exat, ok := request.GetArguments()["exat"].(float64); ok {
		args = append(args, "EXAT", strconv.FormatInt(int64(exat), 10))
	}

	if pxat, ok := request.GetArguments()["pxat"].(float64); ok {
		args = append(args, "PXAT", strconv.FormatInt(int64(pxat), 10))
	}

	if xx, ok := request.GetArguments()["xx"].(bool); ok && xx {
		args = append(args, "XX")
	}

	if nx, ok := request.GetArguments()["nx"].(bool); ok && nx {
		args = append(args, "NX")
	}

	if keepttl, ok := request.GetArguments()["keepttl"].(bool); ok && keepttl {
		args = append(args, "KEEPTTL")
	}

	if get, ok := request.GetArguments()["get"].(bool); ok && get {
		args = append(args, "GET")
	}

	client, err := utils.GetClientFromRequest(request)

	if err != nil {
		return nil, err
	}

	resp := client.Fire(&wire.Command{
		Cmd:  "SET",
		Args: args,
	})

	if resp.Status == wire.Status_ERR {
		return nil, fmt.Errorf("DiceDB error: %s", resp.Message)
	}

	formattedResponse := utils.FormatDiceDBResponse(resp)

	var resultMessage string
	if formattedResponse == "OK" {
		resultMessage = fmt.Sprintf("Successfully set key '%s'", key)
	} else if formattedResponse == "(nil)" {
		resultMessage = fmt.Sprintf("Key '%s' was not set (condition not met)", key)
	} else {
		// Display the value of the key also
		resultMessage = fmt.Sprintf("Set key '%s' with value: %s", key, formattedResponse)
	}

	return mcp.NewToolResultText(resultMessage), nil
}
