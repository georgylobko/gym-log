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
	"github.com/pressly/goose/v3"
)

// var embedMigrations embed.FS

func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	if portString == "" {
		portString = "5000"
	}

	dbHost := os.Getenv("RDS_HOSTNAME")
	dbPort := os.Getenv("RDS_PORT")
	dbUser := os.Getenv("RDS_USERNAME")
	dbPassword := os.Getenv("RDS_PASSWORD")
	dbName := os.Getenv("RDS_DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=enable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	conn.Ping()

	queries := database.New(conn)
	apiCfg := handlers.ApiConfig{
		DB: queries,
	}

	// run migrations
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}
	if err := goose.Up(conn, "sql/schema"); err != nil {
		fmt.Println("Could not run migrations: ", connStr, err)
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	v1Router.Post("/muscle-groups", middlewares.MiddlewareAuth(apiCfg.HandlerCreateMuscleGroup))
	v1Router.Get("/muscle-groups", middlewares.MiddlewareAuth(apiCfg.HandlerGetMuscleGroups))

	v1Router.Post("/exercises", middlewares.MiddlewareAuth(apiCfg.HandlerCreateExercise))
	v1Router.Get("/exercises/{exerciseID}", middlewares.MiddlewareAuth(apiCfg.HandlerGetExercise))
	v1Router.Get("/exercises", middlewares.MiddlewareAuth(apiCfg.HandlerGetExercises))

	v1Router.Post("/workouts", middlewares.MiddlewareAuth(apiCfg.HandlerCreateWorkout))
	v1Router.Get("/workouts", middlewares.MiddlewareAuth(apiCfg.HandlerGetWorkouts))
	v1Router.Put("/workouts", middlewares.MiddlewareAuth(apiCfg.HandlerUpdateWorkout))

	v1Router.Post("/sets", middlewares.MiddlewareAuth(apiCfg.HandlerCreateSet))
	v1Router.Get("/sets", middlewares.MiddlewareAuth(apiCfg.HandlerGetSets))

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
