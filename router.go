package expressgo

type Router struct {
	Get  map[string]Handler
	Post map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		Get:  make(map[string]Handler),
		Post: make(map[string]Handler),
	}
}

func (r *Router) Match(method string, path string) (Handler, error) {
	switch method {
	case "GET":
		if handler, ok := r.Get[path]; ok {
			return handler, nil
		}
	case "POST":
		if handler, ok := r.Post[path]; ok {
			return handler, nil
		}
	}

	return defaultHandler, nil
}

func defaultHandler(req *Request, res *Response) {
	res.Write404()
}
