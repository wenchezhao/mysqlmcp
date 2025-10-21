package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino/components/model"
	"go.uber.org/zap"

	"github.com/wenchezhao/aiagent/pkg/config"
	"github.com/wenchezhao/aiagent/pkg/utils/log"
)

func CreateOllamaChatModel(ctx context.Context, conf config.OllamaConfig) model.ToolCallingChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: conf.BaseURL,
		Model:   conf.Model,
	})
	if err != nil {
		log.Logger.Fatal("load config failed", zap.Error(err))
	}
	return chatModel
}
