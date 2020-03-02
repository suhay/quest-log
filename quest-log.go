package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	quest "github.com/suhay/quest-log/query"

	jwt "github.com/dgrijalva/jwt-go"
)

const defaultPort = "8080"

func main() {
	args := os.Args

	var params quest.Params
	json.Unmarshal([]byte(args[1]), &params)
	_, responseJSON, err := quest.Query(params)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"game": "yourstate",
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		jwterr := fmt.Errorf("Something Went Wrong validating the request: %s", err.Error())
		log.Println(jwterr)
		return
	}

	fmt.Printf(`{ "response": %s, "state": "%s" }`, responseJSON, tokenString)
}
