# DiceDB MCP

A [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server implementation for DiceDB to enable interactions between AI applications (hosts/clients) and DiceDB database servers.

This implementation uses the [DiceDB Go SDK](https://github.com/DiceDB/dicedb-go) to communicate with DiceDB.

Check out the [demo video](./demo.mov) to see it in action!

## Features

- PING DiceDB to check connectivity.
- ECHO a message through DiceDB.
- GET a value from DiceDB by key.
- SET a key-value pair in DiceDB.
- DEL one or more keys from DiceDB.
- INCR the integer value of a key by one.
- DECR the integer value of a key by one.
- HGET returns the value of field present in the string-string map held at key.
- HSET sets the field and value for the key in the string-string map.
- HGETALL returns all the field-value pairs (we call it HElements) from the string-string map stored at key. 

## Installation

### Download Binary

You can [download](https://github.com/pottekkat/dicedb-mcp/releases) and use the appropriate binary for your operating system and processor archetecture from the "Releases" page.

### Install via Go

Prerequisites:

- Go 1.24 or higher

```bash
go install github.com/pottekkat/dicedb-mcp@latest
```

Get the path to the `dicedb-mcp` binary:

```bash
which dicedb-mcp
```

### Build from Source

See [Development](#development) section below.

## Usage

### With MCP Hosts/Clients

Add this to your `claude_desktop_config.json` for Claude Desktop or `mcp.json` for Cursor:

```json
{
    "mcpServers": {
        "dicedb-mcp": {
            "command": "path/to/dicedb-mcp"
        }
    }
}
```

### With OpenAI Agents SDK

The example below shows how to use the `dicedb-mcp` server with the [OpenAI Agents SDK](https://openai.github.io/openai-agents-python/):

```python
from agents import Agent, Runner, trace
from agents.mcp import MCPServer, MCPServerStdio
from dotenv import load_dotenv
import os
import openai
import asyncio

load_dotenv()


async def run(mcp_server: MCPServer, prompt: str, server_url: str):
    agent = Agent(name="DiceDB MCP",
                  instructions=f"""You can interact with a DiceDB database
                                  running at {server_url}, use
                                  this for url.""",
                  mcp_servers=[mcp_server],)
    result = await Runner.run(starting_agent=agent, input=prompt)
    print(result.final_output)


async def main():
    openai.api_key = os.getenv("OPENAI_API_KEY")

    prompt = "Can you change the value of the 'name' key to 'Rachel Green'?"
    server_url = "localhost:7379"

    async with MCPServerStdio(
        cache_tools_list=True,
        params={"command": "path/to/dicedb-mcp", "args": [""]},
    ) as server:
        with trace(workflow_name="DiceDB MCP"):
            await run(server, prompt, server_url)


if __name__ == "__main__":
    asyncio.run(main())
```

## Available Tools

### ping

Pings a DiceDB server to check connectivity.

### echo

Echoes a message through the DiceDB server.

### get

Retrieves a value from DiceDB by key.

### set

Sets a key-value pair in DiceDB.

### del

Deletes one or more keys from DiceDB.

### incr

Increments the integer value of a key by one.

### decr

Decrements the integer value of a key by one.

## Development

Fork and clone the repository:

```bash
git clone https://github.com/username/dicedb-mcp.git
```

Change into the directory:

```bash
cd dicedb-mcp
```

Install dependencies:

```bash
make deps
```

Build the project:

```bash
make build
```

Update your MCP servers configuration to point to the local build:

```json
{
    "mcpServers": {
        "dicedb-mcp": {
            "command": "/path/to/dicedb-mcp/dist/dicedb-mcp"
        }
    }
}
```

## License

[MIT License](LICENSE)
