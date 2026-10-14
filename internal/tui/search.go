package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/estaples/lsproj/internal/scanner"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			PaddingLeft(1)
)

type projectItem struct {
	project scanner.Project
}

func (i projectItem) Title() string       { return i.project.Name }
func (i projectItem) Description() string { return i.project.Path }
func (i projectItem) FilterValue() string { return i.project.Name }

type searchModel struct {
	list     list.Model
	selected *scanner.Project
}

func newSearchModel(projects []scanner.Project) searchModel {
	items := make([]list.Item, len(projects))
	for i, p := range projects {
		items[i] = projectItem{project: p}
	}

	l := list.New(items, list.NewDefaultDelegate(), 80, 24)
	l.Title = "lsproj"
	l.Styles.Title = titleStyle
	l.FilterInput.Prompt = "Search: "
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.FilterInput.Focus()

	return searchModel{list: l}
}

func (m searchModel) Init() tea.Cmd {
	return nil
}

func (m searchModel) Update(msg tea.Msg) (searchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if item, ok := m.list.SelectedItem().(projectItem); ok {
				m.selected = &item.project
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m searchModel) View() string {
	return m.list.View()
}

func (m searchModel) resize(w, h int) searchModel {
	m.list.SetSize(w, h)
	return m
}
