/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/currycan/kube-certs/pkg/version"

	"github.com/spf13/cobra"
)

var shortPrint bool

func newVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "version",
		Args:    cobra.NoArgs,
		Example: `kube-certs version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			marshalled, err := json.Marshal(version.Get())
			if err != nil {
				return err
			}
			if shortPrint {
				fmt.Println(version.Get().String())
			} else {
				fmt.Println(string(marshalled))
			}
			return nil
		},
	}
	versionCmd.Flags().BoolVar(&shortPrint, "short", false, "if true, print just the version number.")
	return versionCmd
}

func init() {
	rootCmd.AddCommand(newVersionCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
