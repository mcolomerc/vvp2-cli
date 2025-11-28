# Quick Reference: Kubernetes Pod Options

## Basic Structure

```yaml
spec:
  template:
    spec:
      kubernetes:
        labels: {}
        pods: {}
        jobManagerPodTemplate: {}
        taskManagerPodTemplate: {}
```

## Common Use Cases

### 1. Mount Secrets for TLS/Kafka Configuration

```yaml
kubernetes:
  jobManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    spec:
      containers:
        - name: flink-main-container
          volumeMounts:
            - name: kafka-certs
              mountPath: /mnt/certs
              readOnly: true
      volumes:
        - name: kafka-certs
          secret:
            secretName: kafka-tls-secret
  taskManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    spec:
      containers:
        - name: flink-main-container
          volumeMounts:
            - name: kafka-certs
              mountPath: /mnt/certs
              readOnly: true
      volumes:
        - name: kafka-certs
          secret:
            secretName: kafka-tls-secret
```

### 2. Node Selection and Scheduling

```yaml
kubernetes:
  pods:
    nodeSelector:
      workload: flink
      zone: us-west-2a
    tolerations:
      - key: "dedicated"
        operator: "Equal"
        value: "flink"
        effect: "NoSchedule"
```

### 3. Environment Variables

```yaml
kubernetes:
  pods:
    envVars:
      # Direct value
      - name: ENV_NAME
        value: production
      # From secret
      - name: DB_PASSWORD
        valueFrom:
          secretKeyRef:
            name: db-credentials
            key: password
      # From ConfigMap
      - name: CONFIG_VALUE
        valueFrom:
          configMapKeyRef:
            name: app-config
            key: config-key
```

### 4. Resource Limits (in Pod Template)

```yaml
kubernetes:
  jobManagerPodTemplate:
    spec:
      containers:
        - name: flink-main-container
          resources:
            requests:
              cpu: "1"
              memory: "2Gi"
            limits:
              cpu: "2"
              memory: "4Gi"
```

### 5. Service Account

```yaml
kubernetes:
  jobManagerPodTemplate:
    spec:
      serviceAccountName: flink-service-account
  taskManagerPodTemplate:
    spec:
      serviceAccountName: flink-service-account
```

### 6. Labels and Annotations

```yaml
kubernetes:
  labels:
    team: data-engineering
    cost-center: analytics
  pods:
    labels:
      monitoring: enabled
    annotations:
      prometheus.io/scrape: "true"
      prometheus.io/port: "9249"
```

## Examples Directory

- `deployment-simple-k8s-templates.yaml` - Basic secret mounting
- `deployment-with-kubernetes-spec.yaml` - Full featured example
- `deployment-kafka-tls-k8s-spec.yaml` - Kafka TLS configuration

## CLI Commands

```bash
# Create deployment
./vvp2 deployment create -f deployment.yaml -n my-namespace

# Update deployment
./vvp2 deployment update my-deployment -f deployment.yaml -n my-namespace

# Get deployment (view configuration)
./vvp2 deployment get my-deployment -n my-namespace

# Delete deployment
./vvp2 deployment delete my-deployment -n my-namespace
```

## Important Notes

1. **Container Name**: Always use `flink-main-container` for the main Flink container
2. **Volume References**: VolumeMounts must reference volumes defined at pod spec level
3. **Both JM and TM**: Usually both JobManager and TaskManager need the same volumes/secrets
4. **API Version**: Use `apiVersion: v1` and `kind: Pod` for pod templates
5. **Backward Compatibility**: Old `flinkConfiguration` approach still works

## Migration from flinkConfiguration

**Old:**
```yaml
flinkConfiguration:
  kubernetes.pod-template-file.jobmanager: |
    apiVersion: v1
    kind: Pod
    spec:
      volumes:
        - name: secret
          secret:
            secretName: my-secret
```

**New (Recommended):**
```yaml
kubernetes:
  jobManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    spec:
      volumes:
        - name: secret
          secret:
            secretName: my-secret
```

## Full Documentation

See [KUBERNETES_POD_OPTIONS.md](KUBERNETES_POD_OPTIONS.md) for complete documentation.
