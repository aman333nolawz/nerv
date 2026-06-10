<p align="center">
  <img src="https://i.pinimg.com/736x/ed/41/b9/ed41b9da2d4f50a0f5d99a871d84fdb3.jpg" alt="Nerv banner" width="300">
</p>

# Nerv

Nerv is a small terminal todo application written in Go.

There are already a lot of todo apps, but this project exists for a more personal reason: building a custom one makes it easy to change the workflow, behavior, and interface whenever needed. It is also a learning project for practicing Go and improving by building something useful from scratch.

## Features

- Add todo items from the command line.
- List todo items by status.
- Toggle tasks between todo and done.
- Remove tasks by ID.
- Use interactive terminal views when toggling or removing without IDs.
- Store tasks locally in `tasks.json`.

## Requirements

- Go `1.26.3` or newer, based on the version in `go.mod`.

## Installation

Clone the repository and build the binary:

```sh
go build -o nerv .
```

Run it from the project directory:

```sh
./nerv --help
```

You can also run it directly without building:

```sh
go run . --help
```

## Usage

```text
Usage: nerv <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  add <task> ...
    Add a new task

  list [<filter>]
    List all tasks

  toggle [<ids> ...]
    Toggle the status of tasks

  rm [<ids> ...]
    Remove tasks

Run "nerv <command> --help" for more information on a command.
```

## Commands

### `add`

Adds a new task.

```sh
nerv add Buy milk
```

The full text after `add` becomes the task description.

```text
Usage: nerv add <task> ...

Add a new task

Arguments:
  <task> ...    Task description
```

### `list`

Lists tasks. The optional filter can be `todo`, `done`, or `all`.

```sh
nerv list
nerv list todo
nerv list done
nerv list all
```

If no filter is provided, Nerv lists `todo` tasks by default.

```text
Usage: nerv list [<filter>]

List all tasks

Arguments:
  [<filter>]    Filter listing
```

### `toggle`

Toggles tasks between todo and done.

```sh
nerv toggle 1
nerv toggle 1 3 4
```

If you run `toggle` without IDs, Nerv opens an interactive terminal view where you can select tasks.

```sh
nerv toggle
```

Interactive controls:

- `space` or `enter`: toggle the selected task
- `q` or `ctrl+c`: quit

```text
Usage: nerv toggle [<ids> ...]

Toggle the status of tasks

Arguments:
  [<ids> ...]    IDs of tasks to toggle
```

### `rm`

Removes tasks by ID.

```sh
nerv rm 2
nerv rm 1 3 5
```

If you run `rm` without IDs, Nerv opens an interactive terminal view where you can mark tasks for removal.

```sh
nerv rm
```

Interactive controls:

- `space` or `enter`: mark or unmark the selected task for removal
- `q` or `ctrl+c`: quit and remove the marked tasks

```text
Usage: nerv rm [<ids> ...]

Remove tasks

Arguments:
  [<ids> ...]    IDs of tasks to remove
```

## Data

Tasks are saved in a local `tasks.json` file in the directory where Nerv is run. Each task stores:

- `desc`: the task description
- `done`: whether the task is completed

Example:

```json
[
  {
    "desc": "Buy milk",
    "done": false
  }
]
```

## Tech Stack

- [Go](https://go.dev/)
- [Kong](https://github.com/alecthomas/kong) for CLI parsing
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Bubbles](https://github.com/charmbracelet/bubbles) for interactive terminal views
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) for terminal styling
