# Installation

## With Helm

1. Install the `bitflow` service account:
    ```
    kubectl create -f bitflow-controller/deploy/service_account.yaml
    ```
2. Optionally, create a file `helm/extra-values.yml` based on `helm/values.yml`, which overwrites relevant configuration values.
3. Install the Helm chart:
    ```
    helm install -f helm/extra-values.yml bitflow ./helm
    ```
