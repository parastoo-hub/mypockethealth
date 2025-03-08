package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var serverStorageDir = "./testdata/server/storage/"

func UploadDicomFile(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to read file.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if handler.Size == 0 {
		http.Error(w, "Uploaded file is empty.", http.StatusBadRequest)
		return
	}

	if err := createServerStorageDirectoryIfNeeded(serverStorageDir); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Println("Could not create storage directory.")
		return
	}

	filePath := serverStorageDir + handler.Filename
	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Unable to save file.", http.StatusInternalServerError)
		logger.Println("Could not create storage directory: " + err.Error())
		return
	}
	defer outFile.Close()
	io.Copy(outFile, file)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "File uploaded successfully: %s.\n", handler.Filename)
}

func createServerStorageDirectoryIfNeeded(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Directory does not exist, so create it
		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("failed to create server directory: %v", err)
		}
	}
	return nil
}
