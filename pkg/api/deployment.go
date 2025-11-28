package api

import (
	"fmt"
	"time"
)

// Deployment represents a VVP deployment
type Deployment struct {
	Metadata DeploymentMetadata `json:"metadata" yaml:"metadata"`
	Spec     DeploymentSpec     `json:"spec" yaml:"spec"`
	Status   *DeploymentStatus  `json:"status,omitempty" yaml:"status,omitempty"`
}

// DeploymentMetadata holds deployment metadata
type DeploymentMetadata struct {
	ID          string            `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string            `json:"name" yaml:"name"`
	Namespace   string            `json:"namespace" yaml:"namespace"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	CreatedAt   time.Time         `json:"createdAt,omitempty" yaml:"createdAt,omitempty"`
	ModifiedAt  time.Time         `json:"modifiedAt,omitempty" yaml:"modifiedAt,omitempty"`
}

// DeploymentSpec holds deployment specification
type DeploymentSpec struct {
	State                string          `json:"state" yaml:"state"`
	UpgradeStrategy      UpgradeStrategy `json:"upgradeStrategy,omitempty" yaml:"upgradeStrategy,omitempty"`
	RestoreStrategy      RestoreStrategy `json:"restoreStrategy,omitempty" yaml:"restoreStrategy,omitempty"`
	DeploymentTargetName string          `json:"deploymentTargetName,omitempty" yaml:"deploymentTargetName,omitempty"`
	Template             Template        `json:"template" yaml:"template"`
	MaxSavepointAge      string          `json:"maxSavepointCreationTime,omitempty" yaml:"maxSavepointCreationTime,omitempty"`
	MaxJobCreationTime   string          `json:"maxJobCreationTime,omitempty" yaml:"maxJobCreationTime,omitempty"`
}

// DeploymentStatus holds deployment status
type DeploymentStatus struct {
	State   string      `json:"state" yaml:"state"`
	Running interface{} `json:"running,omitempty" yaml:"running,omitempty"`
}

// UpgradeStrategy defines upgrade strategy
type UpgradeStrategy struct {
	Kind string `json:"kind" yaml:"kind"`
}

// RestoreStrategy defines restore strategy
type RestoreStrategy struct {
	Kind                  string `json:"kind" yaml:"kind"`
	AllowNonRestoredState bool   `json:"allowNonRestoredState,omitempty" yaml:"allowNonRestoredState,omitempty"`
}

// (Removed old DeploymentTarget struct; API uses deploymentTargetName)

// Template defines the Flink job template
type Template struct {
	Spec TemplateSpec `json:"spec" yaml:"spec"`
}

// TemplateSpec holds template specification
type TemplateSpec struct {
	Artifact             Artifact          `json:"artifact" yaml:"artifact"`
	Parallelism          int               `json:"parallelism,omitempty" yaml:"parallelism,omitempty"`
	NumberOfTaskManagers int               `json:"numberOfTaskManagers,omitempty" yaml:"numberOfTaskManagers,omitempty"`
	Resources            Resources         `json:"resources,omitempty" yaml:"resources,omitempty"`
	FlinkVersion         string            `json:"flinkVersion,omitempty" yaml:"flinkVersion,omitempty"`
	FlinkImageTag        string            `json:"flinkImageTag,omitempty" yaml:"flinkImageTag,omitempty"`
	Logging              Logging           `json:"logging,omitempty" yaml:"logging,omitempty"`
	FlinkConfiguration   map[string]string `json:"flinkConfiguration,omitempty" yaml:"flinkConfiguration,omitempty"`
	Kubernetes           *KubernetesSpec   `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty"`
}

// Artifact represents a JAR artifact or SQL script
type Artifact struct {
	Kind          string `json:"kind" yaml:"kind"`
	JarURI        string `json:"jarUri,omitempty" yaml:"jarUri,omitempty"`
	MainClass     string `json:"mainClass,omitempty" yaml:"mainClass,omitempty"`
	EntryClass    string `json:"entryClass,omitempty" yaml:"entryClass,omitempty"`
	MainArgs      string `json:"mainArgs,omitempty" yaml:"mainArgs,omitempty"`
	FlinkVersion  string `json:"flinkVersion,omitempty" yaml:"flinkVersion,omitempty"`
	FlinkImageTag string `json:"flinkImageTag,omitempty" yaml:"flinkImageTag,omitempty"`
	SQLScript     string `json:"sqlScript,omitempty" yaml:"sqlScript,omitempty"`
}

// Resources defines resource requirements
type Resources struct {
	JobManager  ResourceSpec `json:"jobmanager,omitempty" yaml:"jobmanager,omitempty"`
	TaskManager ResourceSpec `json:"taskmanager,omitempty" yaml:"taskmanager,omitempty"`
}

// ResourceSpec defines CPU and memory resources
// CPU and Memory can be strings or numbers from the API, so we use interface{}
type ResourceSpec struct {
	CPU    interface{} `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory interface{} `json:"memory,omitempty" yaml:"memory,omitempty"`
}

// Logging defines logging configuration
type Logging struct {
	Log4jLoggers map[string]string `json:"log4jLoggers,omitempty" yaml:"log4jLoggers,omitempty"`
}

// KubernetesSpec defines Kubernetes-specific configuration
type KubernetesSpec struct {
	Labels                 map[string]string     `json:"labels,omitempty" yaml:"labels,omitempty"`
	Pods                   *KubernetesPodOptions `json:"pods,omitempty" yaml:"pods,omitempty"`
	JobManagerPodTemplate  *PodTemplateSpec      `json:"jobManagerPodTemplate,omitempty" yaml:"jobManagerPodTemplate,omitempty"`
	TaskManagerPodTemplate *PodTemplateSpec      `json:"taskManagerPodTemplate,omitempty" yaml:"taskManagerPodTemplate,omitempty"`
}

// KubernetesPodOptions defines pod-level options
type KubernetesPodOptions struct {
	Labels       map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations  map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	NodeSelector map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Tolerations  []Toleration      `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Affinity     *Affinity         `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	EnvVars      []EnvVar          `json:"envVars,omitempty" yaml:"envVars,omitempty"`
	Volumes      []Volume          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	VolumeMounts []VolumeMount     `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
}

// PodTemplateSpec represents a Kubernetes V1 PodTemplateSpec
type PodTemplateSpec struct {
	APIVersion string       `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind       string       `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata   *PodMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec       *PodSpec     `json:"spec,omitempty" yaml:"spec,omitempty"`
}

// PodMetadata represents pod metadata
type PodMetadata struct {
	Name        string            `json:"name,omitempty" yaml:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}

// PodSpec represents pod specification
type PodSpec struct {
	Containers         []Container       `json:"containers,omitempty" yaml:"containers,omitempty"`
	InitContainers     []Container       `json:"initContainers,omitempty" yaml:"initContainers,omitempty"`
	Volumes            []Volume          `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	NodeSelector       map[string]string `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	Affinity           *Affinity         `json:"affinity,omitempty" yaml:"affinity,omitempty"`
	Tolerations        []Toleration      `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	ServiceAccountName string            `json:"serviceAccountName,omitempty" yaml:"serviceAccountName,omitempty"`
	SecurityContext    interface{}       `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`
}

// Container represents a container specification
type Container struct {
	Name            string        `json:"name,omitempty" yaml:"name,omitempty"`
	Image           string        `json:"image,omitempty" yaml:"image,omitempty"`
	Command         []string      `json:"command,omitempty" yaml:"command,omitempty"`
	Args            []string      `json:"args,omitempty" yaml:"args,omitempty"`
	Env             []EnvVar      `json:"env,omitempty" yaml:"env,omitempty"`
	VolumeMounts    []VolumeMount `json:"volumeMounts,omitempty" yaml:"volumeMounts,omitempty"`
	Resources       interface{}   `json:"resources,omitempty" yaml:"resources,omitempty"`
	SecurityContext interface{}   `json:"securityContext,omitempty" yaml:"securityContext,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name      string        `json:"name" yaml:"name"`
	Value     string        `json:"value,omitempty" yaml:"value,omitempty"`
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

// EnvVarSource represents the source for an environment variable's value
type EnvVarSource struct {
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty" yaml:"configMapKeyRef,omitempty"`
	SecretKeyRef    *SecretKeySelector    `json:"secretKeyRef,omitempty" yaml:"secretKeyRef,omitempty"`
	FieldRef        *FieldSelector        `json:"fieldRef,omitempty" yaml:"fieldRef,omitempty"`
}

// ConfigMapKeySelector selects a key from a ConfigMap
type ConfigMapKeySelector struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

// SecretKeySelector selects a key from a Secret
type SecretKeySelector struct {
	Name string `json:"name" yaml:"name"`
	Key  string `json:"key" yaml:"key"`
}

// FieldSelector selects a field from the pod
type FieldSelector struct {
	FieldPath string `json:"fieldPath" yaml:"fieldPath"`
}

// Volume represents a volume that can be mounted
type Volume struct {
	Name                  string           `json:"name" yaml:"name"`
	Secret                *SecretVolume    `json:"secret,omitempty" yaml:"secret,omitempty"`
	ConfigMap             *ConfigMapVolume `json:"configMap,omitempty" yaml:"configMap,omitempty"`
	EmptyDir              interface{}      `json:"emptyDir,omitempty" yaml:"emptyDir,omitempty"`
	HostPath              *HostPathVolume  `json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
	PersistentVolumeClaim *PVCVolume       `json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
}

// SecretVolume represents a secret-backed volume
type SecretVolume struct {
	SecretName string `json:"secretName" yaml:"secretName"`
}

// ConfigMapVolume represents a configmap-backed volume
type ConfigMapVolume struct {
	Name string `json:"name" yaml:"name"`
}

// HostPathVolume represents a host path volume
type HostPathVolume struct {
	Path string `json:"path" yaml:"path"`
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
}

// PVCVolume represents a persistent volume claim
type PVCVolume struct {
	ClaimName string `json:"claimName" yaml:"claimName"`
}

// VolumeMount represents a volume mount
type VolumeMount struct {
	Name      string `json:"name" yaml:"name"`
	MountPath string `json:"mountPath" yaml:"mountPath"`
	ReadOnly  bool   `json:"readOnly,omitempty" yaml:"readOnly,omitempty"`
	SubPath   string `json:"subPath,omitempty" yaml:"subPath,omitempty"`
}

// Toleration represents a pod toleration
type Toleration struct {
	Key               string `json:"key,omitempty" yaml:"key,omitempty"`
	Operator          string `json:"operator,omitempty" yaml:"operator,omitempty"`
	Value             string `json:"value,omitempty" yaml:"value,omitempty"`
	Effect            string `json:"effect,omitempty" yaml:"effect,omitempty"`
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty" yaml:"tolerationSeconds,omitempty"`
}

// Affinity represents pod affinity/anti-affinity rules
type Affinity struct {
	NodeAffinity    interface{} `json:"nodeAffinity,omitempty" yaml:"nodeAffinity,omitempty"`
	PodAffinity     interface{} `json:"podAffinity,omitempty" yaml:"podAffinity,omitempty"`
	PodAntiAffinity interface{} `json:"podAntiAffinity,omitempty" yaml:"podAntiAffinity,omitempty"`
}

// DeploymentList represents a list of deployments
type DeploymentList struct {
	Items []DeploymentWithInfo `json:"items"`
}

// DeploymentWithInfo wraps deployment with operator info
type DeploymentWithInfo struct {
	Deployment Deployment `json:"deployment"`
}

// ListDeployments lists all deployments in a namespace
func (c *Client) ListDeployments(namespace string) (*DeploymentList, error) {
	var result DeploymentList
	resp, err := c.httpClient.R().
		SetResult(&result).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployments/with-cr", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDeployment gets a deployment by name
func (c *Client) GetDeployment(namespace, name string) (*Deployment, error) {
	var wrapper DeploymentWithInfo
	resp, err := c.httpClient.R().
		SetResult(&wrapper).
		Get(fmt.Sprintf("/api/v1/namespaces/%s/deployments/with-cr/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &wrapper.Deployment, nil
}

// CreateDeployment creates a new deployment
func (c *Client) CreateDeployment(namespace string, deployment *Deployment) (*Deployment, error) {
	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(deployment).
		SetResult(&result).
		Post(fmt.Sprintf("/api/v1/namespaces/%s/deployments", namespace))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateDeployment updates an existing deployment
func (c *Client) UpdateDeployment(namespace, name string, deployment *Deployment) (*Deployment, error) {
	// NOTE: remove debug logging; could be re-enabled with a verbose flag later
	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(deployment).
		SetResult(&result).
		Put(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteDeployment deletes a deployment
func (c *Client) DeleteDeployment(namespace, name string) error {
	resp, err := c.httpClient.R().
		Delete(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	return handleResponse(resp, err)
}

// UpdateDeploymentState updates the state of a deployment
func (c *Client) UpdateDeploymentState(namespace, name, state string) (*Deployment, error) {
	patch := map[string]interface{}{
		"spec": map[string]interface{}{
			"state": state,
		},
	}

	var result Deployment
	resp, err := c.httpClient.R().
		SetBody(patch).
		SetResult(&result).
		Patch(fmt.Sprintf("/api/v1/namespaces/%s/deployments/%s", namespace, name))

	if err := handleResponse(resp, err); err != nil {
		return nil, err
	}

	return &result, nil
}
