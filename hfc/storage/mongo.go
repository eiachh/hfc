package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/gommon/log"

	"github.com/eiachh/hfc/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStorage struct {
	username      string
	password      string
	Host          string
	Port          string
	authDB        string
	CacheDatabase string
	OFFDatabase   string

	ProdCollName     string
	HsCollName       string
	UnrevImgCollName string

	ConnString string
	Ctx        context.Context
	Client     mongo.Client
}

func NewMongoStorage(uname string, pwd string, host string, port string, offDb string, cacheDb string, authDB string) *MongoStorage {
	log.Infof("Mongo init with Username: %s, Password: %s, Host: %s, Port: %s, OffDB: %s, CacheDB: %s, AuthDB: %s",
		uname, pwd, host, port, offDb, cacheDb, authDB)

	// Create a MongoDB URI
	connString := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s", uname, pwd, host, port, cacheDb, authDB)

	ctx := context.Background()
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatal(err)
	}

	return &MongoStorage{
		username:         uname,
		password:         pwd,
		Host:             host,
		Port:             port,
		authDB:           authDB,
		CacheDatabase:    cacheDb,
		OFFDatabase:      offDb,
		ProdCollName:     "products",
		HsCollName:       "hs",
		UnrevImgCollName: "unreviewedImg",
		ConnString:       connString,
		Ctx:              ctx,
		Client:           *client,
	}
}

func (s *MongoStorage) NewProduct(prod *types.Product) error {
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.ProdCollName)
	filter := bson.M{"code": prod.Code}

	_, err := cacheDBCollection.DeleteMany(s.Ctx, filter)
	if err != nil {
		return err
	}

	// TODO HS can save correctly this can only save str?
	_, err = cacheDBCollection.InsertOne(s.Ctx, prod)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) NewUnreviewedImg(barCode int64, imgAsBase64 string) error {
	unrevImgColl := s.Client.Database(s.CacheDatabase).Collection(s.UnrevImgCollName)
	filter := bson.M{"code": barCode}
	res, err := s.FindInCollection(filter, unrevImgColl)
	if err != nil {
		return err
	}

	if len(*res) != 0 {
		return errors.New("img already exists for this barcode")
	}

	_, err = unrevImgColl.InsertOne(s.Ctx, bson.M{"code": barCode, "imgAsBase64": imgAsBase64})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) GetUnreviewedImg(barCode int64) (string, error) {
	unrevImgColl := s.Client.Database(s.CacheDatabase).Collection(s.UnrevImgCollName)
	filter := bson.M{"code": barCode}
	res, err := s.FindInCollection(filter, unrevImgColl)
	if err != nil {
		return "", err
	}

	if len(*res) == 0 {
		return "", errors.New("img does not exists for this barcode")
	} else if len(*res) > 1 {
		return "", errors.New("multiple images are stored for a single barcode")
	}

	img64, ok := (*res)[0]["imgAsBase64"].(string)
	if !ok {
		return "", errors.New("img does not exists for this barcode")
	}
	return img64, nil
}

func (s *MongoStorage) SaveHomeStorage(hs *types.HomeStorage) error {
	hsDbCollection := s.Client.Database(s.CacheDatabase).Collection(s.HsCollName)
	filter := bson.M{}
	_, err := hsDbCollection.ReplaceOne(s.Ctx, filter, hs, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoStorage) GetAllProduct() *[]types.Product {
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.ProdCollName)
	cacheDBResults, findErr := s.FindInCollection(bson.M{}, cacheDBCollection)

	if findErr != nil {
		log.Error(findErr)
		return nil
	}

	var product []types.Product
	jsonData, _ := json.Marshal(*cacheDBResults)
	json.Unmarshal([]byte(jsonData), &product)
	return &product
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
func (s *MongoStorage) GetByBarCode(code int64) (*types.Product, []byte) {
	filterCache := bson.M{"code": code}
	filterOff := bson.M{"code": strconv.FormatInt(code, 10)}
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.ProdCollName)
	offDBCollection := s.Client.Database(s.OFFDatabase).Collection(s.ProdCollName)

	cacheDBResults, findErrCache := s.FindInCollection(filterCache, cacheDBCollection)
	if findErrCache != nil {
		log.Error(findErrCache)
		return nil, nil
	}
	if len(*cacheDBResults) == 1 {
		var product types.Product
		jsonData, _ := json.Marshal((*cacheDBResults)[0])
		json.Unmarshal([]byte(jsonData), &product)
		return &product, nil
	}

	offDBResults, findErrOff := s.FindInCollection(filterOff, offDBCollection)
	if findErrOff != nil {
		log.Error(findErrOff)
		return nil, nil
	}
	if len(*offDBResults) == 1 {
		jsonData, _ := json.Marshal((*offDBResults)[0])
		return nil, jsonData
	}
	return nil, nil
}

func (s *MongoStorage) LoadHomeStorage() (*types.HomeStorage, error) {
	var homeStorage types.HomeStorage

	hsDbCollection := s.Client.Database(s.CacheDatabase).Collection(s.HsCollName)
	err := hsDbCollection.FindOne(s.Ctx, bson.M{}).Decode(&homeStorage)
	return &homeStorage, err
}

func (s *MongoStorage) GetUnreviewedProducts() (*[]types.Product, error) {
	var unreviewed []types.Product
	filter := bson.M{"reviewed": false}
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.ProdCollName)
	primitives, err := s.FindInCollection(filter, cacheDBCollection)
	if err != nil {
		return nil, err
	}
	jsonData, _ := json.Marshal((*primitives))
	err = json.Unmarshal(jsonData, &unreviewed)
	if err != nil {
		return nil, err
	}
	return &unreviewed, nil
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

// TODO This needs to be an ordered list not just random strings
func (s *MongoStorage) GetCatListDistinct() (*[]byte, error) {
	cacheDBCollection := s.Client.Database(s.CacheDatabase).Collection(s.ProdCollName)

	distCategories, err := cacheDBCollection.Distinct(s.Ctx, "categories_hierarchy", bson.M{})
	if err != nil {

	}
	jsonCat, _ := json.Marshal(distCategories)
	return &jsonCat, nil
}
