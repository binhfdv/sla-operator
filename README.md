# sla_operator

kubernetes v1.23.6

# 1. Install kubebuilder
## golang 1.20.7

```
# remove go
sudo rm -rf /usr/local/go
dpkg -l | grep golang
sudo apt-get remove golang*
```

```
curl -OL https://golang.org/dl/go1.20.7.linux-amd64.tar.gz

sudo tar -C /usr/local -xvf go1.20.7.linux-amd64.tar.gz

# add to ~/.profile:
export PATH=$PATH:/usr/local/go/bin
go version
```

## Download kubebuilder and install locally.
```
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/v3.11.1/$(go env GOOS)/$(go env GOARCH)"

sudo su

sudo chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/

vim ~/.profile

add this to ~/.profile: export PATH=$PATH:/usr/local/kubebuilder/bin

source ~/.profile
```

# 2. Install Volcano v1.8.0

```
kubectl create namespace volcano-system
kubectl apply -f https://raw.githubusercontent.com/volcano-sh/volcano/v1.8.0/installer/volcano-development.yaml
```

# 3. Install helm-chart and prometheus, grfana

```
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
```

```
cd installation/prometheus-chart
helm install kube-prometheus-stack kube-prometheus-stack/

kubectl edit svc kube-prometheus-stack-prometheus
'change type: ClusterIP -> type: NodePort'

kubectl edit svc kube-prometheus-stack-grafana
'change type: ClusterIP -> type: NodePort'

# grafana account:
user: admin
password: prom-operator
```

```
prometheus/grafana web ui:
<master ip>:<port>
# port obtained from
kubectl get svc -A | grep NodePort
```

# Semantic commit
```
- feat: (new feature for the user, not a new feature for build script)
- fix: (bug fix for the user, not a fix to a build script)
- docs: (changes to the documentation)
- style: (formatting, missing semi colons, etc; no production code change)
- refactor`: (refactoring production code, eg. renaming a variable)
- test: (adding missing tests, refactoring tests; no production code change)
- chore: (updating grunt tasks etc; no production code change)
```

# References
- [Semantic Commit Messages](https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716)