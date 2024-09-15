package document_req

import (
	document_handler "annotater/internal/http-server/handlers/document"
	response "annotater/internal/lib/api"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	documentPathUrl = "http://localhost:8080/document/"
	formFileName    = "file"
	EXT             = ".pdf"
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
	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var respJson response.Response
		err = json.NewDecoder(resp.Body).Decode(&respJson)
		if err != nil {
			return err
		}
		if respJson.Status == response.StatusError {
			return errors.New(respJson.Error)
		}
	}

	out, err := os.Create(folderPath + "/" + "err_report_" + filepath.Base(documentPath))
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

func GetDocument(client *http.Client, documentPath string, jwtToken string, id uuid.UUID) error {
	url := documentPathUrl + "getDocument"

	reqBody := document_handler.RequestID{ID: id}

	marhsalledBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	var req *http.Request
	req, err = http.NewRequest("GET", url, bytes.NewReader(marhsalledBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)

	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var respJson response.Response
		err = json.NewDecoder(resp.Body).Decode(&respJson)
		if err != nil {
			return err
		}
		if respJson.Status == response.StatusError {
			return errors.New(respJson.Error)
		}
	}

	if err != nil {
		return err
	}

	out, err := os.Create(documentPath + "/" + "document_" + id.String() + EXT)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return err
}

func GetReport(client *http.Client, documentPath string, jwtToken string, id uuid.UUID) error {
	url := documentPathUrl + "getReport"

	reqBody := document_handler.RequestID{ID: id}

	marhsalledBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}
	var req *http.Request
	req, err = http.NewRequest("GET", url, bytes.NewReader(marhsalledBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	resp, err := client.Do(req)

	ct := resp.Header.Get("Content-Type")
	if strings.Contains(ct, "application/json") {
		var respJson response.Response
		err = json.NewDecoder(resp.Body).Decode(&respJson)
		if err != nil {
			return err
		}
		if respJson.Status == response.StatusError {
			return errors.New(respJson.Error)
		}
	}

	if err != nil {
		return err
	}

	out, err := os.Create(documentPath + "/" + "report_" + id.String() + EXT)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return err
}

func GetDocumentsMetaData(client *http.Client, jwtToken string) (*document_handler.ResponseGettingMetaData, error) {
	url := documentPathUrl + "getDocumentsMeta"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	respJson, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	var resp document_handler.ResponseGettingMetaData
	err = render.DecodeJSON(respJson.Body, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
