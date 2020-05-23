package query

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"sort"
	"time"

	bolt "go.etcd.io/bbolt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Filter struct
type Filter struct {
	Name *string
	Tags *[]*string
}

// Collection struct
type Collection struct {
	Kind    string
	Bolt    *bolt.DB
	Mongo   *mongo.Collection
	Context context.Context
	Cancel  context.CancelFunc
}

// NewCollection returns a Collection struct and connects it to the DB based on 'kind'
func NewCollection(kind string) Collection {
	if mongodbURL := os.Getenv("MONGODB_HOST"); len(mongodbURL) > 0 {
		mctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		log.Println("Connecting to MongoDB - Quest Log")
		client, err := mongo.Connect(mctx, options.Client().SetRetryWrites(true).ApplyURI("mongodb+srv://"+os.Getenv("MONGODB_USERNAME")+":"+os.Getenv("MONGODB_PASSWORD")+"@"+mongodbURL+"/"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connected")

		col := Collection{
			Mongo:   client.Database(os.Getenv("MONGODB_GAME_DB")).Collection(kind),
			Context: mctx,
			Cancel:  cancel,
			Kind:    kind,
		}

		return col
	}

	if localPath := os.Getenv("DB_PATH"); len(localPath) > 0 {
		db, _ := bolt.Open(os.Getenv("DB_PATH"), 0600, &bolt.Options{ReadOnly: true})

		col := Collection{
			Bolt: db,
			Kind: kind,
		}

		return col
	}

	panic("No db connection information provided (in-memory or Mongo).")
}

// Query runs the supplied filter against the configured storage
func (c Collection) Query(filter *Filter, result interface{}) error {
	if ok := filter.Name; ok != nil {
		return c.byName(filter, result)
	} else if ok := filter.Tags; ok != nil {
		return c.byHighestTag(filter, result)
	} else {
		return c.random(result)
	}
}

func (c Collection) random(result interface{}) error {
	if c.Mongo != nil {
		pipeline := []bson.M{bson.M{"$sample": bson.M{"size": 1}}}
		cursor, err := c.Mongo.Aggregate(c.Context, pipeline)
		defer cursor.Close(c.Context)
		if err != nil {
			return err
		}

		if cursor.Next(c.Context) {
			return cursor.Decode(&result)
		}
	} else if c.Bolt != nil {
		defer c.Bolt.Close()

		return c.Bolt.View(func(tx *bolt.Tx) error {
			manifest := tx.Bucket([]byte("manifest"))
			typeManifest := manifest.Get([]byte(c.Kind))

			var typeManifestSlice []string
			json.Unmarshal(typeManifest, &typeManifestSlice)

			count := len(typeManifestSlice)
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)

			b := tx.Bucket([]byte(c.Kind))
			get := b.Get([]byte(typeManifestSlice[r1.Intn(count)]))
			json.Unmarshal(get, &result)

			return nil
		})
	}

	return nil
}

func (c Collection) byName(filter *Filter, result interface{}) error {
	if c.Mongo != nil {
		return c.Mongo.FindOne(c.Context, filter).Decode(&result)
	} else if c.Bolt != nil {
		defer c.Bolt.Close()

		return c.Bolt.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(c.Kind))
			get := b.Get([]byte(*filter.Name))
			json.Unmarshal(get, &result)

			return nil
		})
	}

	return nil
}

func (c Collection) byHighestTag(filter *Filter, result interface{}) error {
	tags := *filter.Tags

	if c.Mongo != nil {
		return c.Mongo.FindOne(c.Context, filter).Decode(&result)
	} else if c.Bolt != nil {
		defer c.Bolt.Close()

		return c.Bolt.View(func(tx *bolt.Tx) error {
			tagBucket := make(map[string]int)
			kind := tx.Bucket([]byte(c.Kind))
			indexTags := kind.Bucket([]byte("index-tags"))

			for _, tag := range tags {
				var names []string
				byteTag := indexTags.Get([]byte(*tag))
				json.Unmarshal(byteTag, &names)

				for _, name := range names {
					// TODO: Add in a bias multiplier brought in from the game state. Can be greater than, equal to, or less than 1.
					if val, ok := tagBucket[name]; ok {
						tagBucket[name] = val + 1
					} else {
						tagBucket[name] = 1
					}
				}
			}

			type kv struct {
				Key   string
				Value int
			}

			var ss []kv
			for k, v := range tagBucket {
				ss = append(ss, kv{k, v})
			}

			sort.Slice(ss, func(i, j int) bool {
				return ss[i].Value > ss[j].Value
			})

			var highest []kv
			for _, v := range ss {
				if len(highest) == 0 || (len(highest) > 0 && v.Value == highest[0].Value) {
					highest = append(highest, v)
				} else {
					break
				}
			}

			count := len(highest)
			s1 := rand.NewSource(time.Now().UnixNano())
			r1 := rand.New(s1)

			get := kind.Get([]byte(highest[r1.Intn(count)].Key))
			json.Unmarshal(get, &result)

			return nil
		})
	}

	return nil
}
