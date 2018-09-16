package arukas

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

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

func (port *Port) MarshalJSON() ([]byte, error) {
	str := strconv.Itoa(port.Number)

	if port.Protocol == "udp" {
		str += "/udp"
	}

	return []byte(`"` + str + `"`), nil
}

type Ports []*Port

func (ports *Ports) UnmarshalJSON(data []byte) error {
	np := new(APIPortFormat)
	var err error
	if err = json.Unmarshal(data, np); err != nil {
		return err
	}
	if *ports, err = np.toPorts(); err != nil {
		return err
	}
	return nil
}

func parsePortFormat(str string) (*Port, error) {
	var protocol string
	var parsedInt int64
	var number int
	var err error
	splitted := strings.Split(str, "/")
	if len(splitted) <= 1 {
		protocol = "tcp"
	} else {
		protocol = splitted[1]
	}
	if parsedInt, err = strconv.ParseInt(splitted[0], 10, 32); err != nil {
		return nil, err
	}
	number = int(parsedInt)
	return &Port{Protocol: protocol, Number: number}, nil
}

type APIPortFormat []string

func (pf APIPortFormat) toPorts() (Ports, error) {
	ports := make(Ports, 0)
	for _, p := range pf {
		var (
			parsedPort *Port
			err        error
		)
		if parsedPort, err = parsePortFormat(p); err != nil {
			return nil, err
		}
		ports = append(ports, parsedPort)
	}
	return ports, nil
}

type Service struct {
	ClientID     string       `jsonapi:"lid"`
	Command      string       `jsonapi:"attr,command"`
	Endpoint     string       `jsonapi:"attr,endpoint,readonly"`
	Environment  Environment  `jsonapi:"attr,environment"`
	ID           string       `jsonapi:"primary,services"`
	Image        string       `jsonapi:"attr,image"`
	Instances    int          `jsonapi:"attr,instances"`
	Ports        Ports        `jsonapi:"attr,ports"`
	PortMappings PortMappings `jsonapi:"attr,port-mappings,readonly"`
	ServicePlan  *ServicePlan `jsonapi:"relation,service-plan"`
	Subdomain    string       `jsonapi:"attr,subdomain"`
	Status       string       `jsonapi:"attr,status,readonly"`
	CreatedAt    time.Time    `jsonapi:"attr,created-at,iso8601,readonly"`
	UpdatedAt    time.Time    `jsonapi:"attr,updated-at,iso8601,readonly"`
}

type Services []*Service
