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

	DATABASE_URI, found := os.LookupEnv("DATABASE_URI")
	if !found || DATABASE_URI == "" {
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
	// <<<<HandleFunc matches any method, hence .Get instead>>>>
	v1Router.Get("/ready", handlerServerReadiness)
	v1Router.Get("/error", handlerError)

	//<<<<CRUD Path Testing>>>>
	// Create card
	// v1Router.Post("/card", controller.CreateOneCard)
	v1Router.Post("/api/flashcards/new", controller.CreateOneCard)

	// Grab all cards
	// v1Router.Get("/cards", controller.GetAllCards)
	v1Router.Get("/api/flashcards", controller.GetAllCards)

	// // Find card
	// // v1Router.Get("/cards", controller.FindCard)
	// v1Router.Get("/api/flashcards/:id", controller.FindOneCard)

	// Update card card/:id
	v1Router.Put("/api/flashcards/:id", controller.UpdateOneCard)

	// Delete card card/:id
	// v1Router.Delete("/api/flashcards/:id", controller.DeleteOneCard)

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
