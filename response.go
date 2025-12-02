package expressgo

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net"
	"os"
	"path"
)

type Response struct {
	conn       net.Conn
	statusCode int
	headers    map[string]string
}

func NewRes(conn net.Conn) *Response {
	return &Response{
		conn:       conn,
		statusCode: 200,
		headers:    make(map[string]string),
	}
}

func (res *Response) WriteString(str string) error {
	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", res.statusCode)

	res.headers["Content-Type"] = "text/plain"

	for key, value := range res.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response += "\n\n" + str

	_, err := res.conn.Write([]byte(response))
	return err
}

func (res *Response) WriteJson(data interface{}) error {
	json, err := json.Marshal(data)
	if err != nil {
		return err
	}

	res.headers["Content-Type"] = "application/json"

	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", res.statusCode)

	for key, value := range res.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response += "\n\n" + string(json)

	_, err = res.conn.Write([]byte(response))
	return err
}

func (res *Response) Status(statusCode int) *Response {
	res.statusCode = statusCode
	return res
}

func (res *Response) WriteFile(pathname string) error {
	file, err := os.Open(pathname)
	if err != nil {
		fmt.Println(err)
		res.Write404()
		return err
	}
	contents, err := io.ReadAll(file)
	if err != nil {
		res.Write500()
		return err
	}

	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", res.statusCode)

	contentType := mime.TypeByExtension(path.Ext(pathname))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	res.headers["Content-Type"] = contentType

	for key, value := range res.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	response += "\n\n"

	_, err = res.conn.Write([]byte(response))
	res.conn.Write(contents)
	return err
}

func (res *Response) Write404() {
	res.Status(404)
	response := fmt.Sprintf("HTTP/1.1 %d Not Found\r\n", res.statusCode)
	errorPage, err := os.Open("public/404.html")
	if err != nil {
		res.headers["Content-Type"] = "text/plain"
	} else {
		res.headers["Content-Type"] = "text/html"
	}

	for key, value := range res.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	contents, err := io.ReadAll(errorPage)
	if err != nil {
		response += "\n\n" + "404 - Page not found"
	} else {
		response += "\n\n" + string(contents)
	}

	res.conn.Write([]byte(response))
}

func (res *Response) Write500() {
	res.Status(500)
	response := fmt.Sprintf("HTTP/1.1 %d Internal Server Error\r\n", res.statusCode)
	errorPage, err := os.Open("public/error.html")
	if err != nil {
		res.headers["Content-Type"] = "text/plain"
	} else {
		res.headers["Content-Type"] = "text/html"
	}

	for key, value := range res.headers {
		response += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	contents, err := io.ReadAll(errorPage)
	if err != nil {
		response += "\n\n" + "500 - Internal Server Error"
	} else {
		response += "\n\n" + string(contents)
	}

	res.conn.Write([]byte(response))
}
