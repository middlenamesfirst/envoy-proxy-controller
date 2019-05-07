#!/bin/sh

set -e # fail immediately on errors

if [ -z "${KUBECONFIG}" ]; then
  echo "KUBECONFIG variable not defined." >&2
  exit 1
fi

if [ -z "${NAMESPACE}" ]; then
  echo "NAMESPACE variable not defined." >&2
  exit 1
fi

if [ -z "${SERVER}" ]; then
    echo "SERVER variable not defined." >&2
    exit 1
fi

kubectl --kubeconfig="${KUBECONFIG}" --namespace="${NAMESPACE}" apply -f ./envoy-proxy-controller-service-account.yaml

< ./envoy-proxy-controller-cluster-role-binding.yaml.template \
sed "s#{{NAMESPACE}}#${NAMESPACE}#g" | \
kubectl --kubeconfig="${KUBECONFIG}" --namespace="${NAMESPACE}" apply -f -

SECRET_NAME=$(kubectl --kubeconfig="${KUBECONFIG}" --namespace="${NAMESPACE}" get serviceaccount envoy-proxy-controller -o jsonpath='{.secrets[0].name}')
CA=$(kubectl --kubeconfig="${KUBECONFIG}" --namespace="${NAMESPACE}" get secret/"${SECRET_NAME}" -o jsonpath='{.data.ca\.crt}')
TOKEN=$(kubectl --kubeconfig="${KUBECONFIG}" --namespace="${NAMESPACE}" get secret/"${SECRET_NAME}" -o jsonpath='{.data.token}' | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: default-cluster
  cluster:
    certificate-authority-data: ${CA}
    server: ${SERVER}
contexts:
- name: default-context
  context:
    cluster: default-cluster
    namespace: default
    user: default-user
current-context: default-context
users:
- name: default-user
  user:
    token: ${TOKEN}" \
> envoy-proxy-controller.kubeconfig
