package certs

import (
	"crypto"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"path"

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
			AltIPs:       AltIPs{},
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
			AltIPs:       AltIPs{},
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
			AltIPs:       AltIPs{},
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
			},
			AltIPs: AltIPs{
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
			AltIPs:       AltIPs{},
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
			AltIPs:       AltIPs{},
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
			AltIPs:       AltIPs{},
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
			},
			AltIPs: AltIPs{
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
			},
			AltIPs: AltIPs{
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
			AltIPs:       AltIPs{},
			Usages:       []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		},
	}
}

// 证书中需要用到的一些信息,传入的参数得提前验证
type CertMetaData struct {
	NodeName  string
	NodeIP    string
	DNSDomain string
	//kubernetes证书配置
	APIServerNames AltNames
	APIServerIPs   AltIPs
	//etcd证书配置
	EtcdServerNames AltNames
	EtcdServerIPs   AltIPs
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
func NewCertMetaData(certPATH, certEtcdPATH string, apiServerDomains, apiServerIPs []string, SvcCIDR, nodeName, nodeIP, DNSDomain string, etcdServerDomains, etcdServerIPs []string) (*CertMetaData, error) {
	data := &CertMetaData{}
	data.CertPath = certPATH
	data.CertEtcdPath = certEtcdPATH
	data.DNSDomain = DNSDomain
	data.APIServerIPs.IPs = make(map[string]net.IP)
	data.APIServerNames.DNSNames = make(map[string]string)
	svcFirstIP, _, err := net.ParseCIDR(SvcCIDR)
	if err != nil {
		return nil, err
	}
	svcFirstIP[len(svcFirstIP)-1]++ //取svc第一个ip
	data.APIServerIPs.IPs[svcFirstIP.String()] = svcFirstIP

	data.EtcdServerNames.DNSNames = make(map[string]string)
	data.EtcdServerIPs.IPs = make(map[string]net.IP)

	//kubernetes
	for _, altName := range apiServerDomains {
		data.APIServerNames.DNSNames[altName] = altName
	}
	for _, altIP := range apiServerIPs {
		ip := net.ParseIP(altIP)
		if ip != nil {
			data.APIServerIPs.IPs[ip.String()] = ip
			continue
		}
	}

	//etcd
	for _, etcdAltName := range etcdServerDomains {
		data.EtcdServerNames.DNSNames[etcdAltName] = etcdAltName
	}
	for _, etcdAltIP := range etcdServerIPs {
		ip := net.ParseIP(etcdAltIP)
		if ip != nil {
			data.EtcdServerIPs.IPs[ip.String()] = ip
			continue
		}
	}

	if ip := net.ParseIP(nodeIP); ip != nil {
		data.APIServerIPs.IPs[ip.String()] = ip
		data.EtcdServerIPs.IPs[ip.String()] = ip
	}

	data.NodeIP = nodeIP
	data.NodeName = nodeName

	return data, nil
}

func (meta *CertMetaData) apiServerNameAndIP(certList *[]Config) {
	for _, dns := range meta.APIServerNames.DNSNames {
		(*certList)[APIserverCert].AltNames.DNSNames[dns] = dns
	}

	svcDNS := fmt.Sprintf("kubernetes.default.svc.%s", meta.DNSDomain)
	(*certList)[APIserverCert].AltNames.DNSNames[svcDNS] = svcDNS
	(*certList)[APIserverCert].AltNames.DNSNames[meta.NodeName] = meta.NodeName

	logger.Info("apiserver Domains : %v", (*certList)[APIserverCert].AltNames)

	for _, ip := range meta.APIServerIPs.IPs {
		(*certList)[APIserverCert].AltIPs.IPs[ip.String()] = ip
	}
	logger.Info("apiserver IPs : %v", (*certList)[APIserverCert].AltIPs)
}

func (meta *CertMetaData) etcdServerNameAndIP(certList *[]Config) {
	for _, etcdDns := range meta.EtcdServerNames.DNSNames {
		(*certList)[EtcdServerCert].AltNames.DNSNames[etcdDns] = etcdDns
		(*certList)[EtcdPeerCert].AltNames.DNSNames[etcdDns] = etcdDns
	}

	for _, ip := range meta.EtcdServerIPs.IPs {
		(*certList)[EtcdServerCert].AltIPs.IPs[ip.String()] = ip
		(*certList)[EtcdPeerCert].AltIPs.IPs[ip.String()] = ip
	}

	////修改 CommonName 为节点名
	//(*certList)[EtcdServerCert].CommonName = meta.NodeName
	//(*certList)[EtcdPeerCert].CommonName = meta.NodeName

	logger.Info("Etcd Domains : %v, commonName : %s", (*certList)[EtcdPeerCert].AltNames, (*certList)[EtcdPeerCert].CommonName)
	logger.Info("Etcd IPs : %v, commonName : %s", (*certList)[EtcdPeerCert].AltIPs, (*certList)[EtcdPeerCert].CommonName)
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
