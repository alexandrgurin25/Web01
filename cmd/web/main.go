package main

import (
	"AuthService/auth"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	password := os.Getenv("DB_PASSWORD") // Получаем пароль из переменной окружения
	connStr := "user=postgres password=" + password + " dbname=postgres sslmode=disable"
	// Подключение к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Используется функция http.NewServeMux() для инициализации нового рутера, затем
	// функцию "home" регистрируется как обработчик для URL-шаблона "/".
	mux := http.NewServeMux()
	mux.HandleFunc("/login",
		func(w http.ResponseWriter, r *http.Request) {
			auth.Login(w, r, db)
		})
	mux.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			home(w, r, db)
		})

	// Используется функция http.ListenAndServe() для запуска нового веб-сервера.
	// Мы передаем два параметра: TCP-адрес сети для прослушивания (в данном случае это "localhost:4000")
	// и созданный рутер. Если вызов http.ListenAndServe() возвращает ошибку
	// мы используем функцию log.Fatal() для логирования ошибок. Обратите внимание
	// что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
