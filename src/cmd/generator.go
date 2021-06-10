// Package cmd represents Command Line Interface
//
// Copyright © 2020 Vitalii Noha vitalii.noga@gmail.com
package cmd

import (
	"github.com/sapplications/sbuilder/src/services/cmd"
	"github.com/spf13/cobra"
)

type Generator struct {
	Generator cmd.Generator
	cobra.Command
}

func (v *Generator) init() {
	v.Command.RunE = func(cmd *cobra.Command, args []string) error {
		if v.Generator == nil {
			return nil
		}
		if len(args) > 0 {
			return v.Generator.Generate(args[0])
		} else {
			return v.Generator.Generate("")
		}
	}
}
