package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

/*********
 * Memos *
 *********/

func AddMemo(ui *Ui, config *Config) {
	if len(os.Args) < 3 {
		cliError("No memo title given")
	}
	title := strings.TrimSpace(os.Args[2])

	memos := LoadMemos(config.SavesDir)
	for _, memo := range memos {
		if memo.Title == title {
			response := ui.GetResponse(
				fmt.Sprintf("Memo '%s' already exists.\nEdit? (y/n) ", title),
				"Invalid response. Try again: ",
				[]string{"y", "n"},
			)
			if response == "y" {
				EditMemo(ui, config) // inefficient but simple
			}

			return
		}
	}

	var content string
	if len(os.Args) < 4 {
		content = ui.EditContent("")
	} else {
		content = strings.TrimSpace(os.Args[3])
	}
	memo := CreateMemo(title, content)
	hash := memo.Save(config.SavesDir)
	fmt.Println(hash[0:8])
}

func EditMemo(ui *Ui, config *Config) {
	identifier := ""
	new_content := ""
	auto_accept := false
	for i := 2; i < len(os.Args); i++ {
		arg := strings.TrimSpace(os.Args[i])
		if arg == "-a" || arg == "--accept" {
			auto_accept = true
		} else if identifier == "" {
			identifier = arg
		} else {
			new_content = arg
		}
	}

	if identifier == "" {
		cliError("No memo hash/title given")
	}

	memos := LoadMemos(config.SavesDir)
	var memo_to_edit *Memo = nil
	for hash, memo := range memos {
		if hash[0:8] == identifier || memo.Title == identifier {
			memo_to_edit = memo
			break
		}
	}

	if memo_to_edit == nil {
		dataError(fmt.Sprintf("Unknown memo identifier '%s'\n", identifier))
	}

	if new_content == "" {
		new_content = ui.EditContent(memo_to_edit.Content)
	}

	if !auto_accept {
		// Unmaintained package
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(memo_to_edit.Content),
			B:        difflib.SplitLines(new_content),
			FromFile: "Original",
			ToFile:   "New",
			Context:  3,
		}
		text, _ := difflib.GetUnifiedDiffString(diff)
		fmt.Println(text)
		fmt.Println("Changes:")
		response := ui.GetResponse(
			"Accept changes? (y/n) ",
			"Try again: ",
			[]string{"y", "n"},
		)
		if response == "n" {
			fmt.Println("Changes scrapped")
			return
		}
	}

	memo_to_edit.Content = new_content
	memo_to_edit.Save(config.SavesDir)
}

func RemoveMemo(ui *Ui, config *Config) {
	if len(os.Args) < 3 {
		cliError("No memo identifier given")
	}
	identifier := strings.TrimSpace(os.Args[2])

	if identifier == "" {
		cliError("No memo hash/title given")
	}

	memos := LoadMemos(config.SavesDir)
	var memo_to_remove *Memo = nil
	for hash, memo := range memos {
		if hash[0:8] == identifier || memo.Title == identifier {
			memo_to_remove = memo
			break
		}
	}

	if memo_to_remove == nil {
		dataError(fmt.Sprintf("Unknown memo identifier '%s'\n", identifier))
	}

	memo_to_remove.Delete(config.SavesDir)
}

func ShowMemo(ui *Ui, config *Config) {
	skip_formatting := false
	identifier := ""
	for i := 2; i < len(os.Args); i++ {
		arg := strings.TrimSpace(os.Args[i])
		if arg == "-n" || arg == "--no-format" {
			skip_formatting = true
		} else {
			identifier = arg
		}
	}

	if identifier == "" {
		cliError("No memo hash/title given")
	}

	memos := LoadMemos(config.SavesDir)
	var memo_to_print *Memo = nil
	var hash_to_print HASH = ""
	for hash, memo := range memos {
		if hash[0:8] == identifier || memo.Title == identifier {
			hash_to_print = hash
			memo_to_print = memo
			break
		}
	}

	if memo_to_print == nil {
		dataError(fmt.Sprintf("Unknown memo identifier '%s'\n", identifier))
	}
	memos_to_print := make(map[string]*Memo)
	memos_to_print[hash_to_print] = memo_to_print

	ui.PrintMemos(memos_to_print, skip_formatting)
}

func ShowMemos(ui *Ui, config *Config) {
	skip_formatting := false
	search_tags_map := make(map[string]bool)
	for i := 2; i < len(os.Args); i++ {
		arg := strings.TrimSpace(os.Args[i])
		var tag string
		if arg == "-n" || arg == "--no-format" {
			skip_formatting = true
		} else if arg == "-t" || arg == "--tag" {
			if len(os.Args) < i+1 {
				cliError("Invalid tag search a")
			}
			tag = strings.TrimSpace(os.Args[i+1])
			if tag != "-t" && tag != "--tag" {
				search_tags_map[tag] = true
				i++
			} else {
				cliError("Invalid tag search")
			}
		}
	}
	search_tags := []string{}
	for s := range search_tags_map {
		search_tags = append(search_tags, s)
	}
	memos := LoadMemos(config.SavesDir)
	memos_to_print := make(map[string]*Memo)
	for hash, memo := range memos {
		if len(search_tags) == 0 || AnyIntersection(search_tags, memo.Tags) {
			memos_to_print[hash] = memo
		}
	}

	ui.PrintMemos(memos_to_print, skip_formatting)
}

/********
 * Tags *
 ********/

func AddTag(config *Config) {
	if len(os.Args) < 4 {
		cliError("No memo hash")
	}
	memo_hash := strings.TrimSpace(os.Args[3])
	memo := LoadMemoByHash(config.SavesDir, memo_hash)
	if memo == nil {
		dataError(fmt.Sprintf("No memo identifier '%s'", memo_hash))
	}
	if len(os.Args) < 5 {
		cliError("No tag")
	}
	tag := strings.TrimSpace(os.Args[4])
	memo.Tags = append(memo.Tags, tag)
	memo.Save(config.SavesDir)
}

func RemoveTag(config *Config) {
	if len(os.Args) < 4 {
		cliError("No memo hash")
	}
	memo_hash := strings.TrimSpace(os.Args[3])
	memo := LoadMemoByHash(config.SavesDir, memo_hash)
	if memo == nil {
		dataError(fmt.Sprintf("No memo identifier '%s'", memo_hash))
	}
	if len(os.Args) < 5 {
		cliError("No tag")
	}
	tag := strings.TrimSpace(os.Args[4])
	var i int
	var found_tag string
	for i, found_tag = range memo.Tags {
		if found_tag == tag {
			break
		}
	}

	memo.Tags = append(memo.Tags[:i], memo.Tags[i+1:]...)
	memo.Save(config.SavesDir)
}

func ShowTags(config *Config) {
	tags := make(map[string]bool)
	memos := LoadMemos(config.SavesDir)
	for _, memo := range memos {
		for _, tag := range memo.Tags {
			tags[tag] = true
		}
	}

	all_tags := make([]string, 0)
	for tag := range tags {
		all_tags = append(all_tags, tag)
	}

	sort.Strings(all_tags)
	for _, tag := range all_tags {
		fmt.Println(tag)
	}
}
