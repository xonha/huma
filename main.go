package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/aws/aws-lambda-go/lambda"
	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func hello(ctx context.Context, input *HelloInput) (*HelloOutput, error) {
	resp := &HelloOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
	return resp, nil
}

var chiLambda *chiadapter.ChiLambda

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	router.Get("/", templ.Handler(pages.Page()).ServeHTTP)
	router.Get("/hello", templ.Handler(pages.Response()).ServeHTTP)

	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
	huma.Get(api, "/hello/{name}", hello)

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		chiLambda = chiadapter.New(router)
		lambda.Start(chiLambda.ProxyWithContext)
	} else {
		// Local development
		fmt.Println("Server starting locally on port 3000...")
		http.ListenAndServe(":3000", router)
	}
}
