package eino

import (
	"context"
	"edustate/internal/conf"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"io"
	"log"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
)

func Client(conf *conf.LLM) {
	ctx := context.Background()

	// 使用模版创建messages
	log.Printf("===create messages===\n")
	messages := createMessagesFromTemplate()
	log.Printf("messages: %+v\n\n", messages)

	// 创建llm
	log.Printf("===create llm===\n")
	cm := createArkChatModel(ctx, conf)
	// cm := createOllamaChatModel(ctx)
	log.Printf("create llm success\n\n")

	log.Printf("===llm generate===\n")
	result := generate(ctx, cm, messages)
	log.Printf("result: %+v\n\n", result)

	log.Printf("===llm stream generate===\n")
	streamResult := stream(ctx, cm, messages)
	reportStream(streamResult)
}

func generate(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.Message {
	result, err := llm.Generate(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func stream(ctx context.Context, llm model.ChatModel, in []*schema.Message) *schema.StreamReader[*schema.Message] {
	result, err := llm.Stream(ctx, in)
	if err != nil {
		log.Fatalf("llm generate failed: %v", err)
	}
	return result
}

func createTemplate() prompt.ChatTemplate {
	// 创建模板，使用 FString 格式
	return prompt.FromMessages(schema.FString,
		// 系统消息模板
		schema.SystemMessage("你是一个{role}。你需要用{style}的语气回答问题。你的目标是帮助程序员保持积极乐观的心态，提供技术建议的同时也要关注他们的心理健康。"),

		// 插入需要的对话历史（新对话的话这里不填）
		schema.MessagesPlaceholder("chat_history", true),

		// 用户消息模板
		schema.UserMessage("问题: {question}"),
	)
}

func createMessagesFromTemplate() []*schema.Message {
	template := createTemplate()

	// 使用模板生成消息
	messages, err := template.Format(context.Background(), map[string]any{
		"role":     "程序员鼓励师",
		"style":    "积极、温暖且专业",
		"question": "我的代码一直报错，感觉好沮丧，该怎么办？",
		// 对话历史（这个例子里模拟两轮对话历史）
		"chat_history": []*schema.Message{
			schema.UserMessage("你好"),
			schema.AssistantMessage("嘿！我是你的程序员鼓励师！记住，每个优秀的程序员都是从 Debug 中成长起来的。有什么我可以帮你的吗？", nil),
			schema.UserMessage("我觉得自己写的代码太烂了"),
			schema.AssistantMessage("每个程序员都经历过这个阶段！重要的是你在不断学习和进步。让我们一起看看代码，我相信通过重构和优化，它会变得更好。记住，Rome wasn't built in a day，代码质量是通过持续改进来提升的。", nil),
		},
	})
	if err != nil {
		log.Fatalf("format template failed: %v\n", err)
	}
	return messages
}

func reportStream(sr *schema.StreamReader[*schema.Message]) {
	defer sr.Close()
	i := 0
	for {
		message, err := sr.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("recv failed: %v", err)
		}
		log.Printf("message[%d]: %+v\n", i, message)
		i++
	}
}

func createOllamaChatModel(ctx context.Context) model.ChatModel {
	chatModel, err := ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
		BaseURL: "http://localhost:11434", // Ollama 服务地址
		Model:   "llama2",                 // 模型名称
	})
	if err != nil {
		log.Fatalf("create ollama chat model failed: %v", err)
	}
	return chatModel
}

func createOpenAIChatModel(ctx context.Context) model.ChatModel {
	key := os.Getenv("OPENAI_API_KEY")
	modelName := os.Getenv("OPENAI_MODEL_NAME")
	baseURL := os.Getenv("OPENAI_BASE_URL")
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  key,
	})
	if err != nil {
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}

func createArkChatModel(ctx context.Context, conf *conf.LLM) model.ChatModel {
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
		log.Fatalf("create openai chat model failed, err=%v", err)
	}
	return chatModel
}
