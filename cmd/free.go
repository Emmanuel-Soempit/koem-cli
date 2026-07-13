package cmd

import (
	"fmt"
	"os"

	"github.com/Emmanuel-Soempit/koem-cli/impls"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var freeLabelName string

var freeCmd = &cobra.Command{
	Use:   "free",
	Short: "Suggest 3 free ports for a label",
	Long:  `Find and suggest 3 free ports within a label's port range for Production, Preview, and Development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if freeLabelName == "" {
			return fmt.Errorf("-l <label> is required")
		}
		return suggestFreePorts(freeLabelName)
	},
}

func suggestFreePorts(name string) error {
	labels, err := impls.LoadAll()
	if err != nil {
		return err
	}
	var label *impls.Label
	for _, l := range labels {
		if l.Name == name {
			l := l
			label = &l
			break
		}
	}
	if label == nil {
		return fmt.Errorf("label %q not found", name)
	}
	freePorts, err := impls.FindFreePorts(label.Min, label.Max, 3)
	if err != nil {
		return err
	}
	if len(freePorts) == 0 {
		fmt.Println("No free ports available in this range.")
		return nil
	}
	envs := []string{"Production", "Preview", "Development"}
	fmt.Printf("Suggested free ports for %q:\n", label.Name)
	table := tablewriter.NewTable(os.Stdout)
	table.Header("Environment", "Port")
	for i, port := range freePorts {
		env := ""
		if i < len(envs) {
			env = envs[i]
		}
		table.Append([]any{env, fmt.Sprintf("%d", port)})
	}
	table.Render()
	return nil
}

func init() {
	freeCmd.Flags().StringVarP(&freeLabelName, "label", "l", "", "Label name")
	getCmd.AddCommand(freeCmd)
}
