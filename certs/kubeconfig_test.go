package certs

import (
	"testing"
)

func TestGenerateKubeconfig(t *testing.T) {
	certPath := "/tmp/kubernetes/pki"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			"generate k8s kubeconfig",
			false,
		},
	}
	certConfig := Config{
		Path:     certPath,
		BaseName: "ca",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateJoinControlPlaneKubeConfigFiles("/tmp/kubernetes",
				certConfig, "192.168.1.1", "https://apiserver.k8s.local:6443", "kubernetes"); (err != nil) != tt.wantErr {
				t.Errorf("GenerateAll() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
