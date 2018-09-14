package arukas

import "time"

type PortMapping struct {
	ContainerPort int    `json:"container_port"`
	ServicePort   int    `json:"service_port"`
	Host          string `json:"host"`
}

type TaskPorts []PortMapping

type PortMappings []TaskPorts

type Env struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Environment []*Env

type Port struct {
	Protocol string `json:"protocol"`
	Number   int    `json:"number"`
}

type Ports []*Port

type Service struct {
	ArukasDomain string       `jsonapi:"attr,arukas-domain"`
	Command      string       `jsonapi:"attr,command"`
	Endpoint     string       `jsonapi:"attr,endpoint"`
	Environment  Environment  `jsonapi:"attr,environment"`
	ID           string       `jsonapi:"primary,services"`
	Image        string       `jsonapi:"attr,image"`
	Instances    int          `jsonapi:"attr,instances"`
	Ports        Ports        `jsonapi:"attr,ports"`
	PortMappings PortMappings `jsonapi:"attr,port-mappings,omitempty"`
	ServicePlan  *ServicePlan `jsonapi:"relation,service-plan"`
	Subdomain    string       `jsonapi:"attr,subdomain"`
	Status       string       `jsonapi:"attr,status"`
	CreatedAt    time.Time    `jsonapi:"attr,created-at,iso8601"`
	UpdatedAt    time.Time    `jsonapi:"attr,updated-at,iso8601"`
}

type Services []*Service
