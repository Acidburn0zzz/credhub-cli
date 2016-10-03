package models

import (
	"strings"
)

type AllPathResponseBody struct {
	Paths []Path `json:"paths"`
}

type Path struct {
	Path string `json:"path,omitempty"`
}

func (allPathResponseBody AllPathResponseBody) String() string {
	lines := []string{}
	lines = append(lines, "Path")
	for _, path := range allPathResponseBody.Paths {
		lines = append(lines, path.Path)
	}
	return strings.Join(lines, "\n")
}
