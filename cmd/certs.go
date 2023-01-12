/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/currycan/kube-certs/certs"
	"github.com/currycan/kube-certs/pkg/logger"

	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
)

type Flag struct {
	APIServerAltNames    []string
	NodeName             string
	ServiceCIDR          string
	NodeIP               string
	DNSDomain            string
	CertPath             string
	EtcdAltNames         []string
	CertEtcdPath         string
	KubeConfigPath       string
	CertConfig           cert.Config
	ControlPlaneEndpoint string
	ClusterName          string
}

var config *Flag

// certsCmd represents the certs command
var certsCmd = &cobra.Command{
	Use:   "certs",
	Short: "generate kubernetes certs",
	Long:  `generate kubernetes certs: etcd and kubernetes component, expire time can be specified`,
	Run: func(cmd *cobra.Command, args []string) {
		certConfig := certs.Config{
			Path:     config.CertPath,
			BaseName: "ca",
		}

		certs.GenerateCert(config.CertPath, config.CertEtcdPath, config.APIServerAltNames, config.NodeIP, config.NodeName, config.ServiceCIDR, config.DNSDomain, config.APIServerAltNames)
		err := certs.CreateJoinControlPlaneKubeConfigFiles(config.KubeConfigPath, certConfig, config.NodeName, config.ControlPlaneEndpoint, config.ClusterName)
		if err != nil {
			logger.Error("generator kubeconfig failed %s", err)
			os.Exit(-1)
		}
	},
}

func init() {
	config = &Flag{}
	rootCmd.AddCommand(certsCmd)

	// kubernetes certs
	certsCmd.Flags().StringVar(&config.NodeName, "node-name", "", "like master0 or 10.88.88.1")
	certsCmd.Flags().StringVar(&config.NodeIP, "node-ip", "", "like 10.88.88.1")
	certsCmd.Flags().StringVar(&config.ServiceCIDR, "service-cidr", "", "like 172.31.0.0/16")
	certsCmd.Flags().StringVar(&config.DNSDomain, "dns-domain", "cluster.local", "cluster dns domain")
	certsCmd.Flags().StringSliceVar(&config.APIServerAltNames, "apiserver-alt-names", []string{}, "like example.com or 10.88.88.88")
	certsCmd.Flags().StringVar(&config.CertPath, "cert-path", "/etc/kubernetes/pki", "kubernetes cert file path")

	// etcd certs
	certsCmd.Flags().StringSliceVar(&config.EtcdAltNames, "etcd-alt-names", []string{}, "like example.com or 10.88.88.99")
	certsCmd.Flags().StringVar(&config.CertEtcdPath, "cert-etcd-path", "/etc/kubernetes/pki/etcd", "kubernetes etcd cert file path")

	// kubernetes kubeconfig
	certsCmd.Flags().StringVar(&config.KubeConfigPath, "kube-config-path", "/etc/kubernetes/", "kubernetes kubeconfig files path")
	certsCmd.Flags().StringVar(&config.ControlPlaneEndpoint, "control-plane-endpoint", "", "like https://apiserver.k8s.local:6443")
	certsCmd.Flags().StringVar(&config.ClusterName, "cluster-name", "kubernetes", "kubernetes cluster name")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// certsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// certsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
