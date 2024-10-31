package entity

import "github.com/dgrijalva/jwt-go"

var User struct {
	Login    string
	Password string
}

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}