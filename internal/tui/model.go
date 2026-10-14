package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/estaples/lsproj/internal/action"
	"github.com/estaples/lsproj/internal/scanner"
)

type view int

const (
	searchView view = iota
	actionView
)

type Result struct {
	Action  action.Action
	Project scanner.Project
}

type model struct {
	projects []scanner.Project
	actions  []action.Action
	view     view
	search   searchModel
	picker   actionModel
	result   *Result
	width    int
	height   int
	quitting bool
}

func New(projects []scanner.Project) model {
	actions := action.All()
	m := model{
		projects: projects,
		actions:  actions,
		view:     searchView,
	}
	m.search = newSearchModel(projects)
	m.picker = newActionModel(actions)
	return m
}

func (m model) Init() tea.Cmd {
	return m.search.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.search = m.search.resize(msg.Width, msg.Height)
		m.picker = m.picker.resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.view == actionView {
				m.view = searchView
				return m, nil
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.view {
	case searchView:
		s, cmd := m.search.Update(msg)
		m.search = s
		if m.search.selected != nil {
			m.view = actionView
			m.picker = m.picker.reset()
			return m, nil
		}
		return m, cmd

	case actionView:
		p, cmd := m.picker.Update(msg)
		m.picker = p
		if m.picker.chosen != nil {
			m.result = &Result{
				Action:  *m.picker.chosen,
				Project: *m.search.selected,
			}
			m.quitting = true
			return m, tea.Quit
		}
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.quitting && m.result != nil {
		return ""
	}
	switch m.view {
	case searchView:
		return m.search.View()
	case actionView:
		return m.picker.View()
	}
	return ""
}

func Run(projects []scanner.Project) (*Result, error) {
	m := New(projects)
	p := tea.NewProgram(m, tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		return nil, err
	}
	fm := final.(model)
	return fm.result, nil
}
