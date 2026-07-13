package cmd

import (
	"fmt"
	"os"

	"github.com/Emmanuel-Soempit/koem/impls"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var reservesLabelName string

var reservesCmd = &cobra.Command{
	Use:   "reserves",
	Short: "Manage port reserves for a label",
	Long:  `List or add port reserves for a label.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listReserves(reservesLabelName)
	},
}

var reservesAddCmd = &cobra.Command{
	Use:   "add <app_name>",
	Short: "Reserve 3 ports for an app under a label",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if reservesLabelName == "" {
			return fmt.Errorf("-l <label> is required")
		}
		return addReserve(reservesLabelName, args[0])
	},
}

func listReserves(labelName string) error {
	labels, err := impls.LoadAll()
	if err != nil {
		return err
	}

	if labelName != "" {
		var found bool
		for _, l := range labels {
			if l.Name == labelName {
				found = true
				labels = []impls.Label{l}
				break
			}
		}
		if !found {
			return fmt.Errorf("label %q not found", labelName)
		}
	}

	table := tablewriter.NewTable(os.Stdout)
	table.Header("Label", "App", "Production", "Preview", "Development")

	total := 0
	for _, l := range labels {
		reserves, err := impls.LoadReserves(l.Name)
		if err != nil {
			return err
		}
		for _, r := range reserves {
			table.Append([]any{l.Name, r.AppName, r.Production, r.Preview, r.Development})
			total++
		}
	}

	if total == 0 {
		fmt.Println("No reserves found.")
		return nil
	}

	table.Render()
	return nil
}

func addReserve(labelName, appName string) error {
	labels, err := impls.LoadAll()
	if err != nil {
		return err
	}

	var label *impls.Label
	for _, l := range labels {
		if l.Name == labelName {
			l := l
			label = &l
			break
		}
	}
	if label == nil {
		return fmt.Errorf("label %q not found", labelName)
	}

	freePorts, err := impls.FindFreePorts(label.Min, label.Max, 3)
	if err != nil {
		return err
	}
	if len(freePorts) < 3 {
		return fmt.Errorf("not enough free ports in range %s-%s", label.Min, label.Max)
	}

	if err := impls.SaveReserve(labelName, appName, freePorts); err != nil {
		return err
	}
	if err := impls.EnsureDaemon(); err != nil {
		return fmt.Errorf("could not start daemon: %w", err)
	}
	if err := impls.SendReserve(freePorts); err != nil {
		return fmt.Errorf("daemon failed to hold ports: %w", err)
	}

	envs := []string{"Production", "Preview", "Development"}
	fmt.Printf("\nReserved ports for %q under %q:\n", appName, labelName)
	table := tablewriter.NewTable(os.Stdout)
	table.Header("Environment", "Port")
	for i, port := range freePorts {
		table.Append([]any{envs[i], fmt.Sprintf("%d", port)})
	}
	table.Render()
	return nil
}

var reservesClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all reserves and free their ports",
	Long:  `Release all reserved ports and remove reservations. Use -l to target a specific label.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return clearReserves(reservesLabelName)
	},
}

func clearReserves(labelName string) error {
	labels, err := impls.LoadAll()
	if err != nil {
		return err
	}

	if labelName != "" {
		var found bool
		for _, l := range labels {
			if l.Name == labelName {
				found = true
				labels = []impls.Label{l}
				break
			}
		}
		if !found {
			return fmt.Errorf("label %q not found", labelName)
		}
	}

	total := 0
	for _, l := range labels {
		reserves, err := impls.LoadReserves(l.Name)
		if err != nil {
			return err
		}
		var ports []int
		for _, r := range reserves {
			p, err := impls.PortsFromReserve(r)
			if err == nil {
				ports = append(ports, p...)
			}
		}
		if len(ports) > 0 {
			impls.SendRelease(ports)
		}
		if err := impls.ClearReserves(l.Name); err != nil {
			return err
		}
		total += len(reserves)
	}

	if labelName != "" {
		fmt.Printf("Cleared %d reserve(s) for label %q.\n", total, labelName)
	} else {
		fmt.Printf("Cleared %d reserve(s) across all labels.\n", total)
	}
	return nil
}

func init() {
	reservesCmd.PersistentFlags().StringVarP(&reservesLabelName, "label", "l", "", "Label name")
	reservesCmd.AddCommand(reservesAddCmd)
	reservesCmd.AddCommand(reservesClearCmd)
	getCmd.AddCommand(reservesCmd)
}
