package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

/*********
 * Memos *
 *********/

func AddMemo(ui *Ui) {
	if len(os.Args) < 3 {
		cliError("No memo title given")
	}
	title := strings.TrimSpace(os.Args[2])

	memos := LoadMemos(saves_dir)
	for _, memo := range memos {
		if memo.Title == title {
			response := ui.GetResponse(
				fmt.Sprintf("Memo '%s' already exists.\nEdit? (y/n) ", title),
				"Invalid response. Try again: ",
				[]string{"y", "n"},
			)
			if response == "y" {
				EditMemo(ui) // inefficient but simple
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
	err := memo.Save(saves_dir)
	if err != nil {
		fmt.Printf("ERRR %v\n", err)
	}
}

func EditMemo(ui *Ui) {
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

	memos := LoadMemos(saves_dir)
	var memo_to_edit *Memo = nil
	for hash, memo := range memos {
		if hash[0:8] == identifier || memo.Title == identifier {
			memo_to_edit = memo
			break
		}
	}

	if memo_to_edit == nil {
		cliError(fmt.Sprintf("Unknown memo identifier '%s'\n", identifier))
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
	memo_to_edit.Save(saves_dir)
}

func ShowMemo(ui *Ui) {
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

	memos := LoadMemos(saves_dir)
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
		cliError(fmt.Sprintf("Unknown memo identifier '%s'\n", identifier))
	}
	memos_to_print := make(map[string]*Memo)
	memos_to_print[hash_to_print] = memo_to_print

	ui.PrintMemos(memos_to_print, skip_formatting)
}

func ShowMemos(ui *Ui) {
	skip_formatting := false
	search_labels_map := make(map[string]bool)
	for i := 2; i < len(os.Args); i++ {
		arg := strings.TrimSpace(os.Args[i])
		var label string
		if arg == "-n" || arg == "--no-format" {
			skip_formatting = true
		} else if arg == "-l" || arg == "--label" {
			if len(os.Args) < i+1 {
				cliError("Invalid label search a")
			}
			label = strings.TrimSpace(os.Args[i+1])
			if label != "-l" && label != "--label" {
				search_labels_map[label] = true
				i++
			} else {
				cliError("Invalid label search")
			}
		}
	}
	search_labels := []string{}
	for s := range search_labels_map {
		search_labels = append(search_labels, s)
	}
	memos := LoadMemos(saves_dir)
	memos_to_print := make(map[string]*Memo)
	for hash, memo := range memos {
		if len(search_labels) == 0 || AnyIntersection(search_labels, memo.Labels) {
			memos_to_print[hash] = memo
		}
	}

	ui.PrintMemos(memos_to_print, skip_formatting)
}

/**********
 * Labels *
 **********/

func AddLabel() {
	if len(os.Args) < 4 {
		cliError("No memo hash")
	}
	memo_hash := strings.TrimSpace(os.Args[3])
	memo := LoadMemoByHash(saves_dir, memo_hash)
	if memo == nil {
		cliError(fmt.Sprintf("No comment identifier '%s'", memo_hash))
	}
	if len(os.Args) < 5 {
		cliError("No label")
	}
	label := strings.TrimSpace(os.Args[4])
	memo.Labels = append(memo.Labels, label)
	memo.Save(saves_dir)
}

func RemoveLabel() {
	if len(os.Args) < 4 {
		cliError("No memo hash")
	}
	memo_hash := strings.TrimSpace(os.Args[3])
	memo := LoadMemoByHash(saves_dir, memo_hash)
	if memo == nil {
		cliError(fmt.Sprintf("No comment identifier '%s'", memo_hash))
	}
	if len(os.Args) < 5 {
		cliError("No label")
	}
	label := strings.TrimSpace(os.Args[4])
	var i int
	var found_label string
	for i, found_label = range memo.Labels {
		if found_label == label {
			break
		}
	}

	memo.Labels = append(memo.Labels[:i], memo.Labels[i+1:]...)
	memo.Save(saves_dir)
}

func ShowLabels() {
	labels := make(map[string]bool)
	memos := LoadMemos(saves_dir)
	for _, memo := range memos {
		for _, label := range memo.Labels {
			labels[label] = true
		}
	}
	for label, _ := range labels {
		fmt.Println(label)
	}
}
