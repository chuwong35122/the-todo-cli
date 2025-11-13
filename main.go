package main

import (
	"os"
	"todo/commands"
	"todo/libs"
	"todo/models"
)

func main() {

	db, err := libs.InitDB()

	if err != nil {
		libs.Logger.Fatal().
			Err(err).
			Msg("Error initializing the database")
	}

	db.AutoMigrate(&models.Todo{}, &models.TodoTag{})

	cmd := commands.Prepare(db)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
