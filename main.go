package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	listCommand := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Subcommand not supplied")
		os.Exit(1)
	}

	loadTasks()
	switch os.Args[1] {
	case "add":
		addCommand.Parse(os.Args[2:])
		task := strings.Join(addCommand.Args(), " ")
		addTask(task)
	case "list":
		listCommand.Parse(os.Args[2:])
		if strings.ToLower(listCommand.Arg(0)) == "all" {
			listTasks(All)
		} else if strings.ToLower(listCommand.Arg(0)) == "done" {
			listTasks(Done)
		} else {
			listTasks(Todo)
		}
	case "toggle":
		toggleTask()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}

	saveTasks()
}
