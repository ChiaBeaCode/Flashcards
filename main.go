package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)
func main(){
	godotenv.Load(".env")
	port, found := os.LookupEnv("PORT")
	if !found || port == ""{
		fmt.Printf("PORT was not found or is empty")
	}
	// fmt.Printf("Port: %v", port)
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
	log.Printf("Server starting on port: %v", port)
} 