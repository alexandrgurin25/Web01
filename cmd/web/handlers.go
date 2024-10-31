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
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	_, err := r.Cookie("access_token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Используем функцию template.ParseFiles() для чтения файла шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles("../../ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func writeInDb(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("GET-Метод запрещен!"))
		return
	}

	// Проверка токена
	cookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)

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
	gpt.Handler(w, r, db)
	// Теперь мы можем использовать claims.Username для получения имени пользователя
	question := r.FormValue("question")
	if question == "" {
		http.Error(w, "Вопрос не может быть пустым", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`INSERT INTO questions (login, question) VALUES ($1, $2)`, claims.Username, question)
	if err != nil {
		log.Println("Ошибка при выполнении запроса:", err)
		http.Error(w, "Не удалось выполнить запрос", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Данные успешно добавлены"))
}
