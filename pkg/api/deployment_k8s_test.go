package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestKubernetesSpecParsing(t *testing.T) {
	yamlContent := `
kind: Deployment
apiVersion: v1
metadata:
  name: test-deployment
  namespace: default
spec:
  state: RUNNING
  template:
    spec:
      artifact:
        kind: JAR
        jarUri: "http://example.com/app.jar"
        flinkVersion: "1.20"
      parallelism: 1
      kubernetes:
        labels:
          team: test
        pods:
          nodeSelector:
            node-type: worker
          tolerations:
            - key: test
              operator: Equal
              value: "true"
              effect: NoSchedule
          envVars:
            - name: ENV_VAR
              value: test-value
            - name: SECRET_VAR
              valueFrom:
                secretKeyRef:
                  name: my-secret
                  key: secret-key
        jobManagerPodTemplate:
          apiVersion: v1
          kind: Pod
          metadata:
            name: jm-pod
          spec:
            containers:
              - name: flink-main-container
                volumeMounts:
                  - name: config
                    mountPath: /etc/config
                    readOnly: true
            volumes:
              - name: config
                secret:
                  secretName: my-config
        taskManagerPodTemplate:
          apiVersion: v1
          kind: Pod
          spec:
            containers:
              - name: flink-main-container
                volumeMounts:
                  - name: data
                    mountPath: /mnt/data
            volumes:
              - name: data
                emptyDir: {}
`

	var deployment Deployment
	err := yaml.Unmarshal([]byte(yamlContent), &deployment)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test basic fields
	if deployment.Metadata.Name != "test-deployment" {
		t.Errorf("Expected name 'test-deployment', got '%s'", deployment.Metadata.Name)
	}

	// Test Kubernetes spec exists
	if deployment.Spec.Template.Spec.Kubernetes == nil {
		t.Fatal("Kubernetes spec should not be nil")
	}

	k8s := deployment.Spec.Template.Spec.Kubernetes

	// Test labels
	if k8s.Labels == nil {
		t.Fatal("Kubernetes labels should not be nil")
	}
	if k8s.Labels["team"] != "test" {
		t.Errorf("Expected label 'team' to be 'test', got '%s'", k8s.Labels["team"])
	}

	// Test pods options
	if k8s.Pods == nil {
		t.Fatal("Pods options should not be nil")
	}
	if k8s.Pods.NodeSelector["node-type"] != "worker" {
		t.Errorf("Expected nodeSelector 'node-type' to be 'worker'")
	}

	// Test tolerations
	if len(k8s.Pods.Tolerations) != 1 {
		t.Fatalf("Expected 1 toleration, got %d", len(k8s.Pods.Tolerations))
	}
	if k8s.Pods.Tolerations[0].Key != "test" {
		t.Errorf("Expected toleration key 'test', got '%s'", k8s.Pods.Tolerations[0].Key)
	}

	// Test environment variables
	if len(k8s.Pods.EnvVars) != 2 {
		t.Fatalf("Expected 2 env vars, got %d", len(k8s.Pods.EnvVars))
	}
	if k8s.Pods.EnvVars[0].Name != "ENV_VAR" {
		t.Errorf("Expected env var name 'ENV_VAR', got '%s'", k8s.Pods.EnvVars[0].Name)
	}
	if k8s.Pods.EnvVars[1].ValueFrom == nil {
		t.Error("Expected second env var to have valueFrom")
	}

	// Test JobManager pod template
	if k8s.JobManagerPodTemplate == nil {
		t.Fatal("JobManager pod template should not be nil")
	}
	if k8s.JobManagerPodTemplate.Metadata.Name != "jm-pod" {
		t.Errorf("Expected JM pod name 'jm-pod', got '%s'", k8s.JobManagerPodTemplate.Metadata.Name)
	}
	if len(k8s.JobManagerPodTemplate.Spec.Volumes) != 1 {
		t.Fatalf("Expected 1 volume in JM pod, got %d", len(k8s.JobManagerPodTemplate.Spec.Volumes))
	}

	// Test TaskManager pod template
	if k8s.TaskManagerPodTemplate == nil {
		t.Fatal("TaskManager pod template should not be nil")
	}
	if len(k8s.TaskManagerPodTemplate.Spec.Volumes) != 1 {
		t.Fatalf("Expected 1 volume in TM pod, got %d", len(k8s.TaskManagerPodTemplate.Spec.Volumes))
	}
	if k8s.TaskManagerPodTemplate.Spec.Volumes[0].EmptyDir == nil {
		t.Error("Expected TM volume to be emptyDir")
	}

	t.Log("All Kubernetes spec parsing tests passed!")
}

func TestSimpleKubernetesPodTemplate(t *testing.T) {
	yamlContent := `
kind: Deployment
apiVersion: v1
metadata:
  name: simple-deployment
spec:
  state: RUNNING
  template:
    spec:
      artifact:
        kind: JAR
        jarUri: "http://example.com/app.jar"
      kubernetes:
        jobManagerPodTemplate:
          apiVersion: v1
          kind: Pod
          spec:
            volumes:
              - name: secret-vol
                secret:
                  secretName: my-secret
`

	var deployment Deployment
	err := yaml.Unmarshal([]byte(yamlContent), &deployment)
	if err != nil {
		t.Fatalf("Failed to unmarshal simple YAML: %v", err)
	}

	if deployment.Spec.Template.Spec.Kubernetes == nil {
		t.Fatal("Kubernetes spec should not be nil")
	}

	if deployment.Spec.Template.Spec.Kubernetes.JobManagerPodTemplate == nil {
		t.Fatal("JobManager pod template should not be nil")
	}

	jmTemplate := deployment.Spec.Template.Spec.Kubernetes.JobManagerPodTemplate
	if len(jmTemplate.Spec.Volumes) != 1 {
		t.Fatalf("Expected 1 volume, got %d", len(jmTemplate.Spec.Volumes))
	}

	if jmTemplate.Spec.Volumes[0].Name != "secret-vol" {
		t.Errorf("Expected volume name 'secret-vol', got '%s'", jmTemplate.Spec.Volumes[0].Name)
	}

	if jmTemplate.Spec.Volumes[0].Secret == nil {
		t.Fatal("Expected secret volume")
	}

	if jmTemplate.Spec.Volumes[0].Secret.SecretName != "my-secret" {
		t.Errorf("Expected secret name 'my-secret', got '%s'", jmTemplate.Spec.Volumes[0].Secret.SecretName)
	}

	t.Log("Simple Kubernetes pod template test passed!")
}
