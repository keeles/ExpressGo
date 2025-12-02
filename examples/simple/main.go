package main

import (
	"fmt"
	"strconv"

	eg "github.com/keeles/expressgo"
)

func main() {
	srv := eg.NewServer()
	srv.Use(eg.FileLogging("logs/logs.txt"))
	srv.Use(eg.StaticDirectory("public"))
	srv.Get("/", pageHandler)
	srv.Post("/api", apiHandler)
	srv.Post("/form-endpoint", formHandler)
	srv.Listen(":8080")
}

type User struct {
	Name string
	Age  int
}

type Form struct {
	Name  string
	Email string
	Age   int
}

func apiHandler(req *eg.Request, res *eg.Response) {
	var user User
	if err := req.ParseJson(&user); err != nil {
		fmt.Println(err)
		res.Status(500)
		res.WriteString("Error!")
	}
	res.WriteJson(user)
}

func formHandler(req *eg.Request, res *eg.Response) {
	var form Form
	values, err := req.ParseFormData()
	if err != nil {
		fmt.Println(err)
		res.Status(500)
		res.WriteString("Error!")
		return
	}
	form.Name = values.Get("name")
	form.Email = values.Get("email")
	form.Age, _ = strconv.Atoi(values.Get("age"))
	res.WriteJson(form)
}

func pageHandler(req *eg.Request, res *eg.Response) {
	res.WriteFile("index.html")
}
