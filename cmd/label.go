/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// labelCmd represents the label command
var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage label configurations",
	Long: `Manage label configurations of Koem-cli
	Add new labels to your configuration with a name and port range`,
	// RunE: func(cmd *cobra.Command, args []string) error {
	// 	fmt.Println("label called", args)
	// 	return nil
	// },
}

func init() {
	rootCmd.AddCommand(labelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// labelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// labelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
