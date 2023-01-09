package core

import (
	"regexp"
	"strconv"

	"github.com/sealyun/lvscare/care"

	"github.com/currycan/supkube/cert"
	"github.com/currycan/supkube/ipvs"
	"github.com/currycan/supkube/pkg/sshcmd/sshutil"
)

var (
	MasterIPs         []string
	NodeIPs           []string
	CertSANS          []string
	DNSDomain         string
	APIServerCertSANs []string
	SSHConfig         sshutil.SSH
	APIServer         string
	CertPath          = cert.ConfigDir + "/pki"
	CertEtcdPath      = cert.ConfigDir + "/pki/etcd"
	EtcdCacart        = cert.ConfigDir + "/pki/etcd/ca.crt"
	EtcdCert          = cert.ConfigDir + "/pki/etcd/healthcheck-client.crt"
	EtcdKey           = cert.ConfigDir + "/pki/etcd/healthcheck-client.key"

	CriSocket    string
	CgroupDriver string
	KubeadmAPI   string

	VIP     string
	PkgURL  string
	Version string
	Repo    string
	PodCIDR string
	SvcCIDR string

	Envs          []string // read env from -e
	PackageConfig string   // install/delete package config
	Values        string   // values for  install package values.yaml
	WorkDir       string   // workdir for install/delete package home

	Ipvs         care.LvsCare
	LvscareImage ipvs.LvscareImage
	KubeadmFile  string

	Network string // network type, calico or flannel etc..

	WithoutCNI bool // if true don't install cni plugin

	Interface string //network interface name, like "eth.*|en.*"

	BGP bool // the ipip mode of the calico

	MTU string // mtu size

	YesRx = regexp.MustCompile("^(?i:y(?:es)?)$")

	CleanForce bool
	CleanAll   bool

	Vlog int

	InDocker     bool
	SnapshotName string
	EtcdBackDir  string
	RestorePath  string

	OssEndpoint      string
	AccessKeyID      string
	AccessKeySecrets string
	BucketName       string
	ObjectPath       string
)

func vlogToStr() string {
	str := strconv.Itoa(Vlog)
	return " -v " + str
}

type metadata struct {
	K8sVersion string `json:"k8sVersion"`
	CniVersion string `json:"cniVersion"`
	CniName    string `json:"cniName"`
}
