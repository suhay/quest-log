package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	graphql "github.com/graph-gophers/graphql-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/joho/godotenv"
  jwt "github.com/dgrijalva/jwt-go"
)

const defaultPort = "8080"

var graphqlSchema *graphql.Schema

func getCollection(db string, col string) (*mongo.Collection, context.Context, context.CancelFunc) {
	if mongodbURL := os.Getenv("MONGODB_HOST"); len(mongodbURL) > 0 {
		mctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		log.Println("Connecting to MongoDB - Quest Log")
		client, err := mongo.Connect(mctx, options.Client().SetRetryWrites(true).ApplyURI("mongodb+srv://"+os.Getenv("MONGODB_USERNAME")+":"+os.Getenv("MONGODB_PASSWORD")+"@"+mongodbURL+"/"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connected")

		collection := client.Database(db).Collection(col)

		return collection, mctx, cancel
	}

	panic("No mongodb connection information provided.")
}

func init() {
	buf, _ := ioutil.ReadFile("schema.graphql")
	Schema := string(buf)
	graphqlSchema = graphql.MustParseSchema(Schema, &Resolver{})
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, defaulting to local files.")
	}

	args := os.Args

	var params struct {
		Query         string
		OperationName string
		Variables     map[string]interface{}
	}

	json.Unmarshal([]byte(args[1]), &params)

	response := graphqlSchema.Exec(context.Background(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return
	}

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
