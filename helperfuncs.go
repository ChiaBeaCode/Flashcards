package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, code int, payload interface{}){
	data, err := json.Marshal(payload)
	if err != nil{
		//Logs on the server side
		log.Printf("Unable to marshal JSON response for: %v ", payload)
		w.WriteHeader(500)
		return	
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}


func errorReponse(w http.ResponseWriter, code int, message string){
	if code > 499 {
		log.Println("Responding with 5XX error: ", message)
	}
	type errResponse struct {
		Error string `json:"error"`
			/*
			field type json reflect tag
			"I have an Error field of type string and I want the field of the key to be error"
			Ex. {
				"error": "something went wrong"
			}
			Specifies how you want the json.Marshal how to convert this struct in a JSON object
			*/
			
	}
	jsonResponse(w, code, errResponse{
		Error: message,
	})

}

//Should only respond if the server is alive and ready
//function signature needed to define an http handler that Go expects
func handlerServerReadiness(w http.ResponseWriter, r *http.Request){
	//Sends an empty JSON object
	jsonResponse(w, 200, struct{}{})
}
func handlerError(w http.ResponseWriter, r *http.Request){
	errorReponse(w, 400, "something went wrong")
}