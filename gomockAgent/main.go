package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

// --- Configuration ---
// Read from environment variables
var (
	openaiAPIKey      = os.Getenv("OPENAI_API_KEY")                           // Use OPENAI_API_KEY now
	openaiAPIEndpoint = os.Getenv("OPENAI_API_BASE") + "/v1/chat/completions" // Allow overriding base URL
	openaiModel       = os.Getenv("OPENAI_MODEL")                             // Allow specifying model
)

type Agent struct {
	restyClient    *resty.Client // Use standard HTTP client
	getUserMessage func() (string, bool)
	model          string                    // Store the target model name
	tools          map[string]ToolDefinition // Map of tool names to tool definitions
	systemPrompt   string                    // Store the system prompt
}

func NewAgent(
	getUserMessage func() (string, bool),
	model string,
	tools []ToolDefinition,
) *Agent {

	restyClient := resty.New().SetTimeout(60 * time.Second)
	toolMap := make(map[string]ToolDefinition)
	for _, tool := range tools {
		toolMap[tool.Name] = tool
	}
	return &Agent{
		restyClient:    restyClient, // Add a timeout
		getUserMessage: getUserMessage,
		model:          model,
		tools:          toolMap,
		// Define the system prompt here or pass it in
		systemPrompt: "You are a helpful Go programmer assistant. You have access to tools to interact with the local filesystem (read, list, edit files). Use them when appropriate to fulfill the user's request. When editing, be precise about the changes. Respond ONLY with tool calls if you need to use tools, otherwise respond with text.",
	}
}

// callOpenAICompletion uses standard library http client
func (a *Agent) callOpenAICompletion(ctx context.Context, conversation []OpenAIChatCompletionMessage) (*OpenAIChatCompletionResponse, error) {

	// Prepare tools in OpenAI format
	openaiTools := []OpenAIChatCompletionTool{}
	for _, toolDef := range a.tools {
		openaiTools = append(openaiTools, OpenAIChatCompletionTool{
			Type: "function",
			Function: OpenAIChatCompletionFunctionDefinition{
				Name:        toolDef.Name,
				Description: toolDef.Description,
				Parameters:  toolDef.InputSchema,
			},
		})
	}

	// Build request payload
	requestPayload := OpenAIChatCompletionRequest{
		Model:       a.model,
		Messages:    conversation,
		Tools:       openaiTools,
		ToolChoice:  "auto", // Let the model decide when to use tools
		MaxTokens:   2048,   // Or make configurable
		Temperature: 0.7,    // Reasonable default
	}
	reply := &OpenAIChatCompletionResponse{}
	resp, err := a.restyClient.SetDebug(true).R().
		SetContext(ctx).
		SetBody(requestPayload).
		SetAuthToken(openaiAPIKey).
		SetResult(reply).
		Post(openaiAPIEndpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	return reply, nil
}

func (a *Agent) Run(ctx context.Context) error {

	conversation := []OpenAIChatCompletionMessage{
		{Role: "system", Content: a.systemPrompt}, // Start with system prompt
	}
	fmt.Println("Chat with AI (use 'ctrl-c' to quit)")
	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ") // Blue prompt for user
		userMessage, ok := a.getUserMessage()
		if !ok {
			return fmt.Errorf("failed to get user message")
		}
		if userMessage == "" {
			continue
		}

		conversation = append(conversation, OpenAIChatCompletionMessage{Role: "user", Content: userMessage})
		// 如果需要使用工具，则需要多次调用OpenAI API
		for {
			resp, err := a.callOpenAICompletion(ctx, conversation)
			if err != nil {
				fmt.Printf("\u001b[91mAPI Error\u001b[0m: %s\n", err.Error())
				break
			}
			// Process the response
			if len(resp.Choices) == 0 {
				fmt.Println("\u001b[91mError\u001b[0m: OpenAI response contained no choices.")
				break // Break inner loop, let user re-prompt
			}
			assistantMessage := resp.Choices[0].Message

			// Add assistant's message (text and/or tool calls) to conversation
			conversation = append(conversation, assistantMessage)
			// 如果返回了内容，则直接输出
			if assistantMessage.Content != "" {
				fmt.Printf("\u001b[93mAI\u001b[0m: %s\n", assistantMessage.Content) // Yellow for AI
			}
			if len(assistantMessage.ToolCalls) == 0 {
				// No tools called, break inner loop and wait for next user input
				break
			}

			// 如果返回了工具调用，则需要调用工具
			// Execute tools and collect results
			toolResults := []OpenAIChatCompletionMessage{}
			for _, toolCall := range assistantMessage.ToolCalls {
				if toolCall.Type != "function" {
					continue // Skip non-function tool calls if any
				}

				toolName := toolCall.Function.Name
				toolArgs := toolCall.Function.Arguments // This is a JSON *string*

				fmt.Printf("\u001b[92mTool Call\u001b[0m: %s(%s)\n", toolName, toolArgs) // Green
				// 获取函数名称
				toolDef, found := a.tools[toolName]
				var resultMsg OpenAIChatCompletionMessage
				if !found {
					errorMsg := fmt.Sprintf("tool '%s' not found by agent", toolName)
					fmt.Printf("\u001b[91mTool Error\u001b[0m: %s\n", errorMsg)
					resultMsg = OpenAIChatCompletionMessage{
						Role:       "tool",
						ToolCallID: toolCall.ID,
						Content:    errorMsg, // Report error back to OpenAI
						Name:       toolName,
					}
				} else {
					// Execute the actual tool function
					// Note: toolArgs is a JSON string, pass it as json.RawMessage
					// 执行函数
					toolOutput, err := toolDef.Function(json.RawMessage(toolArgs))
					if err != nil {
						errorMsg := fmt.Sprintf("error executing tool '%s': %s", toolName, err.Error())
						fmt.Printf("\u001b[91mTool Error\u001b[0m: %s\n", errorMsg)
						resultMsg = OpenAIChatCompletionMessage{
							Role:       "tool",
							ToolCallID: toolCall.ID,
							Content:    errorMsg, // Report error back to OpenAI
							Name:       toolName,
						}
					} else {
						// Log successful tool execution result (optional)
						// fmt.Printf("\u001b[92mTool Result\u001b[0m: %s\n", toolOutput)
						resultMsg = OpenAIChatCompletionMessage{
							Role:       "tool",
							ToolCallID: toolCall.ID,
							Content:    toolOutput, // Send success result back to OpenAI
							Name:       toolName,
						}
					}
				}
				toolResults = append(toolResults, resultMsg)
			} // End of processing tool calls for one response
			conversation = append(conversation, toolResults...)
		}
	}
}

func main() {
	// --- Configuration Checks ---
	if openaiAPIKey == "" {
		fmt.Fprintln(os.Stderr, "\u001b[91mError: OPENAI_API_KEY environment variable not set.\u001b[0m")
		os.Exit(1)
	}
	if os.Getenv("OPENAI_API_BASE") == "" {
		// Default to official OpenAI endpoint if base URL not set
		openaiAPIEndpoint = "https://api.openai.com/v1/chat/completions"
		fmt.Println("Info: OPENAI_API_BASE not set, defaulting to https://api.openai.com")
	}
	if openaiModel == "" {
		// Default model if not set
		openaiModel = "gpt-3.5-turbo" // Or "gpt-3.5-turbo" or another compatible model
		fmt.Printf("Info: OPENAI_MODEL not set, defaulting to %s\n", openaiModel)
	}
	fmt.Printf(" openapibase: %s\n", openaiAPIEndpoint)
	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "\u001b[91mError reading input: %v\u001b[0m\n", err)
				return "", false
			}
			return "", false // EOF
		}
		if scanner.Text() == "exit" {
			return "", false
		}
		return scanner.Text(), true
	}
	// 工具定义
	tools := []ToolDefinition{
		ReadFileDefinition,
		ListFilesDefinition,
		GetMergeDiffDefinition,
	}
	agent := NewAgent(getUserMessage, openaiModel, tools)
	err := agent.Run(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "\u001b[91mAgent exited with error: %s\u001b[0m\n", err.Error())
		os.Exit(1)
	}
}
