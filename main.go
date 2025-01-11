package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/akrylysov/algnhsa"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/xonha/huma-chi/pages"
)

//go:embed all:public
var embeddedFiles embed.FS

type HelloInput struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}
type HelloOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func hello(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
	resp := &HelloOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
	return resp, nil
}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	assets, _ := fs.Sub(embeddedFiles, "public")
	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.FS(assets))))

	router.Get("/", templ.Handler(pages.Page()).ServeHTTP)
	router.Get("/hello", templ.Handler(pages.Response()).ServeHTTP)

	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
	huma.Get(api, "/hello/{name}", hello)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		algnhsa.ListenAndServe(router, nil)
	} else {
		fmt.Println("Server starting locally on port 3000...")
		http.ListenAndServe(":3000", router)
	}
}
