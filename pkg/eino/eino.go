package eino

import (
	"context"
	"edustate/internal/conf"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/go-kratos/kratos/v2/log"
)

var LLMClient *LLM

type LLM struct {
	Client model.ChatModel
	Log    *log.Helper
}

func Init(logger *log.Helper, conf *conf.LLM) error {
	chatModel, err := createArkChatModel(context.Background(), conf)
	LLMClient = &LLM{
		Client: chatModel,
		Log:    logger,
	}
	return err
}

func (llm *LLM) NLToArgs(nlInputStr string) (string, error) {
	ctx := context.Background()
	messages := llm.createMessagesFromTemplate(nlInputStr)
	result, err := llm.Client.Generate(ctx, messages)
	if err != nil {
		llm.Log.Error("generate err:", err)
		return "", err
	}
	// 解析 JSON 响应
	var response struct {
		StudentID string `json:"student_id"`
	}
	if err := json.Unmarshal([]byte(result.Content), &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %v", err)
	}
	return response.StudentID, nil
}

func (llm *LLM) createMessagesFromTemplate(nlInputStr string) []*schema.Message {
	template := prompt.FromMessages(schema.FString,
		schema.SystemMessage(`你是一个数据分析师。请从用户输入中提取学生的ID。请以JSON格式返回。`),
		schema.MessagesPlaceholder("chat_history", false),
		schema.UserMessage("输入文本: {nlInputStr}"),
	)
	messages, err := template.Format(context.Background(), map[string]any{
		"chat_history": []*schema.Message{},
		"nlInputStr":   nlInputStr,
	})
	if err != nil {
		log.Log(log.LevelError, "msg", "format template failed", "err", err)
		return nil
	}
	return messages
}

func createArkChatModel(ctx context.Context, conf *conf.LLM) (model.ChatModel, error) {
	chatModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		Timeout:    nil,
		HTTPClient: nil,
		RetryTimes: nil,
		BaseURL:    conf.GetApiBase(),
		Region:     "",
		APIKey:     conf.GetApiKey(),
		AccessKey:  "",
		SecretKey:  "",
		Model:      conf.GetModel(),
	})
	if err != nil {
		log.Context(ctx).Infof("create openai chat model failed, err=%v", err)
	}
	return chatModel, err
}
