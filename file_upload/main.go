package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func setupRoutes() {

	http.HandleFunc("/upload", uploadFileHandler)
	http.ListenAndServe(":8000", nil)

}

func main() {

	fmt.Println(">file upload server started")
	setupRoutes()

}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	// 1. Parse the input from request, type multipart/form-data

	r.ParseMultipartForm(10 << 20) //10MB File

	// 2. Retrieve file from posted form data

	file, handler, err := r.FormFile("myFile")

	if err != nil {

		fmt.Println("error writing form-data")
		fmt.Println(err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer file.Close()

	fmt.Printf("File Name: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("File Headers: %+v\n", handler.Header)

	// 3. Write temporary file on server

	tempfile, err := ioutil.TempFile("temp-images", "upload-*.png")

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer tempfile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tempfile.Write(fileBytes)

	// 4. Return success/error response

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully uploaded the file!\n")

}
