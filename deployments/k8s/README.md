# jobs-collector Helm chart

Production-oriented Helm chart for the jobs collector application and its optional in-cluster MySQL dependency.

## Install

```sh
helm upgrade --install jobs . --namespace job --create-namespace
```

## Common overrides

```sh
helm upgrade --install jobs . \
  --namespace job \
  --create-namespace \
  --set image.repository=registry.example.com/jobs-collector \
  --set image.tag=1.0.0 \
  --set mysql.auth.password='change-me' \
  --set mysql.auth.rootPassword='change-root-me' \
  --set ingress.hosts[0].host=jobs.example.com
```

## Notes

- Set `mysql.auth.existingSecret` to reuse an externally managed secret with `db-username`, `db-password`, and `db-root-password` keys.
- Application HTTP probes are disabled by default because the current binary initializes MySQL and exits instead of serving HTTP endpoints.
- Enable `probes.enabled` after the app exposes `GET /healthz` and `GET /readyz`.
- Enable `autoscaling.enabled` and `podDisruptionBudget.enabled` for multi-replica production deployments.
- Enable `migrations.enabled` only after providing a migration image command/args and any required migration volumes.
