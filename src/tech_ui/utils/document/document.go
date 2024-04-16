package document_ui

import (
	document_handler "annotater/internal/http-server/handlers/document"
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	documentPathUrl = "http://localhost:8080/document/"
	formFileName    = "file"
)

func CheckDocument(client *http.Client, documentPath string, jwtToken string) (*document_handler.ResponseCheckDoucment, error) {

	url := documentPathUrl + "check"

	file, err := os.Open(documentPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Prepare the multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", filepath.Base(documentPath))
	if err != nil {
		return nil, err
	}

	// Copy the file content to the part
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create a POST request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Set the content type header
	req.Header.Set("Content-Type", writer.FormDataContentType())
	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Decode the response
	var response document_handler.ResponseCheckDoucment
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
