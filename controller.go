package controller

import (
	"html/template"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

type Route struct {
	Path    string
	Method  string
	Handler func()
}

type Controller struct {
	// Configuration
	Prefix   string
	Data     interface{}
	Res      http.ResponseWriter
	Req      *http.Request
	Routes   []Route
	Partials string

	// Oncall
	Params  map[string]string
	Options map[string][]string
}

func hasFieldHelper(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

// Render view to HTTP response writer.
func (c *Controller) Render(name string) error {
	t := template.New(name + ".html").Funcs(template.FuncMap{
		"hasField": hasFieldHelper,
	})
	t = template.Must(t.ParseFiles("views/" + name + ".html"))
	if c.Partials != "" {
		t = template.Must(t.ParseGlob(c.Partials + "/*.html"))
	}
	return t.Execute(c, c.Data)
}

func (c *Controller) RenderError(view string, err error) error {
	t := template.Must(template.ParseFiles("views/errors/" + view + ".html"))
	return t.Execute(c, struct {
		Error string
	}{Error: err.Error()})
}

func (c *Controller) Raw(data string) {
	c.Write([]byte(data))
}

func (c *Controller) Write(data []byte) (int, error) {
	return c.Res.Write(data)
}

func (c *Controller) toHandleFunc(pred func()) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Req = r
		c.Res = w
		c.Params = mux.Vars(r)
		c.Options = r.URL.Query()
		pred()
	}
}

func (c *Controller) Mount(router *mux.Router) {
	for _, route := range c.Routes {
		// TODO: Add prefix field to controller for path prefix
		router.HandleFunc(route.Path, c.toHandleFunc(route.Handler)).
			Methods(route.Method)
	}
}
