# sla_operator

kubernetes v1.23.6

# 1. Install kubebuilder
## golang
```
curl -OL https://golang.org/dl/go1.20.7.linux-amd64.tar.gz

sudo tar -C /usr/local -xvf go1.20.7.linux-amd64.tar.gz

add to ~/.profile: export PATH=$PATH:/usr/local/go/bin
```

## download kubebuilder and install locally.
```
curl -L -o kubebuilder "https://go.kubebuilder.io/dl/v3.11.1/$(go env GOOS)/$(go env GOARCH)"

sudo su

sudo chmod +x kubebuilder && sudo mv kubebuilder /usr/local/bin/

vim ~/.profile

add this to ~/.profile: export PATH=$PATH:/usr/local/kubebuilder/bin

source ~/.profile
```

# 2. Install Volcano
kubectl apply -f https://raw.githubusercontent.com/volcano-sh/volcano/master/installer/volcano-development.yaml

version: 1.8.0


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
- [Semantic Commit Messages]: https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716