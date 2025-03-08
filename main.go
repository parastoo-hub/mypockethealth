package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	logger.Println("Application started")

	r := mux.NewRouter()
	v1 := r.PathPrefix("/dicom/v1").Subrouter()

	v1.HandleFunc("/upload", UploadDicomFile).Methods("POST")

	v1.HandleFunc("/metadata", GetDicomMetadata).Methods("GET")

	v1.HandleFunc("/png/conversion", ConvertDicomToPNG).Methods("GET")

	log.Println("Server is listening to you at port 8080 :)")
	http.ListenAndServe(":8080", r)
}
