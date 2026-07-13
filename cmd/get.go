package cmd

import (
	"fmt"
	"os"

	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var getLabelName string

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources from your koem configuration",
}

var getLabelsCmd = &cobra.Command{
	Use:     "labels",
	Aliases: []string{"label"},
	Short:   "Get labels and their port ranges",
	Long:    `Display all configured labels with their port ranges in a tabular format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return getLabels(getLabelName)
	},
}

func getLabels(name string) error {
	labels, err := impls.LoadAll()
	if err != nil {
		return err
	}

	if name != "" {
		filtered := make([]impls.Label, 0)
		for _, l := range labels {
			if l.Name == name {
				filtered = append(filtered, l)
				break
			}
		}
		if len(filtered) == 0 {
			return fmt.Errorf("label %q not found", name)
		}
		labels = filtered
	}

	table := tablewriter.NewTable(os.Stdout)
	table.Header("Label", "Min", "Max", "In Use", "Reserves")

	for _, l := range labels {
		inUse, err := impls.CountPortsInUse(l.Min, l.Max)
		if err != nil {
			return err
		}
		reserves, _ := impls.LoadReserves(l.Name)
		table.Append([]any{l.Name, l.Min, l.Max, fmt.Sprintf("%d", inUse), fmt.Sprintf("%d", len(reserves))})
	}
	table.Render()
	return nil
}

func init() {
	getLabelsCmd.Flags().StringVarP(&getLabelName, "label", "l", "", "Filter by label name")
	getCmd.AddCommand(getLabelsCmd)
	rootCmd.AddCommand(getCmd)
}
