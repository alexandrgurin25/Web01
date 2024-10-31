package gpt

import (
	"AuthService/constant"
	"AuthService/entity"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetResponseFromGPT(prompt string) (string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": *&constant.StartPrompt},
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		return "", fmt.Errorf("ошибка при создании JSON: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.proxyapi.ru/openai/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("ошибка при создании запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+constant.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка при чтении ответа: %w", err)
	}

	var completionResponse entity.ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResponse); err != nil {
		return "", fmt.Errorf("ошибка при парсинге JSON: %w", err)
	}

	if len(completionResponse.Choices) > 0 {
		return completionResponse.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("нет доступных ответов")
}
