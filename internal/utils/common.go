// pkg/utils/common.go
package utils

import (
	"fmt"

	"github.com/dicedb/dicedb-go"
	"github.com/mark3labs/mcp-go/mcp"
)

// CommonURLParam returns the common URL parameter definition used across tools
func CommonURLParam() mcp.ToolOption {
	return mcp.WithString("url",
		mcp.Description("The URL of the DiceDB server in format 'host:port'"),
		mcp.DefaultString("localhost:7379"),
	)
}

// GetClientFromRequest creates a DiceDB client from the request parameters
func GetClientFromRequest(request mcp.CallToolRequest) (*dicedb.Client, error) {
	// Get the URL with fallback to default
	var url string = "localhost:7379"
	if urlArg, ok := request.GetArguments()["url"]; ok && urlArg != nil {
		if urlStr, ok := urlArg.(string); ok && urlStr != "" {
			url = urlStr
		}
	}

	host, port := parseHostAndPort(url)

	// Create a new DiceDB client
	client, err := dicedb.NewClient(host, port)
	if err != nil {
		return nil, fmt.Errorf("error connecting to DiceDB: %w", err)
	}

	return client, nil
}
