package auth

import (
	"AuthService/entity"
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

func Login(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		if r.FormValue("username") == "" || r.FormValue("password") == "" {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			log.Println("Пустые поля авторизации", r.Form)
			return
		}

		// Обрабатываем логику входа
		err := db.QueryRow(
			`SELECT login, password FROM users WHERE login = $1 AND password = $2`,
			r.FormValue("username"),
			r.FormValue("password"),
		).Scan(&entity.User.Login, &entity.User.Password)

		if err != nil {
			if err == sql.ErrNoRows {
				// Пользователь не найден
				http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
				return
			}
			// Обработка других ошибок
			log.Println(err)
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Создание токена
		expirationTime := time.Now().Add(24 * time.Hour) // Токен будет действителен 24 часа
		claims := &entity.UserClaims{
			Username: entity.User.Login,
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
