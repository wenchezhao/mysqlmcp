package agent

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"go.uber.org/zap"

	"github.com/wenchezhao/aiagent/pkg/utils/log"
)

const (
	nodeModel     = "node_model"
	nodeTools     = "node_tools"
	nodeTemplate  = "node_template"
	nodeConverter = "node_converter"
)

type Agent struct {
	ctx       context.Context
	graph     *compose.Graph[map[string]any, *schema.Message]
	chatTpl   *prompt.DefaultChatTemplate
	takeOne   *compose.Lambda
	branch    *compose.GraphBranch
	chatModel model.ToolCallingChatModel
	tools     []tool.BaseTool
}

func NewAgent(ctx context.Context, chatModel model.ToolCallingChatModel) *Agent {
	return &Agent{
		ctx:   ctx,
		graph: compose.NewGraph[map[string]any, *schema.Message](),
		chatTpl: prompt.FromMessages(schema.FString,
			schema.SystemMessage(`你是一名数据库管理员，请根据输入的问题进行解答。`),
			schema.MessagesPlaceholder("message_histories", true),
			schema.UserMessage("{query}"),
		),
		takeOne: compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) (*schema.Message, error) {
			if len(input) == 0 {
				return nil, errors.New("input is empty")
			}
			return input[0], nil
		}),
		branch: compose.NewStreamGraphBranch(func(ctx context.Context, input *schema.StreamReader[*schema.Message]) (string, error) {
			defer input.Close()
			msg, err := input.Recv()
			if err != nil {
				return "", err
			}
			if len(msg.ToolCalls) > 0 {
				return nodeTools, nil
			}
			return compose.END, nil
		}, map[string]bool{compose.END: true, nodeTools: true}),
		chatModel: chatModel,
	}
}

func (a *Agent) BindTools(tools []tool.BaseTool) {

	a.tools = append(a.tools, tools...)

	toolsInfo := make([]*schema.ToolInfo, 0, len(tools))
	for _, tool := range a.tools {
		info, err := tool.Info(a.ctx)
		if err != nil {
			log.Logger.Fatal("get ToolInfo failed", zap.Error(err))
			return
		}
		toolsInfo = append(toolsInfo, info)
	}
	var err error
	if a.chatModel, err = a.chatModel.WithTools(toolsInfo); err != nil {
		log.Logger.Fatal("WithTools failed", zap.Error(err))
		os.Exit(0)
	}

}

func (a *Agent) Run() error {

	toolsNode, err := compose.NewToolNode(a.ctx, &compose.ToolsNodeConfig{
		Tools: a.tools,
	})
	if err != nil {
		log.Logger.Fatal("NewToolNode failed ", zap.Error(err))
		return err
	}

	_ = a.graph.AddChatTemplateNode(nodeTemplate, a.chatTpl)
	_ = a.graph.AddChatModelNode(nodeModel, a.chatModel)
	_ = a.graph.AddToolsNode(nodeTools, toolsNode)
	_ = a.graph.AddLambdaNode(nodeConverter, a.takeOne)

	_ = a.graph.AddEdge(compose.START, nodeTemplate)
	_ = a.graph.AddEdge(nodeTemplate, nodeModel)
	_ = a.graph.AddBranch(nodeModel, a.branch)
	_ = a.graph.AddEdge(nodeTools, nodeConverter)
	_ = a.graph.AddEdge(nodeConverter, compose.END)

	r, err := a.graph.Compile(a.ctx)
	if err != nil {
		log.Logger.Fatal("Compile failed", zap.Error(err))
		os.Exit(0)
	}

	// 持续从标准输入读取，逐条调用 Invoke 并输出结果
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("输入问题后回车（输入 exit 或 Ctrl+D 退出）：")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			if scanErr := scanner.Err(); scanErr != nil {
				fmt.Printf("input scan error: %v\n", scanErr)
			}
			break
		}
		query := strings.TrimSpace(scanner.Text())
		if query == "" {
			continue
		}
		if query == "exit" || query == "quit" {
			break
		}

		out, err := r.Invoke(a.ctx, map[string]any{"query": query})
		//out, err := r.Invoke(a.ctx, map[string]any{"query": "获取数据库中所有的表名"})
		if err != nil {
			fmt.Printf("Invoke failed, err=%v\n", err)
			continue
		}
		fmt.Printf("result content: %v\n", out.Content)
	}
	return nil
}
