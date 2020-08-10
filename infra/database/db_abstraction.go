package database

import (
	"context"
	"fmt"

	"b2w/swapi-challenge/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	Find(context.Context, interface{}) (CursorHelper, error)
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (*mongo.InsertOneResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
}

type CursorHelper interface {
	Close(context.Context) error
	All(context.Context, interface{}) error
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect(context.Context) error
	Disconnect(context.Context) error
	Ping(context.Context, *readpref.ReadPref) error
}

type mongoClient struct {
	cl *mongo.Client
}
type mongoDatabase struct {
	db *mongo.Database
}
type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	crs *mongo.Cursor
}

func NewClient(cfg config.Database) (ClientHelper, error) {
	userInfo := ""
	if cfg.User != "" {
		userInfo += cfg.User
		if cfg.Password != "" {
			userInfo += ":" + cfg.Password
		}
		userInfo += "@"
	}

	dbUri := fmt.Sprintf("mongodb://%s%s", userInfo, cfg.Host)

	c, err := mongo.NewClient(options.Client().ApplyURI(dbUri))
	return &mongoClient{cl: c}, err

}

func NewDatabase(dbName string, client ClientHelper) DatabaseHelper {
	return client.Database(dbName)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (mc *mongoClient) Disconnect(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (mc *mongoClient) Ping(ctx context.Context, pref *readpref.ReadPref) error {
	return mc.cl.Ping(ctx, pref)
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) (CursorHelper, error) {
	cursor, err := mc.coll.Find(ctx, filter)
	return &mongoCursor{crs: cursor}, err
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return mc.coll.InsertOne(ctx, document)
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return mc.coll.DeleteOne(ctx, filter)
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (mc *mongoCursor) Close(ctx context.Context) error {
	return mc.crs.Close(ctx)
}

func (mc *mongoCursor) All(ctx context.Context, results interface{}) error {
	return mc.crs.All(ctx, results)
}
