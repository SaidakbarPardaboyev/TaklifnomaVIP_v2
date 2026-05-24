package mongo

import "go.mongodb.org/mongo-driver/mongo"

type Database interface {
	Disconnect() error
	TempleteCollection() *mongo.Collection
	// TODO: add one method per collection your service owns
}
