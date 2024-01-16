package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ChiaBeaCode/GoWebServer/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)
type ClientStruct struct{
	client *mongo.Client
}
func NewClientStruct(c *mongo.Client) *ClientStruct {
	return &ClientStruct{
		client: c,
	}
}
// func ( c *ClientStruct) beginConnection() error{
// 	collection := c.client.Database("baz").Collection("qux")
// 	fmt.Printf("collection=>>>> %v\n", *collection)
// 	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Printf("RESP =>>>>> %v\n", resp)
// 	c.client.Disconnect(context.TODO())
// 	return nil
// }
func main(){
	godotenv.Load(".env")
	port, found := os.LookupEnv("PORT")
	if !found || port == ""{
		fmt.Println("PORT was not found or is empty")
	}


	URI, found := os.LookupEnv("DATABASE_URI")
	if !found || URI == ""{
		fmt.Println("URI was not found or is empty")
	}


	router := chi.NewRouter()
		//Users can make a request to server from browser
		//Very open/poor security so modify later
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*"},
		AllowedMethods: []string{"GET, POST, PUT, DELETE, OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,

	}))

	
		//Nesting routers in case of breaking changes, you havea handler in each path to take care of the action
		//Full path == ..../v1/ready
	v1Router := chi.NewRouter()
		//HandleFunc matches any method, hence .Get instead
	v1Router.Get("/ready", handlerServerReadiness)
	v1Router.Get("/error", handlerError)

	//<<<<New-Testing>>>>
	v1Router.Post("/api/item", controller.CreateItem)
	// v1Router.Get("/items", controller.GetAllItems)
	// v1Router.Post("/item", controller.CreateItem)
	// v1Router.Put("/item/{id}", controller.UpdateCompleted)
	//<<<<>>>>>
	router.Mount("/v1", v1Router)

	server := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

		/* Blocker
			if err := http.ListenAndServe("3000", nil); err !=nil{
		 		log.Fatal(err) }
		*/
	err := server.ListenAndServe()
	if err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Server starting on port: %v", port)

	
	
		//Client returns a reference to a struct(binary slice, addresses, booleans, map[] and nil) 
	// fmt.Println("CLIENT=>>>>", client)
	// cc := NewClientStruct(client)
	// cc.beginConnection()
	// defer cc.client.Disconnect(ctx)

} 