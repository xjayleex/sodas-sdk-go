package devops

import (
	"time"
)

type ListEnvironmentParams struct {
	Type       string `json:"type"`
	TemplateId string `json:"templateId"`
	VNodeId    string `json:"vnodeId"`
	State      string `json:"state"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Sort       string `json:"sort"`
	Ordered    string `json:"ordered"`
}

type ListEnvironmentResult struct {
	Total      int           `json:"total"`
	Aggregated []Aggregation `json:"aggregated"`
	Results    []Environment `json:"results"`
}

type GetEnvironmentParams struct {
	Id string `json:"id"`
}

type GetEnvironmentResult Environment

type Aggregation struct {
	State string `json:"state"`
	Count string `json:"count"`
}

type Environment struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Description  string      `json:"description"`
	TemplateID   string      `json:"templateId"`
	TemplateName string      `json:"templateName"`
	State        string      `json:"state"`
	SandboxID    string      `json:"sandboxId"`
	VNodeID      string      `json:"vnodeId"`
	Resources    Resources   `json:"resources"`
	Others       interface{} `json:"others"`
	Config       Config      `json:"config"`
	VNode        VNode       `json:"vnode"`
	IssuerID     string      `json:"issuerId"`
	Issued       time.Time   `json:"issued"`
	ModifierID   string      `json:"modifierId"`
	Modified     time.Time   `json:"modified"`
	Sandbox      Sandbox     `json:"sandbox"`
}

type Resources struct {
	CPU    float32 `json:"cpu"`
	Memory int     `json:"memory"`
	GPU    int     `json:"gpu"`
}

type Config struct {
	Labels       map[string]string `json:"labels"`
	Ports        []Port            `json:"ports"`
	Args         []string          `json:"args"`
	Command      []string          `json:"command"`
	Environments []EnvVariable     `json:"environments"`
	VolumeMounts []VolumeMount     `json:"volumeMounts"`
}

type Port struct {
	Description string `json:"description"`
	Port        int    `json:"port"`
	Protocol    string `json:"protocol"`
	Required    bool   `json:"required"`
	HTTPIngress bool   `json:"httpIngress"`
}

type VolumeMount struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type VNode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EnvVariable struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Secret      bool   `json:"secret"`
}

type Sandbox struct {
	ClusterID           string             `json:"clusterId"`
	VNodeID             string             `json:"vnodeId"`
	ID                  string             `json:"id"`
	Namespace           string             `json:"namespace"`
	UUID                string             `json:"uuid"`
	CreationTimestamp   time.Time          `json:"creationTimestamp"`
	Labels              map[string]string  `json:"labels"`
	Replicas            int                `json:"replicas"`
	AvailableReplicas   int                `json:"availableReplicas"`
	UnavailableReplicas int                `json:"unavailableReplicas"`
	RestartPolicy       string             `json:"restartPolicy"`
	ClusterIP           string             `json:"clusterIP"`
	Ports               []int              `json:"ports"`
	Pod                 []Pod              `json:"pod"`
	Conditions          []SandboxCondition `json:"conditions"`
	Ingress             []Ingress          `json:"ingress"`
}

type Pod struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SandboxCondition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastUpdateTime     time.Time `json:"lastUpdateTime"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Reason             string    `json:"reason"`
	Message            string    `json:"message"`
}

type Ingress struct {
	ClusterAccessPort int    `json:"clusterAccessPort"`
	Host              string `json:"host"`
	ServicePort       int    `json:"servicePort"`
}
