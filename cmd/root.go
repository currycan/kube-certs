/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/currycan/kube-certs/certs"
	"github.com/currycan/kube-certs/core"
	"github.com/currycan/kube-certs/pkg/logger"
)

var (
	cfgFile string
	Info    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-certs",
	Short: "generate kubernetes cers and kubeconfig files",
	Long:  `generate kubernetes cers and kubeconfig files`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kube-certs/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&Info, "info", false, "logger ture for Info, false for Debug")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home := certs.GetUserHomeDir()
	logFile := fmt.Sprintf("%s/.kube-certs/kube-certs.log", home)
	if !core.FileExist(home + "/.kube-certs") {
		err := os.MkdirAll(home+"/.kube-certs", os.ModePerm)
		if err != nil {
			fmt.Println("create default kube-certs config dir failed, please create it by your self mkdir -p /root/.kube-certs && touch /root/.kube-certs/config.yaml")
		}
	}
	if Info {
		logger.Cfg(5, logFile)
	} else {
		logger.Cfg(6, logFile)
	}
}
