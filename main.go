package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/crisp-coder/chirpy/internal/api"
	"github.com/crisp-coder/chirpy/internal/config"
	"github.com/crisp-coder/chirpy/internal/database"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	api_cfg := api.ApiConfig{
		Db:         dbQueries,
		JWT_SECRET: cfg.JWT_SECRET,
	}

	logFile, err := api.SetupLogging("application.log")
	if err != nil {
		log.Fatalf("Failed to set up logging: %v", err)
		os.Exit(1)
	}
	defer func() {
		err := logFile.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing log file")
		}
	}()
	log.Println("log start")

	server := api.MakeServer(&api_cfg)
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
