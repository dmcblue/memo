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
	APP_NAME          = "memo"
	CMD_ADD           = "add"
	CMD_EDIT          = "edit"
	CMD_TAG           = "tag"
	CMD_TAGS          = "tags"
	CMD_LIST          = "ls"
	CMD_REMOVE        = "rm"
	CMD_SEARCH        = "search"
	CMD_SHOW          = "show"
	CMD_VERSION       = "version"
	CMD_VERSION_LONG  = "--version"
	CMD_VERSION_SHORT = "-v"
	HELP              = "--help"
	HELP_SHORT        = "-h"
	VERSION           = "1.0.0"
)

type HelpCommand struct {
	Text    string
	SubText string
}

func help() {
	message := []HelpCommand{
		{
			Text:    fmt.Sprintf("%s <COMMAND>", APP_NAME),
			SubText: "",
		},
		{
			Text:    fmt.Sprintf("%s %s <TITLE> (<CONTENTS>)", APP_NAME, CMD_ADD),
			SubText: "Creates a new memo. If no CONTENTS is given, the system text editor will be opened for input.",
		},
		{
			Text:    fmt.Sprintf("%s %s (-a/--accept) <IDENTIFIER> (<CONTENTS>)", APP_NAME, CMD_EDIT),
			SubText: "Edits a memo. IDENTIFIER is either the memo title or the memo hash. If no CONTENTS is given, the system text editor will be opened for input. If the (-a/--accept) flag is provided, changes are auto-accepted. Otherwise, a diff will be presented for confirmation.",
		},
		{
			Text:    fmt.Sprintf("%s %s (-n/--no-format) (...-t/--tag <TAG>)", APP_NAME, CMD_LIST),
			SubText: "Prints memos. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated. Multiple (-t/--tag) options can be used to limit the results by memos with ANY of the listed tags",
		},
		{
			Text:    fmt.Sprintf("%s %s <IDENTIFIER>", APP_NAME, CMD_REMOVE),
			SubText: "Deletes a memo. IDENTIFIER is either the memo title or the memo hash.",
		},
		{
			Text:    fmt.Sprintf("%s %s  (-t/--title OR -c/--content) (-n/--no-format) <SEARCH_TERM>", APP_NAME, CMD_SEARCH),
			SubText: "Searches memos. The (-t/--title) limits the search to memo titles. The (-c/--content) limits the search to memo contents. You can only limit the search with one flag at a time. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated.",
		},
		{
			Text:    fmt.Sprintf("%s %s (-n/--no-format) <IDENTIFIER>", APP_NAME, CMD_SHOW),
			SubText: "Prints a memo. IDENTIFIER is either the memo title or the memo hash. The (-n/--no-format) flag prints each memo as a single-line with its values tab-separated.",
		},
		{
			Text:    fmt.Sprintf("%s %s %s <IDENTIFIER> <TAG>", APP_NAME, CMD_TAG, CMD_ADD),
			SubText: "Adds a tag to a memo. IDENTIFIER is either the memo title or the memo hash.",
		},
		{
			Text:    fmt.Sprintf("%s %s %s", APP_NAME, CMD_TAG, CMD_LIST),
			SubText: "Lists all existing tags",
		},
		{
			Text:    fmt.Sprintf("%s %s", APP_NAME, CMD_TAGS),
			SubText: fmt.Sprintf("Alias for `%s %s %s`", APP_NAME, CMD_TAG, CMD_LIST),
		},
		{
			Text:    fmt.Sprintf("%s %s %s <IDENTIFIER> <TAG>", APP_NAME, CMD_TAG, CMD_REMOVE),
			SubText: "Removes a tag to a memo. IDENTIFIER is either the memo title or the memo hash.",
		},
		{
			Text:    fmt.Sprintf("%s (%s/%s/%s)", APP_NAME, CMD_VERSION, CMD_VERSION_LONG, CMD_VERSION_SHORT),
			SubText: "Prints this current version.",
		},
		{
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

func dataError(msg string, optional_error_status ...int) {
	error_status := 1
	if len(optional_error_status) > 0 && optional_error_status[0] != 0 {
		error_status = optional_error_status[0]
	}
	fmt.Println(msg)
	os.Exit(error_status)
}

func PrintVersion() {
	fmt.Println(VERSION)
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
	case CMD_TAG:
		if len(os.Args) < 3 {
			cliError("No arguments given")
		}
		tagCommand := strings.TrimSpace(os.Args[2])
		switch tagCommand {
		case CMD_ADD:
			AddTag(config)
		case CMD_LIST:
			ShowTags(config)
		case CMD_REMOVE:
			RemoveTag(config)
		default:
			cliError(fmt.Sprintf("Unknown argument '%s'", tagCommand))
		}
	case CMD_TAGS:
		ShowTags(config)
	case CMD_SEARCH:
		SearchMemos(ui, config)
	case CMD_LIST:
		ShowMemos(ui, config)
	case CMD_REMOVE:
		RemoveMemo(ui, config)
	case CMD_SHOW:
		ShowMemo(ui, config)
	case CMD_VERSION:
		PrintVersion()
	case CMD_VERSION_LONG:
		PrintVersion()
	case CMD_VERSION_SHORT:
		PrintVersion()
	default:
		cliError(fmt.Sprintf("Unknown argument '%s'", command))
	}
}
