package arukas

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/dopin/jsonapi"
)

const singleAppResponseJson = `
  {
	"data": {
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
			"id": "51d27e5b-533a-4603-aa09-450182a49193",
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
	}
  }
`

func TestUnmarshalSingleAppResponse(t *testing.T) {
	r := strings.NewReader(singleAppResponseJson)
	app := new(App)
	err := jsonapi.UnmarshalPayload(r, app)

	if err != nil {
		t.Fatal(err)
	}
}

func TestMarshalAppForCreation(t *testing.T) {
	app := &App{
		Name: "My New App",
		Services: Services{
			&Service{
				Image: "nginx",
				Ports: Ports{
					&Port{
						Protocol: "tcp",
						Number:   80,
					},
				},
				Environment: Environment{
					&Env{
						Key:   "FOO",
						Value: "bar",
					},
				},
				ServicePlan: &ServicePlan{
					ID: "jp-tokyo/hobby",
				},
			},
		},
	}

	out := bytes.NewBuffer(nil)
	var jsonData map[string]interface{}
	var ok bool

	err := jsonapi.MarshalPayload(out, app)

	if err != nil {
		t.Fatal(err)
	}

	if err = json.Unmarshal(out.Bytes(), &jsonData); err != nil {
		t.Fatal(err)
	}

	var attrs map[string]interface{}

	if attrs, ok = jsonData["data"].(map[string]interface{}); !ok {
		t.Fatalf("data key did not contain an Hash/Dict/Map")
	}

	var includedSlice []interface{}
	included := make([]map[string]interface{}, 0)

	if includedSlice, ok = jsonData["included"].([]interface{}); !ok {
		t.Fatalf("data.relationships key did not contain a slice")
	}

	for _, v := range includedSlice {
		if data, ok := v.(map[string]interface{}); !ok {
			t.Fatalf("data.relationships did not contain a slice of map")
		} else {
			included = append(included, data)
		}
	}

	if _, ok := attrs["id"]; ok {
		t.Error("ID must not be included.")
	}

	if attrs["type"] != typeApp {
		t.Fatalf("Wrong type. Expected %s, Got %s", typeApp, attrs["type"])
	}

	if attrs, ok = attrs["attributes"].(map[string]interface{}); !ok {
		t.Fatalf("data.attributes key did not contain an Hash/Dict/Map")
	}

	keys := []string{}

	for k := range attrs {
		keys = append(keys, k)
	}

	if len(keys) > 1 {
		t.Fatalf("data.attributes must have 1 key. Got %d", len(keys))
	}

	if attrs["name"] != "My New App" {
		t.Fatalf("data.attributes[\"name\"] was wrong. Expected \"My New App\", Got \"%s\"", attrs["name"])
	}

	var service map[string]interface{}
	var servicePlan map[string]interface{}

	for _, data := range included {
		if data["type"] == typeService {
			service = data
		} else if data["type"] == typeServicePlan {
			servicePlan = data
		} else {
			t.Fatalf("Unexpected relationship was included. %v", data)
		}
	}

	attrs, ok = service["attributes"].(map[string]interface{})

	if attrs["image"] != "nginx" {
		t.Errorf("Wrong image name. Expected nginx, Got %s.", attrs["image"])
	}

	ports := make([]string, 0)
	pValue := reflect.ValueOf(attrs["ports"])
	for i := 0; i < pValue.Len(); i++ {
		str, ok := pValue.Index(i).Interface().(string)
		if !ok {
			t.Fatalf("Ports is not a slice of string: Got %v", pValue.Index(i).Interface())
		}
		ports = append(ports, str)
	}

	if ports[0] != "80" {
		t.Errorf(`ports[0] must be "80". Got %s`, ports[0])
	}

	envs := make([]map[string]string, 0)
	envValue := reflect.ValueOf(attrs["environment"])

	if envValue.Len() != 1 {
		t.Errorf("environment length was wrong. Expected 1, Got %d", envValue.Len())
	}

	for i := 0; i < envValue.Len(); i++ {
		temp, ok := envValue.Index(i).Interface().(map[string]interface{})
		if !ok {
			t.Fatalf("Environment is not a slice of map[string]string: Got %v", envValue.Index(i))
		}

		env := make(map[string]string)

		for k, v := range temp {
			val, ok := v.(string)

			if !ok {
				t.Fatalf("One of environment is not a slice of map[string]string: Got %v", v)
			}

			env[k] = val
		}

		envs = append(envs, env)
	}

	if envs[0]["key"] != "FOO" || envs[0]["value"] != "bar" {
		t.Errorf(`Env[0] was wrong. Expected key = "FOO" and value = "bar", Got key = "%s", value = "%s"`, envs[0]["key"], envs[0]["value"])
	}

	if servicePlan["id"] != "jp-tokyo/hobby" {
		t.Errorf("relationship service-plans contains wrong id. Expected jp-tokyo/hobby, Got %s", servicePlan["id"])
	}

	if servicePlan["type"] != typeServicePlan {
		t.Errorf("relationship service-plans contains wrong type. Expected %s, Got %s", typeServicePlan, servicePlan["id"])
	}
}
