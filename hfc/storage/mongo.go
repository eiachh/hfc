package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/eiachh/hfc/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	username       string
	password       string
	Host           string
	Port           string
	authDB         string
	CacheDatabase  string
	OFFDatabase    string
	CollectionName string

	ConnString string
	Ctx        context.Context
	Client     mongo.Client
}

func NewMongoStorage(uname string, pwd string, host string, port string, db string, authDB string) *MongoStorage {
	// Create a MongoDB URI
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", uname, pwd, host, port, db, authDB)

	ctx := context.Background()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatal(err)
	}

	return &MongoStorage{
		username:       uname,
		password:       pwd,
		Host:           host,
		Port:           port,
		authDB:         authDB,
		CacheDatabase:  db,
		CollectionName: "products",
		ConnString:     connString,
		Ctx:            ctx,
		Client:         *client,
	}
}

func (s *MongoStorage) New(prod *types.Product) bool {
	collection := s.Client.Database(s.CacheDatabase).Collection(s.CollectionName)
	_, err := collection.InsertOne(s.Ctx, prod.AsStr())
	if err != nil {
		log.Panic("Error inserting product:", err)
		return false
	}
	return true
}

// GetByBarCode searches for a product in both the cache and OFF databases based on a barcode.
// It first checks the cache database, and if no product is found, it queries the OFF database.
//
// Parameters:
// - code: The barcode string used to filter products in both databases.
//
// Return Values:
// - *types.Product: A pointer to the product if found in the cache database.
// - []byte: A JSON-encoded byte slice of the product if found in the OFF database but not in the cache.
// - nil, nil: Returned if no product is found in either database.
func (s *MongoStorage) GetByBarCode(code string) (*types.Product, []byte) {
	var (
		cacheDBResults []bson.M
		offDBResults   []bson.M
	)

	filter := bson.M{"code": code}

	// DB prep
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.CollectionName)
	offDBCollection := s.Client.Database(s.OFFDatabase).Collection(s.CollectionName)
	cacheCursor, err := cacheDBCollection.Find(s.Ctx, filter)
	offCursor, err2 := offDBCollection.Find(s.Ctx, filter)
	if err != nil || err2 != nil {
		if err != nil {
			log.Fatal(err)
		} else {
			log.Fatal(err2)
		}
	}
	defer cacheCursor.Close(s.Ctx)
	defer offCursor.Close(s.Ctx)

	// Search
	for cacheCursor.Next(s.Ctx) {
		var result bson.M
		if err := cacheCursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		cacheDBResults = append(cacheDBResults, result)
	}
	if len(cacheDBResults) < 1 {
		for offCursor.Next(s.Ctx) {
			var result bson.M
			if err := offCursor.Decode(&result); err != nil {
				log.Fatal(err)
			}
			offDBResults = append(offDBResults, result)
		}
	}
	if err := cacheCursor.Err(); err != nil {
		log.Fatal(err)
	}
	if err := offCursor.Err(); err != nil {
		log.Fatal(err)
	}

	if len(cacheDBResults) == 1 {
		var product types.Product
		jsonData, _ := json.Marshal(cacheDBResults[0])
		json.Unmarshal([]byte(jsonData), &product)
		return &product, nil
	}
	if len(offDBResults) == 1 {
		jsonData, _ := json.Marshal(offDBResults[0])
		return nil, jsonData
	}

	return nil, nil
}

func (s *MongoStorage) GenerateNew() bool {
	//TODO Generate new prod somehow
	if false {
		log.Panic("Error inserting product:")
		return false
	}
	return true
}
func (s *MongoStorage) RegisterAsMissing(barC int) bool {
	//TODO Register as missing somehow
	if false {
		log.Panic("Error inserting product:")
		return false
	}
	return true
}