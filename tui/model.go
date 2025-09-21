package tui

import (
	"fmt"
	"github.com/GuPoroca/HexTest/pkg/typeDefines"
	list "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"strconv"
	"strings"
)

type itemKind int

const (
	kindSuite itemKind = iota
	kindTest
	kindAssert
	kindCheck
)

type treeItem struct {
	kind        itemKind
	title       string
	desc        string
	operand     string
	expected    []any
	checkStatus []int // assuming you store status codes as ints: -2, -1, 0, 1
}

func (it treeItem) Title() string       { return it.title }
func (it treeItem) Description() string { return it.desc }
func (it treeItem) FilterValue() string { return it.title }

// ---------- Styling ----------
var (
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00c8ff"))
	styleSuite   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#d1d5db"))
	styleTest    = lipgloss.NewStyle().Foreground(lipgloss.Color("#e5e7eb"))
	styleAssert  = lipgloss.NewStyle().Foreground(lipgloss.Color("#facc15")) // yellow-ish
	styleCheck   = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ca3af"))
	styleSel     = lipgloss.NewStyle().Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color("#00c8ff"))
	styleBulletS = "📁 "
	styleBulletT = "🧪 "
	styleBulletA = "🔎 "
	styleBulletC = "✔ "
)

var detailStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1).
	Align(lipgloss.Center).        // horizontal center
	AlignVertical(lipgloss.Center) // vertical center

// ---------- Custom Delegate ----------
type delegate struct{}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	it, ok := listItem.(treeItem)
	if !ok {
		fmt.Fprint(w, "?")
		return
	}

	var line string
	switch it.kind {
	case kindSuite:
		line = styleSuite.Render(styleBulletS + it.title)
	case kindTest:
		txt := "  " + styleBulletT + it.title
		if it.desc != "" {
			txt += "  " + dim("— "+it.desc)
		}
		line = styleTest.Render(txt)
	case kindAssert:
		txt := "    " + styleBulletA + it.title
		line = styleAssert.Render(txt)
	case kindCheck:
		icons := statusesToIcons(it.checkStatus)
		expected := ""
		if len(it.expected) > 0 {
			stringsIt := make([]string, len(it.expected))
			for i, a := range it.expected {
				stringsIt[i] = typeDefines.StringifyMyAny(a)
			}
			expected = " → [" + strings.Join(stringsIt, ", ") + "]"
		}
		txt := fmt.Sprintf("      %s%s %s%s", styleBulletC, it.operand, it.title, expected)
		if icons != "" {
			txt += " " + icons
		}
		line = styleCheck.Render(txt)
	}

	if index == m.Index() && it.kind != kindSuite {
		line = styleSel.Render(stripANSI(line))
	}
	fmt.Fprint(w, line)
}

func dim(s string) string {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")).Render(s)
}

func stripANSI(s string) string {
	return s
}

// ---------- Model ----------
type Model struct {
	project typeDefines.Project
	list    list.Model
	results viewport.Model
	focus   int
}

func New(project typeDefines.Project, width, height int) Model {
	items := buildItems(project)

	l := list.New(items, delegate{}, width, height)
	l.Title = styleTitle.Render(project.Name)
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()

	vp := viewport.New(width/3, height)
	vp.SetContent(renderResults(project))

	return Model{
		project: project,
		list:    l,
		results: vp,
		focus:   0,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		colWidth := msg.Width / 3
		m.list.SetSize(colWidth, msg.Height)
		m.results.Width = colWidth
		m.results.Height = msg.Height

	case TreeUpdateMsg:
		m.project = msg.Project
		m.list.SetItems(buildItems(m.project))
		m.results.SetContent(renderResults(m.project))

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.focus = (m.focus + 1) % 2 // toggle focus
		}
	}
	var cmd tea.Cmd
	if m.focus == 0 {
		m.list, cmd = m.list.Update(msg)
	} else {
		m.results, cmd = m.results.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	totalWidth := m.list.Width() * 3
	colWidth := totalWidth / 3

	colStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1).
		Width(colWidth)

	treeView := colStyle.Render(m.list.View())
	detailView := detailStyle.Width(colWidth).Height(m.list.Height()).Render(m.selectedDetails())
	resultsView := colStyle.Render(m.results.View())

	return lipgloss.JoinHorizontal(lipgloss.Top, treeView, detailView, resultsView)
}

// ---------- Helpers ----------
func buildItems(p typeDefines.Project) []list.Item {
	var items []list.Item
	for _, s := range p.Suites {
		items = append(items, treeItem{kind: kindSuite, title: s.Name, desc: strconv.Itoa(len(s.Tests))})

		for _, t := range s.Tests {
			desc := fmt.Sprintf("%d asserts", len(t.Asserts))
			items = append(items, treeItem{kind: kindTest, title: t.Name, desc: desc})

			for _, a := range t.Asserts {
				items = append(items, treeItem{kind: kindAssert, title: a.Field, desc: typeDefines.StringifyMyAny(a.FieldResponseValue)})

				for _, c := range a.Checks {
					items = append(items, treeItem{
						kind:        kindCheck,
						operand:     c.Operand,
						expected:    c.Expected,
						checkStatus: c.Passed, // assume []int
						desc:        typeDefines.StringifyMyAny(a.FieldResponseValue),
					})
				}
			}
		}
	}
	return items
}

func statusesToIcons(ss []int) string {
	var b strings.Builder
	for _, s := range ss {
		switch s {
		case -1:
			b.WriteString("💥")
		case 0:
			b.WriteString("❌")
		case 1:
			b.WriteString("✅")
		}
	}
	return b.String()
}

func (m Model) selectedDetails() string {
	it, ok := m.list.SelectedItem().(treeItem)
	if !ok {
		return "\n\nNo item selected"
	}
	switch it.kind {
	case kindSuite:
		return fmt.Sprintf("\n\nSuite: %s\nContains %s tests", it.title, it.desc)
	case kindTest:
		return fmt.Sprintf("\n\nTest: %s\n%s", it.title, it.desc)
	case kindAssert:
		return fmt.Sprintf("\n\nAssert: %s,\n Value: %s", it.title, it.desc)
	case kindCheck:
		expected := "none"
		if len(it.expected) > 0 {
			stringsIt := make([]string, len(it.expected))
			for i, a := range it.expected {
				stringsIt[i] = typeDefines.StringifyMyAny(a)
			}
			expected = strings.Join(stringsIt, ", ")
		}
		return fmt.Sprintf(
			"\n\n\nResponse Value: %s\nOperand: %s\nExpected: %s\nStatuses: %s",
			it.desc,
			it.operand,
			expected,
			statusesToIcons(it.checkStatus),
		)
	}
	return "Unknown item"
}

func renderResults(p typeDefines.Project) string {
	checks, ch_ps, ch_fl, ch_bk := 0, 0, 0, 0
	var b strings.Builder

	for _, s := range p.Suites {
		for _, t := range s.Tests {
			for _, a := range t.Asserts {
				for _, c := range a.Checks {
					for i, status := range c.Passed {
						checks++
						switch status {
						case 1:
							ch_ps++
						case 0:
							ch_fl++
							fmt.Fprintf(&b, "Comparisson: %s\n", fmt.Sprintf("%s %s %s",
								typeDefines.StringifyMyAny(a.FieldResponseValue), c.Operand, typeDefines.StringifyMyAny(c.Expected[i])))
							fmt.Fprintf(&b, "On Assert: %s\n", a.Field)
							fmt.Fprintf(&b, "On Test: %s\n", t.Name)
							fmt.Fprintf(&b, "On Suite: %s\n", s.Name)
							fmt.Fprintf(&b, "FAILED\n\n")
						case -1:
							ch_bk++
							fmt.Fprintf(&b, "Comparisson: %s\n", fmt.Sprintf("%s %s %s",
								typeDefines.StringifyMyAny(a.FieldResponseValue), c.Operand, typeDefines.StringifyMyAny(c.Expected[i])))
							fmt.Fprintf(&b, "On Assert: %s\n", a.Field)
							fmt.Fprintf(&b, "On Test: %s\n", t.Name)
							fmt.Fprintf(&b, "On Suite: %s\n", s.Name)
							fmt.Fprintf(&b, "BROKEN\n\n")
						}
					}
				}
			}
		}
	}

	fmt.Fprintf(&b, "Number of Checks Made: %v\n", checks)
	fmt.Fprintf(&b, "Number of Checks Passed: %v\n", ch_ps)
	fmt.Fprintf(&b, "Number of Checks Failed: %v\n", ch_fl)
	fmt.Fprintf(&b, "Number of Checks Broken: %v\n", ch_bk)

	return b.String()
}
