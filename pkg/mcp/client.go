package mcp

import (
	"context"

	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/wenchezhao/aiagent/pkg/utils/log"
)

func GetMCPTool(ctx context.Context, clientUrl string) []tool.BaseTool {
	cli, err := client.NewSSEMCPClient(clientUrl)
	if err != nil {
		log.Logger.Fatal(err.Error())

	}
	err = cli.Start(ctx)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "mcp-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	return tools
}
