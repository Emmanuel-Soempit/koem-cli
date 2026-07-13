package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var docsCmd = &cobra.Command{
	Use:    "docs",
	Short:  "Generate markdown documentation for all commands",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := os.MkdirAll("./docs", 0755); err != nil {
			return err
		}
		if err := doc.GenMarkdownTree(rootCmd, "./docs"); err != nil {
			return err
		}
		fmt.Println("Documentation generated in ./docs/")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docsCmd)
}
