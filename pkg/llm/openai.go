package llm

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"

	"github.com/wenchezhao/aiagent/pkg/config"
)

func CreateOpenAIChatModel(ctx context.Context, conf config.OpenAIConfig) model.ToolCallingChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: conf.BaseURL,
		Model:   conf.Model,
		APIKey:  conf.APIKey,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
