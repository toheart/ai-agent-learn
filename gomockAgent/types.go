package main

type OpenAIChatCompletionFunctionDefinition struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters"` // Use map[string]any to represent JSON schema object
}

type OpenAIChatCompletionTool struct {
	Type     string                                 `json:"type"` // Always "function" for now
	Function OpenAIChatCompletionFunctionDefinition `json:"function"`
}

type OpenAIChatCompletionRequest struct {
	Model       string                        `json:"model"`
	Messages    []OpenAIChatCompletionMessage `json:"messages"`
	ToolChoice  any                           `json:"tool_choice,omitempty"` // "auto" or specific tool
	MaxTokens   int                           `json:"max_tokens,omitempty"`
	Temperature float32                       `json:"temperature,omitempty"`
	// Add other OpenAI parameters as needed (top_p, stream, etc.)
	Tools []OpenAIChatCompletionTool `json:"tools,omitempty"`
}

type OpenAIChatCompletionMessage struct {
	Role       string                         `json:"role"`                   // "system", "user", "assistant", "tool"
	Content    string                         `json:"content,omitempty"`      // For text content or tool result
	ToolCalls  []OpenAIChatCompletionToolCall `json:"tool_calls,omitempty"`   // For assistant requesting tools
	ToolCallID string                         `json:"tool_call_id,omitempty"` // For tool role messages
	Name       string                         `json:"name,omitempty"`         // For tool role messages (function name) - Optional by OpenAI spec but sometimes useful
}
type OpenAIChatCompletionToolCall struct {
	ID       string                           `json:"id"`   // ID to match with tool response
	Type     string                           `json:"type"` // Always "function"
	Function OpenAIChatCompletionFunctionCall `json:"function"`
}

type OpenAIChatCompletionFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // Arguments are a *string* containing JSON
}

type OpenAIChatCompletionResponse struct {
	ID      string                       `json:"id"`
	Object  string                       `json:"object"`
	Created int64                        `json:"created"`
	Model   string                       `json:"model"`
	Choices []OpenAIChatCompletionChoice `json:"choices"`
}

type OpenAIChatCompletionChoice struct {
	Index        int                         `json:"index"`
	Message      OpenAIChatCompletionMessage `json:"message"`
	FinishReason string                      `json:"finish_reason"`
}
