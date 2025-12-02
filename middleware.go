package expressgo

import (
	"fmt"
	"os"
)

type Middleware func(Handler) Handler

func StaticDirectory(pathname string) Middleware {
	return func(next Handler) Handler {
		return func(req *Request, res *Response) {
			filepath := pathname + req.Path
			_, err := os.Stat(filepath)
			if err != nil {
				next(req, res)
				return
			}
			res.WriteFile(filepath)
		}
	}
}

func FileLogging(pathname string) Middleware {
	return func(next Handler) Handler {
		return func(req *Request, res *Response) {
			_, err := os.Stat(pathname)
			if err != nil {
				fmt.Printf("Error: File not found at %s, cannot initiate file logging \n", pathname)
				next(req, res)
				return
			}
			req.LogToFile(pathname)
			next(req, res)
		}
	}
}
