#!/usr/bin/env bash
set -euo pipefail
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
CRD_DIR="${1}"
HELM_DIR="${2}"

if [[ "$OSTYPE" == "darwin"* ]]; then
  SEDPRG="gsed"
else
  SEDPRG="sed"
fi

cd "${SCRIPT_DIR}"/../

cp ${CRD_DIR}/*.yaml ${HELM_DIR}/crds/

# Remove extra header lines in transformed CRDs
for f in "${HELM_DIR}"/crds/*.yaml; do
  tail -n +2 < "$f" > "$f.bkp"
  cp "$f.bkp" "$f"
  rm "$f.bkp"
done

# Add helm if statement for controlling the install of CRDs
for i in "${HELM_DIR}"/crds/*.yaml; do
  cp "$i" "$i.bkp"
  echo "{{- if .Values.installCRDs }}" > "$i"
  cat "$i.bkp" >> "$i"
  echo "{{- end }}" >> "$i"
  rm "$i.bkp"
done
