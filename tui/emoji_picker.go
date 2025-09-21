// tui/emoji_picker.go
package tui

import (
	"os"
	"path/filepath"
	"strings"

	list "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FileChosenMsg string

type fileItem struct {
	name  string
	path  string
	isDir bool
	ext   string
}

func (f fileItem) Title() string {
	switch {
	case f.isDir:
		return "📁 " + f.name
	case strings.EqualFold(f.ext, ".json"):
		return "{}" + f.name
	default:
		return "📄 " + f.name
	}
}
func (f fileItem) Description() string { return f.path }
func (f fileItem) FilterValue() string { return f.name }

type EmojiPicker struct {
	cwd     string
	allowed string // e.g. ".json"
	lst     list.Model
	width   int
	height  int
}

func NewEmojiPicker(startDir, allowed string, w, h int) EmojiPicker {
	it := readDirItems(startDir, allowed)
	l := list.New(it, list.NewDefaultDelegate(), w, h)
	l.Title = lipgloss.NewStyle().Bold(true).Render("Select a JSON file")
	return EmojiPicker{cwd: startDir, allowed: allowed, lst: l, width: w, height: h}
}

func readDirItems(dir, allowed string) []list.Item {
	ents, _ := os.ReadDir(dir)
	var items []list.Item
	// add “..” to go up
	items = append(items, fileItem{name: "..", path: filepath.Dir(dir), isDir: true})
	for _, e := range ents {
		info, _ := e.Info()
		fi := fileItem{
			name:  info.Name(),
			path:  filepath.Join(dir, info.Name()),
			isDir: e.IsDir(),
			ext:   strings.ToLower(filepath.Ext(info.Name())),
		}
		// allow all dirs; for files, optionally filter by extension (still show but we’ll block on enter)
		items = append(items, fi)
	}
	return items
}

func (m EmojiPicker) Init() tea.Cmd { return nil }

func (m EmojiPicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.lst.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			it, ok := m.lst.SelectedItem().(fileItem)
			if !ok {
				return m, nil
			}
			if it.isDir {
				m.cwd = it.path
				m.lst.SetItems(readDirItems(m.cwd, m.allowed))
				return m, nil
			}
			// file
			if m.allowed == "" || strings.EqualFold(it.ext, m.allowed) {
				// send selection to outer program
				return m, func() tea.Msg { return FileChosenMsg(it.path) }
			}
			// ignore non-allowed file types
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.lst, cmd = m.lst.Update(msg)
	return m, cmd
}

func (m EmojiPicker) View() string { return m.lst.View() }
