package query

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/suhay/quest-log/types"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/joho/godotenv"
	bolt "go.etcd.io/bbolt"
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
	Thread types.Thread `json:"thread,omitempty"`
	Entry  types.Entry  `json:"entry,omitempty"`
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(err)
		log.Println("Error loading .env file, defaulting to local files.")
	}

	buf, _ := ioutil.ReadFile("../schema.graphql")
	Schema := string(buf)

	graphqlSchema = graphql.MustParseSchema(Schema, resolver)

	if mongodbURL := os.Getenv("MONGODB_HOST"); len(mongodbURL) == 0 {
		db, err := bolt.Open(os.Getenv("DB_PATH"), 0600, &bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			log.Printf("create error bucket: %s", err)
		}

		localPath := os.Getenv("LOCAL_PATH")
		if db != nil && localPath != "" {
			db.Update(func(tx *bolt.Tx) error {
				var entryManifest []string
				var threadManifest []string

				manifestBkt, err := tx.CreateBucketIfNotExists([]byte("manifest"))
				if err != nil {
					panic(fmt.Errorf("create bucket: %s", err))
				}

				entryBkt, err := tx.CreateBucketIfNotExists([]byte("entry"))
				if err != nil {
					panic(fmt.Errorf("create bucket: %s", err))
				}

				entries, err := ioutil.ReadDir(localPath + "/entry")
				if err != nil {
					panic(fmt.Errorf("Error reading directory: %s", err))
				}

				entryIndexTags, err := entryBkt.CreateBucketIfNotExists([]byte("index-tags"))
				entryTags := make(map[string][]string)

				for _, file := range entries {
					if dat, err := ioutil.ReadFile(localPath + "/entry/" + file.Name()); err != nil {
						panic(fmt.Errorf("Error reading local files: %s", err))
					} else {
						v := types.Entry{}
						if err := json.Unmarshal(dat, &v); err != nil {
							panic(fmt.Errorf("Error parsing JSON file: %s", err))
						} else if err := entryBkt.Put([]byte(v.Name), dat); err != nil {
							panic(fmt.Errorf("Error saving: %s", err))
						} else {
							entryManifest = append(entryManifest, v.Name)
							for _, tag := range v.Tags {
								entryTags[*tag] = append(entryTags[*tag], v.Name)
							}
						}
					}
				}

				for tag, entryTag := range entryTags {
					entryTagBytes, _ := json.Marshal(entryTag)
					if err := entryIndexTags.Put([]byte(tag), entryTagBytes); err != nil {
						panic(fmt.Errorf("Error saving: %s", err))
					}
				}

				entryBytes, _ := json.Marshal(entryManifest)
				if err := manifestBkt.Put([]byte("entry"), entryBytes); err != nil {
					panic(fmt.Errorf("Error saving: %s", err))
				}

				/////////

				threadBkt, err := tx.CreateBucketIfNotExists([]byte("thread"))
				if err != nil {
					panic(fmt.Errorf("create bucket: %s", err))
				}

				threads, err := ioutil.ReadDir(localPath + "/thread")
				if err != nil {
					panic(fmt.Errorf("Error reading directory: %s", err))
				}

				threadTags := make(map[string][]string)

				for _, file := range threads {
					if dat, err := ioutil.ReadFile(localPath + "/thread/" + file.Name()); err != nil {
						panic(fmt.Errorf("Error reading local files: %s", err))
					} else {
						v := types.Thread{}
						if err := json.Unmarshal(dat, &v); err != nil {
							panic(fmt.Errorf("Error parsing JSON file: %s", err))
						} else if err := threadBkt.Put([]byte(v.Name), dat); err != nil {
							panic(fmt.Errorf("Error saving: %s", err))
						} else {
							threadManifest = append(threadManifest, v.Name)
							for _, tag := range v.Tags {
								threadTags[*tag] = append(threadTags[*tag], v.Name)
							}
						}
					}
				}

				threadIndexTags, err := threadBkt.CreateBucketIfNotExists([]byte("index-tags"))

				for tag, threadTag := range threadTags {
					threadTagBytes, _ := json.Marshal(threadTag)
					if err := threadIndexTags.Put([]byte(tag), threadTagBytes); err != nil {
						panic(fmt.Errorf("Error saving: %s", err))
					}
				}

				if threadBytes, err := json.Marshal(threadManifest); err != nil {
					panic(fmt.Errorf("Error saving: %s", err))
				} else if err := manifestBkt.Put([]byte("thread"), threadBytes); err != nil {
					panic(fmt.Errorf("Error saving: %s", err))
				}

				return nil
			})
		}

		defer db.Close()
	}
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
