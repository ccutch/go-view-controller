package main

import (
	"fmt"
	"net/http"

	"github.com/ccutch/view-controller"
)

type Application struct {
	*controller.Controller
}

// [GET] / => Render homepage
func (this Application) Home() {
	this.Render("spec/views/application.html")
}

func main() {
	a := Application{new(controller.Controller)}

	a.Routes = []controller.Route{
		controller.Route{
			Path:    "/",
			Method:  "GET",
			Handler: a.Home,
		},
	}

	http.Handle("/", a)
	fmt.Println("Server online")
	http.ListenAndServe(":8080", nil)
}
