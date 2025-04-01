# kubecon-2025-ssd

This repository has some notional demo capabilities for putting in place some supply chain security best practices. 

First, it provides a simple image retag capability, that integrates container signing with [notation](https://github.com/notaryproject/notation), and optionally [Project Copacetic](https://github.com/project-copacetic/copacetic) to patch the images. The signing will rely on the Azure Keyvault integration with Notation. 

Next, it provides example configuration for bootstrapping a Kubernetes cluster with Flux, specifically for using [Flux's OCI Artifact](https://fluxcd.io/flux/cheatsheets/oci-artifacts/) support. This also provides an example of pushing a Flux repository as an OCI Artifact, and then signing it. Flux will be configured in cluster to validate the signature. 

Finally, we leverage Kyverno to put policies in place for how the images are consumed. All images must come from our internal registry and must be signed. 

Once this is all in place, we are free to deploy a notional voting application using the Flux support above.

## Create a Kubernetes Cluster

You should be able to use kind for this!

## Install Zot

Zot can be installed on [bare metal](https://zotregistry.dev/v2.1.2/install-guides/install-guide-linux/) or on a [Kubernetes cluster](https://zotregistry.dev/v2.1.2/install-guides/install-guide-k8s/). Install Zot. You can reference the sample config in `zot/config.json`. This config enables anonymous read access and enables support for Docker images, in addition to OCI artifacts.

## Github Setup

First, you'll need to create an [environment](https://docs.github.com/en/actions/managing-workflow-runs-and-deployments/managing-deployments/managing-environments-for-deployment). This environment should have two secrets and two variables.

### Environment Secret

In order to access Azure Key Vault for signing, you will need to create a service principal with permisisions. [This] guide provides a good walkthrough. Once finished, create an environment secret named `AZURE_CREDENTIALS`. This should be of the form:

```
  {
      "clientId": "<Client ID>",
      "clientSecret": "<Client Secret>",
      "subscriptionId": "<Subscription ID>",
      "tenantId": "<Tenant ID>"
  }
```

You can find more information on creating a service principal [here](https://learn.microsoft.com/en-us/entra/identity-platform/howto-create-service-principal-portal).

The values should correlate to the service princiap created above.

NOTE: This is NOT a best practice. 

You also need to create a secret called `DOCKERHUB_TOKEN` that corresponds to the password of your Zot instance.

Next create two variables in the environment. One should be the username for your Zot server and should be called `DOCKERHUB_USERNAME`, while the other should be called `REGISTRY`, which is where your Zot server listens for requests. 

## Run Pipelines

Now you can run the `Retag Images` and `Publish Flux Config` pipelines. This should populate your zot cluster!

## Install Flux.

You should edit the `manifests/flux/kustomization.yml` file to replace `zot.jeremyrickard.com` with your registry of choice. Once you have done that, you can run `kubectl apply -k manifests/flux` to use Kustomize to install Flux.

## Install Kyverno

We will install Kyverno with Helm! 

helm install kyverno kyverno/kyverno -n kyverno --create-namespace \
  --set global.image.registry=<REPLACE-WITH-YOUR-REGISTRY> \
  --set webhooksCleanup.image.repository=mirror/docker.io/bitnami/kubectl  \
  --set policyReportsCleanup.image.repository=mirror/docker.io/bitnami/kubectl \
  --set crd.image.repository=mirror/ghcr.io/kyverno/kyverno-cli \
  --set test.image.repository=mirror/docker.io/library/busybox \
  --set admissionController.initContainer.image.repository=mirror/ghcr.io/kyverno/kyvernopre \
  --set admissionController.container.image.repository=mirror/ghcr.io/kyverno/kyverno \
  --set cleanupController.image.repository=mirror/ghcr.io/kyverno/cleanup-controller \
  --set backgroundController.image.repository=mirror/ghcr.io/kyverno/background-controller \
  --set reportsController.image.repository=mirror/ghcr.io/kyverno/reports-controller

for example:

```
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
```

## Apply cluster policies

kubectl apply -f manifests/policies/restrict-images.yml
kubectl apply -f manifests/policies/signed-images.yml

You will need to update both policies to reflect your environment. For `restrict-images.yml` you should update:

```
    spec:
          =(ephemeralContainers):
          - image: "zot.jeremyrickard.com/*"
          =(initContainers):
          - image: "zot.jeremyrickard.com/*"
          containers:
          - image: "zot.jeremyrickard.com/*"
```

For `signed-images.yml`, update the cert object with the corresponding cert from your AKV walkthrough above.

## Deploy your application

At this point, you will have a cluster with Flux and Kyverno installed, using the retagged images. Now we can deploy our voting app via Flux and OCI artifacts. 

`kubectl apply -f manifests/flux-repos/notation-config.yml` and `kubectl apply -f manifests/flux-repos/emoji-vote.yml`. notation-config.yml also needs to be updated to reflect the correct cert.

If everything is deployed correctly, you should be able to port-forward to the voting service.
