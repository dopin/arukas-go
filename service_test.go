package arukas

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dopin/jsonapi"
)

const json = `{
  "data": [
    {
      "id": "126cd38b-371d-451a-8965-b84f1215ddfb",
      "type": "apps",
      "attributes": {
        "name": "redis",
        "created-at": "2018-03-29T09:20:01.704Z",
        "updated-at": "2018-03-29T09:20:01.704Z"
      },
      "relationships": {
        "user": {
          "data": {
            "id": "64f2383d-a84b-4370-8714-30eb6804cd8f",
            "type": "users"
          }
        },
        "services": {
          "data": [
            {
              "id": "a0154b23-d980-4b36-8c48-b22e350640d6",
              "type": "services"
            }
          ]
        }
      }
    },
    {
      "id": "b05892c0-1c88-4a67-88b6-c102e8c92f06",
      "type": "apps",
      "attributes": {
        "name": "nginx",
        "created-at": "2018-03-29T09:12:29.656Z",
        "updated-at": "2018-03-29T09:12:29.656Z"
      },
      "relationships": {
        "user": {
          "data": {
            "id": "64f2383d-a84b-4370-8714-30eb6804cd8f",
            "type": "users"
          }
        },
        "services": {
          "data": [
            {
              "id": "aec9b943-6901-4c4b-b9ed-57524bb8d286",
              "type": "services"
            }
          ]
        }
      }
    }
  ],
  "included": [
    {
      "id": "a0154b23-d980-4b36-8c48-b22e350640d6",
      "type": "services",
      "attributes": {
        "app-id": "126cd38b-371d-451a-8965-b84f1215ddfb",
        "image": "redis",
        "command": null,
        "instances": 1,
        "environment": [
          {
            "key": "aaaa",
            "value": "aaaaa"
          }
        ],
        "ports": [
          {
            "number": 6379,
            "protocol": "tcp"
          },
          {
            "number": 80,
            "protocol": "udp"
          }
        ],
        "port-mappings": null,
        "created-at": "2018-03-29T09:20:01.715Z",
        "updated-at": "2018-06-08T02:18:11.942Z",
        "status": "stopped",
        "subdomain": "goofy-bhaskara-1634",
        "endpoint": "goofy-bhaskara-1634.arukascloud.io",
        "custom-domain": null,
        "last-instance-failed-at": null,
        "last-instance-failed-status": null
      },
      "relationships": {
        "app": {
          "data": {
            "id": "126cd38b-371d-451a-8965-b84f1215ddfb",
            "type": "apps"
          }
        },
        "service-plan": {
          "data": {
            "id": "jp-tokyo/hobby",
            "type": "service-plans"
          }
        }
      }
    },
    {
      "id": "aec9b943-6901-4c4b-b9ed-57524bb8d286",
      "type": "services",
      "attributes": {
        "app-id": "b05892c0-1c88-4a67-88b6-c102e8c92f06",
        "image": "nginx:stable",
        "command": null,
        "instances": 1,
        "cpus": 0.1,
        "memory": null,
        "environment": null,
        "ports": [
          {
            "number": 80,
            "protocol": "tcp"
          },
          {
            "number": 443,
            "protocol": "tcp"
          }
        ],
        "port-mappings": null,
        "created-at": "2018-03-29T09:12:29.687Z",
        "updated-at": "2018-06-07T08:39:57.936Z",
        "status": "stopped",
        "subdomain": "mad-spence-8124",
        "endpoint": "mad-spence-8124.arukascloud.io",
        "custom-domain": null,
        "last-instance-failed-at": null,
        "last-instance-failed-status": null
      },
      "relationships": {
        "app": {
          "data": {
            "id": "b05892c0-1c88-4a67-88b6-c102e8c92f06",
            "type": "apps"
          }
        },
        "service-plan": {
          "data": {
            "id": "jp-tokyo/hobby",
            "type": "service-plans"
          }
        }
      }
    },
    {
        "id": "jp-tokyo/hobby",
        "type": "service-plans",
        "attributes": {
          "category": "hobby"
       }
    }
  ]
}
`

func TestUnmarshaServiceManyPayloadSideLoaded(t *testing.T) {
	r := strings.NewReader(json)
	data, err := jsonapi.UnmarshalManyPayload(r, reflect.TypeOf(new(App)))

	if err != nil {
		t.Fatal(err)
	}

	app1, ok := data[0].(*App)

	if !ok {
		t.Fatal("Type assertion failed at data[0]")
	}

	if app1.Name != "redis" {
		t.Fatalf("App name is wrong: expected redis, got %s", app1.Name)
	}

	time := time.Date(2018, 3, 29, 9, 20, 1, 704000000, time.UTC)

	if app1.CreatedAt.Sub(time) != 0 {
		t.Fatalf("App CreatedAt is wrong: expected %s, got %s", time, app1.CreatedAt)
	}

	if app1.UpdatedAt != time {
		t.Fatalf("App UreatedAt is wrong: expected %s, got %s", time, app1.UpdatedAt)
	}

	if app1.Services[0].Image != "redis" {
		t.Fatal("app1.Services.[0].Image is wrong.")
	}

	if app1.Services[0].Instances != 1 {
		t.Fatal("app1.Services.[0].Instances is wrong.")
	}

	if len(app1.Services[0].Environment) != 1 {
		t.Fatal("app1.Services[0].Environment length is wrong.")
	}

	if app1.Services[0].Environment[0].Key != "aaaa" ||
		app1.Services[0].Environment[0].Value != "aaaaa" {
		t.Fatal("app1.Services.[0].Environment[0] is wrong.")
	}

	if len(app1.Services[0].Ports) != 2 {
		t.Fatalf("app1.Services[0].Ports length is wrong. Expected 2, Got %d\n", len(app1.Services[0].Ports))
	}

	if app1.Services[0].Ports[0].Number != 6379 {
		t.Fatal("app1.Services[0].Ports[0].Number is wrong.")
	}

	if app1.Services[0].Ports[0].Protocol != "tcp" {
		t.Fatal("app1.Services[0].Ports[0].Protocol is wrong.")
	}

	if app1.Services[0].Status != "stopped" {
		t.Fatal("app1.Services[0].Status is wrong.")
	}

	if app1.Services[0].Subdomain != "goofy-bhaskara-1634" {
		t.Fatal("app1.Services[0].Subdomain is wrong.")
	}

	if app1.Services[0].Endpoint != "goofy-bhaskara-1634.arukascloud.io" {
		t.Fatal("app1.Services[0].Endpoint is wrong.")
	}

	if app1.Services[0].ServicePlan.Category != "hobby" {
		t.Fatalf("app1.Services[0].ServicePlan.Category is wrong. Expected: hobby, Got: %s\n", app1.Services[0].ServicePlan.Category)
	}
}

func TestUnmarshalServicePayload(t *testing.T) {
	body := `{
    "data": {
      "id": "a0154b23-d980-4b36-8c48-b22e350640d6",
      "type": "services",
      "attributes": {
        "app-id": "126cd38b-371d-451a-8965-b84f1215ddfb",
        "image": "redis",
        "command": null,
        "instances": 1,
        "environment": [
          {
            "key": "aaaa",
            "value": "aaaaa"
          }
        ],
        "ports": [
          {
            "number": 6379,
            "protocol": "tcp"
          },
          {
            "number": 80,
            "protocol": "udp"
          }
        ],
        "port-mappings": null,
        "created-at": "2018-03-29T09:20:01.715Z",
        "updated-at": "2018-06-08T02:18:11.942Z",
        "status": "stopped",
        "subdomain": "goofy-bhaskara-1634",
        "endpoint": "goofy-bhaskara-1634.arukascloud.io",
        "custom-domain": null,
        "last-instance-failed-at": null,
        "last-instance-failed-status": null
      },
      "relationships": {
        "app": {
          "data": {
            "id": "126cd38b-371d-451a-8965-b84f1215ddfb",
            "type": "apps"
          }
        },
        "service-plan": {
          "data": {
            "id": "jp-tokyo/hobby",
            "type": "service-plans"
          }
        }
      }
    }
  }`

	r := strings.NewReader(body)
	service := new(Service)
	err := jsonapi.UnmarshalPayload(r, service)

	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(service.Image)
	}

}
