# Local Kubernetes Deployment With k3d and Helm

This guide helps developers run the `jobs-collector` Helm chart locally on a k3d cluster.

## Prerequisites

Install these tools:

- Docker
- k3d
- kubectl
- Helm

Verify them:

```sh
docker version
k3d version
kubectl version --client
helm version
```

Run all commands from this chart directory:

```sh
cd deployments/k8s
```

## 1. Create a k3d Cluster

Create a local cluster with HTTP traffic mapped to your machine:

```sh
k3d cluster create jobs-local \
  --api-port 6550 \
  -p "80:80@loadbalancer" \
  -p "443:443@loadbalancer" \
  --agents 1
```

Check the cluster:

```sh
kubectl cluster-info
kubectl get nodes
```

k3d includes Traefik by default, which matches the chart default `ingress.className=traefik`.

## 2. Build the Application Image

From the repository root, build the Docker image:

```sh
cd ../..
docker build -f deployments/docker/Dockerfile -t jobs-collector:dev .
```

Import the image into k3d:

```sh
k3d image import jobs-collector:dev -c jobs-local
```

Return to the chart directory:

```sh
cd deployments/k8s
```

## 3. Configure Local Hostname

The chart uses `job.localhost` by default.

Most systems already resolve `*.localhost` to `127.0.0.1`. If yours does not, add this line to `/etc/hosts`:

```txt
127.0.0.1 job.localhost
```

## 4. Install With Helm

Install or upgrade the chart:

```sh
kubectl config set-cluster k3d-jobs-local --server=https://127.0.0.1:6550

helm upgrade --install jobs . \
  --namespace job \
  --create-namespace \
  --set image.repository=jobs-collector \
  --set image.tag=dev \
  --set image.pullPolicy=IfNotPresent \
  --set mysql.auth.password=jobs \
  --set mysql.auth.rootPassword=root
```

Check deployed resources:

```sh
kubectl get all -n job
kubectl get ingress -n job
kubectl get pvc -n job
```

## 5. Verify Pods and Logs

Watch pods:

```sh
kubectl get pods -n job -w
```

View application logs:

```sh
kubectl logs -n job deploy/jobs-jobs-collector
```

View MySQL logs:

```sh
kubectl logs -n job statefulset/jobs-jobs-collector-mysql
```

Describe a failing pod:

```sh
kubectl describe pod -n job -l app.kubernetes.io/name=jobs-collector
```

## 6. Access the App

If the app exposes an HTTP server on `API_PORT`, open:

```sh
curl http://job.localhost
```

Current application note: the Go binary currently initializes MySQL and exits instead of serving HTTP endpoints. Until an HTTP server is added, the deployment may restart and ingress will not return an app response.

## 7. Connect to MySQL Locally

Port-forward MySQL:

```sh
kubectl port-forward -n job svc/jobs-jobs-collector-mysql 3306:3306
```

In another terminal:

```sh
mysql -h 127.0.0.1 -P 3306 -ujobs -pjobs jobs
```

If port `3306` is already used locally, forward another port:

```sh
kubectl port-forward -n job svc/jobs-jobs-collector-mysql 3307:3306
mysql -h 127.0.0.1 -P 3307 -ujobs -pjobs jobs
```

## 8. Render or Lint Before Deploying

Lint the chart:

```sh
helm lint .
```

Render manifests locally:

```sh
helm template jobs . --namespace job
```

Inspect effective values:

```sh
helm get values jobs -n job
helm get manifest jobs -n job
```

## 9. Upgrade After Code Changes

Rebuild and re-import the image:

```sh
cd ../..
docker build -f deployments/docker/Dockerfile -t jobs-collector:dev .
k3d image import jobs-collector:dev -c jobs-local
cd deployments/k8s
```

Restart the app deployment so Kubernetes uses the updated local image:

```sh
kubectl rollout restart deployment/jobs-jobs-collector -n job
kubectl rollout status deployment/jobs-jobs-collector -n job
```

If chart values changed, run Helm again:

```sh
helm upgrade --install jobs . \
  --namespace job \
  --create-namespace \
  --set migrations.enabled=true \
  --set migrations.image.repository=linkedin-job \
  --set migrations.image.tag=1.0.0 \
  --set migrations.image.pullPolicy=IfNotPresent \
  --set mysql.auth.password=jobs \
  --set mysql.auth.rootPassword=root
```

## 10. Useful Helm Values for Local Development

Disable ingress:

```sh
helm upgrade --install jobs . -n job --set ingress.enabled=false
```

Disable persistent MySQL storage:

```sh
helm upgrade --install jobs . -n job --set mysql.persistence.enabled=false
```

Use one app replica:

```sh
helm upgrade --install jobs . -n job --set replicaCount=1
```

Use a custom local hostname:

```sh
helm upgrade --install jobs . -n job --set ingress.hosts[0].host=jobs.localhost
```

## 11. Troubleshooting

Check Helm release status:

```sh
helm status jobs -n job
```

Check Kubernetes events:

```sh
kubectl get events -n job --sort-by=.lastTimestamp
```

Check Traefik:

```sh
kubectl get pods -n kube-system -l app.kubernetes.io/name=traefik
kubectl logs -n kube-system -l app.kubernetes.io/name=traefik
```

Common issues:

- `ImagePullBackOff`: rebuild the image and run `k3d image import jobs-collector:dev -c jobs-local`.
- `CrashLoopBackOff`: check app logs; the current binary exits after DB initialization.
- Ingress not responding: verify `kubectl get ingress -n job` and confirm `job.localhost` resolves to `127.0.0.1`.
- MySQL password mismatch: uninstall the release and delete the PVC if you intentionally want a fresh database.

## 12. Cleanup

Uninstall the Helm release:

```sh
helm uninstall jobs -n job
```

Delete the namespace:

```sh
kubectl delete namespace job
```

Delete the k3d cluster:

```sh
k3d cluster delete jobs-local
```
