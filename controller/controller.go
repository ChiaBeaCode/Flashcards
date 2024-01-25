package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ChiaBeaCode/GoWebServer/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Model
//
//	type CardModel struct {
//		ID primitive.ObjectID `bson:"_id,omitempty"`
//		// ID primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
//		Title      string `json:"title,omitempty"`
//		Definition string `json:"definition,omitempty"`
//	}
// var cardModel models.CardModel

// Stating collection type
var collection *mongo.Collection

// Connecting to MongoClient
func init() {
	godotenv.Load(".env")
	DATABASE_URI := os.Getenv("DATABASE_URI")
	DATABASE := os.Getenv("DATABASE")
	COLLECTION := os.Getenv("COLLECTION")
	clientOptions := options.Client().ApplyURI(DATABASE_URI)

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
}

// TODO: Move helpers to separate file
// Lowercase: it's a helper method and we are not exporting it
func createOneCard(card models.CardModel) {
	if err := findCard(card); err != nil {
		res, err := collection.InsertOne(context.Background(), card)
		if err != nil {
			fmt.Println("Error Occurred When Adding Card =>>>>", err)
			log.Fatal(err)
		}
		fmt.Println("\n\nCard Successfully Added With Error =>>>> ", res.InsertedID)
		fmt.Println("\n\nError=>>>> ", err)
	}
}

func findCard(card models.CardModel) error {
	fmt.Println("\n\nWhat Appeared: ", card)
	filter := bson.M{"title": card.Title}
	fmt.Println("\n\nFiltered: ", filter)
	var newRes bson.M
	res := collection.FindOne(context.Background(), filter).Decode(&newRes)
	fmt.Println("\n\n\n Found it! ", res == nil)
	return res

}

func updateOneCard(cardId string) {
	id, err := primitive.ObjectIDFromHex(cardId)
	if err != nil {
		fmt.Println("Error Occurred When Updating Card =>>>> ", err)
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"definition": true}}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("Error Occurred Trying To Add Card To Collection")
		log.Fatal(err)
	}
	fmt.Println("\nUpdate Completed! =>>>>", res)
	fmt.Println("\nres.ModifiedCount =>>>>", res.ModifiedCount)
}

// func deleteOneCard(cardId string ){
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
// 	fmt.Println("Delete Definition", res)
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

func getAllCards() []bson.M {
	//Cursor are pointers to the documents in the collection
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		fmt.Println("Error Occurred When Gathering All Cards =>>>> ", err)
		log.Fatal(err)
	}
	defer func() {
		cursor.Close(context.Background())
		fmt.Println("<<<<Cursor is closed>>>>")
	}()

	var cards []bson.M
	//Next allows us to loop through the cursor, if there is a next value
	//"While true"
	for cursor.Next(context.Background()) {
		var card bson.M
		err := cursor.Decode(&card)
		if err != nil {
			fmt.Println("Error Occurred When Gathering All Cards =>>>> ", err)
			log.Fatal(err)
		}
		fmt.Printf("Results =>>>> %v\n", card)
		cards = append(cards, card)
	}
	if err != nil {
		fmt.Println("Error occurred after loop FIND ALL DATABASE =>>>> ", err)
		log.Fatal(err)
	}
	fmt.Printf("CARDS, %v\n", cards)
	return cards
}

// Exported functions
func GetAllCards(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allCards := getAllCards()
	fmt.Println("Everything:", allCards)
	json.NewEncoder(w).Encode(allCards)
}

func CreateOneCard(w http.ResponseWriter, r *http.Request) {
	fmt.Println("1")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	var card models.CardModel
	json.NewDecoder(r.Body).Decode(&card)
	/*
		Not needed for when parsing form, but when dealing with json, for ex. APIs
		json.NewDecoder(r.Body).Decode(&card)
		json.NewEncoder(w).Encode(card)
	*/
	createOneCard(card)
}

func UpdateOneCard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	if err := r.ParseForm(); err != nil {
		fmt.Println("Error Occurred when Parsing Update Form")
		log.Fatal(err)
	}
	id := r.FormValue("id")
	updateOneCard(id)
	json.NewEncoder(w).Encode(id)
}
