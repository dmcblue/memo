package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
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

type HelpCommand struct {
	Text    string
	SubText string
}

func help() {
	message := []HelpCommand{
		HelpCommand{
			Text:    fmt.Sprintf("%s <COMMAND>", APP_NAME),
			SubText: "",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s <TITLE> (<CONTENTS>)", APP_NAME, CMD_ADD),
			SubText: "Creates a new memo. If no CONTENTS is given, the system text editor will be opened for input.",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s (-a/--accept) <IDENTIFIER> (<CONTENTS>)", APP_NAME, CMD_EDIT),
			SubText: "Edits a memo. IDENTIFIER is either the memo title or the memo hash. If no CONTENTS is given, the system text editor will be opened for input. If the (-a/--accept) flag is provided, changes are auto-accepted. Otherwise, a diff will be presented for confirmation.",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s (-n/--no-format) (...-l/--label <LABEL>)", APP_NAME, CMD_LIST),
			SubText: "Prints memos. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated. Multiple (-l/--label) options can be used to limit the results by memos with ANY of the listed labels",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s (-n/--no-format) <IDENTIFIER>", APP_NAME, CMD_SHOW),
			SubText: "Prints a memo. IDENTIFIER is either the memo title or the memo hash. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated.",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s %s <IDENTIFIER> <LABEL>", APP_NAME, CMD_LABEL, CMD_ADD),
			SubText: "Adds a label/tag to a memo. IDENTIFIER is either the memo title or the memo hash.",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s %s", APP_NAME, CMD_LABEL, CMD_LIST),
			SubText: "Lists all existing labels/tags",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s", APP_NAME, CMD_LABELS),
			SubText: fmt.Sprintf("Alias for `%s %s %s`", APP_NAME, CMD_LABEL, CMD_LIST),
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s %s %s <IDENTIFIER> <LABEL>", APP_NAME, CMD_LABEL, CMD_REMOVE),
			SubText: "Removes a label/tag to a memo. IDENTIFIER is either the memo title or the memo hash.",
		},
		HelpCommand{
			Text:    fmt.Sprintf("%s (%s/%s)", APP_NAME, HELP, HELP_SHORT),
			SubText: "Prints this message.",
		},
	}

	width := GetTermWidth()
	if width == 0 {
		for _, line := range message {
			fmt.Println(line.Text)
			fmt.Println("    " + line.SubText)
		}
	} else {
		for _, line := range message {
			textChunks := Chunks(line.Text, width-1)
			subTextChunks := Chunks(line.SubText, width-5)
			for _, t := range textChunks {
				fmt.Println(t)
			}
			for _, s := range subTextChunks {
				fmt.Println("    " + s)
			}
		}
	}
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

type Config struct {
	SavesDir string
}

var config *Config
var ui *Ui

func LoadConfig() {
	config_path := strings.TrimSpace(os.Getenv("MEMO_CONF_PATH"))
	user_dir, _ := os.UserHomeDir()
	user_config_dir, _ := os.UserConfigDir()
	config_dir := strings.Replace(user_config_dir, "~", user_dir, 1)
	if config_path == "" {
		config_path = path.Join(config_dir, "memo.conf")
	}

	_, err := os.Stat(config_path)
	if errors.Is(err, os.ErrNotExist) {
		default_config := &Config{
			SavesDir: path.Join(config_dir, "memo", "saves"),
		}
		if err := ToJson(default_config, config_path); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		cliError(
			"Unknown config access error\n" +
				fmt.Sprintf("\tLoaded Config Path: '%s'\n", config_path) +
				fmt.Sprintf("\tSystem Config Dir: '%s'\n", config_dir) +
				fmt.Sprintf("\tMEMO_CONF_PATH: '%s'\n", os.Getenv("MEMO_CONF_PATH")) +
				"Consider adjusting/unsetting MEMO_CONF_PATH\n\n" +
				fmt.Sprintf("%v\n", config_path),
		)
	}

	if err = FromJson(config, config_path); err != nil {
		log.Fatal(err)
	}
}

func init() {
	config = &Config{
		SavesDir: "",
	}
	LoadConfig()

	if err := os.MkdirAll(config.SavesDir, os.ModePerm); err != nil {
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
		AddMemo(ui, config)
	case CMD_EDIT:
		EditMemo(ui, config)
	case CMD_LABEL:
		if len(os.Args) < 3 {
			cliError("No arguments given")
		}
		labelCommand := strings.TrimSpace(os.Args[2])
		switch labelCommand {
		case CMD_ADD:
			AddLabel(config)
		case CMD_LIST:
			ShowLabels(config)
		case CMD_REMOVE:
			RemoveLabel(config)
		default:
			cliError(fmt.Sprintf("Unknown argument '%s'", command))
		}
	case CMD_LABELS:
		ShowLabels(config)
	case CMD_LIST:
		ShowMemos(ui, config)
	case CMD_SHOW:
		ShowMemo(ui, config)
	default:
		cliError(fmt.Sprintf("Unknown argument '%s'", command))
	}
}
