# kubecon-2025-ssd



export AKV_SUB_ID=9017db12-b4b5-460a-8edd-8961f6c863ec
export AKV_RG=kubecon-eu
export AKV_NAME=kubecon-eu-2025-akv
export CERT_SUBJECT=CN=jeremyrickard.com,O=kubecon,L=London,ST=London,C=UK

USER_ID=$(az ad signed-in-user show --query id -o tsv)
az role assignment create --role "Key Vault Secrets User" --role "Key Vault Certificates User" --role "Key Vault Crypto User" --assignee $USER_ID --scope "/subscriptions/$AKV_SUB_ID/resourceGroups/$AKV_RG/providers/Microsoft.KeyVault/vaults/$AKV_NAME"

CERT_NAME=kubecon-eu-cert

KEY_ID=https://kubecon-eu-2025-akv.vault.azure.net/keys/kubecon-eu/bece1cef5ec54479a56dc54b223f5a9f


Make a SP
https://learn.microsoft.com/en-us/entra/identity-platform/howto-create-service-principal-portal


Install Kyverno with Helm
helm install kyverno kyverno/kyverno -n kyverno --create-namespace \
  --set global.image.registry=zot.jeremyrickard.com \
  --set webhooksCleanup.image.repository=mirror/docker.io/bitnami/kubectl  \
  --set policyReportsCleanup.image.repository=mirror/docker.io/bitnami/kubectl \
  --set crd.image.repository=mirror/ghcr.io/kyverno/kyverno-cli \
  --set test.image.repository=mirror/docker.io/library/busybox \
  --set admissionController.initContainer.image.repository=mirror/ghcr.io/kyverno/kyvernopre \
  --set admissionController.container.image.repository=mirror/ghcr.io/kyverno/kyverno \
  --set cleanupController.image.repository=mirror/ghcr.io/kyverno/cleanup-controller \
  --set backgroundController.image.repository=mirror/ghcr.io/kyverno/background-controller \
  --set reportsController.image.repository=mirror/ghcr.io/kyverno/reports-controller


kubectl apply -f manifests/policies/restrict-images.yml
kubectl describe ClusterPolicy restrict-image-registries

kubectl apply -k manifests/flux
