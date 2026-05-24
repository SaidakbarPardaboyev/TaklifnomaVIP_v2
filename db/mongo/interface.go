package mongo

import "go.mongodb.org/mongo-driver/mongo"

type Database interface {
	Disconnect() error
	TempleteCollection() *mongo.Collection
	OrdersCollection() *mongo.Collection
}
