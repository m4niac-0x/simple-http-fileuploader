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

func renderUploadForm(w http.ResponseWriter) {
	tmpl, err := template.New("uploadForm").Parse(uploadForm)
	if err != nil {
		http.Error(w, "Error rendering form", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100 << 20) // 100 MB limit
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, "Unable to create file on server", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File uploaded successfully: %s\n", handler.Filename)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderUploadForm(w)
	case http.MethodPost:
		handleFileUpload(w, r)
	default:
		http.Error(w, "Only GET and POST methods are allowed", http.StatusMethodNotAllowed)
	}
}

func startServer() {
	http.HandleFunc("/", uploadHandler)

	fmt.Println("Starting server on :11337...")
	if err := http.ListenAndServe(":11337", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}

func main() {
	startServer()
}
