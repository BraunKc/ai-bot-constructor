package openrouter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	chatdomain "github.com/braunkc/ai-bot-constructor/bot-container/internal/domain/chat"
)

type OpenRouter struct {
	authToken string
	model     string
	log       *slog.Logger
}

func New(authToken, model string, log *slog.Logger) (*OpenRouter, error) {
	if authToken == "" {
		return nil, errors.New("empty auth_token")
	}

	return &OpenRouter{
		authToken: authToken,
		model:     model,
		log:       log,
	}, nil
}

func (or *OpenRouter) CreateResponse(messages []chatdomain.Message) (*chatdomain.Message, error) {
	or.log.Debug("creating open router response")
	reqBody := struct {
		Messages []chatdomain.Message `json:"input"`
		Model    string               `json:"model"`
	}{
		Messages: messages,
		Model:    or.model,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal history: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/responses",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", or.authToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}

	defer resp.Body.Close()

	or.log.Debug("received resp", slog.Any("resp", resp))

	var respBody struct {
		Output []struct {
			Role    string `json:"role"`
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal resp: %w", err)
	}

	if err := json.Unmarshal(respBytes, &respBody); err != nil {
		return nil, fmt.Errorf("failed to unmarshal resp: %w", err)
	}

	if len(respBody.Output) < 1 || len(respBody.Output[0].Content) < 1 {
		return nil, errors.New("failed to parse resp")
	}

	return &chatdomain.Message{
		Role:    chatdomain.RoleAssistant,
		Content: respBody.Output[0].Content[0].Text,
	}, nil
}
