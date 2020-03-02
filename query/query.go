package query

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/suhay/quest-log/types"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/joho/godotenv"
)

var graphqlSchema *graphql.Schema
var resolver *Resolver

// Params is a graphql query struct
type Params struct {
	Query         string
	OperationName string
	Variables     map[string]interface{}
}

// Results is a Query result
type Results struct {
	Thread    types.Thread `json:"thread,omitempty"`
	GetThread types.Thread `json:"GetThread,omitempty"`
	Entry     types.Entry  `json:"entry,omitempty"`
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(err)
		log.Println("Error loading .env file, defaulting to local files.")
	}

	buf, _ := ioutil.ReadFile("../schema.graphql")
	Schema := string(buf)

	graphqlSchema = graphql.MustParseSchema(Schema, resolver)
}

// Query returns a graphql resolved entry
func Query(params Params) (Results, *graphql.Response, error) {
	response := graphqlSchema.Exec(context.Background(), params.Query, params.OperationName, params.Variables)

	var result Results
	err := json.Unmarshal(response.Data, &result)
	if err != nil {
		return result, nil, err
	}

	return result, response, nil
}
