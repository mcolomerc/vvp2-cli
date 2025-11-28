package api

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSessionClusterParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SessionCluster
metadata:
  name: test-session-cluster
  namespace: default
  labels:
    env: test
spec:
  deploymentTargetName: kubernetes-target
  flinkVersion: "1.20"
  flinkImageRegistry: registry.ververica.com/v2.15
  flinkImageRepository: flink
  flinkImageTag: 1.20.1-stream1-scala_2.12-java11
  numberOfTaskManagers: 2
  state: RUNNING
  resources:
    jobmanager:
      cpu: 1.0
      memory: 1024m
    taskmanager:
      cpu: 2.0
      memory: 2048m
  flinkConfiguration:
    taskmanager.numberOfTaskSlots: "2"
    execution.checkpointing.interval: "60s"
`

	var sessionCluster SessionCluster
	err := yaml.Unmarshal([]byte(yamlContent), &sessionCluster)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Test metadata
	if sessionCluster.Metadata.Name != "test-session-cluster" {
		t.Errorf("Expected name 'test-session-cluster', got '%s'", sessionCluster.Metadata.Name)
	}
	if sessionCluster.Metadata.Namespace != "default" {
		t.Errorf("Expected namespace 'default', got '%s'", sessionCluster.Metadata.Namespace)
	}
	if sessionCluster.Metadata.Labels["env"] != "test" {
		t.Errorf("Expected label 'env' to be 'test'")
	}

	// Test spec
	if sessionCluster.Spec.DeploymentTargetName != "kubernetes-target" {
		t.Errorf("Expected deploymentTargetName 'kubernetes-target', got '%s'", sessionCluster.Spec.DeploymentTargetName)
	}
	if sessionCluster.Spec.FlinkVersion != "1.20" {
		t.Errorf("Expected flinkVersion '1.20', got '%s'", sessionCluster.Spec.FlinkVersion)
	}
	if sessionCluster.Spec.NumberOfTaskManagers != 2 {
		t.Errorf("Expected numberOfTaskManagers 2, got %d", sessionCluster.Spec.NumberOfTaskManagers)
	}
	if sessionCluster.Spec.State != "RUNNING" {
		t.Errorf("Expected state 'RUNNING', got '%s'", sessionCluster.Spec.State)
	}

	// Test resources
	if len(sessionCluster.Spec.Resources) != 2 {
		t.Fatalf("Expected 2 resources, got %d", len(sessionCluster.Spec.Resources))
	}

	// Test flink configuration
	if sessionCluster.Spec.FlinkConfiguration["taskmanager.numberOfTaskSlots"] != "2" {
		t.Errorf("Expected taskSlots '2', got '%s'", sessionCluster.Spec.FlinkConfiguration["taskmanager.numberOfTaskSlots"])
	}

	t.Log("SessionCluster parsing test passed!")
}

func TestSessionClusterListParsing(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SessionClusterList
items:
  - apiVersion: v1
    kind: SessionCluster
    metadata:
      name: cluster1
      namespace: default
    spec:
      deploymentTargetName: target1
      flinkVersion: "1.20"
      flinkImageRegistry: registry.ververica.com/v2.15
      flinkImageRepository: flink
      flinkImageTag: 1.20.1-stream1-scala_2.12-java11
      numberOfTaskManagers: 1
      state: RUNNING
      resources:
        jobmanager:
          cpu: 1
          memory: 1024m
        taskmanager:
          cpu: 1
          memory: 1024m
  - apiVersion: v1
    kind: SessionCluster
    metadata:
      name: cluster2
      namespace: default
    spec:
      deploymentTargetName: target2
      flinkVersion: "1.19"
      flinkImageRegistry: registry.ververica.com/v2.15
      flinkImageRepository: flink
      flinkImageTag: 1.19.1-stream1-scala_2.12-java11
      numberOfTaskManagers: 2
      state: STOPPED
      resources:
        jobmanager:
          cpu: 1
          memory: 1024m
        taskmanager:
          cpu: 2
          memory: 2048m
`

	var sessionClusterList SessionClusterList
	err := yaml.Unmarshal([]byte(yamlContent), &sessionClusterList)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if len(sessionClusterList.Items) != 2 {
		t.Fatalf("Expected 2 session clusters, got %d", len(sessionClusterList.Items))
	}

	// Test first cluster
	if sessionClusterList.Items[0].Metadata.Name != "cluster1" {
		t.Errorf("Expected first cluster name 'cluster1', got '%s'", sessionClusterList.Items[0].Metadata.Name)
	}
	if sessionClusterList.Items[0].Spec.State != "RUNNING" {
		t.Errorf("Expected first cluster state 'RUNNING', got '%s'", sessionClusterList.Items[0].Spec.State)
	}

	// Test second cluster
	if sessionClusterList.Items[1].Metadata.Name != "cluster2" {
		t.Errorf("Expected second cluster name 'cluster2', got '%s'", sessionClusterList.Items[1].Metadata.Name)
	}
	if sessionClusterList.Items[1].Spec.State != "STOPPED" {
		t.Errorf("Expected second cluster state 'STOPPED', got '%s'", sessionClusterList.Items[1].Spec.State)
	}

	t.Log("SessionClusterList parsing test passed!")
}

func TestSessionClusterWithKubernetes(t *testing.T) {
	yamlContent := `
apiVersion: v1
kind: SessionCluster
metadata:
  name: session-with-k8s
  namespace: default
spec:
  deploymentTargetName: kubernetes-target
  flinkVersion: "1.20"
  flinkImageRegistry: registry.ververica.com/v2.15
  flinkImageRepository: flink
  flinkImageTag: 1.20.1-stream1-scala_2.12-java11
  numberOfTaskManagers: 1
  state: RUNNING
  resources:
    jobmanager:
      cpu: 1
      memory: 1024m
    taskmanager:
      cpu: 1
      memory: 1024m
  kubernetes:
    jobManagerPodTemplate:
      apiVersion: v1
      kind: Pod
      spec:
        containers:
          - name: flink-main-container
            volumeMounts:
              - name: config
                mountPath: /etc/config
        volumes:
          - name: config
            secret:
              secretName: my-config
`

	var sessionCluster SessionCluster
	err := yaml.Unmarshal([]byte(yamlContent), &sessionCluster)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	if sessionCluster.Spec.Kubernetes == nil {
		t.Fatal("Kubernetes spec should not be nil")
	}

	t.Log("SessionCluster with Kubernetes spec parsing test passed!")
}
