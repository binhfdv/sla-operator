# sla_operator
1. Install kubebuilder (pre-installed, go version go1.16.7 linux/amd64)
# download kubebuilder and install locally.
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/v3.11.1/$(go env GOOS)/$(go env GOARCH)"
chmod +x kubebuilder && mv kubebuilder /usr/local/bin/
vim ~/.profile
add this to ~/.profile: export PATH=$PATH:/usr/local/kubebuilder/bin
source ~/.profile
2. 