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

func (t Task) Render() string {
	if t.Done {
		return doneStyle.Render(t.Desc)
	}
	return t.Desc
}

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

func toggleTasks(ids ...int) {
	// Toggle UI to select whichever you like to toggle if no param is supplied
	if len(ids) == 0 {
		if _, err := tea.NewProgram(initialToggleModel()).Run(); err != nil {
			Log(fmt.Sprintf("Failed to toggle tasks: %v", err), Error)
			os.Exit(1)
		}
		return
	}

	for _, id := range ids {
		if id >= 1 && id <= len(tasks) {
			tasks[id-1].Done = !tasks[id-1].Done
			Log(tasks[id-1].Render(), Toggle)
		} else {
			Log(fmt.Sprintf("Invalid task ID: %d", id), Error)
		}
	}
}

func removeTasks(ids ...int) {
	if len(ids) == 0 {
		if _, err := tea.NewProgram(initialRemoveModel()).Run(); err != nil {
			Log(fmt.Sprintf("Failed to remove tasks: %v", err), Error)
			os.Exit(1)
		}
		return
	}

	removeTasksByID(true, ids...)
}

func removeTasksByID(log bool, ids ...int) {
	idsToRemove := make(map[int]bool)
	for _, id := range ids {
		if id >= 1 && id <= len(tasks) {
			idsToRemove[id-1] = true
		} else {
			Log(fmt.Sprintf("Invalid task ID: %d", id), Error)
		}
	}

	if len(idsToRemove) == 0 {
		return
	}

	var remaining []Task
	for i, task := range tasks {
		if idsToRemove[i] {
			if log {
				Log(task.Render(), Remove)
			}
			continue
		}
		remaining = append(remaining, task)
	}
	tasks = remaining
}
