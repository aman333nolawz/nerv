package main

import (
	"fmt"
	"io"
	"strings"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

type removeItem struct {
	index  int
	remove bool
}

func (i removeItem) FilterValue() string {
	return tasks[i.index].Desc
}

type removeItemDelegate struct {
	styles *styles
}

func (d removeItemDelegate) Height() int                             { return 1 }
func (d removeItemDelegate) Spacing() int                            { return 0 }
func (d removeItemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d removeItemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(removeItem)
	if !ok {
		return
	}

	if i.index < 0 || i.index >= len(tasks) {
		return
	}

	task := tasks[i.index]
	text := task.Render()

	if i.remove {
		text = removeSelectedStyle.Render(text)
	}
	str := fmt.Sprintf("%d. %s", index+1, text)

	fn := d.styles.item.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return d.styles.selectedItem.Render(strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type removeModel struct {
	list     list.Model
	choice   string
	styles   styles
	quitting bool
}

func initialRemoveModel() removeModel {
	const defaultWidth = 20

	l := list.New(removeItems(), removeItemDelegate{}, defaultWidth, listHeight)
	l.Title = "Remove tasks"
	l.SetShowStatusBar(false)
	l.SetShowHelp(false)

	m := removeModel{list: l}
	m.updateStyles(true) // default to dark styles.
	return m
}

func removeItems() []list.Item {
	var items []list.Item
	for i := range tasks {
		items = append(items, removeItem{i, false})
	}
	return items
}

func (m *removeModel) updateStyles(isDark bool) {
	m.styles = newStyles(isDark)
	m.list.Styles.PaginationStyle = m.styles.pagination
	m.list.SetDelegate(removeItemDelegate{styles: &m.styles})
}

func (m removeModel) Init() tea.Cmd {
	return nil
}

func (m removeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyPressMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true

			toRemoveIds := []int{}
			for _, item := range m.list.Items() {
				if item.(removeItem).remove {
					toRemoveIds = append(toRemoveIds, item.(removeItem).index+1)
				}
			}
			removeTasksWithoutLog(toRemoveIds...)

			return m, tea.Quit

		case "space", "enter":
			if m.list.FilterState() != list.Filtering {
				selected, ok := m.list.SelectedItem().(removeItem)
				if !ok {
					return m, nil
				}
				cmd := m.list.SetItem(m.list.GlobalIndex(), removeItem{index: selected.index, remove: !selected.remove})
				return m, cmd
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m removeModel) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}
	v := tea.NewView(m.list.View())
	v.AltScreen = true
	return v
}
