package llm

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"

	"github.com/wenchezhao/aiagent/pkg/config"
)

func CreateSiliconflowChatModel(ctx context.Context, conf config.SiliconflowConfig) model.ToolCallingChatModel {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: conf.BaseURL,
		Model:   conf.Model,
		APIKey:  conf.APIKey,
	})
	if err != nil {
		log.Fatalf("create siliconflow chat model failed, err=%v", err)
	}
	return chatModel
}
