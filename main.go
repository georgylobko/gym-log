package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/georgylobko/gym-log/internal/database"
	"github.com/georgylobko/gym-log/internal/handlers"
	"github.com/georgylobko/gym-log/internal/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in env")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the env")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}

	queries := database.New(conn)
	apiCfg := handlers.ApiConfig{
		DB: queries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Post("/muscle-groups", middlewares.MiddlewareAuth(apiCfg.HandlerCreateMuscleGroup))
	v1Router.Get("/muscle-groups", middlewares.MiddlewareAuth(apiCfg.HandlerGetMuscleGroups))

	v1Router.Post("/exercises", middlewares.MiddlewareAuth(apiCfg.HandlerCreateExercise))
	v1Router.Get("/exercises/{exerciseID}", middlewares.MiddlewareAuth(apiCfg.HandlerGetExercise))

	v1Router.Post("/register", apiCfg.HandlerRegister)
	v1Router.Post("/login", apiCfg.HandlerLogin)
	v1Router.Get("/session", middlewares.MiddlewareAuth(apiCfg.HandlerSession))
	v1Router.Get("/logout", apiCfg.HandlerLogout)

	router.Mount("/v1", v1Router)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Printf("Server starting on port %v", portString)
	fmt.Println("")

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
