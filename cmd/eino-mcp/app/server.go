package app

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/wenchezhao/aiagent/pkg/agent"
	"github.com/wenchezhao/aiagent/pkg/config"
	"github.com/wenchezhao/aiagent/pkg/llm"
	"github.com/wenchezhao/aiagent/pkg/mcp"
	signals "github.com/wenchezhao/aiagent/pkg/utils"
	"github.com/wenchezhao/aiagent/pkg/utils/log"
)

func NewAgentCommand() *cobra.Command {
	s := NewServerRunOptions()
	cmd := &cobra.Command{
		Use: "agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(signals.SetupSignalHandler(), s)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments ,got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cmd.Flags().AddFlagSet(s.Flags())
	return cmd
}

func Run(ctx context.Context, s *ServerRunOptions) error {
	conf, err := config.LoadConfig(s.ConfigFilePath)
	if err != nil {
		log.Logger.Fatal("load config failed", zap.Error(err))
	}
	var cm model.ToolCallingChatModel
	switch conf.ModelType {
	case config.OpenAI:
		cm = llm.CreateOpenAIChatModel(ctx, conf.OpenAIConfig)
	case config.Siliconflow:
		cm = llm.CreateSiliconflowChatModel(ctx, conf.SiliconflowConfig)
	case config.Ollama:
		cm = llm.CreateOllamaChatModel(ctx, conf.OllamaConfig)
	default:
		log.Logger.Fatal("unsupported model type", zap.String("modelType", string(conf.ModelType)))
	}
	mcpTools := mcp.GetMCPTool(ctx, conf.McpUrl)
	einoAgent := agent.NewAgent(ctx, cm)
	einoAgent.BindTools(mcpTools)
	einoAgent.Run()

	return nil
}
