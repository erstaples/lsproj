package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/estaples/lsproj/internal/action"
)

var (
	actionTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				PaddingLeft(1).
				PaddingBottom(1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170"))
)

type actionModel struct {
	actions []action.Action
	cursor  int
	chosen  *action.Action
	width   int
	height  int
}

func newActionModel(actions []action.Action) actionModel {
	return actionModel{
		actions: actions,
	}
}

func (m actionModel) Init() tea.Cmd {
	return nil
}

func (m actionModel) Update(msg tea.Msg) (actionModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.actions)-1 {
				m.cursor++
			}
		case "enter":
			chosen := m.actions[m.cursor]
			m.chosen = &chosen
			return m, nil
		default:
			// Number keys select actions directly (1-9)
			if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
				idx := int(msg.String()[0]-'1')
				if idx < len(m.actions) {
					chosen := m.actions[idx]
					m.chosen = &chosen
					return m, nil
				}
			}
		}
	}
	return m, nil
}

func (m actionModel) View() string {
	var b strings.Builder

	b.WriteString(actionTitleStyle.Render("Choose an action"))
	b.WriteString("\n")

	for i, a := range m.actions {
		num := fmt.Sprintf("%d.", i+1)
		if i == m.cursor {
			fmt.Fprintf(&b, "  %s %s %s\n",
				cursorStyle.Render("▸"),
				selectedStyle.Render(num),
				selectedStyle.Render(a.Name()))
		} else {
			fmt.Fprintf(&b, "    %s %s\n",
				normalStyle.Render(num),
				normalStyle.Render(a.Name()))
		}
	}

	b.WriteString("\n")
	b.WriteString(normalStyle.Render("1-9 or ↑/↓ + enter • esc back"))

	return b.String()
}

func (m actionModel) resize(w, h int) actionModel {
	m.width = w
	m.height = h
	return m
}

func (m actionModel) reset() actionModel {
	m.cursor = 0
	m.chosen = nil
	return m
}
