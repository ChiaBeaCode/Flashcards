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

type ClientStruct struct {
	client *mongo.Client
}

// Client returns a reference to a struct(binary slice, addresses, booleans, map[] and nil)
func NewClientStruct(c *mongo.Client) *ClientStruct {
	return &ClientStruct{
		client: c,
	}
}
func main() {
	godotenv.Load(".env")
	PORT, found := os.LookupEnv("PORT")
	if !found || PORT == "" {
		fmt.Println("PORT Not Found/Or Is Empty")
		return
	}

	URI, found := os.LookupEnv("DATABASE_URI")
	if !found || URI == "" {
		fmt.Println("URI Not Found/Or Is Empty")
		return
	}
	router := chi.NewRouter()
	//Users can make a request to server from browser
	//Very open/poor security so modify later
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET, POST, PUT, DELETE, OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//Nesting routers in case of breaking changes, you havea handler in each path to take care of the action
	//Full path == ..../v1/ready
	v1Router := chi.NewRouter()
	//HandleFunc matches any method, hence .Get instead
	v1Router.Get("/ready", handlerServerReadiness)
	v1Router.Get("/error", handlerError)

	//<<<<CRUD Path Testing>>>>
	v1Router.Post("/item", controller.CreateOneItem)
	v1Router.Get("/items", controller.GetAllItems)
	//<<<<>>>>>
	router.Mount("/v1", v1Router)
	router.Handle("/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Handler: router,
		Addr:    ":" + PORT,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("Server starting on port: %v", PORT)
	}
}
