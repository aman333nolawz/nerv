package main

import (
	"os"

	"charm.land/lipgloss/v2"
)

type MsgType string

const (
	Add    MsgType = "add"
	Toggle MsgType = "toggle"
	Remove MsgType = "remove"
	Error  MsgType = "error"
)

var addStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#0A0A0A")).
	Background(lipgloss.Color("#13E68C")).
	PaddingRight(1).PaddingLeft(1)

var toggleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#0A0A0A")).
	Background(lipgloss.Color("#F2E86D")).
	PaddingRight(1).PaddingLeft(1)

var removeStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#D75F5F")).
	PaddingRight(1).PaddingLeft(1)

var removeSelectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Strikethrough(true)

var errorStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFFFFF")).
	Background(lipgloss.Color("#FF5F87")).
	PaddingRight(1).PaddingLeft(1)

var doneStyle = lipgloss.NewStyle().
	Strikethrough(true).
	Foreground(lipgloss.Color("#A0A0A0"))

var listTasksHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).PaddingLeft(1).PaddingRight(1)

func Log(msg string, msgType MsgType) {
	switch msgType {
	case Add:
		lipgloss.Printf("%s %s\n", addStyle.Render("Add"), msg)
	case Toggle:
		lipgloss.Printf("%s %s\n", toggleStyle.Render("Toggle"), msg)
	case Remove:
		lipgloss.Printf("%s %s\n", removeStyle.Render("Remove"), msg)
	case Error:
		lipgloss.Fprintf(os.Stderr, "%s %s\n", errorStyle.Render("Error"), msg)
	}
}
