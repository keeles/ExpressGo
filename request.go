package expressgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

type Request struct {
	method  string
	Path    string
	headers map[string]string
	body    []byte
}

func NewReq(method string, path string, headers map[string]string, body []byte) *Request {
	return &Request{
		method:  method,
		Path:    path,
		headers: headers,
		body:    body,
	}
}

func (req *Request) ParseJson(v interface{}) error {
	if req.headers["Content-Type"] != "application/json" {
		return errors.New("no application/json content type found for this request, unable to parse json")
	}

	return json.Unmarshal(req.body, v)
}

func (req *Request) ParseFormData() (url.Values, error) {
	if req.headers["Content-Type"] != "application/x-www-form-urlencoded" {
		return nil, errors.New("no application/x-www-form-urlencoded content type found for this request, unable to parse form data")
	}

	return url.ParseQuery(string(req.body))
}

func (req *Request) LogToFile(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	var headers string
	for key, value := range req.headers {
		headers += fmt.Sprintf("%s: %s; ", key, value)
	}
	logString := fmt.Sprintf("%s %s %s \n", req.method, req.Path, headers)
	_, err = file.WriteString(logString)
	if err != nil {
		return err
	}
	return nil
}
