package controller

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Path    string
	Method  string
	Handler func()
}

type Controller struct {
	Prefix string
	Data   interface{}
	Res    http.ResponseWriter
	Req    *http.Request
	Routes []Route

	Params  map[string]string
	Options map[string][]string
}

func (c *Controller) Render(name string) {
	temp, err := template.ParseFiles(name)
	if err != nil {
		c.Raw("Internal Server Error" + err.Error())
		return
	}
	temp.Execute(c.Res, c.Data)
}

func (c *Controller) Raw(data string) {
	c.Write([]byte(data))
}

func (c *Controller) Write(data []byte) (int, error) {
	return c.Res.Write(data)
}

func (c *Controller) toHandleFunc(pred func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Params = mux.Vars(r)
		c.Options = r.URL.Query()
		pred()
	}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Req = r
	c.Res = w

	m := mux.NewRouter()
	for _, route := range c.Routes {
		m.HandleFunc(route.Path, c.toHandleFunc(route.Handler)).
			Methods(route.Method)
	}
	m.ServeHTTP(w, r)
}
