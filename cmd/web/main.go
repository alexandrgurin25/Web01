package main
 
import (
    "log"
    "net/http"
)
 


func main() {
    // Используется функция http.NewServeMux() для инициализации нового рутера, затем
    // функцию "home" регистрируется как обработчик для URL-шаблона "/".
    mux := http.NewServeMux()
    mux.HandleFunc("/", home)
    mux.HandleFunc("/submit", createSnippet)
 
    // Используется функция http.ListenAndServe() для запуска нового веб-сервера. 
    // Мы передаем два параметра: TCP-адрес сети для прослушивания (в данном случае это "localhost:4000")
    // и созданный рутер. Если вызов http.ListenAndServe() возвращает ошибку
    // мы используем функцию log.Fatal() для логирования ошибок. Обратите внимание
    // что любая ошибка, возвращаемая от http.ListenAndServe(), всегда non-nil.
    log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
    err := http.ListenAndServe(":8080", mux)
    log.Fatal(err)
}