package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/term"
)

type Ui struct {
	Scanner *bufio.Scanner
}

func CreateUi() *Ui {
	return &Ui{
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

/*********
 * Input *
 *********/

func (ui *Ui) GetResponse(
	prompt string,
	followUp string,
	acceptableResponses []string,
) string {
	fmt.Print(prompt)
	var text string
	for {
		ui.Scanner.Scan()
		text = strings.TrimSpace(ui.Scanner.Text())
		if !slices.Contains(acceptableResponses, text) {
			fmt.Print(followUp)
		} else {
			break
		}
	}
	return text
}

func (ui *Ui) GetText(
	prompt string,
	followUp string,
) string {
	fmt.Print(prompt)
	var text string
	for {
		ui.Scanner.Scan()
		text = strings.TrimSpace(ui.Scanner.Text())
		if len(text) == 0 && followUp != "" {
			fmt.Print(followUp)
		} else {
			break
		}
	}
	return text
}

func (ui *Ui) GetInt(
	prompt string,
	followUp string,
	min int,
	max int, // inclusive
) int {
	fmt.Print(prompt)
	var text string
	for {
		ui.Scanner.Scan()
		text = strings.TrimSpace(ui.Scanner.Text())
		i, err := strconv.Atoi(text)
		if err != nil || i < min || i > max {
			fmt.Print(followUp)
		} else {
			return i
		}
	}
}

func (ui *Ui) GetMultilineText(prompt string, doneText string) string {
	lineText := ""
	totalText := ""
	for {
		lineText = ui.GetText(
			prompt,
			"",
		)
		if lineText == doneText {
			break
		}
		totalText += lineText + "\n"
	}

	return totalText
}

/*********
 * Memos *
 *********/

// https://github.com/msemjan/go-external-editor/blob/9e2e6ee617d8dcb9a41a86e49282170327ba524d/main.go#L22C2-L66C3
func (ui *Ui) EditContent(fileContents string) string {
	// Create a temporary file
	f, err1 := os.CreateTemp("", "memo-*")
	if err1 != nil {
		log.Fatal(err1)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	// Fill the termporary file with some data
	f.WriteString(fileContents)

	// Open the temporary file in Vim for editing
	editor := strings.TrimSpace(os.Getenv("EDITOR"))
	cmd := exec.Command(editor, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Open the temporary file in an external editor
	err2 := cmd.Start()
	if err2 != nil {
		log.Printf("2")
		log.Fatal(err2)
	}

	// Wait for the user to do their thing...
	err2 = cmd.Wait()

	// Check if they didn't mess up everything
	if err2 != nil {
		log.Printf("Error while editing. Error: %v\n", err2)
	}

	content, err := os.ReadFile(f.Name())
	if err != nil {
		fmt.Println("Err")
	}

	return string(content)
}

func (ui *Ui) PrintMemos(memos map[string]*Memo, skip_formatting bool) {
	width := GetTermWidth()

	if width == 0 || skip_formatting {
		for hash, memo := range memos {
			ui.PrintMemo(hash, memo)
			fmt.Println()
		}
	} else {
		// Assuming tab width == 4 for now
		// sha 8 + 4 space + title + 4 space + content
		var max_title_length float64 = 0
		var max_content_length float64 = 0
		var max_label_length float64 = 0
		for _, memo := range memos {
			max_title_length = math.Max(max_title_length, float64(len(memo.Title)))
			max_content_length = math.Max(max_content_length, LongestOfMultiline(memo.Content))
			for _, label := range memo.Labels {
				max_label_length = math.Max(max_label_length, float64(len(label)))
			}
		}
		title_memo := &Memo{
			Title:   "TITLE",
			Content: "CONTENT",
			Labels:  []string{"LABELS"},
		}
		max_title_length = math.Max(max_title_length, float64(len(title_memo.Title)))
		max_content_length = math.Max(max_content_length, LongestOfMultiline(title_memo.Content))
		max_label_length = math.Max(max_label_length, float64(len(title_memo.Labels[0])))

		title_length := int(math.Min(max_title_length, float64((width-20)/3)))
		label_width := int(max_label_length)
		content_length := int(math.Min(max_content_length, float64(width-20-title_length-label_width-1))) // 1 for right side padding
		// If there's extra room, expand labels
		label_width = width - content_length - 20 - title_length - 1
		ui.PrintMemoFancy(
			"HASH    ",
			title_memo,
			title_length,
			content_length,
			label_width,
		)
		fmt.Println()

		hashes := make([]string, 0)
		for hash := range memos {
			hashes = append(hashes, hash)
		}

		// Now sort the slice
		sort.Strings(hashes)

		// Iterate over all keys in a sorted order
		for _, hash := range hashes {
			memo := memos[hash]
			ui.PrintMemoFancy(
				hash,
				memo,
				title_length,
				content_length,
				label_width,
			)
			fmt.Println()
		}
	}
}

func (ui *Ui) PrintMemoFancy(hash string, memo *Memo, title_length int, content_length int, label_length int) {
	contents := Chunks(memo.Content, content_length)
	titles := Chunks(memo.Title, title_length)
	labels := Chunks(strings.Join(memo.Labels, ", "), label_length)
	lines := int(math.Max(float64(len(contents)), math.Max(float64(len(titles)), float64(len(labels)))))
	for i := range lines {
		if i == 0 {
			fmt.Print(hash[0:8])
		} else {
			fmt.Print(strings.Repeat(" ", 8))
		}

		fmt.Print(strings.Repeat(" ", 4))

		if len(titles) > i {
			fmt.Printf("%-*s", title_length, titles[i])
		} else {
			fmt.Print(strings.Repeat(" ", title_length))
		}

		fmt.Print(strings.Repeat(" ", 4))

		if len(contents) > i {
			fmt.Printf("%-*s", content_length, contents[i])
		} else {
			fmt.Print(strings.Repeat(" ", content_length))
		}

		fmt.Print(strings.Repeat(" ", 4))

		if len(labels) > i {
			fmt.Printf("%-*s", label_length, labels[i])
		} else {
			fmt.Print(strings.Repeat(" ", label_length))
		}

		fmt.Println()
	}
}

func (ui *Ui) PrintMemo(hash string, memo *Memo) {
	fmt.Printf("%s\t%s\t%s", hash[0:8], memo.Title, strings.ReplaceAll(memo.Content, "\n", "\\n"))
	fmt.Printf("\t%s", strings.Join(memo.Labels, ", "))
}

/************
 * UI Utils *
 ************/

// https://stackoverflow.com/a/61469854
func Chunks(str string, chunkSize int) []string {
	if len(str) == 0 {
		return nil
	}

	chunks := []string{}
	currentChunk := ""

	lines := strings.Split(str, "\n")
	for _, line := range lines {

		strs := strings.Split(line, " ")
		for _, t := range strs {

			if currentChunk == "" {
				currentChunk = t
			} else if len(currentChunk)+1+len(t) > chunkSize {
				chunks = append(chunks, currentChunk)
				currentChunk = t
			} else {
				currentChunk += " " + t
			}
		}
		chunks = append(chunks, currentChunk)
		currentChunk = ""
	}

	return chunks
}

func GetTermWidth() int {
	if term.IsTerminal(0) {
		width, _, err := term.GetSize(0)
		if err != nil {
			return 0
		}

		return width
	} else {
		return 0
	}
}

func LongestOfMultiline(str string) float64 {
	lines := strings.Split(str, "\n")
	var max_len float64 = 0
	for _, line := range lines {
		max_len = math.Max(max_len, float64(len(line)))
	}

	return max_len
}
