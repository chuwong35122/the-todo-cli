package commands

import (
	"fmt"
	"io"
	"time"
	"todo/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

const listHeight = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	tagWidth          = 12
	todoWidth         = 40
)

type item struct {
	Checkbox       string
	Tag            string
	DisplayedTitle string
}

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	todoCol := fmt.Sprintf("%-*s", todoWidth, i.DisplayedTitle)
	tagCol := fmt.Sprintf("%-*s", tagWidth, i.Tag)

	str := fmt.Sprintf("%s  %s  %s", i.Checkbox, todoCol, tagCol)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	list     list.Model
	quitting bool
	todos    *[]models.Todo
	db       *gorm.DB
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case " ":
			selectedIndex := m.list.Index()
			if selectedIndex < 0 || selectedIndex >= len(*m.todos) {
				// nothing selected; let list handle it
				break
			}

			todo := &(*m.todos)[selectedIndex]

			if todo.CompletedAt == nil {
				now := time.Now()
				todo.CompletedAt = &now
			} else {
				todo.CompletedAt = nil
			}

			cb := "[ ]"
			if todo.CompletedAt != nil {
				cb = "[X]"
			}

			tag := todo.GetTag()
			if tag == "" {
				tag = "—"
			}

			m.list.SetItem(selectedIndex, item{
				Checkbox:       cb,
				Tag:            tag,
				DisplayedTitle: todo.Title,
			})

			return m, nil

		case "enter":
			idx := m.list.Index()
			if idx < 0 || idx >= len(*m.todos) {
				break
			}

			todo := &(*m.todos)[idx]
			if m.db != nil {
				_ = m.db.Save(todo).Error
			}

			return m, tea.Quit
		}

		// IMPORTANT: don't return here for unhandled keys to be handled by list.Update
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Exiting program.")
	}
	return "\n" + m.list.View()
}

func NewModel(todos *[]models.Todo, db *gorm.DB) Model {
	items := make([]list.Item, len(*todos))

	for i, t := range *todos {
		checkbox := "[ ]"
		if t.CompletedAt != nil {
			checkbox = "[X]"
		}

		tag := "—"
		if t.Tag != nil {
			tag = t.Tag.Tag
		}

		items[i] = item{
			Checkbox:       checkbox,
			Tag:            tag,
			DisplayedTitle: t.GetDisplayTitle(i == 0), // DB order DESC
		}
	}

	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Title = "Your Todos"
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)

	return Model{list: l, todos: todos, db: db}
}
