package planet

import (
	"b2w/swapi-challenge/config"
	"b2w/swapi-challenge/domain"
	"b2w/swapi-challenge/infra/database"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoRepo struct {
	db             database.DatabaseHelper
	commandTimeout time.Duration
}

func (r mongoRepo) CollectionName() string { return "planets" }

func NewMongoRepository(db database.DatabaseHelper) *mongoRepo {
	timeout := config.Data.Database.CommandTimeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	return &mongoRepo{
		db:             db,
		commandTimeout: timeout,
	}
}

func (r *mongoRepo) Insert(p *Planet) error {
	collection := r.db.Collection(r.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), r.commandTimeout)
	defer cancel()

	res, err := collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	p.ID = res.InsertedID.(primitive.ObjectID)

	return nil
}

func (r *mongoRepo) FindAll() ([]Planet, error) {
	collection := r.db.Collection(r.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), r.commandTimeout)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []Planet
	if err = cursor.All(ctx, &result); err != nil {
		return result, err
	}

	return result, nil
}

func (r *mongoRepo) GetById(id primitive.ObjectID) (Planet, error) {
	return r.findOne(bson.M{"_id": id})
}

func (r *mongoRepo) GetByName(name string) (Planet, error) {
	return r.findOne(bson.M{"name": name})
}

func (r *mongoRepo) findOne(filter bson.M) (Planet, error) {
	collection := r.db.Collection(r.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), r.commandTimeout)
	defer cancel()

	var result Planet
	if err := collection.FindOne(ctx, filter).Decode(&result); err != nil {
		if err == mongo.ErrNoDocuments {
			return result, domain.ErrNotFound
		}
		return result, err
	}

	return result, nil
}

func (r *mongoRepo) Delete(id primitive.ObjectID) error {
	collection := r.db.Collection(r.CollectionName())

	ctx, cancel := context.WithTimeout(context.Background(), r.commandTimeout)
	defer cancel()

	if _, err := collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}

	return nil
}
