package query

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/suhay/quest-log/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection struct
type Collection struct {
	Threads []types.Thread `json:"threads,omitempty"`
	Entries []types.Entry  `json:"entries,omitempty"`
	Mongo   *mongo.Collection
}

func getCollection(db string, col string) (*Collection, context.Context, context.CancelFunc) {
	var collection *Collection

	if localPath := os.Getenv("LOCAL_PATH"); len(localPath) > 0 {
		if col == "thread" {
			dat, err := ioutil.ReadFile(localPath + "/" + col + ".json")
			if err != nil {
				log.Fatal(err)
			}

			json.Unmarshal([]byte(dat), &collection)
		} else {
			dat, err := ioutil.ReadFile(localPath + "/" + col + ".json")
			if err != nil {
				log.Fatal(err)
			}

			json.Unmarshal([]byte(dat), &collection)
		}

		return collection, nil, nil
	}

	if mongodbURL := os.Getenv("MONGODB_HOST"); len(mongodbURL) > 0 {
		mctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		log.Println("Connecting to MongoDB - Quest Log")
		client, err := mongo.Connect(mctx, options.Client().SetRetryWrites(true).ApplyURI("mongodb+srv://"+os.Getenv("MONGODB_USERNAME")+":"+os.Getenv("MONGODB_PASSWORD")+"@"+mongodbURL+"/"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connected")

		collection = &Collection{
			Mongo: client.Database(db).Collection(col),
		}
		return collection, mctx, cancel
	}

	panic("No mongodb connection information provided.")
}

func process(db string, coll string, args EntryArgs, v interface{}) interface{} {
	collection, mctx, cancel := getCollection(db, coll)
	if cancel != nil {
		defer cancel()
	}

	if collection.Mongo != nil {
		filter := make(map[string]interface{})
		if len(args.Name) > 0 {
			filter["name"] = args.Name
		}
		if len(filter) > 0 {
			collection.Mongo.FindOne(mctx, filter).Decode(v)
		} else {
			pipeline := []bson.M{bson.M{"$sample": bson.M{"size": 1}}}
			cursor, err := collection.Mongo.Aggregate(mctx, pipeline)
			if err != nil {
				log.Println("Finding random document ERROR:", err)
				defer cursor.Close(mctx)
			} else {
				if cursor.Next(mctx) {
					cursor.Decode(v)
				}
			}
		}
	} else if len(collection.Threads) > 0 {
		if args == (EntryArgs{}) {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			k := collection.Threads[r1.Intn(len(collection.Threads))]
			return k
		}

		for i := 0; i < len(collection.Threads); i++ {
			if len(args.Name) > 0 && collection.Threads[i].Name == args.Name {
				k := collection.Threads[i]
				return k
			}
		}
	} else if len(collection.Entries) > 0 {
		if args == (EntryArgs{}) {
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)
			k := collection.Entries[r1.Intn(len(collection.Entries))]
			return k
		}

		for i := 0; i < len(collection.Entries); i++ {
			if len(args.Name) > 0 && collection.Entries[i].Name == args.Name {
				k := collection.Entries[i]
				return k
			}
		}
	}
	return nil
}

// Resolver struct
type Resolver struct{}

// EntryArgs struct
type EntryArgs struct {
	Name string
}

// Entry returns an entry by name
func (*Resolver) Entry() *types.EntryResolver {
	var result types.Entry
	v := process(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_ENTRY_COLLECTION"), EntryArgs{Name: "tutorial"}, &result)
	if v != nil {
		result = v.(types.Entry)
	}

	return &types.EntryResolver{
		Entry: &result,
	}
}

// Thread returns a random thread
func (*Resolver) Thread() *types.ThreadResolver {
	var result types.Thread
	v := process(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_THREAD_COLLECTION"), EntryArgs{}, &result)
	if v != nil {
		result = v.(types.Thread)
	}

	return &types.ThreadResolver{
		Thread: &result,
	}
}

// GetThread returns a thread by name
func (*Resolver) GetThread(filter EntryArgs) *types.ThreadResolver {
	var result types.Thread
	v := process(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_THREAD_COLLECTION"), filter, &result)
	if v != nil {
		result = v.(types.Thread)
	}

	return &types.ThreadResolver{
		Thread: &result,
	}
}
