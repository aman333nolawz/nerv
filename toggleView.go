package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const listHeight = 14

type styles struct {
	item         lipgloss.Style
	selectedItem lipgloss.Style
	pagination   lipgloss.Style
}

func newStyles(darkBG bool) styles {
	var s styles
	s.item = lipgloss.NewStyle().PaddingLeft(2)
	s.selectedItem = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).PaddingLeft(2)
	s.pagination = list.DefaultStyles(darkBG).PaginationStyle.PaddingLeft(2)
	return s
}

type item struct {
	index int
}

func (i item) FilterValue() string {
	return tasks[i.index].Desc
}

type itemDelegate struct {
	styles *styles
}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	task := tasks[i.index]
	text := task.Render()
	str := fmt.Sprintf("%d. %s", index+1, text)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	styles   styles
	quitting bool
}

func initialModel() model {
	var items []list.Item

	for i := range tasks {
		items = append(items, item{i})
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, min(len(tasks)+3, listHeight))
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetShowHelp(false)

	m := model{list: l}
	m.updateStyles(true) // default to dark styles.
	return m
}

func (m *model) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.SetDelegate(itemDelegate{styles: &m.styles})
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "space", "enter":
			// prevent accidentally toggling while pressing enter or space during filtering
			if m.list.FilterState() != list.Filtering {
				selected, ok := m.list.SelectedItem().(item)
				if !ok {
					return m, nil
				}
				task := &tasks[selected.index]
				task.Done = !task.Done
				m.list.SetItem(m.list.GlobalIndex(), selected)
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	return tea.NewView(m.list.View())
}
