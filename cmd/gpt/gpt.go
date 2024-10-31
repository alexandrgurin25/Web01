package gpt

import (
	"AuthService/constant"
	"AuthService/entity"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	prompt := r.FormValue("question")
	requestBody, err := json.Marshal(map[string]interface{}{
		"model": "gpt-3.5-turbo", // Укажите нужную модель
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	})
	if err != nil {
		fmt.Println("Ошибка при создании JSON:", err)
		return
	}
	// Отправляем POST-запрос
	req, err := http.NewRequest("POST", "https://api.proxyapi.ru/openai/v1/chat/completions", bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		return
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+constant.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем ответ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Парсим JSON-ответ
	var completionResponse entity.ChatCompletionResponse
	if err := json.Unmarshal(body, &completionResponse); err != nil {
		fmt.Println("Ошибка при парсинге JSON:", err)
		return
	}

	// Извлекаем текст ответа
	if len(completionResponse.Choices) > 0 {
		responseContent := completionResponse.Choices[0].Message.Content
		fmt.Println("Ответ от GPT:", responseContent)
	} else {
		fmt.Println("Нет доступных ответов.")
	}
}
