package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/xonha/huma-chi/pages"
)

type HelloInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}
type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func helloHTML(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	templ.Handler(pages.Layout(name)).ServeHTTP(w, r)
}

func main() {
	router := chi.NewRouter()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	huma.Get(api, "/hello/{name}", func(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
		resp := &HelloOutput{}
		resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
		return resp, nil
	})

	router.Get("/", helloHTML)

	fmt.Println("Server starting/reloading...")
	http.ListenAndServe(":3000", router)
}
