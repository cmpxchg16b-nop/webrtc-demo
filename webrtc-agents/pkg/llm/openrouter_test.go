package llm

// How to invoke the tests with appropriate APIKey ?
//
// First install godotenv
// go install github.com/joho/godotenv/cmd/godotenv@latest
//
// Then run
// godotenv -f .env.test -- go test -v ./pkg/llm/

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	testAPIKeyEnv = "OPENROUTER_APIKEY_TEST"
	testModel     = "deepseek/deepseek-v3.2"
	testTimeout   = 30 * time.Second
)

func TestGetAPIKey_DirectKey(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{
		APIKey: "test-direct-key",
	}

	got := proxy.getAPIKey()
	if got != "test-direct-key" {
		t.Errorf("getAPIKey() = %q, want %q", got, "test-direct-key")
	}
}

func TestGetAPIKey_FromEnv(t *testing.T) {
	t.Setenv(testAPIKeyEnv, "test-env-key")

	proxy := &OpenRouterCompletionProxy{
		APIKeyFromEnv: testAPIKeyEnv,
	}

	got := proxy.getAPIKey()
	if got != "test-env-key" {
		t.Errorf("getAPIKey() = %q, want %q", got, "test-env-key")
	}
}

func TestGetAPIKey_DirectKeyTakesPrecedence(t *testing.T) {
	t.Setenv(testAPIKeyEnv, "test-env-key")

	proxy := &OpenRouterCompletionProxy{
		APIKey:        "test-direct-key",
		APIKeyFromEnv: testAPIKeyEnv,
	}

	got := proxy.getAPIKey()
	if got != "test-direct-key" {
		t.Errorf("getAPIKey() = %q, want %q", got, "test-direct-key")
	}
}

func TestGetAPIKey_EmptyEnvVar(t *testing.T) {
	// Set an empty value for the env var
	t.Setenv(testAPIKeyEnv, "")

	proxy := &OpenRouterCompletionProxy{
		APIKeyFromEnv: testAPIKeyEnv,
	}

	got := proxy.getAPIKey()
	if got != "" {
		t.Errorf("getAPIKey() = %q, want empty string", got)
	}
}

func TestGetAPIKey_Empty(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	got := proxy.getAPIKey()
	if got != "" {
		t.Errorf("getAPIKey() = %q, want empty string", got)
	}
}

func TestGetHttpClient_Default(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	got := proxy.getHttpClient()
	if got == nil {
		t.Error("getHttpClient() returned nil, want non-nil client")
	}
	if got != http.DefaultClient {
		t.Error("getHttpClient() should return http.DefaultClient when HttpClient is nil")
	}
}

func TestGetHttpClient_Custom(t *testing.T) {
	customClient := &http.Client{Timeout: testTimeout}
	proxy := &OpenRouterCompletionProxy{
		HttpClient: customClient,
	}

	got := proxy.getHttpClient()
	if got != customClient {
		t.Error("getHttpClient() should return the custom client")
	}
}

func TestGetBaseURL_Default(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	got := proxy.getBaseURL()
	if got != openRouterBaseURL {
		t.Errorf("getBaseURL() = %q, want %q", got, openRouterBaseURL)
	}
}

func TestGetBaseURL_Custom(t *testing.T) {
	customURL := "https://custom.openrouter.example.com/v1"
	proxy := &OpenRouterCompletionProxy{
		BaseURL: customURL,
	}

	got := proxy.getBaseURL()
	if got != customURL {
		t.Errorf("getBaseURL() = %q, want %q", got, customURL)
	}
}

func TestGetCompletionURL_Default(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	got, err := proxy.getCompletionURL()
	if err != nil {
		t.Fatalf("getCompletionURL() returned error: %v", err)
	}

	expected := openRouterBaseURL + "/chat/completions"
	if got.String() != expected {
		t.Errorf("getCompletionURL() = %q, want %q", got.String(), expected)
	}
}

func TestGetCompletionURL_CustomBaseURL(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{
		BaseURL: "https://custom.example.com/api",
	}

	got, err := proxy.getCompletionURL()
	if err != nil {
		t.Fatalf("getCompletionURL() returned error: %v", err)
	}

	expected := "https://custom.example.com/api/chat/completions"
	if got.String() != expected {
		t.Errorf("getCompletionURL() = %q, want %q", got.String(), expected)
	}
}

func TestGetCompletionURL_InvalidURL(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{
		BaseURL: "://invalid-url",
	}

	_, err := proxy.getCompletionURL()
	if err == nil {
		t.Error("getCompletionURL() should return error for invalid URL")
	}
}

func TestGetErrorResponse(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	got := proxy.getErrorResponse(testModel, "test error message")

	if got.ID != "" {
		t.Errorf("ID = %q, want empty", got.ID)
	}
	if got.Provider != "openrouter" {
		t.Errorf("Provider = %q, want %q", got.Provider, "openrouter")
	}
	if got.Model != testModel {
		t.Errorf("Model = %q, want %q", got.Model, testModel)
	}
	if got.Object != "error" {
		t.Errorf("Object = %q, want %q", got.Object, "error")
	}
	if len(got.Choices) != 1 {
		t.Fatalf("len(Choices) = %d, want 1", len(got.Choices))
	}
	if got.Choices[0].Message == nil {
		t.Fatal("Message is nil")
	}
	if got.Choices[0].Message.Content != "test error message" {
		t.Errorf("Message.Content = %q, want %q", got.Choices[0].Message.Content, "test error message")
	}
	if got.Choices[0].Message.Role != LLMRoleAssistant {
		t.Errorf("Message.Role = %q, want %q", got.Choices[0].Message.Role, LLMRoleAssistant)
	}
}

func TestAuthorizationHeaderBuilder_Custom(t *testing.T) {
	customBuilder := func(apiKey string) string {
		return "CustomAuth " + apiKey
	}

	proxy := &OpenRouterCompletionProxy{
		APIKey:                 "test-key",
		GetAuthorizationHeader: customBuilder,
	}

	got := proxy.GetAuthorizationHeader(proxy.getAPIKey())
	if got != "CustomAuth test-key" {
		t.Errorf("Authorization header = %q, want %q", got, "CustomAuth test-key")
	}
}

func TestAuthorizationHeaderBuilder_Default(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{
		APIKey: "test-key",
	}

	var authHeaderGetter AuthorizationHeaderBuilder = defaultAuthHeaderGetter
	if proxy.GetAuthorizationHeader != nil {
		authHeaderGetter = proxy.GetAuthorizationHeader
	}

	got := authHeaderGetter(proxy.getAPIKey())
	if got != "Bearer test-key" {
		t.Errorf("Authorization header = %q, want %q", got, "Bearer test-key")
	}
}

// Integration test - requires actual API key
func TestGenerate_Integration(t *testing.T) {
	apiKey := os.Getenv(testAPIKeyEnv)
	if apiKey == "" {
		t.Skipf("Skipping integration test: %s not set", testAPIKeyEnv)
	}

	proxy := &OpenRouterCompletionProxy{
		APIKey: apiKey,
		HttpClient: &http.Client{
			Timeout: testTimeout,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), testTimeout)
	defer cancel()

	request := OpenRouterCompletionRequest{
		Model: testModel,
		Messages: []OpenRouterCompletionRequestMessage{
			{
				Role:    LLMRoleUser,
				Content: "Say 'hello' in exactly one word.",
			},
		},
		Reasoning: OpenRouterCompletionRequestReasoning{
			Enabled: false,
		},
	}

	resp := proxy.Generate(ctx, request)

	if resp.Object == "error" {
		t.Fatalf("Generate() returned error: %s", resp.Choices[0].Message.Content)
	}

	if resp.Model == "" {
		t.Error("Model is empty in response")
	}

	if len(resp.Choices) == 0 {
		t.Fatal("No choices in response")
	}

	if resp.Choices[0].Message == nil {
		t.Fatal("Message is nil in response")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		t.Error("Message content is empty")
	}

	t.Logf("Response: %s", content)
}

func TestGenerate_EmptyAPIKey(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{}

	ctx := context.Background()
	request := OpenRouterCompletionRequest{
		Model: testModel,
		Messages: []OpenRouterCompletionRequestMessage{
			{Role: LLMRoleUser, Content: "test"},
		},
	}

	resp := proxy.Generate(ctx, request)

	if resp.Object != "error" {
		t.Error("Generate() should return error response for empty API key")
	}

	if !strings.Contains(resp.Choices[0].Message.Content, "401") &&
		!strings.Contains(resp.Choices[0].Message.Content, "Unauthorized") &&
		!strings.Contains(resp.Choices[0].Message.Content, "error") {
		t.Logf("Note: Expected auth error, got: %s", resp.Choices[0].Message.Content)
	}
}

func TestGenerate_InvalidBaseURL(t *testing.T) {
	proxy := &OpenRouterCompletionProxy{
		APIKey:  "test-key",
		BaseURL: "://invalid-url",
	}

	ctx := context.Background()
	request := OpenRouterCompletionRequest{
		Model:    testModel,
		Messages: []OpenRouterCompletionRequestMessage{{Role: LLMRoleUser, Content: "test"}},
	}

	resp := proxy.Generate(ctx, request)

	if resp.Object != "error" {
		t.Error("Generate() should return error response for invalid base URL")
	}
}

func TestGenerate_ContextCancellation(t *testing.T) {
	apiKey := os.Getenv(testAPIKeyEnv)
	if apiKey == "" {
		t.Skipf("Skipping integration test: %s not set", testAPIKeyEnv)
	}

	proxy := &OpenRouterCompletionProxy{
		APIKey: apiKey,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	request := OpenRouterCompletionRequest{
		Model:    testModel,
		Messages: []OpenRouterCompletionRequestMessage{{Role: LLMRoleUser, Content: "test"}},
	}

	resp := proxy.Generate(ctx, request)

	if resp.Object != "error" {
		t.Error("Generate() should return error response for cancelled context")
	}
}

func TestDefaultAuthHeaderGetter(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
		want   string
	}{
		{
			name:   "normal key",
			apiKey: "sk-test-123",
			want:   "Bearer sk-test-123",
		},
		{
			name:   "empty key",
			apiKey: "",
			want:   "Bearer ",
		},
		{
			name:   "key with special chars",
			apiKey: "sk-test_123!@#",
			want:   "Bearer sk-test_123!@#",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := defaultAuthHeaderGetter(tt.apiKey)
			if got != tt.want {
				t.Errorf("defaultAuthHeaderGetter(%q) = %q, want %q", tt.apiKey, got, tt.want)
			}
		})
	}
}

func TestOpenRouterCompletionRequest_MarshalJSON(t *testing.T) {
	request := OpenRouterCompletionRequest{
		Model: testModel,
		Messages: []OpenRouterCompletionRequestMessage{
			{Role: LLMRoleUser, Content: "Hello"},
			{Role: LLMRoleAssistant, Content: "Hi there!"},
		},
		Reasoning: OpenRouterCompletionRequestReasoning{
			Enabled: true,
		},
	}

	data, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("json.Marshal() error: %v", err)
	}

	// Verify the JSON contains expected fields
	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"model":"`+testModel+`"`) {
		t.Errorf("JSON missing model: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"role":"user"`) {
		t.Errorf("JSON missing user role: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"content":"Hello"`) {
		t.Errorf("JSON missing content: %s", jsonStr)
	}
}

// Benchmark for Generate method (skipped if no API key)
func BenchmarkGenerate(b *testing.B) {
	apiKey := os.Getenv(testAPIKeyEnv)
	if apiKey == "" {
		b.Skipf("Skipping benchmark: %s not set", testAPIKeyEnv)
	}

	proxy := &OpenRouterCompletionProxy{
		APIKey: apiKey,
		HttpClient: &http.Client{
			Timeout: testTimeout,
		},
	}

	request := OpenRouterCompletionRequest{
		Model: testModel,
		Messages: []OpenRouterCompletionRequestMessage{
			{Role: LLMRoleUser, Content: "Say hi"},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		proxy.Generate(ctx, request)
	}
}
