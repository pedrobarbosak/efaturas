package db

import (
	"context"
	"time"

	"efaturas-xtreme/pkg/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func New(databaseURI string, databaseName string) (*DB, error) {
	client := options.Client()
	client.ApplyURI(databaseURI)
	client.SetMaxConnIdleTime(15 * time.Second)
	client.SetSocketTimeout(30 * time.Second)
	client.SetServerSelectionTimeout(15 * time.Second)

	cl, err := mongo.Connect(context.Background(), client)
	if err != nil {
		return nil, errors.New("failed to connect:", err)
	}

	if err = cl.Ping(context.Background(), nil); err != nil {
		return nil, errors.New("failed to ping db:", err)
	}

	return &DB{Client: cl, Database: cl.Database(databaseName)}, nil
}

func (db *DB) Disconnect() {
	_ = db.Client.Disconnect(context.TODO())
}

func (db *DB) Create(object interface{}, collectionName string, context context.Context) (string, error) {
	collection := db.Database.Collection(collectionName)
	response, err := collection.InsertOne(context, object)
	if err != nil {
		return "", errors.New("failed to create an object in collection:", collectionName)
	}

	id, ok := response.InsertedID.(primitive.ObjectID)
	if !ok {
		if stringID, valid := response.InsertedID.(string); valid {
			return stringID, nil
		}

		return "", errors.New("failed to retrieve id from created object in collection:", collectionName)
	}

	return id.Hex(), nil
}

func (db *DB) ReplaceByID(objectID primitive.ObjectID, object interface{}, collectionName string, context context.Context) error {
	collection := db.Database.Collection(collectionName)

	var item bson.M
	b, _ := bson.Marshal(object)
	err := bson.Unmarshal(b, &item)
	if err != nil {
		return err
	}
	item["_id"] = objectID

	result, err := collection.ReplaceOne(context, bson.M{"_id": objectID}, item)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.NewNotFound("cannot find object id from collection:", collectionName, objectID)
	}
	return nil
}

func (db *DB) Fetch(ctx context.Context, col string, filter interface{}, findOptions *options.FindOptions) ([]map[string]interface{}, error) {
	collection := db.Database.Collection(col)

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var objects []map[string]interface{}
	for cursor.Next(ctx) {
		var object map[string]interface{}
		err := cursor.Decode(&object)
		if err != nil {
			return nil, err
		}
		objects = append(objects, object)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}

	return objects, nil
}

func (db *DB) Aggregate(ctx context.Context, col string, aggregation string, countFilter interface{}) ([]interface{}, int64, error) {
	collection := db.Database.Collection(col)

	var pipeline interface{}
	err := bson.UnmarshalExtJSON([]byte(aggregation), true, &pipeline)
	if err != nil {
		return nil, 0, err
	}

	aggrOptions := options.Aggregate().SetAllowDiskUse(true).SetBatchSize(10000000).SetCollation(&options.Collation{
		Locale: "en", Strength: 3,
	})

	cursor, err := collection.Aggregate(ctx, pipeline, aggrOptions)
	if err != nil {
		return nil, 0, err
	}

	var objects []interface{}
	for cursor.Next(ctx) {
		var object interface{}
		err := cursor.Decode(&object)
		if err != nil {
			return nil, 0, err
		}
		objects = append(objects, object)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	err = cursor.Close(ctx)
	if err != nil {
		return nil, 0, err
	}

	totalRecords, err := collection.CountDocuments(ctx, countFilter)
	if err != nil {
		return nil, 0, err
	}
	return objects, totalRecords, nil
}

func (db *DB) GetByID(id primitive.ObjectID, collectionName string, context context.Context) (*interface{}, error) {
	collection := db.Database.Collection(collectionName)
	filter := bson.M{"_id": id}
	cursor, err := collection.Find(context, filter)
	if err != nil {
		return nil, err
	}
	if !cursor.Next(context) {
		err = cursor.Close(context)
		if err != nil {
			return nil, err
		}

		return nil, errors.NewNotFound("failed to getByID:", collectionName, id.Hex())
	}
	var object interface{}
	err = cursor.Decode(&object)
	if err != nil {
		cursor.Close(context)
		return nil, err
	}

	err = cursor.Close(context)
	if err != nil {
		return nil, err
	}

	return &object, nil
}

func (db *DB) WithTransaction(ctx context.Context, fn func(sc context.Context) error) error {
	session, err := db.Client.StartSession()
	if err != nil {
		return err
	}

	defer session.EndSession(ctx)

	if err = session.StartTransaction(); err != nil {
		return err
	}

	return mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = fn(sc); err != nil {
			_ = session.AbortTransaction(sc)
			return err
		}

		return session.CommitTransaction(sc)
	})
}
