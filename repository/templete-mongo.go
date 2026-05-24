package repository

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	database "saidakbar.origin/db/mongo"
	"saidakbar.origin/db/mongo/entity"
	"saidakbar.origin/repository/models"
)

type TempleteMongoRepository interface {
	CreateTemplete(templete *entity.TempleteEntity) (err error)
	UpdateTemplete(templete *entity.TempleteEntity) (err error)
	DeleteTemplete(templete *entity.TempleteEntity) (err error)
	GetTemplete(filter bson.M) (templete *entity.TempleteEntity, err error)
	GetAllTempletes(filters *models.GetAllTempleteFilter) (templetes []*entity.TempleteEntity, err error)
	GetTempleteCount(filters *models.GetTempleteCountFilter) (count int64, err error)
}

type templeteMongoRepository struct {
	collection *mongo.Collection
}

func NewTempleteMongoRepository(db database.Database) TempleteMongoRepository {
	return &templeteMongoRepository{
		collection: db.TempleteCollection(),
	}
}

func (r *templeteMongoRepository) CreateTemplete(templete *entity.TempleteEntity) (err error) {
	_, err = r.collection.InsertOne(context.Background(), templete)
	return
}

func (r *templeteMongoRepository) UpdateTemplete(templete *entity.TempleteEntity) (err error) {
	filter := bson.M{
		entity.FieldTempleteID: templete.ID,
	}
	update := bson.M{
		"$set": bson.M{
			entity.FieldTempleteName:      templete.Name,
			entity.FieldTempleteUpdatedAt: templete.UpdatedAt,
			// TODO: add more fields to update
		},
	}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return
}

func (r *templeteMongoRepository) DeleteTemplete(templete *entity.TempleteEntity) (err error) {
	filter := bson.M{
		entity.FieldTempleteID: templete.ID,
	}
	update := bson.M{
		"$set": bson.M{
			entity.FieldTempleteIsDeleted: true,
			entity.FieldTempleteDeletedAt: templete.DeletedAt,
		},
	}
	_, err = r.collection.UpdateOne(context.Background(), filter, update)
	return
}

func (r *templeteMongoRepository) GetTemplete(filter bson.M) (templete *entity.TempleteEntity, err error) {
	if err = r.collection.FindOne(context.Background(), filter).Decode(&templete); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return
}

func (r *templeteMongoRepository) GetAllTempletes(filters *models.GetAllTempleteFilter) (templetes []*entity.TempleteEntity, err error) {
	filter := bson.M{
		entity.FieldTempleteOrganizationID: filters.OrganizationID,
		entity.FieldTempleteIsDeleted:      false,
	}

	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		filter[entity.FieldTempleteName] = bson.M{
			"$regex":   *filters.Search,
			"$options": "i",
		}
	}

	if filters.StartDate != nil || filters.EndDate != nil {
		dateFilter := bson.M{}
		if filters.StartDate != nil {
			dateFilter["$gte"] = *filters.StartDate
		}
		if filters.EndDate != nil {
			dateFilter["$lte"] = *filters.EndDate
		}
		filter[entity.FieldTempleteCreatedAt] = dateFilter
	}

	findOptions := options.Find()
	if filters.Skip != nil {
		findOptions.SetSkip(int64(*filters.Skip))
	}
	if filters.Limit != nil {
		findOptions.SetLimit(int64(*filters.Limit))
	}
	findOptions.SetSort(bson.D{{Key: entity.FieldTempleteCreatedAt, Value: -1}})

	cursor, err := r.collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer func() { _ = cursor.Close(context.Background()) }()
	if err = cursor.All(context.Background(), &templetes); err != nil {
		return nil, err
	}
	return
}

func (r *templeteMongoRepository) GetTempleteCount(filters *models.GetTempleteCountFilter) (count int64, err error) {
	filter := bson.M{
		entity.FieldTempleteOrganizationID: filters.OrganizationID,
		entity.FieldTempleteIsDeleted:      false,
	}

	if filters.Search != nil && strings.TrimSpace(*filters.Search) != "" {
		filter[entity.FieldTempleteName] = bson.M{
			"$regex":   *filters.Search,
			"$options": "i",
		}
	}

	if filters.StartDate != nil || filters.EndDate != nil {
		dateFilter := bson.M{}
		if filters.StartDate != nil {
			dateFilter["$gte"] = *filters.StartDate
		}
		if filters.EndDate != nil {
			dateFilter["$lte"] = *filters.EndDate
		}
		filter[entity.FieldTempleteCreatedAt] = dateFilter
	}

	count, err = r.collection.CountDocuments(context.Background(), filter)
	return
}
