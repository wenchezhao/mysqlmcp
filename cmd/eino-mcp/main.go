package main

import (
	"github.com/wenchezhao/aiagent/cmd/eino-mcp/app"
	"github.com/wenchezhao/aiagent/pkg/utils/log"
)

func main() {
	defer log.Logger.Sync()
	app.NewAgentCommand().Execute()
}
