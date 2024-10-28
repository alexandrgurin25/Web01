package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Обработчик главной страницы.
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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

var jwtKey = []byte("secret_key")

func login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
		return
	}

	if r.Method == http.MethodGet {
		// Отображаем страницу входа
		ts, err := template.ParseFiles("../../ui/html/login.page.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		// Обрабатываем логику входа
		err := db.QueryRow(
			`SELECT login, password FROM users WHERE login = $1 AND password = $2`,
			r.FormValue("username"),
			r.FormValue("password"),
		).Scan(&user.login, &user.password)

		if err != nil {
			if err == sql.ErrNoRows {
				// Пользователь не найден
				http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				return
			}
			// Обработка других ошибок
			http.Error(w, "Ошибка при выполнении запроса", http.StatusInternalServerError)
			return
		}

		// Создание токена
		expirationTime := time.Now().Add(24 * time.Hour) // Токен будет действителен 24 часа
		claims := &UserClaims{
			Username: user.login,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
			return
		}

		// Запись токена в куки
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    tokenString,
			Expires:  expirationTime,
			Path:     "/",
			HttpOnly: true, // Защита от XSS
		})

		// Логика успешного входа (например, установка сессии, перенаправление и т.д.)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
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
	claims := &UserClaims{}

	// Проверка и парсинг токена
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Необходима аутентификация", http.StatusUnauthorized)
		return
	}

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
