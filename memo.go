package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Memo struct {
	Title   string
	Content string
	Labels  []string
}

const (
	SAVES_DIR = "saves"
)

type HASH = string

func CreateMemo(title string, content string) *Memo {
	return &Memo{
		Title:   title,
		Content: content,
		Labels:  []string{},
	}
}

func (memo *Memo) Save(saves_dir string) error {
	filename := ToFilename(memo.Title, "")
	fullpath := filepath.Join(
		saves_dir,
		filename,
	)

	return ToJson(memo, fullpath)
}

func LoadMemo(filename string, memo *Memo, saves_dir string) {
	filePath := filepath.Join(
		saves_dir,
		filename,
	)
	err := FromJson(memo, filePath)
	if err != nil {
		fmt.Printf("Err: %v", err)
	}
}

func LoadMemoByHash(saves_dir, hash string) *Memo {
	memos := LoadMemos(saves_dir)
	for full_hash, memo := range memos {
		if full_hash[0:8] == hash {
			return memo
		}
	}
	return nil
}

func LoadMemos(saves_dir string) map[HASH]*Memo {
	files, err := os.ReadDir(saves_dir)
	if err != nil {
		log.Fatal(err)
	}

	memos := make(map[HASH]*Memo)
	for _, fileEntry := range files {
		if !fileEntry.IsDir() {
			var memo *Memo = CreateMemo("", "")
			LoadMemo(fileEntry.Name(), memo, saves_dir)
			hash := sha1.Sum([]byte(fileEntry.Name()))
			hash_str := fmt.Sprintf("%x", hash)
			memos[hash_str] = memo
		}
	}

	return memos
}
