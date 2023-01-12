# 一键创建 kubernetes 100 年过期时间证书和认证文件

## 编译二进制文件

生成 kube-certs 文件，放到可执行目录下（如：/usr/bin/）

```bash
git clone https://github.com/currycan/kube-certs.git
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -ldflags "-s -w"
```

## 创建集群证书文件

```bash
kube-certs certs \
  --cert-path=/etc/kubernetes/pki \
  --cert-etcd-path=./etc/kubernetes/pki/etcd \
  --kube-config-path=./etc/kubernetes \
  --node-ip=10.10.10.3 \
  # hostoverribe 成节点 IP，改成 IP 地址
  --node-name=10.10.10.3 \
  --dns-domain=cluster.local \
  # service 子网
  --service-cidr=172.31.0.0/16 \
  # 高可用
  --control-plane-endpoint="https://apiserver.k8s.local:6443" \
  --cluster-name="kubernetes" \
  # master 节点
  --apiserver-alt-names="apiserver.k8s.local, 8.8.8.8, 10.10.10.1, k8s-master01, 10.10.10.2, k8s-master02, 10.10.10.3, k8s-master03, 192.168.100.101" \
  # etcd 节点
  --etcd-alt-names="k8s.etcd.local, 10.10.10.1, k8s-master01, 10.10.10.2, k8s-master02, 10.10.10.3, k8s-master03, 192.168.100.101"
```
