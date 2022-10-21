package middleware

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go-server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

// collection object/instance
var taskRegistry TaskRegistry

// create connection with mongo db
func InitTaskRegistry() {
	taskRegistry = newMongoTaskRegistry()
}

func SetTaskRegistryForTests(tr TaskRegistry) {
	taskRegistry = tr
}

func newMongoTaskRegistry() *MongoTaskRegistry {
	loadTheEnv()
	return createDBInstance()
}

func loadTheEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func createDBInstance() *MongoTaskRegistry {
	// DB connection string
	connectionString := os.Getenv("DB_URI")

	// Database Name
	dbName := os.Getenv("DB_NAME")

	// Collection name
	collName := os.Getenv("DB_COLLECTION_NAME")

	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	var taskRegistry = MongoTaskRegistry{Collection: client.Database(dbName).Collection(collName)}

	fmt.Println("Collection instance created!")

	return &taskRegistry
}

type TaskRegistry interface {
	GetAllTask() []models.ToDoList
	InsertOneTask(task models.ToDoList)
	TaskComplete(task string)
	UndoTask(task string)
	DeleteOneTask(task string)
}

type MongoTaskRegistry struct {
	*mongo.Collection
}

// get all task from the DB and return it
func (tr *MongoTaskRegistry) GetAllTask() []models.ToDoList {
	cur, err := tr.Collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var results []models.ToDoList
	for cur.Next(context.Background()) {
		var result models.ToDoList
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		// fmt.Println("cur..>", cur, "result", reflect.TypeOf(result), reflect.TypeOf(result["_id"]))
		results = append(results, result)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.Background())
	return results
}

// Insert one task in the DB
func (tr *MongoTaskRegistry) InsertOneTask(task models.ToDoList) {
	insertResult, err := tr.Collection.InsertOne(context.Background(), task)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}

// task complete method, update task's status to true
func (tr *MongoTaskRegistry) TaskComplete(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": true}}
	result, err := tr.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// task undo method, update task's status to false
func (tr *MongoTaskRegistry) UndoTask(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"status": false}}
	result, err := tr.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

// delete one task from the DB, delete by ID
func (tr *MongoTaskRegistry) DeleteOneTask(task string) {
	fmt.Println(task)
	id, _ := primitive.ObjectIDFromHex(task)
	filter := bson.M{"_id": id}
	d, err := tr.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Document %s\n", d.DeletedCount)
}
