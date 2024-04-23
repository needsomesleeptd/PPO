package document_req

import (
	document_handler "annotater/internal/http-server/handlers/document"
	response "annotater/internal/lib/api"
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
	reportFileName  = "report.pdf"
)

func CheckDocument(client *http.Client, documentPath string, jwtToken string) (*document_handler.ResponseCheckDoucment, error) {

	url := documentPathUrl + "check"

	file, err := os.Open(documentPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(formFileName, filepath.Base(documentPath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var response document_handler.ResponseCheckDoucment
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func ReportDocument(client *http.Client, documentPath string, folderPath string, jwtToken string) error {

	url := documentPathUrl + "report"

	file, err := os.Open(documentPath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(formFileName, filepath.Base(documentPath))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	out, err := os.Create(folderPath + "/" + "Err_report" + filepath.Base(documentPath))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func LoadDocument(client *http.Client, documentPath string, jwtToken string) (*response.Response, error) {
	url := documentPathUrl + "load"

	file, err := os.Open(documentPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(documentPath))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var response response.Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}
