package certs

import (
	"fmt"
	"os"

	"github.com/currycan/kube-certs/pkg/logger"
)

func CMD(altNames []string, hostIP, hostName, serviceCIRD, DNSDomain string) string {
	cmd := "kube-certs cert "
	if hostIP != "" {
		cmd += fmt.Sprintf(" --node-ip %s", hostIP)
	}

	if hostName != "" {
		cmd += fmt.Sprintf(" --node-name %s", hostName)
	}

	if serviceCIRD != "" {
		cmd += fmt.Sprintf(" --service-cidr %s", serviceCIRD)
	}

	if DNSDomain != "" {
		cmd += fmt.Sprintf(" --dns-domain %s", DNSDomain)
	}

	for _, name := range altNames {
		if name != "" {
			cmd += fmt.Sprintf(" --alt-names %s", name)
		}
	}

	return cmd
}

// GenerateCert generate all cert.
func GenerateCert(certPATH, certEtcdPATH string, apiServerDomains, apiServerIPs []string, hostIP, hostName, serviceCIRD, DNSDomain string, etcdServerDomains, etcdServerIPs []string) {
	certConfig, err := NewCertMetaData(certPATH, certEtcdPATH, apiServerDomains, apiServerIPs, serviceCIRD, hostName, hostIP, DNSDomain, etcdServerDomains, etcdServerIPs)
	if err != nil {
		logger.Error("generator cert config failed %s", err)
		os.Exit(-1)
	}
	err = certConfig.GenerateAll()
	if err != nil {
		logger.Error("GenerateAll err %s", err)
	}
}
