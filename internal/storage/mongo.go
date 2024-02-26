package storage

import (
	"app/main/utils"
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoService struct {
	cli *mongo.Client

	isInitialized bool
}

var mongoService MongoService

func InitMongo() error {

	env := utils.Env()
	uri := fmt.Sprintf("mongodb://%s:%s@%s/?maxPoolSize=%s",
		env.Value("MONGO_INITDB_ROOT_USERNAME"),
		env.Value("MONGO_INITDB_ROOT_PASSWORD"),
		env.Value("MONGO_DB_HOST"),
		env.Value("MONGO_INITDB_MAX_POOL_SIZE"),
	)

	if !mongoService.isInitialized {

		client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
		if err != nil {
			return err
		}

		mongoService.cli = client
		mongoService.isInitialized = true
	}

	return nil
}

func Instance() *MongoService {

	if !mongoService.isInitialized {
		return nil
	}
	return &mongoService
}

func (s *MongoService) ReadValue(connection string, key string, value string) (string, error) {

	env := utils.Env()

	coll := s.cli.Database(env.Value("MONGO_INITDB_DATABASE")).Collection(connection)

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: key, Value: value}}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return fmt.Sprintf("No document was found with the title %s\n", key), err
	}

	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s\n", jsonData), nil
}

func (s *MongoService) WriteValue(connection string, key string, value string) (string, error) {

	return "", nil
}
