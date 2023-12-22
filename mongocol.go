// mongodb.go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client     *mongo.Client
	database   *mongo.Database
	collections map[string]*mongo.Collection
}

const NodeMeta string = "nodemeta"
const PeerMeta string = "peermeta_"
const CollectionPrefix string = "selfcollection_"
const ImportedPrefix string = "importedcollection_"
const PeerCollection string = "peercollection_"

func NewMongoDB(connectionString, dbName string, collectionNames ...string) (*MongoDB, error) {
	
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	db := client.Database(dbName)
	
	collections := make(map[string]*mongo.Collection)
	for _, name := range collectionNames {
		collections[name] = db.Collection(name)
	}

	return &MongoDB{
		client:     client,
		database:   db,
		collections: collections,
	}, nil
}

func (m *MongoDB) DropAllCollections() error {
    ctx := context.Background()
    collections, err := m.database.ListCollectionNames(ctx, bson.M{})
    if err != nil {
        return fmt.Errorf("failed to list collections: %v", err)
    }

    for _, col := range collections {
        if err := m.database.Collection(col).Drop(ctx); err != nil {
            return fmt.Errorf("failed to drop collection '%s': %v", col, err)
        }
    }

    return nil
}

func (m *MongoDB) AddCollection(collectionName string) error {
	// Check if the collection already exists
	fmt.Println("the collection is ", collectionName)
	if _, ok := m.collections[collectionName]; ok {
		return fmt.Errorf("collection '%s' already exists", collectionName)
	}

	// Create the new collection
	newCollection := m.database.Collection(collectionName)

	// Add the new collection to the map
	m.collections[collectionName] = newCollection

	return nil
}

func (m *MongoDB) DropCollection(collectionName string) error {
	// Check if the collection exists
	coll, ok := m.collections[collectionName]
	if !ok {
		return fmt.Errorf("collection '%s' not found", collectionName)
	}

	// Drop (delete) the collection
	err := coll.Drop(context.Background())
	if err != nil {
		return err
	}

	// Remove the collection from the map
	delete(m.collections, collectionName)

	fmt.Printf("Collection '%s' dropped successfully\n", collectionName)
	return nil
}

func (m *MongoDB) Close() {
	if m.client != nil {
		err := m.client.Disconnect(context.Background())
		if err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		}
	}
}


// Functions for handling collections

func (m *MongoDB) InsertDocument(collectionName string, document interface{}) error {
	
	coll, ok := m.collections[collectionName]

	fmt.Println("what is happening")
	if !ok {
		return fmt.Errorf("collection '%s' not found", collectionName)
	}
	
	_, err := coll.InsertOne(context.Background(), document)
	fmt.Println("it was addedd successfully")
	return err
}


func (m *MongoDB) EditDocument(collectionName string,filter bson.M, updatedDocument interface{}) error {
	coll, ok := m.collections[collectionName]

	if !ok {
		return fmt.Errorf("collection '%s' not found", collectionName)
	}

	update := bson.M{"$set": updatedDocument}

	_, err := coll.UpdateOne(context.Background(), filter, update)
	return err
}

func (m *MongoDB) InsertManyDocuments(collectionName string, documents []interface{}) error {
    coll, ok := m.collections[collectionName]
    if !ok {
        return fmt.Errorf("collection '%s' not found", collectionName)
    }

    _, err := coll.InsertMany(context.Background(), documents)
    if err != nil {
        return err
    }

    fmt.Println("Documents were added successfully")
    return nil
}

func (m *MongoDB) DeleteDocument(collectionName string, documentID string) error {
	coll, ok := m.collections[collectionName]
	if !ok {
		return fmt.Errorf("collection '%s' not found", collectionName)
	}

	filter := bson.D{{"url", documentID}}
	result, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("document with _id '%s' not found", documentID)
	}

	return nil
}

func (m *MongoDB) GetDocuments(collectionName string, filter bson.D, result interface{}) error {
	coll, ok := m.collections[collectionName]
	if !ok {
		return fmt.Errorf("collection '%s' not found", collectionName)
	}

	cursor, err := coll.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	return cursor.All(context.Background(), result)
}

func (m *MongoDB) GetDocumentsm(collectionName string, filter bson.M, result interface{}) error {
    coll, ok := m.collections[collectionName]
    if !ok {
        return fmt.Errorf("collection '%s' not found", collectionName)
    }

    cursor, err := coll.Find(context.Background(), filter)
    if err != nil {
        return err
    }
	defer cursor.Close(context.Background())

	return cursor.All(context.Background(), result)
}

func (m *MongoDB) GetDocumentAndEdit(collectionName string, filter bson.M) error {
    coll, ok := m.collections[collectionName]
    if !ok {
        return fmt.Errorf("collection '%s' not found", collectionName)
    }

    update := bson.M{"$set": bson.M{"working": false}}

    _, err := coll.UpdateOne(context.Background(), filter, update)
    if err != nil {
        return err
    }

    return nil
}

// Extend this with more CRUD functions as needed

// Structs representing document structures

type BaseInfo struct {
	// Define fields as needed
}

type NodeInfo struct {
	ID    string `bson:"id"`
	URL string `bson:"url"`
	// Add more fields as needed
}



type MetaInfo struct {
	DataID      int `bson:"data
	id"`
	Title       string `bson:"title"`
	Description string `bson:"description"`
	Seconds string `bson:interval`
	Working bool `bson:working`
	// Add more fields as needed
}

type Data struct {
	DataID int `bson:"data_id"`
	JSONres string `bson:"res"`
	Signature string `bson:"signature"`
	Timestamp int `bson:"timestamp"`
	// Add more fields as needed
}
