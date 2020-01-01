package main

import (
	"context"
  "log"
  "net/http"
  "os"
	"time"
  "flag"
  "io/ioutil"

  "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
  "github.com/go-chi/chi"
  "github.com/go-chi/cors"
  "github.com/go-chi/chi/middleware"
  "github.com/joho/godotenv"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "github.com/friendsofgo/graphiql"
)

const defaultPort = "8080"
var graphqlSchema *graphql.Schema

func setPort(port1, port2 string) string {
	if port1 != "" {
		return port1
	} else if port2 != "" {
		return port2
  }

	return defaultPort
}

func getCollection(db string, col string) (*mongo.Collection, context.Context) {
	if mongodbURL := os.Getenv("MONGODB_HOST"); len(mongodbURL) > 0 {
		mctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

		client, err := mongo.Connect(mctx, options.Client().SetRetryWrites(true).ApplyURI("mongodb+srv://"+os.Getenv("MONGODB_USERNAME")+":"+os.Getenv("MONGODB_PASSWORD")+"@"+mongodbURL+"/"))
		if err != nil {
			log.Fatal(err)
		}

    collection := client.Database(db).Collection(col)

    log.Println("Connecting to MongoDB")
    
		return collection, mctx
	}

	panic("No mongodb connection information provided.")
}

func init() {
  buf, _ := ioutil.ReadFile("schema.graphql")
  Schema := string(buf)
	graphqlSchema = graphql.MustParseSchema(Schema, &Resolver{})
}

func main() {
  flagEnvPath := flag.String("env", "", "Path to .env file to use")
	flagPort := flag.String("port", "", "Port for the quest log to run on")

	flag.Parse()

	if *flagEnvPath != "" {
		if err := godotenv.Load(*flagEnvPath); err != nil {
			log.Println("Error loading .env file, defaulting to local files.")
		}
	} else {
		if err := godotenv.Load(); err != nil {
			log.Println("Error loading .env file, defaulting to local files.")
		}
	}

	port := setPort(*flagPort, os.Getenv("PORT"))

  r := chi.NewRouter()

  cors := cors.New(cors.Options{
    // AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
    AllowedOrigins:   []string{"*"},
    // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
    ExposedHeaders:   []string{"Link"},
    AllowCredentials: true,
    MaxAge:           300, // Maximum value not ignored by any of major browsers
  })

  r.Use(middleware.RequestID)
  // r.Use(middleware.Logger)
  r.Use(middleware.Recoverer)
  r.Use(middleware.Timeout(60 * time.Second))
  r.Use(middleware.ThrottleBacklog(2, 5, time.Second*61))
  r.Use(cors.Handler)

  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to the Quest Log!"))
  })

  graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
  if err != nil {
		panic(err)
  }
  
  r.Handle("/graphiql", graphiqlHandler)
  r.Handle("/query", &relay.Handler{Schema: graphqlSchema})
  
  log.Printf("connect to http://localhost:%s/ for GraphQL playground. PID: %d", port, os.Getpid())
  log.Fatal(http.ListenAndServe(":"+port, r))
}