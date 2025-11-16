package commands

import (
	"io"
	"math"
	"strings"
	"time"
	"todo/constants"
	"todo/models"

	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gorm.io/gorm"
)

const (
	listHeight = 14
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(0).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	tagWidth          = 12
	todoWidth         = 40

	arrowEnabledStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	arrowDisabledStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	pageActiveStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	pageInactiveStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
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
			return selectedItemStyle.Render("＞ " + s[0])
		}
	}

	fmt.Fprint(w, fn(str))
}

type Model struct {
	list        list.Model
	quitting    bool
	todos       []models.Todo
	db          *gorm.DB
	totalPages  int
	currentPage int
}

func NewModel(db *gorm.DB) (Model, error) {
	count, err := countAll(db)
	if err != nil {
		return Model{}, err
	}

	totalPages := int(math.Ceil(float64(count) / float64(constants.DefaultLimit)))

	todos, err := read(db, constants.DefaultLimit, 0, false)
	if err != nil {
		return Model{}, err
	}

	items := make([]list.Item, len(*todos))
	for i, t := range *todos {
		checkbox := "[ ]"
		if t.CompletedAt != nil {
			checkbox = "[✓]"
		}

		tag := "—"
		if t.Tag != nil {
			tag = t.Tag.Tag
		}

		items[i] = item{
			Checkbox:       checkbox,
			Tag:            tag,
			DisplayedTitle: t.GetDisplayTitle(i == 0),
		}
	}

	l := list.New(items, itemDelegate{}, 0, listHeight)

	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)

	return Model{
		list:        l,
		todos:       *todos,
		db:          db,
		totalPages:  totalPages,
		currentPage: 0,
	}, nil
}

func (m *Model) refreshTodos() error {
	offset := m.currentPage * constants.DefaultLimit
	todos, err := read(m.db, constants.DefaultLimit, offset, false)
	if err != nil {
		return err
	}

	m.todos = *todos

	items := make([]list.Item, len(*todos))
	for i, t := range *todos {
		checkbox := "[ ]"
		if t.CompletedAt != nil {
			checkbox = "[✓]"
		}

		tag := "—"
		if t.Tag != nil {
			tag = t.Tag.Tag
		}

		items[i] = item{
			Checkbox:       checkbox,
			Tag:            tag,
			DisplayedTitle: t.GetDisplayTitle(i == 0),
		}
	}

	m.list.SetItems(items)
	m.list.Select(0)
	return nil
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

		case " ":
			selectedIndex := m.list.Index()
			if selectedIndex < 0 || selectedIndex >= len(m.todos) {
				break
			}

			todo := &m.todos[selectedIndex]

			if todo.CompletedAt == nil {
				now := time.Now()
				todo.CompletedAt = &now
			} else {
				todo.CompletedAt = nil
			}

			cb := "[ ]"
			if todo.CompletedAt != nil {
				cb = "[✓]"
			}

			tag := "—"
			if todo.Tag != nil {
				tag = todo.Tag.Tag
			}

			m.list.SetItem(selectedIndex, item{
				Checkbox:       cb,
				Tag:            tag,
				DisplayedTitle: todo.Title,
			})

			return m, nil

		case "enter", "q", "ctrl+c", "esc":
			idx := m.list.Index()
			if idx >= 0 && idx < len(m.todos) {
				todo := &m.todos[idx]
				if m.db != nil {
					_ = m.db.Save(todo).Error
				}
			}

			return m, tea.Quit

		case "left", "h":
			if m.totalPages > 1 && m.currentPage > 0 {
				m.currentPage--
				_ = m.refreshTodos()
			}
			return m, nil

		case "right", "l":
			if m.totalPages > 1 && m.currentPage < m.totalPages-1 {
				m.currentPage++
				_ = m.refreshTodos()
			}
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func renderPagination(currentPage, totalPages int) string {
	current := currentPage + 1

	if totalPages <= 1 {
		return ""
	}

	var pages []string

	makePage := func(n int) string {
		if n == current {
			return pageActiveStyle.Render(fmt.Sprintf("[%d]", n))
		}
		return pageInactiveStyle.Render(fmt.Sprintf("%d", n))
	}

	if current > 3 {
		pages = append(pages, makePage(1))
		pages = append(pages, "...")
	}

	start := max(1, current-1)
	end := min(totalPages, current+1)

	for i := start; i <= end; i++ {
		pages = append(pages, makePage(i))
	}

	if current < totalPages-2 {
		pages = append(pages, "...")
		pages = append(pages, makePage(totalPages))
	}

	return strings.Join(pages, " ")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) View() string {
	if m.quitting {
		return quitTextStyle.Render("Exiting program.")
	}

	leftRaw, rightRaw := "◀", "▶"

	var leftArrow, rightArrow string

	if m.totalPages <= 1 {
		leftArrow = arrowDisabledStyle.Render(leftRaw)
		rightArrow = arrowDisabledStyle.Render(rightRaw)
	} else {
		if m.currentPage == 0 {
			leftArrow = arrowDisabledStyle.Render(leftRaw)
		} else {
			leftArrow = arrowEnabledStyle.Render(leftRaw)
		}

		if m.currentPage >= m.totalPages-1 {
			rightArrow = arrowDisabledStyle.Render(rightRaw)
		} else {
			rightArrow = arrowEnabledStyle.Render(rightRaw)
		}
	}

	m.list.Title = "✏️ My Todos"

	pagination := fmt.Sprintf(
		"\n\n  %s  %s  %s",
		leftArrow,
		renderPagination(m.currentPage, m.totalPages),
		rightArrow,
	)

	return "\n" + m.list.View() + pagination
}
