package certs

import (
	"testing"
)

func TestGenerateAll(t *testing.T) {
	BasePath := "/tmp/kubernetes/pki"
	EtcdBasePath := "/tmp/kubernetes/pki/etcd"
	apiServerIPAndDomains := []string{"apiserver.k8s.local", "8.8.8.8", "10.10.10.4", "k8s-node01", "10.10.10.5", "k8s-node02", "10.10.10.6", "k8s-node03", "10.10.10.1", "k8s-master01", "10.10.10.2", "k8s-master02", "10.10.10.3", "k8s-master03", "192.168.100.101"}
	etcdServerIPAndDomains := []string{"k8s.etcd.local", "10.10.10.1", "k8s-master01", "10.10.10.2", "k8s-master02", "10.10.10.3", "k8s-master03", "192.168.100.101"}
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"generate all certs",
			false,
		},
	}

	certMeta, err := NewCertMetaData(BasePath, EtcdBasePath, apiServerIPAndDomains, "10.10.10.1", "1.1.1.1", "172.31.0.0/16", "cluster.local", etcdServerIPAndDomains)
	if err != nil {
		t.Error(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := certMeta.GenerateAll(); (err != nil) != tt.wantErr {
				t.Errorf("GenerateAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
