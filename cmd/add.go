/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/spf13/cobra"
)

var ports []string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [name] [port-range]",
	Short: "Add a new label",
	Args:  cobra.MaximumNArgs(1),
	Long: `Add a new label to your configuration with a name and port range.
	
This command allows you to add a new label to your configuration file.
The label will be added with a name and port range.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// fmt.Println("add called", args)
		return addLabel(args)
	},
}

func addLabel(args []string) error {
	label := &impls.Label{}

	if err := label.AddLabel(args[0], ports); err != nil {
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

	addCmd.Flags().StringSliceVarP(&ports, "ports", "p", []string{}, "Port range for the label")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
