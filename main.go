package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	APP_NAME   = "memo"
	CMD_ADD    = "add"
	CMD_EDIT   = "edit"
	CMD_LABEL  = "label"
	CMD_LABELS = "labels"
	CMD_LIST   = "ls"
	CMD_REMOVE = "rm"
	CMD_SHOW   = "show"
	HELP       = "--help"
	HELP_SHORT = "-h"
)

func help() {
	fmt.Println(fmt.Sprintf("%s <COMMAND>", APP_NAME))
	fmt.Println(fmt.Sprintf("%s %s <TITLE> (<CONTENTS>)", APP_NAME, CMD_ADD))
	fmt.Println("\tCreates a new memo. If no CONTENTS is given, the system text editor will be opened for input.")
	fmt.Println(fmt.Sprintf("%s %s (-a/--accept) <IDENTIFIER> (<CONTENTS>)", APP_NAME, CMD_EDIT))
	fmt.Println("\tEdits a memo. IDENTIFIER is either the memo title or the memo hash. If no CONTENTS is given, the system text editor will be opened for input. If the (-a/--accept) flag is provided, changes are auto-accepted. Otherwise, a diff will be presented for confirmation.")
	fmt.Println(fmt.Sprintf("%s %s (-n/--no-format) (...-l/--label <LABEL>)", APP_NAME, CMD_LIST))
	fmt.Println("\tPrints memos. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated. Multiple (-l/--label) options can be used to limit the results by memos with ANY of the listed labels")
	fmt.Println(fmt.Sprintf("%s %s (-n/--no-format) <IDENTIFIER>", APP_NAME, CMD_SHOW))
	fmt.Println("\tPrints a memo. IDENTIFIER is either the memo title or the memo hash. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated.")
	fmt.Println(fmt.Sprintf("%s %s %s <IDENTIFIER> <LABEL>", APP_NAME, CMD_LABEL, CMD_ADD))
	fmt.Println("\tAdds a label/tag to a memo. IDENTIFIER is either the memo title or the memo hash.")
	fmt.Println(fmt.Sprintf("%s %s %s", APP_NAME, CMD_LABEL, CMD_LIST))
	fmt.Println("\tLists all existing labels/tags")
	fmt.Println(fmt.Sprintf("%s %s", APP_NAME, CMD_LABELS))
	fmt.Println(fmt.Sprintf("\tAlias for `%s %s %s`", APP_NAME, CMD_LABEL, CMD_LIST))
	fmt.Println(fmt.Sprintf("%s %s %s <IDENTIFIER> <LABEL>", APP_NAME, CMD_LABEL, CMD_REMOVE))
	fmt.Println("\tRemoves a label/tag to a memo. IDENTIFIER is either the memo title or the memo hash.")
	fmt.Println(fmt.Sprintf("%s (%s/%s)", APP_NAME, HELP, HELP_SHORT))
	fmt.Println("\tPrints this message.")
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
	case HELP:
		help()
	case HELP_SHORT:
		help()
	case CMD_ADD:
		AddMemo(ui)
	case CMD_EDIT:
		EditMemo(ui)
	case CMD_LABEL:
		if len(os.Args) < 3 {
			cliError("No arguments given")
		}
		labelCommand := strings.TrimSpace(os.Args[2])
		switch labelCommand {
		case CMD_ADD:
			AddLabel()
		case CMD_LIST:
			ShowLabels()
		case CMD_REMOVE:
			RemoveLabel()
		default:
			cliError(fmt.Sprintf("Unknown argument '%s'", command))
		}
	case CMD_LABELS:
		ShowLabels()
	case CMD_LIST:
		ShowMemos(ui)
	case CMD_SHOW:
		ShowMemo(ui)
	default:
		cliError(fmt.Sprintf("Unknown argument '%s'", command))
	}
}
