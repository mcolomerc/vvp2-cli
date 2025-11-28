# Kubernetes Pod Options Support

The vvp2 CLI now supports the full Kubernetes pod specification as described in the [Ververica Platform documentation](https://docs.ververica.com/vvp/user-guide/application-operations/deployments/k8s-resources/#flink-pod-templates-recommended).

## Overview

You can configure Kubernetes-specific options for your Flink deployments using the `kubernetes` field in the deployment spec:

```yaml
spec:
  template:
    spec:
      kubernetes:
        labels: <Map<String, String>>
        pods: <KubernetesPodOptions>
        jobManagerPodTemplate: <V1PodTemplateSpec>
        taskManagerPodTemplate: <V1PodTemplateSpec>
```

## Configuration Options

### 1. Global Kubernetes Labels

Apply labels to all Flink pods (both JobManager and TaskManager):

```yaml
kubernetes:
  labels:
    team: data-engineering
    cost-center: analytics
```

### 2. Pod-Level Options (`pods`)

Configure common options that apply to both JobManager and TaskManager pods (unless overridden by specific pod templates):

```yaml
kubernetes:
  pods:
    labels:
      monitoring: enabled
    annotations:
      prometheus.io/scrape: "true"
    nodeSelector:
      workload-type: flink
    tolerations:
      - key: "flink-workload"
        operator: "Equal"
        value: "true"
        effect: "NoSchedule"
    envVars:
      - name: ENV_NAME
        value: production
    volumes:
      - name: kafka-config
        secret:
          secretName: kafka-config-properties
    volumeMounts:
      - name: kafka-config
        mountPath: /mnt/kafka-config
        readOnly: true
```

### 3. JobManager Pod Template

Define a specific pod template for the JobManager:

```yaml
kubernetes:
  jobManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    metadata:
      name: jobmanager-pod
      labels:
        component: jobmanager
    spec:
      serviceAccountName: flink-service-account
      containers:
        - name: flink-main-container
          resources:
            requests:
              cpu: "1"
              memory: "2Gi"
          volumeMounts:
            - name: kafka-config
              mountPath: /mnt/kafka-config
              readOnly: true
      volumes:
        - name: kafka-config
          secret:
            secretName: kafka-config-properties
```

### 4. TaskManager Pod Template

Define a specific pod template for the TaskManager:

```yaml
kubernetes:
  taskManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    metadata:
      name: taskmanager-pod
      labels:
        component: taskmanager
    spec:
      containers:
        - name: flink-main-container
          resources:
            requests:
              cpu: "2"
              memory: "4Gi"
          volumeMounts:
            - name: kafka-config
              mountPath: /mnt/kafka-config
              readOnly: true
      volumes:
        - name: kafka-config
          secret:
            secretName: kafka-config-properties
```

## Examples

### Basic Example with Secrets

```yaml
kind: Deployment
apiVersion: v1
metadata:
  name: simple-k8s-deployment
  namespace: default
spec:
  state: RUNNING
  deploymentTargetName: kubernetes-target
  template:
    spec:
      artifact:
        kind: JAR
        jarUri: "http://minio.vvp.svc:9000/maven/org/example/app/1.0.0/app-1.0.0.jar"
        flinkVersion: "1.20"
      parallelism: 1
      kubernetes:
        jobManagerPodTemplate:
          apiVersion: v1
          kind: Pod
          spec:
            containers:
              - name: flink-main-container
                volumeMounts:
                  - name: config-volume
                    mountPath: /etc/config
                    readOnly: true
            volumes:
              - name: config-volume
                secret:
                  secretName: my-config-secret
        taskManagerPodTemplate:
          apiVersion: v1
          kind: Pod
          spec:
            containers:
              - name: flink-main-container
                volumeMounts:
                  - name: config-volume
                    mountPath: /etc/config
                    readOnly: true
            volumes:
              - name: config-volume
                secret:
                  secretName: my-config-secret
```

See the following example files for more details:
- `examples/deployment-simple-k8s-templates.yaml` - Simple pod template configuration
- `examples/deployment-with-kubernetes-spec.yaml` - Comprehensive example with all options

## Migration from Old Format

### Old Format (Still Supported via flinkConfiguration)

```yaml
flinkConfiguration:
  kubernetes.pod-template-file.jobmanager: |
    apiVersion: v1
    kind: Pod
    metadata:
      name: jobmanager-pod-template
    spec:
      containers:
      - name: flink-main-container
        volumeMounts:
        - name: kafka-tls-certs
          mountPath: /mnt/sslcerts
    volumes:
    - name: kafka-tls-certs
      secret:
        secretName: kafka-generated-jks
```

### New Format (Recommended)

```yaml
kubernetes:
  jobManagerPodTemplate:
    apiVersion: v1
    kind: Pod
    metadata:
      name: jobmanager-pod
    spec:
      containers:
        - name: flink-main-container
          volumeMounts:
            - name: kafka-tls-certs
              mountPath: /mnt/sslcerts
      volumes:
        - name: kafka-tls-certs
          secret:
            secretName: kafka-generated-jks
```

## Supported Fields

The implementation supports the following Kubernetes resources:

- **Labels and Annotations**: Apply custom metadata to pods
- **Node Selection**: `nodeSelector`, `affinity`, `tolerations`
- **Volumes**: `secret`, `configMap`, `emptyDir`, `hostPath`, `persistentVolumeClaim`
- **Volume Mounts**: Mount volumes into containers
- **Environment Variables**: Direct values, ConfigMap refs, Secret refs, field refs
- **Service Accounts**: Assign service accounts to pods
- **Security Contexts**: Pod and container security contexts
- **Resource Limits**: CPU and memory requests/limits
- **Init Containers**: Run initialization containers before the main container

## Usage with CLI

Create a deployment with Kubernetes pod options:

```bash
./vvp2 deployment create -f examples/deployment-with-kubernetes-spec.yaml -n default
```

Update an existing deployment:

```bash
./vvp2 deployment update my-deployment -f examples/deployment-with-kubernetes-spec.yaml -n default
```

## Notes

- Pod templates are applied at the VVP API level and will be translated to Kubernetes resources by the Ververica Platform operator
- The `flink-main-container` is the standard container name that must be used for Flink JobManager and TaskManager containers
- Volume mounts in the main container should reference volumes defined at the pod spec level
- Both the new `kubernetes` spec format and the legacy `flinkConfiguration` approach are supported
