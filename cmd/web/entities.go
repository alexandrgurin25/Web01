package main

import "github.com/dgrijalva/jwt-go"

var user struct {
	login    string
	password string
}

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
