package types

type AiReqBody struct {
	Model               string      `json:"model"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	FrequencyPenalty    float32     `json:"frequency_penalty"`
	TopP                float32     `json:"top_p"`
	Messages            []AiMessage `json:"messages"`
	Tools               []AiTool    `json:"tools"`
}

type AiMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls"`
	Refusal    *string    `json:"refusal,omitempty"`
	Name       string     `json:"name,omitempty"`
	ToolCallId string     `json:"tool_call_id,omitempty"`
}

type AiTool struct {
	Type     string     `json:"type"`
	Function AiFunction `json:"function"`
}

type AiFunction struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Parameters  AiParameters `json:"parameters"`
}

type AiParameters struct {
	Type       string                `json:"type"`
	Properties map[string]AiProperty `json:"properties"`
	Required   []string              `json:"required"`
}

type AiProperty struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Enum        []string    `json:"enum,omitempty"`
	Items       *AiProperty `json:"items,omitempty"`
}

// RESPONSE PARSER
type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      AiMessage `json:"message"`
	Logprobs     *string   `json:"logprobs,omitempty"`
	FinishReason string    `json:"finish_reason"`
}

type ToolCall struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

type Function struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Usage struct {
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

type PromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}
