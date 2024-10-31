package constant

import "os"

var JwtKey = []byte("secret_key")

var ApiKey = os.Getenv("YOUR_API_KEY")
