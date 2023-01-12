package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/currycan/kube-certs/pkg/logger"
)

var (
	ConfigDir               = GetUserHomeDir() + "/.kube-certs"
	KubernetesDir           = "/etc/kubernetes"
	KubeDefaultCertPath     = "/etc/kubernetes/pki"
	kubeDefaultCertEtcdPath = "/etc/kubernetes/pki/etcd"
)

func GetUserHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home
}

func CaList(CertPath, CertEtcdPath string) []Config {
	return []Config{
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "ca",
			CommonName:   "kubernetes-ca",
			Organization: nil,
			Year:         100,
			AltNames:     AltNames{},
			Usages:       nil,
		},
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "front-proxy-ca",
			CommonName:   "kubernetes-front-proxy-ca",
			Organization: nil,
			Year:         100,
			AltNames:     AltNames{},
			Usages:       nil,
		},
		{
			Path:         CertEtcdPath,
			DefaultPath:  kubeDefaultCertEtcdPath,
			BaseName:     "ca",
			CommonName:   "etcd-ca",
			Organization: nil,
			Year:         100,
			AltNames:     AltNames{},
			Usages:       nil,
		},
	}
}

func List(CertPath, CertEtcdPath string) []Config {
	return []Config{
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "apiserver",
			CAName:       "kubernetes-ca",
			CommonName:   "kube-apiserver",
			Organization: nil,
			Year:         100,
			AltNames: AltNames{
				DNSNames: map[string]string{
					"localhost":              "localhost",
					"kubernetes":             "kubernetes",
					"kubernetes.default":     "kubernetes.default",
					"kubernetes.default.svc": "kubernetes.default.svc",
				},
				IPs: map[string]net.IP{
					"127.0.0.1": net.IPv4(127, 0, 0, 1),
				},
			},
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		},
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "apiserver-kubelet-client",
			CAName:       "kubernetes-ca",
			CommonName:   "kube-apiserver-kubelet-client",
			Organization: []string{"system:masters"},
			Year:         100,
			AltNames:     AltNames{},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "front-proxy-client",
			CAName:       "kubernetes-front-proxy-ca",
			CommonName:   "front-proxy-client",
			Organization: nil,
			Year:         100,
			AltNames:     AltNames{},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		{
			Path:         CertPath,
			DefaultPath:  KubeDefaultCertPath,
			BaseName:     "apiserver-etcd-client",
			CAName:       "etcd-ca",
			CommonName:   "kube-apiserver-etcd-client",
			Organization: []string{"system:masters"},
			Year:         100,
			AltNames:     AltNames{},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
		{
			Path:         CertEtcdPath,
			DefaultPath:  kubeDefaultCertEtcdPath,
			BaseName:     "server",
			CAName:       "etcd-ca",
			CommonName:   "kube-etcd", // kubeadm using node name as common name cc.CommonName = mc.NodeRegistration.Name
			Organization: nil,
			Year:         100,
			AltNames: AltNames{
				DNSNames: map[string]string{
					"localhost": "localhost",
				},
				IPs: map[string]net.IP{
					"127.0.0.1":                     net.IPv4(127, 0, 0, 1),
					net.IPv4(127, 0, 0, 1).String(): net.IPv4(127, 0, 0, 1),
					net.IPv6loopback.String():       net.IPv6loopback,
				},
			},
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		{
			Path:         CertEtcdPath,
			DefaultPath:  kubeDefaultCertEtcdPath,
			BaseName:     "peer",
			CAName:       "etcd-ca",
			CommonName:   "kube-etcd-peer", // change this in filter
			Organization: nil,
			Year:         100,
			AltNames: AltNames{
				DNSNames: map[string]string{
					"localhost": "localhost",
				},
				IPs: map[string]net.IP{
					"127.0.0.1":                     net.IPv4(127, 0, 0, 1),
					net.IPv4(127, 0, 0, 1).String(): net.IPv4(127, 0, 0, 1),
					net.IPv6loopback.String():       net.IPv6loopback,
				},
			},
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		},
		{
			Path:         CertEtcdPath,
			DefaultPath:  kubeDefaultCertEtcdPath,
			BaseName:     "healthcheck-client",
			CAName:       "etcd-ca",
			CommonName:   "kube-etcd-healthcheck-client",
			Organization: []string{"system:masters"},
			Year:         100,
			AltNames:     AltNames{},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
}

// 证书中需要用到的一些信息,传入的参数得提前验证
type CertMetaData struct {
	NodeName  string
	NodeIP    string
	DNSDomain string
	//kubernetes
	APIServer AltNames
	//etcd证书配置
	EtcdServer AltNames
	//证书生成的位置
	CertPath     string
	CertEtcdPath string
}

const (
	APIserverCert = iota
	APIserverKubeletClientCert
	FrontProxyClientCert
	APIserverEtcdClientCert
	EtcdServerCert
	EtcdPeerCert
	EtcdHealthcheckClientCert
)

// apiServerIPAndDomains = MasterIP + VIP + CertSANS 暂时只有apiserver, 记得把cluster.local后缀加到apiServerIPAndDOmas里先
func NewCertMetaData(certPATH, certEtcdPATH string, apiServerIPAndDomains []string, nodeIP, nodeName, SvcCIDR, DNSDomain string, etcdServerIPsAndDomains []string) (*CertMetaData, error) {
	data := &CertMetaData{}
	data.CertPath = certPATH
	data.CertEtcdPath = certEtcdPATH
	data.DNSDomain = DNSDomain
	data.APIServer.IPs = make(map[string]net.IP)
	data.APIServer.DNSNames = make(map[string]string)
	svcFirstIP, _, err := net.ParseCIDR(SvcCIDR)
	if err != nil {
		return nil, err
	}
	svcFirstIP[len(svcFirstIP)-1]++ //取svc第一个ip
	data.APIServer.IPs[svcFirstIP.String()] = svcFirstIP
	// etcd
	data.EtcdServer.DNSNames = make(map[string]string)
	data.EtcdServer.IPs = make(map[string]net.IP)

	//kubernetes
	fmt.Printf("apiServerIPAndDomains: %T, %v\n", apiServerIPAndDomains, apiServerIPAndDomains)
	for _, altName := range apiServerIPAndDomains {
		fmt.Printf("altName: %T, %v\n", altName, altName)
		ip := net.ParseIP(strings.TrimSpace(altName))
		fmt.Println("ip:", ip)
		if ip != nil {
			data.APIServer.IPs[ip.String()] = ip
			fmt.Println("data.APIServer.IPs:", data.APIServer.IPs)
			continue
		}
		data.APIServer.DNSNames[altName] = altName
		fmt.Println("data.APIServer.DNSNames:", data.APIServer.DNSNames)
	}

	//etcd
	for _, altName := range etcdServerIPsAndDomains {
		ip := net.ParseIP(strings.TrimSpace(altName))
		if ip != nil {
			data.EtcdServer.IPs[ip.String()] = ip
			continue
		}
		data.EtcdServer.DNSNames[altName] = altName
	}

	if ip := net.ParseIP(strings.TrimSpace(nodeIP)); ip != nil {
		data.APIServer.IPs[ip.String()] = ip
		data.EtcdServer.IPs[ip.String()] = ip
	}

	data.NodeIP = nodeIP
	data.NodeName = nodeName

	return data, nil
}

func (meta *CertMetaData) apiServerNameAndIP(certList *[]Config) {
	for _, dns := range meta.APIServer.DNSNames {
		(*certList)[APIserverCert].AltNames.DNSNames[dns] = dns
	}

	svcDNS := fmt.Sprintf("kubernetes.default.svc.%s", meta.DNSDomain)
	(*certList)[APIserverCert].AltNames.DNSNames[svcDNS] = svcDNS
	ip := net.ParseIP(strings.TrimSpace(meta.NodeName))
	if ip == nil {
		(*certList)[APIserverCert].AltNames.DNSNames[meta.NodeName] = meta.NodeName
	}

	for _, ip := range meta.APIServer.IPs {
		(*certList)[APIserverCert].AltNames.IPs[ip.String()] = ip
	}
	logger.Info("apiserver cert alt-names : %v", (*certList)[APIserverCert].AltNames)
}

func (meta *CertMetaData) etcdServerNameAndIP(certList *[]Config) {
	for _, etcdDns := range meta.EtcdServer.DNSNames {
		(*certList)[EtcdServerCert].AltNames.DNSNames[etcdDns] = etcdDns
		ip := net.ParseIP(strings.TrimSpace(meta.NodeName))
		if ip == nil {
			(*certList)[APIserverCert].AltNames.DNSNames[meta.NodeName] = meta.NodeName
		}
		(*certList)[EtcdPeerCert].AltNames.DNSNames[etcdDns] = etcdDns
	}

	for _, ip := range meta.EtcdServer.IPs {
		(*certList)[EtcdServerCert].AltNames.IPs[ip.String()] = ip
		(*certList)[EtcdPeerCert].AltNames.IPs[ip.String()] = ip
	}

	////修改 CommonName 为节点名
	//(*certList)[EtcdServerCert].CommonName = meta.NodeName
	//(*certList)[EtcdPeerCert].CommonName = meta.NodeName

	logger.Info("etcd cert alt-names  : %v, commonName : %s", (*certList)[EtcdPeerCert].AltNames, (*certList)[EtcdPeerCert].CommonName)
}

// create sa.key sa.pub for service Account
func (meta *CertMetaData) generatorServiceAccountKeyPaire() error {
	dir := meta.CertPath
	_, err := os.Stat(path.Join(dir, "sa.key"))
	if !os.IsNotExist(err) {
		logger.Info("sa.key sa.pub already exist")
		return nil
	}

	key, err := NewPrivateKey(x509.RSA)
	if err != nil {
		return err
	}
	pub := key.Public()

	err = WriteKey(dir, "sa", key)
	if err != nil {
		return err
	}

	return WritePublicKey(dir, "sa", pub)
}

func (meta *CertMetaData) GenerateAll() error {
	cas := CaList(meta.CertPath, meta.CertEtcdPath)
	certs := List(meta.CertPath, meta.CertEtcdPath)
	meta.apiServerNameAndIP(&certs)
	meta.etcdServerNameAndIP(&certs)
	_ = meta.generatorServiceAccountKeyPaire()

	CACerts := map[string]*x509.Certificate{}
	CAKeys := map[string]crypto.Signer{}
	for _, ca := range cas {
		caCert, caKey, err := NewCaCertAndKey(ca)
		if err != nil {
			return err
		}

		CACerts[ca.CommonName] = caCert
		CAKeys[ca.CommonName] = caKey

		err = WriteCertAndKey(ca.Path, ca.BaseName, caCert, caKey)
		if err != nil {
			return err
		}
	}

	for _, cert := range certs {
		caCert, ok := CACerts[cert.CAName]
		if !ok {
			return fmt.Errorf("root ca cert not found %s", cert.CAName)
		}
		caKey, ok := CAKeys[cert.CAName]
		if !ok {
			return fmt.Errorf("root ca key not found %s", cert.CAName)
		}

		Cert, Key, err := NewCaCertAndKeyFromRoot(cert, caCert, caKey)
		if err != nil {
			return err
		}
		err = WriteCertAndKey(cert.Path, cert.BaseName, Cert, Key)
		if err != nil {
			return err
		}
	}
	return nil
}
