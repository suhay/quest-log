package main

import (
	"os"

	"github.com/suhay/quest-log/types"

	"go.mongodb.org/mongo-driver/bson"
)

// Resolver struct
type Resolver struct{}

// Entry returns an entry by name
func (r *Resolver) Entry() *types.EntryResolver {
	collection, mctx, cancel := getCollection(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_ENTRY_COLLECTION"))
	if cancel != nil {
		defer cancel()
	}
	result := &types.Entry{}

	filter := bson.M{"name": "tutorial"}
	collection.FindOne(mctx, filter).Decode(&result)

	if result != nil {
		return &types.EntryResolver{
			E: result,
		}
	}

	return nil
}

// Thread returns a thread by name
func (r *Resolver) Thread() *types.ThreadResolver {
	collection, mctx, cancel := getCollection(os.Getenv("MONGODB_GAME_DB"), os.Getenv("MONGODB_THREAD_COLLECTION"))
	if cancel != nil {
		defer cancel()
	}
	result := &types.Thread{}

	filter := bson.M{"name": "you-spin-me-right-round"}
	collection.FindOne(mctx, filter).Decode(&result)

	if result != nil {
		return &types.ThreadResolver{
			E: result,
		}
	}

	return nil
}
