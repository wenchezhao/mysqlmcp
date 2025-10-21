// // package main

// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"fmt"
// // 	"log"
// // 	"time"

// // 	"github.com/mark3labs/mcp-go/client"
// // 	"github.com/mark3labs/mcp-go/mcp"
// // )

// // func main() {
// // 	c, err := client.NewStdioMCPClient("go", nil, "run", "../server/main.go")
// // 	if err != nil {
// // 		log.Fatalf("Failed to create client: %v", err)
// // 	}
// // 	defer c.Close()

// // 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// // 	defer cancel()

// // 	// Initialize
// // 	initRequest := mcp.InitializeRequest{}
// // 	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
// // 	initRequest.Params.ClientInfo = mcp.Implementation{
// // 		Name:    "mysql-mcp-client",
// // 		Version: "1.0.0",
// // 	}

// // 	initRes, err := c.Initialize(ctx, initRequest)
// // 	if err != nil {
// // 		log.Fatalf("Initialize failed: %v", err)
// // 	}
// // 	fmt.Printf("Connected to server: %s %s\n", initRes.ServerInfo.Name, initRes.ServerInfo.Version)

// // 	// List Tools
// // 	toolsRes, err := c.ListTools(ctx, mcp.ListToolsRequest{})
// // 	if err != nil {
// // 		log.Fatalf("List tools failed: %v", err)
// // 	}
// // 	for _, t := range toolsRes.Tools {
// // 		fmt.Printf("- %s: %s\n", t.Name, t.Description)
// // 	}

// // 	mysqlRequest := mcp.CallToolRequest{}
// // 	mysqlRequest.Params.Name = "query_mysql"
// // 	mysqlRequest.Params.Arguments = map[string]interface{}{
// // 		"query": "SELECT id, name, age FROM users WHERE age > 30",
// // 	}
// // 	// Call query_mysql
// // 	callRes, err := c.CallTool(ctx, mysqlRequest)
// // 	if err != nil {
// // 		log.Fatalf("Call tool failed: %v", err)
// // 	}

// // 	// Print result
// // 	printToolResult(callRes)
// // }

// // // Helper function to print tool results
// // func printToolResult(result *mcp.CallToolResult) {
// // 	for _, content := range result.Content {
// // 		if textContent, ok := content.(mcp.TextContent); ok {
// // 			fmt.Println(textContent.Text)
// // 		} else {
// // 			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
// // 			fmt.Println(string(jsonBytes))
// // 		}
// // 	}
// // }

// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strings"
// 	"time"

// 	"github.com/mark3labs/mcp-go/client"
// 	"github.com/mark3labs/mcp-go/mcp"
// )

// // 自定义工具名称
// const ToolName = "query_mysql"

// // Ollama 地址
// const OLLAMA_URL = "http://192.168.31.222:11434"

// type GenerateRequest struct {
// 	Model    string `json:"model"`
// 	Prompt   string `json:"prompt"`
// 	Stream   bool   `json:"stream,omitempty"`
// 	MaxToken int    `json:"max_token,omitempty"`
// }

// type GenerateResponse struct {
// 	Response string `json:"response"`
// 	Done     bool   `json:"done"`
// }

// func callOllamaGenerate(ctx context.Context, prompt string) (string, error) {
// 	reqBody := GenerateRequest{
// 		Model:  "qwen3:4b", // 可替换为 qwen3、mistral 等
// 		Prompt: prompt,
// 		Stream: false,
// 	}

// 	bodyBytes, _ := json.Marshal(reqBody)
// 	req, err := http.NewRequestWithContext(ctx, "POST", OLLAMA_URL+"/api/generate", bytes.NewBuffer(bodyBytes))
// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}

// 	defer resp.Body.Close()
// 	fmt.Println(resp.Body)
// 	var genResp GenerateResponse
// 	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
// 		return "", err
// 	}

// 	return strings.TrimSpace(genResp.Response), nil
// }

// func main() {
// 	mcpClient, err := client.NewSSEMCPClient("http://localhost:8088" + "/sse")
// 	if err != nil {
// 		log.Fatalf("Failed to create client: %v", err)
// 	}
// 	defer mcpClient.Close()

// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()

// 	// Start the client
// 	if err := mcpClient.Start(ctx); err != nil {
// 		log.Fatalf("Failed to start client: %v", err)
// 	}

// 	// Initialize
// 	initRequest := mcp.InitializeRequest{}
// 	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
// 	initRequest.Params.ClientInfo = mcp.Implementation{
// 		Name:    "MySQL DB Server",
// 		Version: "1.0.0",
// 	}

// 	result, err := mcpClient.Initialize(ctx, initRequest)
// 	if err != nil {
// 		log.Fatalf("Failed to initialize: %v", err)
// 	}

// 	if result.ServerInfo.Name != "MySQL DB Server" {
// 		log.Fatalf(
// 			"Expected server name 'MySQL DB Server', got '%s'",
// 			result.ServerInfo.Name,
// 		)
// 	}

// 	// Test Ping
// 	if err := mcpClient.Ping(ctx); err != nil {
// 		log.Fatalf("Ping failed: %v", err)
// 	}

// 	// 获取可用工具列表
// 	toolsRes, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
// 	if err != nil {
// 		log.Fatalf("获取工具列表失败: %v", err)
// 	}
// 	for _, t := range toolsRes.Tools {
// 		fmt.Printf("- 工具: %s - %s\n", t.Name, t.Description)
// 	}

// 	// 用户输入示例
// 	userInput := "找出年龄大于30岁的用户信息。"
// 	fmt.Printf("\n🗣️ 用户提问: %s\n", userInput)

// 	// 构造 Prompt，让 LLM 生成 SQL 并调用工具
// 	sqlGenPrompt := fmt.Sprintf(`
// 你是一个智能助手，能根据用户的自然语言查询生成对应的 SQL 语句，并调用数据库工具进行查询。

// 请按以下步骤处理用户的问题：

// 1. 分析用户的自然语言请求。
// 2. 生成对应的 SQL 查询语句。
// 3. 使用名为 "%s" 的工具执行该 SQL。
// 4. 根据查询结果生成自然语言的回答。

// 用户请求: "%s"

// 请直接输出 SQL 查询语句，不要解释：
// `, ToolName, userInput)

// 	// Step 1: 让 LLM 生成 SQL 查询
// 	sqlQuery, err := callOllamaGenerate(ctx, sqlGenPrompt)
// 	if err != nil {
// 		log.Fatalf("调用 Ollama 失败: %v", err)
// 	}
// 	fmt.Printf("🧠 LLM 生成的 SQL 查询: %s\n", sqlQuery)

// 	// Step 2: 调用 MCP 工具执行 SQL 查询
// 	callToolReq := mcp.CallToolRequest{}

// 	callResult, err := mcpClient.CallTool(ctx, callToolReq)
// 	if err != nil {
// 		log.Fatalf("调用工具失败: %v", err)
// 	}

// 	// Step 3: 解析数据库返回的结果
// 	printToolResult(callResult)
// 	// var dbResult []map[string]interface{}
// 	// err = json.Unmarshal(result.Content, &dbResult)
// 	// if err != nil {
// 	// 	log.Fatalf("解析数据库结果失败: %v", err)
// 	// }

// 	// fmt.Println("📊 数据库返回结果:")
// 	// for i, row := range dbResult {
// 	// 	fmt.Printf("  [%d] %v\n", i+1, row)
// 	// }

// 	// Step 4: 将结果交给 LLM 生成自然语言回复
// 	// 	answerPrompt := fmt.Sprintf(`
// 	// 以下是数据库查询结果，请根据这些数据生成一段自然语言的回答：

// 	// SQL 查询: %s
// 	// 查询结果:
// 	// %v

// 	// 请以简洁易懂的方式总结结果，不需要技术术语。
// 	// `, sqlQuery, dbResult)

// 	// 	finalAnswer, err := callOllamaGenerate(ctx, answerPrompt)
// 	// 	if err != nil {
// 	// 		log.Fatalf("生成最终回答失败: %v", err)
// 	// 	}

// 	// fmt.Printf("\n🤖 模型回答:\n%s\n", finalAnswer)
// }

// func printToolResult(result *mcp.CallToolResult) {
// 	for _, content := range result.Content {
// 		if textContent, ok := content.(mcp.TextContent); ok {
// 			fmt.Println(textContent.Text)
// 		} else {
// 			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
// 			fmt.Println(string(jsonBytes))
// 		}
// 	}
// }

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ollama"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func createOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://192.168.31.222:11434", // Ollama 服务地址
		Model:   "qwen3:4b",                    // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

const (
	nodeModel     = "node_model"
	nodeTools     = "node_tools"
	nodeTemplate  = "node_template"
	nodeConverter = "node_converter"
)

func main() {
	systemTpl := `你是一名计算专家，请根据用户输入的问题进行解答。`
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(systemTpl),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{query}"),
	)
	takeOne := compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (*schema.Message, error) {
		if len(input) == 0 {
			return nil, errors.New("input is empty")
		}
		return input[0], nil
	})
	branch := compose.NewStreamGraphBranch(func(ctx context.Context, input *schema.StreamReader[*schema.Message]) (string, error) {
		defer input.Close()
		msg, err := input.Recv()
		if err != nil {
			return "", err
		}

		if len(msg.ToolCalls) > 0 {
			return nodeTools, nil
		}

		return compose.END, nil
	}, map[string]bool{compose.END: true, nodeTools: true})
	////////////////////////////////////////////////////////////////////////////

	//startMCPServer()
	time.Sleep(1 * time.Second)
	ctx := context.Background()

	mcpTools := getMCPTool(ctx)

	chatModel := createOllamaChatModel(ctx)
	// 获取工具信息, 用于绑定到 ChatModel
	toolInfos := make([]*schema.ToolInfo, 0, len(mcpTools))
	for _, todoTool := range mcpTools {
		info, err := todoTool.Info(ctx)
		if err != nil {
			fmt.Printf("get ToolInfo failed, err=%v", err)
			return
		}
		toolInfos = append(toolInfos, info)
	}
	// 将 tools 绑定到 ChatModel
	err := chatModel.BindTools(toolInfos)
	if err != nil {
		fmt.Printf("BindTools failed, err=%v", err)
		return
	}

	// 创建 tools 节点
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: mcpTools,
	})
	if err != nil {
		fmt.Printf("NewToolNode failed, err=%v", err)
		return
	}

	graph := compose.NewGraph[map[string]any, *schema.Message]()

	_ = graph.AddChatTemplateNode(nodeTemplate, chatTpl)
	_ = graph.AddChatModelNode(nodeModel, chatModel)
	_ = graph.AddToolsNode(nodeTools, toolsNode)
	_ = graph.AddLambdaNode(nodeConverter, takeOne)

	_ = graph.AddEdge(compose.START, nodeTemplate)
	_ = graph.AddEdge(nodeTemplate, nodeModel)
	_ = graph.AddBranch(nodeModel, branch)
	_ = graph.AddEdge(nodeTools, nodeConverter)
	_ = graph.AddEdge(nodeConverter, compose.END)

	r, err := graph.Compile(ctx)
	if err != nil {
		fmt.Printf("Compile failed, err=%v", err)
		os.Exit(0)
	}

	out, err := r.Invoke(ctx, map[string]any{"query": "你好，你是谁"})
	if err != nil {
		fmt.Printf("Invoke failed, err=%v", err)
		os.Exit(0)
	}
	fmt.Printf("result content: %v", out.Content)
	out, err = r.Invoke(ctx, map[string]any{"query": "list table"})
	if err != nil {
		fmt.Printf("Invoke failed, err=%v", err)
		os.Exit(0)
	}
	fmt.Printf("result content: %v", out.Content)
}

func getMCPTool(ctx context.Context) []tool.BaseTool {
	cli, err := client.NewSSEMCPClient("http://localhost:12345/sse")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		log.Fatal(err)
	}

	return tools
}

func startMCPServer() {
	svr := server.NewMCPServer("demo", mcp.LATEST_PROTOCOL_VERSION)
	svr.AddTool(mcp.NewTool("calculate",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		arg := request.Params.Arguments.(map[string]any)
		op := arg["operation"].(string)
		x := arg["x"].(float64)
		y := arg["y"].(float64)

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultText("Cannot divide by zero"), nil
			}
			result = x / y
		}
		log.Printf("Calculated result: %.2f", result)
		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})
	go func() {
		defer func() {
			e := recover()
			if e != nil {
				fmt.Println(e)
			}
		}()

		err := server.NewSSEServer(svr, server.WithBaseURL("http://localhost:12345")).Start("localhost:12345")

		if err != nil {
			log.Fatal(err)
		}
	}()
}
