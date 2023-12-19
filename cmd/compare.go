/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"plugin"

	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	branch1 string
	branch2 string
)

const (
	pluginPath      string = "./branchcomparer.so"
	extComparerType string = "ExtComparer"
)

type Comparer interface {
	Compare(b1 string, b2 string) (string, error)
}

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use: "compare",
	// Short: "...",
	// Long: `...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Working with branches: %s, %s.\n", branch1, branch2)
		fmt.Println("Loading module.")
		plug, err := plugin.Open(pluginPath)
		if err != nil {
			fmt.Printf("Open plugin (%s) is failed. Error: %v\n", pluginPath, err)
			os.Exit(1)
		}

		symComparer, err := plug.Lookup(extComparerType)
		if err != nil {
			return fmt.Errorf("Can't find type {%s} in plugin. Error: %v\n", extComparerType, err)
		}

		var comparer Comparer
		comparer, ok := symComparer.(Comparer)
		if !ok {
			return fmt.Errorf("Unexpected type from module symbol.")
		}

		fmt.Println("Comparer is loaded. Start comrare branches.")

		result, err := comparer.Compare(branch1, branch2)
		if err != nil {
			return err
		}

		fmt.Println("Finish comrare branches. Result:")
		fmt.Println(result)
		return nil
	},
}

func init() {
	compareCmd.Flags().StringVarP(&branch1, "branch1", "f", "", "Name of first branch for comparer (required)")
	compareCmd.Flags().StringVarP(&branch2, "branch2", "s", "", "Name of second branch for comparer (required)")
	compareCmd.MarkFlagRequired("branch1")
	compareCmd.MarkFlagRequired("branch2")
	rootCmd.AddCommand(compareCmd)
}
