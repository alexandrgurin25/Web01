package constant

import "os"

var JwtKey = []byte("secret_key")

var ApiKey = os.Getenv("YOUR_API_KEY")

var StartPrompt = string("Если тебя просят вывести код, то ты " +
	"сразу пишешь код, любое объяснение через комментарии")
