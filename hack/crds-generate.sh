#!/usr/bin/env bash

set -euo pipefail
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CRD_DIR="${1}"
HELM_DIR="${2}"

cd "${SCRIPT_DIR}"/../
cp ${CRD_DIR}/*.yaml ${HELM_DIR}/charts/crds/crds/

{
  for file in "${HELM_DIR}"/charts/crds/crds/*.yaml; do
    cat "${file}"
    echo "---"
  done
} | bzip2 --best --compress --keep --stdout - >"${HELM_DIR}/charts/crds/files/crds.bz2"
