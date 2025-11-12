package commands

import (
	"fmt"
	"os"

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
			if len(args) == 0 {
				todos := read(db)
				if todos == nil {
					return
				}

				m := NewModel(todos)
				p := tea.NewProgram(m)

				if _, err := p.Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
			}

			if createDesc != "" {
				fmt.Printf("Creating todo: %s with tag: %s\n", createDesc, tag)
				err := create(db, createDesc, tag)
				if err != nil {
					fmt.Println("Error creating todo:", err)
				}

				todos := read(db)
				if todos == nil {
					return
				}

				m := NewModel(todos)
				p := tea.NewProgram(m)

				if _, err := p.Run(); err != nil {
					fmt.Println("Error running program:", err)
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
