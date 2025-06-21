package main

import (
	"strings"

	"github.com/flytam/filenamify"
)

func AnyIntersection(as []string, bs []string) bool {
	intersection := make(map[string]int)
	for _, a := range as {
		intersection[a] = 1
	}
	for _, b := range bs {
		if _, ok := intersection[b]; ok {
			return true
		}
	}
	return false
}

func StringInSlice(str string, sl []string) bool {
	for _, s := range sl {
		if str == s {
			return true
		}
	}
	return false
}

func ToFilename(str string, ext string) string {
	filename, _ := filenamify.Filenamify(
		str,
		filenamify.Options{
			Replacement: "_",
		},
	)
	filename = strings.ReplaceAll(filename, " ", "_")
	filename = strings.ToLower(filename)
	if ext != "" {
		filename += "." + ext
	}
	return filename
}
