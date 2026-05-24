package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	database "saidakbar.origin/db/mongo"
	"saidakbar.origin/db/mongo/entity"
	"saidakbar.origin/repository/models"
)

type OrderRepository interface {
	Create(order *entity.OrderEntity) error
	GetByID(id string) (*entity.OrderEntity, error)
	GetAll(filters *models.GetAllOrdersFilter) ([]*entity.OrderEntity, error)
	Update(order *entity.OrderEntity) error
	UpdateStatus(id, status string) error
	IncrementViewCount(id string) error
}

type orderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(db database.Database) OrderRepository {
	return &orderRepository{collection: db.OrdersCollection()}
}

func (r *orderRepository) Create(order *entity.OrderEntity) error {
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}

func (r *orderRepository) GetByID(id string) (*entity.OrderEntity, error) {
	var order entity.OrderEntity
	err := r.collection.FindOne(context.Background(), bson.M{entity.FieldOrderID: id}).Decode(&order)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) GetAll(filters *models.GetAllOrdersFilter) ([]*entity.OrderEntity, error) {
	filter := bson.M{entity.FieldOrderAccountID: filters.AccountID}

	if filters.Status != nil {
		filter[entity.FieldOrderStatus] = *filters.Status
	}
	if filters.TemplateCode != nil {
		filter[entity.FieldOrderTemplateCode] = *filters.TemplateCode
	}
	if filters.FromDate != nil || filters.ToDate != nil {
		dateFilter := bson.M{}
		if filters.FromDate != nil {
			dateFilter["$gte"] = *filters.FromDate
		}
		if filters.ToDate != nil {
			dateFilter["$lte"] = *filters.ToDate
		}
		filter[entity.FieldOrderCreatedAt] = dateFilter
	}

	sortField := entity.FieldOrderCreatedAt
	if filters.SortBy != nil && *filters.SortBy != "" {
		sortField = *filters.SortBy
	}
	sortDir := -1
	if filters.Order != nil && *filters.Order == "asc" {
		sortDir = 1
	}

	findOpts := options.Find().SetSort(bson.D{{Key: sortField, Value: sortDir}})
	if filters.Page != nil && filters.Limit != nil && *filters.Limit > 0 {
		skip := int64((*filters.Page - 1) * *filters.Limit)
		findOpts.SetSkip(skip).SetLimit(int64(*filters.Limit))
	} else if filters.Limit != nil && *filters.Limit > 0 {
		findOpts.SetLimit(int64(*filters.Limit))
	}

	cursor, err := r.collection.Find(context.Background(), filter, findOpts)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()

	var results []*entity.OrderEntity
	if err = cursor.All(context.Background(), &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *orderRepository) Update(order *entity.OrderEntity) error {
	_, err := r.collection.ReplaceOne(
		context.Background(),
		bson.M{entity.FieldOrderID: order.ID},
		order,
	)
	return err
}

func (r *orderRepository) UpdateStatus(id, status string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{entity.FieldOrderID: id},
		bson.M{"$set": bson.M{entity.FieldOrderStatus: status, entity.FieldOrderUpdatedAt: bson.M{"$currentDate": true}}},
	)
	return err
}

func (r *orderRepository) IncrementViewCount(id string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{entity.FieldOrderID: id},
		bson.M{"$inc": bson.M{entity.FieldOrderViewCount: 1}},
	)
	return err
}
