
kubebuilder init --domain dcn.com --repo=github.com/binhfdv/sla-operator --skip-go-version-check
kubebuilder create api --group sla-operator --version v1alpha1 --kind Slaml
