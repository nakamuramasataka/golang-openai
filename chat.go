package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model            string          `json:"model"`
	Messages         []Message       `json:"messages"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[int]float64 `json:"logit_bias,omitempty"`
	LogProbs         *bool           `json:"logprobs,omitempty"`
	MaxTokens        *int            `json:"max_tokens,omitempty"`
	N                *int            `json:"n,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	ResponseFormat   *string         `json:"response_format,omitempty"`
	Seed             *int            `json:"seed,omitempty"`
	Stop             []string        `json:"stop,omitempty"`
	Stream           *bool           `json:"stream,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	Tools            []string        `json:"tools,omitempty"`
	ToolChoice       *string         `json:"tool_choice,omitempty"`
	User             *string         `json:"user,omitempty"`
}

type ChatCompletionResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index        int         `json:"index"`
		Message      Message     `json:"message"`
		LogProbs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func (c *Client) CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making request to chat completion API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("chat completion API returned non-200 status code: %d", resp.StatusCode)
	}

	var chatCompletionResponse ChatCompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&chatCompletionResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding response body: %w", err)
	}

	return &chatCompletionResponse, nil
}
