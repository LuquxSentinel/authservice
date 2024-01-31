package main

import (
	"context"
	"time"

	"github.com/luqus/authservice/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Storage interface {
	Get(ctx context.Context, email string) (*types.User, error)
	Create(ctx context.Context, user *types.User) error
	CountEmail(ctx context.Context, email string) (int64, error)
	GenerateID() string
}

func initStorage(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

type MongoStorage struct {
	userCollection *mongo.Collection
}

func NewMongoStorage(userCollection *mongo.Collection) *MongoStorage {
	return &MongoStorage{
		userCollection: userCollection,
	}
}

func (s *MongoStorage) Get(ctx context.Context, email string) (*types.User, error) {
	filter := primitive.D{primitive.E{Key: "email", Value: email}}

	user := new(types.User)

	err := s.userCollection.FindOne(ctx, filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *MongoStorage) Create(ctx context.Context, user *types.User) error {
	_, err := s.userCollection.InsertOne(ctx, user)
	return err
}

func (s *MongoStorage) CountEmail(ctx context.Context, email string) (int64, error) {
	filter := primitive.D{primitive.E{Key: "email", Value: email}}
	count, err := s.userCollection.CountDocuments(ctx, filter)
	return count, err
}

func (s *MongoStorage) GenerateID() string {
	return primitive.NewObjectID().Hex()
}
