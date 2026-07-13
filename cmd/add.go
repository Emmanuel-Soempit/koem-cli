/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <name> <min> <max>",
	Short: "Add a new label",
	Args:  cobra.ExactArgs(3),
	Long: `Add a new label to your configuration with a name and port range.

Example:
  koem-cli label add backend 8000 9000`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return addLabel(args)
	},
}

func addLabel(args []string) error {
	label := &impls.Label{}

	if err := label.AddLabel(args[0], []string{args[1], args[2]}); err != nil {
		return err
	}

	if err := label.Save(); err != nil {
		return err
	}
	fmt.Println("Label added successfully")
	return nil
}

func init() {
	labelCmd.AddCommand(addCmd)
}
