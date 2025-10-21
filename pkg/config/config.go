package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type LlmType string

const (
	Ollama      LlmType = "ollama"
	OpenAI      LlmType = "openai"
	Siliconflow LlmType = "siliconflow"
)

// Config 结构体定义整体配置
type Config struct {
	ModelType         LlmType           `yaml:"modelType"`
	McpUrl            string            `yaml:"mcpUrl"`
	OllamaConfig      OllamaConfig      `yaml:"ollamaConfig"`
	OpenAIConfig      OpenAIConfig      `yaml:"openaiConfig"`
	SiliconflowConfig SiliconflowConfig `yaml:"siliconflowConfig"`
}

// OllamaConfig 定义 Ollama 相关配置
type OllamaConfig struct {
	Model   string `yaml:"model"`
	BaseURL string `yaml:"baseUrl"`
}

// OpenAIConfig 定义 OpenAI 相关配置
type OpenAIConfig struct {
	APIKey  string `yaml:"apiKey"`
	Model   string `yaml:"model"`
	BaseURL string `yaml:"baseUrl"`
}

// OpenAIConfig 定义 OpenAI 相关配置
type SiliconflowConfig struct {
	APIKey  string `yaml:"apiKey"`
	Model   string `yaml:"model"`
	BaseURL string `yaml:"baseUrl"`
}

// LoadConfig 从文件加载配置
func LoadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
