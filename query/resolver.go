package query

import (
	"log"

	"github.com/suhay/quest-log/types"
)

// Resolver struct
type Resolver struct{}

func process(kind string, filter *Filter, result interface{}) {
	collection := NewCollection(kind)
	if collection.Cancel != nil {
		defer collection.Cancel()
	}

	if err := collection.Query(filter, &result); err != nil {
		log.Fatal(err)
	}
}

// Entry returns a random entry
func (*Resolver) Entry(filter *Filter) *types.EntryResolver {
	var result types.Entry
	process("entry", filter, &result)

	return &types.EntryResolver{
		Entry: &result,
	}
}

// Thread returns a random thread
func (*Resolver) Thread(filter *Filter) *types.ThreadResolver {
	var result types.Thread
	process("thread", filter, &result)

	return &types.ThreadResolver{
		Thread: &result,
	}
}
