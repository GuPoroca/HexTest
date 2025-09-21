package tui

import (
	"github.com/GuPoroca/HexTest/pkg/jsonOperations"
	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	Current tea.Model
	width   int
	height  int
}

func (r RootModel) Init() tea.Cmd {
	return r.Current.Init()
}

func (r RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		r.width, r.height = msg.Width, msg.Height
		// forward to current child
		var cmd tea.Cmd
		r.Current, cmd = r.Current.Update(msg)
		return r, cmd
	case FileChosenMsg:
		// Load project (value or pointer is fine; we need a pointer here)
		proj := jsonOperations.ReadJSON(string(msg)) // change to return *Project if possible
		// If ReadJSON returns a value, take its address:
		// p := proj; projPtr := &p

		dash := New(proj, r.width/3, r.height) // pass pointer

		// Kick off the runner in background (mutates proj)
		go proj.ExecuteProject()

		// Dashboard's Init() already starts listening to bus.CheckEvents
		return dash, dash.Init()
	}

	var cmd tea.Cmd
	r.Current, cmd = r.Current.Update(msg)
	return r, cmd
}

func (r RootModel) View() string {
	return r.Current.View()
}
