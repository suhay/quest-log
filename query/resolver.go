package query

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/suhay/quest-log/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

// Resolver struct
type Resolver struct{}

// EntryArgs struct
type EntryArgs struct {
	Name string
}

// Entry returns an entry by name
func (*Resolver) Entry() *types.EntryResolver {
	collection, mctx, cancel := getCollection(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_ENTRY_COLLECTION"))
	if cancel != nil {
		defer cancel()
	}
	result := &types.Entry{}

	filter := bson.M{"name": "tutorial"}
	collection.FindOne(mctx, filter).Decode(&result)

	if result != nil {
		return &types.EntryResolver{
			Entry: result,
		}
	}

	return nil
}

// Thread returns a random thread
func (*Resolver) Thread() *types.ThreadResolver {
	collection, mctx, cancel := getCollection(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_THREAD_COLLECTION"))
	if cancel != nil {
		defer cancel()
	}
	result := &types.Thread{}

	filter := bson.M{"name": "You spin me right round"}
	collection.FindOne(mctx, filter).Decode(&result)

	if result != nil {
		return &types.ThreadResolver{
			Thread: result,
		}
	}

	return nil
}

// GetThread returns a thread by name
func (*Resolver) GetThread(args EntryArgs) *types.ThreadResolver {
	collection, mctx, cancel := getCollection(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_THREAD_COLLECTION"))
	if cancel != nil {
		defer cancel()
	}
	result := &types.Thread{}

	filter := bson.M{"name": args.Name}
	collection.FindOne(mctx, filter).Decode(&result)

	if result != nil {
		return &types.ThreadResolver{
			Thread: result,
		}
	}

	return nil
}
