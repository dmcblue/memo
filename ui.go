package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"slices"
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

func (ui *Ui) PrintMemos(memos map[string]*Memo) {
	width := GetTermWidth()

	if width == 0 {
		for hash, memo := range memos {
			ui.PrintMemo(hash, memo)
		}
	} else {
		// sha 8 + 4 space + title + 4 space + content
		var max_title_length float64 = 0
		var max_label_length float64 = 0
		for _, memo := range memos {
			max_title_length = math.Max(max_title_length, float64(len(memo.Title)))
			for _, label := range memo.Labels {
				max_label_length = math.Max(max_label_length, float64(len(label)))
			}
		}
		title_length := (width - 20) / 3
		label_width := int(max_label_length)
		content_length := width - 20 - title_length - label_width
		for hash, memo := range memos {
			ui.PrintMemoFancy(
				hash,
				memo,
				title_length,
				content_length,
				content_length,
			)
		}
	}
}

func (ui *Ui) PrintMemoFancy(hash string, memo *Memo, title_length int, content_length int, label_length int) {
	// contents := Chunks(memo.Content, content_length)
	// titles := Chunks(memo.Title, title_length)
	// labels := Chunks(strings.Join(memo.Labels, ", "), label_length)

}

func (ui *Ui) PrintMemo(hash string, memo *Memo) {
	fmt.Printf("%s\n%s\n%s\n", hash[0:8], memo.Title, memo.Content)
	if len(memo.Labels) > 0 {
		fmt.Println(strings.Join(memo.Labels, ", "))
	}
}

/************
 * UI Utils *
 ************/

func GetTermWidth() int {
	if term.IsTerminal(0) {
		println("in a term")
		width, _, err := term.GetSize(0)
		if err != nil {
			return 0
		}

		return width
	} else {
		return 0
	}
}

// https://stackoverflow.com/a/61469854
func Chunks(str string, chunkSize int) []string {
	if len(str) == 0 {
		return nil
	}
	if chunkSize >= len(str) {
		return []string{str}
	}
	chunks := []string{}
	// var chunks []string = make([]string, 0, (len(str)-1)/chunkSize+1)
	// currentLen := 0
	// currentStart := 0
	strs := strings.Split(str, " ")
	currentChunk := ""
	for _, t := range strs {
		if len(currentChunk)+1+len(t) > chunkSize {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		} else {
			currentChunk += " " + t
		}
	}
	chunks = append(chunks, currentChunk)
	// for i := range str {
	// 	if currentLen == chunkSize {
	// 		chunks = append(chunks, str[currentStart:i])
	// 		currentLen = 0
	// 		currentStart = i
	// 	}
	// 	currentLen++
	// }
	// chunks = append(chunks, str[currentStart:])
	return chunks
}
