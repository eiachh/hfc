package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/eiachh/hfc/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func NewMongoStorage(uname string, pwd string, host string, port string, offDb string, cacheDb string, authDB string) *MongoStorage {
	// Create a MongoDB URI
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", uname, pwd, host, port, cacheDb, authDB)

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
		CacheDatabase:  cacheDb,
		OFFDatabase:    offDb,
		CollectionName: "products",
		ConnString:     connString,
		Ctx:            ctx,
		Client:         *client,
	}
}

func (s *MongoStorage) NewProduct(prod *types.Product) error {
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.CollectionName)
	filter := bson.M{"code": prod.Code}

	res, err := s.FindInCollection(filter, cacheDBCollection)
	if err != nil {
		return err
	} else if len(*res) > 0 {
		return errors.New("item already exists in loc-cache")
	}

	_, err = cacheDBCollection.InsertOne(s.Ctx, prod.AsStr())
	if err != nil {
		return err
	}
	return nil
}

// GetByBarCode searches for a product in both the cache and OFF databases based on a barcode.
// It first checks the cache database, and if no product is found, it queries the OFF database.
//
// Parameters:
// - code: The barcode int used to filter products in both databases.
//
// Return Values:
// - *types.Product: A pointer to the product if found in the cache database.
// - []byte: A JSON-encoded byte slice of the product if found in the OFF database but not in the cache.
// - nil, nil: Returned if no product is found in either database.
func (s *MongoStorage) GetByBarCode(code int) (*types.Product, []byte) {
	var (
		cacheDBResults []bson.M
		offDBResults   []bson.M
	)
	codeStr := strconv.Itoa(code)

	filter := bson.M{"code": codeStr}

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

// FindInCollection finds documents in a MongoDB collection that match the given filter criteria.
//
// Parameters:
//   - filter: A MongoDB filter used to specify the search criteria (type: primitive.M).
//   - coll: A pointer to the MongoDB collection where the search will be performed.
//
// Returns:
//   - *[]primitive.M: A pointer to a slice of documents (in BSON format) that match the filter criteria.
//   - error: An error object if any issue occurs during the operation.
//
// The function executes a query using the provided filter on the specified collection, iterates through the results,
// decodes each document into a BSON map, and appends it to a slice. If any error occurs during querying or decoding,
// the function returns the error.
func (s *MongoStorage) FindInCollection(filter primitive.M, coll *mongo.Collection) (*[]primitive.M, error) {
	var cacheDBResults []bson.M

	cacheCursor, err := coll.Find(s.Ctx, filter)
	if err != nil {
		return nil, err
	}

	coll.Find(s.Ctx, filter)
	defer cacheCursor.Close(s.Ctx)

	for cacheCursor.Next(s.Ctx) {
		var result bson.M
		if err := cacheCursor.Decode(&result); err != nil {
			return nil, err
		}
		cacheDBResults = append(cacheDBResults, result)
	}

	if err := cacheCursor.Err(); err != nil {
		return nil, err
	}

	return &cacheDBResults, nil
}
