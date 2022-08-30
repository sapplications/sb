// Copyright 2022 Vitalii Noha vitalii.noga@gmail.com. All rights reserved.

package app

import (
	"fmt"
)

const (
	coderAttrName string = "coder"
)

var modKind = struct {
	sa string
	sb string
	sp string
}{
	"sa",
	"sb",
	"sp",
}

type builder interface {
	Build(app string, sources *map[string]map[string]string) error
	Clean(app string, sources *map[string]map[string]string) error
	Generate(app string, sources *map[string]map[string]string) error
}

func handleError() {
	if r := recover(); r != nil {
		fmt.Printf(ErrorMessageF, r)
	}
}