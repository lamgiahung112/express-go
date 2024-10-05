package main

import (
	"net/http"

	expressgo "github.com/lamgiahung112/express-go"
)

func main() {
	config := &expressgo.ExpressConfig{
		Host: "localhost",
		Port: 8080,
	}
	app := expressgo.NewExpress(config)
	route := app.NewRoute("/hello")
	route.Route("/world").Use(func(r *http.Request, context *expressgo.RequestContext) bool {
		context.Set("User", "Hello")
		return true
	}).Get(func(w http.ResponseWriter, r *http.Request, context *expressgo.RequestContext) {
		w.Write([]byte(context.Get("User").(string)))
	})
	app.StartServer()
}
