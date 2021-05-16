SHELL=/bin/bash -o pipefail

.PHONY: generate
generate: generate-objects generate-crds

.PHONY: generate-crds
generate-crds:
	controller-gen crd paths=./pkg/apis/... output:stdout > deploy/crds.yaml

.PHONY: generate-objects
generate-objects:
	controller-gen object paths=./pkg/apis/...

.PHONY: format
format:
	goimports -l -w .