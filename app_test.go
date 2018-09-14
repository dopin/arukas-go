package arukas

import (
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
