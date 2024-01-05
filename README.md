# k8s-assistant: Your AI-Powered Kubernetes Assistant

## Setup

Obtain an API key from https://platform.openai.com/api-keys. Provide
API key using environment variable `OPENAI_API_KEY`.

Start k8s-assistant:

```bash
go run main.go
```

## Example

```
I am k8s-assistant, your AI-powered Kubernetes assistant. I will assist you in managing your Kubernetes Cluster. Please provide me with your instruction.
>>> get all pods

kubectl get pods --all-namespaces

NAMESPACE            NAME                                                  READY   STATUS    RESTARTS   AGE
kube-system          coredns-5d78c9869d-5h8wb                              1/1     Running   0          18h
kube-system          coredns-5d78c9869d-ndx48                              1/1     Running   0          18h
kube-system          etcd-k8s-assistant-control-plane                      1/1     Running   0          18h
kube-system          kindnet-z5nt2                                         1/1     Running   0          18h
kube-system          kube-apiserver-k8s-assistant-control-plane            1/1     Running   0          18h
kube-system          kube-controller-manager-k8s-assistant-control-plane   1/1     Running   0          18h
kube-system          kube-proxy-xg8nc                                      1/1     Running   0          18h
kube-system          kube-scheduler-k8s-assistant-control-plane            1/1     Running   0          18h
local-path-storage   local-path-provisioner-6bc4bddd6b-mb7dh               1/1     Running   0          18h

>>>
```

## Roadmap

- [ ] Accept flags to customize app behaviour.
- [ ] (UX) Stream executing command output to stdout instead of printing all at once.
