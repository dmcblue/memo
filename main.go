package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func help() {
	fmt.Println("Help")
}

func cliError(msg string, optional_error_status ...int) {
	error_status := 1
	if len(optional_error_status) > 0 && optional_error_status[0] != 0 {
		error_status = optional_error_status[0]
	}
	fmt.Println(msg)
	fmt.Println("")
	help()
	os.Exit(error_status)
}

var saves_dir string
var ui *Ui

func init() {
	saves_dir = strings.TrimSpace(os.Getenv("MEMO_SAVES_DIR"))
	if strings.TrimSpace(saves_dir) == "" {
		saves_dir = SAVES_DIR
	}
	if err := os.MkdirAll(saves_dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	ui = CreateUi()
}

func main() {
	if len(os.Args) < 2 {
		cliError("No arguments given")
	}
	command := strings.TrimSpace(os.Args[1])
	switch command {
	case "add":
		AddMemo(ui)
	case "edit":
		EditMemo(ui)
	case "label":
		if len(os.Args) < 3 {
			cliError("No arguments given")
		}
		labelCommand := strings.TrimSpace(os.Args[2])
		switch labelCommand {
		case "add":
			AddLabel()
		case "ls":
			ShowLabels()
		case "rm":
			RemoveLabel()
		default:
			cliError(fmt.Sprintf("Unknown argument '%s'", command))
		}
	case "labels":
		ShowLabels()
	case "ls":
		ShowMemos(ui)
	default:
		cliError(fmt.Sprintf("Unknown argument '%s'", command))
	}
}
