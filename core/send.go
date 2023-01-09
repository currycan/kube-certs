package core

import (
	"fmt"
	"path"
)

//SendPackage is
func (s *Installer) SendPackage() {
	pkg := path.Base(PkgURL)
	// rm old  in package avoid old version problem. if  not exist in package then skip rm
	kubeHook := fmt.Sprintf("cd /root && rm -rf kube && tar zxvf %s  && cd /root/kube/shell && rm -f ../bin/ && bash init.sh", pkg)
	deletekubectl := `sed -i '/kubectl/d;//d' /root/.bashrc `
	completion := "echo 'command -v kubectl &>/dev/null && source <(kubectl completion bash)' >> /root/.bashrc && echo '[ -x /usr/bin/ ] && source <( completion bash)' >> /root/.bashrc && source /root/.bashrc"
	kubeHook = kubeHook + " && " + deletekubectl + " && " + completion
	PkgURL = SendPackage(PkgURL, s.Hosts, "/root", nil, &kubeHook)
}

// Send is send the exec  to /usr/bin/
func (s *Installer) Send() {
	// send  first to avoid old version
	filePath := FetchAbsPath()
	beforeHook := "ps -ef |grep -v 'grep'|grep  >/dev/null || rm -rf /usr/bin/"
	afterHook := "chmod a+x /usr/bin/"
	SendPackage(filePath, s.Hosts, "/usr/bin", &beforeHook, &afterHook)
}

// SendPackage is send new pkg to all nodes.
func (u *Upgrade) SendPackage() {
	all := append(u.Masters, u.Nodes...)
	pkg := path.Base(u.NewPkgURL)
	// rm old  in package avoid old version problem. if  not exist in package then skip rm
	var kubeHook string
	if For120(Version) {
		// TODO update need load modprobe -- br_netfilter modprobe -- bridge.
		// https://github.com/fanux/cloud-kernel/issues/23
		kubeHook = fmt.Sprintf("cd /root && rm -rf kube && tar zxvf %s  && cd /root/kube/shell && rm -f ../bin/ && (ctr -n=k8s.io image import ../images/images.tar || true) && cp -f ../bin/* /usr/bin/ ", pkg)
	} else {
		kubeHook = fmt.Sprintf("cd /root && rm -rf kube && tar zxvf %s  && cd /root/kube/shell && rm -f ../bin/ && (docker load -i ../images/images.tar || true) && cp -f ../bin/* /usr/bin/ ", pkg)
	}
	PkgURL = SendPackage(pkg, all, "/root", nil, &kubeHook)
}
