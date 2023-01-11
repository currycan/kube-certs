package certs

import (
	"testing"
)

func TestGenerateAll(t *testing.T) {
	BasePath := "/tmp/kubernetes/pki"
	EtcdBasePath := "/tmp/kubernetes/pki/etcd"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"generate all certs",
			false,
		},
	}

	certMeta, err := NewCertMetaData(BasePath, EtcdBasePath, []string{"k8s.local"}, []string{"10.88.88.88"}, "172.31.0.0/16", "master1", "10.88.88.1", "cluster.local", []string{"etcd.local"}, []string{"10.88.88.99"})
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
