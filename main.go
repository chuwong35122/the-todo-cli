package main

import (
	"fmt"
	"os"
	"todo/commands"
	"todo/libs"
	"todo/models"

	"github.com/rs/zerolog"
)

func main() {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "",
		NoColor:    false,
		FormatLevel: func(i any) string {
			return ""
		},
		FormatTimestamp: func(i any) string {
			return ""
		},
		FormatMessage: func(i any) string {
			return fmt.Sprintf("%s", i)
		},
	}

	logger := zerolog.New(output).With().Logger()
	db, err := libs.InitDB()

	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("Error initializing the database")
	}

	db.AutoMigrate(&models.Todo{}, &models.TodoTag{})

	cmd := commands.Prepare(db)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
