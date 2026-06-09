package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type FilterDone string

const (
	All  FilterDone = "all"
	Done FilterDone = "done"
	Todo FilterDone = "todo"
)

type Task struct {
	Desc string `json:"desc"`
	Done bool   `json:"done"`
}

const tasksFilePath = "tasks.json"

var tasks []Task

func initTask(desc string) Task {
	task := Task{
		Desc: desc,
	}
	return task
}

func loadTasks() {
	if _, err := os.Stat(tasksFilePath); errors.Is(err, os.ErrNotExist) {
		saveTasks()
	}

	content, err := os.ReadFile(tasksFilePath)

	if err != nil {
		Log(fmt.Sprintf("Failed to load tasks: %v", err), Error)
		os.Exit(1)
	}

	err = json.Unmarshal(content, &tasks)
	if err != nil {
		Log(fmt.Sprintf("Failed to load tasks: %v", err), Error)
		os.Exit(1)
	}
}

func saveTasks() {
	fileData, err := json.Marshal(tasks)
	if err != nil {
		Log(fmt.Sprintf("Failed to save tasks: %v", err), Error)
		os.Exit(1)
	}
	err = os.WriteFile(tasksFilePath, fileData, 0644)
	if err != nil {
		Log(fmt.Sprintf("Failed to save tasks: %v", err), Error)
		os.Exit(1)
	}
}

func addTask(desc string) {
	task := initTask(desc)
	tasks = append(tasks, task)
	Log(desc, Add)
}

func listTasks(filter FilterDone) {
	lipgloss.Printf("%s\n", listTasksHeaderStyle.Render("Tasks"))
	for i, task := range tasks {
		if task.Done && (filter == All || filter == Done) {
			lipgloss.Printf("%d. %s\n", i+1, doneStyle.Render(task.Desc))
		}
		if !task.Done && (filter == All || filter == Todo) {
			lipgloss.Printf("%d. %s\n", i+1, task.Desc)
		}
	}
}

func toggleTask() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
