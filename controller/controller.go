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

//Model
type TempModel struct {
	ID primitive.ObjectID `bson:"_id,omitempty"` 
	// ID primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Title string `json:"title,omitempty"`
	Completed bool `json:"completed,omitempty"`
}

	//Stating collection type
var collection *mongo.Collection

	//Connecting to MongoClient
func init() {
	godotenv.Load(".env")
	URI := os.Getenv("DATABASE_URI")
	DATABASE := os.Getenv("DATABASE")
	COLLECTION := os.Getenv("COLLECTION")

	clientOptions := options.Client().ApplyURI(URI)

		//Background()
		//TODO for when your not sure
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Error occured: ", err)
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
//will put helpers in seperate file

	//Lowercase: it's a helper method and we are not exporting it
func insertOneItem(item TempModel){
	res, err := collection.InsertOne(context.Background(), item)
	if err != nil{
		fmt.Println("Error Occurred in INSERT ONE ITEM section")
		log.Fatal(err)
	}
	fmt.Println("Item added with ID =>>>> ", res.InsertedID)
	// id := res.InsertedID
	// fmt.Printf("idk=>>> %v", id)
}
func updateOneItem(itemId string){
	id, err := primitive.ObjectIDFromHex(itemId)
	if err != nil{
		fmt.Println("Error Occurred in UPDATE ONE ITEM Id section =>>>> ", err)
		log.Fatal(err)
	}
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"completed": true}}

	res, err := collection.UpdateOne(context.Background(), filter,update)
	if err != nil {
		fmt.Println("Error occurred after UPDATE COLLECTION")
		log.Fatal(err)
	}
	fmt.Println("Update Completed", res)
	fmt.Println("Modified cOunt", res.ModifiedCount)
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
func getAllItems() []bson.D {
		//Cursor are pointers to the documents in the collection
	cursor, err := collection.Find(context.Background(), bson.D{{}})

	if err != nil {
		fmt.Println("Error occurred while gathering FIND ALL DATABASE =>>>> ", err)
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var items []bson.D
	// var items []primitive.D
		//Next allows us to loop through the cursor, if there is a next value
		//"While true"
	for cursor.Next(context.Background()) {
		var item bson.D
		err := cursor.Decode(&item)
		if err != nil {
			fmt.Println("Error occurred after FIND ALL DATABASE =>>>> ", err)
			log.Fatal(err) 
		}
		fmt.Println("Results =>>>> ", item)
		items = append(items, item)
	}
	if err != nil {
		fmt.Println("Error occurred after loop FIND ALL DATABASE =>>>> ", err)
		log.Fatal(err)
	}
	return items

}


	//Exported functions
func GetAllItems(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	allItems := getAllItems()
	json.NewEncoder(w).Encode(allItems)
}

func CreateItem(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var item TempModel
	json.NewDecoder(r.Body).Decode(&item)
	insertOneItem(item)
	json.NewEncoder(w).Encode(item)
}

func UpdateCompleted(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	if err := r.ParseForm(); err != nil{
		fmt.Println("Error in Export Updating Function")
		log.Fatal(err)
	}
	id := r.FormValue("id")
	updateOneItem(id)
	json.NewEncoder(w).Encode(id)
}









/*
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()


	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Printf("Disconnect Error =>>>> %v", err)
			return
		}
	}()


	collection := client.Database("testing").Collection("numbers")
	res, err := collection.InsertOne(context.Background(), bson.M{"hello": "world"})
	if err != nil { return  }
	id := res.InsertedID
	fmt.Printf("idk=>>> %v", id)
*/