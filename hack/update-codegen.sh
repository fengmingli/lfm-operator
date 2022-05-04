#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

bash ../vendor/k8s.io/code-generator/generate-groups.sh "deepcopy" \
  lfm-operator/generated \
  lfm-operator/api \
    :v1 \
  --go-header-file $(pwd)/boilerplate.go.txt \
  --output-base $(pwd)/../../
# To use your own boilerplate text append:
#   --go-header-file "${SCRIPT_ROOT}"/hack/custom-boilerplate.go.txt