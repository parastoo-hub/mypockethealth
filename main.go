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

	// curl -X POST -F "file=@ST000001/IM000001" http://localhost:8080/dicom/v1/upload
	v1.HandleFunc("/upload", UploadDicomFile).Methods("POST")

	// curl -X GET "http://localhost:8080/dicom/v1/metadata?file=IM000001&tag=(0010,0010)"
	v1.HandleFunc("/metadata", GetDicomMetadata).Methods("GET")

	// curl -X GET "http://localhost:8080/dicom/v1/png/conversion?file=IM000001" -o output
	v1.HandleFunc("/png/conversion", ConvertDicomToPNG).Methods("GET")

	log.Println("Server is listening to you at port 8080 :)")
	http.ListenAndServe(":8080", r)
}
