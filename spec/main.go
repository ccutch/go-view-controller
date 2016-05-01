package main

import (
	"fmt"
	"net/http"

	"github.com/ccutch/view-controller"
)

type Application struct {
	*controller.Controller
}

func (this Application) Home() {
	this.Render("spec/views/application.html")
}

func (this Application) Welcome() {
	this.Data = struct {
		Name    string
		Awesome bool
	}{
		Name:    this.Params["name"],
		Awesome: len(this.Options["awesome"]) == 1,
	}
	this.Render("spec/views/welcome.html")
}

func main() {
	a := Application{new(controller.Controller)}

	a.Routes = []controller.Route{
		controller.Route{
			Path:    "/welcome/{name}",
			Method:  "GET",
			Handler: a.Welcome,
		},
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
