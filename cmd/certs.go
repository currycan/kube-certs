/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	cert "github.com/currycan/kube-certs/certs"

	"github.com/spf13/cobra"
)

type Flag struct {
	AltNames     []string
	NodeName     string
	ServiceCIDR  string
	NodeIP       string
	DNSDomain    string
	CertPath     string
	CertEtcdPath string
}

var config *Flag

// certsCmd represents the certs command
var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "generate kubernetes certs",
	Long:  `generate kubernetes certs: etcd and kubernetes component, expire time can be specified`,
	Run: func(cmd *cobra.Command, args []string) {
		cert.GenerateCert(config.CertPath, config.CertEtcdPath, config.AltNames, config.NodeIP, config.NodeName, config.ServiceCIDR, config.DNSDomain)
	},
}

func init() {
	config = &Flag{}
	rootCmd.AddCommand(certsCmd)

	certsCmd.Flags().StringSliceVar(&config.AltNames, "alt-names", []string{}, "like example.com or 10.88.88.88")
	certsCmd.Flags().StringVar(&config.NodeName, "node-name", "", "like master0")
	certsCmd.Flags().StringVar(&config.ServiceCIDR, "service-cidr", "", "like 10.88.88.0/24")
	certsCmd.Flags().StringVar(&config.NodeIP, "node-ip", "", "like 10.88.88.1")
	certsCmd.Flags().StringVar(&config.DNSDomain, "dns-domain", "cluster.local", "cluster dns domain")
	certsCmd.Flags().StringVar(&config.CertPath, "cert-path", "/etc/kubernetes/pki", "kubernetes cert file path")
	certsCmd.Flags().StringVar(&config.CertEtcdPath, "cert-etcd-path", "/etc/kubernetes/pki/etcd", "kubernetes etcd cert file path")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// certsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// certsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
