package cmd

import (
	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:    "daemon",
	Short:  "Manage the koem daemon",
	Hidden: true,
}

var daemonRunCmd = &cobra.Command{
	Use:    "run",
	Short:  "Run the koem daemon (internal use)",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return impls.RunDaemon()
	},
}

func init() {
	daemonCmd.AddCommand(daemonRunCmd)
	rootCmd.AddCommand(daemonCmd)
}
