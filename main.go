package main

import (
	"strings"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Add struct {
		Task []string `arg:"" name:"task" help:"Task description"`
	} `cmd:"" help:"Add a new task"`

	List struct {
		Filter string `arg:"" default:"todo" enum:"all,todo,done" help:"Filter listing"`
	} `cmd:"" help:"List all tasks"`

	Toggle struct {
		Ids []int `arg:"" optional:"" name:"ids" help:"IDs of tasks to toggle"`
	} `cmd:"" help:"Toggle the status of tasks"`

	Rm struct {
		Ids []int `arg:"" optional:"" name:"ids" help:"IDs of tasks to remove"`
	} `cmd:"" help:"Remove tasks"`
}

func main() {
	ctx := kong.Parse(&CLI)
	loadTasks()

	switch ctx.Command() {
	case "add <task>":
		task := strings.Join(CLI.Add.Task, " ")
		addTask(task)
	case "list", "list <filter>":
		switch strings.ToLower(CLI.List.Filter) {
		case string(All):
			listTasks(All)
		case string(Done):
			listTasks(Done)
		case string(Todo):
			listTasks(Todo)
		}
	case "toggle", "toggle <ids>":
		toggleTasks(CLI.Toggle.Ids...)
	case "rm", "rm <ids>":
		removeTasks(CLI.Rm.Ids...)
	}

	saveTasks()
}
