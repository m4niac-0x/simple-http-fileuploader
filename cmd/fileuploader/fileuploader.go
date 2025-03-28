package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

const uploadForm = `
<!DOCTYPE html>
<html>
<head>
	<title>Simple File Uploader</title>
</head>
<body>
	<h1>Upload a File</h1>
	<form enctype="multipart/form-data" action="/" method="post">
		<input type="file" name="file" required>
		<br><br>
		<input type="submit" value="Upload">
	</form>
</body>
</html>
`

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the HTML form
		tmpl := template.New("uploadForm")
		tmpl, _ = tmpl.Parse(uploadForm)
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		// Parse the multipart form
		err := r.ParseMultipartForm(100 << 20) // 100 MB limit
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// Retrieve the file from the form
		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Create a new file on the server
		dst, err := os.Create(handler.Filename)
		if err != nil {
			http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		// Copy the uploaded file to the server
		_, err = io.Copy(dst, file)
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Only GET and POST methods are allowed")
	}
}

func StartServer() {
	http.HandleFunc("/", uploadHandler)

	fmt.Println("Starting server on :11337...")
	err := http.ListenAndServe(":11337", nil)
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func main() {
	// Start the server
	StartServer()
}
