package mongo

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	templete = "templete"
	orders   = "orders"
)

type mongoDatabase struct {
	client      *mongo.Client
	database    *mongo.Database
	_collection sync.Map
}

func (db *mongoDatabase) TempleteCollection() *mongo.Collection { return db.collection(templete) }
func (db *mongoDatabase) OrdersCollection() *mongo.Collection   { return db.collection(orders) }

func NewDatabase(connectionString string, databaseName string) Database {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Printf("Failed to connect mongo db: %v", err)
		return nil
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("Ping failed to mongo db: %v", err)
		return nil
	}

	md := &mongoDatabase{
		client:   client,
		database: client.Database(databaseName),
	}

	{
		collectionNames, _ := md.database.ListCollectionNames(context.TODO(), bson.M{})
		for _, name := range collectionNames {
			md._collection.Store(name, true)
		}
	}

	md.createIndexes()

	return md
}

func (db *mongoDatabase) Disconnect() error {
	return db.client.Disconnect(context.TODO())
}

func (db *mongoDatabase) collection(name string) *mongo.Collection {
	if _, exists := db._collection.Load(name); exists {
		return db.database.Collection(name)
	}

	db._collection.Store(name, true)
	_ = db.database.CreateCollection(context.TODO(), name)

	return db.database.Collection(name)
}

func (db *mongoDatabase) createIndexes() (err error) {
	ordersCol := db.OrdersCollection()
	_, _ = ordersCol.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "account_id", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	return
}

type mongoIndex struct {
	Name string         `json:"name"`
	Key  map[string]int `json:"key"`
}
