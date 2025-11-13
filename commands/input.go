package commands

import (
	"fmt"
	"os"
	"todo/libs"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var (
	createDesc string
	tag        string
)

func Prepare(db *gorm.DB) *cobra.Command {
	RootCmd := &cobra.Command{
		Use:   "todo",
		Short: "The Todo CLI app",
		Run: func(cmd *cobra.Command, args []string) {
			if createDesc != "" {
				err := create(db, createDesc, tag)
				if err != nil {
					libs.Logger.Fatal().
						Err(err).
						Msg("Error creating todo")

				}

				todos := read(db)
				if todos == nil {
					libs.Logger.Info().Msg("Todo List doesn't exist. Create one!")
					return
				}

				for _, todo := range *todos {
					fmt.Println(todo.Title)
					fmt.Println(todo.Tag.Tag)
					fmt.Println(todo.CreatedAt)
				}

				m := NewModel(todos, db)
				p := tea.NewProgram(m)

				if _, err := p.Run(); err != nil {
					libs.Logger.Fatal().
						Err(err).
						Msg("Error running program")

					os.Exit(1)
				}
			}

			if len(args) == 0 {
				todos := read(db)
				if todos == nil {
					return
				}

				m := NewModel(todos, db)
				p := tea.NewProgram(m)

				if _, err := p.Run(); err != nil {
					libs.Logger.Fatal().
						Err(err).
						Msg("Error running program")

					os.Exit(1)
				}
			}

			fmt.Println("No command provided. Use --create <description> to create a todo.")
			cmd.Help()
			os.Exit(1)
		},
	}

	RootCmd.Flags().StringVarP(&createDesc, "create", "c", "", "Create a new todo with the given description")
	RootCmd.Flags().StringVarP(&tag, "tag", "t", "", "Optional tag for categorizing todo item")

	return RootCmd
}
