package main

import (
	"encoding/json"
	"os"
)

func FromJson(typ interface{}, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(typ)
	return err
}
func FromJsonFile(typ interface{}, file *os.File) error {
	err := json.NewDecoder(file).Decode(typ)
	return err
}
func FromJsonBytes(typ interface{}, bytes []byte) error {
	return json.Unmarshal(bytes, typ)
}

func ToJson(typ interface{}, file_name string) error {
	b, err := json.MarshalIndent(typ, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file_name, b, 0644)
}
