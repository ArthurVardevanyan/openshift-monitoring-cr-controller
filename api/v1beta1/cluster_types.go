/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	EnableUserWorkload    bool                  `json:"enableUserWorkload,omitempty"`
	PrometheusOperator    PrometheusOperator    `json:"prometheusOperator,omitempty"`
	PrometheusK8S         PrometheusK8S         `json:"prometheusK8s,omitempty"`
	AlertmanagerMain      AlertmanagerMain      `json:"alertmanagerMain,omitempty"`
	KubeStateMetrics      KubeStateMetrics      `json:"kubeStateMetrics,omitempty"`
	MonitoringPlugin      MonitoringPlugin      `json:"monitoringPlugin,omitempty"`
	OpenshiftStateMetrics OpenshiftStateMetrics `json:"openshiftStateMetrics,omitempty"`
	TelemeterClient       TelemeterClient       `json:"telemeterClient,omitempty"`
	MetricsServer         MetricsServer         `json:"metricsServer,omitempty"`
	ThanosQuerier         ThanosQuerier         `json:"thanosQuerier,omitempty"`
}

type Metadata struct {
}

type Status struct {
}

type AlertmanagerMain struct {
	LogLevel                     string                                `json:"logLevel,omitempty"`
	EnableUserAlertmanagerConfig bool                                  `json:"enableUserAlertmanagerConfig,omitempty"`
	NodeSelector                 map[string]string                     `json:"nodeSelector,omitempty"`
	Resources                    *corev1.ResourceRequirements          `json:"resources,omitempty"`
	Tolerations                  []corev1.Toleration                   `json:"tolerations,omitempty"`
	VolumeClaimTemplate          *corev1.PersistentVolumeClaimTemplate `json:"volumeClaimTemplate,omitempty"`
}
type MonitoringPlugin struct {
	LogLevel     string              `json:"logLevel,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}
type MetricsServer struct {
	LogLevel     string              `json:"logLevel,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}
type KubeStateMetrics struct {
	LogLevel     string              `json:"logLevel,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}
type OpenshiftStateMetrics struct {
	LogLevel     string              `json:"logLevel,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}
type BearerToken struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
type Ca struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}
type TLSConfig struct {
	ServerName         string `json:"ServerName,omitempty"`
	Ca                 Ca     `json:"ca,omitempty"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify,omitempty"`
}
type AdditionalAlertManagerConfigs struct {
	APIVersion    string      `json:"apiVersion,omitempty"`
	BearerToken   BearerToken `json:"bearerToken,omitempty"`
	PathPrefix    string      `json:"pathPrefix,omitempty"`
	Scheme        string      `json:"scheme,omitempty"`
	StaticConfigs []string    `json:"staticConfigs,omitempty"`
	TLSConfig     TLSConfig   `json:"tlsConfig,omitempty"`
}

type PrometheusK8S struct {
	AdditionalAlertManagerConfigs []AdditionalAlertManagerConfigs       `json:"additionalAlertmanagerConfigs,omitempty"`
	ExternalLabels                map[string]string                     `json:"externalLabels,omitempty"`
	LogLevel                      string                                `json:"logLevel,omitempty"`
	NodeSelector                  map[string]string                     `json:"nodeSelector,omitempty"`
	Resources                     *corev1.ResourceRequirements          `json:"resources,omitempty"`
	Retention                     string                                `json:"retention,omitempty"`
	Tolerations                   []corev1.Toleration                   `json:"tolerations,omitempty"`
	VolumeClaimTemplate           *corev1.PersistentVolumeClaimTemplate `json:"volumeClaimTemplate,omitempty"`
}
type PrometheusOperator struct {
	LogLevel     string              `json:"logLevel,omitempty"`
	NodeSelector map[string]string   `json:"nodeSelector,omitempty"`
	Tolerations  []corev1.Toleration `json:"tolerations,omitempty"`
}
type TelemeterClient struct {
	LogLevel           string              `json:"logLevel,omitempty"`
	ClusterID          string              `json:"clusterID,omitempty"`
	NodeSelector       map[string]string   `json:"nodeSelector,omitempty"`
	TelemeterServerURL string              `json:"telemeterServerURL,omitempty"`
	Token              string              `json:"token,omitempty"`
	Tolerations        []corev1.Toleration `json:"tolerations,omitempty"`
}
type ThanosQuerier struct {
	LogLevel     string                       `json:"logLevel,omitempty"`
	NodeSelector map[string]string            `json:"nodeSelector,omitempty"`
	Resources    *corev1.ResourceRequirements `json:"resources,omitempty"`
	Tolerations  []corev1.Toleration          `json:"tolerations,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
