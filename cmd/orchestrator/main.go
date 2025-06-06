package main

import (
	"calculator_go/internal/http/middlewares"
	"calculator_go/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	authHandler "calculator_go/internal/http/handlers/auth"
	calcHandler "calculator_go/internal/http/handlers/expression"

	_ "github.com/joho/godotenv/autoload"
)

// Orchestrator - the main app server, which directly operates with database
// It can be described as "manager" of services

func main() {
	ctx := context.Background()

	db, err := storage.New("./database/storage.db")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mainPageHandler := middlewares.AuthorizeJWTToken(http.FileServer(http.Dir("frontend/main")))
	authPageHandler := http.StripPrefix("/auth", http.FileServer(http.Dir("frontend/auth")))

	mux.Handle("/", mainPageHandler)
	mux.Handle("/auth/", authPageHandler)

	mux.Handle("POST /auth/signup/", authHandler.RegisterUserHandler(ctx, db))
	mux.Handle("POST /auth/login/", authHandler.LoginUserHandler(ctx, db))

	mux.Handle("POST /expression/", middlewares.AuthorizeJWTToken(calcHandler.CreateExpressionHandler(ctx, db)))
	mux.Handle("GET /expression/", middlewares.AuthorizeJWTToken(calcHandler.GetExpressionsHandler(ctx, db)))
	mux.Handle("DELETE /expression/{id}/", middlewares.AuthorizeJWTToken(calcHandler.DeleteExpressionHandler(ctx, db)))

	host, ok := os.LookupEnv("ORCHESTRATOR_HOST")
	if !ok {
		log.Print("ORCHESTRATOR_HOST not set, using 0.0.0.0")
		host = "0.0.0.0"
	}

	port, ok := os.LookupEnv("ORCHESTRATOR_PORT")
	if !ok {
		log.Print("ORCHESTRATOR_PORT not set, using 8080")
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Printf("running Orchestrator server at %s", addr)
	go log.Fatal(server.ListenAndServe())

	log.Print("Something went wrong...")
}
