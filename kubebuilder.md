
kubebuilder init --domain dcn.com --repo=github.com/binhfdv/sla-operator
kubebuilder create api --group sla-operator --version v1alpha1 --kind Slaml

cd sla-operator
make manifests
make install
make run


make docker-build docker-push IMG=<some-registry>/<project-name>:tag
make docker-build docker-push IMG=ddocker122/sla-test:v1