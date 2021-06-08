# Install a three-node k8s cluster with GPU device plugin with kubeadmin

1. $kubeadm init --pod-network-cidr=10.244.0.0/16
2. $export KUBECONFIG=/etc/kubernetes/admin.conf  #if run as root, for non-root copy configure to $HOME/.kube/config
3. $kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
4. $sudo swapoff -a  # on worker node, disable swap, before join 
5. $kubadm join XXXX
6. $kubectl create -f https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/1.0.0-beta4/nvidia-device-plugin.yml

when switch from calico to flannel, coreDNS pod stuck at container creating.
Solution: https://stackoverflow.com/questions/53900779/pods-failed-to-start-after-switch-cni-plugin-from-flannel-to-calico-and-then-fla

**Need to delete calico configure in worker node's /etc/cni/net.d/ as well** otherwise may have network plugin cannot create network for pod, x509 ...

worker node prepare
- nvidia drive 460
- nvidia container runtime
- nvidia docker2
