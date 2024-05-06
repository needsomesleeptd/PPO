package markup_req

import (
	annot_handler "annotater/internal/http-server/handlers/annot"
	response "annotater/internal/lib/api"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

var (
	annotFileFieldName = "annotFile"
	bbsFieldName       = "jsonBbs"
	annotsUrlPath      = "http://localhost:8080/annot/"
)

func AddMarkup(client *http.Client, filePath string, bbs []float32, classLabel uint64, jwtToken string) error {
	url := annotsUrlPath + "add"
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	requestBody := &bytes.Buffer{}
	writer := multipart.NewWriter(requestBody)

	annotFile, err := writer.CreateFormFile(annotFileFieldName, file.Name())
	if err != nil {
		return err
	}
	_, err = io.Copy(annotFile, file)
	if err != nil {
		return err
	}

	jsonReq := annot_handler.RequestAddAnnot{
		ErrorBB:    bbs,
		ClassLabel: classLabel,
	}
	jsonReqMarshalled, err := json.Marshal(jsonReq)
	if err != nil {
		return err
	}
	err = writer.WriteField(bbsFieldName, string(jsonReqMarshalled))
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, requestBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+jwtToken)

	respJson, err := client.Do(req)
	if err != nil {
		return err
	}
	var resp response.Response
	err = json.NewDecoder(respJson.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Status == response.StatusError {
		return errors.New(resp.Error)
	}
	return nil
}

func GetMarkupsByID(client *http.Client, userID uint64, jwtToken string) (*annot_handler.ResponseGetAnnots, error) {
	url := annotsUrlPath + "creatorID"

	req, err := http.NewRequest("GET", url, nil) // there might need to be a paginzing here
	if err != nil {
		return nil, fmt.Errorf("error in request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	respJson, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var resp annot_handler.ResponseGetAnnots
	err = json.NewDecoder(respJson.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("error in request decoding json %v", err)
	}
	if resp.Status == response.StatusError {
		return nil, errors.New(resp.Error)
	}
	return &resp, nil
}

func GetAllMarkups(client *http.Client, jwtToken string) (*annot_handler.ResponseGetAnnots, error) {
	url := annotsUrlPath + "getsAll"

	req, err := http.NewRequest("GET", url, nil) // there might need to be a paginzing here
	if err != nil {
		return nil, fmt.Errorf("error in request %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+jwtToken)

	respJson, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var resp annot_handler.ResponseGetAnnots
	err = json.NewDecoder(respJson.Body).Decode(&resp)
	if err != nil {
		return nil, fmt.Errorf("error in request decoding json %v", err)
	}
	if resp.Status == response.StatusError {
		return nil, errors.New(resp.Error)
	}
	return &resp, nil
}

func DeletingMarkups(client *http.Client, annotID uint64, jwtToken string) error {
	url := annotsUrlPath + "delete"

	jsonReq := annot_handler.RequestID{
		ID: annotID,
	}
	jsonReqMarshalled, err := json.Marshal(jsonReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewReader(jsonReqMarshalled)) // there might need to be a paginzing here
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+jwtToken)
	req.Header.Set("Content-Type", "application/json")

	respJson, err := client.Do(req)
	if err != nil {
		return err
	}
	var resp response.Response
	err = json.NewDecoder(respJson.Body).Decode(&resp)
	if err != nil {
		return err
	}
	if resp.Status == response.StatusError {
		return errors.New(resp.Error)
	}
	return nil
}
