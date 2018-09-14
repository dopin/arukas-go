package arukas

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/dopin/jsonapi"
)

type Client struct {
	APIURL     *url.URL
	HTTP       *http.Client
	Username   string
	Password   string
	UserAgent  string
	Debug      bool
	Output     func(...interface{})
	OutputDest io.Writer
	Timeout    time.Duration
}

const VERSION = "v0.1.0"

func NewClientWithOsExitOnErr() *Client {
	client, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func NewClient() (*Client, error) {
	debug := false
	if os.Getenv("ARUKAS_DEBUG") != "" {
		debug = true
	}
	apiURL := "https://app.arukas.io/api/"
	if os.Getenv("ARUKAS_API_URL") != "" {
		apiURL = os.Getenv("ARUKAS_API_URL")
	}
	client := new(Client)
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")

	client.APIURL = parsedURL
	client.UserAgent = "Arukas Go Client (" + VERSION + ")"
	client.Debug = debug
	client.OutputDest = os.Stdout
	client.Timeout = 30 * time.Second

	if username := os.Getenv("ARUKAS_API_TOKEN"); username != "" {
		client.Username = username
	} else {
		return nil, errors.New("ARUKAS_API_TOKEN is not set")
	}

	if password := os.Getenv("ARUKAS_API_SECRET"); password != "" {
		client.Password = password
	} else {
		return nil, errors.New("ARUKAS_API_SECRET is not set")
	}

	return client, nil
}

func (client *Client) newRequest(method string, path string) (*http.Request, error) {
	var req *http.Request
	var err error

	requestURL := *client.APIURL
	requestURL.Path += path

	if req, err = http.NewRequest(method, requestURL.String(), nil); err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.api+json")
	req.Header.Set("User-Agent", client.UserAgent)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	req.SetBasicAuth(client.Username, client.Password)

	if client.Debug {
		for k, v := range req.Header {
			fmt.Fprint(os.Stderr, "[header] "+k)
			fmt.Fprintln(os.Stderr, ": "+strings.Join(v, ","))
		}
		fmt.Fprintln(os.Stderr, requestURL.String())
	}

	return req, nil
}

func (client *Client) request(req *http.Request) ([]byte, error) {
	httpClient := client.HTTP
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: client.Timeout,
		}
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if client.Debug {
		fmt.Fprintln(os.Stderr, "Status:", res.StatusCode)
		headers := make([]string, len(res.Header))
		for k := range res.Header {
			headers = append(headers, k)
		}
		sort.Strings(headers)
		for _, k := range headers {
			if k != "" {
				fmt.Fprintln(os.Stderr, k+":", strings.Join(res.Header[k], " "))
			}
		}
		fmt.Fprintln(os.Stderr, string(body))
	}
	return body, nil
}

func unmarshalManyPayload(responseBody []byte, resource reflect.Type) ([]interface{}, error) {
	return jsonapi.UnmarshalManyPayload(bytes.NewReader(responseBody), resource)
}

// GetApps returns a list of Arukas apps.
func (client *Client) GetApps() (Apps, error) {
	var err error
	req, err := client.newRequest("GET", endpointApps)

	if err != nil {
		return nil, err
	}

	responseBody, err := client.request(req)
	data, err := unmarshalManyPayload(responseBody, reflect.TypeOf(new(App)))

	if err != nil {
		return nil, err
	}

	var apps Apps

	for _, app := range data {
		if a, ok := app.(*App); ok {
			apps = append(apps, a)
		} else {
			return nil, fmt.Errorf("Failed to parse apps response")
		}
	}

	return apps, nil
}

// GetApp returns an Arukas app.
func (client *Client) GetApp(uuid string) (*App, error) {
	var err error
	req, err := client.newRequest("GET", fmt.Sprintf(endpointApp, uuid))

	if err != nil {
		return nil, err
	}

	responseBody, err := client.request(req)
	app := new(App)
	err = jsonapi.UnmarshalPayload(bytes.NewReader(responseBody), app)

	if err != nil {
		return nil, err
	}

	return app, nil
}

// GetService returns an Arukas service.
func (client *Client) GetService(uuid string) (*Service, error) {
	var err error
	req, err := client.newRequest("GET", fmt.Sprintf(endpointService, uuid))

	if err != nil {
		return nil, err
	}

	responseBody, err := client.request(req)
	service := new(Service)
	err = jsonapi.UnmarshalPayload(bytes.NewReader(responseBody), service)

	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetServices returns a list of Arukas services.
func (client *Client) GetServices() (Services, error) {
	var err error
	req, err := client.newRequest("GET", endpointServices)

	if err != nil {
		return nil, err
	}

	responseBody, err := client.request(req)
	data, err := unmarshalManyPayload(responseBody, reflect.TypeOf(new(Service)))

	if err != nil {
		return nil, err
	}

	var services Services

	for _, service := range data {
		if s, ok := service.(*Service); ok {
			services = append(services, s)
		} else {
			return nil, fmt.Errorf("Failed to parse services response")
		}
	}

	return services, nil
}
