package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	openRouterBaseURL = "https://openrouter.ai/api/v1"
	openRouterTimeout = 60 * time.Second
)

type OpenRouterCompletionProxy struct {
	apiKey     string
	httpClient *http.Client
}

func (p *OpenRouterCompletionProxy) Generate(ctx context.Context, request OpenRouterCompletionRequest) OpenRouterResponse {
	url := openRouterBaseURL + "/chat/completions"

	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("Failed to marshal request: %v", err),
					},
				},
			},
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("Failed to create request: %v", err),
					},
				},
			},
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("Failed to send request: %v", err),
					},
				},
			},
		}
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("Failed to read response: %v", err),
					},
				},
			},
		}
	}

	if resp.StatusCode != http.StatusOK {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("API error (status %d): %s", resp.StatusCode, string(respBody)),
					},
				},
			},
		}
	}

	var openRouterResp OpenRouterResponse
	if err := json.Unmarshal(respBody, &openRouterResp); err != nil {
		return OpenRouterResponse{
			ID:       "",
			Provider: "openrouter",
			Model:    request.Model,
			Object:   "error",
			Created:  time.Now().Unix(),
			Choices: []OpenRouterChoice{
				{
					Index:        0,
					FinishReason: "error",
					Message: &OpenRouterMessage{
						Role:    LLMRoleAssistant,
						Content: fmt.Sprintf("Failed to parse response: %v", err),
					},
				},
			},
		}
	}

	return openRouterResp
}
