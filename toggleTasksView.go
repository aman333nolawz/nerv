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

type toggleItem struct {
	index int
}

func (i toggleItem) FilterValue() string {
	return tasks[i.index].Desc
}

type toggleItemDelegate struct {
	styles *styles
}

func (d toggleItemDelegate) Height() int                             { return 1 }
func (d toggleItemDelegate) Spacing() int                            { return 0 }
func (d toggleItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d toggleItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(toggleItem)
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

type toggleModel struct {
	list     list.Model
	choice   string
	styles   styles
	quitting bool
}

func initialToggleModel() toggleModel {
	var items []list.Item

	for i := range tasks {
		items = append(items, toggleItem{i})
	}

	const defaultWidth = 20

	l := list.New(items, toggleItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Toggle tasks"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	m := toggleModel{list: l}
	m.updateStyles(true) // default to dark styles.
	return m
}

func (m *toggleModel) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.SetDelegate(toggleItemDelegate{styles: &m.styles})
}

func (m toggleModel) Init() tea.Cmd {
	return nil
}

func (m toggleModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				selected, ok := m.list.SelectedItem().(toggleItem)
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

func (m toggleModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}
