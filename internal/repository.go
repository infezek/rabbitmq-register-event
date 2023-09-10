package internal

import (
	"context"
	"encoding/json"
	"log"
	"rabbitmain/pkg/entity"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type RepoMongo struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type RepoMongoInterface interface {
	Register(entity interface{}) error
	FindEvent(pipe string) ([]RegisterEvent, error)
}

func NewMongo() (*RepoMongo, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27021"))
	if err != nil {
		return nil, err
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, err
	}
	databaseName := strings.ReplaceAll(entity.ExchangeRequestCollects, ".", "-")
	collectionName := strings.ReplaceAll(entity.QueueRequestCollects, ".", "-")
	collection := client.Database(databaseName).Collection(collectionName)
	return &RepoMongo{
		Client:     client,
		Collection: collection,
	}, nil
}

type RegisterEvent struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Body         interface{}        `json:"body" bson:"body"`
	RegisteredAt time.Time          `json:"registeredAt" bson:"registered_at"`
}

func (m *RepoMongo) Register(value interface{}) error {
	event := RegisterEvent{
		ID:           primitive.NewObjectID(),
		Body:         value,
		RegisteredAt: time.Now(),
	}
	_, err := m.Collection.InsertOne(context.Background(), event)
	return err
}

func (m *RepoMongo) FindEvent(pipeline string) ([]RegisterEvent, error) {
	var filterMap map[string]interface{}
	if err := json.Unmarshal([]byte(pipeline), &filterMap); err != nil {
		log.Fatal("Erro ao deserializar JSON:", err)
	}
	filterBson, err := bson.Marshal(filterMap)
	if err != nil {
		log.Fatal("Erro ao converter para BSON:", err)
	}
	var filter bson.M
	if err := bson.Unmarshal(filterBson, &filter); err != nil {
		log.Fatal("Erro ao desconverter BSON:", err)
	}
	cur, err := m.Collection.Find(context.Background(), filter)
	if err != nil {
		log.Fatal("Erro ao realizar a pesquisa:", err)
	}
	var event []RegisterEvent
	if err := cur.All(context.Background(), &event); err != nil {
		return []RegisterEvent{}, err
	}
	return event, nil
}
