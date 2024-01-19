package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model
type TempModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// ID primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title     string `json:"title,omitempty"`
	Completed string `json:"completed,omitempty"`
}

// Stating collection type
var collection *mongo.Collection

// Connecting to MongoClient
func init() {
	godotenv.Load(".env")
	URI := os.Getenv("DATABASE_URI")
	DATABASE := os.Getenv("DATABASE")
	COLLECTION := os.Getenv("COLLECTION")
	clientOptions := options.Client().ApplyURI(URI)

	//TODO() for when your not sure
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error Occurred When Connecting MongoDB Client =>>>> ", err)
		log.Fatal(err)
	}
	fmt.Println("<<<<Connected!>>>>")
	collection = client.Database(DATABASE).Collection(COLLECTION)
	fmt.Println("<<<<Connection instance is ready>>>>")

	//>>>>>>>>>>>>>>>>>>>>>>>>>
	// res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	// if err != nil { return  }
	// id := res.InsertedID
	// fmt.Printf("idk=>>> %v", id)

}

// TODO: Move helpers to separate file
// Lowercase: it's a helper method and we are not exporting it
func createOneItem(item TempModel) {
	res, err := collection.InsertOne(context.Background(), item)
	if err != nil {
		fmt.Println("Error Occurred When Adding Item =>>>>", err)
		log.Fatal(err)
	}
	fmt.Println("\nItem Successfully Added With ID =>>>> ", res.InsertedID)
}

func updateOneItem(itemId string) {
	id, err := primitive.ObjectIDFromHex(itemId)
	if err != nil {
		fmt.Println("Error Occurred When Updating Item =>>>> ", err)
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": true}}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error Occurred Trying To Add Item To Collection")
		log.Fatal(err)
	}
	fmt.Println("\nUpdate Completed! =>>>>", res)
	fmt.Println("\nres.ModifiedCount =>>>>", res.ModifiedCount)
}

// func deleteOneItem(itemId string ){
// 	id, err := primitive.ObjectIDFromHex(itemId)
// 	if err != nil{
// 		fmt.Println("Error Occurred in DELETE Item Id section =>>>> ", err)
// 		log.Fatal(err)
// 	}
// 	filter := bson.M{"_id": id}

// 	res, err := collection.DeleteOne(context.TODO(), filter)
// 	if err != nil{
// 		fmt.Println("Error occurred after DELETE COLLECTION")
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Delete Completed", res)
// 	fmt.Println("DElete modifed cOunt", res.DeletedCount)
// }
// func deleteAllItems(){
// 	res, err := collection.DeleteMany(context.Background(), bson.D{{}})
// 	if err != nil {
// 		fmt.Println("Error occurred after DELETE ALL COLLECTION =>>>> ", err)
// 		log.Fatal(err)
// 	}
// 	fmt.Println("DELETE ALL Completed", res)
// 	fmt.Println("Modified DELETE COUNT", res.DeletedCount)
// }

func getAllItems() ([]bson.M, int) {
	//Cursor are pointers to the documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("Error Occurred When Gathering All Items =>>>> ", err)
		log.Fatal(err)
	}
	defer func() {
		cursor.Close(context.Background())
		fmt.Println("<<<<Cursor is closed>>>>")
	}()

	var items []bson.M
	var amount int
	//Next allows us to loop through the cursor, if there is a next value
	//"While true"
	for cursor.Next(context.Background()) {
		var item bson.M
		err := cursor.Decode(&item)
		if err != nil {
			fmt.Println("Error Occurred When Gathering All Items =>>>> ", err)
			log.Fatal(err)
		}
		fmt.Printf("Results =>>>> %v\n", item)
		items = append(items, item)
		amount = amount + 1
	}
	if err != nil {
		fmt.Println("Error occurred after loop FIND ALL DATABASE =>>>> ", err)
		log.Fatal(err)
	}
	fmt.Printf("ITEMS, %v\n", items)
	return items, amount
}

// Exported functions
func GetAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allItems, amount := getAllItems()
	fmt.Println("Everything:", allItems)
	fmt.Println("Amount:", amount)
	json.NewEncoder(w).Encode(allItems)
}

func CreateOneItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var item TempModel
	// idk := make([]byte, 10000)
	// boo, _ := r.Body.Read(idk)
	// fmt.Printf("BODY: %v", string(idk[:boo]))

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error while parsing Create Item form")
		log.Fatal(err)
	}
	title := r.FormValue("title")
	completed := r.FormValue("completed")
	item = TempModel{
		Title:     title,
		Completed: completed,
	}
	json.NewDecoder(r.Body).Decode(&item)

	/*
		Not needed for when parsing form, but when dealing with json, for ex. APIs
		json.NewDecoder(r.Body).Decode(&item)
		json.NewEncoder(w).Encode(item)
	*/
	createOneItem(item)
}

func UpdateOneItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error Occurred when Parsing Update Form")
		log.Fatal(err)
	}
	id := r.FormValue("id")
	updateOneItem(id)
	json.NewEncoder(w).Encode(id)
}
