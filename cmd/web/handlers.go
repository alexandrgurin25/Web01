package main

import (
	"AuthService/constant"
	"AuthService/entity"
	"AuthService/gpt"
	"database/sql"
	"log"
	"net/http"
	"text/template"

	"github.com/dgrijalva/jwt-go"
)

// Обработчик главной страницы.
func home(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Проверка токена
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	tokenStr := cookie.Value
	claims := &entity.UserClaims{}

	// Проверка и парсинг токена
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return constant.JwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)
		return
	}

	if r.Method == http.MethodGet {
		// Обработка GET-запроса для отображения главной страницы
		ts, err := template.ParseFiles("../../ui/html/home.page.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		// Обработка POST-запроса для получения вопроса
		question := r.FormValue("question")
		if question == "" {
			http.Error(w, "Вопрос не может быть пустым", http.StatusBadRequest)
			return
		}

		// Вызываем обработчик GPT для получения ответа
		responseContent, err := gpt.GetResponseFromGPT(question)
		if err != nil {
			http.Error(w, "Ошибка при получении ответа от GPT", http.StatusInternalServerError)
			return
		}

		// Сохраняем вопрос и ответ в базе данных
		if len(responseContent) > 50 {
			_, err = db.Exec(`INSERT INTO questions (login, question, answer) VALUES ($1, $2, $3)`, claims.Username, question, "Большой ответ...")

		} else {
			_, err = db.Exec(`INSERT INTO questions (login, question, answer) VALUES ($1, $2, $3)`, claims.Username, question, responseContent)

		}
		
		if err != nil {
			log.Println("Ошибка при выполнении запроса:", err)
			http.Error(w, "Не удалось выполнить запрос", http.StatusInternalServerError)
			return
		}

		// Возвращаем ответ пользователю
		w.Write([]byte(responseContent))
		return
	}

	// Если метод не GET и не POST, возвращаем ошибку
	http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
}
